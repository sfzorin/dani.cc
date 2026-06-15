// Command gensvg generates schematic, animated exercise SVGs into
// assets/static/ex/<id>.svg, one per exercise in the content pool.
//
// Figures are simple stick figures animated with SMIL (limbs rotate around
// their joints). Colors are themed by slot (warmup / main / cooldown).
//
// Run: go run ./cmd/gensvg
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"morning30/internal/content"
)

// theme returns (background, figure) colors for a slot.
func theme(slot content.Slot) (string, string) {
	switch slot {
	case content.Warmup:
		return "#fff7ed", "#ea580c"
	case content.Cooldown:
		return "#ecfdf5", "#059669"
	default:
		return "#eff6ff", "#0284c7"
	}
}

// archetype maps each exercise ID to a (schematic) drawing function key.
var archetype = map[string]string{
	// Warm-up
	"W01": "armcircle", "W02": "neck", "W03": "armcircle", "W04": "armcircle",
	"W05": "armcircle", "W06": "armcircle", "W07": "armcircle", "W08": "twist",
	"W09": "fold", "W10": "calf",
	// Push-ups
	"M01": "pushup", "M02": "pushup", "M03": "pushup",
	// Core
	"C01": "legraise", "C02": "plank", "C03": "legraise", "C04": "crunch",
	"C05": "sideplank", "C06": "crunch", "C07": "crunch", "C08": "plank",
	// Back (prone)
	"B01": "superman", "B02": "superman", "B03": "superman", "B04": "superman", "B05": "superman",
	// Glutes / legs
	"G01": "bridge", "G02": "bridge", "L01": "squat",
	// Arms / shoulders
	"A01": "armcircle", "A02": "armcircle", "S01": "armcircle",
	// Cool-down
	"CD01": "kneechest", "CD02": "cobra", "CD03": "twist", "CD04": "fold",
	"CD05": "kneechest", "CD06": "cobra",
}

// figures holds the inner markup for each archetype. {{anim}} placeholders are
// filled inline below; we keep them as plain strings.
func figure(kind string) string {
	switch kind {
	case "march":
		// Standing, arms + knees alternate.
		return head(60, 26) + line(60, 35, 60, 66) +
			rotLine(60, 40, 46, 58, "25 60 40", "-15 60 40", 1.0) + // left arm
			rotLine(60, 40, 74, 58, "-25 60 40", "15 60 40", 1.0) + // right arm
			rotLine(60, 66, 50, 96, "-20 60 66", "20 60 66", 1.0) + // left leg
			rotLine(60, 66, 70, 96, "20 60 66", "-20 60 66", 1.0) // right leg
	case "squat":
		// Whole figure dips down and up.
		return grp("translate", "0 0", "0 16", 1.6,
			head(60, 26)+line(60, 35, 60, 64)+
				line(60, 42, 44, 56)+line(60, 42, 76, 56)+
				line(60, 64, 46, 90)+line(60, 64, 74, 90))
	case "lunge":
		return grp("translate", "0 0", "0 12", 1.6,
			head(58, 24)+line(58, 33, 60, 60)+
				line(60, 44, 46, 56)+line(60, 44, 74, 56)+
				line(60, 60, 40, 96)+line(60, 60, 84, 92))
	case "calf":
		return grp("translate", "0 0", "0 -8", 1.0,
			head(60, 30)+line(60, 39, 60, 70)+
				line(60, 46, 46, 62)+line(60, 46, 74, 62)+
				line(60, 70, 50, 98)+line(60, 70, 70, 98))
	case "twist":
		// Torso + arms rotate around the spine.
		return head(60, 26) + line(60, 70, 50, 98) + line(60, 70, 70, 98) +
			grp("rotate", "-22 60 70", "22 60 70", 1.4,
				line(60, 35, 60, 70)+line(40, 50, 80, 50))
	case "sidebend":
		// Lean side to side, one arm overhead.
		return line(60, 70, 50, 98) + line(60, 70, 70, 98) +
			grp("rotate", "-18 60 70", "18 60 70", 1.6,
				head(60, 26)+line(60, 35, 60, 70)+
					line(60, 40, 40, 30)+line(60, 44, 76, 60))
	case "armcircle":
		// Both arms sweep a full circle.
		return head(60, 30) + line(60, 39, 60, 74) +
			line(60, 74, 50, 100) + line(60, 74, 70, 100) +
			rotLine(60, 44, 60, 18, "0 60 44", "360 60 44", 2.2) +
			rotLine(60, 44, 60, 70, "180 60 44", "540 60 44", 2.2)
	case "neck":
		// Head tilts side to side.
		return line(60, 40, 60, 74) +
			line(60, 48, 44, 64) + line(60, 48, 76, 64) +
			line(60, 74, 50, 100) + line(60, 74, 70, 100) +
			grp("rotate", "-16 60 40", "16 60 40", 1.8, head(60, 28))
	case "hipcircle":
		return head(60, 26) + line(60, 35, 60, 60) +
			line(60, 42, 46, 56) + line(60, 42, 74, 56) +
			grp("rotate", "-14 60 60", "14 60 60", 1.6,
				line(60, 60, 48, 96)+line(60, 60, 72, 96))
	case "catcow":
		// On all fours, back arches up/down.
		return head(30, 54) +
			grp("translate", "0 0", "0 -8", 1.8, line(38, 56, 86, 56)) +
			line(86, 56, 86, 96) + line(40, 56, 40, 96) + head(94, 50)
	case "plank":
		// Horizontal hold, gentle bob.
		return grp("translate", "0 0", "0 4", 1.6,
			head(30, 56)+line(38, 58, 92, 70)+
				line(40, 58, 40, 92)+line(92, 70, 92, 96))
	case "sideplank":
		// Diagonal body, one arm down.
		return head(34, 44) + line(40, 48, 96, 92) +
			grp("translate", "0 0", "0 4", 1.6,
				line(54, 56, 54, 30))
	case "pushup":
		// Horizontal body moves up and down.
		return grp("translate", "0 0", "0 10", 1.4,
			head(30, 52)+line(38, 54, 92, 66))+
			line(44, 96, 44, 60) + line(92, 96, 92, 70)
	case "birddog":
		// All fours, opposite arm + leg extend.
		return line(48, 56, 48, 92) + line(80, 56, 80, 92) +
			head(36, 52) +
			grp("rotate", "0 58 56", "-8 58 56", 1.6,
				line(40, 56, 78, 56))+
			rotLine(36, 52, 18, 44, "10 36 52", "-20 36 52", 1.6) +
			rotLine(80, 92, 102, 100, "-10 80 92", "20 80 92", 1.6)
	case "superman":
		// Prone, arms + legs lift.
		return line(34, 76, 86, 76) +
			rotLine(34, 76, 16, 66, "0 34 76", "-22 34 76", 1.6) +
			rotLine(86, 76, 104, 66, "0 86 76", "22 86 76", 1.6) +
			head(28, 72)
	case "crunch":
		// Lying, knees bent, torso curls up.
		return line(96, 96, 80, 70) + line(80, 70, 96, 64) + // bent legs
			grp("rotate", "0 40 92", "-28 40 92", 1.4,
				head(30, 78)+line(38, 80, 60, 92))
	case "legraise":
		// Lying flat, legs lift up and down.
		return head(28, 86) + line(36, 88, 60, 88) +
			rotLine(60, 88, 100, 88, "0 60 88", "-70 60 88", 1.8)
	case "cobra":
		// Prone, chest lifts (rotate upper body up).
		return line(60, 92, 96, 92) +
			grp("rotate", "0 60 92", "-34 60 92", 1.8,
				head(34, 84)+line(42, 86, 60, 92))
	case "fold":
		// Seated/kneeling forward fold (slow lean).
		return line(40, 96, 88, 96) +
			grp("rotate", "-6 40 96", "-28 40 96", 2.2,
				head(40, 60)+line(44, 66, 40, 96)+line(44, 70, 84, 92))
	case "butterfly":
		// Seated, knees flap up and down.
		return head(60, 50) + line(60, 58, 60, 86) +
			grp("rotate", "10 60 86", "-6 60 86", 1.2, line(60, 86, 38, 96)) +
			grp("rotate", "-10 60 86", "6 60 86", 1.2, line(60, 86, 82, 96))
	case "kneechest":
		// Lying, knees hugged to chest (gentle pulse).
		return head(28, 80) +
			grp("translate", "0 0", "6 0", 1.8,
				line(36, 82, 64, 78)+line(64, 78, 50, 96))
	default:
		return head(60, 30) + line(60, 39, 60, 74) +
			line(60, 46, 46, 62) + line(60, 46, 74, 62) +
			line(60, 74, 50, 100) + line(60, 74, 70, 100)
	}
}

