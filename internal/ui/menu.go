package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"secure-code-game/internal/challenges"
	"secure-code-game/internal/game"
	"secure-code-game/internal/utils"
)

// MenuType represents different types of menus
type MenuType int

const (
	MainMenu MenuType = iota
	CategoryMenu
	ChallengeMenu
	ProgressMenu
	SettingsMenu
)

// defines keybindings for menus
type MenuKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// reduced help view.
func (k MenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// expanded help view.
func (k MenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},     // Navigation
		{k.Select, k.Back}, // Actions
		{k.Help, k.Quit},
	}
}

// holds the key mappings for menus
var MenuKeys = MenuKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
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

// sent when a category is selected
type SelectCategoryMsg struct {
	CategoryIndex int
}

// sent when a challenge is selected
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
}

// main menu with options
func NewMainMenu(gs *game.GameState, width, height int) *MenuView {
	items := []MenuItem{
		{Title: "Start Game", Description: "Begin playing from where you left off"},
		{Title: "Categories", Description: "Browse security challenge categories"},
		{Title: "Progress", Description: "View your progress statistics"},
		{Title: "Settings", Description: "Configure game preferences"},
		{Title: "Exit", Description: "Save and exit the game"},
	}

	return &MenuView{
		type_:       MainMenu,
		title:       "Security Code Game",
		items:       items,
		gameState:   gs,
		help:        help.New(),
		width:       width,
		height:      height,
		description: utils.WrapText("Train your eye to find and fix insecure coding practices through challenges!\nIdentify common security vulnerabilities based on the OWASP Top 10.", width),
	}
}

// New category menu for categories display
func NewCategoryMenu(gs *game.GameState, width, height int, source MenuType) *MenuView {
	items := make([]MenuItem, len(gs.ChallengeSets))

	for i, set := range gs.ChallengeSets {
		// completion percentage calculation
		completed := gs.GetCategoryCompletionPercentage(set.Category)
		completionText := fmt.Sprintf("[%d%% Complete]", completed)

		// difficulty levels/indicator
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

		// full category information for "Categories" containing: description, difficulties and completion percentage.
		enhancedDescription := utils.WrapText(set.Description, width) + "\n" + difficultyIndicator + completionText

		items[i] = MenuItem{
			Title:       set.Category,
			Description: enhancedDescription,
			Completed:   completed == 100,
			ID:          fmt.Sprintf("category-%d", i),
		}
	}

	return &MenuView{
		type_:       CategoryMenu,
		title:       "Challenge Categories",
		items:       items,
		gameState:   gs,
		help:        help.New(),
		width:       width,
		height:      height,
		description: "Select any category to view its challenges. Categories contain challenges of various difficulty levels based on the OWASP Top 10.",
		sourceMenu:  source,
	}
}

