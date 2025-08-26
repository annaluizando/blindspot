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

// Navigation messages
type nextChallengeMsg struct{}
type backToMenuMsg struct{}

// ExplanationKeyMap defines keybindings for the explanation view
type ExplanationKeyMap struct {
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Next       key.Binding
	Back       key.Binding
	Help       key.Binding
	Quit       key.Binding
}

// ExplanationView displays vulnerability explanations with scrolling support
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

// ExplanationKeys defines the key bindings for the explanation view
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

// Styles for the explanation view
var (
	explanationSubtitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#0066CC")).Bold(true)      // Primary blue
	explanationHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00AA55")).Bold(true)      // Security green
	explanationTextStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))                 // Light text
	resourceStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#0066CC")).Underline(true) // Primary blue
	completedStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("#00CC44")).Bold(true)      // Success green
)

// NewExplanationView creates a new explanation view
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

	explanationView.updateViewportDimensions()
	explanationView.updateContent()

	return explanationView
}

// Init initializes the explanation view
func (v *ExplanationView) Init() tea.Cmd {
	return nil
}

// Update handles messages and user input
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
			return v.handleNextAction()

		case key.Matches(msg, keys.ScrollUp):
			v.viewport.LineUp(1)

		case key.Matches(msg, keys.ScrollDown):
			v.viewport.LineDown(1)

		case key.Matches(msg, ExplanationKeys.Back):
			return v.handleBackAction()
		}

	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		v.updateViewportDimensions()
	}

	v.viewport, cmd = v.viewport.Update(msg)
	return v, cmd
}

// handleNextAction handles the next action based on context
func (v *ExplanationView) handleNextAction() (tea.Model, tea.Cmd) {
	if v.isFromCompletion {
		return v, func() tea.Msg {
			return nextChallengeMsg{}
		}
	}

	// Return to category view
	return v.navigateToCategoryView()
}

// handleBackAction handles the back action based on context
func (v *ExplanationView) handleBackAction() (tea.Model, tea.Cmd) {
	if v.sourceMenu == ChallengeMenu {
		return v.navigateToCategoryView()
	}

	// Return to main menu
	return v, func() tea.Msg {
		return backToMenuMsg{}
	}
}

// navigateToCategoryView navigates back to the category view
func (v *ExplanationView) navigateToCategoryView() (tea.Model, tea.Cmd) {
	for i, set := range v.gameState.ChallengeSets {
		if set.Category == v.gameState.GetCurrentCategory() {
			return NewCategoryMenu(v.gameState, i, v.width, v.height, v.sourceMenu), nil
		}
	}
	return v, nil
}

// updateViewportDimensions updates the viewport dimensions
func (v *ExplanationView) updateViewportDimensions() {
	viewportHeight := max(v.height-4, 5)
	v.viewport = viewport.New(v.width, viewportHeight)
}

// updateContent updates the content and viewport
func (v *ExplanationView) updateContent() {
	var b strings.Builder

	// Header section
	v.buildHeader(&b)

	if v.explanationFound {
		v.buildExplanationContent(&b)
	} else {
		v.buildNoExplanationContent(&b)
	}

	v.buildHelpSection(&b)

	v.contentStr = b.String()
	v.updateViewportContent()
}

func (v *ExplanationView) buildHeader(b *strings.Builder) {
	if v.isFromCompletion {
		b.WriteString(completedStyle.Render("🎉 Challenge Completed!") + "\n\n")
		b.WriteString(fmt.Sprintf("You've completed: %s\n\n", selectedItemStyle.Render(v.challenge.Title)))
	} else {
		b.WriteString(explanationHighlightStyle.Render("🔍 Category Explanation") + "\n\n")
	}

	b.WriteString(fmt.Sprintf("%s\n\n", explanationHighlightStyle.Render(v.gameState.GetCurrentCategory())))
}

func (v *ExplanationView) buildExplanationContent(b *strings.Builder) {
	b.WriteString(explanationSubtitleStyle.Render("What is this vulnerability?") + "\n")
	wrappedDesc := utils.WrapText(v.explanation.ShortDescription, v.width)
	b.WriteString(descriptionStyle.Render(wrappedDesc) + "\n\n")

	b.WriteString(explanationSubtitleStyle.Render("Learn More:") + "\n")
	wrappedExplanation := utils.WrapText(v.explanation.Explanation, v.width)
	b.WriteString(explanationTextStyle.Render(wrappedExplanation) + "\n\n")

	if len(v.explanation.Resources) > 0 {
		v.buildResourcesSection(b)
	}
}

func (v *ExplanationView) buildNoExplanationContent(b *strings.Builder) {
	b.WriteString(errorStyle.Render("Detailed explanation for this vulnerability category is not available yet.") + "\n\n")
}

func (v *ExplanationView) buildResourcesSection(b *strings.Builder) {
	b.WriteString(explanationSubtitleStyle.Render("Additional Resources:") + "\n")
	for _, resource := range v.explanation.Resources {
		b.WriteString(fmt.Sprintf("- %s: %s\n",
			resource.Title,
			resourceStyle.Render(resource.URL)))
	}
	b.WriteString("\n")
}

func (v *ExplanationView) buildHelpSection(b *strings.Builder) {
	if v.showHelp {
		b.WriteString("\n" + v.help.View(ExplanationKeys))
	} else if v.isFromCompletion {
		b.WriteString("\n" + helpHintStyle.Render("Press 'Enter'/'N' to continue to next challenge"))
		b.WriteString("\n" + helpHintStyle.Render("Press ? for help"))
	} else {
		b.WriteString("\n" + helpHintStyle.Render("Press ? for help"))
	}
}

func (v *ExplanationView) updateViewportContent() {
	helpHeight := 1
	if v.showHelp {
		helpHeight = 4
	}

	contentHeight := strings.Count(v.contentStr, "\n") + 1
	viewportHeight := min(contentHeight, v.height-helpHeight-1)

	v.viewport = viewport.New(v.width, viewportHeight)
	v.viewport.SetContent(v.contentStr)
}

func (v *ExplanationView) View() string {
	var b strings.Builder

	b.WriteString(v.viewport.View())

	v.buildScrollIndicator(&b)

	v.buildHelpFooter(&b)

	return b.String()
}

func (v *ExplanationView) buildScrollIndicator(b *strings.Builder) {
	hasScroll := v.viewport.YOffset > 0 || v.viewport.YOffset+v.viewport.Height < strings.Count(v.contentStr, "\n")+1
	if hasScroll {
		b.WriteString("\n" + dimStyle.Render("j/k to scroll"))
	}
}

func (v *ExplanationView) buildHelpFooter(b *strings.Builder) {
	if v.showHelp {
		b.WriteString("\n" + v.help.View(MenuKeys))
	} else {
		v.buildHelpText(b)
	}
}

func (v *ExplanationView) buildHelpText(b *strings.Builder) {
	hasScroll := v.viewport.YOffset > 0 || v.viewport.YOffset+v.viewport.Height < strings.Count(v.contentStr, "\n")+1

	helpText := "Press ? for help | ↑/↓ to navigate"
	if hasScroll {
		helpText += " | j/k to scroll"
	}

	if v.width < 60 {
		helpText = "? for help | ↑/↓ nav"
		if hasScroll {
			helpText += " | j/k scroll"
		}
	}

	b.WriteString("\n" + helpHintStyle.Render(helpText))
}

// ShortHelp returns the short help key bindings
func (k ExplanationKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns the full help key bindings
func (k ExplanationKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Back},
		{k.ScrollUp, k.ScrollDown},
		{k.Help, k.Quit},
	}
}
