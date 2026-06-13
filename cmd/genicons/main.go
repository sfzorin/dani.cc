// Command genicons writes PNG app icons (192, 512, 180) to assets/static/.
// A simple "sunrise" mark: orange sun over a dark background.
//
// Run: go run ./cmd/genicons
package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

var (
	bg   = color.RGBA{0x0f, 0x17, 0x2a, 0xff}
	sun1 = color.RGBA{0xfb, 0x92, 0x3c, 0xff}
	sun2 = color.RGBA{0xf9, 0x73, 0x16, 0xff}
)

func lerp(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		uint8(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		uint8(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		uint8(float64(a.B) + (float64(b.B)-float64(a.B))*t),
		0xff,
	}
}

func draw(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	cx := float64(size) / 2
	cy := float64(size) * 0.62
	r := float64(size) * 0.26
	horizon := float64(size) * 0.62
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			fx, fy := float64(x), float64(y)
			c := bg
			d := math.Hypot(fx-cx, fy-cy)
			if d <= r {
				c = lerp(sun1, sun2, d/r) // sun disc with vertical-ish shade
			} else {
				// faint rays above horizon
				if fy < horizon {
					ang := math.Atan2(horizon-fy, fx-cx)
					ray := math.Abs(math.Sin(ang * 6))
					if ray > 0.86 && d < r*2.1 {
						c = lerp(bg, sun1, (ray-0.86)/0.14*0.5)
					}
				}
			}
			img.Set(x, y, c)
		}
	}
	return img
}

func write(size int, name string) {
	f, err := os.Create(filepath.Join("assets", "static", name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, draw(size)); err != nil {
		panic(err)
	}
}

func main() {
	write(192, "icon-192.png")
	write(512, "icon-512.png")
	write(180, "icon-180.png")
	println("wrote icons")
}
