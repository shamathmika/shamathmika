package pet

import (
	"math"
	"path/filepath"
	"testing"
	"time"
)

var start = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func near(t *testing.T, got, want float64, what string) {
	t.Helper()
	if math.Abs(got-want) > 0.0001 {
		t.Errorf("%s = %v, want %v", what, got, want)
	}
}

func TestDecayOverElapsedTime(t *testing.T) {
	cases := []struct {
		name                   string
		elapsed                time.Duration
		food, water, affection float64
	}{
		{"no time at all", 0, 100, 100, 100},
		{"one hour", time.Hour, 100 - FoodDrainPerDay/24, 100 - WaterDrainPerDay/24, 100 - AffectionDrainPerDay/24},
		{"one day", 24 * time.Hour, 100 - FoodDrainPerDay, 100 - WaterDrainPerDay, 100 - AffectionDrainPerDay},
		{"three days", 72 * time.Hour, 100 - 3*FoodDrainPerDay, 100 - 3*WaterDrainPerDay, 100 - 3*AffectionDrainPerDay},
		{"a week empties food", 7 * 24 * time.Hour, 0, 0, 100 - 7*AffectionDrainPerDay},
		{"a month is still zero", 30 * 24 * time.Hour, 0, 0, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := New(start)
			s.Food, s.Water, s.Affection = 100, 100, 100
			s.Decay(start.Add(c.elapsed))
			near(t, s.Food, c.food, "food")
			near(t, s.Water, c.water, "water")
			near(t, s.Affection, c.affection, "affection")
		})
	}
}

func TestDecayIgnoresHowOftenItRuns(t *testing.T) {
	once := New(start)
	once.Decay(start.Add(48 * time.Hour))

	stepped := New(start)
	for i := 1; i <= 48; i++ {
		stepped.Decay(start.Add(time.Duration(i) * time.Hour))
	}

	near(t, stepped.Food, once.Food, "food after hourly runs")
	near(t, stepped.Water, once.Water, "water after hourly runs")
	near(t, stepped.Affection, once.Affection, "affection after hourly runs")
}

func TestDecayDrainsWaterFastestAffectionSlowest(t *testing.T) {
	s := New(start)
	s.Decay(start.Add(24 * time.Hour))
	if !(s.Water < s.Food && s.Food < s.Affection) {
		t.Errorf("want water < food < affection, got %v %v %v", s.Water, s.Food, s.Affection)
	}
}

func TestDecayAdvancesTheClockOnce(t *testing.T) {
	s := New(start)
	now := start.Add(24 * time.Hour)
	s.Decay(now)
	after := s.Food
	s.Decay(now)
	near(t, s.Food, after, "food after a repeat decay")
	if !s.LastDecayAt.Equal(now) {
		t.Errorf("lastDecayAt = %v, want %v", s.LastDecayAt, now)
	}
}

func TestDecayIgnoresTimeRunningBackwards(t *testing.T) {
	s := New(start)
	s.Decay(start.Add(-48 * time.Hour))
	near(t, s.Food, StartStat, "food")
	if !s.LastDecayAt.Equal(start) {
		t.Errorf("lastDecayAt moved backwards to %v", s.LastDecayAt)
	}
}

func TestStatsFloorAtZero(t *testing.T) {
	s := New(start)
	s.Decay(start.Add(365 * 24 * time.Hour))
	for _, got := range []float64{s.Food, s.Water, s.Affection} {
		if got != 0 {
			t.Errorf("stat = %v, want 0", got)
		}
	}
}

func TestActionCapsAtHundred(t *testing.T) {
	s := New(start)
	s.Food = 90
	full := s.Apply(Feed, "someone", start)
	if full {
		t.Error("a pet at 90 should not read as full")
	}
	near(t, s.Food, MaxStat, "food")
}

func TestFeedingAFullPetChangesNothingButTheReply(t *testing.T) {
	s := New(start)
	s.Food = MaxStat
	full := s.Apply(Feed, "someone", start)
	if !full {
		t.Error("want full = true when the stat was already 100")
	}
	near(t, s.Food, MaxStat, "food")
	if s.TotalActions != 1 {
		t.Errorf("totalActions = %d, want 1", s.TotalActions)
	}
}

func TestActionDecaysFirstThenAddsBoost(t *testing.T) {
	s := New(start)
	now := start.Add(24 * time.Hour)
	s.Apply(Water, "someone", now)
	near(t, s.Water, StartStat-WaterDrainPerDay+Boost, "water")
	near(t, s.Food, StartStat-FoodDrainPerDay, "food")
}

func TestActionsTouchOnlyTheirOwnStat(t *testing.T) {
	for _, c := range []struct {
		action Action
		need   Need
	}{{Feed, NeedFood}, {Water, NeedWater}, {Pet, NeedAffection}} {
		s := New(start)
		s.Food, s.Water, s.Affection = 10, 10, 10
		s.Apply(c.action, "someone", start)
		near(t, s.Get(c.need), 10+Boost, string(c.need))
		for _, other := range []Need{NeedFood, NeedWater, NeedAffection} {
			if other != c.need {
				near(t, s.Get(other), 10, string(other))
			}
		}
	}
}

