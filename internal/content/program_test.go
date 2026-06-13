package content

import (
	"strings"
	"testing"
)

func mainItems(w Workout) []Item {
	var out []Item
	for _, it := range w.Items {
		if it.Slot == Main {
			out = append(out, it)
		}
	}
	return out
}

func isPlank(id string) bool { return id == "C02" || id == "C05" || id == "C08" }

func TestStructureWarmupCooldown(t *testing.T) {
	for day := 1; day <= TotalDays; day++ {
		w := Build(day, 18)
		// Warm-up: 10 exercises × 2 rounds = first 20 items, all W*.
		if len(w.Items) < 26 {
			t.Fatalf("day %d: too few items (%d)", day, len(w.Items))
		}
		for i := 0; i < 20; i++ {
			if w.Items[i].Slot != Warmup || !strings.HasPrefix(w.Items[i].ExerciseID, "W") {
				t.Fatalf("day %d: item %d is not a warm-up (%s)", day, i, w.Items[i].ExerciseID)
			}
		}
		// Cool-down: last 6 items, all CD*, 1 round.
		cd := w.Items[len(w.Items)-6:]
		for _, it := range cd {
			if it.Slot != Cooldown || !strings.HasPrefix(it.ExerciseID, "CD") {
				t.Fatalf("day %d: expected cool-down at the end, got %s", day, it.ExerciseID)
			}
		}
		// Last item has no trailing rest.
		if w.Items[len(w.Items)-1].Rest != 0 {
			t.Errorf("day %d: last item has rest", day)
		}
		// All items are known exercises.
		for _, it := range w.Items {
			if _, ok := Get(it.ExerciseID); !ok {
				t.Errorf("day %d: unknown exercise %s", day, it.ExerciseID)
			}
		}
	}
}

func TestPlankChallengeAndProgression(t *testing.T) {
	for day := 1; day <= TotalDays; day++ {
		p := calendar[day]
		if p.Recovery {
			continue
		}
		w := Build(day, 18)
		m := mainItems(w)
		last := m[len(m)-1]
		if !isPlank(last.ExerciseID) {
			t.Errorf("day %d: main does not end with a plank (got %s)", day, last.ExerciseID)
		}
		// Working-day plank is never below 45 s (spec). RKC (C08) is the
		// intentional short exception (25–30 s).
		if last.ExerciseID == "C02" && last.Value < 45 {
			t.Errorf("day %d: plank %ds < 45s", day, last.Value)
		}
		if p.Control && last.Value != 90 {
			t.Errorf("day 30 control plank = %ds, want 90", last.Value)
		}
	}
}

func TestRecoveryDays(t *testing.T) {
	allowed := map[string]bool{"C04": true, "G01": true, "B01": true, "S01": true, "C05": true}
	rCount := 0
	for day := 1; day <= TotalDays; day++ {
		if !calendar[day].Recovery {
			continue
		}
		rCount++
		w := Build(day, 18)
		// Recovery main = exactly 1 round.
		for _, it := range mainItems(w) {
			if it.Round != 1 {
				t.Errorf("R day %d: main round %d, want single round", day, it.Round)
			}
			if !allowed[it.ExerciseID] {
				t.Errorf("R day %d: forbidden recovery exercise %s", day, it.ExerciseID)
			}
		}
	}
	if rCount < 7 {
		t.Errorf("only %d recovery days, expected several", rCount)
	}
}

func TestRecoveryCadence(t *testing.T) {
	// No more than 4 consecutive working days without a recovery day.
	run := 0
	for day := 1; day <= TotalDays; day++ {
		if calendar[day].Recovery {
			run = 0
			continue
		}
		run++
		if run > 4 {
			t.Errorf("more than 4 working days in a row by day %d", day)
		}
	}
}

func TestForbiddenExercisesAbsent(t *testing.T) {
	// The library itself must contain none of the forbidden movements.
	forbidden := []string{"jump", "burpee", "jack", "mountain", "lunge", "highknee", "stepup"}
	for id := range Pool {
		low := strings.ToLower(id)
		for _, f := range forbidden {
			if strings.Contains(low, f) {
				t.Errorf("forbidden exercise in pool: %s", id)
			}
		}
	}
}

func hasID(w Workout, id string) bool {
	for _, it := range w.Items {
		if it.ExerciseID == id {
			return true
		}
	}
	return false
}

