package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// ExclusionConfig represents the structure of the YAML exclusion file
type ExclusionConfig struct {
	Exclude []string `yaml:"exclude"`
}

// ParseExclusionConfig reads and parses an exclusion config YAML file, returning a set of excluded paths or an error.
func ParseExclusionConfig(configPath string) (map[string]struct{}, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read exclusion config file: %w", err)
	}

	var config ExclusionConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	exclusions := make(map[string]struct{})
	for _, excludePath := range config.Exclude {
		exclusions[excludePath] = struct{}{}
	}

	return exclusions, nil
}
