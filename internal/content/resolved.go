package content

import "encoding/json"

// Resolved is the editable, flat program: a shared warm-up and cool-down series
// (the same for every day) plus a per-day MAIN block. This is what the editor
// edits and what Workout builds from. Warm-up/cool-down are stored once so
// editing them applies to the whole series; only the main block is per-day.
type Resolved struct {
	Version      int     `json:"version"`
	Name         string  `json:"name"`
	TotalDays    int     `json:"totalDays"`
	WarmupRounds int     `json:"warmupRounds"`
	Warmup       []RItem `json:"warmup"`   // shared, one round (repeated WarmupRounds times)
	Days         []RDay  `json:"days"`     // per-day MAIN items
	Cooldown     []RItem `json:"cooldown"` // shared, one round
}

// RDay is one day's MAIN block.
type RDay struct {
	Day   int     `json:"day"`
	Items []RItem `json:"items"`
}

// RItem is one exercise occurrence with its value.
type RItem struct {
	ID      string `json:"id"`
	Unit    Unit   `json:"unit"`
	Value   int    `json:"value"`
	PerSide bool   `json:"perSide"`
	Slot    Slot   `json:"slot"`
	Round   int    `json:"round"`
}

// ResolveBuiltin flattens the built-in program: shared warm-up/cool-down series
// plus each day's main block (with its baked-in day-to-day progression).
func ResolveBuiltin() Resolved {
	r := Resolved{Version: 1, Name: "Sergey", TotalDays: TotalDays, WarmupRounds: warmupRounds}
	r.Warmup = seqItems(warmupSeq, Warmup)
	r.Cooldown = seqItems(cooldownSeq, Cooldown)
	for day := 1; day <= TotalDays; day++ {
		r.Days = append(r.Days, RDay{Day: day, Items: mainItems(day)})
	}
	return r
}

// ParseResolved parses a stored custom program JSON and normalizes it (defaults +
// migration of legacy full-day programs into shared warm-up/cool-down).
func ParseResolved(data string) (Resolved, error) {
	var r Resolved
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return r, err
	}
	r.normalize()
	return r, nil
}

// normalize fills defaults and migrates legacy programs (where warm-up/cool-down
// lived inside each day's Items) into the shared warm-up/cool-down model.
func (r *Resolved) normalize() {
	if r.WarmupRounds < 1 {
		r.WarmupRounds = warmupRounds
	}
	if len(r.Warmup) == 0 && len(r.Cooldown) == 0 && len(r.Days) > 0 {
		// Legacy: lift the shared warm-up (round 1) and cool-down out of day 1,
		// then strip warm-up/cool-down from every day, leaving main only.
		for _, it := range r.Days[0].Items {
			switch {
			case it.Slot == Warmup && it.Round <= 1:
				r.Warmup = append(r.Warmup, it)
			case it.Slot == Cooldown:
				r.Cooldown = append(r.Cooldown, it)
			}
		}
		if len(r.Warmup) > 0 || len(r.Cooldown) > 0 {
			for di := range r.Days {
				main := make([]RItem, 0, len(r.Days[di].Items))
				for _, it := range r.Days[di].Items {
					if it.Slot == Main {
						main = append(main, it)
					}
				}
				r.Days[di].Items = main
			}
		}
	}
}

// Marshal serializes the program to JSON.
func (r Resolved) Marshal() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// Workout builds a session for a day: the shared warm-up (× WarmupRounds), the
// day's main block scaled by the difficulty level, then the shared cool-down.
// Warm-up/cool-down and breaths are never scaled.
func (r Resolved) Workout(day, lightRest, level int) Workout {
	if len(r.Days) == 0 {
		return ResolveBuiltin().Workout(day, lightRest, level)
	}
	if day < 1 {
		day = 1
	}
	if day > len(r.Days) {
		day = len(r.Days)
	}
	rounds := r.WarmupRounds
	if rounds < 1 {
		rounds = 1
	}

	var items []Item
	emit := func(ri RItem, round int, scale bool) {
		unit, perSide := ri.Unit, ri.PerSide
		if e, ok := Pool[ri.ID]; ok { // built-in def wins for unit/side
			unit, perSide = e.Unit, e.PerSide
		}
		val := ri.Value
		if scale {
			val = ScaleValue(unit, val, level)
		}
		items = append(items, Item{
			ExerciseID: ri.ID, Unit: unit, Value: val,
			Slot: ri.Slot, PerSide: perSide, Round: round,
		})
	}

	for rnd := 1; rnd <= rounds; rnd++ {
		for _, ri := range r.Warmup {
			emit(ri, rnd, false)
		}
	}
	for _, ri := range r.Days[day-1].Items {
		emit(ri, ri.Round, true) // main — scaled by level
	}
	for _, ri := range r.Cooldown {
		emit(ri, 1, false)
	}

	est := assignRestsAndEst(items, lightRest)
	return Workout{Day: day, Items: items, EstSec: est}
}
