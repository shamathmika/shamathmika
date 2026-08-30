package render

import (
	"fmt"
	"os"
	"strings"

	"github.com/shamathmika/shamathmika/internal/pet"
)

type Theme string

const (
	Light Theme = "light"
	Dark  Theme = "dark"
)

const NoReaction = pet.Action("")

const (
	viewX, viewY, viewW, viewH = 192, 225, 820, 826
	pivotX, pivotY             = 602, 1005
	drawHeight                 = 230
)

func LoadArt(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	s := string(b)
	i := strings.Index(s, "<path")
	j := strings.LastIndex(s, "</g>")
	if i < 0 || j < 0 || j < i {
		return "", fmt.Errorf("%s: no drawing found", path)
	}
	return strings.TrimSpace(s[i:j]), nil
}

type palette struct {
	ink   map[pet.Mood]string
	spark string
	blush string
	tear  string
	treat string
	shine string
}

func paletteFor(t Theme) palette {
	if t == Dark {
		return palette{
			ink: map[pet.Mood]string{
				pet.Delighted: "#f6f1e7",
				pet.Content:   "#e7e1d5",
				pet.Fine:      "#cbc5b9",
				pet.Droopy:    "#9d978c",
				pet.Wilting:   "#78736b",
			},
			spark: "#f2c14e",
			blush: "#d98b7c",
			tear:  "#8fc0e0",
			treat: "#e0776b",
			shine: "#fff8ec",
		}
	}
	return palette{
		ink: map[pet.Mood]string{
			pet.Delighted: "#14110d",
			pet.Content:   "#241f19",
			pet.Fine:      "#3d3830",
			pet.Droopy:    "#6b6459",
			pet.Wilting:   "#98928a",
		},
		spark: "#dd9f18",
		blush: "#e59283",
		tear:  "#6f9fc4",
		treat: "#c9524a",
		shine: "#ffffff",
	}
}

