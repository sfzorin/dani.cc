// Package app holds process-wide dependencies shared by the page/segment code.
package app

import (
	"time"

	"danicc/internal/catalog"
	"danicc/internal/content"
	"danicc/internal/db"
)

// DB is the application database, set once at startup in main.
var DB *db.DB

// Today returns the server-local calendar date as YYYY-MM-DD.
func Today() string {
	return time.Now().Format("2006-01-02")
}

// ProgramByKey returns the program for "custom" (the stored one, if any) or a
// standard program by key, falling back to the default standard program.
func ProgramByKey(userID int64, key string) content.Resolved {
	if key == "custom" {
		if j, err := DB.GetProgramJSON(userID); err == nil && j != "" {
			if r, err := content.ParseResolved(j); err == nil && len(r.Days) > 0 {
				return r
			}
		}
	}
	if r, ok := content.StandardByKey(key); ok {
		return r
	}
	r, _ := content.StandardByKey(content.DefaultStandardKey())
	return r
}

// SeedBuiltins makes the global `exercises` table match the built-in library
// (authoritative upsert + prune), so it is the canonical base every user reads.
func SeedBuiltins() {
	ids := make([]string, 0, len(content.Library()))
	for _, ex := range content.Library() {
		if d, ok := catalog.Builtin(ex.ID); ok {
			_ = DB.SetExercise(ex.ID, d.Compact()) // upsert (code is canonical)
			ids = append(ids, ex.ID)
		}
	}
	_ = DB.PruneExercises(ids) // drop exercises removed from the library
}

// UserCatalog returns the user's exercise catalog: the global DB library overlaid
// with their custom exercises. Pages use it to resolve names/content/voice so
// custom exercises are selectable, shown and spoken.
func UserCatalog(userID int64) catalog.Catalog {
	base, _ := DB.AllExercises()
	j, _ := DB.GetCustomExercises(userID)
	return catalog.NewLayered(base, j)
}

// SaveCustomExercise adds or overrides a custom exercise in the user's library.
func SaveCustomExercise(userID int64, d catalog.Doc) error {
	c := UserCatalog(userID)
	c.Put(d)
	return DB.SetCustomExercises(userID, c.Marshal())
}

// HasCustomProgram reports whether the user saved a custom program.
func HasCustomProgram(userID int64) bool {
	j, err := DB.GetProgramJSON(userID)
	return err == nil && j != ""
}

// ActiveProgramKey returns the user's selected program key: "custom" when a saved
// custom program is active, otherwise a valid standard program key (the default
// when the stored value is empty or a legacy/unknown key).
func ActiveProgramKey(userID int64) string {
	k, _ := DB.ActiveProgram(userID)
	if k == "custom" && HasCustomProgram(userID) {
		return "custom"
	}
	if _, ok := content.StandardByKey(k); ok {
		return k
	}
	return content.DefaultStandardKey()
}

// UserProgram returns the user's active program (a standard set or custom).
func UserProgram(userID int64) content.Resolved {
	return ProgramByKey(userID, ActiveProgramKey(userID))
}

// UserWorkout builds a day's workout from the active program (built-in or custom),
// scaled by the user's universal difficulty level.
func UserWorkout(userID int64, day, rest int) content.Workout {
	return UserProgram(userID).Workout(day, rest, DB.GetLevel(userID))
}
