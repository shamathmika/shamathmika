package render

import (
	"fmt"
	"html"
	"math"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/shamathmika/shamathmika/internal/pet"
)

const barWidth = 20

func Section(s *pet.State, now time.Time) string {
	v := s.LastDecayAt.Unix()
	light := fmt.Sprintf("%s?v=%d", path.Join(pet.AssetsDir, pet.LightSVG), v)
	dark := fmt.Sprintf("%s?v=%d", path.Join(pet.AssetsDir, pet.DarkSVG), v)

	var b strings.Builder
	b.WriteString(`<div align="center">` + "\n")
	b.WriteString(`<picture>` + "\n")
	fmt.Fprintf(&b, `<source media="(prefers-color-scheme: dark)" srcset="%s">`+"\n", dark)
	fmt.Fprintf(&b, `<img src="%s" alt="%s" height="200">`+"\n", light, altText(s))
	b.WriteString(`</picture>` + "\n")

	b.WriteString(`<p>` + "\n")
	fmt.Fprintf(&b, `<code>%s</code><br>`+"\n", statLine("food", s.Food))
	fmt.Fprintf(&b, `<code>%s</code><br>`+"\n", statLine("water", s.Water))
	fmt.Fprintf(&b, `<code>%s</code>`+"\n", statLine("affection", s.Affection))
	b.WriteString(`</p>` + "\n")

	b.WriteString(`<p>` + "\n")
	fmt.Fprintf(&b, `<a href="%s">feed</a>`+"\n", ActionURL(pet.Feed))
	b.WriteString(`&nbsp;&nbsp;&middot;&nbsp;&nbsp;` + "\n")
	fmt.Fprintf(&b, `<a href="%s">give water</a>`+"\n", ActionURL(pet.Water))
	b.WriteString(`&nbsp;&nbsp;&middot;&nbsp;&nbsp;` + "\n")
	fmt.Fprintf(&b, `<a href="%s">pet</a>`+"\n", ActionURL(pet.Pet))
	b.WriteString(`</p>` + "\n")

	fmt.Fprintf(&b, `<p><sub>the links open a pre-filled issue, submit it and %s replies in about half a minute</sub></p>`+"\n", pet.PetName)
	fmt.Fprintf(&b, `<p><sub>%s</sub></p>`+"\n", lastLine(s, now))
	fmt.Fprintf(&b, `<p><sub>%s</sub></p>`+"\n", lifetimeLine(s))
	b.WriteString(`</div>` + "\n")
	return b.String()
}

func Insert(readme, section string) (string, error) {
	block := pet.StartMarker + "\n" + section + pet.EndMarker
	i := strings.Index(readme, pet.StartMarker)
	j := strings.Index(readme, pet.EndMarker)

	switch {
	case i < 0 && j < 0:
		if !strings.HasSuffix(readme, "\n") {
			readme += "\n"
		}
		return readme + "\n" + block + "\n", nil
	case i < 0 || j < 0 || j < i:
		return "", fmt.Errorf("readme markers are broken: expected %s then %s", pet.StartMarker, pet.EndMarker)
	default:
		return readme[:i] + block + readme[j+len(pet.EndMarker):], nil
	}
}

func WriteAssets(s *pet.State, reaction pet.Action, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	m := s.Mood()
	files := map[string]Theme{pet.LightSVG: Light, pet.DarkSVG: Dark}
	for name, theme := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(Pet(m, theme, reaction)), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func ActionURL(a pet.Action) string {
	q := url.Values{}
	q.Set("title", a.Title())
	q.Set("body", fmt.Sprintf("Submit this issue as it is. %s gets %s, then replies here and closes it in about half a minute. Nothing else to do.", pet.PetName, a.Past()))
	return "https://github.com/" + pet.Repo + "/issues/new?" + q.Encode()
}

func altText(s *pet.State) string {
	return fmt.Sprintf("%s looks %s. Food %d, water %d, affection %d out of 100.",
		pet.PetName, s.Mood(), whole(s.Food), whole(s.Water), whole(s.Affection))
}

func statLine(label string, v float64) string {
	filled := int(math.Round(v / pet.MaxStat * barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	return nbsp(fmt.Sprintf("%-9s %s %3d", label, bar, whole(v)))
}

func lastLine(s *pet.State, now time.Time) string {
	if s.LastAction == nil {
		return "nobody has come by yet"
	}
	e := s.LastAction
	user := html.EscapeString(e.User)
	return fmt.Sprintf(`last %s by <a href="https://github.com/%s">@%s</a>, %s`,
		e.Action.Past(), user, user, ago(now.Sub(e.At)))
}

func lifetimeLine(s *pet.State) string {
	if s.TotalActions == 0 {
		return "no visits yet"
	}
	visits := "1 visit"
	if s.TotalActions != 1 {
		visits = fmt.Sprintf("%d visits", s.TotalActions)
	}
	people := "1 person"
	if len(s.Helpers) != 1 {
		people = fmt.Sprintf("%d people", len(s.Helpers))
	}
	return fmt.Sprintf("%s from %s", visits, people)
}

func ago(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < 2*time.Minute:
		return "a minute ago"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 2*time.Hour:
		return "an hour ago"
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	case d < 48*time.Hour:
		return "yesterday"
	default:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	}
}

func whole(v float64) int { return int(math.Round(v)) }

func nbsp(s string) string { return strings.ReplaceAll(s, " ", "&nbsp;") }