func Pet(art string, m pet.Mood, t Theme, reaction pet.Action) string {
	p := paletteFor(t)
	width := drawHeight * viewW / viewH

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="%d %d %d %d" width="%d" height="%d" role="img" aria-label="%s looks %s">`+"\n",
		viewX, viewY, viewW, viewH, width, drawHeight, pet.PetName, m)
	fmt.Fprintf(&b, `  <title>%s looks %s</title>`+"\n", pet.PetName, m)
	b.WriteString(style(reaction))
	b.WriteString(`  <g class="react"><g class="idle">` + "\n")
	fmt.Fprintf(&b, `    <g transform="%s" fill="%s" stroke="none">`+"\n", posture(m), p.ink[m])
	b.WriteString(`      <g transform="translate(0,1200) scale(0.1,-0.1)">` + "\n")
	b.WriteString(art)
	b.WriteString("\n      </g>\n    </g>\n  </g></g>\n")
	b.WriteString(sparkles(m, p))
	b.WriteString(props(reaction, p))
	b.WriteString("</svg>\n")
	return b.String()
}

func posture(m pet.Mood) string {
	tilt := map[pet.Mood]float64{
		pet.Delighted: -2,
		pet.Content:   -0.5,
		pet.Fine:      1,
		pet.Droopy:    2.5,
		pet.Wilting:   3.5,
	}[m]
	drop := map[pet.Mood]int{
		pet.Delighted: -6,
		pet.Content:   0,
		pet.Fine:      8,
		pet.Droopy:    20,
		pet.Wilting:   32,
	}[m]
	return fmt.Sprintf("translate(0 %d) rotate(%g %d %d)", drop, tilt, pivotX, pivotY)
}

func style(reaction pet.Action) string {
	var b strings.Builder
	b.WriteString("  <style>\n")
	fmt.Fprintf(&b, "    .idle{animation:breathe 4.5s ease-in-out infinite;transform-origin:%dpx %dpx}\n", pivotX, pivotY)
	b.WriteString("    @keyframes breathe{0%,100%{transform:scale(1,1)}50%{transform:scale(1.015,0.99)}}\n")

	switch reaction {
	case pet.Feed:
		fmt.Fprintf(&b, "    .react{animation:chomp 1.3s ease-in-out both;transform-origin:%dpx %dpx}\n", pivotX, pivotY)
		b.WriteString("    @keyframes chomp{0%,100%{transform:scale(1,1)}25%{transform:scale(1.05,0.95)}50%{transform:scale(0.98,1.03)}75%{transform:scale(1.02,0.99)}}\n")
		b.WriteString("    .prop{animation:drop 1.2s ease-in both}\n")
		b.WriteString(dropFrames)
	case pet.Water:
		fmt.Fprintf(&b, "    .react{animation:hop 1.1s ease-out 0.5s both;transform-origin:%dpx %dpx}\n", pivotX, pivotY)
		b.WriteString("    @keyframes hop{0%{transform:translateY(0)}30%{transform:translateY(-42px)}55%{transform:translateY(0)}75%{transform:translateY(-18px)}100%{transform:translateY(0)}}\n")
		b.WriteString("    .prop{animation:drop 1.2s ease-in both}\n")
		b.WriteString(dropFrames)
	case pet.Pet:
		fmt.Fprintf(&b, "    .react{animation:nuzzle 1.4s ease-in-out both;transform-origin:%dpx %dpx}\n", pivotX, pivotY)
		b.WriteString("    @keyframes nuzzle{0%,100%{transform:rotate(0deg)}20%{transform:rotate(-3deg)}55%{transform:rotate(3deg)}80%{transform:rotate(-1.2deg)}}\n")
		b.WriteString("    .prop{animation:rise 1.8s ease-out both}\n")
		b.WriteString("    .prop2{animation:rise 1.8s ease-out 0.35s both}\n")
		b.WriteString("    @keyframes rise{0%{opacity:0;transform:translateY(28px) scale(0.5)}20%{opacity:0.95}100%{opacity:0;transform:translateY(-160px) scale(1)}}\n")
	}

	b.WriteString("    @media (prefers-reduced-motion:reduce){.idle,.react,.prop,.prop2{animation:none}.prop,.prop2{opacity:0}}\n")
	b.WriteString("  </style>\n")
	return b.String()
}

const dropFrames = "    @keyframes drop{0%{opacity:0;transform:translateY(-340px)}15%{opacity:1}70%{opacity:1;transform:translateY(0)}88%{opacity:0}100%{opacity:0}}\n"

func sparkles(m pet.Mood, p palette) string {
	if m != pet.Delighted {
		return ""
	}
	return sparkle(258, 330, 27, p) + sparkle(948, 305, 22, p) + sparkle(232, 252, 15, p)
}

func props(reaction pet.Action, p palette) string {
	switch reaction {
	case pet.Feed:
		return fmt.Sprintf(
			`  <g class="prop"><circle cx="602" cy="742" r="26" fill="%s"/>`+
				`<circle cx="593" cy="733" r="7" fill="%s" opacity="0.7"/></g>`+"\n",
			p.treat, p.shine)
	case pet.Water:
		return fmt.Sprintf(
			`  <g class="prop"><path d="M 602 700 C 580 728 577 748 589 757 C 602 766 620 751 613 733 Z" fill="%s"/></g>`+"\n",
			p.tear)
	case pet.Pet:
		return fmt.Sprintf(`  <g class="prop">%s</g>`+"\n"+`  <g class="prop2">%s</g>`+"\n",
			heart(915, 470, 32, p), heart(300, 430, 25, p))
	}
	return ""
}

func heart(cx, cy, s float64, p palette) string {
	d := fmt.Sprintf("M %g %g C %g %g %g %g %g %g C %g %g %g %g %g %g Z",
		cx, cy+s*0.6,
		cx-s*1.1, cy-s*0.2, cx-s*0.55, cy-s*1.1, cx, cy-s*0.35,
		cx+s*0.55, cy-s*1.1, cx+s*1.1, cy-s*0.2, cx, cy+s*0.6)
	return fmt.Sprintf(`<path d="%s" fill="%s"/>`, d, p.blush)
}

func sparkle(cx, cy, s float64, p palette) string {
	d := fmt.Sprintf("M %g %g C %g %g %g %g %g %g C %g %g %g %g %g %g C %g %g %g %g %g %g C %g %g %g %g %g %g Z",
		cx, cy-s,
		cx+s*0.18, cy-s*0.3, cx+s*0.3, cy-s*0.18, cx+s, cy,
		cx+s*0.3, cy+s*0.18, cx+s*0.18, cy+s*0.3, cx, cy+s,
		cx-s*0.18, cy+s*0.3, cx-s*0.3, cy+s*0.18, cx-s, cy,
		cx-s*0.3, cy-s*0.18, cx-s*0.18, cy-s*0.3, cx, cy-s)
	return fmt.Sprintf(`  <path d="%s" fill="%s"/>`+"\n", d, p.spark)
}
