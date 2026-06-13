// Managed by GoX v0.1.36

//line not_found.gox:1
// Package components holds reusable UI components shared across pages.
package components

import (
	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

// NotFound renders the 404 fallback page.
// It receives the raw location source and:
//   - Sets the HTTP response status to 404 (initial render only)
//   - Sets the page title to "Not Found"
//   - Displays the unmatched path reactively via Bind
//
// Note: title and status tags can be rendered anywhere in the tree;
// Doors moves <title> into <head> and applies Status to the response.
//line not_found.gox:17
func NotFound(location doors.Source[doors.Location]) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line not_found.gox:18
		__e = __c.Any(doors.Status(404)); if __e != nil { return }
		__e = __c.Init("title"); if __e != nil { return }
		{
			__e = __c.Submit(); if __e != nil { return }
			__e = __c.Text("Not Found"); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
		__e = __c.Init("h1"); if __e != nil { return }
		{
			__e = __c.Submit(); if __e != nil { return }
			__e = __c.Text("Not Found"); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
		__e = __c.Init("p"); if __e != nil { return }
		{
			__e = __c.Submit(); if __e != nil { return }
//line not_found.gox:22
			__e = __c.Any(location.Bind(func(l doors.Location) gox.Elem {
			return gox.Elem(func(__c gox.Cursor) (__e error) {
				ctx := __c.Context(); _ = ctx
				__e = __c.Text("Page `"); if __e != nil { return }
//line not_found.gox:23
				__e = __c.Any(l.Path()); if __e != nil { return }
				__e = __c.Text("` not found. "); if __e != nil { return }
			return })
//line not_found.gox:24
		})); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line not_found.gox:26
}
