package ui

import (
	"fmt"
	"strings"

	"blindspot/internal/challenges"
	"blindspot/internal/game"
	"blindspot/internal/utils"
)

type MenuItem struct {
	Title       string
	Description string
	Completed   bool
	ID          string
}

type MenuBuilder struct {
	width  int
	height int
}

func NewMenuBuilder(width, height int) *MenuBuilder {
	return &MenuBuilder{
		width:  width,
		height: height,
	}
}

func (mb *MenuBuilder) UpdateDimensions(width, height int) {
	mb.width = width
	mb.height = height
}

func (mb *MenuBuilder) BuildMenuContent(title, description string, items []MenuItem, cursor int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(title) + "\n\n")
	wrappedDescription := utils.WrapText(description, mb.width)
	b.WriteString(descriptionStyle.Render(wrappedDescription) + "\n\n")

	for i, item := range items {
		mb.writeMenuItem(&b, item, i, cursor)
	}

	return b.String()
}

func (mb *MenuBuilder) writeMenuItem(b *strings.Builder, item MenuItem, index, cursor int) {
	cursorStr := " "
	if cursor == index {
		cursorStr = ">"
	}

	status := " "
	if item.Completed {
		status = "✓"
	}

	line := fmt.Sprintf("%s %s %s", cursorStr, status, item.Title)
	if cursor == index {
		b.WriteString(selectedItemStyle.Render(line) + "\n")

		wrappedItemDescription := utils.WrapText(item.Description, mb.width)
		b.WriteString(itemDescriptionStyle.Render(wrappedItemDescription) + "\n\n")
	} else {
		b.WriteString(itemStyle.Render(line) + "\n")
	}
}

func (mb *MenuBuilder) BuildCategoryItems(gs *game.GameState) []MenuItem {
	items := make([]MenuItem, len(gs.ChallengeSets))

	for i, set := range gs.ChallengeSets {
		completed := gs.GetCategoryCompletionPercentage(set.Category)
		completionText := fmt.Sprintf("[%d%% Complete]", completed)

		difficultyIndicator := mb.buildDifficultyIndicator(set.Challenges)
		enhancedDescription := utils.WrapText(set.Description, mb.width) + "\n" + difficultyIndicator + completionText

		items[i] = MenuItem{
			Title:       set.Category,
			Description: enhancedDescription,
			Completed:   completed == 100,
			ID:          fmt.Sprintf("category-%d", i),
		}
	}

	return items
}

