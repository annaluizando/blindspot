package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"secure-code-game/internal/challenges"
)

type UserProgress struct {
	CompletedChallenges map[string]bool `json:"completedChallenges"`
	CurrentCategoryIdx  int             `json:"currentCategoryIdx"`
	CurrentChallengeIdx int             `json:"currentChallengeIdx"`
}

type UserSettings struct {
	ShowVulnerabilityNames bool `json:"showVulnerabilityNames"`
}

type GameState struct {
	ChallengeSets             []challenges.ChallengeSet
	CurrentCategoryIdx        int
	CurrentChallengeIdx       int
	Progress                  UserProgress
	Settings                  UserSettings
	ConfigDir                 string
	VulnerabilityExplanations map[string]challenges.VulnerabilityInfo
}

func NewGameState() (*GameState, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	// Load challenges from challenges.yaml
	challengeSets, err := challenges.LoadChallenges()
	if err != nil {
		return nil, err
	}

	// Load or create user progress
	progress, err := loadProgress(configDir)
	if err != nil {
		// If no progress file exists, create a new one
		progress = UserProgress{
			CompletedChallenges: make(map[string]bool),
			CurrentCategoryIdx:  0,
			CurrentChallengeIdx: 0,
		}
	}

	// Load or create user settings
	settings, err := loadSettings(configDir)
	if err != nil {
		// If no settings file exists, create a new one with defaults
		settings = UserSettings{
			ShowVulnerabilityNames: false,
		}
	}

	vulnExplanations, err := challenges.LoadVulnerabilityExplanations()
	if err != nil {
		// Just log the error but continue - this is non-critical
		fmt.Printf("Warning: Could not load vulnerability explanations: %s\n", err)
		vulnExplanations = make(map[string]challenges.VulnerabilityInfo)
	}

	return &GameState{
		ChallengeSets:             challengeSets,
		CurrentCategoryIdx:        progress.CurrentCategoryIdx,
		CurrentChallengeIdx:       progress.CurrentChallengeIdx,
		Progress:                  progress,
		Settings:                  settings,
		ConfigDir:                 configDir,
		VulnerabilityExplanations: vulnExplanations,
	}, nil
}

// Helper method to get explanation for a specific challenge
func (gs *GameState) GetVulnerabilityExplanation(challenge challenges.Challenge) (challenges.VulnerabilityInfo, bool) {
	explanation, found := gs.VulnerabilityExplanations[challenge.Category]
	return explanation, found
}

// Helper method to toggle setting for showing vulnerability names
func (gs *GameState) ToggleShowVulnerabilityNames() {
	gs.Settings.ShowVulnerabilityNames = !gs.Settings.ShowVulnerabilityNames
	gs.SaveSettings()
}

// Checks if a challenge has been completed
func (gs *GameState) IsChallengeCompleted(challengeID string) bool {
	return gs.Progress.CompletedChallenges[challengeID]
}

func (gs *GameState) MarkChallengeCompleted(challengeID string) {
	gs.Progress.CompletedChallenges[challengeID] = true
	// Save progress after marking challenge completed
	gs.SaveProgress()
}

// Calculates completion percentage for a category
func (gs *GameState) GetCategoryCompletionPercentage(category string) int {
	var total, completed int

	// Find the category
	for _, set := range gs.ChallengeSets {
		if set.Category == category {
			total = len(set.Challenges)
			for _, challenge := range set.Challenges {
				if gs.IsChallengeCompleted(challenge.ID) {
					completed++
				}
			}
			break
		}
	}

	if total == 0 {
		return 0
	}

	return (completed * 100) / total
}

// Calculates overall completion percentage
func (gs *GameState) GetTotalCompletionPercentage() int {
	var total, completed int

	for _, set := range gs.ChallengeSets {
		total += len(set.Challenges)
		for _, challenge := range set.Challenges {
			if gs.IsChallengeCompleted(challenge.ID) {
				completed++
			}
		}
	}

	if total == 0 {
		return 0
	}

	return (completed * 100) / total
}

// Returns the current challenge
func (gs *GameState) GetCurrentChallenge() challenges.Challenge {
	return gs.ChallengeSets[gs.CurrentCategoryIdx].Challenges[gs.CurrentChallengeIdx]
}