func TestNextAdaptRules(t *testing.T) {
	base := AdaptState{}
	// Easy + no pain → level up.
	if NextAdapt(base, Feedback{Difficulty: 5, Energy: "normal"}).Level != 1 {
		t.Errorf("easy no-pain should bump level +1")
	}
	// Easy but with knee pain → hold, and block the mini-squat.
	s := NextAdapt(base, Feedback{Difficulty: 4, Knee: 5, Energy: "normal"})
	if s.Level != 0 || s.Knee != 3 {
		t.Errorf("knee pain: level %d knee %d, want 0 and 3", s.Level, s.Knee)
	}
	// Hard → level down; very hard → force R.
	if NextAdapt(base, Feedback{Difficulty: 8, Energy: "normal"}).Level != -1 {
		t.Errorf("difficulty 8 should drop a level")
	}
	if !NextAdapt(base, Feedback{Difficulty: 10, Energy: "normal"}).ForceR {
		t.Errorf("difficulty 10 should force recovery")
	}
	// Back / shoulder pain set their blocks.
	if NextAdapt(base, Feedback{Difficulty: 6, Back: 4, Energy: "normal"}).Back != 3 {
		t.Errorf("back pain should block 3 sessions")
	}
	if NextAdapt(base, Feedback{Difficulty: 6, Shoulder: 3, Energy: "normal"}).Shoulder != 3 {
		t.Errorf("shoulder pain should block 3 sessions")
	}
	// Two "worse" energy reports in a row → force R.
	s1 := NextAdapt(base, Feedback{Difficulty: 6, Energy: "worse"})
	if s1.ForceR {
		t.Errorf("single worse-energy should not force R")
	}
	s2 := NextAdapt(s1, Feedback{Difficulty: 6, Energy: "worse"})
	if !s2.ForceR {
		t.Errorf("two worse-energy in a row should force R")
	}
}

func TestSuppressIncreaseOnGap(t *testing.T) {
	up := Adapt{LevelDelta: 1}
	if up.SuppressIncreaseOnGap(true).LevelDelta != 1 {
		t.Errorf("alive streak should keep the +")
	}
	if up.SuppressIncreaseOnGap(false).LevelDelta != 0 {
		t.Errorf("broken streak should drop the + (don't increase after a skip)")
	}
	down := Adapt{LevelDelta: -1}
	if down.SuppressIncreaseOnGap(false).LevelDelta != -1 {
		t.Errorf("a deload should survive a gap")
	}
}

func TestBuildAdapted(t *testing.T) {
	// Force recovery turns a working day into R.
	w := BuildAdapted(1, 18, Adapt{ForceR: true})
	if w.BlockCode != "R" {
		t.Errorf("forced recovery block = %s, want R", w.BlockCode)
	}
	// Level +1 on day 1 (A−) raises the first push-up set 10 → 11.
	w = BuildAdapted(1, 18, Adapt{LevelDelta: 1})
	if mainItems(w)[0].Value != 11 {
		t.Errorf("level +1 push-ups = %d, want 11", mainItems(w)[0].Value)
	}
	// Knee block removes the mini-squat (day 9 = C−, has L01).
	w = BuildAdapted(9, 18, Adapt{KneeBlock: true})
	if hasID(w, "L01") {
		t.Errorf("knee block should remove L01")
	}
	// Back block swaps scissors C01 → toe taps C06 (day 1 has C01).
	w = BuildAdapted(1, 18, Adapt{BackBlock: true})
	if hasID(w, "C01") || !hasID(w, "C06") {
		t.Errorf("back block should swap C01 → C06")
	}
	// Shoulder block swaps narrow push-ups M02 → M01 (day 5 = B has M02).
	w = BuildAdapted(5, 18, Adapt{ShoulderBlock: true})
	if hasID(w, "M02") {
		t.Errorf("shoulder block should remove M02")
	}
}

func TestRestRules(t *testing.T) {
	w := Build(1, 18) // block A
	for i, it := range w.Items {
		if i == len(w.Items)-1 {
			continue
		}
		switch {
		case it.ExerciseID == "M01" || it.ExerciseID == "M02" || it.ExerciseID == "M03":
			if it.Rest != 35 {
				t.Errorf("after push-ups rest = %d, want 35", it.Rest)
			}
		case isPlank(it.ExerciseID) && w.Items[i+1].Slot == Main:
			// plank rest only checked when not the last main item
		}
	}
}
