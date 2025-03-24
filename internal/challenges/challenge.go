package challenges

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ChallengeType represents different types of security challenges
type ChallengeType int

const (
	MultipleChoice ChallengeType = iota
	CodeFix
)

type DifficultyLevel int

const (
	Beginner DifficultyLevel = iota
	Intermediate
	Advanced
)

// Challenge represents a security coding challenge
type Challenge struct {
	ID            string          `yaml:"id"`            // Unique identifier
	Title         string          `yaml:"title"`         // Challenge title
	Description   string          `yaml:"description"`   // Challenge description
	Type          ChallengeType   `yaml:"type"`          // Multiple choice or code fix
	Difficulty    DifficultyLevel `yaml:"difficulty"`    // Difficulty level
	Category      string          `yaml:"category"`      // Security category (e.g., "SQL Injection")
	Code          string          `yaml:"code"`          // The vulnerable code to display
	Options       []string        `yaml:"options"`       // For multiple choice: possible answers
	CorrectAnswer string          `yaml:"correctAnswer"` // For multiple choice: correct option
	Hint          string          `yaml:"hint"`          // Optional hint for the user
	Solution      string          `yaml:"solution"`      // For code fix: a sample correct solution
	Tags          []string        `yaml:"tags"`          // Tags for filtering challenges
}

// ChallengeSet represents a group of related challenges
type ChallengeSet struct {
	Category    string      `yaml:"category"`    // Category name
	Description string      `yaml:"description"` // Category description
	Challenges  []Challenge `yaml:"challenges"`  // Challenges in this category
}

// YAMLChallenges represents the structure of the YAML file
type YAMLChallenges struct {
	ChallengeSets []ChallengeSet `yaml:"challengeSets"` // Changed from "Sets" to "challengeSets"
}

// LoadChallenges loads all challenges from the YAML file
func LoadChallenges() ([]ChallengeSet, error) {
	// Look for challenges.yaml in multiple locations
	searchPaths := []string{
		"assets/challenges.yaml",       // From project root
		"../assets/challenges.yaml",    // If running from cmd/security-game
		"../../assets/challenges.yaml", // If running from elsewhere
		"./challenges.yaml",            // Current directory
		"challenges.yaml",              // Also try just the filename directly
	}

	var yamlData []byte
	var err error

	for _, path := range searchPaths {
		yamlData, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		// If we still didn't find the file, try to find it relative to the executable
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			yamlPath := filepath.Join(execDir, "assets/challenges.yaml")
			yamlData, err = os.ReadFile(yamlPath)
			if err == nil {
				fmt.Printf("Error in reading yaml file: %s\n", err)
			}
		}

		// If still not found, return error
		if err != nil {
			return nil, err
		}
	}

	// Parse the YAML file
	var challengeData YAMLChallenges
	err = yaml.Unmarshal(yamlData, &challengeData)
	if err != nil {
		return nil, err
	}

	// Add debugging to confirm what was loaded
	if len(challengeData.ChallengeSets) > 0 {
		fmt.Printf("Loaded %d challenge sets\n", len(challengeData.ChallengeSets))
	}

	return challengeData.ChallengeSets, nil
}
