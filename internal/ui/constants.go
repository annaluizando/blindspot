package ui

// UI Constants
const (
	// Default dimensions
	DefaultWidth  = 80
	DefaultHeight = 24

	// Viewport margins
	ViewportMargin = 4
	ContentMargin  = 4

	// Help display
	ShortHelpHeight = 1
	FullHelpHeight  = 4

	// Menu types
	MainMenuType      = "main"
	CategoryMenuType  = "category"
	ChallengeMenuType = "challenge"
	ProgressMenuType  = "progress"
	SettingsMenuType  = "settings"

	// Game modes
	CategoryMode           = "category"
	RandomByDifficultyMode = "random-by-difficulty"

	// Difficulty levels
	BeginnerLevel     = 0
	IntermediateLevel = 1
	AdvancedLevel     = 2

	// Error rate thresholds
	HighErrorRate     = 50
	ModerateErrorRate = 30
	LowErrorRate      = 15

	// Time thresholds for messages
	ErrorMessageTimeout   = 10 // seconds
	SuccessMessageTimeout = 5  // seconds

	// Content markers
	OptionsStartMarker    = "What vulnerability is in this code?"
	ResultCorrectMarker   = "✓ Correct! You've identified the vulnerability."
	ResultIncorrectMarker = "✗ Incorrect. Try another option by moving arrow keys!"
)

// Menu item indices
const (
	StartGameIndex = iota
	CategoriesIndex
	ProgressIndex
	SettingsIndex
	ExitIndex
)

// Settings menu indices
const (
	VulnerabilityNamesIndex = iota
	GameModeIndex
	DeleteProgressIndex
	BackToMainMenuIndex
)

// Difficulty level names
var DifficultyNames = map[int]string{
	BeginnerLevel:     "Beginner",
	IntermediateLevel: "Intermediate",
	AdvancedLevel:     "Advanced",
}

// Difficulty level short names
var DifficultyShortNames = map[int]string{
	BeginnerLevel:     "B",
	IntermediateLevel: "I",
	AdvancedLevel:     "A",
}

// Error level names
var ErrorLevelNames = map[string]string{
	"high":     "High",
	"moderate": "Moderate",
	"low":      "Low",
}
