// Managed by GoX v0.1.36

//line app.gox:1
// App is the root component: HTML shell, head (PWA meta + styles), and the router.
package main

import (
	"morning/assets"
	"morning/components"
	"morning/internal/auth"
	"morning/path"
	"morning/segments/pages"

	"github.com/doors-dev/doors"
	"github.com/doors-dev/gox"
)

// App holds the shared reactive session created by the page factory.
type App struct {
	auth doors.Source[auth.Session]
}

//line app.gox:20
func (a App) Main() gox.Elem {
	return gox.Elem(func(__c gox.Cursor) (__e error) {
		ctx := __c.Context(); _ = ctx
		__e = __c.Raw("<!doctype html>"); if __e != nil { return }
		__e = __c.Init("html"); if __e != nil { return }
		{
//line app.gox:22
			__e = __c.Set("lang", "en"); if __e != nil { return }
			__e = __c.Submit(); if __e != nil { return }
			__e = __c.Init("head"); if __e != nil { return }
			{
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("meta"); if __e != nil { return }
				{
//line app.gox:24
					__e = __c.Set("charset", "utf-8"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("meta"); if __e != nil { return }
				{
//line app.gox:26
					__e = __c.Set("name", "viewport"); if __e != nil { return }
//line app.gox:27
					__e = __c.Set("content", "width=device-width, initial-scale=1, maximum-scale=1, viewport-fit=cover"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("meta"); if __e != nil { return }
				{
//line app.gox:28
					__e = __c.Set("name", "theme-color"); if __e != nil { return }
//line app.gox:28
					__e = __c.Set("content", "#0f172a"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("meta"); if __e != nil { return }
				{
//line app.gox:29
					__e = __c.Set("name", "mobile-web-app-capable"); if __e != nil { return }
//line app.gox:29
					__e = __c.Set("content", "yes"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("meta"); if __e != nil { return }
				{
//line app.gox:30
					__e = __c.Set("name", "apple-mobile-web-app-capable"); if __e != nil { return }
//line app.gox:30
					__e = __c.Set("content", "yes"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("meta"); if __e != nil { return }
				{
//line app.gox:31
					__e = __c.Set("name", "apple-mobile-web-app-status-bar-style"); if __e != nil { return }
//line app.gox:31
					__e = __c.Set("content", "black-translucent"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("meta"); if __e != nil { return }
				{
//line app.gox:32
					__e = __c.Set("name", "apple-mobile-web-app-title"); if __e != nil { return }
//line app.gox:32
					__e = __c.Set("content", "Morning"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("link"); if __e != nil { return }
				{
//line app.gox:33
					__e = __c.Set("rel", "manifest"); if __e != nil { return }
//line app.gox:33
					__e = __c.Set("href", "/static/manifest.webmanifest"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("link"); if __e != nil { return }
				{
//line app.gox:34
					__e = __c.Set("rel", "icon"); if __e != nil { return }
//line app.gox:34
					__e = __c.Set("type", "image/svg+xml"); if __e != nil { return }
//line app.gox:34
					__e = __c.Set("href", "/static/icon.svg"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("link"); if __e != nil { return }
				{
//line app.gox:35
					__e = __c.Set("rel", "apple-touch-icon"); if __e != nil { return }
//line app.gox:35
					__e = __c.Set("href", "/static/icon-180.png"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.InitVoid("link"); if __e != nil { return }
				{
//line app.gox:36
					__e = __c.Set("rel", "stylesheet"); if __e != nil { return }
//line app.gox:36
					__e = __c.Set("href", assets.Style); if __e != nil { return }
//line app.gox:36
					__e = __c.Set("name", "main.css"); if __e != nil { return }
				}
				__e = __c.Submit(); if __e != nil { return }
				__e = __c.Init("title"); if __e != nil { return }
				{
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Text("Morning"); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
			__e = __c.Init("body"); if __e != nil { return }
			{
				__e = __c.Submit(); if __e != nil { return }
//line app.gox:40
				__e = __c.Any(doors.Route(
				doors.RouteModel(func(p doors.Source[path.Path]) gox.Comp {
					return pages.Root{Path: p, Auth: a.auth}
				}),
				doors.RouteLocationDefault(components.NotFound),
			)); if __e != nil { return }
				__e = __c.Init("script"); if __e != nil { return }
				{
					__e = __c.Submit(); if __e != nil { return }
					__e = __c.Raw("if (\"serviceWorker\" in navigator) {\n\t\t\t\t\tnavigator.serviceWorker.register(\"/static/sw.js\").catch(function () {})\n\t\t\t\t}"); if __e != nil { return }
				}
				__e = __c.Close(); if __e != nil { return }
			}
			__e = __c.Close(); if __e != nil { return }
		}
		__e = __c.Close(); if __e != nil { return }
	return })
//line app.gox:54
}
