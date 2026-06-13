package content

// TotalDays is the cycle length.
const TotalDays = 30

// Item is one step the player walks the user through.
type Item struct {
	ExerciseID string `json:"id"`
	Unit       Unit   `json:"unit"`
	Value      int    `json:"value"`
	Slot       Slot   `json:"slot"`
	PerSide    bool   `json:"perSide"`
	Round      int    `json:"round"`
	Rest       int    `json:"rest"` // seconds of rest AFTER this item (0 on the last)
}

// Workout is a full session for one day.
type Workout struct {
	Day       int    `json:"day"`
	BlockCode string `json:"block"`  // display label, e.g. "A+", "R", "Control"
	Items     []Item `json:"items"`
	EstSec    int    `json:"estSec"`
}

// Feedback is the post-workout questionnaire (spec §14). Pain values are 0–10;
// Energy is "worse" | "normal" | "better".
type Feedback struct {
	Difficulty int    `json:"difficulty"` // 1–10
	Knee       int    `json:"knee"`
	Back       int    `json:"back"`
	Shoulder   int    `json:"shoulder"`
	Energy     string `json:"energy"`
}

// AdaptState is the persisted adaptation state that shapes the NEXT workout.
type AdaptState struct {
	Level     int  // variant delta for the next workout (−1/0/+1)
	ForceR    bool // make the next workout a recovery day
	Knee      int  // sessions remaining without the mini-squat
	Shoulder  int  // sessions remaining without narrow push-ups
	Back      int  // sessions remaining swapping scissors/flutter
	LowEnergy int  // consecutive "worse" energy reports
}

// Adapt is the per-build view derived from AdaptState.
type Adapt struct {
	LevelDelta    int
	ForceR        bool
	KneeBlock     bool // replace L01 (mini-squat) -> G01 (glute bridge)
	ShoulderBlock bool // replace M02 (narrow push-ups) -> M01 −30%
	BackBlock     bool // replace C01 -> C06, C03 -> C04
}

// SuppressIncreaseOnGap drops a pending level "+" when the streak is broken
// (the user skipped): repeat the day without increasing the load (spec §14).
func (a Adapt) SuppressIncreaseOnGap(streakAlive bool) Adapt {
	if !streakAlive && a.LevelDelta > 0 {
		a.LevelDelta = 0
	}
	return a
}

// ToAdapt projects the persisted state onto a single build.
func (s AdaptState) ToAdapt() Adapt {
	return Adapt{
		LevelDelta:    s.Level,
		ForceR:        s.ForceR,
		KneeBlock:     s.Knee > 0,
		ShoulderBlock: s.Shoulder > 0,
		BackBlock:     s.Back > 0,
	}
}

func max0(n int) int {
	if n < 0 {
		return 0
	}
	return n
}

// Decay advances the state by one completed session: block counters tick down
// and the one-shot level/force-R (already applied to that session) reset.
func (s AdaptState) Decay() AdaptState {
	s.Knee = max0(s.Knee - 1)
	s.Shoulder = max0(s.Shoulder - 1)
	s.Back = max0(s.Back - 1)
	s.Level = 0
	s.ForceR = false
	return s
}

// NextAdapt applies the spec §14 rules from the latest feedback to produce the
// state for the next workout. Call after Decay().
func NextAdapt(prev AdaptState, fb Feedback) AdaptState {
	next := prev
	noPain := fb.Knee < 3 && fb.Back < 3 && fb.Shoulder < 3
	switch {
	case fb.Difficulty >= 10:
		next.ForceR = true
		next.Level = 0
	case fb.Difficulty >= 8:
		next.Level = -1
	case fb.Difficulty >= 1 && fb.Difficulty <= 5 && noPain:
		next.Level = 1
	default: // 6–7, or easy but with pain → hold
		next.Level = 0
	}
	if fb.Knee >= 3 {
		next.Knee = 3
	}
	if fb.Shoulder >= 3 {
		next.Shoulder = 3
	}
	if fb.Back >= 3 {
		next.Back = 3
	}
	if fb.Energy == "worse" {
		next.LowEnergy = prev.LowEnergy + 1
		if next.LowEnergy >= 2 {
			next.ForceR = true
			next.LowEnergy = 0
		}
	} else {
		next.LowEnergy = 0
	}
	return next
}

// substitute applies the spec's replacement table for active pain blocks.
func substitute(id string, val int, a Adapt) (string, int) {
	switch {
	case a.KneeBlock && id == "L01":
		return "G01", val // mini-squat → glute bridge
	case a.ShoulderBlock && id == "M02":
		return "M01", int(float64(val)*0.7 + 0.5) // narrow → classic −30%
	case a.BackBlock && id == "C01":
		return "C06", 10 // scissors → toe taps (reps/side)
	case a.BackBlock && id == "C03":
		return "C04", 8 // flutter → dead bug (reps/side)
	}
	return id, val
}

