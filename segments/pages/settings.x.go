// Managed by GoX v0.1.36

//line settings.gox:1
package pages

import (
	"context"
	"strings"

	"morning/internal/app"
	"morning/internal/auth"
	"morning/internal/i18n"
	"morning/path"

	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

// settingsForm is the decoded settings submission.
type settingsForm struct {
	Name  string `form:"name"`
	Rest  int    `form:"rest"`
	Voice bool   `form:"voice"`
}

// settingsPage lets the user change rest time, language, voice and reset progress.
type settingsPage struct {
	sess auth.Session
	auth doors.Source[auth.Session]
	path doors.Source[path.Path]
}

//line settings.gox:30
func (s settingsPage) Main() gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line settings.gox:32
		l := i18n.Lang(s.sess.Lang)
		saved := doors.NewSource(false)
		rest := s.sess.Rest
		if rest < 10 {
			rest = 10
		}
		if rest > 40 {
			rest = 40
		}
		restVal := doors.NewSource(rest)
		streak, _ := app.DB.Streak(s.sess.UserID, app.Today())

		__e = __c.Init("title"); if __e != nil { return }
		{
			__e = __c.Submit(); if __e != nil { return }
//line settings.gox:44
			__e = __c.Any(i18n.T(l, "settings.title")); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
		__e = __c.Init("main"); if __e != nil { return }
		{
//line settings.gox:45
			__e = __c.Set("class", "screen"); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
//line settings.gox:46
			__e = __c.Any(header{sess: s.sess, auth: s.auth, path: s.path, streak: streak}); if __e != nil { return }
			__e = __c.Init("section"); if __e != nil { return }
			{
//line settings.gox:47
				__e = __c.Set("class", "settings"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("h1"); if __e != nil { return }
				{
					__e = __c.Submit(); if __e != nil { return }
//line settings.gox:48
					__e = __c.Any(i18n.T(l, "settings.title")); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line settings.gox:51
					__e = __c.Set("class", "field"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("label"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:52
						__e = __c.Any(i18n.T(l, "settings.language")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
//line settings.gox:53
					__e = __c.Any(langSwitch{sess: s.sess, auth: s.auth, persist: true}); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
//line settings.gox:56
				__e = (doors.ASubmit[settingsForm]{On: s.save(saved)}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
					ctx := __c.Context(); _ = ctx
					__e = __c.Init("form"); if __e != nil { return }
					{
//line settings.gox:56
						__e = __c.Set("class", "form"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:57
							__e = __c.Set("class", "field"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Init("label"); if __e != nil { return }
							{
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:58
								__e = __c.Any(i18n.T(l, "auth.name")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
							__e = __c.InitVoid("input"); if __e != nil { return }
							{
//line settings.gox:59
								__e = __c.Set("class", "text-input"); if __e != nil { return }
//line settings.gox:59
								__e = __c.Set("name", "name"); if __e != nil { return }
//line settings.gox:59
								__e = __c.Set("type", "text"); if __e != nil { return }
//line settings.gox:59
								__e = __c.Set("maxlength", "40"); if __e != nil { return }
//line settings.gox:59
								__e = __c.Set("value", s.sess.Name); if __e != nil { return }
//line settings.gox:59
								__e = __c.Set("placeholder", i18n.T(l, "auth.name")); if __e != nil { return }
							}
							__e = __c.Submit(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:61
							__e = __c.Set("class", "field"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Init("label"); if __e != nil { return }
							{
//line settings.gox:62
								__e = __c.Set("class", "row-between"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
								__e = __c.Init("span"); if __e != nil { return }
								{
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:63
									__e = __c.Any(i18n.T(l, "settings.rest")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
								__e = __c.Init("output"); if __e != nil { return }
								{
//line settings.gox:64
									__e = __c.Set("class", "rest-out"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:65
									__e = __c.Any(restVal.Bind(func(v int) gox.Elem {
								return gox.Elem(func(__c gox.Cursor) (__e error) {
									ctx := __c.Context(); _ = ctx
//line settings.gox:66
									__e = __c.Many(v, " ", i18n.T(l, "workout.sec")); if __e != nil { return }
								return })
//line settings.gox:67
							})); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
							__e = __c.InitVoid("input"); if __e != nil { return }
							{
//line settings.gox:71
								__e = __c.Set("type", "range"); if __e != nil { return }
//line settings.gox:72
								__e = __c.Set("name", "rest"); if __e != nil { return }
//line settings.gox:73
								__e = __c.Set("min", "10"); if __e != nil { return }
//line settings.gox:74
								__e = __c.Set("max", "40"); if __e != nil { return }
//line settings.gox:75
								__e = __c.Set("step", "5"); if __e != nil { return }
//line settings.gox:76
								__e = __c.Set("value", rest); if __e != nil { return }
//line settings.gox:77
								__e = __c.Modify(doors.AChange{On: func(ctx context.Context, r doors.RequestEvent[doors.ChangeEvent]) bool {
							if n := r.Event().Number; n != nil {
								restVal.Update(ctx, int(*n))
							}
							return false
						}}); if __e != nil { return }
							}
							__e = __c.Submit(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:85
							__e = __c.Set("class", "field"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Init("label"); if __e != nil { return }
							{
//line settings.gox:86
								__e = __c.Set("class", "switch-label"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
								__e = __c.InitVoid("input"); if __e != nil { return }
								{
//line settings.gox:87
									__e = __c.Set("type", "checkbox"); if __e != nil { return }
//line settings.gox:87
									__e = __c.Set("name", "voice"); if __e != nil { return }
//line settings.gox:87
									__e = __c.Set("value", "true"); if __e != nil { return }
//line settings.gox:87
									__e = __c.Set("checked", func() any { return s.sess.Voice }()); if __e != nil { return }
								}
								__e = __c.Submit(); if __e != nil { return }
								__e = __c.Init("span"); if __e != nil { return }
								{
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:88
									__e = __c.Any(i18n.T(l, "settings.voice")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("button"); if __e != nil { return }
						{
//line settings.gox:92
							__e = __c.Set("class", "btn primary"); if __e != nil { return }
//line settings.gox:92
							__e = __c.Set("type", "submit"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
//line settings.gox:92
							__e = __c.Any(i18n.T(l, "settings.save")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
//line settings.gox:93
						__e = __c.Any(saved.Bind(func(ok bool) gox.Elem {
					return gox.Elem(func(__c gox.Cursor) (__e error) {
						ctx := __c.Context(); _ = ctx
//line settings.gox:94
						if ok {
							__e = __c.Init("span"); if __e != nil { return }
							{
//line settings.gox:95
								__e = __c.Set("class", "saved-msg"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:95
								__e = __c.Any(i18n.T(l, "settings.saved")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						}
					return })
//line settings.gox:97
				})); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				return })); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line settings.gox:100
					__e = __c.Set("class", "field account"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("label"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:101
						__e = __c.Any(i18n.T(l, "settings.account")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
//line settings.gox:102
					if s.sess.IsGuest {
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:103
							__e = __c.Set("class", "email"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
//line settings.gox:103
							__e = __c.Any(i18n.T(l, "guest.note")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
//line settings.gox:104
						__e = (doors.ALink{Model: path.Path{Page: path.Register}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
							ctx := __c.Context(); _ = ctx
							__e = __c.Init("a"); if __e != nil { return }
							{
//line settings.gox:104
								__e = __c.Set("class", "btn primary small"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:104
								__e = __c.Any(i18n.T(l, "auth.register")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						return })); if __e != nil { return }
					} else  {
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:106
							__e = __c.Set("class", "email"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
//line settings.gox:106
							__e = __c.Any(s.sess.Email); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					}
				}
				__e = __c.Close(); if __e != nil { return }
//line settings.gox:110
				__e = (doors.AClick{On: s.reset}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
					ctx := __c.Context(); _ = ctx
					__e = __c.Init("button"); if __e != nil { return }
					{
//line settings.gox:110
						__e = __c.Set("class", "btn danger"); if __e != nil { return }
//line settings.gox:110
						__e = __c.Set("type", "button"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:110
						__e = __c.Any(i18n.T(l, "settings.reset")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				return })); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line settings.gox:113
}

// save persists rest + voice and updates the reactive session.
func (s settingsPage) save(saved doors.Source[bool]) func(context.Context, doors.RequestForm[settingsForm]) bool {
	return func(ctx context.Context, r doors.RequestForm[settingsForm]) bool {
		f := r.Data()
		rest := f.Rest
		if rest < 10 {
			rest = 10
		}
		if rest > 40 {
			rest = 40
		}
		name := strings.TrimSpace(f.Name)
		if len(name) > 40 {
			name = name[:40]
		}
		_ = app.DB.UpdateSettings(s.sess.UserID, name, s.sess.Lang, rest, f.Voice)
		s.auth.Mutate(ctx, func(cur auth.Session) auth.Session {
			cur.Name = name
			cur.Rest = rest
			cur.Voice = f.Voice
			return cur
		})
		saved.Update(ctx, true)
		return false
	}
}

// reset clears all progress and returns to the home screen.
func (s settingsPage) reset(ctx context.Context, r doors.RequestPointer) bool {
	_ = app.DB.ResetProgress(s.sess.UserID)
	s.path.Update(ctx, path.Path{Page: path.Home})
	return false
}
