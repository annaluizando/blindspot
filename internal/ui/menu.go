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

type MenuKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Select     key.Binding
	Back       key.Binding
	Help       key.Binding
	Quit       key.Binding
}

var MenuKeys = MenuKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "scroll up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "scroll down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}

type MenuItem struct {
	Title       string
	Description string
	Completed   bool
	ID          string
}

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
}

func (m *MenuView) Init() tea.Cmd {
	m.updateContent()
	return nil
}

// handles messages and user input
func (m *MenuView) updateContent() {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.title) + "\n\n")
	wrappedDescription := utils.WrapText(m.description, m.width)
	b.WriteString(descriptionStyle.Render(wrappedDescription) + "\n\n")

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		status := " "
		if item.Completed {
			status = "✓"
		}

		line := fmt.Sprintf("%s %s %s", cursor, status, item.Title)
		if m.cursor == i {
			b.WriteString(selectedItemStyle.Render(line) + "\n")

			wrappedItemDescription := utils.WrapText(item.Description, m.width)
			b.WriteString(itemDescriptionStyle.Render(wrappedItemDescription) + "\n\n")
		} else {
			b.WriteString(itemStyle.Render(line) + "\n")
		}
	}

	// Calculate space for help
	helpHeight := 1
	if m.showHelp {
		helpHeight = 4 // Full help takes more space
	}

	// Create or update the viewport
	m.contentStr = b.String()
	contentHeight := strings.Count(m.contentStr, "\n") + 1

	// If content is shorter than available height, no scrolling needed
	viewportHeight := min(contentHeight, m.height-helpHeight-1)

	m.viewport = viewport.New(m.width, viewportHeight)
	m.viewport.SetContent(m.contentStr)
}

func (m *MenuView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, MenuKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, MenuKeys.Back):
			if m.type_ == CategoryMenu {
				m.gameState.ClearMessages()
				newMenu := NewMainMenu(m.gameState, m.width, m.height)
				return newMenu, nil
			} else if m.type_ == ChallengeMenu {
				m.gameState.ClearMessages()
				if m.sourceMenu == ProgressMenu {
					newMenu := NewProgressMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				}
				newMenu := NewCategoriesMenu(m.gameState, m.width, m.height, MainMenu)
				return newMenu, nil
			} else if m.type_ == ProgressMenu {
				m.gameState.ClearMessages()
				newMenu := NewMainMenu(m.gameState, m.width, m.height)
				return newMenu, nil
			} else if m.type_ == SettingsMenu {
				m.gameState.ClearMessages()
				newMenu := NewMainMenu(m.gameState, m.width, m.height)
				return newMenu, nil
			}

		case key.Matches(msg, MenuKeys.Help):
			m.showHelp = !m.showHelp
			m.updateContent()

		case key.Matches(msg, MenuKeys.Up):
			if m.cursor > 0 {
				m.cursor--
				m.updateContent()
				// keep the cursor visible
				cursorPos := strings.Index(m.contentStr, ">")
				if cursorPos > -1 {
					m.viewport.SetYOffset(0) // Reset to top first
					linesBefore := strings.Count(m.contentStr[:cursorPos], "\n")
					if linesBefore > m.viewport.Height/2 {
						m.viewport.SetYOffset(linesBefore - m.viewport.Height/2)
					}
				}
			}

		case key.Matches(msg, MenuKeys.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.updateContent()
				// keep the cursor visible
				cursorPos := strings.Index(m.contentStr, ">")
				if cursorPos > -1 {
					m.viewport.GotoTop()
					linesBefore := strings.Count(m.contentStr[:cursorPos], "\n")
					if linesBefore > m.viewport.Height/2 {
						m.viewport.SetYOffset(linesBefore - m.viewport.Height/2)
					}
				}
			}

		case key.Matches(msg, MenuKeys.ScrollUp):
			m.viewport.LineUp(1)

		case key.Matches(msg, MenuKeys.ScrollDown):
			m.viewport.LineDown(1)

		case key.Matches(msg, MenuKeys.Select):
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

					_, found := m.gameState.GetNextIncompleteChallenge()
					var challenge challenges.Challenge

					if found {
						challenge = m.gameState.GetCurrentChallenge()
						m.gameState.ClearMessages()
						return m, func() tea.Msg {
							return SelectChallengeMsg{Challenge: challenge}
						}
					} else {
						newCompletion := NewCompletionView(m.gameState, m.width, m.height, MainMenu)
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

	// Handle viewport updates
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
		b.WriteString("\n" + m.help.View(MenuKeys))
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
	}

	menu.updateContent()
	return menu
}

