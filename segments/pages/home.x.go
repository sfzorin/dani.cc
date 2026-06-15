// Managed by GoX v0.1.36

//line home.gox:1
package pages

import (
	"context"
	"fmt"
	"time"

	"morning30/internal/app"
	"morning30/internal/auth"
	"morning30/internal/content"
	"morning30/internal/i18n"
	"morning30/path"

	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

// homePage is the dashboard: today's session, the date-based streak calendar,
// and the separate looping 30-day cycle track.
type homePage struct {
	sess auth.Session
	auth doors.Source[auth.Session]
	path doors.Source[path.Path]
}

//line home.gox:26
func (h homePage) Main() gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line home.gox:28
		l := i18n.Lang(h.sess.Lang)
		today := app.Today()
		pos, _ := app.DB.Position(h.sess.UserID)
		streak, _ := app.DB.Streak(h.sess.UserID, today)
		doneToday, _ := app.DB.HasActivity(h.sess.UserID, today)

		cycle := pos/content.TotalDays + 1
		curDay := pos%content.TotalDays + 1
		w := app.UserWorkout(h.sess.UserID, curDay, h.sess.Rest)
		acked, _ := app.DB.SafetyAck(h.sess.UserID)
		safety := doors.NewSource(acked)

		startLabel := i18n.T(l, "home.start")
		if pos > 0 {
			startLabel = i18n.T(l, "home.continue")
		}
		if doneToday {
			startLabel = i18n.T(l, "home.repeat")
		}

		greeting := i18n.T(l, "home.morning")
		if h.sess.Name != "" {
			greeting += ", " + h.sess.Name
		}
		greeting += "!"

		monthOffset := doors.NewSource(0)

		__e = __c.Init("title"); if __e != nil { return }
		{
			__e = __c.Submit(); if __e != nil { return }
//line home.gox:56
			__e = __c.Any(i18n.T(l, "app.name")); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
		__e = __c.Init("main"); if __e != nil { return }
		{
//line home.gox:57
			__e = __c.Set("class", "screen"); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
//line home.gox:58
			__e = __c.Any(header{sess: h.sess, auth: h.auth, path: h.path, streak: streak, title: greeting}); if __e != nil { return }
//line home.gox:60
			__e = __c.Any(safety.Bind(func(ok bool) gox.Elem {
			return gox.Elem(func(__c gox.Cursor) (__e error) {
				ctx := __c.Context(); _ = ctx
//line home.gox:61
				if !ok && h.sess.IsGuest {
					__e = __c.Init("div"); if __e != nil { return }
					{
//line home.gox:62
						__e = __c.Set("class", "safety-card"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("h3"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Text("⚠️ "); if __e != nil { return }
//line home.gox:63
							__e = __c.Any(i18n.T(l, "safety.title")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("p"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
//line home.gox:64
							__e = __c.Any(i18n.T(l, "safety.text")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
//line home.gox:65
						__e = (doors.AClick{On: func(ctx context.Context, r doors.RequestPointer) bool {
						_ = app.DB.AckSafety(h.sess.UserID)
						safety.Update(ctx, true)
						return false
					}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
							ctx := __c.Context(); _ = ctx
							__e = __c.Init("button"); if __e != nil { return }
							{
//line home.gox:69
								__e = __c.Set("class", "btn primary"); if __e != nil { return }
//line home.gox:69
								__e = __c.Set("type", "button"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line home.gox:69
								__e = __c.Any(i18n.T(l, "safety.ok")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						return })); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
			return })
//line home.gox:72
		})); if __e != nil { return }
//line home.gox:74
			if h.sess.IsGuest {
				__e = __c.Init("div"); if __e != nil { return }
				{
//line home.gox:75
					__e = __c.Set("class", "guest-note"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("span"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line home.gox:76
						__e = __c.Any(i18n.T(l, "guest.note")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
//line home.gox:77
					__e = (doors.ALink{Model: path.Path{Page: path.Register}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
						ctx := __c.Context(); _ = ctx
						__e = __c.Init("a"); if __e != nil { return }
						{
//line home.gox:77
							__e = __c.Set("class", "btn primary small"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
//line home.gox:77
							__e = __c.Any(i18n.T(l, "auth.register")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					return })); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Init("section"); if __e != nil { return }
			{
//line home.gox:81
				__e = __c.Set("class", "hero"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
//line home.gox:82
				if doneToday {
					__e = __c.Init("p"); if __e != nil { return }
					{
//line home.gox:83
						__e = __c.Set("class", "done-banner"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line home.gox:83
						__e = __c.Any(i18n.T(l, "home.done_today")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Init("div"); if __e != nil { return }
				{
//line home.gox:85
					__e = __c.Set("class", "hero-cycle"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line home.gox:85
					__e = __c.Many(i18n.T(l, "home.cycle"), " ", cycle, "  ·  ", levelLabel(l, h.sess.Level)); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line home.gox:86
					__e = __c.Set("class", "day-num"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line home.gox:86
					__e = __c.Many(i18n.T(l, "home.day"), " ", curDay); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line home.gox:87
					__e = __c.Set("class", "day-meta"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("span"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line home.gox:88
						__e = __c.Many(len(w.Items), " ", i18n.T(l, "home.exercises")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("span"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line home.gox:89
						__e = __c.Many(i18n.T(l, "home.est"), " ", minutes(w.EstSec), " ", i18n.T(l, "home.min")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
//line home.gox:91
				__e = (doors.ALink{Model: path.Path{Page: path.Day, Day: curDay}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
					ctx := __c.Context(); _ = ctx
					__e = __c.Init("a"); if __e != nil { return }
					{
//line home.gox:91
						__e = __c.Set("class", "btn primary big"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line home.gox:91
						__e = __c.Any(startLabel); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				return })); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("section"); if __e != nil { return }
			{
//line home.gox:94
				__e = __c.Set("class", "streak-row"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("span"); if __e != nil { return }
				{
//line home.gox:95
					__e = __c.Set("class", "fire"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Text("🔥"); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("span"); if __e != nil { return }
				{
//line home.gox:96
					__e = __c.Set("class", "streak-n"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line home.gox:96
					__e = __c.Any(streak); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("span"); if __e != nil { return }
				{
//line home.gox:97
					__e = __c.Set("class", "streak-label"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line home.gox:97
					__e = __c.Many(i18n.T(l, "home.streak"), " · ", streak, " ", i18n.T(l, "home.days")); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("section"); if __e != nil { return }
			{
//line home.gox:101
				__e = __c.Set("class", "calendar"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("h3"); if __e != nil { return }
				{
					__e = __c.Submit(); if __e != nil { return }
//line home.gox:102
					__e = __c.Any(i18n.T(l, "home.calendar")); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
//line home.gox:103
				__e = __c.Any(h.month(monthOffset)); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("section"); if __e != nil { return }
			{
//line home.gox:107
				__e = __c.Set("class", "cycle"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("h3"); if __e != nil { return }
				{
					__e = __c.Submit(); if __e != nil { return }
//line home.gox:108
					__e = __c.Many(i18n.T(l, "home.cycle"), " ", cycle); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line home.gox:109
					__e = __c.Set("class", "diff"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line home.gox:110
						__e = __c.Set("class", "diff-bar"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("span"); if __e != nil { return }
					{
//line home.gox:111
						__e = __c.Set("class", "diff-label"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line home.gox:111
						__e = __c.Any(i18n.T(l, "home.diff")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line home.gox:113
					__e = __c.Set("class", "cyc-grid"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line home.gox:114
					for d := 1; d <= content.TotalDays; d++ {
//line home.gox:115
						__e = __c.Any(h.cycleCell(d, curDay)); if __e != nil { return }
					}
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line home.gox:120
}

// month renders the date calendar for the current month ± offset, reactively.
//line home.gox:123
func (h homePage) month(monthOffset doors.Source[int]) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line home.gox:124
		__e = __c.Any(monthOffset.Bind(func(off int) gox.Elem {
		return gox.Elem(func(__c gox.Cursor) (__e error) {
			ctx := __c.Context(); _ = ctx
//line home.gox:126
			now := time.Now()
			base := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
			disp := base.AddDate(0, off, 0)
			year := disp.Year()
			month := int(disp.Month())
			daysIn := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.Local).Day()
			lead := (int(disp.Weekday()) + 6) % 7 // Monday-first leading blanks
			first := disp.Format("2006-01-02")
			last := time.Date(year, time.Month(month), daysIn, 0, 0, 0, 0, time.Local).Format("2006-01-02")
			acts, _ := app.DB.ActivityInRange(h.sess.UserID, first, last)
			todayStr := app.Today()
			l := i18n.Lang(h.sess.Lang)

			__e = __c.Init("div"); if __e != nil { return }
			{
//line home.gox:139
				__e = __c.Set("class", "cal-head"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
//line home.gox:140
				__e = (doors.AClick{On: func(ctx context.Context, r doors.RequestPointer) bool {
				monthOffset.Mutate(ctx, func(o int) int { return o - 1 })
				return false
			}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
					ctx := __c.Context(); _ = ctx
					__e = __c.Init("button"); if __e != nil { return }
					{
//line home.gox:143
						__e = __c.Set("class", "cal-nav"); if __e != nil { return }
//line home.gox:143
						__e = __c.Set("type", "button"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Text("‹"); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				return })); if __e != nil { return }
				__e = __c.Init("span"); if __e != nil { return }
				{
//line home.gox:144
					__e = __c.Set("class", "cal-title"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line home.gox:144
					__e = __c.Many(i18n.Month(l, month), " ", year); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
//line home.gox:145
				__e = (doors.AClick{On: func(ctx context.Context, r doors.RequestPointer) bool {
				monthOffset.Mutate(ctx, func(o int) int { return o + 1 })
				return false
			}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
					ctx := __c.Context(); _ = ctx
					__e = __c.Init("button"); if __e != nil { return }
					{
//line home.gox:148
						__e = __c.Set("class", "cal-nav"); if __e != nil { return }
//line home.gox:148
						__e = __c.Set("type", "button"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Text("›"); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				return })); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line home.gox:150
				__e = __c.Set("class", "cal-week"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
//line home.gox:151
				for wd := 0; wd < 7; wd++ {
					__e = __c.Init("span"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line home.gox:152
						__e = __c.Any(i18n.Weekday(l, wd)); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line home.gox:155
				__e = __c.Set("class", "cal-grid month"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
//line home.gox:156
				for b := 0; b < lead; b++ {
					__e = __c.Init("div"); if __e != nil { return }
					{
//line home.gox:157
						__e = __c.Set("class", "cal-blank"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
//line home.gox:159
				for dd := 1; dd <= daysIn; dd++ {
//line home.gox:160
					__e = __c.Any(h.dayCell(year, month, dd, acts, todayStr)); if __e != nil { return }
				}
			}
			__e = __c.Close(); if __e != nil { return }
		return })
//line home.gox:163
	})); if __e != nil { return }
	return })
//line home.gox:164
}

// dayCell renders one date cell, marked when the user worked out that day.
//line home.gox:167
func (h homePage) dayCell(year, month, dd int, acts map[string]bool, todayStr string) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line home.gox:169
		ds := fmt.Sprintf("%04d-%02d-%02d", year, month, dd)
		cls := "cal-day"
		if acts[ds] {
			cls += " act"
		}
		if ds == todayStr {
			cls += " today"
		}

		__e = __c.Init("div"); if __e != nil { return }
		{
//line home.gox:178
			__e = __c.Set("class", cls); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
//line home.gox:179
			if acts[ds] {
				__e = __c.Text("🔥"); if __e != nil { return }
			} else  {
//line home.gox:182
				__e = __c.Any(dd); if __e != nil { return }
			}
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line home.gox:185
}

// cycleCell renders one of the 30 program days; sequential, looping each cycle.
//line home.gox:188
func (h homePage) cycleCell(d, current int) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line home.gox:190
		done := d < current
		isCurrent := d == current
		locked := d > current
		cls := "cell"
		if done {
			cls += " done"
		}
		if isCurrent {
			cls += " current"
		}
		if locked {
			cls += " locked"
		}

//line home.gox:204
		if locked {
			__e = __c.Init("div"); if __e != nil { return }
			{
//line home.gox:205
				__e = __c.Set("class", cls); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
//line home.gox:205
				__e = __c.Any(d); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
		} else  {
//line home.gox:207
			__e = (doors.ALink{Model: path.Path{Page: path.Day, Day: d}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
				ctx := __c.Context(); _ = ctx
				__e = __c.Init("a"); if __e != nil { return }
				{
//line home.gox:207
					__e = __c.Set("class", cls); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line home.gox:208
					if done {
						__e = __c.Text("✓"); if __e != nil { return }
					} else  {
//line home.gox:211
						__e = __c.Any(d); if __e != nil { return }
					}
				}
				__e = __c.Close(); if __e != nil { return }
			return })); if __e != nil { return }
		}
	return })
//line home.gox:215
}
