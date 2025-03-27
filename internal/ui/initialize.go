package ui

import (
	"secure-code-game/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

// Represents the entire application
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

	case SelectCategoryMsg:
		// Switch to challenge menu for the selected category
		m.activeView = NewCategoryMenu(m.gameState, m.width, m.height, MainMenu)
		return m, m.activeView.Init()

	case SelectChallengeMsg:
		// Switch to the challenge view for the selected challenge
		m.activeView = NewChallengeView(m.gameState, msg.Challenge, m.width, m.height, MainMenu)
		return m, m.activeView.Init()

	case backToMenuMsg:
		// Handle back action from current view
		if challengeView, ok := m.activeView.(*ChallengeView); ok {
			// If coming from challenge view, check the source menu
			switch challengeView.sourceMenu {
			case MainMenu:
				// If started from main menu, go back to main menu
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			case ProgressMenu:
				// If came from progress, go back to progress
				m.activeView = NewProgressMenu(m.gameState, m.width, m.height)
			case CategoryMenu:
				// If came from category, go back to challenge menu for that category
				m.activeView = NewChallengeMenu(
					m.gameState,
					m.gameState.CurrentCategoryIdx,
					m.width,
					m.height,
					MainMenu,
				)
			default:
				// Default to main menu if source is unknown
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			}
		} else if explanationView, ok := m.activeView.(*ExplanationView); ok {
			// If coming from explanation view, go back to the menu based on source
			switch explanationView.sourceMenu {
			case MainMenu:
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			case ProgressMenu:
				m.activeView = NewProgressMenu(m.gameState, m.width, m.height)
			case CategoryMenu:
				m.activeView = NewCategoryMenu(m.gameState, m.width, m.height, MainMenu)
			default:
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			}
		} else if menuView, ok := m.activeView.(*MenuView); ok {
			// For a MenuView, check its type
			if menuView.type_ == ChallengeMenu {
				// If coming from challenge menu, check its source
				if menuView.sourceMenu == ProgressMenu {
					// If the source was progress, go back to progress
					m.activeView = NewProgressMenu(m.gameState, m.width, m.height)
				} else {
					// Otherwise go back to category menu
					m.activeView = NewCategoryMenu(m.gameState, m.width, m.height, MainMenu)
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
				// If all challenges are complete, go back to main menu
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
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
