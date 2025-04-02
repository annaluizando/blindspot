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

// messages used to navigate after showing vulnerability explanation
type nextChallengeMsg struct{}
type backToMenuMsg struct{}

// defines keybindings for the explanation view
type ExplanationKeyMap struct {
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Next       key.Binding
	Back       key.Binding
	Help       key.Binding
	Quit       key.Binding
}

type ExplanationView struct {
	gameState        *game.GameState
	challenge        challenges.Challenge
	explanation      challenges.VulnerabilityInfo
	explanationFound bool
	width            int
	height           int
	sourceMenu       MenuType
	help             help.Model
	showHelp         bool
	isFromCompletion bool
	viewport         viewport.Model
	contentStr       string
}

var ExplanationKeys = ExplanationKeyMap{
	ScrollUp: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "scroll up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "scroll down"),
	),
	Next: key.NewBinding(
		key.WithKeys("enter", "n"),
		key.WithHelp("enter/n", "next challenge"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back to menu"),
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

var (
	explanationSubtitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#87D7FF")).Bold(true)
	explanationHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFA33")).Bold(true)
	explanationTextStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
	resourceStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#87D7FF")).Underline(true)
	completedStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("#5FFF87")).Bold(true)
)

// vulnerability explanation view
func NewExplanationView(gs *game.GameState, challenge challenges.Challenge, width, height int, sourceMenu MenuType, isFromCompletion bool) *ExplanationView {
	explanation, found := gs.GetVulnerabilityExplanation(gs.GetCurrentCategory())

	explanationView := &ExplanationView{
		gameState:        gs,
		challenge:        challenge,
		explanation:      explanation,
		explanationFound: found,
		width:            width,
		height:           height,
		sourceMenu:       sourceMenu,
		help:             help.New(),
		showHelp:         false,
		isFromCompletion: isFromCompletion,
	}

	viewportHeight := max(height-4, 5)
	explanationView.viewport = viewport.New(width, viewportHeight)
	explanationView.updateContent()

	return explanationView
}

func (v *ExplanationView) Init() tea.Cmd {
	v.updateContent()
	return nil
}

// handles messages and user input
func (v *ExplanationView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ExplanationKeys.Quit):
			return v, tea.Quit

		case key.Matches(msg, ExplanationKeys.Help):
			v.showHelp = !v.showHelp
			v.updateContent()
			return v, nil

		case key.Matches(msg, ExplanationKeys.Next):
			if v.isFromCompletion {
				return v, func() tea.Msg {
					return nextChallengeMsg{}
				}
			} else {
				// If viewing from category menu, go back to category view
				for i, set := range v.gameState.ChallengeSets {
					if set.Category == v.gameState.GetCurrentCategory() {
						return NewCategoryMenu(v.gameState, i, v.width, v.height, v.sourceMenu), nil
					}
				}
			}

			if v.sourceMenu == ChallengeMenu {
				// Find the category index
				for i, set := range v.gameState.ChallengeSets {
					if set.Category == v.gameState.GetCurrentCategory() {
						return NewCategoryMenu(v.gameState, i, v.width, v.height, v.sourceMenu), nil
					}
				}
			}

		case key.Matches(msg, keys.ScrollUp):
			v.viewport.LineUp(1)

		case key.Matches(msg, keys.ScrollDown):
			v.viewport.LineDown(1)

		case key.Matches(msg, ExplanationKeys.Back):
			if v.sourceMenu == ChallengeMenu {
				// Find the category index
				for i, set := range v.gameState.ChallengeSets {
					if set.Category == v.gameState.GetCurrentCategory() {
						return NewCategoryMenu(v.gameState, i, v.width, v.height, v.sourceMenu), nil
					}
				}
			}

			// Otherwise, return to main menu
			return v, func() tea.Msg {
				return backToMenuMsg{}
			}
		}

	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height

		viewportHeight := max(v.height-4, 5)
		v.viewport.Width = msg.Width
		v.viewport.Height = viewportHeight
	}

	v.viewport, cmd = v.viewport.Update(msg)
	return v, cmd
}

