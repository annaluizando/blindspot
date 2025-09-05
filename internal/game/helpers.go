package game

import (
	"time"
)

type GameStateHelpers struct {
	gs *GameState
}

func NewGameStateHelpers(gs *GameState) *GameStateHelpers {
	return &GameStateHelpers{gs: gs}
}

func (h *GameStateHelpers) IsErrorRecent() bool {
	return time.Since(h.gs.ErrorTimestamp) < 10*time.Second
}

func (h *GameStateHelpers) IsSuccessMessageRecent() bool {
	return time.Since(h.gs.SuccessMessageTimestamp) < 5*time.Second
}

func (h *GameStateHelpers) GetTotalChallengesCount() int {
	total := 0
	for _, set := range h.gs.ChallengeSets {
		total += len(set.Challenges)
	}
	return total
}

func (h *GameStateHelpers) GetCompletedChallengesCount() int {
	completed := 0
	for _, set := range h.gs.ChallengeSets {
		for _, challenge := range set.Challenges {
			if h.gs.IsChallengeCompleted(challenge.ID) {
				completed++
			}
		}
	}
	return completed
}

func (h *GameStateHelpers) GetOverallCompletionPercentage() int {
	total := h.GetTotalChallengesCount()
	if total == 0 {
		return 0
	}
	completed := h.GetCompletedChallengesCount()
	return (completed * 100) / total
}