// ---- primitive builders ----

func head(cx, cy int) string {
	return fmt.Sprintf(`<circle cx="%d" cy="%d" r="9" fill="currentFill" stroke="none"/>`, cx, cy)
}

func line(x1, y1, x2, y2 int) string {
	return fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x1, y1, x2, y2)
}

// rotLine draws a line that rotates between two angle specs ("angle cx cy").
func rotLine(x1, y1, x2, y2 int, from, to string, dur float64) string {
	return fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d">%s</line>`,
		x1, y1, x2, y2, anim("rotate", from, to, dur))
}

// grp wraps markup in an animated <g>.
func grp(typ, from, to string, dur float64, inner string) string {
	return fmt.Sprintf(`<g>%s%s</g>`, anim(typ, from, to, dur), inner)
}

func anim(typ, from, to string, dur float64) string {
	return fmt.Sprintf(
		`<animateTransform attributeName="transform" type="%s" values="%s;%s;%s" dur="%.2fs" repeatCount="indefinite" calcMode="spline" keySplines="0.4 0 0.6 1;0.4 0 0.6 1" keyTimes="0;0.5;1"/>`,
		typ, from, to, from, dur)
}

func svg(slot content.Slot, kind string) string {
	bg, fg := theme(slot)
	fig := figure(kind)
	// currentFill is a stand-in we replace with fg for filled head circles.
	fig = replaceAll(fig, "currentFill", fg)
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 120 120" width="120" height="120" role="img">`+
		`<rect x="2" y="2" width="116" height="116" rx="20" fill="%s"/>`+
		`<g fill="none" stroke="%s" stroke-width="5" stroke-linecap="round" stroke-linejoin="round">%s</g></svg>`,
		bg, fg, fig)
}

func replaceAll(s, old, new string) string {
	out := ""
	for {
		i := indexOf(s, old)
		if i < 0 {
			return out + s
		}
		out += s[:i] + new
		s = s[i+len(old):]
	}
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func main() {
	dir := filepath.Join("assets", "static", "ex")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	n := 0
	for id, ex := range content.Pool {
		kind, ok := archetype[id]
		if !ok {
			kind = "default"
		}
		out := svg(ex.Slot, kind)
		if err := os.WriteFile(filepath.Join(dir, id+".svg"), []byte(out), 0o644); err != nil {
			panic(err)
		}
		n++
	}
	fmt.Printf("wrote %d exercise SVGs to %s\n", n, dir)
}
