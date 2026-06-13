// Package db is the SQLite persistence layer (pure-Go driver, no CGO).
package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// DB wraps the connection pool.
type DB struct {
	sql *sql.DB
}

// User is an account row.
type User struct {
	ID      int64
	Email   string
	Name    string
	Lang    string
	Rest    int  // rest seconds between exercises (10..40)
	Voice   bool // voice cues enabled
	IsGuest bool // anonymous cookie-only account
	Created string
}

// Open opens (creating if needed) the SQLite database at path and migrates it.
func Open(path string) (*DB, error) {
	// _pragma busy_timeout avoids "database is locked" under concurrent writes.
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)", path)
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1) // SQLite: serialize writes; simplest correct setting
	d := &DB{sql: sqlDB}
	if err := d.migrate(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return d, nil
}

// Close closes the pool.
func (d *DB) Close() error { return d.sql.Close() }

func (d *DB) migrate() error {
	const schema = `
CREATE TABLE IF NOT EXISTS users (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    email        TEXT NOT NULL UNIQUE,
    pass_hash    TEXT NOT NULL,
    name         TEXT NOT NULL DEFAULT '',
    lang         TEXT NOT NULL DEFAULT 'en',
    rest_seconds INTEGER NOT NULL DEFAULT 20,
    voice        INTEGER NOT NULL DEFAULT 1,
    is_guest     INTEGER NOT NULL DEFAULT 0,    -- anonymous cookie-only account
    position     INTEGER NOT NULL DEFAULT 0,   -- completed sessions in the looping 30-day cycle
    adapt_level    INTEGER NOT NULL DEFAULT 0,  -- variant delta for the next workout
    adapt_force_r  INTEGER NOT NULL DEFAULT 0,  -- force the next workout to recovery
    knee_block     INTEGER NOT NULL DEFAULT 0,  -- sessions left without mini-squat
    shoulder_block INTEGER NOT NULL DEFAULT 0,  -- sessions left without narrow push-ups
    back_block     INTEGER NOT NULL DEFAULT 0,  -- sessions left swapping scissors/flutter
    low_energy     INTEGER NOT NULL DEFAULT 0,  -- consecutive "worse" energy reports
    safety_ack     INTEGER NOT NULL DEFAULT 0,  -- safety disclaimer acknowledged
    created_at   TEXT NOT NULL DEFAULT (datetime('now'))
);

-- One row per submitted post-workout questionnaire.
CREATE TABLE IF NOT EXISTS feedback (
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    on_date    TEXT NOT NULL,
    day        INTEGER NOT NULL,
    difficulty INTEGER NOT NULL DEFAULT 0,
    knee       INTEGER NOT NULL DEFAULT 0,
    back       INTEGER NOT NULL DEFAULT 0,
    shoulder   INTEGER NOT NULL DEFAULT 0,
    energy     TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- One row per calendar day the user worked out (drives streaks + date calendar).
CREATE TABLE IF NOT EXISTS activity (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    done_on TEXT NOT NULL,                     -- YYYY-MM-DD (local calendar day)
    PRIMARY KEY (user_id, done_on)
);

CREATE TABLE IF NOT EXISTS sessions (
    token      TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at INTEGER NOT NULL               -- unix seconds
);
CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);

DROP TABLE IF EXISTS progress;
`
	if _, err := d.sql.Exec(schema); err != nil {
		return err
	}
	// Add columns to pre-existing user tables (ignore if already present).
	for _, alter := range []string{
		`ALTER TABLE users ADD COLUMN position INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN name TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN is_guest INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN adapt_level INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN adapt_force_r INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN knee_block INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN shoulder_block INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN back_block INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN low_energy INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN safety_ack INTEGER NOT NULL DEFAULT 0`,
	} {
		if _, err := d.sql.Exec(alter); err != nil {
			if !strings.Contains(err.Error(), "duplicate column") {
				return err
			}
		}
	}
	return nil
}
