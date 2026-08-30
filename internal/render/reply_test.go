package render

import (
	"strings"
	"testing"
	"time"

	"github.com/shamathmika/shamathmika/internal/pet"
)

func allLines() []string {
	var out []string
	for _, list := range acks {
		out = append(out, list...)
	}
	for _, list := range moodLines {
		out = append(out, list...)
	}
	out = append(out, fullLines...)
	out = append(out, rateLimitLines...)
	out = append(out, rejectLines...)
	return out
}

func TestNoDashesAnywhereInTheCopy(t *testing.T) {
	s := pet.New(now)
	s.Apply(pet.Water, "ana", now)

	copy := append(allLines(), Section(s, NoReaction, now))
	for _, dash := range []string{"—", "–"} {
		for _, line := range copy {
			if strings.Contains(line, dash) {
				t.Errorf("found %q in %q", dash, line)
			}
		}
	}
}

func TestEveryReplyIsOneShortLine(t *testing.T) {
	for _, line := range allLines() {
		if line == "" {
			t.Error("empty variant")
		}
		if strings.ContainsAny(line, "\n\r") {
			t.Errorf("variant spans lines: %q", line)
		}
		if len(line) > 90 {
			t.Errorf("variant is long: %q", line)
		}
	}
}

func TestReplyNamesTheLowestNeedWhenThePetIsStruggling(t *testing.T) {
	s := pet.New(now)
	s.Food, s.Water, s.Affection = 60, 10, 60

	got := Reply(s, pet.Feed, false)
	if !strings.Contains(got, string(pet.NeedWater)) {
		t.Errorf("Reply = %q, want it to mention water", got)
	}
	if !strings.HasPrefix(got, "Thanks for the food") && !strings.Contains(got, "snack") &&
		!strings.Contains(got, "Ate") && !strings.Contains(got, "spot") {
		t.Errorf("Reply = %q, want it to acknowledge the food", got)
	}
}

func TestReplyStaysQuietAboutNeedsWhenThePetIsHappy(t *testing.T) {
	s := pet.New(now)
	s.Food, s.Water, s.Affection = 100, 100, 100

	got := Reply(s, pet.Pet, false)
	for _, need := range []pet.Need{pet.NeedFood, pet.NeedWater, pet.NeedAffection} {
		if strings.Contains(got, string(need)) {
			t.Errorf("Reply = %q, want no complaint about %s", got, need)
		}
	}
}

func TestFullPetGetsItsOwnReply(t *testing.T) {
	s := pet.New(now)
	s.Food = 100
	full := s.Apply(pet.Feed, "ana", now)
	if !full {
		t.Fatal("expected the pet to read as full")
	}
	got := Reply(s, pet.Feed, full)
	if !strings.Contains(strings.ToLower(got), "full") && !strings.Contains(strings.ToLower(got), "any more") &&
		!strings.Contains(strings.ToLower(got), "topped up") {
		t.Errorf("Reply = %q, want it to say the stat was already full", got)
	}
}

func TestRepliesRotate(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < len(rateLimitLines); i++ {
		seen[RateLimitReply(i)] = true
	}
	if len(seen) != len(rateLimitLines) {
		t.Errorf("rate limit replies gave %d distinct lines, want %d", len(seen), len(rateLimitLines))
	}

	seen = map[string]bool{}
	s := pet.New(now)
	s.Food, s.Water, s.Affection = 40, 40, 40
	for i := 0; i < 4; i++ {
		s.TotalActions = i
		seen[Reply(s, pet.Water, false)] = true
	}
	if len(seen) != 4 {
		t.Errorf("four visits gave %d distinct replies, want 4", len(seen))
	}
}

func TestChooseSurvivesOddPicks(t *testing.T) {
	if got := choose(rejectLines, -7); got == "" {
		t.Error("a negative pick returned nothing")
	}
	if got := choose(nil, 3); got != "" {
		t.Errorf("choose(nil) = %q", got)
	}
	if got := RejectReply(int(time.Now().Unix())); got == "" {
		t.Error("RejectReply returned nothing")
	}
}
