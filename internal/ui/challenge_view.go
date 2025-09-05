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

type ChallengeView struct {
	gameState      *game.GameState
	challenge      challenges.Challenge
	cursor         int
	showHint       bool
	showHelp       bool
	help           help.Model
	quizOptions    []string
	result         string
	resultStyle    lipgloss.Style
	width          int
	height         int
	hasAnswered    bool
	isCorrect      bool
	sourceMenu     MenuType
	viewport       viewport.Model
	contentStr     string
	helpHeight     int
	keys           ChallengeKeyMap
	viewportHelper *ViewportHelper
}

func NewChallengeView(gs *game.GameState, challenge challenges.Challenge, width, height int, source MenuType) *ChallengeView {
	challengeView := &ChallengeView{
		gameState:      gs,
		challenge:      challenge,
		help:           help.New(),
		resultStyle:    successStyle,
		showHint:       false,
		cursor:         0,
		width:          width,
		height:         height,
		hasAnswered:    false,
		isCorrect:      false,
		quizOptions:    challenge.Options,
		sourceMenu:     source,
		keys:           NewChallengeKeyMap(),
		viewportHelper: NewViewportHelper(width, height),
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
		difficultyText = difficultyStyle["beginner"].Render("> [Beginner]")
	case challenges.Intermediate:
		difficultyText = difficultyStyle["intermediate"].Render("> [Intermediate]")
	case challenges.Advanced:
		difficultyText = difficultyStyle["advanced"].Render("> [Advanced]")
	}

	b.WriteString(difficultyText + "\n\n")

	showCategory := m.gameState.Settings.ShowVulnerabilityNames || (m.hasAnswered && m.isCorrect)

	if showCategory {
		b.WriteString(titleStyle.Render(m.challenge.Title) + "\n")

		categoryHeader := fmt.Sprintf(" CATEGORY: %s", m.gameState.GetCurrentCategory())
		b.WriteString(categoryStyle.Render(categoryHeader) + "\n\n")
	}

	language := ""
	if len(m.challenge.Lang) > 0 {
		lang := utils.GetLanguageFromChallenge(m.challenge.Lang)
		if lang != "" {
			language = lang
		}
	}

	separator := strings.Repeat("─", m.width)
	b.WriteString(subtleStyle.Render(separator) + "\n\n")

	highlightedCode := utils.HighlightCode(m.challenge.Code, language)

	b.WriteString(codeBoxStyle.Render(highlightedCode) + "\n")

	b.WriteString(subtleStyle.Render(separator) + "\n\n")

	b.WriteString(descStyle.Render("What vulnerability is in this code?") + "\n\n")

	for i, option := range m.challenge.Options {
		var renderedOption string
		cursor := "  "

		if m.hasAnswered && m.isCorrect {
			if option == m.challenge.CorrectAnswer {
				cursor = "✓ "
				renderedOption = cursor + correctAnswerStyle.Render(option)
			} else {
				renderedOption = cursor + unselectedItemStyle.Render(option)
			}
		} else if m.hasAnswered && !m.isCorrect && i == m.cursor {
			cursor = "✗ "
			renderedOption = cursor + incorrectAnswerStyle.Render(option)
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
			if m.challenge.Explanation != "" {
				b.WriteString("\n" + explanationStyle.Render("💡 Why this is correct:") + "\n")
				wrappedExplanation := utils.WrapText(m.challenge.Explanation, m.width)
				b.WriteString(explanationTextStyle.Render(wrappedExplanation) + "\n")
			}

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
		if m.hasAnswered && m.isCorrect && key.Matches(msg, m.keys.Next) {
			currentCategory := m.gameState.GetCurrentCategory()
			currentCategoryIdx := m.gameState.CurrentCategoryIdx
			currentChallengeIdx := m.gameState.CurrentChallengeIdx

			// Check if this was the last challenge in the current category
			currentSet := m.gameState.ChallengeSets[currentCategoryIdx]
			isLastChallengeInCategory := currentChallengeIdx == len(currentSet.Challenges)-1

			if isLastChallengeInCategory {
				// This was the last challenge in the category - check if we should show explanation
				shouldShow := m.gameState.ShouldShowVulnerabilityExplanation(currentCategory)

				if shouldShow {
					m.gameState.SetPendingCategoryExplanation(currentCategory)
					explanationView := NewExplanationView(m.gameState, &m.challenge, m.width, m.height, m.sourceMenu, true)
					return explanationView, explanationView.Init()
				} else {
					// If started via CLI, this means all intended challenges are completed
					if m.gameState.StartedViaCLI {
						return CLICompletionViewScreen(m.gameState, m.width, m.height), nil
					}

					// Move to next category/challenge (normal game mode)
					if m.gameState.MoveToNextChallenge() {
						challenge := m.gameState.GetCurrentChallenge()
						challengeView := NewChallengeView(m.gameState, challenge, m.width, m.height, MainMenu)
						return challengeView, challengeView.Init()
					} else {
						// No more challenges available, show completion view
						return CompletionViewScreen(m.gameState, m.width, m.height, MainMenu), nil
					}
				}
			} else {
				// More challenges in current category, move to next challenge
				if m.gameState.MoveToNextChallenge() {
					challenge := m.gameState.GetCurrentChallenge()
					challengeView := NewChallengeView(m.gameState, challenge, m.width, m.height, MainMenu)
					return challengeView, challengeView.Init()
				}
			}
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Back):
			if m.gameState.StartedViaCLI {
				return m, nil
			}
			return m, func() tea.Msg {
				return backToMenuMsg{}
			}

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			m.updateViewportDimensions()
			m.updateContent()

		case key.Matches(msg, m.keys.ShowHint):
			m.showHint = !m.showHint
			m.updateContent()

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 && (!m.hasAnswered || !m.isCorrect) {
				m.cursor--
				if m.hasAnswered && !m.isCorrect {
					m.hasAnswered = false
					m.result = ""
				}
				m.updateContent()
				m.ensureCursorVisible()
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.challenge.Options)-1 && (!m.hasAnswered || !m.isCorrect) {
				m.cursor++
				if m.hasAnswered && !m.isCorrect {
					m.hasAnswered = false
					m.result = ""
				}
				m.updateContent()
				m.ensureCursorVisible()
			}

		case key.Matches(msg, m.keys.ScrollUp):
			if !m.hasAnswered || m.showHelp {
				m.viewport.LineUp(1)
			}

		case key.Matches(msg, m.keys.ScrollDown):
			if !m.hasAnswered || m.showHelp {
				m.viewport.LineDown(1)
			}

		case key.Matches(msg, m.keys.Select):
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
			// Focus on the result and explanation after answering
			m.focusOnResult()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.updateViewportDimensions()

		m.updateContent()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *ChallengeView) View() string {
	var b strings.Builder
	b.WriteString(m.viewport.View())

	hasScroll := m.viewport.YOffset > 0 ||
		m.viewport.YOffset+m.viewport.Height < strings.Count(m.contentStr, "\n")+1

	b.WriteString("\n")

	notificationDisplay := NewNotificationDisplay()
	b.WriteString(notificationDisplay.RenderAllNotifications(m.gameState))

	if m.showHelp {
		b.WriteString(m.help.View(m))
	} else {
		helpText := "Press ? for help | ↑/↓ to select option"
		if hasScroll {
			helpText += " | j/k to scroll content"
		}
		if !m.gameState.StartedViaCLI {
			helpText += " | ESC to go back"
		}
		b.WriteString(helpHintStyle.Render(helpText))
	}

	return b.String()
}

// Helper methods
func (k ChallengeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k ChallengeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.ScrollUp, k.ScrollDown},
		{k.Select},
		{k.Help, k.Quit},
		{k.ShowHint, k.Next},
	}
}

