package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"blindspot/internal/challenges"
	"blindspot/internal/game"
	"blindspot/internal/utils"
)

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Select     key.Binding
	Back       key.Binding
	Help       key.Binding
	Quit       key.Binding
	ShowHint   key.Binding
	Next       key.Binding
}

var keys = keyMap{
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
	viewport    viewport.Model
	contentStr  string
}

// NewChallengeView creates a new challenge view
func NewChallengeView(gs *game.GameState, challenge challenges.Challenge, width, height int, source MenuType) *ChallengeView {
	helpModel := help.New()
	helpModel.Width = 80

	challengeView := &ChallengeView{
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

	// Initialize the viewport
	viewportHeight := max(height-4, 5)
	challengeView.viewport = viewport.New(width, viewportHeight)

	// Generate initial content
	challengeView.updateContent()

	return challengeView
}

// generates the content for the viewport
func (m *ChallengeView) updateContent() {
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

		categoryHeader := fmt.Sprintf(" CATEGORY: %s", m.gameState.GetCurrentCategory())
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
	wrappedCode := utils.WrapText(highlightedCode, m.width)

	b.WriteString(codeBoxStyle.Render(wrappedCode) + "\n\n")

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

	// Set the content in the viewport
	m.contentStr = b.String()
	m.viewport.SetContent(m.contentStr)
}

// Init initializes the challenge view
func (m *ChallengeView) Init() tea.Cmd {
	return nil
}

// Update handles messages and user input
func (m *ChallengeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// if challenge is answered correctly, completion of category is 100%
		// and user goes to next challenge, category explanation is shown
		if m.hasAnswered &&
			m.isCorrect &&
			key.Matches(msg, keys.Next) &&
			m.gameState.ShouldShowVulnerabilityExplanation(m.gameState.GetCurrentCategory()) {
			explanationView := NewExplanationView(m.gameState, m.challenge, m.width, m.height, m.sourceMenu, true)
			return explanationView, explanationView.Init()
		} else if m.hasAnswered &&
			m.isCorrect &&
			key.Matches(msg, keys.Next) &&
			m.gameState.MoveToNextChallenge() {
			challenge := m.gameState.GetCurrentChallenge()
			challengeView := NewChallengeView(m.gameState, challenge, m.width, m.height, MainMenu)
			return challengeView, challengeView.Init()
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Back):
			return m, func() tea.Msg {
				return backToMenuMsg{}
			}

		case key.Matches(msg, keys.Help):
			m.helpModel.ShowAll = !m.helpModel.ShowAll

		case key.Matches(msg, keys.ShowHint):
			m.showHint = !m.showHint
			m.updateContent()

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 && (!m.hasAnswered || !m.isCorrect) {
				m.cursor--
				// Reset the result message for another try
				if m.hasAnswered && !m.isCorrect {
					m.hasAnswered = false
					m.result = ""
				}
				m.updateContent()

				// Find the cursor position in the content and scroll to it
				cursorPos := strings.Index(m.contentStr, "> ")
				if cursorPos > -1 {
					m.viewport.SetYOffset(0) // Reset to top first
					linesBefore := strings.Count(m.contentStr[:cursorPos], "\n")
					if linesBefore > m.viewport.Height/2 {
						m.viewport.SetYOffset(linesBefore - m.viewport.Height/2)
					}
				}
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.challenge.Options)-1 && (!m.hasAnswered || !m.isCorrect) {
				m.cursor++
				// Reset the result message for another try
				if m.hasAnswered && !m.isCorrect {
					m.hasAnswered = false
					m.result = ""
				}
				m.updateContent()

				// Find the cursor position in the content and scroll to it
				cursorPos := strings.Index(m.contentStr, "> ")
				if cursorPos > -1 {
					linesBefore := strings.Count(m.contentStr[:cursorPos], "\n")

					// If cursor is below visible area, scroll down
					if linesBefore >= m.viewport.YOffset+m.viewport.Height {
						m.viewport.SetYOffset(linesBefore - 3) // Show a few lines of context
					}
				}
			}

		case key.Matches(msg, keys.ScrollUp):
			m.viewport.LineUp(1)

		case key.Matches(msg, keys.ScrollDown):
			m.viewport.LineDown(1)

		case key.Matches(msg, keys.Select):
			m.hasAnswered = true
			selectedOption := m.challenge.Options[m.cursor]
			currentCategory := m.gameState.GetCurrentCategory()
			if selectedOption == m.challenge.CorrectAnswer {
				m.isCorrect = true
				m.result = "✓ Correct! You've identified the vulnerability."
				m.resultStyle = successStyle
				m.gameState.MarkChallengeCompleted(m.challenge.ID)
			} else {
				m.isCorrect = false
				m.result = "✗ Incorrect. Try another option by moving arrow keys!"
				m.resultStyle = errorStyle
				m.gameState.AddErrorCount(currentCategory)
			}
			m.updateContent()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update viewport size
		viewportHeight := max(m.height-4, 5)
		m.viewport.Width = msg.Width
		m.viewport.Height = viewportHeight

		// Regenerate content for new dimensions
		m.updateContent()
	}

	// Handle viewport updates
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the challenge
func (m *ChallengeView) View() string {
	var b strings.Builder

	// Render viewport content
	b.WriteString(m.viewport.View())

	// Add scroll indicator if needed
	hasScroll := m.viewport.YOffset > 0 ||
		m.viewport.YOffset+m.viewport.Height < strings.Count(m.contentStr, "\n")+1

	// Help section
	helpText := m.helpModel.View(keys)
	if !m.helpModel.ShowAll {
		helpText = helpHintStyle.Render("Press ? for help")
		if hasScroll {
			helpText += " | j/k to scroll"
		}
	}
	b.WriteString("\n" + helpStyle.Render(helpText))

	return b.String()
}

// --- helpers ---
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.ScrollUp, k.ScrollDown},
		{k.Back, k.Help, k.Quit},
		{k.ShowHint, k.Next},
	}
}
