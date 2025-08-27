package ui

import (
	"blindspot/internal/challenges"
	"blindspot/internal/game"

	"golang.org/x/sys/unix"

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

func getTerminalSize() (width, height int) {
	width, height = 80, 24

	ws, err := unix.IoctlGetWinsize(unix.Stdin, unix.TIOCGWINSZ)
	if err == nil {
		// Successfully got terminal size
		width = int(ws.Col)
		height = int(ws.Row)

		// Ensure minimum dimensions
		if width < 40 {
			width = 40
		}
		if height < 10 {
			height = 10
		}
	}

	return width, height
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
