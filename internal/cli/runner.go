package cli

import (
	"blindspot/internal/challenges"
	"blindspot/internal/game"
	"blindspot/internal/ui"
	"log"
)

type Runner struct {
	config *Config
}

func NewRunner(config *Config) *Runner {
	return &Runner{
		config: config,
	}
}

func (r *Runner) Run() {
	gameState, err := game.NewGameState()
	if err != nil {
		log.Fatal("Failed to initialize game state: ", err)
	}

	r.configureGameState(gameState)

	program, err := ui.InitializeUIWithChallenge(gameState)
	if err != nil {
		log.Fatal("Error initializing UI: ", err)
	}

	if _, err := program.Run(); err != nil {
		log.Fatal("Error running program: ", err)
	}
}

// applies CLI configuration to the game state
func (r *Runner) configureGameState(gameState *game.GameState) {
	if r.config.WasFlagChanged("difficulty") {
		gameState.Settings.GameMode = "random-by-difficulty"
		gameState.UseRandomizedOrder = true

		allChallenges := gameState.GetChallengesGroupedByDifficulty()
		challenges := filterChallengesByDifficulty(allChallenges, r.config.Difficulty)

		if len(challenges) > 0 {
			gameState.RandomizedChallenges = challenges
			gameState.SaveRandomizedOrder()
		}
	}

	gameState.SaveSettings()
}

func filterChallengesByDifficulty(allChallenges []challenges.Challenge, level int) []challenges.Challenge {
	var diffLevel challenges.DifficultyLevel

	// Convert integer level to enum
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
