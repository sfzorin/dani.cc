// Package assets manages embedded static files for the application.
// Files are embedded at compile time using Go's embed directive,
// then served via Doors middleware or used as managed resources.
package assets

import (
	"embed"
	"io/fs"
)

// Style is the embedded CSS content served as a Doors managed resource.
//
//go:embed style.css
var Style []byte

// Player is the client-side workout engine, served as a managed inline script.
//
//go:embed player.js
var Player []byte

// static holds the static/ directory tree (icons, manifest, sw, exercise SVGs),
// embedded recursively and exposed via Static() for the UseFS middleware.
//
//go:embed static
var static embed.FS

// Static returns an fs.FS rooted at the static/ subdirectory.
func Static() fs.FS {
	sub, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
