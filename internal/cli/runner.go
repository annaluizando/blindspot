package cli

import (
	"blindspot/internal/challenges"
	"blindspot/internal/game"
	"blindspot/internal/ui"
	"log"
	"strings"
)

type Runner struct {
	config *Config
}

func NewRunner(config *Config) *Runner {
	return &Runner{config: config}
}

func (r *Runner) Run() {
	gameState, err := game.NewGameState()
	if err != nil {
		log.Fatal("Failed to initialize game state: ", err)
	}

	r.configureGameState(gameState)

	if r.config.WasFlagSet("difficulty") || r.config.WasFlagSet("category") {
		program, err := ui.InitializeUIWithChallenge(gameState)
		if err != nil {
			log.Fatal("Error initializing UI: ", err)
		}

		if _, err := program.Run(); err != nil {
			log.Fatal("Error running program: ", err)
		}
	} else {
		program, err := ui.InitializeUI(gameState)
		if err != nil {
			log.Fatal("Error initializing UI: ", err)
		}

		if _, err := program.Run(); err != nil {
			log.Fatal("Error running program: ", err)
		}
	}
}

func (r *Runner) configureGameState(gameState *game.GameState) {
	// Handle different flag combinations
	if r.config.WasFlagSet("difficulty") && r.config.WasFlagSet("category") && r.config.Category != "" {
		// Both difficulty and category: filter by difficulty within the specified category
		gameState.Settings.GameMode = "category"
		gameState.UseRandomizedOrder = false
		gameState.StartedViaCLI = true

		categoryIndex := findCategoryIndex(gameState, r.config.Category)
		if categoryIndex >= 0 {
			gameState.CurrentCategoryIdx = categoryIndex
			gameState.CurrentChallengeIdx = 0
			gameState.Progress.CurrentCategoryIdx = categoryIndex
			gameState.Progress.CurrentChallengeIdx = 0

			// Filter challenges in this category by difficulty
			categoryChallenges := gameState.ChallengeSets[categoryIndex].Challenges
			filteredChallenges := filterChallengesByDifficulty(categoryChallenges, r.config.Difficulty)

			if len(filteredChallenges) > 0 {
				// Replace the category challenges with filtered ones
				gameState.ChallengeSets[categoryIndex].Challenges = filteredChallenges
			}

			gameState.SaveProgress()
		}
	} else if r.config.WasFlagSet("difficulty") {
		// Only difficulty: use random mode with difficulty filtering
		gameState.Settings.GameMode = "random-by-difficulty"
		gameState.UseRandomizedOrder = true
		gameState.StartedViaCLI = true

		allChallenges := gameState.GetChallengesGroupedByDifficulty()
		challenges := filterChallengesByDifficulty(allChallenges, r.config.Difficulty)

		if len(challenges) > 0 {
			gameState.RandomizedChallenges = challenges
			gameState.SaveRandomizedOrder()
		}
	} else if r.config.WasFlagSet("category") && r.config.Category != "" {
		// Only category flag
		gameState.Settings.GameMode = "category"
		gameState.UseRandomizedOrder = false
		gameState.StartedViaCLI = true

		categoryIndex := findCategoryIndex(gameState, r.config.Category)
		if categoryIndex >= 0 {
			gameState.CurrentCategoryIdx = categoryIndex
			gameState.CurrentChallengeIdx = 0
			gameState.Progress.CurrentCategoryIdx = categoryIndex
			gameState.Progress.CurrentChallengeIdx = 0
			gameState.SaveProgress()
		}
	}

	gameState.SaveSettings()
}

func filterChallengesByDifficulty(allChallenges []challenges.Challenge, level int) []challenges.Challenge {
	var diffLevel challenges.DifficultyLevel

	switch level {
	case 0:
		diffLevel = challenges.Beginner
	case 1:
		diffLevel = challenges.Intermediate
	case 2:
		diffLevel = challenges.Advanced
	default:
		diffLevel = challenges.Beginner
	}

	var filteredChallenges []challenges.Challenge
	for _, challenge := range allChallenges {
		if challenge.Difficulty == diffLevel {
			filteredChallenges = append(filteredChallenges, challenge)
		}
	}

	return filteredChallenges
}

func findCategoryIndex(gs *game.GameState, categoryName string) int {
	categoryName = strings.ToLower(categoryName)

	for i, set := range gs.ChallengeSets {
		if strings.ToLower(set.Category) == categoryName {
			return i
		}
	}

	return -1
}
