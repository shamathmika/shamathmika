package render

import (
	"fmt"
	"strings"

	"github.com/shamathmika/shamathmika/internal/pet"
)

var acks = map[pet.Action][]string{
	pet.Feed: {
		"Thanks for the food",
		"Ate the whole thing",
		"That hit the spot",
		"Good snack, thank you",
	},
	pet.Water: {
		"Thanks for the water",
		"Cold and clean, thank you",
		"Drank all of it",
		"That was a good drink",
	},
	pet.Pet: {
		"Thanks for the scratch",
		"That was nice",
		"More of that please",
		"Good visit, thank you",
	},
}

var moodLines = map[pet.Mood][]string{
	pet.Delighted: {
		"I feel great.",
		"Nothing to ask for today.",
		"Full and happy.",
		"Best I have felt all week.",
	},
	pet.Content: {
		"Doing well here.",
		"No complaints.",
		"Steady and comfortable.",
		"All in good shape.",
	},
	pet.Fine: {
		"My %s is getting thin though.",
		"Could use some %s before long.",
		"%s is the low one now.",
		"Keeping an eye on my %s.",
	},
	pet.Droopy: {
		"My %s is running low.",
		"%s is what I need next.",
		"Short on %s, if you have a minute.",
		"Could really use some %s.",
	},
	pet.Wilting: {
		"I am very low on %s.",
		"Some %s, please, when you can.",
		"Running on empty for %s.",
		"Almost out of %s here.",
	},
}

var fullLines = []string{
	"Already full on that one, but thank you.",
	"Could not take any more, thank you though.",
	"That one is topped up. Thanks anyway.",
	"Full up there, but I liked the visit.",
}

var rateLimitLines = []string{
	"You were just here. Come back in a little while.",
	"Still going on your last visit. Try again in an hour.",
	"Thank you, but one thing per hour is plenty.",
	"You have already helped this hour. See you soon.",
}

var rejectLines = []string{
	"That is not a title I know. Use the links on the profile.",
	"I could not read that one. The links on the profile work.",
	"Not sure what that means. Try feed, water, or pet from the profile.",
}

func Reply(s *pet.State, a pet.Action, full bool) string {
	if full {
		return choose(fullLines, s.TotalActions)
	}
	need, _ := s.Lowest()
	line := choose(moodLines[s.Mood()], s.TotalActions*3+1)
	if strings.Contains(line, "%s") {
		line = fmt.Sprintf(line, need)
	}
	return choose(acks[a], s.TotalActions) + ". " + line
}

func RateLimitReply(pick int) string { return choose(rateLimitLines, pick) }

func RejectReply(pick int) string { return choose(rejectLines, pick) }

func choose(list []string, n int) string {
	if len(list) == 0 {
		return ""
	}
	if n < 0 {
		n = -n
	}
	return list[n%len(list)]
}