func NewCategoriesMenu(gs *game.GameState, width, height int, source MenuType) *MenuView {
	items := make([]MenuItem, len(gs.ChallengeSets))

	for i, set := range gs.ChallengeSets {
		completed := gs.GetCategoryCompletionPercentage(set.Category)
		completionText := fmt.Sprintf("[%d%% Complete]", completed)

		hasBeginner := false
		hasIntermediate := false
		hasAdvanced := false

		for _, challenge := range set.Challenges {
			switch challenge.Difficulty {
			case challenges.Beginner:
				hasBeginner = true
			case challenges.Intermediate:
				hasIntermediate = true
			case challenges.Advanced:
				hasAdvanced = true
			}
		}

		difficultyIndicator := ""
		if hasBeginner {
			difficultyIndicator += difficultyStyle["beginner"].Render("[B]") + " "
		}
		if hasIntermediate {
			difficultyIndicator += difficultyStyle["intermediate"].Render("[I]") + " "
		}
		if hasAdvanced {
			difficultyIndicator += difficultyStyle["advanced"].Render("[A]") + " "
		}

		enhancedDescription := utils.WrapText(set.Description, width) + "\n" + difficultyIndicator + completionText

		items[i] = MenuItem{
			Title:       set.Category,
			Description: enhancedDescription,
			Completed:   completed == 100,
			ID:          fmt.Sprintf("category-%d", i),
		}
	}

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
	}

	menu.updateContent()
	return menu
}

func NewCategoryMenu(gs *game.GameState, categoryIndex int, width, height int, source MenuType) *MenuView {
	category := gs.ChallengeSets[categoryIndex]

	var items []MenuItem

	items = append(items, MenuItem{
		Title:       "📚 See Explanation: " + category.Category,
		Description: "View detailed explanation about this vulnerability type, its impact, and prevention techniques.",
		ID:          "explanation-" + category.Category,
	})

	for _, challenge := range category.Challenges {
		completed := gs.IsChallengeCompleted(challenge.ID)
		difficultyText := ""
		switch challenge.Difficulty {
		case challenges.Beginner:
			difficultyText = difficultyStyle["beginner"].Render("[Beginner]")
		case challenges.Intermediate:
			difficultyText = difficultyStyle["intermediate"].Render("[Intermediate]")
		case challenges.Advanced:
			difficultyText = difficultyStyle["advanced"].Render("[Advanced]")
		}

		status := ""
		if completed {
			status = completionStyle.Render("[✓ Completed]")
		} else {
			status = "[Not Completed]"
		}

		items = append(items, MenuItem{
			Title:       challenge.Title,
			Description: fmt.Sprintf("%s %s\n%s", difficultyText, status, challenge.Description),
			Completed:   completed,
			ID:          challenge.ID,
		})
	}

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
	}

	menu.updateContent()
	return menu
}

