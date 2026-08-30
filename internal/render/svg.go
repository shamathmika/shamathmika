package render

import (
	"fmt"
	"strings"

	"github.com/shamathmika/shamathmika/internal/pet"
)

type Theme string

const (
	Light Theme = "light"
	Dark  Theme = "dark"
)

const NoReaction = pet.Action("")

type palette struct {
	ink     string
	feature string
	body    string
	belly   string
	shadow  string
	blush   string
	spark   string
	shine   string
	tear    string
	treat   string
}

func paletteFor(t Theme) palette {
	if t == Dark {
		return palette{
			ink:     "#f4ecdd",
			feature: "#2b2419",
			body:    "#8d7546",
			belly:   "#ab9160",
			shadow:  "#6f5a35",
			blush:   "#d98b7c",
			spark:   "#f2c14e",
			shine:   "#fff8ec",
			tear:    "#8fc0e0",
			treat:   "#e0776b",
		}
	}
	return palette{
		ink:     "#3a322a",
		feature: "#3a322a",
		body:    "#f2cd88",
		belly:   "#fbe6bd",
		shadow:  "#e3b96f",
		blush:   "#e59283",
		spark:   "#dd9f18",
		shine:   "#ffffff",
		tear:    "#6f9fc4",
		treat:   "#c9524a",
	}
}

