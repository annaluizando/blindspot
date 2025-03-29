package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"blindspot/internal/game"
	"blindspot/internal/utils"
)

type CompletionKeyMap struct {
	Next key.Binding
	Back key.Binding
	Help key.Binding
	Quit key.Binding
}

type CompletionView struct {
	gameState  *game.GameState
	width      int
	height     int
	sourceMenu MenuType
	help       help.Model
	showHelp   bool
}

var CompletionKeys = CompletionKeyMap{
	Next: key.NewBinding(
		key.WithKeys("enter", "n"),
		key.WithHelp("enter/n", "back to main menu"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back to main menu"),
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
	textStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
	highlightStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFA33"))
	contributionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF87FF")).Italic(true)
)

func NewCompletionView(gs *game.GameState, width, height int, sourceMenu MenuType) *CompletionView {
	return &CompletionView{
		gameState:  gs,
		width:      width,
		height:     height,
		sourceMenu: sourceMenu,
		help:       help.New(),
		showHelp:   false,
	}
}

func (v *CompletionView) Init() tea.Cmd {
	return nil
}

// handles messages and user input
func (v *CompletionView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, CompletionKeys.Quit):
			return v, tea.Quit

		case key.Matches(msg, CompletionKeys.Help):
			v.showHelp = !v.showHelp
			return v, nil

		case key.Matches(msg, CompletionKeys.Next), key.Matches(msg, CompletionKeys.Back):
			// return to main menu
			return v, func() tea.Msg {
				return backToMenuMsg{}
			}
		}

	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
	}

	return v, nil
}

// renders completion view
func (v *CompletionView) View() string {
	var b strings.Builder

	header := titleStyle.Render("🎉 Congratulations! You've Completed All Challenges! 🎉")
	b.WriteString(header + "\n\n")

	message := utils.WrapText("You've demonstrated a strong understanding of secure coding practices and completed all challenges in blindspot! The skills you've developed here will help you write safer, more robust code in your projects.", v.width)
	b.WriteString(textStyle.Render(message) + "\n\n")

	b.WriteString(subtitleStyle.Render("🔦 What's Next?") + "\n")
	b.WriteString(textStyle.Render("Your journey doesn't end here! Security is an ongoing learning process.") + "\n\n")

	b.WriteString(contributionStyle.Render("📚 Keep Learning") + "\n")
	keepLearning := utils.WrapText("Stay updated with the latest security practices and vulnerabilities:\n", v.width)
	b.WriteString(textStyle.Render(keepLearning) + "\n\n")
	b.WriteString("- Follow the OWASP Top 10\n")
	b.WriteString("- Join security communities\n")
	b.WriteString("- Practice on other platforms\n\n")

	b.WriteString(contributionStyle.Render("💡 Contribute Your Own Challenges") + "\n")
	contribution := utils.WrapText("Want to help others learn while reinforcing your own knowledge? Consider creating your own challenge scenarios! Check the #FAQ section in the README.md for instructions on creating challenges.yaml files.", v.width)
	b.WriteString(textStyle.Render(contribution) + "\n\n")

	if v.gameState != nil {
		totalChallenges := 0
		for _, set := range v.gameState.ChallengeSets {
			totalChallenges += len(set.Challenges)
		}

		statsMessage := fmt.Sprintf("You've completed %d challenges across %d vulnerability categories!",
			totalChallenges, len(v.gameState.ChallengeSets))
		b.WriteString(highlightStyle.Render(statsMessage) + "\n\n")
	}

	thankYou := utils.WrapText("Thank you for playing blindspot and joining our mission to create more secure software. Keep learning, keep coding securely!", v.width)
	b.WriteString(textStyle.Render(thankYou) + "\n\n")

	if v.showHelp {
		b.WriteString("\n" + v.help.View(CompletionKeys))
	} else {
		b.WriteString("\n" + helpHintStyle.Render("Press 'Enter'/'N' or 'Esc' to return to main menu"))
		b.WriteString("\n" + helpHintStyle.Render("Press ? for more options"))
	}

	return b.String()
}

// ---- helpers ----
func (k CompletionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k CompletionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Back},
		{k.Help, k.Quit},
	}
}