func (gs *GameState) MoveToNextChallenge() bool {
	currentSet := gs.ChallengeSets[gs.CurrentCategoryIdx]

	// If there are more challenges in current category
	if gs.CurrentChallengeIdx < len(currentSet.Challenges)-1 {
		gs.CurrentChallengeIdx++
		gs.Progress.CurrentChallengeIdx = gs.CurrentChallengeIdx
		gs.SaveProgress()
		return true
	}

	// If there are more categories
	if gs.CurrentCategoryIdx < len(gs.ChallengeSets)-1 {
		gs.CurrentCategoryIdx++
		gs.CurrentChallengeIdx = 0
		gs.Progress.CurrentCategoryIdx = gs.CurrentCategoryIdx
		gs.Progress.CurrentChallengeIdx = gs.CurrentChallengeIdx
		gs.SaveProgress()
		return true
	}

	// No more challenges
	return false
}

func (gs *GameState) GetNextIncompleteChallenge() (challenges.Challenge, bool) {
	// Start from current position
	startCategoryIdx := gs.CurrentCategoryIdx
	startChallengeIdx := gs.CurrentChallengeIdx

	categoryIdx := startCategoryIdx
	for categoryIdx < len(gs.ChallengeSets) {
		challengeIdx := 0
		if categoryIdx == startCategoryIdx {
			challengeIdx = startChallengeIdx
		}

		challenges := gs.ChallengeSets[categoryIdx].Challenges
		for challengeIdx < len(challenges) {
			challenge := challenges[challengeIdx]
			if !gs.IsChallengeCompleted(challenge.ID) {
				return challenge, true
			}
			challengeIdx++
		}
		categoryIdx++
	}

	// If we got here, try from the beginning (in case we started mid-way)
	if startCategoryIdx > 0 || startChallengeIdx > 0 {
		for categoryIdx := 0; categoryIdx <= startCategoryIdx; categoryIdx++ {
			maxChallengeIdx := len(gs.ChallengeSets[categoryIdx].Challenges)
			if categoryIdx == startCategoryIdx {
				maxChallengeIdx = startChallengeIdx
			}

			for challengeIdx := range maxChallengeIdx {
				challenge := gs.ChallengeSets[categoryIdx].Challenges[challengeIdx]
				if !gs.IsChallengeCompleted(challenge.ID) {
					return challenge, true
				}
			}
		}
	}

	// No incomplete challenges found
	return challenges.Challenge{}, false
}

// Loads user progress from file
func loadProgress(configDir string) (UserProgress, error) {
	progressPath := filepath.Join(configDir, "progress.json")

	data, err := os.ReadFile(progressPath)
	if err != nil {
		return UserProgress{
			CompletedChallenges: make(map[string]bool),
		}, err
	}

	var progress UserProgress
	err = json.Unmarshal(data, &progress)
	if err != nil {
		return UserProgress{
			CompletedChallenges: make(map[string]bool),
		}, err
	}

	return progress, nil
}

// Loads user settings from file
func loadSettings(configDir string) (UserSettings, error) {
	settingsPath := filepath.Join(configDir, "settings.json")

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return UserSettings{
			ShowVulnerabilityNames: false, // Default value
		}, err
	}

	var settings UserSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		return UserSettings{
			ShowVulnerabilityNames: false, // Default value
		}, err
	}

	return settings, nil
}

// Saves the current user progress
func (gs *GameState) SaveProgress() error {
	progressPath := filepath.Join(gs.ConfigDir, "progress.json")

	// Update game state in progress struct
	gs.Progress.CurrentCategoryIdx = gs.CurrentCategoryIdx
	gs.Progress.CurrentChallengeIdx = gs.CurrentChallengeIdx

	data, err := json.Marshal(gs.Progress)
	if err != nil {
		return err
	}

	return os.WriteFile(progressPath, data, 0644)
}

// Saves the current user settings
func (gs *GameState) SaveSettings() error {
	settingsPath := filepath.Join(gs.ConfigDir, "settings.json")

	data, err := json.Marshal(gs.Settings)
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
}

// Helper function to get config directory
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".secure-code-game")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}
