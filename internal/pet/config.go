package pet

import "time"

const PetName = "Jack"

const Repo = "shamathmika/shamathmika"

const (
	StatePath   = "state/pet.json"
	HistoryPath = "state/history.jsonl"
	ReadmePath  = "README.md"
	AssetsDir   = "assets"
	LightSVG    = "pet-light.svg"
	DarkSVG     = "pet-dark.svg"
)

const (
	StartMarker = "<!-- PET:START -->"
	EndMarker   = "<!-- PET:END -->"
)

const (
	MinStat   = 0.0
	MaxStat   = 100.0
	StartStat = 80.0
	Boost     = 35.0
)

const RateLimit = time.Hour
