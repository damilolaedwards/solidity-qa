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
	// Out describes the output path for the generated invariants.
	Out string `json:"out"`

	// TargetContracts describes the directory that holds the contracts to be fuzzed.
	TargetContracts DirectoryConfig `json:"targetContracts"`

	// FuzzTests describes the directory that holds the fuzz tests.
	FuzzTests DirectoryConfig `json:"fuzzTests"`

	// UnitTests describes the directory that holds the unit tests.
	UnitTests DirectoryConfig `json:"unitTests"`

	// CoverageReportFile describes the path to the coverage report file
	CoverageReportFile string `json:"coverageReportFile"`

	// SupportingFiles describes the paths to the files that provide additional information about the codebase
	SupportingFiles []string `json:"supportingFiles"`

	// NumInvariants describes the number of invariants to generate at a time
	NumInvariants int `json:"numInvariants"`
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
	if p.Out == "" {
		return errors.New("project configuration must specify an output path")
	}

	if p.TargetContracts.Dir == "" {
		return errors.New("project configuration must specify target contracts directory")
	}

	return nil
}
