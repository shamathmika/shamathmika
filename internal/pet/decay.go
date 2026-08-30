package pet

import "time"

const (
	FoodDrainPerDay      = MaxStat / 7.0
	WaterDrainPerDay     = MaxStat / 6.0
	AffectionDrainPerDay = MaxStat / 8.0
)

func (s *State) Decay(now time.Time) {
	elapsed := now.Sub(s.LastDecayAt)
	if elapsed <= 0 {
		return
	}
	days := elapsed.Hours() / 24
	s.Food = clamp(s.Food - days*FoodDrainPerDay)
	s.Water = clamp(s.Water - days*WaterDrainPerDay)
	s.Affection = clamp(s.Affection - days*AffectionDrainPerDay)
	s.LastDecayAt = now
}

func clamp(v float64) float64 {
	if v < MinStat {
		return MinStat
	}
	if v > MaxStat {
		return MaxStat
	}
	return v
}
