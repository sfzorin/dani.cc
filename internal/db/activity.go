package db

import "time"

// RecordActivity marks that the user worked out on the given calendar day
// (YYYY-MM-DD). Idempotent — one entry per day.
func (d *DB) RecordActivity(userID int64, day string) error {
	_, err := d.sql.Exec(
		`INSERT INTO activity (user_id, done_on) VALUES (?, ?)
		 ON CONFLICT(user_id, done_on) DO NOTHING`,
		userID, day,
	)
	return err
}

// HasActivity reports whether the user worked out on the given day.
func (d *DB) HasActivity(userID int64, day string) (bool, error) {
	var n int
	err := d.sql.QueryRow(
		`SELECT COUNT(*) FROM activity WHERE user_id = ? AND done_on = ?`,
		userID, day,
	).Scan(&n)
	return n > 0, err
}

// ActivityInRange returns the set of YYYY-MM-DD days the user worked out
// between start and end (inclusive), for the calendar view.
func (d *DB) ActivityInRange(userID int64, start, end string) (map[string]bool, error) {
	rows, err := d.sql.Query(
		`SELECT done_on FROM activity WHERE user_id = ? AND done_on >= ? AND done_on <= ?`,
		userID, start, end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]bool{}
	for rows.Next() {
		var on string
		if err := rows.Scan(&on); err != nil {
			return nil, err
		}
		out[on] = true
	}
	return out, rows.Err()
}

// Streak returns consecutive calendar days (ending today or yesterday) with a
// workout. `today` is YYYY-MM-DD.
func (d *DB) Streak(userID int64, today string) (int, error) {
	rows, err := d.sql.Query(`SELECT done_on FROM activity WHERE user_id = ?`, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	done := map[string]bool{}
	for rows.Next() {
		var on string
		if err := rows.Scan(&on); err != nil {
			return 0, err
		}
		done[on] = true
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	t, err := time.Parse("2006-01-02", today)
	if err != nil {
		return 0, err
	}
	cursor := t
	if !done[cursor.Format("2006-01-02")] {
		cursor = cursor.AddDate(0, 0, -1)
		if !done[cursor.Format("2006-01-02")] {
			return 0, nil
		}
	}
	streak := 0
	for done[cursor.Format("2006-01-02")] {
		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}
	return streak, nil
}

// Position returns how many cycle-sessions the user has completed. The current
// program day = position%30 + 1, and the current cycle = position/30 + 1.
func (d *DB) Position(userID int64) (int, error) {
	var p int
	err := d.sql.QueryRow(`SELECT position FROM users WHERE id = ?`, userID).Scan(&p)
	return p, err
}

// AdvancePosition moves the user one step forward in the looping 30-day cycle.
func (d *DB) AdvancePosition(userID int64) error {
	_, err := d.sql.Exec(`UPDATE users SET position = position + 1 WHERE id = ?`, userID)
	return err
}

// ResetProgress clears activity, feedback, cycle position and adaptation state.
func (d *DB) ResetProgress(userID int64) error {
	if _, err := d.sql.Exec(`DELETE FROM activity WHERE user_id = ?`, userID); err != nil {
		return err
	}
	if _, err := d.sql.Exec(`DELETE FROM feedback WHERE user_id = ?`, userID); err != nil {
		return err
	}
	_, err := d.sql.Exec(
		`UPDATE users SET position = 0, adapt_level = 0, adapt_force_r = 0,
		        knee_block = 0, shoulder_block = 0, back_block = 0, low_energy = 0
		 WHERE id = ?`, userID,
	)
	return err
}
