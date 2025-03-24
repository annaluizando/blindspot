package challenges

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Resource represents an educational resource about a vulnerability
type Resource struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

// VulnerabilityInfo contains detailed information about a security vulnerability
type VulnerabilityInfo struct {
	Name             string     `yaml:"name"`
	ShortDescription string     `yaml:"short_description"`
	Explanation      string     `yaml:"explanation"`
	Resources        []Resource `yaml:"resources"`
}

// VulnerabilitiesData represents the structure of the vulnerabilities.yaml file
type VulnerabilitiesData struct {
	Vulnerabilities []VulnerabilityInfo `yaml:"vulnerabilities"`
}

// LoadVulnerabilityExplanations loads vulnerability explanations from the YAML file
func LoadVulnerabilityExplanations() (map[string]VulnerabilityInfo, error) {
	searchPaths := []string{
		"assets/vuln_explanations.yaml",       // From project root
		"../assets/vuln_explanations.yaml",    // If running from cmd/security-game
		"../../assets/vuln_explanations.yaml", // If running from elsewhere
		"./vuln_explanations.yaml",            // Current directory
		"vuln_explanations.yaml",              // Also try just the filename directly
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
		// Try to find it relative to the executable as a fallback
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			yamlPath := filepath.Join(execDir, "assets/vulnerabilities.yaml")
			yamlData, err = os.ReadFile(yamlPath)
		}

		// If still not found, return error
		if err != nil {
			return nil, fmt.Errorf("could not find vulnerabilities.yaml: %w", err)
		}
	}

	// Parse the YAML file
	var vulnData VulnerabilitiesData
	err = yaml.Unmarshal(yamlData, &vulnData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vulnerabilities.yaml: %w", err)
	}

	// Create a map for easy lookup by vulnerability name
	vulnMap := make(map[string]VulnerabilityInfo)
	for _, vuln := range vulnData.Vulnerabilities {
		vulnMap[vuln.Name] = vuln
	}

	fmt.Printf("Loaded %d vulnerability explanations\n", len(vulnMap))
	return vulnMap, nil
}
