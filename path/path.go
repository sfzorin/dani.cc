// Package path defines the URL path model for the app.
package path

// Page is the matched route variant (index into the path pattern).
type Page int

const (
	Home Page = iota
	Login
	Register
	Day
	Settings
)

// Path is the single path model. The Day field is populated for /day/:Day.
type Path struct {
	Page Page `path:"/ | /login | /register | /day/:Day | /settings"`
	Day  int
}
