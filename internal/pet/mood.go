package pet

import "math"

type Mood int

const (
	Wilting Mood = iota
	Droopy
	Fine
	Content
	Delighted
)

const (
	DelightedAt = 80.0
	ContentAt   = 60.0
	FineAt      = 35.0
	DroopyAt    = 15.0
)

func (m Mood) String() string {
	switch m {
	case Delighted:
		return "delighted"
	case Content:
		return "content"
	case Fine:
		return "fine"
	case Droopy:
		return "droopy"
	default:
		return "wilting"
	}
}

func MoodFor(lowest float64) Mood {
	switch {
	case lowest >= DelightedAt:
		return Delighted
	case lowest >= ContentAt:
		return Content
	case lowest >= FineAt:
		return Fine
	case lowest >= DroopyAt:
		return Droopy
	default:
		return Wilting
	}
}

func (s *State) Mood() Mood {
	_, v := s.Lowest()
	return MoodFor(math.Round(v))
}

func (s *State) Lowest() (Need, float64) {
	need, low := NeedFood, s.Food
	if s.Water < low {
		need, low = NeedWater, s.Water
	}
	if s.Affection < low {
		need, low = NeedAffection, s.Affection
	}
	return need, low
}
