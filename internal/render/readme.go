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
		return "eating.webp"
	case pet.Water:
		return "drinking.webp"
	case pet.Pet:
		return "petted.webp"
	}
	if m >= pet.Content {
		return "happy.webp"
	}
	return "waiting.webp"
}

func scene(m pet.Mood, reaction pet.Action) string {
	switch reaction {
	case pet.Feed:
		return "taking a treat from a hand"
	case pet.Water:
		return "drinking from a bowl someone is holding"
	case pet.Pet:
		return "being scratched on the head"
	}
	if m >= pet.Content {
		return "sitting up, tail wagging"
	}
	return "sitting quietly, tail resting"
}

func Section(s *pet.State, reaction pet.Action, now time.Time) string {
	var b strings.Builder
	b.WriteString(`<div align="center">` + "\n")

	fmt.Fprintf(&b, `<p>Say hello to <strong>%s</strong> while you are here!</p>`+"\n", pet.PetName)

	fmt.Fprintf(&b, `<img src="%s" alt="%s" height="%d">`+"\n",
		path.Join(pet.AssetsDir, imageFor(s.Mood(), reaction)), altText(s, reaction), drawHeight)

	b.WriteString(`<p>` + "\n" + `Interact with her by giving her: <br>` + "\n")
	fmt.Fprintf(&b, `a <a href="%s">treat</a>`+"\n<br>\n", ActionURL(pet.Feed))
	fmt.Fprintf(&b, `some <a href="%s">water</a>`+"\n<br>\n", ActionURL(pet.Water))
	fmt.Fprintf(&b, `lots of <a href="%s">pets</a>`+"\n", ActionURL(pet.Pet))
	b.WriteString(`</p>` + "\n")

	b.WriteString(`<p><strong>How she is doing</strong></p>` + "\n")
	b.WriteString(`<p>` + "\n")
	fmt.Fprintf(&b, `<code>%s</code><br>`+"\n", statLine("food", s.Food))
	fmt.Fprintf(&b, `<code>%s</code><br>`+"\n", statLine("water", s.Water))
	fmt.Fprintf(&b, `<code>%s</code>`+"\n", statLine("affection", s.Affection))
	b.WriteString(`</p>` + "\n")

	fmt.Fprintf(&b, `<p><sub>%s<br>%s</sub></p>`+"\n", lastLine(s, now), lifetimeLine(s))
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
	owner, _, _ := strings.Cut(pet.Repo, "/")
	q.Set("body", fmt.Sprintf(
		"Submit this issue and %s gets %s. She replies here and closes it herself. "+
			"Then head back to [the profile](https://github.com/%s) and refresh to see her reaction.",
		pet.PetName, a.Gift(), owner))
	return "https://github.com/" + pet.Repo + "/issues/new?" + q.Encode()
}

const drawHeight = 350

func altText(s *pet.State, reaction pet.Action) string {
	return fmt.Sprintf("%s, a golden retriever, %s. Looking %s. Food %d, water %d, affection %d out of 100.",
		pet.PetName, scene(s.Mood(), reaction), s.Mood(), whole(s.Food), whole(s.Water), whole(s.Affection))
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
	return fmt.Sprintf(`<a href="https://github.com/%s">@%s</a> gave her %s, %s`,
		user, user, e.Action.Gift(), ago(now.Sub(e.At)))
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
