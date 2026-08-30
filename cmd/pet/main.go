package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shamathmika/shamathmika/internal/pet"
	"github.com/shamathmika/shamathmika/internal/render"
)

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
	default:
		fmt.Fprintln(os.Stderr, "usage: pet render | pet act --user NAME --title TITLE")
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

func write(s *pet.State, reaction pet.Action, now time.Time) error {
	b, err := os.ReadFile(pet.ReadmePath)
	if err != nil {
		return err
	}
	out, err := render.Insert(string(b), render.Section(s, reaction, now))
	if err != nil {
		return err
	}
	return os.WriteFile(pet.ReadmePath, []byte(out), 0o644)
}

func report(commit, reply string) error {
	for _, kv := range [][2]string{{"commit", commit}, {"reply", reply}} {
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
