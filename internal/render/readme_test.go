package render

import (
	"strings"
	"testing"
	"time"

	"github.com/shamathmika/shamathmika/internal/pet"
)

var now = time.Date(2026, 1, 8, 12, 0, 0, 0, time.UTC)

func TestInsertAppendsWhenTheMarkersAreMissing(t *testing.T) {
	got, err := Insert("# hello\n", "PET\n")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "# hello\n") {
		t.Errorf("the text above was disturbed: %q", got)
	}
	if !strings.Contains(got, pet.StartMarker+"\nPET\n"+pet.EndMarker) {
		t.Errorf("markers were not appended: %q", got)
	}
}

func TestInsertReplacesBetweenTheMarkersAndLeavesTheRestAlone(t *testing.T) {
	before := "# bio\n![stats](x.png)\n\n" + pet.StartMarker + "\nold pet\n" + pet.EndMarker + "\n\nfooter\n"
	got, err := Insert(before, "new pet\n")
	if err != nil {
		t.Fatal(err)
	}
	want := "# bio\n![stats](x.png)\n\n" + pet.StartMarker + "\nnew pet\n" + pet.EndMarker + "\n\nfooter\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestInsertRefusesBrokenMarkers(t *testing.T) {
	for _, readme := range []string{
		pet.StartMarker + "\nno end\n",
		"no start\n" + pet.EndMarker + "\n",
		pet.EndMarker + "\nbackwards\n" + pet.StartMarker + "\n",
	} {
		if _, err := Insert(readme, "pet\n"); err == nil {
			t.Errorf("Insert(%q) should have failed", readme)
		}
	}
}

func TestStatLine(t *testing.T) {
	cases := []struct {
		v      float64
		bars   int
		number string
	}{
		{100, 20, "100"},
		{97.6, 20, "98"},
		{50, 10, "50"},
		{2.4, 0, "2"},
		{0, 0, "0"},
	}
	for _, c := range cases {
		got := statLine("food", c.v)
		if n := strings.Count(got, "█"); n != c.bars {
			t.Errorf("statLine(food, %v) drew %d filled blocks, want %d", c.v, n, c.bars)
		}
		if !strings.HasSuffix(got, c.number) {
			t.Errorf("statLine(food, %v) = %q, want it to end in %q", c.v, got, c.number)
		}
	}
}

func TestTheNumberAndTheFaceAgree(t *testing.T) {
	s := pet.New(now)
	s.Food, s.Water, s.Affection = 79.6, 90, 90
	if got := s.Mood(); got != pet.Delighted {
		t.Errorf("mood = %v, but the bar reads 80", got)
	}
	if !strings.HasSuffix(statLine("food", s.Food), "80") {
		t.Errorf("bar = %q, want it to read 80", statLine("food", s.Food))
	}
}

func TestAgo(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{-time.Hour, "just now"},
		{20 * time.Second, "just now"},
		{90 * time.Second, "a minute ago"},
		{25 * time.Minute, "25 minutes ago"},
		{95 * time.Minute, "an hour ago"},
		{3 * time.Hour, "3 hours ago"},
		{30 * time.Hour, "yesterday"},
		{9 * 24 * time.Hour, "9 days ago"},
	}
	for _, c := range cases {
		if got := ago(c.d); got != c.want {
			t.Errorf("ago(%v) = %q, want %q", c.d, got, c.want)
		}
	}
}

func TestLifetimeLine(t *testing.T) {
	s := pet.New(now)
	if got := lifetimeLine(s); got != "no visits yet" {
		t.Errorf("lifetimeLine = %q on a fresh pet", got)
	}
	s.Apply(pet.Feed, "ana", now)
	if got := lifetimeLine(s); got != "1 visit from 1 person" {
		t.Errorf("lifetimeLine = %q, want singular", got)
	}
	s.Apply(pet.Water, "ben", now)
	if got := lifetimeLine(s); got != "2 visits from 2 people" {
		t.Errorf("lifetimeLine = %q, want plural", got)
	}
}

func TestSectionCarriesTheWholePetBlock(t *testing.T) {
	s := pet.New(now)
	s.Apply(pet.Water, "ana", now.Add(-3*time.Hour))
	got := Section(s, now)

	for _, want := range []string{
		`prefers-color-scheme: dark`,
		pet.LightSVG,
		pet.DarkSVG,
		"title=pet%7Cfeed",
		"title=pet%7Cwater",
		"title=pet%7Cpet",
		"last watered by",
		`https://github.com/ana`,
		"3 hours ago",
		"1 visit from 1 person",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("section is missing %q:\n%s", want, got)
		}
	}
}
