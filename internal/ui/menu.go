package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"blindspot/internal/challenges"
	"blindspot/internal/game"
	"blindspot/internal/utils"
)

type MenuType int

const (
	MainMenu MenuType = iota
	CategoryMenu
	ChallengeMenu
	ProgressMenu
	SettingsMenu
)

type SelectCategoryMsg struct {
	CategoryIndex int
}

type SelectChallengeMsg struct {
	Challenge challenges.Challenge
}

type MenuView struct {
	type_       MenuType
	title       string
	items       []MenuItem
	cursor      int
	gameState   *game.GameState
	help        help.Model
	showHelp    bool
	width       int
	height      int
	description string
	sourceMenu  MenuType
	viewport    viewport.Model
	contentStr  string
	keys        MenuKeyMap
	builder     *MenuBuilder
}

func (m *MenuView) Init() tea.Cmd {
	m.updateContent()
	return nil
}

func (m *MenuView) updateContent() {
	m.contentStr = m.builder.BuildMenuContent(m.title, m.description, m.items, m.cursor)

	helpHeight := ShortHelpHeight
	if m.showHelp {
		helpHeight = FullHelpHeight
	}

	contentHeight := strings.Count(m.contentStr, "\n") + 1
	viewportHeight := min(contentHeight, m.height-helpHeight-1)

	m.viewport = viewport.New(m.width, viewportHeight)
	m.viewport.SetContent(m.contentStr)
}

