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

type ChallengeView struct {
	gameState   *game.GameState
	challenge   challenges.Challenge
	cursor      int
	showHint    bool
	showHelp    bool
	help        help.Model
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
	helpHeight  int
}

func NewChallengeView(gs *game.GameState, challenge challenges.Challenge, width, height int, source MenuType) *ChallengeView {
	challengeView := &ChallengeView{
		gameState:   gs,
		challenge:   challenge,
		help:        help.New(),
		resultStyle: successStyle,
		showHint:    false,
		cursor:      0,
		width:       width,
		height:      height,
		hasAnswered: false,
		isCorrect:   false,
		quizOptions: challenge.Options,
		sourceMenu:  source,
	}

	challengeView.updateContent()
	return challengeView
}

func (m *ChallengeView) updateViewportDimensions() {
	m.helpHeight = 1
	if m.showHelp {
		// If showing help on small screen, give it a bit more space
		m.helpHeight = 2
	}

	// Minimum height required to show at least 3 options plus question
	minimumOptionsHeight := 8 // 1 for question + 3 options + 1 buffer

	// Calculate viewport height ensuring there's enough room for options
	viewportHeight := max(m.height-m.helpHeight, minimumOptionsHeight)

	m.viewport.Height = viewportHeight - 2
	m.viewport.Width = m.width
}

func (m *ChallengeView) updateContent() {
	var b strings.Builder

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

	// for visual clarity
	separator := strings.Repeat("─", m.width/2)
	b.WriteString(subtleStyle.Render(separator) + "\n\n")

	language := ""
	if len(m.challenge.Lang) > 0 {
		lang := utils.GetLanguageFromExtension(m.challenge.Lang)
		if lang != "" {
			language = lang
		}
	}

	highlightedCode := utils.HighlightCode(m.challenge.Code, language)

	b.WriteString(codeBoxStyle.Render(highlightedCode) + "\n\n")

	b.WriteString(descStyle.Render("What vulnerability is in this code?") + "\n\n")

	for i, option := range m.challenge.Options {
		var renderedOption string
		cursor := "  "

		if m.hasAnswered && m.isCorrect {
			if option == m.challenge.CorrectAnswer {
				cursor = "✓ "
				successStyleCopy := successStyle
				successStyleCopy = successStyleCopy.Bold(false)
				renderedOption = cursor + successStyleCopy.Render(option)
			} else {
				renderedOption = cursor + unselectedItemStyle.Render(option)
			}
		} else if m.hasAnswered && !m.isCorrect && i == m.cursor {
			cursor = "✗ "
			errorStyleCopy := errorStyle
			errorStyleCopy = errorStyleCopy.Bold(false)
			renderedOption = cursor + errorStyleCopy.Render(option)
		} else if !m.hasAnswered {
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

	if m.showHint {
		b.WriteString("\n" + hintStyle.Render("Hint: "+m.challenge.Hint) + "\n")
	}

	if m.result != "" {
		b.WriteString("\n" + m.resultStyle.Render(m.result) + "\n")

		if m.isCorrect {
			b.WriteString("\n" + helpHintStyle.Render("Press 'Enter'/'N' to continue to next challenge"))
		}
	}

	m.contentStr = b.String()
	m.viewport.SetContent(m.contentStr)

	// Calculate options area position to ensure it's visible
	optionsStartMarker := "What vulnerability is in this code?"
	optionsPos := strings.Index(m.contentStr, optionsStartMarker)

	if optionsPos > -1 {
		linesBeforeOptions := strings.Count(m.contentStr[:optionsPos], "\n")

		totalLines := strings.Count(m.contentStr, "\n") + 1

		// Calculate lines for options area
		optionsLines := totalLines - linesBeforeOptions

		// If window is small and not all content fits, prioritize showing options
		if optionsLines+3 > m.viewport.Height {
			newOffset := max(0, linesBeforeOptions-2)
			m.viewport.SetYOffset(newOffset)
		}
	}

	// Ensure cursor is visible - find cursor position in content
	if !m.hasAnswered || !m.isCorrect {
		cursorPos := strings.Index(m.contentStr, "> ")
		if cursorPos > -1 {
			linesBefore := strings.Count(m.contentStr[:cursorPos], "\n")

			// Adjust viewport if cursor would be outside visible area
			if linesBefore < m.viewport.YOffset {
				m.viewport.SetYOffset(linesBefore)
			} else if linesBefore >= m.viewport.YOffset+m.viewport.Height-2 {
				m.viewport.SetYOffset(linesBefore - m.viewport.Height + 3)
			}
		}
	}
}

func (m *ChallengeView) Init() tea.Cmd {
	m.updateContent()
	return tea.Sequence(
		func() tea.Msg {
			return tea.WindowSizeMsg{
				Width:  m.width,
				Height: m.height,
			}
		},
	)
}

func (m *ChallengeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle challenge completion navigation
		if m.hasAnswered && m.isCorrect && key.Matches(msg, keys.Next) {
			if m.gameState.ShouldShowVulnerabilityExplanation(m.gameState.GetCurrentCategory()) {
				explanationView := NewExplanationView(m.gameState, m.challenge, m.width, m.height, m.sourceMenu, true)
				return explanationView, explanationView.Init()
			} else if m.gameState.MoveToNextChallenge() {
				challenge := m.gameState.GetCurrentChallenge()
				challengeView := NewChallengeView(m.gameState, challenge, m.width, m.height, MainMenu)
				return challengeView, challengeView.Init()
			}
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Back):
			return m, func() tea.Msg {
				return backToMenuMsg{}
			}

		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp
			m.updateViewportDimensions()
			m.updateContent()

		case key.Matches(msg, keys.ShowHint):
			m.showHint = !m.showHint
			m.updateContent()

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 && (!m.hasAnswered || !m.isCorrect) {
				m.cursor--
				if m.hasAnswered && !m.isCorrect {
					m.hasAnswered = false
					m.result = ""
				}
				m.updateContent()
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.challenge.Options)-1 && (!m.hasAnswered || !m.isCorrect) {
				m.cursor++
				if m.hasAnswered && !m.isCorrect {
					m.hasAnswered = false
					m.result = ""
				}
				m.updateContent()
			}

		case key.Matches(msg, keys.ScrollUp):
			if m.viewport.YOffset > 0 {
				m.viewport.LineUp(1)
			}

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

		m.updateViewportDimensions()

		m.updateContent()
	}

	// Handle viewport updates
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *ChallengeView) View() string {
	var b strings.Builder
	b.WriteString(m.viewport.View())

	hasScroll := m.viewport.YOffset > 0 ||
		m.viewport.YOffset+m.viewport.Height < strings.Count(m.contentStr, "\n")+1

	b.WriteString("\n")

	// The help text
	if m.showHelp {
		b.WriteString(m.help.View(MenuKeys))
	} else {
		helpText := "Press ? for help | ↑/↓ to navigate"
		if hasScroll {
			helpText += " | j/k to scroll"
		}
		b.WriteString(helpHintStyle.Render(helpText))
	}

	return b.String()
}

// ----helpers----
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.ScrollUp, k.ScrollDown},
		{k.Select, k.Back},
		{k.Help, k.Quit},
		{k.ShowHint, k.Next},
	}
}
