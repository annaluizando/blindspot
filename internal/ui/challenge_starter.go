package ui

import (
	"blindspot/internal/challenges"
	"blindspot/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

func InitializeUIWithChallenge(gs *game.GameState) (*tea.Program, error) {
	width, height := getTerminalSize()
	challenge := determineStartingChallenge(gs)

	app := AppModel{
		gameState:  gs,
		activeView: NewChallengeView(gs, challenge, width, height, MainMenu),
		width:      width,
		height:     height,
	}

	program := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return program, nil
}

func getTerminalSize() (width, height int) {
	width, height = 80, 24
	return width, height
}

func determineStartingChallenge(gs *game.GameState) challenges.Challenge {
	if gs.UseRandomizedOrder && len(gs.RandomizedChallenges) > 0 {
		return gs.RandomizedChallenges[0]
	}

	if gs.CurrentCategoryIdx < len(gs.ChallengeSets) {
		if gs.CurrentChallengeIdx < len(gs.ChallengeSets[gs.CurrentCategoryIdx].Challenges) {
			return gs.ChallengeSets[gs.CurrentCategoryIdx].Challenges[gs.CurrentChallengeIdx]
		}
	}

	if len(gs.ChallengeSets) > 0 && len(gs.ChallengeSets[0].Challenges) > 0 {
		return gs.ChallengeSets[0].Challenges[0]
	}

	return challenges.Challenge{}
}
