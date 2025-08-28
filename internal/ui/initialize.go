package ui

import (
	"blindspot/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	gameState     *game.GameState
	activeView    tea.Model
	width, height int
}

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

func (m AppModel) Init() tea.Cmd {
	return m.activeView.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.activeView, cmd = m.activeView.Update(msg)
		return m, cmd

	case SelectChallengeMsg:
		m.activeView = NewChallengeView(m.gameState, msg.Challenge, m.width, m.height, MainMenu)
		return m, m.activeView.Init()

	case backToMenuMsg:
		m.gameState.ClearMessages()

		if challengeView, ok := m.activeView.(*ChallengeView); ok {
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
			if menuView.type_ == ChallengeMenu {
				if menuView.sourceMenu == ProgressMenu {
					m.activeView = NewProgressMenu(m.gameState, m.width, m.height)
				} else {
					m.activeView = NewCategoriesMenu(m.gameState, m.width, m.height, MainMenu)
				}
			} else {
				m.activeView = NewMainMenu(m.gameState, m.width, m.height)
			}
		} else {
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
				m.activeView = CompletionViewScreen(m.gameState, m.width, m.height, MainMenu)
			}
		}
		return m, m.activeView.Init()
	}

	m.activeView, cmd = m.activeView.Update(msg)
	return m, cmd
}

func (m AppModel) View() string {
	return m.activeView.View()
}
