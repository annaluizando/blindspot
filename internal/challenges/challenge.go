package challenges

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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

type Challenge struct {
	ID            string          `yaml:"id"`
	Title         string          `yaml:"title"`
	Description   string          `yaml:"description"`
	Type          ChallengeType   `yaml:"type"` // Multiple choice or code fix
	Difficulty    DifficultyLevel `yaml:"difficulty"`
	Code          string          `yaml:"code"`          // The vulnerable code to display
	Options       []string        `yaml:"options"`       // For multiple choice: possible answers
	CorrectAnswer string          `yaml:"correctAnswer"` // For multiple choice: correct option
	Hint          string          `yaml:"hint"`          // Optional hint for the user
	Solution      string          `yaml:"solution"`      // For code fix: a sample correct solution
	Lang          string          `yaml:"lang"`          // programming language in challenge's code
	Explanation   string          `yaml:"explanation"`   // Explanation of why the correct answer is right
}

// group of related challenges
type ChallengeSet struct {
	Category    string      `yaml:"category"` // Category name
	ID          string      `yaml:"id"`
	Description string      `yaml:"description"` // Category description
	Challenges  []Challenge `yaml:"challenges"`  // Challenges in this category
}

// structure of the YAML file
type YAMLChallenges struct {
	ChallengeSets []ChallengeSet `yaml:"challengeSets"`
}

// loads all challenges from the YAML file
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

func GetChallengesCategories() ([]string, error) {
	challengeSets, err := LoadChallenges()
	if err != nil {
		return nil, err
	}

	var categories []string
	for _, challenge := range challengeSets {
		categories = append(categories, challenge.Category)
	}

	return categories, nil
}