func (m *MenuView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Back):
			return m.handleBackAction()

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			m.updateContent()

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
				m.updateContent()
				m.scrollToCursor()
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.updateContent()
				m.scrollToCursor()
			}

		case key.Matches(msg, m.keys.ScrollUp):
			m.viewport.LineUp(1)

		case key.Matches(msg, m.keys.ScrollDown):
			m.viewport.LineDown(1)

		case key.Matches(msg, m.keys.Select):
			if m.type_ == MainMenu {
				switch m.cursor {
				case 0: // Start Game
					m.gameState.UseRandomizedOrder = m.gameState.Settings.GameMode == "random-by-difficulty"

					if m.gameState.UseRandomizedOrder && len(m.gameState.RandomizedChallenges) == 0 {
						m.gameState.RandomizedChallenges = m.gameState.GetChallengesGroupedByDifficulty()
						m.gameState.SaveRandomizedOrder()
					}

					if m.gameState.ShouldReturnToCategoryExplanation() {
						explanationView := NewExplanationView(m.gameState, nil, m.width, m.height, MainMenu, false)
						return explanationView, nil
					}

					// Check if all challenges are completed
					helpers := game.NewGameStateHelpers(m.gameState)
					totalChallenges := helpers.GetTotalChallengesCount()
					completedChallenges := helpers.GetCompletedChallengesCount()

					if totalChallenges > 0 && completedChallenges >= totalChallenges {
						// All challenges completed, show completion screen
						newCompletion := CompletionViewScreen(m.gameState, m.width, m.height, MainMenu)
						return newCompletion, nil
					}

					// Try to get next incomplete challenge
					challenge, found := m.gameState.GetNextIncompleteChallenge()
					if found {
						m.gameState.ClearMessages()
						return m, func() tea.Msg {
							return SelectChallengeMsg{Challenge: challenge}
						}
					} else {
						// No incomplete challenges found, show completion screen
						newCompletion := CompletionViewScreen(m.gameState, m.width, m.height, MainMenu)
						return newCompletion, nil
					}

				case 1: // Categories
					m.gameState.ClearMessages()
					newMenu := NewCategoriesMenu(m.gameState, m.width, m.height, MainMenu)
					return newMenu, nil
				case 2: // Progress
					m.gameState.ClearMessages()
					newMenu := NewProgressMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				case 3: // Settings
					m.gameState.ClearMessages()
					newMenu := NewSettingsMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				case 4: // Exit
					return m, tea.Quit
				}
			} else if m.type_ == CategoryMenu {
				m.gameState.ClearMessages()
				newMenu := NewCategoryMenu(m.gameState, m.cursor, m.width, m.height, m.type_)
				return newMenu, nil
			} else if m.type_ == ProgressMenu {
				m.gameState.ClearMessages()
				newMenu := NewCategoryMenu(m.gameState, m.cursor, m.width, m.height, m.type_)
				return newMenu, nil
			} else if m.type_ == ChallengeMenu {
				if strings.HasPrefix(m.items[m.cursor].ID, "explanation-") {
					m.gameState.ClearMessages()
					explanationView := NewExplanationView(m.gameState, nil, m.width, m.height, m.type_, false)
					return explanationView, nil
				}

				for _, set := range m.gameState.ChallengeSets {
					for _, challenge := range set.Challenges {
						if challenge.ID == m.items[m.cursor].ID {
							m.gameState.ClearMessages()
							return m, func() tea.Msg {
								return SelectChallengeMsg{Challenge: challenge}
							}
						}
					}
				}
			} else if m.type_ == SettingsMenu {
				if m.cursor == 0 { // Vulnerability Names toggle
					m.gameState.ToggleShowVulnerabilityNames()
					newMenu := NewSettingsMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				} else if m.cursor == 1 { // Challenge Order toggle
					m.gameState.ToggleGameMode()
					newMenu := NewSettingsMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				} else if m.cursor == 2 { // Delete Progress Data
					m.gameState.EraseProgressData()
					newMenu := NewSettingsMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				} else if m.cursor == 3 { // Back to Main Menu
					m.gameState.ClearMessages()
					newMenu := NewMainMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.updateContent()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *MenuView) View() string {
	var b strings.Builder
	b.WriteString(m.viewport.View())

	notificationDisplay := NewNotificationDisplay()
	b.WriteString(notificationDisplay.RenderAllNotifications(m.gameState))

	// Help
	if m.showHelp {
		b.WriteString("\n" + m.help.View(m.keys))
	} else {
		helpText := "Press ? for help | ↑/↓ to navigate"
		if m.width < 60 {
			helpText = "? for help | ↑/↓ nav"
		}
		b.WriteString("\n" + helpHintStyle.Render(helpText))
	}

	return b.String()
}

func NewMainMenu(gs *game.GameState, width, height int) *MenuView {
	items := []MenuItem{
		{Title: "Start Game", Description: "Begin playing from where you left off"},
		{Title: "Categories", Description: "Browse security challenge categories"},
		{Title: "Progress", Description: "View your progress statistics"},
		{Title: "Settings", Description: "Configure game preferences"},
		{Title: "Exit", Description: "Save and exit the game"},
	}

	ascii := `
	 	╭━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╮
	    │                                               │
	    │                                               │
	    │                                               │
	    │       █▄▄ █   █ █▄ █ █▀▄ █▀ █▀█ █▀█ ▀█▀       │
	    │       █▄█ █▄▄ █ █ ▀█ █▄▀ ▄█ █▀▀ █▄█  █        │
	    │                                               │
	    │                                               │
	    │                                               │
	    ╰━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╯

	             ✧ find insecure code practices ✧
	         ✧ level up your code security training ✧
			 

	・．．・゜゜・．．・゜゜・．．・゜゜・．．・゜゜・．．・゜゜
	 	`

	menu := &MenuView{
		type_:       MainMenu,
		title:       ascii,
		items:       items,
		gameState:   gs,
		help:        help.New(),
		width:       width,
		height:      height,
		description: "Train your eye to find and fix insecure coding practices through challenges!\nIdentify common security vulnerabilities based on the OWASP Top 10.",
		keys:        NewMenuKeyMap(),
		builder:     NewMenuBuilder(width, height),
	}

	menu.updateContent()
	return menu
}

func NewCategoriesMenu(gs *game.GameState, width, height int, source MenuType) *MenuView {
	builder := NewMenuBuilder(width, height)
	items := builder.BuildCategoryItems(gs)

	menu := &MenuView{
		type_:       CategoryMenu,
		title:       "Challenge Categories",
		items:       items,
		gameState:   gs,
		help:        help.New(),
		width:       width,
		height:      height,
		description: "Select any category to view its challenges.\nCategories contain challenges of various difficulty levels based on the OWASP Top 10.",
		sourceMenu:  source,
		keys:        NewMenuKeyMap(),
		builder:     builder,
	}

	menu.updateContent()
	return menu
}

func NewCategoryMenu(gs *game.GameState, categoryIndex int, width, height int, source MenuType) *MenuView {
	category := gs.ChallengeSets[categoryIndex]
	builder := NewMenuBuilder(width, height)
	items := builder.BuildChallengeItems(gs, category)

	menu := &MenuView{
		type_:       ChallengeMenu,
		title:       fmt.Sprintf("%s Challenges", category.Category),
		items:       items,
		gameState:   gs,
		help:        help.New(),
		width:       width,
		height:      height,
		description: utils.WrapText(category.Description, width),
		sourceMenu:  source,
		keys:        NewMenuKeyMap(),
		builder:     builder,
	}

	menu.updateContent()
	return menu
}

func NewProgressMenu(gs *game.GameState, width, height int) *MenuView {
	builder := NewMenuBuilder(width, height)
	items := builder.BuildProgressItems(gs)

	helpers := game.NewGameStateHelpers(gs)
	overallPercentage := helpers.GetOverallCompletionPercentage()
	totalChallenges := helpers.GetTotalChallengesCount()
	completedChallenges := helpers.GetCompletedChallengesCount()

	description := fmt.Sprintf("Overall Progress: %d of %d challenges completed (%d%%)\n\n",
		completedChallenges, totalChallenges, overallPercentage)
	description += "Select a category to see detailed progress statistics.\nPress Enter on a category to view its challenges."

	menu := &MenuView{
		type_:       ProgressMenu,
		title:       "Your Progress",
		items:       items,
		gameState:   gs,
		help:        help.New(),
		width:       width,
		height:      height,
		description: utils.WrapText(description, width),
		sourceMenu:  MainMenu,
		keys:        NewMenuKeyMap(),
		builder:     builder,
	}

	menu.updateContent()
	return menu
}

func NewSettingsMenu(gs *game.GameState, width, height int) *MenuView {
	builder := NewMenuBuilder(width, height)
	items := builder.BuildSettingsItems(gs)

	menu := &MenuView{
		type_:       SettingsMenu,
		title:       "Game Settings",
		items:       items,
		gameState:   gs,
		help:        help.New(),
		width:       width,
		height:      height,
		description: "Configure your game preferences. These settings will be saved for future sessions.",
		sourceMenu:  MainMenu,
		keys:        NewMenuKeyMap(),
		builder:     builder,
	}

	menu.updateContent()
	return menu
}

func (k MenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k MenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.ScrollUp, k.ScrollDown},
		{k.Select, k.Back},
		{k.Help, k.Quit},
	}
}

// Helper methods for MenuView
func (m *MenuView) handleBackAction() (tea.Model, tea.Cmd) {
	m.gameState.ClearMessages()

	switch m.type_ {
	case CategoryMenu:
		return NewMainMenu(m.gameState, m.width, m.height), nil
	case ChallengeMenu:
		if m.sourceMenu == ProgressMenu {
			return NewProgressMenu(m.gameState, m.width, m.height), nil
		}
		return NewCategoriesMenu(m.gameState, m.width, m.height, MainMenu), nil
	case ProgressMenu, SettingsMenu:
		return NewMainMenu(m.gameState, m.width, m.height), nil
	default:
		return NewMainMenu(m.gameState, m.width, m.height), nil
	}
}

func (m *MenuView) scrollToCursor() {
	cursorPos := strings.Index(m.contentStr, ">")
	if cursorPos > -1 {
		linesBefore := strings.Count(m.contentStr[:cursorPos], "\n")
		if linesBefore > m.viewport.Height/2 {
			m.viewport.SetYOffset(linesBefore - m.viewport.Height/2)
		}
	}
}
