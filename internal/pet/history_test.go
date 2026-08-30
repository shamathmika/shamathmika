package pet

import (
	"path/filepath"
	"testing"
	"time"
)

func TestRateLimitWindow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "history.jsonl")

	if recent, err := RecentlyActed(path, "ana", start); err != nil || recent {
		t.Fatalf("empty history said recent=%v err=%v", recent, err)
	}

	s := New(start)
	s.Apply(Feed, "ana", start)
	if err := AppendHistory(path, s.Record(Feed, "ana", start)); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		user string
		when time.Time
		want bool
	}{
		{"straight away", "ana", start.Add(RateLimit / 2), true},
		{"same login in another case", "ANA", start.Add(RateLimit / 2), true},
		{"someone else", "ben", start.Add(RateLimit / 2), false},
		{"a moment short of the limit", "ana", start.Add(RateLimit - time.Second), true},
		{"exactly at the limit", "ana", start.Add(RateLimit), false},
		{"the next day", "ana", start.Add(24 * time.Hour), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := RecentlyActed(path, c.user, c.when)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.want {
				t.Errorf("RecentlyActed = %v, want %v", got, c.want)
			}
		})
	}
}

func TestHistoryKeepsOneLinePerAction(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.jsonl")
	s := New(start)

	for i, user := range []string{"ana", "ben", "ana"} {
		at := start.Add(time.Duration(i) * 2 * RateLimit)
		s.Apply(Water, user, at)
		if err := AppendHistory(path, s.Record(Water, user, at)); err != nil {
			t.Fatal(err)
		}
	}

	when := start.Add(4*RateLimit + RateLimit/2)
	if recent, _ := RecentlyActed(path, "ben", when); recent {
		t.Error("ben acted three windows ago and should be clear")
	}
	if recent, _ := RecentlyActed(path, "ana", when); !recent {
		t.Error("ana acted half a window ago and should be held")
	}
}

func TestRecordCarriesTheResultingStats(t *testing.T) {
	s := New(start)
	s.Apply(Feed, "ana", start)
	r := s.Record(Feed, "ana", start)

	if r.Food != round2(s.Food) || r.Water != round2(s.Water) || r.Affection != round2(s.Affection) {
		t.Errorf("record %+v does not match state %+v", r, s)
	}
	if r.User != "ana" || r.Action != Feed || !r.At.Equal(start) {
		t.Errorf("record = %+v", r)
	}
}