func NewChallengeMenu(gs *game.GameState, categoryIndex int, width, height int, source MenuType) *MenuView {
	category := gs.ChallengeSets[categoryIndex]

	// Create items array - add space for explanation option if available
	var items []MenuItem

	items = append(items, MenuItem{
		Title:       "📚 See Explanation: " + category.Category,
		Description: "View detailed explanation about this vulnerability type, its impact, and prevention techniques.",
		ID:          "explanation-" + category.Category,
	})

	// Add all challenges for this category
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

	return &MenuView{
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
}

// Creates a new progress view showing completion across all categories
func NewProgressMenu(gs *game.GameState, width, height int) *MenuView {
	items := make([]MenuItem, len(gs.ChallengeSets))

	var totalChallenges, completedChallenges int

	// Create menu items for each category with completion details
	for i, set := range gs.ChallengeSets {
		completed := gs.GetCategoryCompletionPercentage(set.Category)

		// Count challenges for overall statistics
		categoryCompleted := 0
		for _, challenge := range set.Challenges {
			totalChallenges++
			if gs.IsChallengeCompleted(challenge.ID) {
				completedChallenges++
				categoryCompleted++
			}
		}

		// Build a detailed description with completion information
		description := fmt.Sprintf("%d of %d challenges completed (%d%%)\n",
			categoryCompleted, len(set.Challenges), completed)

		// Add difficulty distribution information
		beginnerCount, intermediateCount, advancedCount := 0, 0, 0
		beginnerCompleted, intermediateCompleted, advancedCompleted := 0, 0, 0

		for _, challenge := range set.Challenges {
			switch challenge.Difficulty {
			case challenges.Beginner:
				beginnerCount++
				if gs.IsChallengeCompleted(challenge.ID) {
					beginnerCompleted++
				}
			case challenges.Intermediate:
				intermediateCount++
				if gs.IsChallengeCompleted(challenge.ID) {
					intermediateCompleted++
				}
			case challenges.Advanced:
				advancedCount++
				if gs.IsChallengeCompleted(challenge.ID) {
					advancedCompleted++
				}
			}
		}

		// Add difficulty breakdown to description
		if beginnerCount > 0 {
			description += fmt.Sprintf("Beginner: %d/%d completed\n",
				beginnerCompleted, beginnerCount)
		}
		if intermediateCount > 0 {
			description += fmt.Sprintf("Intermediate: %d/%d completed\n",
				intermediateCompleted, intermediateCount)
		}
		if advancedCount > 0 {
			description += fmt.Sprintf("Advanced: %d/%d completed\n",
				advancedCompleted, advancedCount)
		}

		items[i] = MenuItem{
			Title:       set.Category,
			Description: description,
			Completed:   completed == 100,
			ID:          fmt.Sprintf("progress-category-%d", i),
		}
	}

	// Overall completion percentage
	overallPercentage := 0
	if totalChallenges > 0 {
		overallPercentage = (completedChallenges * 100) / totalChallenges
	}

	description := fmt.Sprintf("Overall Progress: %d of %d challenges completed (%d%%)\n\n",
		completedChallenges, totalChallenges, overallPercentage)
	description += "Select a category to see detailed progress statistics. Press Enter on a category to view its challenges."

	return &MenuView{
		type_:       ProgressMenu,
		title:       "Your Progress",
		items:       items,
		gameState:   gs,
		help:        help.New(),
		width:       width,
		height:      height,
		description: description,
		sourceMenu:  MainMenu,
	}
}

// SettingsMenu is the menu component for game settings
// menu component for game settings
func NewSettingsMenu(gs *game.GameState, width, height int) *MenuView {
	vulnerabilityNamesStatus := "Show"
	if !gs.Settings.ShowVulnerabilityNames {
		vulnerabilityNamesStatus = "Hide"
	}

	orderModeText := "Category Order"
	if gs.Settings.ChallengeOrderMode == "random-by-difficulty" {
		orderModeText = "Random by Difficulty"
	}

	items := []MenuItem{
		{
			Title:       "Vulnerability Names: " + vulnerabilityNamesStatus,
			Description: "Toggle whether vulnerability names are shown during challenges.",
			ID:          "setting-vulnnames",
		},
		{
			Title: "Challenge Order: " + orderModeText,
			Description: "Choose how challenges are ordered when playing the game.\n" +
				"Category Order: Play challenges grouped by vulnerabilty category.\n" +
				"Random by Difficulty: Play challenges in random order but grouped by difficulty level (beginner, intermediate, advanced).",
			ID: "setting-ordermode",
		},
		{
			Title:       "Back to Main Menu",
			Description: "Return to the main menu",
			ID:          "setting-back",
		},
	}

	return &MenuView{
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
}

func (m *MenuView) Init() tea.Cmd {
	return nil
}

// Update handles messages and user input
func (m *MenuView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, MenuKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, MenuKeys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, MenuKeys.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case key.Matches(msg, MenuKeys.Help):
			m.showHelp = !m.showHelp

		case key.Matches(msg, MenuKeys.Back):
			// Handle back navigation based on current menu type
			if m.type_ == CategoryMenu {
				// Go back to main menu from categories
				newMenu := NewMainMenu(m.gameState, m.width, m.height)
				return newMenu, nil
			} else if m.type_ == ChallengeMenu {
				// For challenge menu, check the source
				if m.sourceMenu == ProgressMenu {
					// If came from progress, go back to progress
					newMenu := NewProgressMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				}
				newMenu := NewCategoryMenu(m.gameState, m.width, m.height, MainMenu)
				return newMenu, nil
			} else if m.type_ == ProgressMenu {
				// When in progress menu, go back to main menu
				newMenu := NewMainMenu(m.gameState, m.width, m.height)
				return newMenu, nil
			} else if m.type_ == SettingsMenu {
				// When in settings menu, go back to main menu
				newMenu := NewMainMenu(m.gameState, m.width, m.height)
				return newMenu, nil
			}

		case key.Matches(msg, MenuKeys.Select):
			if m.type_ == MainMenu {
				switch m.cursor {
				case 0: // Start Game
					m.gameState.UseRandomizedOrder = m.gameState.Settings.ChallengeOrderMode == "random-by-difficulty"

					// if in random mode and don't have randomized challenges yet, generate them
					if m.gameState.UseRandomizedOrder && len(m.gameState.RandomizedChallenges) == 0 {
						m.gameState.RandomizedChallenges = m.gameState.GetRandomizedChallengesByDifficulty()
						m.gameState.SaveRandomizedOrder()
					}

					challenge := m.gameState.GetCurrentChallenge()

					// Go directly to current challenge
					return m, func() tea.Msg {
						return SelectChallengeMsg{Challenge: challenge}
					}
				case 1: // Categories
					newMenu := NewCategoryMenu(m.gameState, m.width, m.height, MainMenu)
					return newMenu, nil
				case 2: // Progress
					newMenu := NewProgressMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				case 3: // Settings
					newMenu := NewSettingsMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				case 4: // Exit
					return m, tea.Quit
				}
			} else if m.type_ == CategoryMenu {
				newMenu := NewChallengeMenu(m.gameState, m.cursor, m.width, m.height, m.type_)
				return newMenu, nil
			} else if m.type_ == ProgressMenu {
				// Create a new challenge menu when coming from Progress Menu
				newMenu := NewChallengeMenu(m.gameState, m.cursor, m.width, m.height, m.type_)
				return newMenu, nil
			} else if m.type_ == ChallengeMenu {
				if strings.HasPrefix(m.items[m.cursor].ID, "explanation-") {
					categoryName := strings.TrimPrefix(m.items[m.cursor].ID, "explanation-")

					var sampleChallenge challenges.Challenge
					for _, set := range m.gameState.ChallengeSets {
						if set.Category == categoryName && len(set.Challenges) > 0 {
							sampleChallenge = set.Challenges[0]
							break
						}
					}

					// Create and return explanation view
					explanationView := NewExplanationView(m.gameState, sampleChallenge, m.width, m.height, m.type_, false)
					return explanationView, nil
				}

				// Handle regular challenge selection (existing code)
				for _, set := range m.gameState.ChallengeSets {
					for _, challenge := range set.Challenges {
						if challenge.ID == m.items[m.cursor].ID {
							return m, func() tea.Msg {
								return SelectChallengeMsg{Challenge: challenge}
							}
						}
					}
				}
			} else if m.type_ == SettingsMenu {
				if m.cursor == 0 { // Vulnerability Names toggle
					// Toggle the setting
					m.gameState.ToggleShowVulnerabilityNames()

					// Create a new settings menu to refresh the display
					newMenu := NewSettingsMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				} else if m.cursor == 1 { // Challenge Order toggle
					// Toggle the challenge order mode
					m.gameState.ToggleChallengeOrderMode()

					// Create a new settings menu to refresh the display
					newMenu := NewSettingsMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				} else if m.cursor == 2 { // Back to Main Menu
					newMenu := NewMainMenu(m.gameState, m.width, m.height)
					return newMenu, nil
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the menu
func (m *MenuView) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.title) + "\n\n")

	wrappedDescription := utils.WrapText(m.description, m.width)
	b.WriteString(descriptionStyle.Render(wrappedDescription) + "\n\n")

	// Menu items
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

	// Help
	if m.showHelp {
		b.WriteString("\n" + m.help.View(MenuKeys))
	} else {
		b.WriteString("\n" + helpHintStyle.Render("Press ? for help"))
	}

	return b.String()
}
