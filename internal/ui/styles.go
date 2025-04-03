package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Main color palette
var (
	primaryColor   = lipgloss.Color("#4BA8FF")
	secondaryColor = lipgloss.Color("#05B3FF")
	accentColor    = lipgloss.Color("#FFDB58")
	successColor   = lipgloss.Color("#00FF00")
	errorColor     = lipgloss.Color("#FF4040")
	textColor      = lipgloss.Color("#FFFFFF")
	mutedColor     = lipgloss.Color("#888888")
)

var (
	categoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")). // Purple
			Bold(true)

	// Subtle style for separators
	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

// Style for challenge and menu titles
var titleStyle = lipgloss.NewStyle().
	Foreground(primaryColor).
	Bold(true).
	BorderForeground(secondaryColor).
	Padding(0, 1).
	MarginBottom(1)

// Style for challenge descriptions
var descStyle = lipgloss.NewStyle().
	Foreground(textColor).
	MarginBottom(1).
	MaxWidth(100).
	Padding(0, 1)

// Style for menu descriptions
var descriptionStyle = lipgloss.NewStyle().
	Foreground(textColor).
	MaxWidth(100).
	Italic(true).
	MarginBottom(1)

// Style for code display
var codeBoxStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(secondaryColor).
	Padding(1, 2).
	MarginBottom(1).
	MaxWidth(100)

// Style for selected items in multiple choice or menu
var selectedItemStyle = lipgloss.NewStyle().
	Foreground(primaryColor).
	Bold(true)

// Style for unselected items in menu or multiple choice
var itemStyle = lipgloss.NewStyle().
	Foreground(textColor)

// Style for unselected items in list
var unselectedItemStyle = lipgloss.NewStyle().
	Foreground(textColor)

// Style for menu item descriptions when selected
var itemDescriptionStyle = lipgloss.NewStyle().
	Foreground(mutedColor).
	Italic(true).
	Padding(0, 4)

// Style for success message
var successStyle = lipgloss.NewStyle().
	Foreground(successColor).
	Bold(true).
	Padding(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(successColor)

// Style for error/failure message
var errorStyle = lipgloss.NewStyle().
	Foreground(errorColor).
	Bold(true).
	Padding(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(errorColor)

// Style for hints
var hintStyle = lipgloss.NewStyle().
	Foreground(accentColor).
	Italic(true).
	Padding(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(accentColor)

var helpHintStyle = lipgloss.NewStyle().
	Foreground(mutedColor).
	Italic(true)

// Style for completion indicator
var completionStyle = lipgloss.NewStyle().
	Foreground(successColor).
	Bold(true)

// Style for difficulty indicators
var difficultyStyle = map[string]lipgloss.Style{
	"beginner":     lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")), // Green
	"intermediate": lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")), // Yellow
	"advanced":     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")), // Red
}

func RenderBox(content string, title string) string {
	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2)

	if title != "" {
		return boxStyle.BorderTop(true).Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(" "+title+" "),
				content,
			),
		)
	}

	return boxStyle.Render(content)
}
