// Managed by GoX v0.1.36

//line settings.gox:1
package pages

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"morning30/assets"
	"morning30/internal/app"
	"morning30/internal/auth"
	"morning30/internal/catalog"
	"morning30/internal/content"
	"morning30/internal/i18n"
	"morning30/path"

	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

// settingsForm is the decoded settings submission.
type settingsForm struct {
	Name      string `form:"name"`
	Rest      int    `form:"rest"`
	VoiceMode string `form:"voice_mode"`
	Sex       string `form:"sex"`
}

// editState is the editor's working copy: the selected day plus the program as
// JSON. It is a single comparable value so one Source/Bind drives the whole
// editor (day switch + every edit) — nested Binds don't patch reliably.
type editState struct {
	Day  int
	JSON string
	Open string // which row's content preview is expanded ("target:index"), "" = none
}

// settingsPage lets the user change rest time, language, voice and reset progress.
type settingsPage struct {
	sess auth.Session
	auth doors.Source[auth.Session]
	path doors.Source[path.Path]
}

//line settings.gox:45
func (s settingsPage) Main() gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line settings.gox:47
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
		active := app.ActiveProgramKey(s.sess.UserID)
		hasCustom := app.HasCustomProgram(s.sess.UserID)
		cat := app.UserCatalog(s.sess.UserID)
		entries := cat.List(l)
		prog := app.UserProgram(s.sess.UserID) // the active program (built-in or custom)
		editSrc := doors.NewSource(editState{Day: 1, JSON: prog.Marshal()})
		nameSrc := doors.NewSource(prog.Name)
		firstID := "W01"
		if len(entries) > 0 {
			firstID = entries[0].ID
		}
		selSrc := doors.NewSource(firstID)
		levelSrc := doors.NewSource(s.sess.Level)
		curSex := app.DB.GetSex(s.sess.UserID)

		__e = __c.Init("title"); if __e != nil { return }
		{
			__e = __c.Submit(); if __e != nil { return }
//line settings.gox:73
			__e = __c.Any(i18n.T(l, "settings.title")); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
		__e = __c.Init("main"); if __e != nil { return }
		{
//line settings.gox:74
			__e = __c.Set("class", "screen"); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
//line settings.gox:75
			__e = __c.Any(header{sess: s.sess, auth: s.auth, path: s.path, streak: streak}); if __e != nil { return }
			__e = __c.Init("section"); if __e != nil { return }
			{
//line settings.gox:76
				__e = __c.Set("class", "settings"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("div"); if __e != nil { return }
				{
//line settings.gox:77
					__e = __c.Set("class", "settings-head"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("h1"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:78
						__e = __c.Any(i18n.T(l, "settings.title")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
//line settings.gox:79
					__e = (doors.ALink{Model: path.Path{Page: path.Home}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
						ctx := __c.Context(); _ = ctx
						__e = __c.Init("a"); if __e != nil { return }
						{
//line settings.gox:79
							__e = __c.Set("class", "quit"); if __e != nil { return }
//line settings.gox:79
							__e = __c.Set("aria-label", i18n.T(l, "nav.close")); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Text("✕"); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					return })); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
//line settings.gox:82
				vm := s.sess.VoiceMode
				if vm == "" {
					vm = "normal"
				}

				__e = __c.Init("div"); if __e != nil { return }
				{
//line settings.gox:87
					__e = __c.Set("class", "settings-card"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line settings.gox:89
						__e = __c.Set("class", "srow"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("span"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
//line settings.gox:90
							__e = __c.Any(i18n.T(l, "settings.language")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
//line settings.gox:91
						__e = __c.Any(langSwitch{sess: s.sess, auth: s.auth, persist: true}); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
//line settings.gox:94
					__e = (doors.ASubmit[settingsForm]{On: s.save(saved)}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
						ctx := __c.Context(); _ = ctx
						__e = __c.Init("form"); if __e != nil { return }
						{
//line settings.gox:94
							__e = __c.Set("class", "settings-form"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Init("div"); if __e != nil { return }
							{
//line settings.gox:95
								__e = __c.Set("class", "srow"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
								__e = __c.Init("span"); if __e != nil { return }
								{
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:96
									__e = __c.Any(i18n.T(l, "auth.name")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
								__e = __c.InitVoid("input"); if __e != nil { return }
								{
//line settings.gox:97
									__e = __c.Set("class", "text-input compact"); if __e != nil { return }
//line settings.gox:97
									__e = __c.Set("name", "name"); if __e != nil { return }
//line settings.gox:97
									__e = __c.Set("type", "text"); if __e != nil { return }
//line settings.gox:97
									__e = __c.Set("maxlength", "40"); if __e != nil { return }
//line settings.gox:97
									__e = __c.Set("value", s.sess.Name); if __e != nil { return }
								}
								__e = __c.Submit(); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
							__e = __c.Init("div"); if __e != nil { return }
							{
//line settings.gox:99
								__e = __c.Set("class", "srow col"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
								__e = __c.Init("div"); if __e != nil { return }
								{
//line settings.gox:100
									__e = __c.Set("class", "row-between"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
									__e = __c.Init("span"); if __e != nil { return }
									{
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:101
										__e = __c.Any(i18n.T(l, "settings.rest")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
									__e = __c.Init("output"); if __e != nil { return }
									{
//line settings.gox:102
										__e = __c.Set("class", "rest-out"); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:103
										__e = __c.Any(restVal.Bind(func(v int) gox.Elem {
									return gox.Elem(func(__c gox.Cursor) (__e error) {
										ctx := __c.Context(); _ = ctx
//line settings.gox:104
										__e = __c.Many(v, " ", i18n.T(l, "workout.sec")); if __e != nil { return }
									return })
//line settings.gox:105
								})); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
								__e = __c.InitVoid("input"); if __e != nil { return }
								{
//line settings.gox:109
									__e = __c.Set("type", "range"); if __e != nil { return }
//line settings.gox:110
									__e = __c.Set("name", "rest"); if __e != nil { return }
//line settings.gox:111
									__e = __c.Set("min", "10"); if __e != nil { return }
//line settings.gox:112
									__e = __c.Set("max", "40"); if __e != nil { return }
//line settings.gox:113
									__e = __c.Set("step", "5"); if __e != nil { return }
//line settings.gox:114
									__e = __c.Set("value", rest); if __e != nil { return }
//line settings.gox:115
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
//line settings.gox:122
								__e = __c.Set("class", "srow"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
								__e = __c.Init("span"); if __e != nil { return }
								{
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:123
									__e = __c.Any(i18n.T(l, "settings.voice")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
								__e = __c.Init("select"); if __e != nil { return }
								{
//line settings.gox:124
									__e = __c.Set("class", "lang-select compact"); if __e != nil { return }
//line settings.gox:124
									__e = __c.Set("name", "voice_mode"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
									__e = __c.Init("option"); if __e != nil { return }
									{
//line settings.gox:125
										__e = __c.Set("value", "off"); if __e != nil { return }
//line settings.gox:125
										__e = __c.Set("selected", func() any { return vm == "off" }()); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:125
										__e = __c.Any(i18n.T(l, "voice.off")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
									__e = __c.Init("option"); if __e != nil { return }
									{
//line settings.gox:126
										__e = __c.Set("value", "min"); if __e != nil { return }
//line settings.gox:126
										__e = __c.Set("selected", func() any { return vm == "min" }()); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:126
										__e = __c.Any(i18n.T(l, "voice.min")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
									__e = __c.Init("option"); if __e != nil { return }
									{
//line settings.gox:127
										__e = __c.Set("value", "normal"); if __e != nil { return }
//line settings.gox:127
										__e = __c.Set("selected", func() any { return vm == "normal" }()); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:127
										__e = __c.Any(i18n.T(l, "voice.normal")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
									__e = __c.Init("option"); if __e != nil { return }
									{
//line settings.gox:128
										__e = __c.Set("value", "detailed"); if __e != nil { return }
//line settings.gox:128
										__e = __c.Set("selected", func() any { return vm == "detailed" }()); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:128
										__e = __c.Any(i18n.T(l, "voice.detailed")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
							__e = __c.Init("div"); if __e != nil { return }
							{
//line settings.gox:131
								__e = __c.Set("class", "srow"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
								__e = __c.Init("span"); if __e != nil { return }
								{
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:132
									__e = __c.Any(i18n.T(l, "settings.sex")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
								__e = __c.Init("select"); if __e != nil { return }
								{
//line settings.gox:133
									__e = __c.Set("class", "lang-select compact"); if __e != nil { return }
//line settings.gox:133
									__e = __c.Set("name", "sex"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
									__e = __c.Init("option"); if __e != nil { return }
									{
//line settings.gox:134
										__e = __c.Set("value", "m"); if __e != nil { return }
//line settings.gox:134
										__e = __c.Set("selected", func() any { return curSex != "f" }()); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:134
										__e = __c.Any(i18n.T(l, "sex.male")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
									__e = __c.Init("option"); if __e != nil { return }
									{
//line settings.gox:135
										__e = __c.Set("value", "f"); if __e != nil { return }
//line settings.gox:135
										__e = __c.Set("selected", func() any { return curSex == "f" }()); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:135
										__e = __c.Any(i18n.T(l, "sex.female")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
							__e = __c.Init("div"); if __e != nil { return }
							{
//line settings.gox:138
								__e = __c.Set("class", "srow save-row"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
								__e = __c.Init("button"); if __e != nil { return }
								{
//line settings.gox:139
									__e = __c.Set("class", "btn primary"); if __e != nil { return }
//line settings.gox:139
									__e = __c.Set("type", "submit"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:139
									__e = __c.Any(i18n.T(l, "settings.save")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
//line settings.gox:140
								__e = __c.Any(saved.Bind(func(ok bool) gox.Elem {
							return gox.Elem(func(__c gox.Cursor) (__e error) {
								ctx := __c.Context(); _ = ctx
//line settings.gox:141
								if ok {
									__e = __c.Init("span"); if __e != nil { return }
									{
//line settings.gox:142
										__e = __c.Set("class", "saved-msg"); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:142
										__e = __c.Any(i18n.T(l, "settings.saved")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
								}
							return })
//line settings.gox:144
						})); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					return })); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line settings.gox:149
						__e = __c.Set("class", "srow col"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:150
							__e = __c.Set("class", "row-between"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Init("span"); if __e != nil { return }
							{
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:151
								__e = __c.Any(i18n.T(l, "settings.difficulty")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
							__e = __c.Init("output"); if __e != nil { return }
							{
//line settings.gox:152
								__e = __c.Set("class", "rest-out"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:153
								__e = __c.Any(levelSrc.Bind(func(v int) gox.Elem {
								return gox.Elem(func(__c gox.Cursor) (__e error) {
									ctx := __c.Context(); _ = ctx
//line settings.gox:154
									__e = __c.Any(levelLabel(l, v)); if __e != nil { return }
								return })
//line settings.gox:155
							})); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.InitVoid("input"); if __e != nil { return }
						{
//line settings.gox:158
							__e = __c.Set("type", "range"); if __e != nil { return }
//line settings.gox:158
							__e = __c.Set("min", "-3"); if __e != nil { return }
//line settings.gox:158
							__e = __c.Set("max", "3"); if __e != nil { return }
//line settings.gox:158
							__e = __c.Set("step", "1"); if __e != nil { return }
//line settings.gox:158
							__e = __c.Set("value", s.sess.Level); if __e != nil { return }
//line settings.gox:158
							__e = __c.Modify(doors.AChange{On: s.setLevel(levelSrc)}); if __e != nil { return }
						}
						__e = __c.Submit(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("details"); if __e != nil { return }
					{
//line settings.gox:162
						__e = __c.Set("class", "cycle-details"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("summary"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
//line settings.gox:163
							__e = __c.Any(i18n.T(l, "settings.cycle")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:166
							__e = __c.Set("class", "srow"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Init("span"); if __e != nil { return }
							{
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:167
								__e = __c.Any(i18n.T(l, "settings.program")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
							__e = __c.Init("select"); if __e != nil { return }
							{
//line settings.gox:168
								__e = __c.Set("class", "lang-select compact"); if __e != nil { return }
//line settings.gox:168
								__e = __c.Modify(doors.AChange{On: s.selectProgram}); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:169
								for _, p := range content.StandardPrograms {
									__e = __c.Init("option"); if __e != nil { return }
									{
//line settings.gox:170
										__e = __c.Set("value", p.Key); if __e != nil { return }
//line settings.gox:170
										__e = __c.Set("selected", func() any { return active == p.Key }()); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:170
										__e = __c.Any(p.Name); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
								}
//line settings.gox:172
								if hasCustom {
									__e = __c.Init("option"); if __e != nil { return }
									{
//line settings.gox:173
										__e = __c.Set("value", "custom"); if __e != nil { return }
//line settings.gox:173
										__e = __c.Set("selected", func() any { return active == "custom" }()); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:173
										__e = __c.Any(i18n.T(l, "settings.my_version")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
								}
							}
							__e = __c.Close(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:178
							__e = __c.Set("class", "srow"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Init("span"); if __e != nil { return }
							{
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:179
								__e = __c.Any(i18n.T(l, "home.day")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
							__e = __c.Init("select"); if __e != nil { return }
							{
//line settings.gox:180
								__e = __c.Set("class", "lang-select compact"); if __e != nil { return }
//line settings.gox:180
								__e = __c.Modify(doors.AChange{On: s.setDay(editSrc)}); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:181
								for d := 1; d <= len(prog.Days); d++ {
									__e = __c.Init("option"); if __e != nil { return }
									{
//line settings.gox:182
										__e = __c.Set("value", d); if __e != nil { return }
//line settings.gox:182
										__e = __c.Set("selected", func() any { return d == 1 }()); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:182
										__e = __c.Any(d); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
								}
							}
							__e = __c.Close(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
//line settings.gox:189
						__e = __c.Any(editSrc.Bind(func(st editState) gox.Elem {
						return gox.Elem(func(__c gox.Cursor) (__e error) {
							ctx := __c.Context(); _ = ctx
//line settings.gox:191
							p, _ := content.ParseResolved(st.JSON)
							idx := clampDay(st.Day, len(p.Days)) - 1
							if idx < 0 {
								idx = 0
							}

							__e = __c.Init("div"); if __e != nil { return }
							{
//line settings.gox:197
								__e = __c.Set("class", "cyc-head"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:197
								__e = __c.Many(i18n.T(l, "settings.warmup_series"), " · ", p.WarmupRounds, "×"); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
//line settings.gox:198
							__e = __c.Any(s.editList(editSrc, cat, entries, l, p.Warmup, "warmup", st.Open)); if __e != nil { return }
//line settings.gox:200
							if len(p.Days) > 0 {
//line settings.gox:201
								rd := p.Days[idx]
								__e = __c.Init("div"); if __e != nil { return }
								{
//line settings.gox:202
									__e = __c.Set("class", "cyc-head"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:202
									__e = __c.Many(i18n.T(l, "settings.main_block"), " · ", i18n.T(l, "home.day"), " ", st.Day, " · ", len(rd.Items), " ", i18n.T(l, "home.exercises")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
//line settings.gox:203
								__e = __c.Any(s.editList(editSrc, cat, entries, l, rd.Items, "day", st.Open)); if __e != nil { return }
							}
							__e = __c.Init("div"); if __e != nil { return }
							{
//line settings.gox:206
								__e = __c.Set("class", "cyc-head"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:206
								__e = __c.Any(i18n.T(l, "settings.cooldown_series")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
//line settings.gox:207
							__e = __c.Any(s.editList(editSrc, cat, entries, l, p.Cooldown, "cooldown", st.Open)); if __e != nil { return }
						return })
//line settings.gox:208
					})); if __e != nil { return }
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:210
							__e = __c.Set("class", "srow cyc-save"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.InitVoid("input"); if __e != nil { return }
							{
//line settings.gox:211
								__e = __c.Set("class", "text-input compact"); if __e != nil { return }
//line settings.gox:211
								__e = __c.Set("type", "text"); if __e != nil { return }
//line settings.gox:211
								__e = __c.Set("maxlength", "40"); if __e != nil { return }
//line settings.gox:211
								__e = __c.Set("value", prog.Name); if __e != nil { return }
//line settings.gox:211
								__e = __c.Set("placeholder", i18n.T(l, "settings.version_name")); if __e != nil { return }
//line settings.gox:212
								__e = __c.Modify(doors.AChange{On: func(ctx context.Context, r doors.RequestEvent[doors.ChangeEvent]) bool {
								nameSrc.Update(ctx, r.Event().Value)
								return false
							}}); if __e != nil { return }
							}
							__e = __c.Submit(); if __e != nil { return }
//line settings.gox:216
							__e = (doors.AClick{On: s.saveEditor(editSrc, nameSrc)}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
								ctx := __c.Context(); _ = ctx
								__e = __c.Init("button"); if __e != nil { return }
								{
//line settings.gox:216
									__e = __c.Set("class", "btn primary small"); if __e != nil { return }
//line settings.gox:216
									__e = __c.Set("type", "button"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:216
									__e = __c.Any(i18n.T(l, "settings.save_version")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
							return })); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
//line settings.gox:219
						if hasCustom {
//line settings.gox:220
							__e = (doors.AClick{On: s.deleteProgram}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
								ctx := __c.Context(); _ = ctx
								__e = __c.Init("button"); if __e != nil { return }
								{
//line settings.gox:220
									__e = __c.Set("class", "btn danger small"); if __e != nil { return }
//line settings.gox:220
									__e = __c.Set("type", "button"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:220
									__e = __c.Any(i18n.T(l, "settings.delete_version")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
							return })); if __e != nil { return }
						}
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("details"); if __e != nil { return }
					{
//line settings.gox:226
						__e = __c.Set("class", "cycle-details"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("summary"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
//line settings.gox:227
							__e = __c.Any(i18n.T(l, "settings.exercise_content")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:228
							__e = __c.Set("class", "srow"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Init("span"); if __e != nil { return }
							{
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:229
								__e = __c.Any(i18n.T(l, "settings.exercise")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
							__e = __c.Init("select"); if __e != nil { return }
							{
//line settings.gox:230
								__e = __c.Set("class", "lang-select compact"); if __e != nil { return }
//line settings.gox:230
								__e = __c.Modify(doors.AChange{On: func(ctx context.Context, r doors.RequestEvent[doors.ChangeEvent]) bool {
							selSrc.Update(ctx, r.Event().Value)
							return false
						}}); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:234
								__e = __c.Any(s.exOptions(entries, l, firstID)); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
//line settings.gox:238
						__e = __c.Any(selSrc.Bind(func(id string) gox.Elem {
						return gox.Elem(func(__c gox.Cursor) (__e error) {
							ctx := __c.Context(); _ = ctx
//line settings.gox:239
							__e = __c.Any(s.exDetail(l, cat, id)); if __e != nil { return }
						return })
//line settings.gox:240
					})); if __e != nil { return }
						__e = __c.Init("div"); if __e != nil { return }
						{
//line settings.gox:241
							__e = __c.Set("class", "srow"); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Init("span"); if __e != nil { return }
							{
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:242
								__e = __c.Any(i18n.T(l, "settings.download")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
//line settings.gox:243
							__e = __c.Any(selSrc.Bind(func(id string) gox.Elem {
							return gox.Elem(func(__c gox.Cursor) (__e error) {
								ctx := __c.Context(); _ = ctx
//line settings.gox:245
								doc, _ := cat.Doc(id)
								b := doc.Marshal()

								__e = __c.Init("a"); if __e != nil { return }
								{
//line settings.gox:248
									__e = __c.Set("class", "btn small"); if __e != nil { return }
//line settings.gox:248
									__e = __c.Modify(doors.ResourceBytes(b)); if __e != nil { return }
//line settings.gox:248
									__e = __c.Set("name", id + ".json"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:248
									__e = __c.Many(id, ".json"); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
							return })
//line settings.gox:249
						})); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
//line settings.gox:251
						__e = (doors.ARawSubmit{On: s.importExercise}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
							ctx := __c.Context(); _ = ctx
							__e = __c.Init("form"); if __e != nil { return }
							{
//line settings.gox:251
								__e = __c.Set("class", "ex-import-form"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
								__e = __c.Init("textarea"); if __e != nil { return }
								{
//line settings.gox:252
									__e = __c.Set("class", "ex-import"); if __e != nil { return }
//line settings.gox:252
									__e = __c.Set("name", "exjson"); if __e != nil { return }
//line settings.gox:252
									__e = __c.Set("rows", "3"); if __e != nil { return }
//line settings.gox:252
									__e = __c.Set("placeholder", i18n.T(l, "settings.paste_json")); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
								__e = __c.Init("div"); if __e != nil { return }
								{
//line settings.gox:253
									__e = __c.Set("class", "srow ex-import-row"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
									__e = __c.InitVoid("input"); if __e != nil { return }
									{
//line settings.gox:254
										__e = __c.Set("class", "ex-file"); if __e != nil { return }
//line settings.gox:254
										__e = __c.Set("type", "file"); if __e != nil { return }
//line settings.gox:254
										__e = __c.Set("accept", ".json,application/json"); if __e != nil { return }
									}
									__e = __c.Submit(); if __e != nil { return }
									__e = __c.Init("button"); if __e != nil { return }
									{
//line settings.gox:255
										__e = __c.Set("class", "btn primary small"); if __e != nil { return }
//line settings.gox:255
										__e = __c.Set("type", "submit"); if __e != nil { return }
										__e = __c.Submit(); if __e != nil { return }
//line settings.gox:255
										__e = __c.Any(i18n.T(l, "settings.import")); if __e != nil { return }
									}
									__e = __c.Close(); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						return })); if __e != nil { return }
						__e = __c.Init("script"); if __e != nil { return }
						{
//line settings.gox:258
							__e = __c.Set("src", assets.ExUpload); if __e != nil { return }
							__e = __c.Set("inline", true); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Raw(""); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line settings.gox:261
						__e = __c.Set("class", "srow"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Init("span"); if __e != nil { return }
						{
							__e = __c.Submit(); if __e != nil { return }
//line settings.gox:262
							__e = __c.Any(i18n.T(l, "settings.account")); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
//line settings.gox:263
						if s.sess.IsGuest {
//line settings.gox:264
							__e = (doors.ALink{Model: path.Path{Page: path.Register}}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
								ctx := __c.Context(); _ = ctx
								__e = __c.Init("a"); if __e != nil { return }
								{
//line settings.gox:264
									__e = __c.Set("class", "btn primary small"); if __e != nil { return }
									__e = __c.Submit(); if __e != nil { return }
//line settings.gox:264
									__e = __c.Any(i18n.T(l, "auth.register")); if __e != nil { return }
								}
								__e = __c.Close(); if __e != nil { return }
							return })); if __e != nil { return }
						} else  {
							__e = __c.Init("span"); if __e != nil { return }
							{
//line settings.gox:266
								__e = __c.Set("class", "email"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:266
								__e = __c.Any(s.sess.Email); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						}
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("div"); if __e != nil { return }
					{
//line settings.gox:270
						__e = __c.Set("class", "srow save-row"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:271
						__e = (doors.AClick{On: s.reset}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
							ctx := __c.Context(); _ = ctx
							__e = __c.Init("button"); if __e != nil { return }
							{
//line settings.gox:271
								__e = __c.Set("class", "btn danger small"); if __e != nil { return }
//line settings.gox:271
								__e = __c.Set("type", "button"); if __e != nil { return }
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:271
								__e = __c.Any(i18n.T(l, "settings.reset")); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						return })); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line settings.gox:276
}

// exDetail renders the read-only content card for an exercise: description,
// how-to, correct, mistakes, breathing and safety, localized.
//line settings.gox:280
func (s settingsPage) exDetail(l i18n.Lang, cat catalog.Catalog, id string) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line settings.gox:281
		d, _ := cat.Detail(l, id)
		__e = __c.Init("div"); if __e != nil { return }
		{
//line settings.gox:282
			__e = __c.Set("class", "ex-detail"); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
			__e = __c.Init("div"); if __e != nil { return }
			{
//line settings.gox:283
				__e = __c.Set("class", "ex-detail-name"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
//line settings.gox:283
				__e = __c.Any(cat.Name(l, id)); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
//line settings.gox:284
			if d.Desc != "" {
				__e = __c.Init("p"); if __e != nil { return }
				{
//line settings.gox:285
					__e = __c.Set("class", "ex-detail-desc"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
//line settings.gox:285
					__e = __c.Any(d.Desc); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
//line settings.gox:287
			if len(d.HowTo) > 0 {
				__e = __c.Init("div"); if __e != nil { return }
				{
//line settings.gox:288
					__e = __c.Set("class", "ex-detail-sec"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("h4"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:289
						__e = __c.Any(i18n.T(l, "info.how")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("ol"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:290
						for _, x := range d.HowTo {
							__e = __c.Init("li"); if __e != nil { return }
							{
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:290
								__e = __c.Any(x); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						}
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
//line settings.gox:293
			if len(d.Correct) > 0 {
				__e = __c.Init("div"); if __e != nil { return }
				{
//line settings.gox:294
					__e = __c.Set("class", "ex-detail-sec"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("h4"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:295
						__e = __c.Any(i18n.T(l, "info.correct")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("ul"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:296
						for _, x := range d.Correct {
							__e = __c.Init("li"); if __e != nil { return }
							{
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:296
								__e = __c.Any(x); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						}
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
//line settings.gox:299
			if len(d.Wrong) > 0 {
				__e = __c.Init("div"); if __e != nil { return }
				{
//line settings.gox:300
					__e = __c.Set("class", "ex-detail-sec"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("h4"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:301
						__e = __c.Any(i18n.T(l, "info.wrong")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("ul"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:302
						for _, x := range d.Wrong {
							__e = __c.Init("li"); if __e != nil { return }
							{
								__e = __c.Submit(); if __e != nil { return }
//line settings.gox:302
								__e = __c.Any(x); if __e != nil { return }
							}
							__e = __c.Close(); if __e != nil { return }
						}
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
//line settings.gox:305
			if d.Breathing != "" {
				__e = __c.Init("div"); if __e != nil { return }
				{
//line settings.gox:306
					__e = __c.Set("class", "ex-detail-sec"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("h4"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:307
						__e = __c.Any(i18n.T(l, "info.breathing")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("p"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:308
						__e = __c.Any(d.Breathing); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
//line settings.gox:311
			if d.Safety != "" {
				__e = __c.Init("div"); if __e != nil { return }
				{
//line settings.gox:312
					__e = __c.Set("class", "ex-detail-sec"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("h4"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
						__e = __c.Text("⚠️ "); if __e != nil { return }
//line settings.gox:313
						__e = __c.Any(i18n.T(l, "info.warning")); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("p"); if __e != nil { return }
					{
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:314
						__e = __c.Any(d.Safety); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line settings.gox:318
}

// exOptions renders the exercise dropdown options grouped by slot.
//line settings.gox:321
func (s settingsPage) exOptions(entries []catalog.Entry, l i18n.Lang, cur string) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line settings.gox:322
		__e = __c.Any(s.exGroup(entries, l, cur, "warmup", "workout.warmup")); if __e != nil { return }
//line settings.gox:323
		__e = __c.Any(s.exGroup(entries, l, cur, "main", "workout.main")); if __e != nil { return }
//line settings.gox:324
		__e = __c.Any(s.exGroup(entries, l, cur, "cooldown", "workout.cooldown")); if __e != nil { return }
	return })
//line settings.gox:325
}

//line settings.gox:327
func (s settingsPage) exGroup(entries []catalog.Entry, l i18n.Lang, cur, slot, labelKey string) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
		__e = __c.Init("optgroup"); if __e != nil { return }
		{
//line settings.gox:328
			__e = __c.Set("label", i18n.T(l, labelKey)); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
//line settings.gox:329
			for _, e := range entries {
//line settings.gox:330
				if e.Slot == slot {
					__e = __c.Init("option"); if __e != nil { return }
					{
//line settings.gox:331
						__e = __c.Set("value", e.ID); if __e != nil { return }
//line settings.gox:331
						__e = __c.Set("selected", func() any { return e.ID == cur }()); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:331
						__e = __c.Any(e.Name); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
			}
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line settings.gox:335
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
		mode := f.VoiceMode
		switch mode {
		case "off", "min", "normal", "detailed":
		default:
			mode = "normal"
		}
		_ = app.DB.UpdateSettings(s.sess.UserID, name, s.sess.Lang, rest, mode)
		_ = app.DB.SetSex(s.sess.UserID, f.Sex)
		s.auth.Mutate(ctx, func(cur auth.Session) auth.Session {
			cur.Name = name
			cur.Rest = rest
			cur.VoiceMode = mode
			cur.Voice = mode != "off"
			return cur
		})
		saved.Update(ctx, true)
		return false
	}
}

// unitWord is the unit label shown next to an editable value.
func unitWord(l i18n.Lang, u content.Unit, perSide bool) string {
	var w string
	switch u {
	case content.Seconds:
		w = i18n.T(l, "workout.sec")
	case content.Breaths:
		w = i18n.T(l, "workout.breaths")
	default:
		w = i18n.T(l, "workout.x_reps")
	}
	if perSide {
		w += " /" + i18n.T(l, "workout.side")
	}
	return w
}

// clampDay clamps a day number to [1, n] (0 when n == 0).
func clampDay(d, n int) int {
	if n == 0 {
		return 0
	}
	if d < 1 {
		d = 1
	}
	if d > n {
		d = n
	}
	return d
}

// setDay switches the edited day, keeping the working JSON.
func (s settingsPage) setDay(editSrc doors.Source[editState]) func(context.Context, doors.RequestEvent[doors.ChangeEvent]) bool {
	return func(ctx context.Context, r doors.RequestEvent[doors.ChangeEvent]) bool {
		d, err := strconv.Atoi(r.Event().Value)
		if err != nil {
			return false
		}
		st := editSrc.Get()
		st.Day = d
		editSrc.Update(ctx, st)
		return false
	}
}

// slotForTarget maps an editor section to the slot its items belong to.
func slotForTarget(target string) content.Slot {
	switch target {
	case "warmup":
		return content.Warmup
	case "cooldown":
		return content.Cooldown
	default:
		return content.Main
	}
}

// mutateList parses the working-copy JSON, applies fn to the target list — the
// shared warm-up / cool-down (whole series) or the selected day's main block —
// and stores it back (triggering a re-render).
func mutateList(ctx context.Context, editSrc doors.Source[editState], target string, fn func(items []content.RItem) []content.RItem) {
	st := editSrc.Get()
	p, err := content.ParseResolved(st.JSON)
	if err != nil {
		return
	}
	switch target {
	case "warmup":
		p.Warmup = fn(p.Warmup)
	case "cooldown":
		p.Cooldown = fn(p.Cooldown)
	default:
		d := clampDay(st.Day, len(p.Days))
		if d < 1 {
			return
		}
		p.Days[d-1].Items = fn(p.Days[d-1].Items)
	}
	st.JSON = p.Marshal()
	editSrc.Update(ctx, st)
}

// editList renders one CRUD section (warm-up / day main / cool-down).
//line settings.gox:455
func (s settingsPage) editList(editSrc doors.Source[editState], cat catalog.Catalog, entries []catalog.Entry, l i18n.Lang, items []content.RItem, target, open string) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
		__e = __c.Init("ol"); if __e != nil { return }
		{
//line settings.gox:456
			__e = __c.Set("class", "cyc-list edit"); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
//line settings.gox:457
			for ii := range items {
//line settings.gox:459
				it := items[ii]
				key := target + ":" + strconv.Itoa(ii)

				__e = __c.Init("li"); if __e != nil { return }
				{
//line settings.gox:462
					__e = __c.Set("class", "cyc-row"); if __e != nil { return }
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("span"); if __e != nil { return }
					{
//line settings.gox:463
						__e = __c.Set("class", "cyc-num"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.Init("select"); if __e != nil { return }
					{
//line settings.gox:464
						__e = __c.Set("class", "cyc-select"); if __e != nil { return }
//line settings.gox:464
						__e = __c.Modify(doors.AChange{On: s.editReplace(editSrc, cat, target, ii)}); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:465
						__e = __c.Any(s.exOptions(entries, l, it.ID)); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
					__e = __c.InitVoid("input"); if __e != nil { return }
					{
//line settings.gox:467
						__e = __c.Set("class", "cyc-val"); if __e != nil { return }
//line settings.gox:467
						__e = __c.Set("type", "number"); if __e != nil { return }
//line settings.gox:467
						__e = __c.Set("inputmode", "numeric"); if __e != nil { return }
//line settings.gox:467
						__e = __c.Set("min", "1"); if __e != nil { return }
//line settings.gox:467
						__e = __c.Set("max", "600"); if __e != nil { return }
//line settings.gox:467
						__e = __c.Set("value", it.Value); if __e != nil { return }
//line settings.gox:467
						__e = __c.Modify(doors.AChange{On: s.editSetValue(editSrc, target, ii)}); if __e != nil { return }
					}
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Init("span"); if __e != nil { return }
					{
//line settings.gox:468
						__e = __c.Set("class", "cyc-unit"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:468
						__e = __c.Any(unitWord(l, it.Unit, false)); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
//line settings.gox:469
					__e = (doors.AClick{On: s.editPeek(editSrc, target, ii)}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
						ctx := __c.Context(); _ = ctx
						__e = __c.Init("button"); if __e != nil { return }
						{
//line settings.gox:469
							__e = __c.Set("class", "cyc-op cyc-info"); if __e != nil { return }
//line settings.gox:469
							__e = __c.Set("type", "button"); if __e != nil { return }
//line settings.gox:469
							__e = __c.Set("aria-label", i18n.T(l, "settings.view")); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Text("ⓘ"); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					return })); if __e != nil { return }
//line settings.gox:470
					__e = (doors.AClick{On: s.editMove(editSrc, target, ii, -1)}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
						ctx := __c.Context(); _ = ctx
						__e = __c.Init("button"); if __e != nil { return }
						{
//line settings.gox:470
							__e = __c.Set("class", "cyc-op"); if __e != nil { return }
//line settings.gox:470
							__e = __c.Set("type", "button"); if __e != nil { return }
//line settings.gox:470
							__e = __c.Set("aria-label", i18n.T(l, "settings.move_up")); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Text("↑"); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					return })); if __e != nil { return }
//line settings.gox:471
					__e = (doors.AClick{On: s.editMove(editSrc, target, ii, 1)}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
						ctx := __c.Context(); _ = ctx
						__e = __c.Init("button"); if __e != nil { return }
						{
//line settings.gox:471
							__e = __c.Set("class", "cyc-op"); if __e != nil { return }
//line settings.gox:471
							__e = __c.Set("type", "button"); if __e != nil { return }
//line settings.gox:471
							__e = __c.Set("aria-label", i18n.T(l, "settings.move_down")); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Text("↓"); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					return })); if __e != nil { return }
//line settings.gox:472
					__e = (doors.AClick{On: s.editDelete(editSrc, target, ii)}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
						ctx := __c.Context(); _ = ctx
						__e = __c.Init("button"); if __e != nil { return }
						{
//line settings.gox:472
							__e = __c.Set("class", "cyc-op danger"); if __e != nil { return }
//line settings.gox:472
							__e = __c.Set("type", "button"); if __e != nil { return }
//line settings.gox:472
							__e = __c.Set("aria-label", i18n.T(l, "settings.remove")); if __e != nil { return }
							__e = __c.Submit(); if __e != nil { return }
							__e = __c.Text("✕"); if __e != nil { return }
						}
						__e = __c.Close(); if __e != nil { return }
					return })); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
//line settings.gox:474
				if open == key {
					__e = __c.Init("li"); if __e != nil { return }
					{
//line settings.gox:475
						__e = __c.Set("class", "cyc-peek"); if __e != nil { return }
						__e = __c.Submit(); if __e != nil { return }
//line settings.gox:475
						__e = __c.Any(s.exDetail(l, cat, it.ID)); if __e != nil { return }
					}
					__e = __c.Close(); if __e != nil { return }
				}
			}
		}
		__e = __c.Close(); if __e != nil { return }
//line settings.gox:479
		__e = (doors.AClick{On: s.editAdd(editSrc, cat, target)}).Proxy(__c, gox.Elem(func(__c gox.Cursor) (__e error) {
			ctx := __c.Context(); _ = ctx
			__e = __c.Init("button"); if __e != nil { return }
			{
//line settings.gox:479
				__e = __c.Set("class", "btn small cyc-add"); if __e != nil { return }
//line settings.gox:479
				__e = __c.Set("type", "button"); if __e != nil { return }
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Text("＋ "); if __e != nil { return }
//line settings.gox:479
				__e = __c.Any(i18n.T(l, "settings.add_exercise")); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
		return })); if __e != nil { return }
	return })
//line settings.gox:480
}

// levelLabel formats the difficulty level as a percent (or the "base" word).
func levelLabel(l i18n.Lang, level int) string {
	if level == 0 {
		return i18n.T(l, "settings.level_base")
	}
	if level > 0 {
		return "+" + strconv.Itoa(level*10) + "%"
	}
	return strconv.Itoa(level*10) + "%"
}

// setLevel persists the universal difficulty level immediately and updates the session.
func (s settingsPage) setLevel(levelSrc doors.Source[int]) func(context.Context, doors.RequestEvent[doors.ChangeEvent]) bool {
	return func(ctx context.Context, r doors.RequestEvent[doors.ChangeEvent]) bool {
		n := r.Event().Number
		if n == nil {
			return false
		}
		lv := int(*n)
		if lv < -3 {
			lv = -3
		}
		if lv > 3 {
			lv = 3
		}
		_ = app.DB.SetLevel(s.sess.UserID, lv)
		s.auth.Mutate(ctx, func(cur auth.Session) auth.Session { cur.Level = lv; return cur })
		levelSrc.Update(ctx, lv)
		return false
	}
}

// editReplace swaps the exercise at index ii for the one picked in the dropdown,
// resetting unit/side/slot from the chosen exercise and converting the value.
func (s settingsPage) editReplace(editSrc doors.Source[editState], cat catalog.Catalog, target string, ii int) func(context.Context, doors.RequestEvent[doors.ChangeEvent]) bool {
	return func(ctx context.Context, r doors.RequestEvent[doors.ChangeEvent]) bool {
		newID := r.Event().Value
		def, ok := cat.Def(newID)
		if !ok {
			return false
		}
		mutateList(ctx, editSrc, target, func(its []content.RItem) []content.RItem {
			if ii < 0 || ii >= len(its) {
				return its
			}
			old := its[ii]
			its[ii] = content.RItem{
				ID: newID, Unit: def.Unit, PerSide: def.PerSide, Slot: slotForTarget(target),
				Value: content.ConvertValue(old.Unit, def.Unit, old.Value), Round: old.Round,
			}
			return its
		})
		return false
	}
}

// editSetValue updates the value (reps/seconds) at index ii.
func (s settingsPage) editSetValue(editSrc doors.Source[editState], target string, ii int) func(context.Context, doors.RequestEvent[doors.ChangeEvent]) bool {
	return func(ctx context.Context, r doors.RequestEvent[doors.ChangeEvent]) bool {
		n := r.Event().Number
		if n == nil {
			return false
		}
		v := int(*n)
		if v < 1 {
			v = 1
		}
		if v > 600 {
			v = 600
		}
		mutateList(ctx, editSrc, target, func(its []content.RItem) []content.RItem {
			if ii >= 0 && ii < len(its) {
				its[ii].Value = v
			}
			return its
		})
		return false
	}
}

// editMove shifts the item at index ii by delta (-1 up / +1 down).
func (s settingsPage) editMove(editSrc doors.Source[editState], target string, ii, delta int) func(context.Context, doors.RequestPointer) bool {
	return func(ctx context.Context, r doors.RequestPointer) bool {
		mutateList(ctx, editSrc, target, func(its []content.RItem) []content.RItem {
			j := ii + delta
			if ii < 0 || ii >= len(its) || j < 0 || j >= len(its) {
				return its
			}
			its[ii], its[j] = its[j], its[ii]
			return its
		})
		return false
	}
}

// editDelete removes the item at index ii.
func (s settingsPage) editDelete(editSrc doors.Source[editState], target string, ii int) func(context.Context, doors.RequestPointer) bool {
	return func(ctx context.Context, r doors.RequestPointer) bool {
		mutateList(ctx, editSrc, target, func(its []content.RItem) []content.RItem {
			if ii < 0 || ii >= len(its) {
				return its
			}
			nw := make([]content.RItem, 0, len(its)-1)
			nw = append(nw, its[:ii]...)
			nw = append(nw, its[ii+1:]...)
			return nw
		})
		return false
	}
}

// editPeek toggles the inline content preview for a row.
func (s settingsPage) editPeek(editSrc doors.Source[editState], target string, ii int) func(context.Context, doors.RequestPointer) bool {
	key := target + ":" + strconv.Itoa(ii)
	return func(ctx context.Context, r doors.RequestPointer) bool {
		st := editSrc.Get()
		if st.Open == key {
			st.Open = ""
		} else {
			st.Open = key
		}
		editSrc.Update(ctx, st)
		return false
	}
}

// editAdd appends a default exercise (by section) to the target list.
func (s settingsPage) editAdd(editSrc doors.Source[editState], cat catalog.Catalog, target string) func(context.Context, doors.RequestPointer) bool {
	return func(ctx context.Context, r doors.RequestPointer) bool {
		it := defaultItem(cat, target)
		mutateList(ctx, editSrc, target, func(its []content.RItem) []content.RItem {
			return append(append([]content.RItem{}, its...), it)
		})
		return false
	}
}

// defaultItem is the exercise added by a section's "add" button.
func defaultItem(cat catalog.Catalog, target string) content.RItem {
	id := "M01"
	switch target {
	case "warmup":
		id = "W03"
	case "cooldown":
		id = "CD02"
	}
	unit, perSide := content.Reps, false
	if def, ok := cat.Def(id); ok {
		unit, perSide = def.Unit, def.PerSide
	}
	val := 10
	if unit == content.Seconds {
		val = 30
	} else if unit == content.Breaths {
		val = 5
	}
	return content.RItem{ID: id, Unit: unit, Value: val, PerSide: perSide, Slot: slotForTarget(target), Round: 1}
}

// saveEditor stores the edited program as the user's custom version and reloads.
func (s settingsPage) saveEditor(editSrc doors.Source[editState], nameSrc doors.Source[string]) func(context.Context, doors.RequestPointer) bool {
	return func(ctx context.Context, r doors.RequestPointer) bool {
		p, err := content.ParseResolved(editSrc.Get().JSON)
		if err != nil {
			return false
		}
		name := strings.TrimSpace(nameSrc.Get())
		if name == "" {
			name = "My program"
		}
		if len(name) > 40 {
			name = name[:40]
		}
		p.Name = name
		_ = app.DB.SetProgramJSON(s.sess.UserID, p.Marshal())
		_ = app.DB.SetActiveProgram(s.sess.UserID, "custom")
		_ = r.After(doors.ActionLocationReload{})
		return false
	}
}

// importExercise stores an uploaded/pasted per-exercise JSON into the user's
// custom library (it then appears in the dropdowns and is shown/spoken), reloads.
func (s settingsPage) importExercise(ctx context.Context, r doors.RequestRawForm) bool {
	pf, err := r.ParseForm(1 << 20)
	if err != nil {
		return false
	}
	raw := strings.TrimSpace(pf.FormValue("exjson"))
	if raw == "" {
		return false
	}
	var d catalog.Doc
	if json.Unmarshal([]byte(raw), &d) != nil {
		return false
	}
	if !validDoc(d) {
		return false
	}
	_ = app.SaveCustomExercise(s.sess.UserID, d)
	_ = r.After(doors.ActionLocationReload{})
	return false
}

// validDoc sanity-checks an imported exercise document.
func validDoc(d catalog.Doc) bool {
	if strings.TrimSpace(d.ID) == "" {
		return false
	}
	switch content.Slot(d.Slot) {
	case content.Warmup, content.Main, content.Cooldown:
	default:
		return false
	}
	switch content.Unit(d.Unit) {
	case content.Reps, content.Seconds, content.Breaths:
	default:
		return false
	}
	return true
}

// selectProgram switches the active program (a standard key or "custom") and reloads.
func (s settingsPage) selectProgram(ctx context.Context, r doors.RequestEvent[doors.ChangeEvent]) bool {
	app.DB.SetActiveProgram(s.sess.UserID, r.Event().Value)
	_ = r.After(doors.ActionLocationReload{})
	return false
}

// deleteProgram removes the custom program and reverts to the default standard
// program (which can never be deleted), then reloads.
func (s settingsPage) deleteProgram(ctx context.Context, r doors.RequestPointer) bool {
	_ = app.DB.SetProgramJSON(s.sess.UserID, "")
	_ = app.DB.SetActiveProgram(s.sess.UserID, content.DefaultStandardKey())
	_ = r.After(doors.ActionLocationReload{})
	return false
}

// reset clears all progress and returns to the home screen.
func (s settingsPage) reset(ctx context.Context, r doors.RequestPointer) bool {
	_ = app.DB.ResetProgress(s.sess.UserID)
	s.path.Update(ctx, path.Path{Page: path.Home})
	return false
}
