package config

import (
	"encoding/json"
	"errors"
	"os"
)

type DirectoryConfig struct {
	// Dir describes the path to the directory
	Dir string `json:"directory"`

	// ExcludePaths describes the paths that should be excluded when parsing the directory
	ExcludePaths []string `json:"excludePaths"`
}

type ProjectConfig struct {
	// TargetContracts describes the directory that holds the contracts to be fuzzed.
	TargetContracts DirectoryConfig `json:"targetContracts"`

	// Port describes the port that the API will be running on
	Port int `json:"port"`

	// IncludeInterfaces describes whether interfaces will be included in the slither output
	IncludeInterfaces bool `json:"includeInterfaces"`

	// IncludeAbstract describes whether abstract contracts will be included in the slither output
	IncludeAbstract bool `json:"includeAbstract"`

	// IncludeLibraries describes whether libraries will be included in the slither output
	IncludeLibraries bool `json:"includeLibraries"`
}

// ReadProjectConfigFromFile reads a JSON-serialized ProjectConfig from a provided file path.
// Returns the ProjectConfig if it succeeds, or an error if one occurs.
func ReadProjectConfigFromFile(path string) (*ProjectConfig, error) {
	// Read our project configuration file data
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse the project configuration
	projectConfig, err := GetDefaultProjectConfig()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, projectConfig)
	if err != nil {
		return nil, err
	}

	return projectConfig, nil
}

// WriteToFile writes the ProjectConfig to a provided file path in a JSON-serialized format.
// Returns an error if one occurs.
func (p *ProjectConfig) WriteToFile(path string) error {
	// Serialize the configuration
	b, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err
	}

	// Save it to the provided output path and return the result
	err = os.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Validate validates that the ProjectConfig meets certain requirements.
// Returns an error if one occurs.
func (p *ProjectConfig) Validate() error {
	if p.TargetContracts.Dir == "" {
		return errors.New("project configuration must specify target contracts directory")
	}

	return nil
}
