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
	StartedViaCLI             bool
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

	settings, err := loadSettings(configDir)
	if err != nil {
		settings = UserSettings{
			ShowVulnerabilityNames: false,
			GameMode:               "category",
		}
	}

	vulnExplanations, err := challenges.LoadVulnerabilityExplanations()
	if err != nil {
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

	if len(progress.RandomizedChallengeIDs) > 0 {
		gs.RandomizedChallenges = gs.restoreRandomizedChallenges(progress.RandomizedChallengeIDs)
	} else {
		gs.RandomizedChallenges = gs.GetChallengesGroupedByDifficulty()
		gs.SaveRandomizedOrder()
	}

	return gs, nil
}

func (gs *GameState) GetVulnerabilityExplanation(category string) (challenges.VulnerabilityInfo, bool) {
	explanation, found := gs.VulnerabilityExplanations[category]
	return explanation, found
}

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
	helpers := NewGameStateHelpers(gs)
	return helpers.IsErrorRecent()
}

func (gs *GameState) IsSuccessMessageRecent() bool {
	helpers := NewGameStateHelpers(gs)
	return helpers.IsSuccessMessageRecent()
}

func (gs *GameState) ToggleShowVulnerabilityNames() {
	gs.Settings.ShowVulnerabilityNames = !gs.Settings.ShowVulnerabilityNames
	if err := gs.SaveSettings(); err != nil {
		gs.SetError(fmt.Errorf("failed to save vulnerability names setting: %w", err))
	}
}

func (gs *GameState) ToggleGameMode() {
	if gs.Settings.GameMode == "category" {
		gs.Settings.GameMode = "random-by-difficulty"
	} else {
		gs.Settings.GameMode = "category"
	}

	gs.UseRandomizedOrder = gs.Settings.GameMode == "random-by-difficulty"

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
		if gs.CurrentChallengeIdx < len(gs.RandomizedChallenges) {
			return gs.RandomizedChallenges[gs.CurrentChallengeIdx]
		}
		return gs.RandomizedChallenges[0]
	}
	if gs.CurrentCategoryIdx >= len(gs.ChallengeSets) {
		gs.CurrentCategoryIdx = 0
		gs.Progress.CurrentCategoryIdx = 0
	}

	if gs.CurrentChallengeIdx >= len(gs.ChallengeSets[gs.CurrentCategoryIdx].Challenges) {
		gs.CurrentChallengeIdx = 0
		gs.Progress.CurrentChallengeIdx = 0
	}

	return gs.ChallengeSets[gs.CurrentCategoryIdx].Challenges[gs.CurrentChallengeIdx]
}

func (gs *GameState) GetCurrentCategory() string {
	if gs.CurrentCategoryIdx >= len(gs.ChallengeSets) {
		gs.CurrentCategoryIdx = 0
		gs.Progress.CurrentCategoryIdx = 0
	}
	return gs.ChallengeSets[gs.CurrentCategoryIdx].Category
}

