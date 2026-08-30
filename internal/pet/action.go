package pet

import (
	"fmt"
	"regexp"
	"strings"
)

type Action string

const (
	Feed  Action = "feed"
	Water Action = "water"
	Pet   Action = "pet"
)

const TitlePrefix = "pet|"

var titleRE = regexp.MustCompile(`^pet\|(feed|water|pet)$`)

func ParseTitle(title string) (Action, error) {
	m := titleRE.FindStringSubmatch(strings.ToLower(strings.TrimSpace(title)))
	if m == nil {
		return "", fmt.Errorf("not a pet action: %q", title)
	}
	return Action(m[1]), nil
}

func (a Action) Need() Need {
	switch a {
	case Water:
		return NeedWater
	case Pet:
		return NeedAffection
	default:
		return NeedFood
	}
}

func (a Action) Past() string {
	switch a {
	case Water:
		return "watered"
	case Pet:
		return "petted"
	default:
		return "fed"
	}
}

func (a Action) Gift() string {
	switch a {
	case Water:
		return "some water"
	case Pet:
		return "lots of pets"
	default:
		return "a treat"
	}
}

func (a Action) Title() string { return TitlePrefix + string(a) }
