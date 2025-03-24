package main

import (
	"fmt"
	"log"
	"os"

	"secure-code-game/internal/game"
	"secure-code-game/internal/ui"
)

func main() {
	gameState, err := game.NewGameState()
	if err != nil {
		log.Fatal("Failed to initialize game state: ", err)
	}

	program, err := ui.InitializeUI(gameState)
	if err != nil {
		fmt.Println("Error initializing UI:", err)
		os.Exit(1)
	}

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
