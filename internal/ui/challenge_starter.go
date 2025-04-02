package ui

import (
	"blindspot/internal/challenges"
	"blindspot/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

func InitializeUIWithChallenge(gs *game.GameState) (*tea.Program, error) {
	width, height := getTerminalSize()

	challenge := determineStartingChallenge(gs)

	challengeView := NewChallengeView(gs, challenge, width, height, MainMenu)

	program := tea.NewProgram(
		challengeView,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return program, nil
}

// this will be updated by the WindowSizeMsg when the program starts
// to-do: needs improvement
func getTerminalSize() (width, height int) {
	return 80, 24
}

// selects the appropriate challenge to start with
// may be deleted
func determineStartingChallenge(gs *game.GameState) challenges.Challenge {
	if len(gs.RandomizedChallenges) > 0 {
		if gs.CurrentChallengeIdx >= len(gs.RandomizedChallenges) {
			gs.CurrentChallengeIdx = 0
			gs.Progress.CurrentChallengeIdx = 0
			gs.SaveProgress()
		}
		return gs.RandomizedChallenges[gs.CurrentChallengeIdx]
	} // to-do: add else where randomized challenges are generated

	challenge, found := gs.GetNextIncompleteChallenge()
	if found {
		return challenge
	}

	return gs.ChallengeSets[0].Challenges[0]
}
