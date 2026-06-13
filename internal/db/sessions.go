package db

import (
	"database/sql"
	"errors"
	"time"
)

// CreateSession stores a session token for a user with the given lifetime.
func (d *DB) CreateSession(token string, userID int64, ttl time.Duration) error {
	exp := time.Now().Add(ttl).Unix()
	_, err := d.sql.Exec(
		`INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)`,
		token, userID, exp,
	)
	return err
}

// SessionUser returns the user for a valid, unexpired session token.
func (d *DB) SessionUser(token string) (User, bool, error) {
	var userID int64
	var exp int64
	row := d.sql.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE token = ?`, token)
	if err := row.Scan(&userID, &exp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, false, nil
		}
		return User{}, false, err
	}
	if time.Now().Unix() > exp {
		_ = d.DeleteSession(token)
		return User{}, false, nil
	}
	u, err := d.UserByID(userID)
	if err != nil {
		return User{}, false, err
	}
	return u, true, nil
}

// DeleteSession removes a session token (logout).
func (d *DB) DeleteSession(token string) error {
	_, err := d.sql.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}
