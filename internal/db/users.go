package db

import (
	"database/sql"
	"errors"
)

// ErrEmailTaken is returned when registering an email that already exists.
var ErrEmailTaken = errors.New("email already registered")

func scanUser(row interface {
	Scan(dest ...any) error
}) (User, error) {
	var u User
	var voice, guest int
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Lang, &u.Rest, &voice, &guest, &u.VoiceMode, &u.Level, &u.Created)
	u.Voice = voice != 0
	u.IsGuest = guest != 0
	return u, err
}

const userCols = `id, email, name, lang, rest_seconds, voice, is_guest, voice_mode, level, created_at`

// CreateUser inserts a new account. Returns ErrEmailTaken on duplicate email.
func (d *DB) CreateUser(email, passHash, lang, name string) (User, error) {
	res, err := d.sql.Exec(
		`INSERT INTO users (email, pass_hash, lang, name) VALUES (?, ?, ?, ?)`,
		email, passHash, lang, name,
	)
	if err != nil {
		// modernc sqlite reports UNIQUE violations in the error text.
		return User{}, ErrEmailTaken
	}
	id, _ := res.LastInsertId()
	return d.UserByID(id)
}

// CreateGuest inserts an anonymous account identified only by a cookie session.
// email must be a unique placeholder (e.g. "guest:<token>").
func (d *DB) CreateGuest(email, lang string) (User, error) {
	res, err := d.sql.Exec(
		`INSERT INTO users (email, pass_hash, lang, is_guest) VALUES (?, '', ?, 1)`,
		email, lang,
	)
	if err != nil {
		return User{}, err
	}
	id, _ := res.LastInsertId()
	return d.UserByID(id)
}

// ConvertGuest upgrades a guest account into a real one in place, preserving all
// progress. Returns ErrEmailTaken if the email already belongs to someone else.
func (d *DB) ConvertGuest(userID int64, email, passHash, name string) error {
	if _, err := d.sql.Exec(
		`UPDATE users SET email = ?, pass_hash = ?, name = ?, is_guest = 0 WHERE id = ?`,
		email, passHash, name, userID,
	); err != nil {
		return ErrEmailTaken
	}
	return nil
}

// UserByID loads a user by primary key.
func (d *DB) UserByID(id int64) (User, error) {
	row := d.sql.QueryRow(`SELECT `+userCols+` FROM users WHERE id = ?`, id)
	return scanUser(row)
}

// UserByEmail loads a user by email (ok=false if not found).
func (d *DB) UserByEmail(email string) (User, bool, error) {
	row := d.sql.QueryRow(`SELECT `+userCols+` FROM users WHERE email = ?`, email)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, err
	}
	return u, true, nil
}

// PassHash returns the stored bcrypt hash for an email (ok=false if no user).
func (d *DB) PassHash(email string) (hash string, ok bool, err error) {
	row := d.sql.QueryRow(`SELECT pass_hash FROM users WHERE email = ?`, email)
	err = row.Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return hash, true, nil
}

// UpdateSettings updates per-user preferences. voiceMode is off|min|normal|detailed.
func (d *DB) UpdateSettings(userID int64, name, lang string, rest int, voiceMode string) error {
	v := 1
	if voiceMode == "off" {
		v = 0
	}
	_, err := d.sql.Exec(
		`UPDATE users SET name = ?, lang = ?, rest_seconds = ?, voice = ?, voice_mode = ? WHERE id = ?`,
		name, lang, rest, v, voiceMode, userID,
	)
	return err
}

// SafetyAck reports whether the user acknowledged the safety disclaimer.
func (d *DB) SafetyAck(userID int64) (bool, error) {
	var n int
	err := d.sql.QueryRow(`SELECT safety_ack FROM users WHERE id = ?`, userID).Scan(&n)
	return n != 0, err
}

// AckSafety records that the user acknowledged the safety disclaimer.
func (d *DB) AckSafety(userID int64) error {
	_, err := d.sql.Exec(`UPDATE users SET safety_ack = 1 WHERE id = ?`, userID)
	return err
}

// GetProgramJSON returns the user's custom program JSON ("" = built-in).
func (d *DB) GetProgramJSON(userID int64) (string, error) {
	var s string
	err := d.sql.QueryRow(`SELECT program_json FROM users WHERE id = ?`, userID).Scan(&s)
	return s, err
}

// SetProgramJSON stores (or clears, with "") the user's custom program JSON.
func (d *DB) SetProgramJSON(userID int64, j string) error {
	_, err := d.sql.Exec(`UPDATE users SET program_json = ? WHERE id = ?`, j, userID)
	return err
}

// ActiveProgram returns the stored active-program key ("custom" or a standard
// program key); the app layer normalizes legacy/unknown values.
func (d *DB) ActiveProgram(userID int64) (string, error) {
	var s string
	err := d.sql.QueryRow(`SELECT active_program FROM users WHERE id = ?`, userID).Scan(&s)
	return s, err
}

// SetActiveProgram stores which program is active ("custom" or a standard key).
func (d *DB) SetActiveProgram(userID int64, which string) error {
	_, err := d.sql.Exec(`UPDATE users SET active_program = ? WHERE id = ?`, which, userID)
	return err
}

// GetLevel returns the user's universal difficulty level (clamped −3..+3).
func (d *DB) GetLevel(userID int64) int {
	var lvl int
	_ = d.sql.QueryRow(`SELECT level FROM users WHERE id = ?`, userID).Scan(&lvl)
	return clampLevel(lvl)
}

// SetLevel stores the user's universal difficulty level (clamped −3..+3).
func (d *DB) SetLevel(userID int64, level int) error {
	_, err := d.sql.Exec(`UPDATE users SET level = ? WHERE id = ?`, clampLevel(level), userID)
	return err
}

func clampLevel(l int) int {
	if l < -3 {
		return -3
	}
	if l > 3 {
		return 3
	}
	return l
}

// GetCustomExercises returns the user's custom exercise library JSON ("" = none).
func (d *DB) GetCustomExercises(userID int64) (string, error) {
	var s string
	err := d.sql.QueryRow(`SELECT custom_exercises FROM users WHERE id = ?`, userID).Scan(&s)
	return s, err
}

// SetCustomExercises stores (or clears, with "") the user's custom exercise library.
func (d *DB) SetCustomExercises(userID int64, j string) error {
	_, err := d.sql.Exec(`UPDATE users SET custom_exercises = ? WHERE id = ?`, j, userID)
	return err
}

// SetLang updates only the language (used by the quick switcher).
func (d *DB) SetLang(userID int64, lang string) error {
	_, err := d.sql.Exec(`UPDATE users SET lang = ? WHERE id = ?`, lang, userID)
	return err
}
