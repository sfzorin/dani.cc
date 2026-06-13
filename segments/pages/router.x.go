// Managed by GoX v0.1.36

//line router.gox:1
// Package pages holds the application's page components (router, auth, home,
// workout, settings) and their server handlers.
package pages

import (
	"morning/internal/auth"
	"morning/path"

	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

// Root is the top dispatcher. It binds the shared session: unauthenticated
// visitors see the auth pages; authenticated users see the app.
type Root struct {
	Path doors.Source[path.Path]
	Auth doors.Source[auth.Session]
}

//line router.gox:20
func (r Root) Main() gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line router.gox:21
		__e = __c.Any(r.Auth.Bind(func(s auth.Session) gox.Elem {
		return gox.Elem(func(__c gox.Cursor) (__e error) {
			ctx := __c.Context(); _ = ctx
//line router.gox:22
			if !s.Authorized {
//line router.gox:23
				__e = __c.Any(r.publicRoutes(s)); if __e != nil { return }
			} else  {
//line router.gox:25
				__e = __c.Any(r.privateRoutes(s)); if __e != nil { return }
			}
		return })
//line router.gox:27
	})); if __e != nil { return }
	return })
//line router.gox:28
}

// publicRoutes: login (default) and register.
//line router.gox:31
func (r Root) publicRoutes(s auth.Session) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line router.gox:32
		__e = __c.Any(r.Path.Route(
		doors.RouteMatch(func(p path.Path) bool { return p.Page == path.Register }).
			Comp(authPage{sess: s, auth: r.Auth, path: r.Path, register: true}),
		doors.RouteDefaultComp[path.Path](authPage{sess: s, auth: r.Auth, path: r.Path, register: false}),
	)); if __e != nil { return }
	return })
//line router.gox:37
}

// privateRoutes: settings, the workout player (/day/:n), and home (default).
// Guests can also reach the auth pages to upgrade to a real account.
//line router.gox:41
func (r Root) privateRoutes(s auth.Session) gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
//line router.gox:42
		__e = __c.Any(r.Path.Route(
		doors.RouteMatch(func(p path.Path) bool { return p.Page == path.Register }).
			Comp(authPage{sess: s, auth: r.Auth, path: r.Path, register: true}),
		doors.RouteMatch(func(p path.Path) bool { return p.Page == path.Login }).
			Comp(authPage{sess: s, auth: r.Auth, path: r.Path, register: false}),
		doors.RouteMatch(func(p path.Path) bool { return p.Page == path.Settings }).
			Comp(settingsPage{sess: s, auth: r.Auth, path: r.Path}),
		doors.RouteDerive(func(p path.Path) (int, bool) {
			return p.Day, p.Page == path.Day
		}).Bind(func(day int) gox.Elem {
			return workoutPage{sess: s, day: day, path: r.Path}.Main()
		}),
		doors.RouteDefaultComp[path.Path](homePage{sess: s, auth: r.Auth, path: r.Path}),
	)); if __e != nil { return }
	return })
//line router.gox:56
}
