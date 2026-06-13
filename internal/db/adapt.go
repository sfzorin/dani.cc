package db

import "morning/internal/content"

// GetAdaptState loads the user's adaptation state (shapes the next workout).
func (d *DB) GetAdaptState(userID int64) (content.AdaptState, error) {
	var s content.AdaptState
	var force int
	err := d.sql.QueryRow(
		`SELECT adapt_level, adapt_force_r, knee_block, shoulder_block, back_block, low_energy
		 FROM users WHERE id = ?`, userID,
	).Scan(&s.Level, &force, &s.Knee, &s.Shoulder, &s.Back, &s.LowEnergy)
	s.ForceR = force != 0
	return s, err
}

// SetAdaptState persists the user's adaptation state.
func (d *DB) SetAdaptState(userID int64, s content.AdaptState) error {
	force := 0
	if s.ForceR {
		force = 1
	}
	_, err := d.sql.Exec(
		`UPDATE users SET adapt_level = ?, adapt_force_r = ?, knee_block = ?,
		        shoulder_block = ?, back_block = ?, low_energy = ? WHERE id = ?`,
		s.Level, force, s.Knee, s.Shoulder, s.Back, s.LowEnergy, userID,
	)
	return err
}

// AddFeedback logs one post-workout questionnaire.
func (d *DB) AddFeedback(userID int64, date string, day int, fb content.Feedback) error {
	_, err := d.sql.Exec(
		`INSERT INTO feedback (user_id, on_date, day, difficulty, knee, back, shoulder, energy)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, date, day, fb.Difficulty, fb.Knee, fb.Back, fb.Shoulder, fb.Energy,
	)
	return err
}
