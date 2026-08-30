package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shamathmika/shamathmika/internal/pet"
	"github.com/shamathmika/shamathmika/internal/render"
)

var moods = []pet.Mood{pet.Delighted, pet.Content, pet.Fine, pet.Droopy, pet.Wilting}

func main() {
	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	var err error
	switch cmd {
	case "render":
		err = renderNow()
	case "act":
		err = act(os.Args[2:])
	case "preview":
		err = preview("preview")
	default:
		fmt.Fprintln(os.Stderr, "usage: pet render | pet act --user NAME --title TITLE | pet preview")
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "pet:", err)
		os.Exit(1)
	}
}

func renderNow() error {
	now := time.Now().UTC().Truncate(time.Second)

	s, err := pet.LoadOrNew(pet.StatePath, now)
	if err != nil {
		return err
	}
	s.Decay(now)
	if err := s.Save(pet.StatePath); err != nil {
		return err
	}
	if err := write(s, render.NoReaction, now); err != nil {
		return err
	}

	fmt.Printf("%s is %s: food %.0f, water %.0f, affection %.0f\n",
		pet.PetName, s.Mood(), s.Food, s.Water, s.Affection)
	return output("commit", "pet: time passes")
}

func act(args []string) error {
	fs := flag.NewFlagSet("act", flag.ExitOnError)
	user := fs.String("user", "", "github login of the visitor")
	title := fs.String("title", "", "issue title")
	if err := fs.Parse(args); err != nil {
		return err
	}

	now := time.Now().UTC().Truncate(time.Second)

	a, err := pet.ParseTitle(*title)
	if err != nil {
		return report("", render.RejectReply(int(now.Unix())))
	}

	limited, err := pet.RecentlyActed(pet.HistoryPath, *user, now)
	if err != nil {
		return err
	}
	if limited {
		return report("", render.RateLimitReply(int(now.Unix())))
	}

	s, err := pet.LoadOrNew(pet.StatePath, now)
	if err != nil {
		return err
	}
	full := s.Apply(a, *user, now)
	if err := s.Save(pet.StatePath); err != nil {
		return err
	}
	if err := pet.AppendHistory(pet.HistoryPath, s.Record(a, *user, now)); err != nil {
		return err
	}
	if err := write(s, a, now); err != nil {
		return err
	}

	return report(fmt.Sprintf("pet: %s by @%s", a.Past(), *user), render.Reply(s, a, full))
}

func report(commit, reply string) error {
	for _, kv := range [][2]string{
		{"commit", commit},
		{"reply", reply},
	} {
		if err := output(kv[0], kv[1]); err != nil {
			return err
		}
	}
	return nil
}

func output(key, value string) error {
	line := key + "=" + strings.Join(strings.Fields(value), " ")
	path := os.Getenv("GITHUB_OUTPUT")
	if path == "" {
		fmt.Println(line)
		return nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, line)
	return err
}

func write(s *pet.State, reaction pet.Action, now time.Time) error {
	if err := render.WriteAssets(s, reaction, pet.AssetsDir); err != nil {
		return err
	}
	b, err := os.ReadFile(pet.ReadmePath)
	if err != nil {
		return err
	}
	out, err := render.Insert(string(b), render.Section(s, now))
	if err != nil {
		return err
	}
	return os.WriteFile(pet.ReadmePath, []byte(out), 0o644)
}

func preview(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	art, err := render.LoadArt(pet.ArtPath)
	if err != nil {
		return err
	}
	for _, m := range moods {
		for _, t := range []render.Theme{render.Light, render.Dark} {
			name := fmt.Sprintf("%s-%s.svg", m, t)
			if err := os.WriteFile(filepath.Join(dir, name), []byte(render.Pet(art, m, t, render.NoReaction)), 0o644); err != nil {
				return err
			}
		}
	}
	for _, a := range []pet.Action{pet.Feed, pet.Water, pet.Pet} {
		name := fmt.Sprintf("reaction-%s.svg", a)
		if err := os.WriteFile(filepath.Join(dir, name), []byte(render.Pet(art, pet.Content, render.Light, a)), 0o644); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "faces.html"), []byte(contactSheet(art)), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", filepath.Join(dir, "faces.html"))
	return nil
}

func contactSheet(art string) string {
	var b strings.Builder
	b.WriteString(`<!doctype html><meta charset="utf-8"><title>pet faces</title>`)
	b.WriteString(`<style>body{font:14px -apple-system,system-ui,sans-serif;margin:0}` +
		`section{padding:28px}.light{background:#ffffff;color:#24292f}.dark{background:#0d1117;color:#c9d1d9}` +
		`.row{display:flex;gap:24px;flex-wrap:wrap;align-items:flex-end}` +
		`figure{margin:0;text-align:center}figcaption{margin-top:8px;opacity:.75}` +
		`button{margin-left:12px;font:inherit}</style>`)
	for _, t := range []render.Theme{render.Light, render.Dark} {
		fmt.Fprintf(&b, `<section class="%s"><h2>%s mode</h2><div class="row">`, t, t)
		for _, m := range moods {
			fmt.Fprintf(&b, `<figure>%s<figcaption>%s</figcaption></figure>`, render.Pet(art, m, t, render.NoReaction), m)
		}
		b.WriteString(`</div></section>`)
	}
	b.WriteString(`<section class="light"><h2>reactions<button onclick="location.reload()">replay</button></h2><div class="row">`)
	for _, a := range []pet.Action{pet.Feed, pet.Water, pet.Pet} {
		fmt.Fprintf(&b, `<figure>%s<figcaption>%s</figcaption></figure>`, render.Pet(art, pet.Content, render.Light, a), a)
	}
	b.WriteString(`</div></section>`)
	return b.String()
}
