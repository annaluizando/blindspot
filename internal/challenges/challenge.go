package challenges

import (
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
	Type          ChallengeType   `yaml:"type"`
	Difficulty    DifficultyLevel `yaml:"difficulty"`
	Code          string          `yaml:"code"`
	Options       []string        `yaml:"options"`
	CorrectAnswer string          `yaml:"correctAnswer"`
	Hint          string          `yaml:"hint"`
	Solution      string          `yaml:"solution"`
	Lang          string          `yaml:"lang"`
	Explanation   string          `yaml:"explanation"`
}

type ChallengeSet struct {
	Category    string      `yaml:"category"`
	ID          string      `yaml:"id"`
	Description string      `yaml:"description"`
	Challenges  []Challenge `yaml:"challenges"`
}

type YAMLChallenges struct {
	ChallengeSets []ChallengeSet `yaml:"challengeSets"`
}

func LoadChallenges() ([]ChallengeSet, error) {
	searchPaths := []string{
		"assets/challenges.yaml",
		"../assets/challenges.yaml",
		"../../assets/challenges.yaml",
		"./challenges.yaml",
		"challenges.yaml",
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
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			yamlPath := filepath.Join(execDir, "assets/challenges.yaml")
			yamlData, err = os.ReadFile(yamlPath)
		}

		if err != nil {
			return nil, err
		}
	}

	var challengeData YAMLChallenges
	err = yaml.Unmarshal(yamlData, &challengeData)
	if err != nil {
		return nil, err
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
