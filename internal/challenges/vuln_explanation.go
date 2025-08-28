package challenges

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Resource struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

type VulnerabilityInfo struct {
	Name             string     `yaml:"name"`
	ShortDescription string     `yaml:"short_description"`
	Explanation      string     `yaml:"explanation"`
	Resources        []Resource `yaml:"resources"`
}

type VulnerabilitiesData struct {
	Vulnerabilities []VulnerabilityInfo `yaml:"vulnerabilities"`
}

func LoadVulnerabilityExplanations() (map[string]VulnerabilityInfo, error) {
	searchPaths := []string{
		"assets/vuln_explanations.yaml",
		"../assets/vuln_explanations.yaml",
		"../../assets/vuln_explanations.yaml",
		"./vuln_explanations.yaml",
		"vuln_explanations.yaml",
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
			yamlPath := filepath.Join(execDir, "assets/vulnerabilities.yaml")
			yamlData, err = os.ReadFile(yamlPath)
		}

		if err != nil {
			return nil, fmt.Errorf("could not find vulnerabilities.yaml: %w", err)
		}
	}

	var vulnData VulnerabilitiesData
	err = yaml.Unmarshal(yamlData, &vulnData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vulnerabilities.yaml: %w", err)
	}

	vulnMap := make(map[string]VulnerabilityInfo)
	for _, vuln := range vulnData.Vulnerabilities {
		vulnMap[vuln.Name] = vuln
	}

	return vulnMap, nil
}
