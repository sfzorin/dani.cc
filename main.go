// Command morning is a mobile-first PWA for a 30-day morning-exercise program,
// built on the doors framework. See internal/* for content, i18n, db and auth.
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"morning/assets"
	"morning/internal/app"
	"morning/internal/auth"
	"morning/internal/db"
	"morning/internal/i18n"

	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

func main() {
	dbPath := os.Getenv("MORNING_DB")
	if dbPath == "" {
		dbPath = "morning.db"
	}
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()
	app.DB = database

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

			// No valid session → create an anonymous guest, identified by a
			// cookie, so progress is saved without registration.
			if tok, err := auth.NewToken(); err == nil {
				if emailTok, err := auth.NewToken(); err == nil {
					if gu, err := database.CreateGuest("guest:"+emailTok, lang); err == nil {
						if err := database.CreateSession(tok, gu.ID, auth.SessionTTL); err == nil {
							r.SetCookie(&http.Cookie{
								Name: "session", Value: tok, Path: "/",
								HttpOnly: true, SameSite: http.SameSiteLaxMode,
								MaxAge: int(auth.SessionTTL.Seconds()),
							})
							doors.SessionExpire(ctx, auth.SessionTTL)
							return doors.NewSource(sessionFromUser(gu))
						}
					}
				}
			}

			// Fallback (DB error): unauthenticated shell.
			return doors.NewSource(auth.Session{Lang: lang, Rest: 20})
		}).(doors.Source[auth.Session])
		return App{auth: src}
	},
		// Safari on http://localhost needs a non-secure session cookie.
		doors.WithConf(doors.Conf{ServerSessionCookieNoSecure: true}),
	)

	// Serve embedded static files (icons, manifest, service worker, exercise SVGs).
	a.Use(doors.UseFS("/static", assets.Static(), doors.CacheControlStatic))

	addr := ":8080"
	log.Printf("morning listening on %s", addr)
	if err := http.ListenAndServe(addr, a); err != nil {
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
		IsGuest:    u.IsGuest,
	}
}