// setDef is a fixed-value entry (warm-up, cool-down, recovery, control).
type setDef struct {
	ID  string
	Val int
}

// blockSet carries the three variant magnitudes [−, base, +].
type blockSet struct {
	ID   string
	Vals [3]int
}

// ---- Fixed sequences ----

var warmupSeq = []setDef{
	{"W01", 5}, {"W02", 5}, {"W03", 5}, {"W04", 5}, {"W05", 10},
	{"W06", 10}, {"W07", 5}, {"W08", 5}, {"W09", 5}, {"W10", 5},
}

var cooldownSeq = []setDef{
	{"CD01", 5}, {"CD02", 25}, {"CD03", 25}, {"CD04", 25}, {"CD05", 25}, {"CD06", 25},
}

// Recovery (R) block — 1 round, no long plank / scissors / flutter / narrow
// push-ups / mini-squat / RKC, per the spec.
var recoverySeq = []setDef{
	{"C04", 6}, {"G01", 12}, {"B01", 10}, {"S01", 15}, {"C05", 20},
}

// Day-30 control block — 2 rounds, then a 90 s plank challenge.
var controlSeq = []setDef{
	{"M01", 16}, {"C01", 30}, {"C03", 30}, {"C04", 10},
	{"B01", 15}, {"B02", 15}, {"L01", 10}, {"A01", 12},
}

const controlPlankSec = 90

// ---- Variant blocks (A–E), 2 rounds each ----

var blocks = map[string][]blockSet{
	"A": {
		{"M01", [3]int{10, 11, 12}},
		{"C01", [3]int{20, 22, 25}},
		{"G01", [3]int{14, 14, 15}},
		{"B01", [3]int{12, 12, 13}},
		{"A01", [3]int{8, 9, 10}},
		{"S01", [3]int{20, 22, 25}},
		{"C02", [3]int{45, 50, 55}}, // plank — last
	},
	"B": {
		{"M01", [3]int{12, 13, 14}},
		{"C03", [3]int{25, 27, 30}},
		{"C04", [3]int{8, 9, 10}},
		{"B02", [3]int{12, 13, 14}},
		{"M02", [3]int{6, 7, 8}},
		{"A01", [3]int{10, 11, 12}},
		{"C05", [3]int{30, 32, 35}},
		{"C02", [3]int{55, 60, 65}}, // plank — last
	},
	"C": {
		{"M03", [3]int{10, 11, 12}},
		{"C01", [3]int{25, 28, 30}},
		{"C06", [3]int{10, 11, 12}},
		{"B03", [3]int{20, 22, 25}},
		{"B04", [3]int{12, 13, 14}},
		{"G02", [3]int{16, 17, 18}},
		{"L01", [3]int{8, 9, 10}},
		{"C05", [3]int{35, 38, 40}},
		{"C02", [3]int{60, 65, 70}}, // plank — last
	},
	"D": {
		{"M01", [3]int{14, 15, 16}},
		{"M02", [3]int{8, 9, 10}},
		{"C03", [3]int{30, 33, 35}},
		{"C07", [3]int{20, 23, 25}},
		{"B02", [3]int{14, 15, 16}},
		{"B05", [3]int{10, 11, 12}},
		{"A01", [3]int{12, 13, 14}},
		{"A02", [3]int{15, 18, 20}},
		{"C08", [3]int{25, 28, 30}}, // RKC plank — last
	},
	"E": {
		{"M01", [3]int{16, 16, 16}},
		{"M02", [3]int{10, 10, 12}},
		{"C01", [3]int{35, 38, 40}},
		{"C03", [3]int{35, 38, 40}},
		{"C04", [3]int{10, 11, 12}},
		{"B01", [3]int{15, 16, 18}},
		{"B02", [3]int{15, 16, 18}},
		{"L01", [3]int{10, 12, 12}},
		{"A01", [3]int{12, 14, 15}},
		{"C02", [3]int{75, 80, 90}}, // plank — last
	},
}

// ---- 30-day calendar ----

type dayPlan struct {
	Code     string // block letter
	Variant  int    // 0=−, 1=base, 2=+
	Recovery bool
	Control  bool
}

var variantSuffix = [3]string{"−", "", "+"}

// calendar[day] (1-indexed) maps each day to its block/variant per spec §8.
var calendar = map[int]dayPlan{
	1: {Code: "A", Variant: 0}, 2: {Code: "A", Variant: 1}, 3: {Code: "A", Variant: 2},
	4: {Recovery: true},
	5: {Code: "B", Variant: 0}, 6: {Code: "B", Variant: 1},
	7: {Recovery: true},
	8: {Code: "B", Variant: 2}, 9: {Code: "C", Variant: 0},
	10: {Recovery: true},
	11: {Code: "C", Variant: 1}, 12: {Code: "C", Variant: 2},
	13: {Recovery: true},
	14: {Code: "B", Variant: 2}, 15: {Code: "C", Variant: 1}, 16: {Code: "D", Variant: 0},
	17: {Recovery: true},
	18: {Code: "D", Variant: 1}, 19: {Code: "C", Variant: 2},
	20: {Recovery: true},
	21: {Code: "D", Variant: 2}, 22: {Code: "D", Variant: 1},
	23: {Recovery: true},
	24: {Code: "E", Variant: 0}, 25: {Code: "D", Variant: 2},
	26: {Recovery: true},
	27: {Code: "E", Variant: 1},
	28: {Recovery: true},
	29: {Code: "E", Variant: 2},
	30: {Control: true},
}

