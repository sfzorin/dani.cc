// Package catalog is the per-user exercise library: it overlays a user's custom
// exercises (uploaded/edited per-exercise JSON) on top of the built-in pool and
// resolves an exercise's definition, localized name/content and voice. A custom
// exercise with the same ID overrides the built-in one. It is the single source
// the pages use so custom exercises are selectable, shown and spoken.
package catalog

import (
	"encoding/json"
	"sort"

	"morning30/internal/content"
	"morning30/internal/i18n"
)

// VoiceDoc are the non-blocking voice cues for one language.
type VoiceDoc struct {
	Start  string   `json:"start,omitempty"`
	Mid    []string `json:"mid,omitempty"`
	Last   string   `json:"last,omitempty"`
	Finish string   `json:"finish,omitempty"`
}

// Loc is the localized text for an exercise in one language.
type Loc struct {
	Name      string   `json:"name"`
	Hint      string   `json:"hint,omitempty"`
	Warning   string   `json:"warning,omitempty"`
	Desc      string   `json:"desc,omitempty"`
	HowTo     []string `json:"howTo,omitempty"`
	Correct   []string `json:"correct,omitempty"`
	Wrong     []string `json:"wrong,omitempty"`
	Breathing string   `json:"breathing,omitempty"`
	Safety    string   `json:"safety,omitempty"`
	Voice     VoiceDoc `json:"voice,omitempty"`
}

// Doc is a self-contained exercise: definition + localized content for every
// language. This is the unit of per-exercise JSON (download / upload).
type Doc struct {
	ID          string         `json:"id"`
	Slot        string         `json:"slot"`
	Unit        string         `json:"unit"`
	PerSide     bool           `json:"perSide"`
	KneeLoad    string         `json:"kneeLoad,omitempty"`
	Difficulty  int            `json:"difficulty"`
	Replacement string         `json:"replacement,omitempty"`
	SVG         string         `json:"svg,omitempty"`
	Content     map[string]Loc `json:"content"`
}

// loc returns the best localized entry: the requested language, then English,
// then any present, then a zero value.
func (d Doc) loc(l i18n.Lang) Loc {
	if x, ok := d.Content[string(l)]; ok {
		return x
	}
	if x, ok := d.Content["en"]; ok {
		return x
	}
	for _, x := range d.Content {
		return x
	}
	return Loc{}
}

// Marshal pretty-prints the doc (for download).
func (d Doc) Marshal() []byte {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return []byte("{}")
	}
	return b
}

