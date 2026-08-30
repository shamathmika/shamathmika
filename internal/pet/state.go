package pet

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Need string

const (
	NeedFood      Need = "food"
	NeedWater     Need = "water"
	NeedAffection Need = "affection"
)

type State struct {
	Food      float64 `json:"food"`
	Water     float64 `json:"water"`
	Affection float64 `json:"affection"`

	LastDecayAt  time.Time `json:"lastDecayAt"`
	TotalActions int       `json:"totalActions"`
	BornAt       time.Time `json:"bornAt"`

	LastAction *Event   `json:"lastAction,omitempty"`
	Helpers    []string `json:"helpers,omitempty"`
}

type Event struct {
	User   string    `json:"user"`
	Action Action    `json:"action"`
	At     time.Time `json:"at"`
}

func New(now time.Time) *State {
	now = now.UTC().Truncate(time.Second)
	return &State{
		Food:        StartStat,
		Water:       StartStat,
		Affection:   StartStat,
		LastDecayAt: now,
		BornAt:      now,
	}
}

func LoadOrNew(path string, now time.Time) (*State, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return New(now), nil
	}
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &s, nil
}

func (s *State) Save(path string) error {
	out := *s
	out.Food = round2(out.Food)
	out.Water = round2(out.Water)
	out.Affection = round2(out.Affection)

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')

	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *State) Apply(a Action, user string, now time.Time) bool {
	s.Decay(now)

	need := a.Need()
	before := s.Get(need)
	s.set(need, clamp(before+Boost))

	s.TotalActions++
	s.LastAction = &Event{User: user, Action: a, At: now.UTC().Truncate(time.Second)}
	s.addHelper(user)

	return before >= MaxStat
}

func (s *State) Get(n Need) float64 {
	switch n {
	case NeedWater:
		return s.Water
	case NeedAffection:
		return s.Affection
	default:
		return s.Food
	}
}

func (s *State) set(n Need, v float64) {
	switch n {
	case NeedWater:
		s.Water = v
	case NeedAffection:
		s.Affection = v
	default:
		s.Food = v
	}
}

func (s *State) addHelper(user string) {
	if user == "" {
		return
	}
	for _, h := range s.Helpers {
		if strings.EqualFold(h, user) {
			return
		}
	}
	s.Helpers = append(s.Helpers, user)
}

func round2(v float64) float64 { return math.Round(v*100) / 100 }