// ShortHelp returns the short help key bindings for ChallengeView
func (m *ChallengeView) ShortHelp() []key.Binding {
	return m.keys.ShortHelp()
}

// FullHelp returns the full help key bindings for ChallengeView
func (m *ChallengeView) FullHelp() [][]key.Binding {
	helpKeys := m.keys.FullHelp()

	if !m.gameState.StartedViaCLI {
		helpKeys[2] = append(helpKeys[2], m.keys.Back)
	}

	return helpKeys
}

// ensureCursorVisible ensures the selected option is visible in the viewport
func (m *ChallengeView) ensureCursorVisible() {
	// Find the position of the options section in the content
	optionsStartMarker := "What vulnerability is in this code?"
	optionsPos := strings.Index(m.contentStr, optionsStartMarker)

	if optionsPos == -1 {
		return
	}

	// Calculate the line number where options start
	linesBeforeOptions := strings.Count(m.contentStr[:optionsPos], "\n")

	// Calculate the line number of the selected option
	selectedOptionLine := linesBeforeOptions + 2 + m.cursor

	// Get current viewport dimensions
	viewportHeight := m.viewport.Height
	currentOffset := m.viewport.YOffset

	// Position the selected option with padding for context
	padding := 2
	idealOffset := max(selectedOptionLine-padding, 0)

	// Ensure we don't scroll past the end of content
	maxOffset := max(strings.Count(m.contentStr, "\n")-viewportHeight+1, 0)
	if idealOffset > maxOffset {
		idealOffset = maxOffset
	}

	// Only scroll if the selected option is not properly visible
	if selectedOptionLine < currentOffset || selectedOptionLine >= currentOffset+viewportHeight {
		m.viewport.SetYOffset(idealOffset)
	}
}

// focusOnResult focuses the viewport on the result and explanation after answering
func (m *ChallengeView) focusOnResult() {
	// Find the position of the result section
	resultMarker := "✓ Correct! You've identified the vulnerability."
	if !m.isCorrect {
		resultMarker = "✗ Incorrect. Try another option by moving arrow keys!"
	}

	resultPos := strings.Index(m.contentStr, resultMarker)
	if resultPos == -1 {
		return
	}

	// Calculate the line number where the result starts
	resultLine := strings.Count(m.contentStr[:resultPos], "\n")

	// Get current viewport dimensions
	viewportHeight := m.viewport.Height

	// Position the result near the top of the viewport with some padding
	padding := 1
	idealOffset := max(resultLine-padding, 0)

	// Ensure we don't scroll past the end of content
	maxOffset := max(strings.Count(m.contentStr, "\n")-viewportHeight+1, 0)
	if idealOffset > maxOffset {
		idealOffset = maxOffset
	}

	// Scroll to show the result
	m.viewport.SetYOffset(idealOffset)
}