func (mb *MenuBuilder) BuildChallengeItems(gs *game.GameState, category challenges.ChallengeSet) []MenuItem {
	var items []MenuItem

	// Add explanation item
	items = append(items, MenuItem{
		Title:       "📚 See Explanation: " + category.Category,
		Description: "View detailed explanation about this vulnerability type, its impact, and prevention techniques.",
		ID:          "explanation-" + category.Category,
	})

	// Add challenge items
	for _, challenge := range category.Challenges {
		completed := gs.IsChallengeCompleted(challenge.ID)
		difficultyText := mb.getDifficultyText(challenge.Difficulty)

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

	return items
}

func (mb *MenuBuilder) BuildProgressItems(gs *game.GameState) []MenuItem {
	items := make([]MenuItem, len(gs.ChallengeSets))

	for i, set := range gs.ChallengeSets {
		completed := gs.GetCategoryCompletionPercentage(set.Category)
		description := mb.buildProgressDescription(gs, set, completed)

		items[i] = MenuItem{
			Title:       set.Category,
			Description: description,
			Completed:   completed == 100,
			ID:          fmt.Sprintf("progress-category-%d", i),
		}
	}

	return items
}

func (mb *MenuBuilder) BuildSettingsItems(gs *game.GameState) []MenuItem {
	vulnerabilityNamesStatus := "Show"
	if !gs.Settings.ShowVulnerabilityNames {
		vulnerabilityNamesStatus = "Hide"
	}

	orderModeText := "Category Order"
	if gs.Settings.GameMode == RandomByDifficultyMode {
		orderModeText = "Random by Difficulty"
	}

	return []MenuItem{
		{
			Title:       "Vulnerability Names: " + vulnerabilityNamesStatus,
			Description: "Toggle whether vulnerability names are shown during challenges.",
			ID:          "setting-vulnnames",
		},
		{
			Title: "Game Mode: " + orderModeText,
			Description: "Choose how challenges are ordered when playing the game.\n" +
				"Category Order: Play challenges grouped by vulnerabilty category. (Standard Mode)\n" +
				"Random by Difficulty: Play challenges in random order but grouped by difficulty level. (More challenging mode, specially if combined with 'Vulnerability Names: Hide')",
			ID: "setting-ordermode",
		},
		{
			Title: "Delete all progress data",
			Description: "Erases ALL progress data and begin game from start.\n" +
				"!!! Be aware this will make you loose ALL your current progress. \n",
			ID: "setting-deleteprogress",
		},
		{
			Title:       "Back to Main Menu",
			Description: "Return to the main menu",
			ID:          "setting-back",
		},
	}
}

func (mb *MenuBuilder) buildDifficultyIndicator(challengeList []challenges.Challenge) string {
	hasBeginner := false
	hasIntermediate := false
	hasAdvanced := false

	for _, challenge := range challengeList {
		switch challenge.Difficulty {
		case challenges.Beginner:
			hasBeginner = true
		case challenges.Intermediate:
			hasIntermediate = true
		case challenges.Advanced:
			hasAdvanced = true
		}
	}

	var indicator strings.Builder
	if hasBeginner {
		indicator.WriteString(difficultyStyle["beginner"].Render("[B]") + " ")
	}
	if hasIntermediate {
		indicator.WriteString(difficultyStyle["intermediate"].Render("[I]") + " ")
	}
	if hasAdvanced {
		indicator.WriteString(difficultyStyle["advanced"].Render("[A]") + " ")
	}

	return indicator.String()
}

func (mb *MenuBuilder) getDifficultyText(difficulty challenges.DifficultyLevel) string {
	switch difficulty {
	case challenges.Beginner:
		return difficultyStyle["beginner"].Render("[Beginner]")
	case challenges.Intermediate:
		return difficultyStyle["intermediate"].Render("[Intermediate]")
	case challenges.Advanced:
		return difficultyStyle["advanced"].Render("[Advanced]")
	default:
		return ""
	}
}

func (mb *MenuBuilder) buildProgressDescription(gs *game.GameState, set challenges.ChallengeSet, completed int) string {
	var description strings.Builder

	categoryCompleted := 0
	for _, challenge := range set.Challenges {
		if gs.IsChallengeCompleted(challenge.ID) {
			categoryCompleted++
		}
	}

	description.WriteString(fmt.Sprintf("%d of %d challenges completed (%d%%)\n",
		categoryCompleted, len(set.Challenges), completed))

	mb.addDifficultyBreakdown(&description, gs, set.Challenges)

	mb.addErrorStatistics(&description, gs, set.Category, categoryCompleted)

	return description.String()
}

func (mb *MenuBuilder) addDifficultyBreakdown(b *strings.Builder, gs *game.GameState, challengeList []challenges.Challenge) {
	beginnerCount, intermediateCount, advancedCount := 0, 0, 0
	beginnerCompleted, intermediateCompleted, advancedCompleted := 0, 0, 0

	for _, challenge := range challengeList {
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

	hasDifficultyBreakdown := beginnerCount > 0 || intermediateCount > 0 || advancedCount > 0
	if hasDifficultyBreakdown {
		b.WriteString("By Difficulty:\n")

		if beginnerCount > 0 {
			b.WriteString(fmt.Sprintf("    Beginner: %d/%d completed\n",
				beginnerCompleted, beginnerCount))
		}
		if intermediateCount > 0 {
			b.WriteString(fmt.Sprintf("    Intermediate: %d/%d completed\n",
				intermediateCompleted, intermediateCount))
		}
		if advancedCount > 0 {
			b.WriteString(fmt.Sprintf("    Advanced: %d/%d completed\n",
				advancedCompleted, advancedCount))
		}
	}
}

func (mb *MenuBuilder) addErrorStatistics(b *strings.Builder, gs *game.GameState, category string, categoryCompleted int) {
	categoryErrorCount := 0
	if gs.Progress.CategoryErrorCounts != nil {
		categoryErrorCount = gs.Progress.CategoryErrorCounts[category]
	}

	if categoryErrorCount > 0 {
		errorRate := 0
		totalAttempts := categoryCompleted + categoryErrorCount
		if totalAttempts > 0 {
			errorRate = (categoryErrorCount * 100) / totalAttempts
		}

		var errorLevel string
		if errorRate > HighErrorRate {
			errorLevel = ErrorLevelNames["high"]
		} else if errorRate > ModerateErrorRate {
			errorLevel = ErrorLevelNames["moderate"]
		} else if errorRate > LowErrorRate {
			errorLevel = ErrorLevelNames["low"]
		}

		if errorLevel != "" {
			b.WriteString(fmt.Sprintf("Errors in category: %d (%s - %d%% error rate)\n",
				categoryErrorCount, errorLevel, errorRate))
		} else {
			b.WriteString(fmt.Sprintf("Errors in category: %d (%d%% error rate)\n",
				categoryErrorCount, errorRate))
		}
	} else if categoryCompleted > 0 {
		b.WriteString("No errors in this category. Great job!\n")
	}
}
