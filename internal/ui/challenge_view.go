package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"secure-code-game/internal/challenges"
	"secure-code-game/internal/game"
	"secure-code-game/internal/utils"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Select   key.Binding
	Back     key.Binding
	Help     key.Binding
	Quit     key.Binding
	ShowHint key.Binding
	Next     key.Binding
}

var keys = keyMap{
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
		key.WithKeys("esc"),
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
	ShowHint: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "show hint"),
	),
	Next: key.NewBinding(
		key.WithKeys("n", "enter"),
		key.WithHelp("n", "next challenge"),
		key.WithHelp("enter", "next challenge"),
	),
}

// displays a multiple choice vulnerability challenge
type ChallengeView struct {
	gameState   *game.GameState
	challenge   challenges.Challenge
	cursor      int
	selected    bool
	showHint    bool
	helpModel   help.Model
	quizOptions []string
	result      string
	resultStyle lipgloss.Style
	width       int
	height      int
	hasAnswered bool
	isCorrect   bool
	sourceMenu  MenuType
}

// NewChallengeView creates a new challenge view
func NewChallengeView(gs *game.GameState, challenge challenges.Challenge, width, height int, source MenuType) *ChallengeView {
	helpModel := help.New()
	helpModel.Width = 80

	return &ChallengeView{
		gameState:   gs,
		challenge:   challenge,
		helpModel:   helpModel,
		resultStyle: successStyle,
		showHint:    false,
		cursor:      0,
		selected:    false,
		width:       width,
		height:      height,
		hasAnswered: false,
		isCorrect:   false,
		quizOptions: challenge.Options,
		sourceMenu:  source,
	}
}

// Init initializes the challenge view
func (m *ChallengeView) Init() tea.Cmd {
	return nil
}