// BlockLabel returns a short human label for the day's block (e.g. "A+", "R").
func BlockLabel(day int) string {
	p := calendar[clampDay(day)]
	switch {
	case p.Control:
		return "Control"
	case p.Recovery:
		return "R"
	default:
		return p.Code + variantSuffix[p.Variant]
	}
}

func clampDay(day int) int {
	if day < 1 {
		return 1
	}
	if day > TotalDays {
		return TotalDays
	}
	return day
}

// restAfter applies the spec's rest rules.
func restAfter(e Exercise, lightRest int) int {
	switch e.Slot {
	case Warmup:
		return 5
	case Cooldown:
		return 0
	}
	switch e.ID {
	case "M01", "M02", "M03": // after push-ups
		return 35
	case "C02", "C05", "C08": // after a plank
		return 50
	}
	if lightRest < 10 {
		lightRest = 10
	}
	if lightRest > 40 {
		lightRest = 40
	}
	return lightRest
}

const restBetweenRounds = 75

// Build assembles the workout for a day (1..30) with no adaptation.
func Build(day, lightRest int) Workout {
	return BuildAdapted(day, lightRest, Adapt{})
}

// BuildAdapted assembles the workout for a day applying the adaptation state
// (variant level shift, forced recovery, pain substitutions). lightRest tunes
// the rest after light exercises (10..40 s); push-up/plank/between-round rests
// follow the spec.
func BuildAdapted(day, lightRest int, a Adapt) Workout {
	day = clampDay(day)
	plan := calendar[day]
	if a.ForceR && !plan.Control {
		plan = dayPlan{Recovery: true} // deload: next session becomes recovery
	}
	// Variant shifted by the adaptation level, clamped to [0,2].
	variant := plan.Variant + a.LevelDelta
	if variant < 0 {
		variant = 0
	}
	if variant > 2 {
		variant = 2
	}

	var items []Item
	add := func(id string, slot Slot, val, round int) {
		if slot == Main {
			id, val = substitute(id, val, a)
		}
		e := Pool[id]
		items = append(items, Item{
			ExerciseID: id,
			Unit:       e.Unit,
			Value:      val,
			Slot:       slot,
			PerSide:    e.PerSide,
			Round:      round,
		})
	}

	// Warm-up: 2 rounds.
	for r := 1; r <= 2; r++ {
		for _, s := range warmupSeq {
			add(s.ID, Warmup, s.Val, r)
		}
	}

	// Main block.
	switch {
	case plan.Control:
		for r := 1; r <= 2; r++ {
			for _, s := range controlSeq {
				add(s.ID, Main, s.Val, r)
			}
		}
		add("C02", Main, controlPlankSec, 3) // plank challenge
	case plan.Recovery:
		for _, s := range recoverySeq {
			add(s.ID, Main, s.Val, 1)
		}
	default:
		for r := 1; r <= 2; r++ {
			for _, bs := range blocks[plan.Code] {
				add(bs.ID, Main, bs.Vals[variant], r)
			}
		}
	}

	// Cool-down: 1 round.
	for _, s := range cooldownSeq {
		add(s.ID, Cooldown, s.Val, 1)
	}

	// Rest + duration estimate.
	est := 0
	for i := range items {
		it := items[i]
		e := Pool[it.ExerciseID]
		// rest after this item
		if i < len(items)-1 {
			next := items[i+1]
			if it.Slot == Main && next.Slot == Main && next.Round != it.Round {
				items[i].Rest = restBetweenRounds // between main rounds
			} else {
				items[i].Rest = restAfter(e, lightRest)
			}
		}
		// duration: per-side work counts double
		work := it.Value
		if it.Unit == Reps {
			work *= 3 // ~3s per rep
		} else if it.Unit == Breaths {
			work *= 4 // ~4s per breath
		}
		if it.PerSide {
			work *= 2
		}
		est += work + items[i].Rest
	}

	// Label reflects the actually-built plan (forced recovery / shifted variant).
	label := BlockLabel(day)
	switch {
	case plan.Control:
		label = "Control"
	case plan.Recovery:
		label = "R"
	default:
		label = plan.Code + variantSuffix[variant]
	}

	return Workout{Day: day, BlockCode: label, Items: items, EstSec: est}
}
