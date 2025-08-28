package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"blindspot/internal/game"
	"blindspot/internal/utils"
)

type CLICompletionKeyMap struct {
	Quit key.Binding
	Help key.Binding
}

var CLICompletionKeys = CLICompletionKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

type CLICompletionView struct {
	gameState *game.GameState
	width     int
	height    int
	help      help.Model
	showHelp  bool
}

func CLICompletionViewScreen(gs *game.GameState, width, height int) *CLICompletionView {
	return &CLICompletionView{
		gameState: gs,
		width:     width,
		height:    height,
		help:      help.New(),
		showHelp:  false,
	}
}

func (v *CLICompletionView) Init() tea.Cmd {
	return nil
}

func (v *CLICompletionView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, CLICompletionKeys.Quit):
			return v, tea.Quit

		case key.Matches(msg, CLICompletionKeys.Help):
			v.showHelp = !v.showHelp
			return v, nil
		}

	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
	}

	return v, nil
}

func (v *CLICompletionView) View() string {
	var b strings.Builder

	header := titleStyle.Render("🎯 CLI Challenge Session Completed! 🎯")
	b.WriteString(header + "\n\n")

	if v.gameState.StartedViaCLI {
		var filterDescription strings.Builder

		if v.gameState.UseRandomizedOrder {
			// Difficulty-only mode
			difficultyText := "beginner"
			if v.gameState.Settings.GameMode == "random-by-difficulty" {
				// Try to determine difficulty from randomized challenges
				if len(v.gameState.RandomizedChallenges) > 0 {
					switch v.gameState.RandomizedChallenges[0].Difficulty {
					case 0:
						difficultyText = "beginner"
					case 1:
						difficultyText = "intermediate"
					case 2:
						difficultyText = "advanced"
					}
				}
			}
			filterDescription.WriteString(fmt.Sprintf("You've completed all %s difficulty challenges!", difficultyText))
		} else {
			// Category mode (with or without difficulty filtering)
			currentCategory := v.gameState.GetCurrentCategory()
			filterDescription.WriteString(fmt.Sprintf("You've completed all challenges in the '%s' category!", currentCategory))

			// If there was difficulty filtering, mention it
			if v.gameState.Settings.GameMode == "category" && v.gameState.CurrentCategoryIdx < len(v.gameState.ChallengeSets) {
				challenges := v.gameState.ChallengeSets[v.gameState.CurrentCategoryIdx].Challenges
				if len(challenges) > 0 {
					difficultyText := "all difficulty levels"
					switch challenges[0].Difficulty {
					case 0:
						difficultyText = "beginner level"
					case 1:
						difficultyText = "intermediate level"
					case 2:
						difficultyText = "advanced level"
					}
					filterDescription.WriteString(fmt.Sprintf(" (filtered to %s)", difficultyText))
				}
			}
		}

		b.WriteString(textStyle.Render(filterDescription.String()) + "\n\n")
	}

	message := utils.WrapText("Great job! You've successfully completed your focused challenge session. This targeted practice helps reinforce specific security concepts and builds your expertise in particular areas.", v.width)
	b.WriteString(textStyle.Render(message) + "\n\n")

	b.WriteString(subtitleStyle.Render("🔦 What You've Accomplished") + "\n")
	b.WriteString(textStyle.Render("• Completed focused challenge set") + "\n")
	b.WriteString(textStyle.Render("• Reinforced specific security knowledge") + "\n")
	b.WriteString(textStyle.Render("• Built targeted expertise") + "\n\n")

	b.WriteString(contributionStyle.Render("📚 Keep Learning") + "\n")
	keepLearning := utils.WrapText("Continue building your security knowledge:", v.width)
	b.WriteString(textStyle.Render(keepLearning) + "\n\n")
	b.WriteString("- Try different difficulty levels\n")
	b.WriteString("- Explore other vulnerability categories\n")
	b.WriteString("- Join security communities\n\n")

	if v.gameState != nil {
		completedCount := 0
		if v.gameState.UseRandomizedOrder && len(v.gameState.RandomizedChallenges) > 0 {
			completedCount = len(v.gameState.RandomizedChallenges)
		} else if v.gameState.CurrentCategoryIdx < len(v.gameState.ChallengeSets) {
			completedCount = len(v.gameState.ChallengeSets[v.gameState.CurrentCategoryIdx].Challenges)
		}

		if completedCount > 0 {
			statsMessage := fmt.Sprintf("You completed %d challenges in this focused session!", completedCount)
			b.WriteString(highlightStyle.Render(statsMessage) + "\n\n")
		}
	}

	if v.showHelp {
		b.WriteString("\n" + v.help.View(CLICompletionKeys))
	} else {
		b.WriteString("\n" + helpHintStyle.Render("Press 'Ctrl+C' or 'Q' to quit"))
		b.WriteString("\n" + helpHintStyle.Render("Press ? for more options"))
	}

	return b.String()
}

func (k CLICompletionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k CLICompletionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
	}
}