func NewProgressMenu(gs *game.GameState, width, height int) *MenuView {
	items := make([]MenuItem, len(gs.ChallengeSets))

	var totalChallenges, completedChallenges int

	for i, set := range gs.ChallengeSets {
		completed := gs.GetCategoryCompletionPercentage(set.Category)

		categoryCompleted := 0
		for _, challenge := range set.Challenges {
			totalChallenges++
			if gs.IsChallengeCompleted(challenge.ID) {
				completedChallenges++
				categoryCompleted++
			}
		}

		description := fmt.Sprintf("%d of %d challenges completed (%d%%)\n",
			categoryCompleted, len(set.Challenges), completed)

		hasDifficultyBreakdown := false
		beginnerCount, intermediateCount, advancedCount := 0, 0, 0
		beginnerCompleted, intermediateCompleted, advancedCompleted := 0, 0, 0

		for _, challenge := range set.Challenges {
			switch challenge.Difficulty {
			case challenges.Beginner:
				beginnerCount++
				hasDifficultyBreakdown = true
				if gs.IsChallengeCompleted(challenge.ID) {
					beginnerCompleted++
				}
			case challenges.Intermediate:
				intermediateCount++
				hasDifficultyBreakdown = true
				if gs.IsChallengeCompleted(challenge.ID) {
					intermediateCompleted++
				}
			case challenges.Advanced:
				advancedCount++
				hasDifficultyBreakdown = true
				if gs.IsChallengeCompleted(challenge.ID) {
					advancedCompleted++
				}
			}
		}

		if hasDifficultyBreakdown {
			description += "By Difficulty:\n"

			if beginnerCount > 0 {
				description += fmt.Sprintf("    Beginner: %d/%d completed\n",
					beginnerCompleted, beginnerCount)
			}
			if intermediateCount > 0 {
				description += fmt.Sprintf("    Intermediate: %d/%d completed\n",
					intermediateCompleted, intermediateCount)
			}
			if advancedCount > 0 {
				description += fmt.Sprintf("    Advanced: %d/%d completed\n",
					advancedCompleted, advancedCount)
			}
		}

		categoryErrorCount := 0
		if gs.Progress.CategoryErrorCounts != nil {
			categoryErrorCount = gs.Progress.CategoryErrorCounts[set.Category]
		}

		if categoryErrorCount > 0 {
			errorRate := 0
			totalAttempts := categoryCompleted + categoryErrorCount
			if totalAttempts > 0 {
				errorRate = (categoryErrorCount * 100) / totalAttempts
			}

			var errorLevel string
			if errorRate > 50 {
				errorLevel = "High"
			} else if errorRate > 30 {
				errorLevel = "Moderate"
			} else if errorRate > 15 {
				errorLevel = "Low"
			} else {
				errorLevel = ""
			}

			if errorLevel != "" {
				description += fmt.Sprintf("Errors in category: %d (%s - %d%% error rate)\n",
					categoryErrorCount, errorLevel, errorRate)
			} else {
				description += fmt.Sprintf("Errors in category: %d (%d%% error rate)\n",
					categoryErrorCount, errorRate)
			}
		} else if categoryCompleted > 0 {
			description += "No errors in this category. Great job!\n"
		}

		items[i] = MenuItem{
			Title:       set.Category,
			Description: description,
			Completed:   completed == 100,
			ID:          fmt.Sprintf("progress-category-%d", i),
		}
	}

	overallPercentage := 0
	if totalChallenges > 0 {
		overallPercentage = (completedChallenges * 100) / totalChallenges
	}

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
	}

	menu.updateContent()
	return menu
}

func NewSettingsMenu(gs *game.GameState, width, height int) *MenuView {
	vulnerabilityNamesStatus := "Show"
	if !gs.Settings.ShowVulnerabilityNames {
		vulnerabilityNamesStatus = "Hide"
	}

	orderModeText := "Category Order"
	if gs.Settings.GameMode == "random-by-difficulty" {
		orderModeText = "Random by Difficulty"
	}

	items := []MenuItem{
		{
			Title:       "Vulnerability Names: " + vulnerabilityNamesStatus,
			Description: "Toggle whether vulnerability names are shown during challenges.",
			ID:          "setting-vulnnames",
		},
		{
			Title: "Game Mode: " + orderModeText,
			Description: "Choose how challenges are ordered when playing the game.\n" +
				"Category Order: Play challenges grouped by vulnerabilty category. (Standard Mode)\n" +
				"Random by Difficulty: Play challenges in random order but grouped by difficulty level. (More challenging mode, specially if combined with 'Vulnerability Names: Hide')",
			ID: "setting-ordermode",
		},
		{
			Title: "Delete all progress data",
			Description: "Erases ALL progress data and begin game from start.\n" +
				"!!! Be aware this will make you loose ALL your current progress. \n",
			ID: "setting-deleteprogress",
		},
		{
			Title:       "Back to Main Menu",
			Description: "Return to the main menu",
			ID:          "setting-back",
		},
	}

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
	}

	menu.updateContent()
	return menu
}

// ---- helpers ----
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
