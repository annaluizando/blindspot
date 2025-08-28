package game

import (
	"blindspot/internal/challenges"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type UserProgress struct {
	CompletedChallenges        map[string]bool `json:"completedChallenges"`
	CurrentCategoryIdx         int             `json:"currentCategoryIdx"`
	CurrentChallengeIdx        int             `json:"currentChallengeIdx"`
	RandomizedChallengeIDs     []string        `json:"randomizedChallengeIDs"`
	CategoryErrorCounts        map[string]int  `json:"categoryErrorCounts"`
	IsRandomMode               bool            `json:"isRandomMode"`
	PendingCategoryExplanation string          `json:"pendingCategoryExplanation"`
}

type UserSettings struct {
	ShowVulnerabilityNames bool   `json:"showVulnerabilityNames"`
	GameMode               string `json:"gameMode"`
}

type GameState struct {
	ChallengeSets             []challenges.ChallengeSet
	CurrentCategoryIdx        int
	CurrentChallengeIdx       int
	Progress                  UserProgress
	Settings                  UserSettings
	ConfigDir                 string
	VulnerabilityExplanations map[string]challenges.VulnerabilityInfo
	RandomizedChallenges      []challenges.Challenge
	UseRandomizedOrder        bool
	LastError                 string
	ErrorTimestamp            time.Time
	LastSuccessMessage        string
	SuccessMessageTimestamp   time.Time
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

	progress, err := loadProgress(configDir)
	if err != nil {
		progress = UserProgress{
			CompletedChallenges:    make(map[string]bool),
			CurrentCategoryIdx:     0,
			CurrentChallengeIdx:    0,
			RandomizedChallengeIDs: []string{},
		}
	}

	// Load or create user settings
	settings, err := loadSettings(configDir)
	if err != nil {
		// If no settings file exists, create a new one with defaults
		settings = UserSettings{
			ShowVulnerabilityNames: false,
			GameMode:               "category",
		}
	}

	vulnExplanations, err := challenges.LoadVulnerabilityExplanations()
	if err != nil {
		// Just log the error but continue - this is non-critical
		fmt.Printf("Warning: Could not load vulnerability explanations: %s\n", err)
		vulnExplanations = make(map[string]challenges.VulnerabilityInfo)
	}

	gs := &GameState{
		ChallengeSets:             challengeSets,
		CurrentCategoryIdx:        progress.CurrentCategoryIdx,
		CurrentChallengeIdx:       progress.CurrentChallengeIdx,
		Progress:                  progress,
		Settings:                  settings,
		ConfigDir:                 configDir,
		VulnerabilityExplanations: vulnExplanations,
		UseRandomizedOrder:        settings.GameMode == "random-by-difficulty",
	}

	// Generate or restore the randomized challenges
	if len(progress.RandomizedChallengeIDs) > 0 {
		// Restore the previously saved randomized order
		gs.RandomizedChallenges = gs.restoreRandomizedChallenges(progress.RandomizedChallengeIDs)
	} else {
		// Generate a new randomized order
		gs.RandomizedChallenges = gs.GetChallengesGroupedByDifficulty()
		// Save the order of IDs to progress
		gs.SaveRandomizedOrder()
	}

	return gs, nil
}

func (gs *GameState) GetVulnerabilityExplanation(category string) (challenges.VulnerabilityInfo, bool) {
	explanation, found := gs.VulnerabilityExplanations[category]
	return explanation, found
}

// SetPendingCategoryExplanation marks that the user should return to a category explanation
func (gs *GameState) SetPendingCategoryExplanation(category string) {
	gs.Progress.PendingCategoryExplanation = category
	gs.SaveProgress()
}

func (gs *GameState) ClearPendingCategoryExplanation() {
	gs.Progress.PendingCategoryExplanation = ""
	gs.SaveProgress()
}

func (gs *GameState) GetPendingCategoryExplanation() string {
	return gs.Progress.PendingCategoryExplanation
}

func (gs *GameState) ShouldReturnToCategoryExplanation() bool {
	return gs.Progress.PendingCategoryExplanation != ""
}

func (gs *GameState) SetError(err error) {
	if err != nil {
		gs.LastError = err.Error()
		gs.ErrorTimestamp = time.Now()
	}
}

func (gs *GameState) SetSuccessMessage(message string) {
	gs.LastSuccessMessage = message
	gs.SuccessMessageTimestamp = time.Now()
}

func (gs *GameState) ClearError() {
	gs.LastError = ""
}

func (gs *GameState) ClearSuccessMessage() {
	gs.LastSuccessMessage = ""
}

func (gs *GameState) ClearMessages() {
	gs.ClearError()
	gs.ClearSuccessMessage()
}

func (gs *GameState) HasError() bool {
	return gs.LastError != ""
}

func (gs *GameState) HasSuccessMessage() bool {
	return gs.LastSuccessMessage != ""
}

func (gs *GameState) GetError() string {
	return gs.LastError
}

func (gs *GameState) GetSuccessMessage() string {
	return gs.LastSuccessMessage
}

func (gs *GameState) IsErrorRecent() bool {
	return time.Since(gs.ErrorTimestamp) < 10*time.Second
}

func (gs *GameState) IsSuccessMessageRecent() bool {
	return time.Since(gs.SuccessMessageTimestamp) < 5*time.Second
}

func (gs *GameState) ToggleShowVulnerabilityNames() {
	gs.Settings.ShowVulnerabilityNames = !gs.Settings.ShowVulnerabilityNames
	if err := gs.SaveSettings(); err != nil {
		gs.SetError(fmt.Errorf("failed to save vulnerability names setting: %w", err))
	}
}

// Helper method to toggle challenge order setting
func (gs *GameState) ToggleGameMode() {
	if gs.Settings.GameMode == "category" {
		gs.Settings.GameMode = "random-by-difficulty"
	} else {
		gs.Settings.GameMode = "category"
	}

	// Update the UseRandomizedOrder flag to match the setting
	gs.UseRandomizedOrder = gs.Settings.GameMode == "random-by-difficulty"

	// If we're switching to random mode and don't have randomized challenges yet, generate them
	if gs.UseRandomizedOrder && len(gs.RandomizedChallenges) == 0 {
		gs.RandomizedChallenges = gs.GetChallengesGroupedByDifficulty()
		gs.SaveRandomizedOrder()
	}

	if err := gs.SaveSettings(); err != nil {
		gs.SetError(fmt.Errorf("failed to save game mode setting: %w", err))
	}
}

func (gs *GameState) IsChallengeCompleted(challengeID string) bool {
	return gs.Progress.CompletedChallenges[challengeID]
}

func (gs *GameState) MarkChallengeCompleted(challengeID string) {
	gs.Progress.CompletedChallenges[challengeID] = true
	if err := gs.SaveProgress(); err != nil {
		gs.SetError(fmt.Errorf("failed to save challenge completion: %w", err))
	}
}

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

func (gs *GameState) GetCurrentChallenge() challenges.Challenge {
	if gs.UseRandomizedOrder && len(gs.RandomizedChallenges) > 0 {
		// When in randomized mode, use the currentChallengeIdx directly
		// on the randomized list (making sure not to go out of bounds)
		if gs.CurrentChallengeIdx < len(gs.RandomizedChallenges) {
			return gs.RandomizedChallenges[gs.CurrentChallengeIdx]
		}
		// If out of bounds, return the first challenge
		return gs.RandomizedChallenges[0]
	}
	// Otherwise, use the original order by category
	// Add bounds checking to prevent panics
	if gs.CurrentCategoryIdx >= len(gs.ChallengeSets) {
		// If category index is out of bounds, reset to first category
		gs.CurrentCategoryIdx = 0
		gs.Progress.CurrentCategoryIdx = 0
	}

	if gs.CurrentChallengeIdx >= len(gs.ChallengeSets[gs.CurrentCategoryIdx].Challenges) {
		// If challenge index is out of bounds, reset to first challenge
		gs.CurrentChallengeIdx = 0
		gs.Progress.CurrentChallengeIdx = 0
	}

	return gs.ChallengeSets[gs.CurrentCategoryIdx].Challenges[gs.CurrentChallengeIdx]
}

func (gs *GameState) GetCurrentCategory() string {
	if gs.CurrentCategoryIdx >= len(gs.ChallengeSets) {
		// If category index is out of bounds, reset to first category
		gs.CurrentCategoryIdx = 0
		gs.Progress.CurrentCategoryIdx = 0
	}
	return gs.ChallengeSets[gs.CurrentCategoryIdx].Category
}

func (gs *GameState) MoveToNextChallenge() bool {
	if gs.UseRandomizedOrder {
		// In randomized mode, just increment the challenge index
		if gs.CurrentChallengeIdx < len(gs.RandomizedChallenges)-1 {
			gs.CurrentChallengeIdx++
			gs.Progress.CurrentChallengeIdx = gs.CurrentChallengeIdx
			gs.Progress.IsRandomMode = true
			if err := gs.SaveProgress(); err != nil {
				gs.SetError(fmt.Errorf("failed to save progress: %w", err))
			}
			return true
		}
		// No more challenges in randomized list
		return false
	}

	// Original behavior for category-based navigation
	currentSet := gs.ChallengeSets[gs.CurrentCategoryIdx]

	// If there are more challenges in current category
	if gs.CurrentChallengeIdx < len(currentSet.Challenges)-1 {
		gs.CurrentChallengeIdx++
		gs.Progress.CurrentChallengeIdx = gs.CurrentChallengeIdx
		gs.Progress.IsRandomMode = false
		if err := gs.SaveProgress(); err != nil {
			gs.SetError(fmt.Errorf("failed to save progress: %w", err))
		}
		return true
	}

	// If there are more categories
	if gs.CurrentCategoryIdx < len(gs.ChallengeSets)-1 {
		gs.CurrentCategoryIdx++
		gs.CurrentChallengeIdx = 0
		gs.Progress.CurrentCategoryIdx = gs.CurrentCategoryIdx
		gs.Progress.CurrentChallengeIdx = gs.CurrentChallengeIdx
		gs.Progress.IsRandomMode = false
		if err := gs.SaveProgress(); err != nil {
			gs.SetError(fmt.Errorf("failed to save progress: %w", err))
		}
		return true
	}

	// No more challenges
	return false
}

func (gs *GameState) GetNextIncompleteChallenge() (challenges.Challenge, bool) {
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

func (gs *GameState) EraseProgressData() {
	configDir, err := getConfigDir()
	if err != nil {
		gs.SetError(fmt.Errorf("failed to get config directory: %w", err))
		return
	}

	progressPath := filepath.Join(configDir, "progress.json")
	err = os.Remove(progressPath)
	if err != nil {
		gs.SetError(fmt.Errorf("failed to remove progress file: %w", err))
		return
	}

	gs.resetProgress()
	gs.SetSuccessMessage("Progress data cleared successfully")
}

func (gs *GameState) resetProgress() {
	gs.Progress = UserProgress{
		CompletedChallenges:    make(map[string]bool),
		CurrentCategoryIdx:     0,
		CurrentChallengeIdx:    0,
		RandomizedChallengeIDs: []string{},
		IsRandomMode:           false,
	}

	gs.CurrentCategoryIdx = 0
	gs.CurrentChallengeIdx = 0

	gs.RandomizedChallenges = []challenges.Challenge{}
	gs.UseRandomizedOrder = false
}

func loadProgress(configDir string) (UserProgress, error) {
	progressPath := filepath.Join(configDir, "progress.json")

	data, err := os.ReadFile(progressPath)
	if err != nil {
		return UserProgress{
			CompletedChallenges: make(map[string]bool),
			CategoryErrorCounts: make(map[string]int),
		}, err
	}

	var progress UserProgress
	err = json.Unmarshal(data, &progress)
	if err != nil {
		return UserProgress{
			CompletedChallenges: make(map[string]bool),
			CategoryErrorCounts: make(map[string]int),
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
			ShowVulnerabilityNames: false,      // Default value
			GameMode:               "category", // Default to category mode
		}, err
	}

	var settings UserSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		return UserSettings{
			ShowVulnerabilityNames: false,      // Default value
			GameMode:               "category", // Default to category mode
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

	configDir := filepath.Join(homeDir, ".blindspot-game")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// Helper method to save the current randomized order to progress
func (gs *GameState) SaveRandomizedOrder() {
	// Extract challenge IDs in the current randomized order
	ids := make([]string, len(gs.RandomizedChallenges))
	for i, challenge := range gs.RandomizedChallenges {
		ids[i] = challenge.ID
	}

	// Save to progress
	gs.Progress.RandomizedChallengeIDs = ids
	gs.Progress.IsRandomMode = gs.UseRandomizedOrder
	if err := gs.SaveProgress(); err != nil {
		gs.SetError(fmt.Errorf("failed to save randomized order: %w", err))
	}
}

// Helper method to restore challenges from saved IDs
func (gs *GameState) restoreRandomizedChallenges(ids []string) []challenges.Challenge {
	result := make([]challenges.Challenge, 0, len(ids))
	challengeMap := make(map[string]challenges.Challenge)

	// Create a map of all challenges by ID for quick lookup
	for _, set := range gs.ChallengeSets {
		for _, challenge := range set.Challenges {
			challengeMap[challenge.ID] = challenge
		}
	}

	// Restore the challenges in the saved order
	for _, id := range ids {
		if challenge, found := challengeMap[id]; found {
			result = append(result, challenge)
		}
	}

	// If the restored list is empty (which shouldn't happen), generate a new one
	if len(result) == 0 {
		return gs.GetChallengesGroupedByDifficulty()
	}

	return result
}

// Helper function to shuffle a slice of challenges
func shuffleChallenge(challenges []challenges.Challenge) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(challenges), func(i, j int) {
		challenges[i], challenges[j] = challenges[j], challenges[i]
	})
}

// returns all challenges sorted by difficulty
// and randomized within each difficulty group
func (gs *GameState) GetChallengesGroupedByDifficulty() []challenges.Challenge {
	beginnerChallenges := []challenges.Challenge{}
	intermediateChallenges := []challenges.Challenge{}
	advancedChallenges := []challenges.Challenge{}

	for _, set := range gs.ChallengeSets {
		for _, challenge := range set.Challenges {
			switch challenge.Difficulty {
			case challenges.Beginner:
				beginnerChallenges = append(beginnerChallenges, challenge)
			case challenges.Intermediate:
				intermediateChallenges = append(intermediateChallenges, challenge)
			case challenges.Advanced:
				advancedChallenges = append(advancedChallenges, challenge)
			}
		}
	}

	shuffleChallenge(beginnerChallenges)
	shuffleChallenge(intermediateChallenges)
	shuffleChallenge(advancedChallenges)

	// Combine challenges in order of difficulty (beginner -> intermediate -> advanced)
	result := append(beginnerChallenges, intermediateChallenges...)
	result = append(result, advancedChallenges...)

	return result
}

func (gs *GameState) ShouldShowVulnerabilityExplanation(category string) bool {
	if gs.UseRandomizedOrder {
		return false
	}

	return gs.GetCategoryCompletionPercentage(category) == 100
}

func (gs *GameState) AddErrorCount(challengeCategory string) {
	if gs.Progress.CategoryErrorCounts == nil {
		gs.Progress.CategoryErrorCounts = make(map[string]int)
	}

	gs.Progress.CategoryErrorCounts[challengeCategory]++

	gs.SaveProgress()
}
