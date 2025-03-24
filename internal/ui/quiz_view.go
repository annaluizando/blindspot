package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// QuizKeyMap defines keybindings for quizzes
type QuizKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Hint   key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k QuizKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k QuizKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select}, // Navigation
		{k.Back, k.Hint, k.Help}, // Actions
		{k.Quit},                 // System
	}
}

// QuizKeys holds the quiz key mappings
var QuizKeys = QuizKeyMap{
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
		key.WithHelp("enter", "select option"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Hint: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "show hint"),
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

// QuizOption represents a single quiz option
type QuizOption struct {
	Text string
	ID   int
}

// QuizAnswerMsg is sent when an answer is selected
type QuizAnswerMsg struct {
	SelectedID int
	Correct    bool
}

// QuizView is the multiple-choice quiz component
type QuizView struct {
	question     string
	options      []QuizOption
	correctID    int
	cursor       int
	selected     bool
	showHint     bool
	hint         string
	help         help.Model
	showHelp     bool
	result       string
	resultStyle  lipgloss.Style
	width        int
	height       int
	explanation  string
	hasAnswered  bool
}

// NewQuizView creates a new quiz view
func NewQuizView(question string, options []string, correctID int, hint string, width, height int) *QuizView {
	quizOptions := make([]QuizOption, len(options))
	for i, opt := range options {
		quizOptions[i] = QuizOption{
			Text: opt,
			ID:   i,
		}
	}

	return &QuizView{
		question:    question,
		options:     quizOptions,
		correctID:   correctID,
		help:        help.New(),
		hint:        hint,
		width:       width,
		height:      height,
		resultStyle: lipgloss.NewStyle().Padding(1).BorderStyle(lipgloss.RoundedBorder()),
		hasAnswered: false,
	}
}

// Init initializes the quiz view
func (q *QuizView) Init() tea.Cmd {
	return nil
}

// Update handles messages and user input
func (q *QuizView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, QuizKeys.Quit):
			return q, tea.Quit

		case key.Matches(msg, QuizKeys.Back):
			return q, func() tea.Msg {
				return backToMenuMsg{}
			}

		case key.Matches(msg, QuizKeys.Up):
			if !q.hasAnswered && q.cursor > 0 {
				q.cursor--
			}

		case key.Matches(msg, QuizKeys.Down):
			if !q.hasAnswered && q.cursor < len(q.options)-1 {
				q.cursor++
			}

		case key.Matches(msg, QuizKeys.Help):
			q.showHelp = !q.showHelp

		case key.Matches(msg, QuizKeys.Hint):
			q.showHint = !q.showHint

		case key.Matches(msg, QuizKeys.Select):
			if !q.hasAnswered {
				q.hasAnswered = true
				q.selected = true
				isCorrect := q.cursor == q.correctID

				if isCorrect {
					q.result = "✓ Correct! You've identified the security vulnerability."
					q.resultStyle = successStyle
				} else {
					q.result = "✗ Incorrect. Try again by pressing ESC to go back."
					q.resultStyle = errorStyle
				}

				return q, func() tea.Msg {
					return QuizAnswerMsg{
						SelectedID: q.cursor,
						Correct:    isCorrect,
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		q.width = msg.Width
		q.height = msg.Height
	}

	return q, nil
}

// View renders the quiz view
func (q *QuizView) View() string {
	var b strings.Builder

	// Question
	b.WriteString(titleStyle.Render("Security Quiz") + "\n\n")
	b.WriteString(descStyle.Render(q.question) + "\n\n")

	// Options
	for i, option := range q.options {
		var renderedOption string
		cursor := "  "

		// Show selection cursor or answer indicators
		if q.hasAnswered {
			if i == q.correctID {
				cursor = "✓ "
				correctStyle := successStyle
				correctStyle = correctStyle.Bold(false)
				renderedOption = cursor + correctStyle.Render(option.Text)
			} else if i == q.cursor && i != q.correctID {
				cursor = "✗ "
				incorrectStyle := errorStyle
				incorrectStyle = incorrectStyle.Bold(false)
				renderedOption = cursor + incorrectStyle.Render(option.Text)
			} else {
				renderedOption = cursor + unselectedItemStyle.Render(option.Text)
			}
		} else {
			if q.cursor == i {
				cursor = "▶ "
				renderedOption = cursor + selectedItemStyle.Render(option.Text)
			} else {
				renderedOption = cursor + unselectedItemStyle.Render(option.Text)
			}
		}

		b.WriteString(renderedOption + "\n")
	}

	// Hint
	if q.showHint {
		b.WriteString("\n" + hintStyle.Render("Hint: " + q.hint) + "\n")
	}

	// Result
	if q.result != "" {
		b.WriteString("\n" + q.resultStyle.Render(q.result) + "\n")
	}

	// Explanation (shown after answering)
	if q.hasAnswered && q.explanation != "" {
		b.WriteString("\n" + descStyle.Render("Explanation: " + q.explanation) + "\n")
	}

	// Help
	helpText := ""
	if q.showHelp {
		helpText = q.help.View(QuizKeys)
	} else {
		helpText = helpHintStyle.Render("Press ? for help, h for hint")
	}
	b.WriteString("\n" + helpText)

	return b.String()
}

// IsCorrect returns whether the selected answer is correct
func (q *QuizView) IsCorrect() bool {
	return q.hasAnswered && q.cursor == q.correctID
}

// SetExplanation sets an explanation text to show after answering
func (q *QuizView) SetExplanation(explanation string) {
	q.explanation = explanation
}
