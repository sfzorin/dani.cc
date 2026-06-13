// Package content defines the exercise library and the 30-day program,
// implemented to the product spec: warm-up (W) -> main blocks (A–E / R /
// control) with a plank challenge -> cool-down (CD). Floor/mat only; no jumps,
// no equipment, no lunges/deep-squats/mountain-climbers (see forbidden list).
package content

// Unit is how an exercise is measured.
type Unit string

const (
	Reps    Unit = "reps"    // counted repetitions
	Seconds Unit = "seconds" // held/timed
	Breaths Unit = "breaths" // counted breaths (paced by the user)
)

// Slot is where an exercise belongs in a session.
type Slot string

const (
	Warmup   Slot = "warmup"
	Main     Slot = "main"
	Cooldown Slot = "cooldown"
)

// KneeLoad flags exercises that load the knee, for the (phase-3) adaptation engine.
type KneeLoad string

const (
	KneeNone   KneeLoad = "none"
	KneeLow    KneeLoad = "low"
	KneeMedium KneeLoad = "medium"
)

// Exercise is one movement in the library. Magnitudes are NOT stored here — they
// come from the block tables per day. The ID doubles as the SVG asset name and
// the i18n key.
type Exercise struct {
	ID          string
	Slot        Slot
	Unit        Unit
	PerSide     bool
	KneeLoad    KneeLoad
	Difficulty  int    // 1–5
	Replacement string // substitute ID used by the adaptation engine
	Video       string // optional video/GIF URL (reserved; falls back to the SVG)
	Thumb       string // optional thumbnail URL (reserved)
}

// Pool is the library keyed by ID.
var Pool = func() map[string]Exercise {
	m := make(map[string]Exercise, len(library))
	for _, e := range library {
		m[e.ID] = e
	}
	return m
}()

// Get returns an exercise by ID.
func Get(id string) (Exercise, bool) {
	e, ok := Pool[id]
	return e, ok
}

var library = []Exercise{
	// ---- Warm-up (fixed, 2 rounds) ----
	{ID: "W01", Slot: Warmup, Unit: Breaths, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "W02", Slot: Warmup, Unit: Reps, PerSide: true, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "W03", Slot: Warmup, Unit: Reps, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "W04", Slot: Warmup, Unit: Reps, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "W05", Slot: Warmup, Unit: Seconds, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "W06", Slot: Warmup, Unit: Seconds, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "W07", Slot: Warmup, Unit: Reps, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "W08", Slot: Warmup, Unit: Reps, PerSide: true, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "W09", Slot: Warmup, Unit: Reps, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "W10", Slot: Warmup, Unit: Reps, Difficulty: 1, KneeLoad: KneeLow},

	// ---- Main library ----
	{ID: "M01", Slot: Main, Unit: Reps, Difficulty: 3, KneeLoad: KneeNone},
	{ID: "M02", Slot: Main, Unit: Reps, Difficulty: 4, KneeLoad: KneeNone, Replacement: "M01"},
	{ID: "M03", Slot: Main, Unit: Reps, Difficulty: 3, KneeLoad: KneeNone, Replacement: "M01"},

	{ID: "C01", Slot: Main, Unit: Seconds, Difficulty: 3, KneeLoad: KneeNone, Replacement: "C06"},
	{ID: "C02", Slot: Main, Unit: Seconds, Difficulty: 3, KneeLoad: KneeNone},
	{ID: "C03", Slot: Main, Unit: Seconds, Difficulty: 3, KneeLoad: KneeNone, Replacement: "C04"},
	{ID: "C04", Slot: Main, Unit: Reps, PerSide: true, Difficulty: 2, KneeLoad: KneeNone},
	{ID: "C05", Slot: Main, Unit: Seconds, PerSide: true, Difficulty: 3, KneeLoad: KneeNone},
	{ID: "C06", Slot: Main, Unit: Reps, PerSide: true, Difficulty: 2, KneeLoad: KneeNone},
	{ID: "C07", Slot: Main, Unit: Seconds, Difficulty: 3, KneeLoad: KneeNone, Replacement: "C04"},
	{ID: "C08", Slot: Main, Unit: Seconds, Difficulty: 4, KneeLoad: KneeNone, Replacement: "C02"},

	{ID: "B01", Slot: Main, Unit: Reps, Difficulty: 2, KneeLoad: KneeNone},
	{ID: "B02", Slot: Main, Unit: Reps, Difficulty: 2, KneeLoad: KneeNone},
	{ID: "B03", Slot: Main, Unit: Seconds, Difficulty: 2, KneeLoad: KneeNone},
	{ID: "B04", Slot: Main, Unit: Reps, Difficulty: 2, KneeLoad: KneeNone},
	{ID: "B05", Slot: Main, Unit: Reps, Difficulty: 2, KneeLoad: KneeNone},

	{ID: "G01", Slot: Main, Unit: Reps, Difficulty: 1, KneeLoad: KneeLow},
	{ID: "G02", Slot: Main, Unit: Reps, Difficulty: 2, KneeLoad: KneeLow},

	{ID: "L01", Slot: Main, Unit: Reps, Difficulty: 2, KneeLoad: KneeMedium, Replacement: "G01"},

	{ID: "A01", Slot: Main, Unit: Reps, PerSide: true, Difficulty: 2, KneeLoad: KneeNone},
	{ID: "A02", Slot: Main, Unit: Seconds, PerSide: true, Difficulty: 2, KneeLoad: KneeNone},

	{ID: "S01", Slot: Main, Unit: Seconds, Difficulty: 1, KneeLoad: KneeNone},

	// ---- Cool-down (fixed, 1 round) ----
	{ID: "CD01", Slot: Cooldown, Unit: Breaths, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "CD02", Slot: Cooldown, Unit: Seconds, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "CD03", Slot: Cooldown, Unit: Seconds, PerSide: true, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "CD04", Slot: Cooldown, Unit: Seconds, PerSide: true, Difficulty: 1, KneeLoad: KneeNone},
	{ID: "CD05", Slot: Cooldown, Unit: Seconds, PerSide: true, Difficulty: 1, KneeLoad: KneeLow},
	{ID: "CD06", Slot: Cooldown, Unit: Seconds, Difficulty: 1, KneeLoad: KneeNone},
}