// Update handles messages and user input
func (m *ChallengeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If we've already answered correctly and goes to next challenge
		if m.hasAnswered && m.isCorrect && key.Matches(msg, keys.Next) {
			// Show category explanation before going to next challenge
			explanationView := NewExplanationView(m.gameState, m.challenge, m.width, m.height, m.sourceMenu, true)
			return explanationView, explanationView.Init()
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Back):
			// Return to menu
			return m, func() tea.Msg {
				return backToMenuMsg{}
			}

		case key.Matches(msg, keys.Help):
			m.helpModel.ShowAll = !m.helpModel.ShowAll

		case key.Matches(msg, keys.ShowHint):
			m.showHint = !m.showHint

		case key.Matches(msg, keys.Up):
			// Can navigate up/down even after a wrong answer
			if m.cursor > 0 && (!m.hasAnswered || !m.isCorrect) {
				m.cursor--
				// Reset the result message for another try
				if m.hasAnswered && !m.isCorrect {
					m.hasAnswered = false
					m.result = ""
				}
			}

		case key.Matches(msg, keys.Down):
			// Can navigate up/down even after a wrong answer
			if m.cursor < len(m.challenge.Options)-1 && (!m.hasAnswered || !m.isCorrect) {
				m.cursor++
				// Reset the result message for another try
				if m.hasAnswered && !m.isCorrect {
					m.hasAnswered = false
					m.result = ""
				}
			}

		case key.Matches(msg, keys.Select):
			// Check if the answer is correct
			selectedOption := m.challenge.Options[m.cursor]
			m.hasAnswered = true
			if selectedOption == m.challenge.CorrectAnswer {
				m.isCorrect = true
				m.result = "✓ Correct! You've identified the vulnerability."
				m.resultStyle = successStyle

				// Mark challenge as completed
				m.gameState.MarkChallengeCompleted(m.challenge.ID)

				if m.gameState.ShouldShowVulnerabilityExplanation(m.challenge.Category) {
					// Show explanation view immediately after marking correct if not in random mode
					explanationView := NewExplanationView(m.gameState, m.challenge, m.width, m.height, m.sourceMenu, true)
					return explanationView, explanationView.Init()
				}
			} else {
				m.isCorrect = false
				m.result = "✗ Incorrect. Try another option by moving arrow keys!"
				m.resultStyle = errorStyle
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the challenge
func (m *ChallengeView) View() string {
	var b strings.Builder

	// Challenge title with difficulty indicator
	difficultyText := ""
	switch m.challenge.Difficulty {
	case challenges.Beginner:
		difficultyText = difficultyStyle["beginner"].Render("[Beginner]")
	case challenges.Intermediate:
		difficultyText = difficultyStyle["intermediate"].Render("[Intermediate]")
	case challenges.Advanced:
		difficultyText = difficultyStyle["advanced"].Render("[Advanced]")
	}

	b.WriteString(difficultyText + "\n\n")

	showCategory := m.gameState.Settings.ShowVulnerabilityNames || (m.hasAnswered && m.isCorrect)

	if showCategory {
		b.WriteString(subtitleStyle.Render("CHALLENGE NAME:") + "\n")
		b.WriteString(titleStyle.Render(m.challenge.Title) + "\n")

		categoryHeader := fmt.Sprintf("VULNERABILITY TYPE: %s", m.challenge.Category)
		b.WriteString(categoryStyle.Render(categoryHeader) + "\n\n")
	}

	// Create a separator for visual clarity
	separator := strings.Repeat("─", m.width/2)
	b.WriteString(subtleStyle.Render(separator) + "\n\n")

	// Detect language from the challenge code to highlight
	language := ""
	if len(m.challenge.Lang) > 0 {
		lang := utils.GetLanguageFromExtension(m.challenge.Lang)
		if lang != "" {
			language = lang
		}
	}

	// Apply syntax highlighting to the code
	highlightedCode := utils.HighlightCode(m.challenge.Code, language)

	// b.WriteString("Vulnerable Code:\n")
	b.WriteString(codeBoxStyle.Render(highlightedCode) + "\n\n")

	b.WriteString(descStyle.Render("What vulnerability is in this code?") + "\n\n")

	for i, option := range m.challenge.Options {
		var renderedOption string
		cursor := "  "

		if m.hasAnswered && m.isCorrect {
			// Only show the correct answer when the user get it right
			if option == m.challenge.CorrectAnswer {
				cursor = "✓ "
				successStyleCopy := successStyle
				successStyleCopy = successStyleCopy.Bold(false)
				renderedOption = cursor + successStyleCopy.Render(option)
			} else {
				renderedOption = cursor + unselectedItemStyle.Render(option)
			}
		} else if m.hasAnswered && !m.isCorrect && i == m.cursor {
			// Mark the user's wrong selection with an X
			cursor = "✗ "
			errorStyleCopy := errorStyle
			errorStyleCopy = errorStyleCopy.Bold(false)
			renderedOption = cursor + errorStyleCopy.Render(option)
		} else if !m.hasAnswered {
			// Not answered yet - normal cursor
			if m.cursor == i {
				cursor = "> "
				renderedOption = cursor + selectedItemStyle.Render(option)
			} else {
				renderedOption = cursor + unselectedItemStyle.Render(option)
			}
		} else {
			// Default case - keep normal styling
			renderedOption = cursor + unselectedItemStyle.Render(option)
		}

		b.WriteString(renderedOption + "\n")
	}

	// Hint section
	if m.showHint {
		b.WriteString("\n" + hintStyle.Render("Hint: "+m.challenge.Hint) + "\n")
	}

	// Result of submission
	if m.result != "" {
		b.WriteString("\n" + m.resultStyle.Render(m.result) + "\n")

		// Show "next challenge" prompt if the answer is correct
		if m.isCorrect {
			b.WriteString("\n" + helpHintStyle.Render("Press 'Enter'/'N' to continue to next challenge"))
		}
	}

	// Help section
	helpText := m.helpModel.View(keys)
	if !m.helpModel.ShowAll {
		helpText = helpHintStyle.Render("Press ? for help")
	}
	b.WriteString("\n" + helpStyle.Render(helpText))

	return b.String()
}

// --- helpers ---
// returns keybindings to be shown in the mini help view.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select}, // First column
		{k.Back, k.Help, k.Quit}, // Second column
		{k.ShowHint, k.Next},     // Third column
	}
}
