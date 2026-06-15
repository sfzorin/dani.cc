// Managed by GoX v0.1.36

//line common.gox:1
package pages

import (
	"context"
	"net/http"

	"morning30/internal/app"
	"morning30/internal/auth"
	"morning30/internal/i18n"
	"morning30/path"

	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

// langSwitch is a compact language selector. It updates the shared session and,
// when persist is set, saves the choice to the user's account.
type langSwitch struct {
	sess    auth.Session
	auth    doors.Source[auth.Session]
	persist bool
}

//line common.gox:24
func (ls langSwitch) Main() gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line common.gox:26
		cur := i18n.Lang(ls.sess.Lang)

		__e = __c.Init("select"); if __e != nil { return }
		{
//line common.gox:28
			__e = __c.Set("class", "lang-select"); if __e != nil { return }
//line common.gox:28
			__e = __c.Modify(doors.AChange{On: ls.onChange}); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
//line common.gox:29
			for _, m := range i18n.Languages {
				__e = __c.Init("option"); if __e != nil { return }
				{
//line common.gox:30
					__e = __c.Set("value", string(m.Code)); if __e != nil { return }
//line common.gox:30
					__e = __c.Set("selected", func() any { return m.Code == cur }()); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line common.gox:31
					__e = __c.Any(m.Name); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line common.gox:35
}

func (ls langSwitch) onChange(ctx context.Context, r doors.RequestEvent[doors.ChangeEvent]) bool {
	code := r.Event().Value
	if !i18n.Valid(code) {
		return false
	}
	ls.auth.Mutate(ctx, func(s auth.Session) auth.Session {
		s.Lang = code
		return s
	})
	if ls.persist && ls.sess.UserID != 0 {
		_ = app.DB.SetLang(ls.sess.UserID, code)
	}
	return false
}

// header is the top bar for in-app pages: streak, settings + logout, and a
// left-side title. When title is set it replaces the brand (e.g. the home
// greeting), so the greeting is the very first line — no extra brand row.
type header struct {
	sess   auth.Session
	auth   doors.Source[auth.Session]
	path   doors.Source[path.Path]
	streak int
	title  string
}

//line common.gox:63
func (h header) Main() gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line common.gox:65
		l := i18n.Lang(h.sess.Lang)

		__e = __c.Init("header"); if __e != nil { return }
		{
//line common.gox:67
			__e = __c.Set("class", "top"); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
//line common.gox:68
			if h.title != "" {
				__e = __c.Init("span"); if __e != nil { return }
				{
//line common.gox:69
					__e = __c.Set("class", "greet"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line common.gox:69
					__e = __c.Any(h.title); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			} else  {
//line common.gox:71
				__e = (doors.ALink{Model: path.Path{Page: path.Home}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
					ctx := __c.Context(); _ = ctx
					__e = __c.Init("a"); if __e != nil { return }
					{
//line common.gox:71
						__e = __c.Set("class", "brand"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line common.gox:71
						__e = __c.Any(i18n.T(l, "app.name")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				return })); if __e != nil { return }
			}
			__e = __c.Init("div"); if __e != nil { return }
			{
//line common.gox:73
				__e = __c.Set("class", "top-right"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("span"); if __e != nil { return }
				{
//line common.gox:74
					__e = __c.Set("class", "streak-badge"); if __e != nil { return }
//line common.gox:74
					__e = __c.Set("title", i18n.T(l, "home.streak")); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Text("🔥 "); if __e != nil { return }
//line common.gox:74
					__e = __c.Any(h.streak); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
//line common.gox:75
				__e = (doors.ALink{Model: path.Path{Page: path.Settings}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
					ctx := __c.Context(); _ = ctx
					__e = __c.Init("a"); if __e != nil { return }
					{
//line common.gox:75
						__e = __c.Set("class", "icon-btn"); if __e != nil { return }
//line common.gox:75
						__e = __c.Set("aria-label", i18n.T(l, "nav.settings")); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Text("⚙️"); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				return })); if __e != nil { return }
//line common.gox:76
				if !h.sess.IsGuest {
//line common.gox:77
					__e = (doors.AClick{On: h.logout}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
						ctx := __c.Context(); _ = ctx
						__e = __c.Init("button"); if __e != nil { return }
						{
//line common.gox:77
							__e = __c.Set("class", "icon-btn"); if __e != nil { return }
//line common.gox:77
							__e = __c.Set("type", "button"); if __e != nil { return }
//line common.gox:77
							__e = __c.Set("aria-label", i18n.T(l, "nav.logout")); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Text("⎋"); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					return })); if __e != nil { return }
				}
			}
			__e = __c.Close(); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line common.gox:81
}

func (h header) logout(ctx context.Context, r doors.RequestPointer) bool {
	if c, err := r.GetCookie("session"); err == nil {
		_ = app.DB.DeleteSession(c.Value)
	}
	r.SetCookie(&http.Cookie{Name: "session", Path: "/", MaxAge: -1, HttpOnly: true})
	doors.SessionExpire(ctx, 0)
	h.auth.Update(ctx, auth.Session{Lang: h.sess.Lang, Rest: 20})
	h.path.Update(ctx, path.Path{Page: path.Login})
	return false
}

// minutes renders a seconds count as a rounded minute number (min 1).
func minutes(sec int) int {
	m := (sec + 59) / 60
	if m < 1 {
		m = 1
	}
	return m
}
