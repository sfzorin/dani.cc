package db

import (
	"path/filepath"
	"testing"
	"time"
)

func openTest(t *testing.T) *DB {
	t.Helper()
	d, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

func TestUserAndSettings(t *testing.T) {
	d := openTest(t)
	u, err := d.CreateUser("a@b.com", "hash", "ru", "Test")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if u.Lang != "ru" || u.Rest != 20 || !u.Voice {
		t.Errorf("defaults wrong: %+v", u)
	}
	if _, err := d.CreateUser("a@b.com", "h2", "en", "Test"); err != ErrEmailTaken {
		t.Errorf("expected ErrEmailTaken, got %v", err)
	}
	if err := d.UpdateSettings(u.ID, "Sergey", "de", 35, false); err != nil {
		t.Fatal(err)
	}
	u2, _ := d.UserByID(u.ID)
	if u2.Lang != "de" || u2.Rest != 35 || u2.Voice || u2.Name != "Sergey" {
		t.Errorf("settings not saved: %+v", u2)
	}
}

func TestStreakConsecutive(t *testing.T) {
	d := openTest(t)
	u, _ := d.CreateUser("s@b.com", "h", "en", "Test")
	today := "2026-06-13"
	for _, day := range []string{"2026-06-11", "2026-06-12", "2026-06-13"} {
		if err := d.RecordActivity(u.ID, day); err != nil {
			t.Fatal(err)
		}
	}
	// duplicate same-day record is a no-op
	_ = d.RecordActivity(u.ID, today)
	streak, err := d.Streak(u.ID, today)
	if err != nil {
		t.Fatal(err)
	}
	if streak != 3 {
		t.Errorf("streak = %d, want 3", streak)
	}

	// A gap breaks the streak (only today counts).
	d2 := openTest(t)
	u2, _ := d2.CreateUser("g@b.com", "h", "en", "Test")
	_ = d2.RecordActivity(u2.ID, "2026-06-10")
	_ = d2.RecordActivity(u2.ID, today)
	if s, _ := d2.Streak(u2.ID, today); s != 1 {
		t.Errorf("gapped streak = %d, want 1", s)
	}

	// Streak survives if yesterday (not today) was the last workout.
	if s, _ := d.Streak(u.ID, "2026-06-14"); s != 3 {
		t.Errorf("streak ending yesterday = %d, want 3", s)
	}
	// But not if neither today nor yesterday.
	if s, _ := d.Streak(u.ID, "2026-06-20"); s != 0 {
		t.Errorf("stale streak = %d, want 0", s)
	}
}

func TestCyclePositionLoops(t *testing.T) {
	d := openTest(t)
	u, _ := d.CreateUser("c@b.com", "h", "en", "Test")

	if p, _ := d.Position(u.ID); p != 0 {
		t.Fatalf("initial position = %d, want 0", p)
	}
	// Advance through a full 30-day cycle and into the next.
	for i := 0; i < 31; i++ {
		if err := d.AdvancePosition(u.ID); err != nil {
			t.Fatal(err)
		}
	}
	p, _ := d.Position(u.ID)
	if p != 31 {
		t.Fatalf("position = %d, want 31", p)
	}
	day := p%30 + 1
	cycle := p/30 + 1
	if day != 2 || cycle != 2 {
		t.Errorf("after 31 sessions: day %d cycle %d, want day 2 cycle 2", day, cycle)
	}

	// Reset clears activity and position.
	_ = d.RecordActivity(u.ID, "2026-06-13")
	if err := d.ResetProgress(u.ID); err != nil {
		t.Fatal(err)
	}
	if p, _ := d.Position(u.ID); p != 0 {
		t.Errorf("position after reset = %d, want 0", p)
	}
	if has, _ := d.HasActivity(u.ID, "2026-06-13"); has {
		t.Errorf("activity not cleared after reset")
	}
}

func TestActivityInRange(t *testing.T) {
	d := openTest(t)
	u, _ := d.CreateUser("r@b.com", "h", "en", "Test")
	for _, day := range []string{"2026-05-31", "2026-06-01", "2026-06-15", "2026-07-01"} {
		_ = d.RecordActivity(u.ID, day)
	}
	got, err := d.ActivityInRange(u.ID, "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || !got["2026-06-01"] || !got["2026-06-15"] {
		t.Errorf("range = %v, want only the two June days", got)
	}
}

func TestSessions(t *testing.T) {
	d := openTest(t)
	u, _ := d.CreateUser("se@b.com", "h", "en", "Test")
	if err := d.CreateSession("tok", u.ID, time.Hour); err != nil {
		t.Fatal(err)
	}
	got, ok, err := d.SessionUser("tok")
	if err != nil || !ok || got.ID != u.ID {
		t.Fatalf("session lookup failed: ok=%v err=%v", ok, err)
	}
	_ = d.DeleteSession("tok")
	if _, ok, _ := d.SessionUser("tok"); ok {
		t.Errorf("session still valid after delete")
	}
	// Expired session is rejected.
	_ = d.CreateSession("old", u.ID, -time.Hour)
	if _, ok, _ := d.SessionUser("old"); ok {
		t.Errorf("expired session accepted")
	}
}
