package ui

import (
	"secure-code-game/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	gameState     *game.GameState
	activeView    tea.Model
	width, height int
}

// Creates and configures a new BubbleTea program
func InitializeUI(gameState *game.GameState) (*tea.Program, error) {
	app := AppModel{
		gameState:  gameState,
		activeView: NewMainMenu(gameState, 80, 24), // Default size until window size message
		width:      80,
		height:     24,
	}

	program := tea.NewProgram(
		app,
		tea.WithAltScreen(),       // Use the full terminal screen
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	return program, nil
}

// initializes the application
func (m AppModel) Init() tea.Cmd {
	return m.activeView.Init()
}

// handles messages and user input
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update active view with new size
		m.activeView, cmd = m.activeView.Update(msg)
		return m, cmd

	case SelectChallengeMsg:
		// Switch to the challenge view for the selected challenge
		m.activeView = NewChallengeView(m.gameState, msg.Challenge, m.width, m.height, MainMenu)
		return m, m.activeView.Init()

	case backToMenuMsg:
		if challengeView, ok := m.activeView.(*ChallengeView); ok {
			// if challenge view is active, check the source menu and redirect to specific menu
			switch challengeView.sourceMenu {
			case MainMenu:
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			case ProgressMenu:
				m.activeView = NewProgressMenu(m.gameState, m.width, m.height)
			case CategoryMenu:
				m.activeView = NewCategoryMenu(
					m.gameState,
					m.gameState.CurrentCategoryIdx,
					m.width,
					m.height,
					MainMenu,
				)
			default:
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			}
		} else if explanationView, ok := m.activeView.(*ExplanationView); ok {
			// if explanation view is active, go back to specific menu based on source
			switch explanationView.sourceMenu {
			case MainMenu:
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			case ProgressMenu:
				m.activeView = NewProgressMenu(m.gameState, m.width, m.height)
			case CategoryMenu:
				m.activeView = NewCategoriesMenu(m.gameState, m.width, m.height, MainMenu)
			case ChallengeMenu:
				m.activeView = NewCategoriesMenu(m.gameState, m.width, m.height, MainMenu)
			default:
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			}
		} else if menuView, ok := m.activeView.(*MenuView); ok {
			// if coming from challenge menu, check its source
			if menuView.type_ == ChallengeMenu {
				if menuView.sourceMenu == ProgressMenu {
					m.activeView = NewProgressMenu(m.gameState, m.width, m.height)
				} else {
					m.activeView = NewCategoriesMenu(m.gameState, m.width, m.height, MainMenu)
				}
			} else {
				// For other menu types, go back to main menu
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			}
		} else {
			// For any other view type, default to main menu
			m.activeView = NewMainMenu(m.gameState, m.width, m.height)
		}
		return m, m.activeView.Init()

	case nextChallengeMsg:
		if m.gameState.MoveToNextChallenge() {
			challenge := m.gameState.GetCurrentChallenge()
			m.activeView = NewChallengeView(m.gameState, challenge, m.width, m.height, MainMenu)
		} else {
			challenge, found := m.gameState.GetNextIncompleteChallenge()
			if found {
				m.activeView = NewChallengeView(m.gameState, challenge, m.width, m.height, MainMenu)
			} else {
				m.activeView = NewMainMenu(m.gameState, m.width, m.height) // should modify here NewMainMenu for congratulations view when no next challenge is found
			}
		}
		return m, m.activeView.Init()
	}

	// Handle updates in the active view
	m.activeView, cmd = m.activeView.Update(msg)
	return m, cmd
}

// renders the application
func (m AppModel) View() string {
	return m.activeView.View()
}