func (v *ExplanationView) updateContent() {
	var b strings.Builder

	// Add scroll indicator if needed
	if v.isFromCompletion {
		b.WriteString(completedStyle.Render("🎉 Challenge Completed!") + "\n\n")
		b.WriteString(fmt.Sprintf("You've completed: %s\n\n", selectedItemStyle.Render(v.challenge.Title)))
	} else {
		b.WriteString(explanationHighlightStyle.Render("🔍 Category Explanation") + "\n\n")
	}

	b.WriteString(fmt.Sprintf("%s\n\n", explanationHighlightStyle.Render(v.gameState.GetCurrentCategory())))

	if v.explanationFound {
		b.WriteString(explanationSubtitleStyle.Render("What is this vulnerability?") + "\n")
		wrappedDesc := utils.WrapText(v.explanation.ShortDescription, v.width)
		b.WriteString(descriptionStyle.Render(wrappedDesc) + "\n\n")

		b.WriteString(explanationSubtitleStyle.Render("Learn More:") + "\n")
		wrappedExplanation := utils.WrapText(v.explanation.Explanation, v.width)
		b.WriteString(explanationTextStyle.Render(wrappedExplanation) + "\n\n")

		if len(v.explanation.Resources) > 0 {
			b.WriteString(explanationSubtitleStyle.Render("Additional Resources:") + "\n")
			for _, resource := range v.explanation.Resources {
				b.WriteString(fmt.Sprintf("- %s: %s\n",
					resource.Title,
					resourceStyle.Render(resource.URL)))
			}
			b.WriteString("\n")
		}
	} else {
		b.WriteString(errorStyle.Render("Detailed explanation for this vulnerability category is not available yet.") + "\n\n")
	}

	if v.showHelp {
		b.WriteString("\n" + v.help.View(ExplanationKeys))
	} else if v.isFromCompletion {
		b.WriteString("\n" + helpHintStyle.Render("Press 'Enter'/'N' to continue to next challenge"))
		b.WriteString("\n" + helpHintStyle.Render("Press ? for help"))
	} else {
		b.WriteString("\n" + helpHintStyle.Render("Press ? for help"))
	}

	helpHeight := 1
	if v.showHelp {
		helpHeight = 4 // Full help takes vore space
	}

	// Create or update the viewport
	v.contentStr = b.String()
	contentHeight := strings.Count(v.contentStr, "\n") + 1

	// If content is shorter than available height, no scrolling needed
	viewportHeight := min(contentHeight, v.height-helpHeight-1)

	v.viewport = viewport.New(v.width, viewportHeight)
	v.viewport.SetContent(v.contentStr)
}

// renders vulnerability explanation
func (v *ExplanationView) View() string {
	var b strings.Builder
	// Render viewport content
	b.WriteString(v.viewport.View())

	hasScroll := v.viewport.YOffset > 0 || v.viewport.YOffset+v.viewport.Height < strings.Count(v.contentStr, "\n")+1
	if hasScroll {
		scrollInfo := fmt.Sprintf(" Line %d of %d ",
			v.viewport.YOffset+1,
			strings.Count(v.contentStr, "\n")+1)
		b.WriteString("\n" + dimStyle.Render("j/k to scroll, currently at") + scrollInfo)
	}

	// Help
	if v.showHelp {
		b.WriteString("\n" + v.help.View(MenuKeys))
	} else {
		helpText := "Press ? for help | ↑/↓ to navigate | j/k to scroll"
		if v.width < 60 {
			helpText = "? for help | ↑/↓ nav | j/k scroll"
		}
		if !hasScroll {
			helpText = strings.Replace(helpText, " | j/k to scroll", "", 1)
		}
		b.WriteString("\n" + helpHintStyle.Render(helpText))
	}

	return b.String()
}

// ---- helpers ----
func (k ExplanationKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k ExplanationKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Back},
		{k.ScrollUp, k.ScrollDown},
		{k.Help, k.Quit},
	}
}
