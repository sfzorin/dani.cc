// Managed by GoX v0.1.36

//line workout.gox:1
package pages

import (
	"context"
	"encoding/json"
	"strconv"

	"morning/assets"
	"morning/internal/app"
	"morning/internal/auth"
	"morning/internal/content"
	"morning/internal/i18n"
	"morning/path"

	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

// completeReq is the body sent by the client when a day is finished, optionally
// carrying the post-workout feedback.
type completeReq struct {
	Day         int              `json:"day"`
	HasFeedback bool             `json:"hasFeedback"`
	Feedback    content.Feedback `json:"feedback"`
}

type playerItem struct {
	ID      string `json:"id"`
	Unit    string `json:"unit"`
	Value   int    `json:"value"`
	Slot    string `json:"slot"`
	PerSide bool   `json:"perSide"`
	Round   int    `json:"round"`
	Rest    int    `json:"rest"`
	Name      string `json:"name"`
	Hint      string `json:"hint"`
	Warning   string `json:"warning"`
	Technique string `json:"technique"`
	Mistake   string `json:"mistake"`
	SVG       string `json:"svg"`
	Video     string `json:"video"`
}

type playerPayload struct {
	Day    int               `json:"day"`
	Lang   string            `json:"lang"`
	BCP47  string            `json:"bcp47"`
	Voice  bool              `json:"voice"`
	Name   string            `json:"name"`
	Items  []playerItem      `json:"items"`
	Labels map[string]string `json:"labels"`
	Cues   map[string]string `json:"cues"`
	Enc    []string          `json:"enc"`
}

// workoutPage renders the client-driven player for one day.
type workoutPage struct {
	sess auth.Session
	day  int
	path doors.Source[path.Path]
}

func (w workoutPage) clampDay() int {
	d := w.day
	if d < 1 {
		d = 1
	}
	if d > content.TotalDays {
		d = content.TotalDays
	}
	return d
}

// payloadJSON serializes everything the client player needs.
func (w workoutPage) payloadJSON() string {
	l := i18n.Lang(w.sess.Lang)
	day := w.clampDay()
	adapt, _ := app.DB.GetAdaptState(w.sess.UserID)
	streak, _ := app.DB.Streak(w.sess.UserID, app.Today())
	wk := content.BuildAdapted(day, w.sess.Rest, adapt.ToAdapt().SuppressIncreaseOnGap(streak > 0))

	items := make([]playerItem, 0, len(wk.Items))
	for _, it := range wk.Items {
		items = append(items, playerItem{
			ID:      it.ExerciseID,
			Unit:    string(it.Unit),
			Value:   it.Value,
			Slot:    string(it.Slot),
			PerSide: it.PerSide,
			Round:   it.Round,
			Rest:    it.Rest,
			Name:      i18n.Exercise(l, it.ExerciseID),
			Hint:      i18n.Hint(l, it.ExerciseID),
			Warning:   i18n.Warning(l, it.ExerciseID),
			Technique: i18n.Technique(l, it.ExerciseID),
			Mistake:   i18n.Mistake(l, it.ExerciseID),
			SVG:       "/static/ex/" + it.ExerciseID + ".svg",
			Video:     content.Pool[it.ExerciseID].Video,
		})
	}

	labels := map[string]string{
		"warmup":     i18n.T(l, "workout.warmup"),
		"main":       i18n.T(l, "workout.main"),
		"cooldown":   i18n.T(l, "workout.cooldown"),
		"round":      i18n.T(l, "workout.round"),
		"rest":       i18n.T(l, "workout.rest"),
		"ready":      i18n.T(l, "workout.get_ready"),
		"next":       i18n.T(l, "workout.next"),
		"pause":      i18n.T(l, "workout.pause"),
		"resume":     i18n.T(l, "workout.resume"),
		"skip":       i18n.T(l, "workout.skip"),
		"prev":       i18n.T(l, "workout.prev"),
		"done":       i18n.T(l, "workout.done"),
		"sec":        i18n.T(l, "workout.sec"),
		"reps":       i18n.T(l, "workout.x_reps"),
		"breaths":    i18n.T(l, "workout.breaths"),
		"side":       i18n.T(l, "workout.side"),
		"done_title": i18n.T(l, "done.title"),
		"bravo":      i18n.T(l, "done.bravo"),
		"day_done":   i18n.T(l, "done.day_done"),
		"streak":     i18n.T(l, "done.streak"),
		"back_home":  i18n.T(l, "done.back_home"),
	}
	cues := map[string]string{
		"ready":       i18n.Cue(l, "ready"),
		"go":          i18n.Cue(l, "go"),
		"rest":        i18n.Cue(l, "rest"),
		"next":        i18n.Cue(l, "next"),
		"halfway":     i18n.Cue(l, "halfway"),
		"last":        i18n.Cue(l, "last"),
		"well_done":   i18n.Cue(l, "well_done"),
		"workout_end": i18n.Cue(l, "workout_end"),
		"switch_side": i18n.Cue(l, "switch_side"),
	}
	enc := make([]string, 0, 6)
	for n := 0; n < 6; n++ {
		enc = append(enc, i18n.T(l, "enc."+strconv.Itoa(n)))
	}

	p := playerPayload{
		Day:    day,
		Lang:   w.sess.Lang,
		BCP47:  i18n.BCP47(l),
		Voice:  w.sess.Voice,
		Name:   w.sess.Name,
		Items:  items,
		Labels: labels,
		Cues:   cues,
		Enc:    enc,
	}
	b, _ := json.Marshal(p)
	return string(b)
}

// complete is the JS hook fired when a day finishes. It records today's
// activity (streak/calendar), advances the looping 30-day cycle when the
// finished day is the current one, and folds the optional post-workout feedback
// into the adaptation state that shapes the next workout. Returns the streak.
func (w workoutPage) complete(ctx context.Context, r doors.RequestHook[completeReq]) (any, bool) {
	uid := w.sess.UserID
	day := w.clampDay()
	today := app.Today()
	data := r.Data()

	_ = app.DB.RecordActivity(uid, today)

	// One session passed: decay block counters and clear the one-shot
	// level/force-R that applied to this workout.
	st, _ := app.DB.GetAdaptState(uid)
	st = st.Decay()

	if pos, err := app.DB.Position(uid); err == nil {
		if day == pos%content.TotalDays+1 {
			_ = app.DB.AdvancePosition(uid)
		}
	}

	if data.HasFeedback {
		st = content.NextAdapt(st, data.Feedback)
		_ = app.DB.AddFeedback(uid, today, day, data.Feedback)
	}
	_ = app.DB.SetAdaptState(uid, st)

	streak, _ := app.DB.Streak(uid, today)
	return streak, false
}

// painGroup renders a 3-chip pain question (None / A bit / Pain → 0 / 3 / 6).
//line workout.gox:190
func (w workoutPage) painGroup(l i18n.Lang, key, labelKey string) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
		__e = __c.Init("div"); if __e != nil { return }
		{
//line workout.gox:191
			__e = __c.Set("class", "fb-q fb-group"); if __e != nil { return }
//line workout.gox:191
			__e = __c.Set("data-key", key); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
			__e = __c.Init("label"); if __e != nil { return }
			{
//line workout.gox:192
				__e = __c.Set("class", "fb-label"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
//line workout.gox:192
				__e = __c.Any(i18n.T(l, labelKey)); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line workout.gox:193
				__e = __c.Set("class", "fb-chips"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("button"); if __e != nil { return }
				{
//line workout.gox:194
					__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:194
					__e = __c.Set("class", "fb-chip active"); if __e != nil { return }
//line workout.gox:194
					__e = __c.Set("data-val", "0"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line workout.gox:194
					__e = __c.Any(i18n.T(l, "fb.none")); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("button"); if __e != nil { return }
				{
//line workout.gox:195
					__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:195
					__e = __c.Set("class", "fb-chip"); if __e != nil { return }
//line workout.gox:195
					__e = __c.Set("data-val", "3"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line workout.gox:195
					__e = __c.Any(i18n.T(l, "fb.some")); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("button"); if __e != nil { return }
				{
//line workout.gox:196
					__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:196
					__e = __c.Set("class", "fb-chip"); if __e != nil { return }
//line workout.gox:196
					__e = __c.Set("data-val", "6"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line workout.gox:196
					__e = __c.Any(i18n.T(l, "fb.pain")); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line workout.gox:199
}

//line workout.gox:201
func (w workoutPage) Main() gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line workout.gox:203
		l := i18n.Lang(w.sess.Lang)
		payloadStr := w.payloadJSON()

		__e = __c.Init("title"); if __e != nil { return }
		{
			__e = __c.Submit(); if __e != nil { return }
//line workout.gox:206
			__e = __c.Any(i18n.T(l, "app.name")); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
		__e = __c.Init("div"); if __e != nil { return }
		{
//line workout.gox:207
			__e = __c.Set("class", "player"); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line workout.gox:208
				__e = __c.Set("class", "player-top"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
//line workout.gox:209
				__e = (doors.ALink{Model: path.Path{Page: path.Home}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
					ctx := __c.Context(); _ = ctx
					__e = __c.Init("a"); if __e != nil { return }
					{
//line workout.gox:209
						__e = __c.Set("class", "quit"); if __e != nil { return }
//line workout.gox:209
						__e = __c.Set("aria-label", i18n.T(l, "workout.quit")); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Text("✕"); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				return })); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:210
					__e = __c.Set("class", "pbar"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line workout.gox:210
						__e = __c.Set("class", "pbar-fill"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:211
					__e = __c.Set("class", "phase"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("button"); if __e != nil { return }
				{
//line workout.gox:212
					__e = __c.Set("class", "info-btn"); if __e != nil { return }
//line workout.gox:212
					__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:212
					__e = __c.Set("aria-label", i18n.T(l, "info.technique")); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Text("ⓘ"); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line workout.gox:215
				__e = __c.Set("class", "stage"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:216
					__e = __c.Set("class", "ex-index"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.InitVoid("img"); if __e != nil { return }
				{
//line workout.gox:217
					__e = __c.Set("class", "ex-svg"); if __e != nil { return }
//line workout.gox:217
					__e = __c.Set("alt", ""); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:218
					__e = __c.Set("class", "ex-name"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:219
					__e = __c.Set("class", "ex-side"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:220
					__e = __c.Set("class", "ex-value"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:221
					__e = __c.Set("class", "ex-hint"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:222
					__e = __c.Set("class", "ex-warn"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line workout.gox:225
				__e = __c.Set("class", "controls"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("button"); if __e != nil { return }
				{
//line workout.gox:226
					__e = __c.Set("class", "ctl-prev btn"); if __e != nil { return }
//line workout.gox:226
					__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:226
					__e = __c.Set("aria-label", i18n.T(l, "workout.prev")); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Text("◀"); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("button"); if __e != nil { return }
				{
//line workout.gox:227
					__e = __c.Set("class", "ctl-pause btn"); if __e != nil { return }
//line workout.gox:227
					__e = __c.Set("type", "button"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("button"); if __e != nil { return }
				{
//line workout.gox:228
					__e = __c.Set("class", "ctl-skip btn"); if __e != nil { return }
//line workout.gox:228
					__e = __c.Set("type", "button"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("button"); if __e != nil { return }
				{
//line workout.gox:229
					__e = __c.Set("class", "ctl-done btn primary"); if __e != nil { return }
//line workout.gox:229
					__e = __c.Set("type", "button"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line workout.gox:232
				__e = __c.Set("class", "rest-overlay hidden"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:233
					__e = __c.Set("class", "rest-title"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:234
					__e = __c.Set("class", "rest-count"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:235
					__e = __c.Set("class", "rest-next"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:236
					__e = __c.Set("class", "rest-actions"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("button"); if __e != nil { return }
					{
//line workout.gox:237
						__e = __c.Set("class", "rest-back"); if __e != nil { return }
//line workout.gox:237
						__e = __c.Set("type", "button"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("button"); if __e != nil { return }
					{
//line workout.gox:238
						__e = __c.Set("class", "rest-skip"); if __e != nil { return }
//line workout.gox:238
						__e = __c.Set("type", "button"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line workout.gox:243
				__e = __c.Set("class", "info-overlay hidden"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:244
					__e = __c.Set("class", "info-card"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.InitVoid("img"); if __e != nil { return }
					{
//line workout.gox:245
						__e = __c.Set("class", "info-media"); if __e != nil { return }
//line workout.gox:245
						__e = __c.Set("alt", ""); if __e != nil { return }
					}
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("h2"); if __e != nil { return }
					{
//line workout.gox:246
						__e = __c.Set("class", "info-name"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line workout.gox:247
						__e = __c.Set("class", "info-sec info-tech-sec"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("h3"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
//line workout.gox:248
							__e = __c.Any(i18n.T(l, "info.technique")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("p"); if __e != nil { return }
						{
//line workout.gox:249
							__e = __c.Set("class", "info-tech"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line workout.gox:251
						__e = __c.Set("class", "info-sec info-mistake-sec"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("h3"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
//line workout.gox:252
							__e = __c.Any(i18n.T(l, "info.mistake")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("p"); if __e != nil { return }
						{
//line workout.gox:253
							__e = __c.Set("class", "info-mistake"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line workout.gox:255
						__e = __c.Set("class", "info-sec info-warn-sec"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("h3"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Text("⚠️ "); if __e != nil { return }
//line workout.gox:256
							__e = __c.Any(i18n.T(l, "info.warning")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("p"); if __e != nil { return }
						{
//line workout.gox:257
							__e = __c.Set("class", "info-warn"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("button"); if __e != nil { return }
					{
//line workout.gox:259
						__e = __c.Set("class", "info-close btn primary"); if __e != nil { return }
//line workout.gox:259
						__e = __c.Set("type", "button"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line workout.gox:259
						__e = __c.Any(i18n.T(l, "info.close")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line workout.gox:264
				__e = __c.Set("class", "feedback-overlay hidden"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("h2"); if __e != nil { return }
				{
					__e = __c.Submit(); if __e != nil { return }
//line workout.gox:265
					__e = __c.Any(i18n.T(l, "fb.title")); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:266
					__e = __c.Set("class", "fb-q"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("label"); if __e != nil { return }
					{
//line workout.gox:267
						__e = __c.Set("class", "row-between"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("span"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
//line workout.gox:268
							__e = __c.Any(i18n.T(l, "fb.difficulty")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("output"); if __e != nil { return }
						{
//line workout.gox:269
							__e = __c.Set("class", "fb-diff-out"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Text("6"); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.InitVoid("input"); if __e != nil { return }
					{
//line workout.gox:271
						__e = __c.Set("class", "fb-diff"); if __e != nil { return }
//line workout.gox:271
						__e = __c.Set("type", "range"); if __e != nil { return }
//line workout.gox:271
						__e = __c.Set("min", "1"); if __e != nil { return }
//line workout.gox:271
						__e = __c.Set("max", "10"); if __e != nil { return }
//line workout.gox:271
						__e = __c.Set("value", "6"); if __e != nil { return }
					}
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
//line workout.gox:273
				__e = __c.Any(w.painGroup(l, "knee", "fb.knee")); if __e != nil { return }
//line workout.gox:274
				__e = __c.Any(w.painGroup(l, "back", "fb.back")); if __e != nil { return }
//line workout.gox:275
				__e = __c.Any(w.painGroup(l, "shoulder", "fb.shoulder")); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:276
					__e = __c.Set("class", "fb-q fb-group"); if __e != nil { return }
//line workout.gox:276
					__e = __c.Set("data-key", "energy"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("label"); if __e != nil { return }
					{
//line workout.gox:277
						__e = __c.Set("class", "fb-label"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line workout.gox:277
						__e = __c.Any(i18n.T(l, "fb.energy")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line workout.gox:278
						__e = __c.Set("class", "fb-chips"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("button"); if __e != nil { return }
						{
//line workout.gox:279
							__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:279
							__e = __c.Set("class", "fb-chip"); if __e != nil { return }
//line workout.gox:279
							__e = __c.Set("data-val", "worse"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
//line workout.gox:279
							__e = __c.Any(i18n.T(l, "fb.worse")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("button"); if __e != nil { return }
						{
//line workout.gox:280
							__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:280
							__e = __c.Set("class", "fb-chip active"); if __e != nil { return }
//line workout.gox:280
							__e = __c.Set("data-val", "normal"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
//line workout.gox:280
							__e = __c.Any(i18n.T(l, "fb.normal")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("button"); if __e != nil { return }
						{
//line workout.gox:281
							__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:281
							__e = __c.Set("class", "fb-chip"); if __e != nil { return }
//line workout.gox:281
							__e = __c.Set("data-val", "better"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
//line workout.gox:281
							__e = __c.Any(i18n.T(l, "fb.better")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:284
					__e = __c.Set("class", "fb-actions"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("button"); if __e != nil { return }
					{
//line workout.gox:285
						__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:285
						__e = __c.Set("class", "fb-save btn primary"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line workout.gox:285
						__e = __c.Any(i18n.T(l, "fb.save")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("button"); if __e != nil { return }
					{
//line workout.gox:286
						__e = __c.Set("type", "button"); if __e != nil { return }
//line workout.gox:286
						__e = __c.Set("class", "fb-skip"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line workout.gox:286
						__e = __c.Any(i18n.T(l, "fb.skip")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line workout.gox:290
				__e = __c.Set("class", "done-overlay hidden"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:291
					__e = __c.Set("class", "applause"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Text("👏"); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:292
					__e = __c.Set("class", "done-title"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:293
					__e = __c.Set("class", "done-enc"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line workout.gox:294
					__e = __c.Set("class", "done-streak"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
//line workout.gox:295
				__e = (doors.ALink{Model: path.Path{Page: path.Home}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
					ctx := __c.Context(); _ = ctx
					__e = __c.Init("a"); if __e != nil { return }
					{
//line workout.gox:295
						__e = __c.Set("class", "btn primary back-home"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line workout.gox:295
						__e = __c.Any(i18n.T(l, "done.back_home")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				return })); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("script"); if __e != nil { return }
			{
//line workout.gox:299
				__e = __c.Set("src", assets.Player); if __e != nil { return }
				__e = __c.Set("inline", true); if __e != nil { return }
//line workout.gox:301
				__e = __c.Set("data:payload", payloadStr); if __e != nil { return }
//line workout.gox:302
				__e = __c.Modify(doors.AHook[completeReq]{Name: "complete", On: w.complete}); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Raw(""); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line workout.gox:305
}
