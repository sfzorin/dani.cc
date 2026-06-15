// Command morning is a mobile-first PWA for a 30-day morning-exercise program,
// built on the doors framework. See internal/* for content, i18n, db and auth.
package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"

	"morning30/assets"
	"morning30/internal/app"
	"morning30/internal/auth"
	"morning30/internal/db"
	"morning30/internal/i18n"

	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

// bodyCore is the shared figure-math ES module: joints()/constraints (used by the 3D
// editor's mini-projections) + the legacy SVG figure. Served at /body/core.mjs.
//
//go:embed body/core.mjs
var bodyCore []byte

// bodyR3D is the Three.js + Mixamo-rig 3D mannequin editor, served at /body (and /body3d).
//
//go:embed body/web/r3d.html
var bodyR3D []byte

// bodyShot is the headless single-pose figure renderer, served at /body/shot.
//
//go:embed body/web/shot.html
var bodyShot []byte

// bodyPoses is the per-exercise pose table used by the shot renderer.
//
//go:embed body/poses.mjs
var bodyPoses []byte

// bodyShotGrid renders every pose to a contact sheet (review tool), served at /body/shotgrid.
//
//go:embed body/web/shotgrid.html
var bodyShotGrid []byte

func main() {
	dbPath := os.Getenv("DANI_DB")
	if dbPath == "" {
		dbPath = "dani.db"
	}
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()
	app.DB = database
	app.SeedBuiltins() // populate the global exercise library (idempotent)

	a := doors.NewApp(func(ctx context.Context, r doors.Request) gox.Comp {
		// Bootstrap the shared reactive session once per browser session.
		src := doors.SessionStore(ctx).Init(auth.StoreKey{}, func() any {
			lang := string(i18n.FromAcceptLanguage(r.RequestHeader().Get("Accept-Language")))

			// Existing (real or guest) session via cookie.
			if c, err := r.GetCookie("session"); err == nil {
				if u, ok, _ := database.SessionUser(c.Value); ok {
					return doors.NewSource(sessionFromUser(u))
				}
			}

			// No valid session → land on the sign-in screen (unauthenticated).
			// A guest session is created on demand only when the visitor picks
			// "continue as guest" on the auth page.
			return doors.NewSource(auth.Session{Lang: lang, Rest: 20})
		}).(doors.Source[auth.Session])
		return App{auth: src}
	},
		// Safari on http://localhost needs a non-secure session cookie.
		doors.WithConf(doors.Conf{ServerSessionCookieNoSecure: true}),
	)

	// Serve embedded static files (icons, manifest, service worker, exercise SVGs).
	a.Use(doors.UseFS("/static", assets.Static(), doors.CacheControlStatic))

	// /body — standalone interactive 3D mannequin editor (separate from the workout app).
	mux := http.NewServeMux()
	mux.HandleFunc("/body", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(bodyR3D)
	})
	mux.HandleFunc("/body/core.mjs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(bodyCore)
	})
	mux.HandleFunc("/body3d", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(bodyR3D)
	})
	mux.HandleFunc("/body/shot", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(bodyShot)
	})
	mux.HandleFunc("/body/poses.mjs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(bodyPoses)
	})
	mux.HandleFunc("/body/shotgrid", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(bodyShotGrid)
	})
	// /body/models/ — drop-in GLB mannequins (e.g. m.glb / f.glb exported from Mixamo),
	// served from disk so swapping a model needs no rebuild. The 3D editor loads these
	// first and falls back to CDN stand-ins if absent.
	mux.Handle("/body/models/", http.StripPrefix("/body/models/", http.FileServer(http.Dir("body/web/models"))))
	mux.Handle("/", a)

	addr := ":8080"
	log.Printf("morning30 listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// sessionFromUser builds the reactive session value from a DB user.
func sessionFromUser(u db.User) auth.Session {
	return auth.Session{
		Authorized: true,
		UserID:     u.ID,
		Email:      u.Email,
		Name:       u.Name,
		Lang:       u.Lang,
		Rest:       u.Rest,
		Voice:      u.Voice,
		VoiceMode:  u.VoiceMode,
		IsGuest:    u.IsGuest,
		Level:      u.Level,
	}
}