func (gs *GameState) MoveToNextChallenge() bool {
	if gs.UseRandomizedOrder {
		if gs.CurrentChallengeIdx < len(gs.RandomizedChallenges)-1 {
			gs.CurrentChallengeIdx++
			gs.Progress.CurrentChallengeIdx = gs.CurrentChallengeIdx
			gs.Progress.IsRandomMode = true
			if err := gs.SaveProgress(); err != nil {
				gs.SetError(fmt.Errorf("failed to save progress: %w", err))
			}
			return true
		}
		return false
	}

	currentSet := gs.ChallengeSets[gs.CurrentCategoryIdx]

	if gs.CurrentChallengeIdx < len(currentSet.Challenges)-1 {
		gs.CurrentChallengeIdx++
		gs.Progress.CurrentChallengeIdx = gs.CurrentChallengeIdx
		gs.Progress.IsRandomMode = false
		if err := gs.SaveProgress(); err != nil {
			gs.SetError(fmt.Errorf("failed to save progress: %w", err))
		}
		return true
	}

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

	if startCategoryIdx > 0 || startChallengeIdx > 0 {
		for categoryIdx := range startCategoryIdx {
			if categoryIdx >= len(gs.ChallengeSets) {
				break
			}
			challenges := gs.ChallengeSets[categoryIdx].Challenges
			for challengeIdx := range challenges {
				challenge := challenges[challengeIdx]
				if !gs.IsChallengeCompleted(challenge.ID) {
					return challenge, true
				}
			}
		}

		// Check the start category from startChallengeIdx onwards
		if startCategoryIdx < len(gs.ChallengeSets) {
			challenges := gs.ChallengeSets[startCategoryIdx].Challenges
			for challengeIdx := startChallengeIdx; challengeIdx < len(challenges); challengeIdx++ {
				challenge := challenges[challengeIdx]
				if !gs.IsChallengeCompleted(challenge.ID) {
					return challenge, true
				}
			}
		}
	}

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

func loadSettings(configDir string) (UserSettings, error) {
	settingsPath := filepath.Join(configDir, "settings.json")

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return UserSettings{
			ShowVulnerabilityNames: false,
			GameMode:               "category",
		}, err
	}

	var settings UserSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		return UserSettings{
			ShowVulnerabilityNames: false,
			GameMode:               "category",
		}, err
	}

	return settings, nil
}

func (gs *GameState) SaveProgress() error {
	progressPath := filepath.Join(gs.ConfigDir, "progress.json")

	gs.Progress.CurrentCategoryIdx = gs.CurrentCategoryIdx
	gs.Progress.CurrentChallengeIdx = gs.CurrentChallengeIdx

	data, err := json.Marshal(gs.Progress)
	if err != nil {
		return err
	}

	return os.WriteFile(progressPath, data, 0644)
}

func (gs *GameState) SaveSettings() error {
	settingsPath := filepath.Join(gs.ConfigDir, "settings.json")

	data, err := json.Marshal(gs.Settings)
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
}

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

func (gs *GameState) SaveRandomizedOrder() {
	ids := make([]string, len(gs.RandomizedChallenges))
	for i, challenge := range gs.RandomizedChallenges {
		ids[i] = challenge.ID
	}

	gs.Progress.RandomizedChallengeIDs = ids
	gs.Progress.IsRandomMode = gs.UseRandomizedOrder
	if err := gs.SaveProgress(); err != nil {
		gs.SetError(fmt.Errorf("failed to save randomized order: %w", err))
	}
}

func (gs *GameState) restoreRandomizedChallenges(ids []string) []challenges.Challenge {
	result := make([]challenges.Challenge, 0, len(ids))
	challengeMap := make(map[string]challenges.Challenge)

	for _, set := range gs.ChallengeSets {
		for _, challenge := range set.Challenges {
			challengeMap[challenge.ID] = challenge
		}
	}

	for _, id := range ids {
		if challenge, found := challengeMap[id]; found {
			result = append(result, challenge)
		}
	}

	if len(result) == 0 {
		return gs.GetChallengesGroupedByDifficulty()
	}

	return result
}

func shuffleChallenge(challenges []challenges.Challenge) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(challenges), func(i, j int) {
		challenges[i], challenges[j] = challenges[j], challenges[i]
	})
}

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

	result := append(beginnerChallenges, intermediateChallenges...)
	result = append(result, advancedChallenges...)

	return result
}

func (gs *GameState) ShouldShowVulnerabilityExplanation(category string) bool {
	if gs.UseRandomizedOrder || gs.StartedViaCLI {
		return false
	}
	return gs.GetCategoryCompletionPercentage(category) == 100
}

func (gs *GameState) ShouldReturnToCategoryExplanation() bool {
	if gs.UseRandomizedOrder || gs.StartedViaCLI {
		return false
	}
	return gs.Progress.PendingCategoryExplanation != ""
}

func (gs *GameState) AddErrorCount(challengeCategory string) {
	if gs.Progress.CategoryErrorCounts == nil {
		gs.Progress.CategoryErrorCounts = make(map[string]int)
	}

	gs.Progress.CategoryErrorCounts[challengeCategory]++

	gs.SaveProgress()
}