func TestHelpersAreCountedOnce(t *testing.T) {
	s := New(start)
	s.Apply(Feed, "ana", start)
	s.Apply(Water, "Ana", start)
	s.Apply(Pet, "ben", start)
	if len(s.Helpers) != 2 {
		t.Errorf("helpers = %v, want two distinct people", s.Helpers)
	}
	if s.TotalActions != 3 {
		t.Errorf("totalActions = %d, want 3", s.TotalActions)
	}
	if s.LastAction == nil || s.LastAction.User != "ben" || s.LastAction.Action != Pet {
		t.Errorf("lastAction = %+v, want ben petting", s.LastAction)
	}
}

func TestMoodThresholds(t *testing.T) {
	cases := []struct {
		lowest float64
		want   Mood
	}{
		{100, Delighted},
		{80, Delighted},
		{79.99, Content},
		{79, Content},
		{60, Content},
		{59.99, Fine},
		{35, Fine},
		{34.99, Droopy},
		{15, Droopy},
		{14.99, Wilting},
		{0, Wilting},
	}
	for _, c := range cases {
		if got := MoodFor(c.lowest); got != c.want {
			t.Errorf("MoodFor(%v) = %v, want %v", c.lowest, got, c.want)
		}
	}
}

func TestMoodFollowsTheLowestStat(t *testing.T) {
	s := New(start)
	s.Food, s.Water, s.Affection = 100, 100, 10
	if got := s.Mood(); got != Wilting {
		t.Errorf("mood = %v, want wilting", got)
	}
	need, low := s.Lowest()
	if need != NeedAffection || low != 10 {
		t.Errorf("lowest = %v %v, want affection 10", need, low)
	}
}

func TestLowestBreaksTiesInAStableOrder(t *testing.T) {
	s := New(start)
	s.Food, s.Water, s.Affection = 20, 20, 20
	if need, _ := s.Lowest(); need != NeedFood {
		t.Errorf("lowest = %v, want food on a three way tie", need)
	}
	s.Food = 50
	if need, _ := s.Lowest(); need != NeedWater {
		t.Errorf("lowest = %v, want water", need)
	}
}

func TestParseTitle(t *testing.T) {
	good := map[string]Action{
		"pet|feed":    Feed,
		"pet|water":   Water,
		"pet|pet":     Pet,
		"  pet|pet  ": Pet,
		"PET|FEED":    Feed,
	}
	for title, want := range good {
		got, err := ParseTitle(title)
		if err != nil {
			t.Errorf("ParseTitle(%q) errored: %v", title, err)
			continue
		}
		if got != want {
			t.Errorf("ParseTitle(%q) = %v, want %v", title, got, want)
		}
	}

	junk := []string{
		"",
		"pet",
		"pet|",
		"pet|hug",
		"pet|feed|feed",
		"pet |feed",
		"pet|feed now",
		"please pet|feed",
		"pet|feed\nrm -rf /",
		"pet|feed; echo hi",
		"petfeed",
		"Bug: the CI is broken",
		"pet|../../etc/passwd",
	}
	for _, title := range junk {
		if got, err := ParseTitle(title); err == nil {
			t.Errorf("ParseTitle(%q) = %v, want an error", title, got)
		}
	}
}

func TestActionMapping(t *testing.T) {
	for _, c := range []struct {
		action Action
		need   Need
		past   string
		title  string
	}{
		{Feed, NeedFood, "fed", "pet|feed"},
		{Water, NeedWater, "watered", "pet|water"},
		{Pet, NeedAffection, "petted", "pet|pet"},
	} {
		if got := c.action.Need(); got != c.need {
			t.Errorf("%v.Need() = %v, want %v", c.action, got, c.need)
		}
		if got := c.action.Past(); got != c.past {
			t.Errorf("%v.Past() = %q, want %q", c.action, got, c.past)
		}
		if got := c.action.Title(); got != c.title {
			t.Errorf("%v.Title() = %q, want %q", c.action, got, c.title)
		}
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "pet.json")

	fresh, err := LoadOrNew(path, start)
	if err != nil {
		t.Fatal(err)
	}
	near(t, fresh.Food, StartStat, "food on a new pet")
	if !fresh.BornAt.Equal(start) {
		t.Errorf("bornAt = %v, want %v", fresh.BornAt, start)
	}

	fresh.Apply(Water, "ana", start.Add(30*time.Hour))
	if err := fresh.Save(path); err != nil {
		t.Fatal(err)
	}

	back, err := LoadOrNew(path, start)
	if err != nil {
		t.Fatal(err)
	}
	near(t, back.Water, round2(fresh.Water), "water after a round trip")
	if back.TotalActions != 1 || len(back.Helpers) != 1 || back.LastAction == nil {
		t.Errorf("state did not survive the round trip: %+v", back)
	}
	if !back.LastDecayAt.Equal(fresh.LastDecayAt) {
		t.Errorf("lastDecayAt = %v, want %v", back.LastDecayAt, fresh.LastDecayAt)
	}
}