func Pet(m pet.Mood, t Theme, reaction pet.Action) string {
	p := paletteFor(t)

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="-16 0 212 200" width="212" height="200" role="img" aria-label="%s looks %s">`+"\n", pet.PetName, m)
	fmt.Fprintf(&b, `  <title>%s looks %s</title>`+"\n", pet.PetName, m)
	b.WriteString(style(reaction))
	b.WriteString(`  <g fill="none" stroke="` + p.ink + `" stroke-width="4" stroke-linecap="round" stroke-linejoin="round">` + "\n")
	b.WriteString(`    <g class="react"><g class="idle">` + "\n")
	b.WriteString(ears(m, p))
	b.WriteString(feet(p))
	b.WriteString(body(p))
	b.WriteString(face(m, p))
	b.WriteString("    </g></g>\n")
	b.WriteString(props(reaction, p))
	b.WriteString("  </g>\n</svg>\n")
	return b.String()
}

func style(reaction pet.Action) string {
	var b strings.Builder
	b.WriteString("  <style>\n")
	b.WriteString("    .idle{animation:breathe 4.5s ease-in-out infinite;transform-origin:90px 176px}\n")
	b.WriteString("    @keyframes breathe{0%,100%{transform:scale(1,1)}50%{transform:scale(1.025,0.98)}}\n")

	switch reaction {
	case pet.Feed:
		b.WriteString("    .react{animation:chomp 1.3s ease-in-out both;transform-origin:90px 176px}\n")
		b.WriteString("    @keyframes chomp{0%,100%{transform:scale(1,1)}25%{transform:scale(1.07,0.93)}50%{transform:scale(0.96,1.05)}75%{transform:scale(1.03,0.98)}}\n")
		b.WriteString("    .prop{animation:drop 1.2s ease-in both}\n")
		b.WriteString(dropFrames)
	case pet.Water:
		b.WriteString("    .react{animation:hop 1.1s ease-out 0.5s both;transform-origin:90px 176px}\n")
		b.WriteString("    @keyframes hop{0%{transform:translateY(0)}30%{transform:translateY(-12px)}55%{transform:translateY(0)}75%{transform:translateY(-5px)}100%{transform:translateY(0)}}\n")
		b.WriteString("    .prop{animation:drop 1.2s ease-in both}\n")
		b.WriteString(dropFrames)
	case pet.Pet:
		b.WriteString("    .react{animation:nuzzle 1.4s ease-in-out both;transform-origin:90px 176px}\n")
		b.WriteString("    @keyframes nuzzle{0%,100%{transform:rotate(0deg)}20%{transform:rotate(-3.5deg)}55%{transform:rotate(3.5deg)}80%{transform:rotate(-1.5deg)}}\n")
		b.WriteString("    .prop{animation:rise 1.8s ease-out both}\n")
		b.WriteString("    .prop2{animation:rise 1.8s ease-out 0.35s both}\n")
		b.WriteString("    @keyframes rise{0%{opacity:0;transform:translateY(8px) scale(0.5)}20%{opacity:0.95}100%{opacity:0;transform:translateY(-46px) scale(1)}}\n")
	}

	b.WriteString("    @media (prefers-reduced-motion:reduce){.idle,.react,.prop,.prop2{animation:none}.prop,.prop2{opacity:0}}\n")
	b.WriteString("  </style>\n")
	return b.String()
}

const dropFrames = "    @keyframes drop{0%{opacity:0;transform:translateY(-96px)}15%{opacity:1}70%{opacity:1;transform:translateY(0)}88%{opacity:0}100%{opacity:0}}\n"

func props(reaction pet.Action, p palette) string {
	switch reaction {
	case pet.Feed:
		return fmt.Sprintf(
			`    <g class="prop"><circle cx="90" cy="119" r="7" fill="%s" stroke="none"/>`+
				`<circle cx="87.5" cy="116.5" r="2" fill="%s" stroke="none" opacity="0.7"/></g>`+"\n",
			p.treat, p.shine)
	case pet.Water:
		return fmt.Sprintf(
			`    <g class="prop"><path d="M 90 108 C 84 116 83 122 86.5 124.5 C 90 127 95 123 93 118 Z" fill="%s" stroke="none"/></g>`+"\n",
			p.tear)
	case pet.Pet:
		return fmt.Sprintf(
			`    <g class="prop">%s</g>`+"\n"+`    <g class="prop2">%s</g>`+"\n",
			heart(142, 108, 9, p), heart(42, 100, 7, p))
	}
	return ""
}

func heart(cx, cy, s float64, p palette) string {
	d := fmt.Sprintf("M %g %g C %g %g %g %g %g %g C %g %g %g %g %g %g Z",
		cx, cy+s*0.6,
		cx-s*1.1, cy-s*0.2, cx-s*0.55, cy-s*1.1, cx, cy-s*0.35,
		cx+s*0.55, cy-s*1.1, cx+s*1.1, cy-s*0.2, cx, cy+s*0.6)
	return fmt.Sprintf(`<path d="%s" fill="%s" stroke="none"/>`, d, p.blush)
}

func body(p palette) string {
	return fmt.Sprintf(
		`      <path d="M 90 38 C 132 38 152 70 150 112 C 148 154 124 176 90 176 C 56 176 32 154 30 112 C 28 70 48 38 90 38 Z" fill="%s"/>`+"\n"+
			`      <ellipse cx="90" cy="130" rx="40" ry="32" fill="%s" stroke="none"/>`+"\n", p.body, p.belly)
}

func feet(p palette) string {
	return fmt.Sprintf(
		`      <ellipse cx="64" cy="178" rx="16" ry="9" fill="%s"/>`+"\n"+
			`      <ellipse cx="116" cy="178" rx="16" ry="9" fill="%s"/>`+"\n", p.body, p.body)
}

func ears(m pet.Mood, p palette) string {
	droop := map[pet.Mood]int{
		pet.Delighted: -8,
		pet.Content:   4,
		pet.Fine:      32,
		pet.Droopy:    88,
		pet.Wilting:   128,
	}[m]

	shape := fmt.Sprintf(
		`<path d="M 0 6 C -16 -12 -14 -38 2 -48 C 17 -39 14 -12 9 4 Z" fill="%s"/>`+
			`<path d="M 1 -4 C -8 -17 -7 -31 2 -39 C 10 -31 8 -17 5 -4 Z" fill="%s" stroke="none"/>`,
		p.body, p.shadow)

	return fmt.Sprintf(
		`      <g transform="translate(48,64) rotate(%d)">%s</g>`+"\n"+
			`      <g transform="translate(132,64) scale(-1,1) rotate(%d)">%s</g>`+"\n",
		-droop, shape, -droop, shape)
}

func face(m pet.Mood, p palette) string {
	drop := map[pet.Mood]int{pet.Delighted: 0, pet.Content: 0, pet.Fine: 2, pet.Droopy: 5, pet.Wilting: 8}[m]

	var b strings.Builder
	fmt.Fprintf(&b, `      <g transform="translate(0,%d)" stroke="%s">`+"\n", drop, p.feature)

	switch m {
	case pet.Delighted:
		b.WriteString(openEyes(9, 3.2, p))
		fmt.Fprintf(&b, `        <path d="M 76 122 Q 90 143 104 122 Z" fill="%s"/>`+"\n", p.feature)
		b.WriteString(blush(p, "0.6"))
		b.WriteString(sparkle(22, 76, 8, p))
		b.WriteString(sparkle(158, 66, 7, p))
		b.WriteString(sparkle(26, 44, 5, p))
	case pet.Content:
		b.WriteString(openEyes(8, 2.6, p))
		b.WriteString(`        <path d="M 78 126 Q 90 138 102 126"/>` + "\n")
		b.WriteString(blush(p, "0.45"))
	case pet.Fine:
		b.WriteString(openEyes(7, 0, p))
		b.WriteString(`        <path d="M 80 129 Q 90 135 100 129"/>` + "\n")
	case pet.Droopy:
		fmt.Fprintf(&b, `        <path d="M 57.5 100 A 8.5 8.5 0 0 0 74.5 100 Z" fill="%s" stroke="none"/>`+"\n", p.feature)
		fmt.Fprintf(&b, `        <path d="M 105.5 100 A 8.5 8.5 0 0 0 122.5 100 Z" fill="%s" stroke="none"/>`+"\n", p.feature)
		b.WriteString(`        <path d="M 56 96 Q 66 89 76 96"/>` + "\n")
		b.WriteString(`        <path d="M 104 96 Q 114 89 124 96"/>` + "\n")
		b.WriteString(`        <path d="M 80 133 Q 90 127 100 133"/>` + "\n")
	default:
		b.WriteString(`        <path d="M 58 101 Q 66 92 74 101"/>` + "\n")
		b.WriteString(`        <path d="M 106 101 Q 114 92 122 101"/>` + "\n")
		b.WriteString(`        <path d="M 77 137 Q 84 127 90 134 Q 96 142 103 131"/>` + "\n")
		fmt.Fprintf(&b, `        <path d="M 62 108 C 57 116 56 121 59 123 C 63 125 66 121 64 116 Z" fill="%s" stroke="none"/>`+"\n", p.tear)
	}

	b.WriteString("      </g>\n")
	return b.String()
}

func openEyes(r, shine float64, p palette) string {
	var b strings.Builder
	for _, cx := range []int{66, 114} {
		fmt.Fprintf(&b, `        <circle cx="%d" cy="99" r="%g" fill="%s" stroke="none"/>`+"\n", cx, r, p.feature)
		if shine > 0 {
			fmt.Fprintf(&b, `        <circle cx="%g" cy="%g" r="%g" fill="%s" stroke="none"/>`+"\n",
				float64(cx)+r*0.34, 99-r*0.36, shine, p.shine)
		}
	}
	return b.String()
}

func blush(p palette, opacity string) string {
	return fmt.Sprintf(
		`        <ellipse cx="48" cy="122" rx="9" ry="5.5" fill="%s" stroke="none" opacity="%s"/>`+"\n"+
			`        <ellipse cx="132" cy="122" rx="9" ry="5.5" fill="%s" stroke="none" opacity="%s"/>`+"\n",
		p.blush, opacity, p.blush, opacity)
}

func sparkle(cx, cy, s float64, p palette) string {
	d := fmt.Sprintf("M %g %g C %g %g %g %g %g %g C %g %g %g %g %g %g C %g %g %g %g %g %g C %g %g %g %g %g %g Z",
		cx, cy-s,
		cx+s*0.18, cy-s*0.3, cx+s*0.3, cy-s*0.18, cx+s, cy,
		cx+s*0.3, cy+s*0.18, cx+s*0.18, cy+s*0.3, cx, cy+s,
		cx-s*0.18, cy+s*0.3, cx-s*0.3, cy+s*0.18, cx-s, cy,
		cx-s*0.3, cy-s*0.18, cx-s*0.18, cy-s*0.3, cx, cy-s)
	return fmt.Sprintf(`        <path d="%s" fill="%s" stroke="none"/>`+"\n", d, p.spark)
}