// Compact serializes the doc densely (for DB storage).
func (d Doc) Compact() string {
	b, err := json.Marshal(d)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// Builtin assembles the canonical Doc for a built-in exercise: its definition
// plus name/hint/warning/technique/voice in all 7 languages.
func Builtin(id string) (Doc, bool) {
	ex, ok := content.Get(id)
	if !ok {
		return Doc{}, false
	}
	d := Doc{
		ID: id, Slot: string(ex.Slot), Unit: string(ex.Unit), PerSide: ex.PerSide,
		KneeLoad: string(ex.KneeLoad), Difficulty: ex.Difficulty, Replacement: ex.Replacement,
		SVG:     "/static/ex/" + id + ".png",
		Content: map[string]Loc{},
	}
	for _, m := range i18n.Languages {
		l := m.Code
		loc := Loc{
			Name:    i18n.Exercise(l, id),
			Hint:    i18n.Hint(l, id),
			Warning: i18n.Warning(l, id),
		}
		if det, ok := content.GetDetailL(string(l), id); ok {
			loc.Desc = det.Desc
			loc.HowTo = det.HowTo
			loc.Correct = det.Correct
			loc.Wrong = det.Wrong
			loc.Breathing = det.Breathing
			loc.Safety = det.Safety
			loc.Voice = VoiceDoc{Start: det.Voice.Start, Mid: det.Voice.Mid, Last: det.Voice.Last, Finish: det.Voice.Finish}
		}
		d.Content[string(l)] = loc
	}
	return d, true
}

// Catalog resolves an exercise through three layers, custom winning:
//
//	custom (per-user)  ->  base (global DB library)  ->  Builtin (code fallback)
//
// The code fallback keeps it working when the DB hasn't been seeded (e.g. tests).
type Catalog struct {
	base   map[string]Doc // global library loaded from the DB
	custom map[string]Doc // per-user overrides / additions
}

func parseDocs(j string) map[string]Doc {
	out := map[string]Doc{}
	if j == "" {
		return out
	}
	var m map[string]Doc
	if json.Unmarshal([]byte(j), &m) == nil {
		for k, v := range m {
			if v.ID == "" {
				v.ID = k
			}
			out[v.ID] = v
		}
	}
	return out
}

// New parses the user's stored custom-exercises JSON (a map id -> Doc; "" ok).
// The global base is empty, so built-ins resolve from code.
func New(j string) Catalog {
	return Catalog{base: map[string]Doc{}, custom: parseDocs(j)}
}

// NewLayered builds a catalog from the global library (id -> Doc JSON, from the
// DB) plus the user's custom-exercises JSON.
func NewLayered(base map[string]string, customJSON string) Catalog {
	c := Catalog{base: make(map[string]Doc, len(base)), custom: parseDocs(customJSON)}
	for id, j := range base {
		var d Doc
		if json.Unmarshal([]byte(j), &d) == nil {
			if d.ID == "" {
				d.ID = id
			}
			c.base[d.ID] = d
		}
	}
	return c
}

// resolve returns the effective Doc for an ID (custom -> base -> code).
func (c Catalog) resolve(id string) (Doc, bool) {
	if d, ok := c.custom[id]; ok {
		return d, true
	}
	if d, ok := c.base[id]; ok {
		return d, true
	}
	return Builtin(id)
}

// Put adds or overrides a custom exercise.
func (c *Catalog) Put(d Doc) {
	if c.custom == nil {
		c.custom = map[string]Doc{}
	}
	c.custom[d.ID] = d
}

// Marshal serializes the custom exercises back to storage ("" when none).
func (c Catalog) Marshal() string {
	if len(c.custom) == 0 {
		return ""
	}
	b, err := json.Marshal(c.custom)
	if err != nil {
		return ""
	}
	return string(b)
}

// IsCustom reports whether the ID resolves to a user custom exercise.
func (c Catalog) IsCustom(id string) bool {
	_, ok := c.custom[id]
	return ok
}

// Def returns the exercise definition (custom -> base -> built-in).
func (c Catalog) Def(id string) (content.Exercise, bool) {
	d, ok := c.resolve(id)
	if !ok {
		return content.Exercise{}, false
	}
	return content.Exercise{
		ID: id, Slot: content.Slot(d.Slot), Unit: content.Unit(d.Unit),
		PerSide: d.PerSide, KneeLoad: content.KneeLoad(d.KneeLoad),
		Difficulty: d.Difficulty, Replacement: d.Replacement,
	}, true
}

// Name returns the localized exercise name (custom -> base -> built-in).
func (c Catalog) Name(l i18n.Lang, id string) string {
	d, ok := c.resolve(id)
	if !ok {
		return id
	}
	if n := d.loc(l).Name; n != "" {
		return n
	}
	return id
}

// Detail returns the localized technique card (custom -> base -> built-in).
func (c Catalog) Detail(l i18n.Lang, id string) (content.Detail, bool) {
	d, ok := c.resolve(id)
	if !ok {
		return content.Detail{}, false
	}
	x := d.loc(l)
	return content.Detail{
		Desc: x.Desc, HowTo: x.HowTo, Correct: x.Correct, Wrong: x.Wrong,
		Breathing: x.Breathing, Safety: x.Safety,
		Voice: content.Voice{Start: x.Voice.Start, Mid: x.Voice.Mid, Last: x.Voice.Last, Finish: x.Voice.Finish},
	}, true
}

// HasVoiceLang reports whether per-exercise voice exists for this exercise+lang.
func (c Catalog) HasVoiceLang(l i18n.Lang, id string) bool {
	d, ok := c.resolve(id)
	if !ok {
		return false
	}
	return d.loc(l).Voice.Start != ""
}

// Hint returns the short on-screen cue (custom -> base -> built-in).
func (c Catalog) Hint(l i18n.Lang, id string) string {
	if d, ok := c.resolve(id); ok {
		return d.loc(l).Hint
	}
	return ""
}

// Warning returns the inline safety warning (custom -> base -> built-in).
func (c Catalog) Warning(l i18n.Lang, id string) string {
	if d, ok := c.resolve(id); ok {
		return d.loc(l).Warning
	}
	return ""
}

// SVG returns the figure URL for an exercise (a static PNG under /static/ex).
func (c Catalog) SVG(id string) string {
	if d, ok := c.resolve(id); ok && d.SVG != "" {
		return d.SVG
	}
	return "/static/ex/" + id + ".png"
}

// Figure returns the figure URL for an exercise, choosing the male/female
// variant by sex ("m"/"f"). A custom (per-user) exercise keeps its own SVG; the
// built-in figures live under /static/ex/<sex>/<id>.png.
func (c Catalog) Figure(id, sex string) string {
	if c.IsCustom(id) {
		if d, ok := c.resolve(id); ok && d.SVG != "" {
			return d.SVG
		}
	}
	if sex != "f" {
		sex = "m"
	}
	return "/static/ex/" + sex + "/" + id + ".png"
}

// Doc returns the self-contained doc for an exercise (custom -> base -> built-in).
func (c Catalog) Doc(id string) (Doc, bool) {
	return c.resolve(id)
}

// Entry is one selectable exercise for the editor dropdown.
type Entry struct {
	ID   string
	Slot string
	Name string
}

// List returns every selectable exercise (built-in in canonical order, then
// custom-only ones), with the localized name, for the editor's dropdown.
func (c Catalog) List(l i18n.Lang) []Entry {
	out := make([]Entry, 0, len(content.Library())+len(c.base)+len(c.custom))
	seen := map[string]bool{}
	for _, ex := range content.Library() {
		out = append(out, Entry{ID: ex.ID, Slot: c.slotOf(ex.ID), Name: c.Name(l, ex.ID)})
		seen[ex.ID] = true
	}
	extra := make([]string, 0)
	add := func(id string) {
		if !seen[id] {
			extra = append(extra, id)
			seen[id] = true
		}
	}
	for id := range c.base {
		add(id)
	}
	for id := range c.custom {
		add(id)
	}
	sort.Strings(extra)
	for _, id := range extra {
		out = append(out, Entry{ID: id, Slot: c.slotOf(id), Name: c.Name(l, id)})
	}
	return out
}

func (c Catalog) slotOf(id string) string {
	if d, ok := c.Def(id); ok && d.Slot != "" {
		return string(d.Slot)
	}
	return string(content.Main)
}
