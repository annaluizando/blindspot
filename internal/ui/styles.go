package ui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor   = lipgloss.Color("#0066CC") // Deep blue
	secondaryColor = lipgloss.Color("#00AA55")
	accentColor    = lipgloss.Color("#FF8800") // Warning orange
	textColor      = lipgloss.Color("#E0E0E0") // Light gray (works on both light/dark)
	mutedColor     = lipgloss.Color("#888888") // Medium gray
	successColor   = lipgloss.Color("#00CC44") // Success green
	errorColor     = lipgloss.Color("#CC3333") // Error red
	warningColor   = lipgloss.Color("#FFAA00") // Warning yellow
)

// Common styles
var (
	// Base text styles
	baseTextStyle = lipgloss.NewStyle().
			Foreground(textColor).
			MarginBottom(1)

	// Title and header styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			MarginBottom(1)

	// Description styles
	descriptionStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Italic(true).
				MarginBottom(1)

	descStyle = lipgloss.NewStyle().
			Foreground(textColor).
			MarginBottom(1).
			Padding(0, 1)

	itemDescriptionStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				Padding(0, 4)

	// Code display
	codeBoxStyle = lipgloss.NewStyle().
			Padding(1, 5).
			MarginBottom(1)

	// Selection styles
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	unselectedItemStyle = lipgloss.NewStyle().
				Foreground(textColor)

	itemStyle = lipgloss.NewStyle().
			Foreground(textColor)

	correctAnswerStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	incorrectAnswerStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true)

	// Status styles
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(successColor)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(errorColor)

	hintStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Italic(true).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(accentColor)

	// Navigation and help styles
	helpHintStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	explanationStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				Padding(0, 1)

	// Completion and progress styles
	completionStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// Difficulty indicators - Security-themed colors
	difficultyStyle = map[string]lipgloss.Style{
		"beginner":     lipgloss.NewStyle().Foreground(lipgloss.Color("#00CC44")), // Security green
		"intermediate": lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAA00")), // Warning orange
		"advanced":     lipgloss.NewStyle().Foreground(lipgloss.Color("#CC3333")), // Alert red
	}

	// Category and navigation styles
	categoryStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Padding(0, 1)

	// Subtle elements
	subtleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	// Dimmed text for secondary information
	dimStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Faint(true)
)
