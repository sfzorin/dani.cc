// Package auth holds password hashing, session tokens and the session value
// shared across a user's pages. It is framework-agnostic (no doors import) so
// it can be used from both main and segment packages.
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SessionTTL is how long a login session (cookie + DB row) lives.
const SessionTTL = 90 * 24 * time.Hour

// Session is the per-user auth state stored reactively in the doors session store.
type Session struct {
	Authorized bool
	UserID     int64
	Email      string
	Name       string
	Lang       string
	Rest       int
	Voice      bool
	IsGuest    bool
}

// StoreKey keys the Source[Session] in the doors SessionStore.
type StoreKey struct{}

// HashPassword returns a bcrypt hash of the plaintext password.
func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword reports whether the plaintext matches the stored hash.
func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

// NewToken returns a cryptographically random URL-safe session token.
func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// NormalizeEmail trims and lowercases an email for storage/lookup.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ValidEmail reports whether email parses as an address.
func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
