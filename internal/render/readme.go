package render

import (
	"fmt"
	"html"
	"math"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/shamathmika/shamathmika/internal/pet"
)

const barWidth = 20

const NoReaction = pet.Action("")

func imageFor(m pet.Mood, reaction pet.Action) string {
	switch reaction {
	case pet.Feed:
		return "eating.png"
	case pet.Water:
		return "drinking.png"
	case pet.Pet:
		return "petted.png"
	}
	if m >= pet.Content {
		return "happy.png"
	}
	return "waiting.png"
}

func scene(reaction pet.Action) string {
	switch reaction {
	case pet.Feed:
		return "taking a treat from a hand"
	case pet.Water:
		return "drinking from a bowl someone is holding"
	case pet.Pet:
		return "being scratched on the head"
	}
	return "sitting on the beach"
}

func Section(s *pet.State, reaction pet.Action, now time.Time) string {
	var b strings.Builder
	b.WriteString(`<div align="center">` + "\n")
	fmt.Fprintf(&b, `<img src="%s" alt="%s" height="%d">`+"\n",
		path.Join(pet.AssetsDir, imageFor(s.Mood(), reaction)), altText(s, reaction), drawHeight)

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

func ActionURL(a pet.Action) string {
	q := url.Values{}
	q.Set("title", a.Title())
	q.Set("body", fmt.Sprintf("Submit this issue. %s gets %s, then replies here and closes it.", pet.PetName, a.Past()))
	return "https://github.com/" + pet.Repo + "/issues/new?" + q.Encode()
}

const drawHeight = 300

func altText(s *pet.State, reaction pet.Action) string {
	return fmt.Sprintf("%s, a labrador, %s. Looking %s. Food %d, water %d, affection %d out of 100.",
		pet.PetName, scene(reaction), s.Mood(), whole(s.Food), whole(s.Water), whole(s.Affection))
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
