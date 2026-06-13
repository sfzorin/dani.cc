// Package app holds process-wide dependencies shared by the page/segment code.
package app

import (
	"time"

	"morning/internal/db"
)

// DB is the application database, set once at startup in main.
var DB *db.DB

// Today returns the server-local calendar date as YYYY-MM-DD.
func Today() string {
	return time.Now().Format("2006-01-02")
}
