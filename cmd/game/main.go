package main

import (
	"log"

	"blindspot/internal/game"
	"blindspot/internal/ui"
)

func main() {
	gameState, err := game.NewGameState()
	if err != nil {
		log.Fatal("Failed to initialize game state: ", err)
	}

	program, err := ui.InitializeUI(gameState)
	if err != nil {
		log.Fatal("Error initializing UI:", err)
	}

	if _, err := program.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}
