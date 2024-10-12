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

type OnChainConfig struct {
	// Enabled describes whether an onchain contract is to be used
	Enabled bool `json:"enabled"`

	// Address describes the address of the onchain contract
	Address string `json:"address"`

	// ApiKey describes the API key to be used for the network
	ApiKey string `json:"apiKey"`

	// NetworkPrefix describes the network prefix of the onchain contract
	NetworkPrefix string `json:"networkPrefix"`

	// ExcludeInterfaces describes whether interfaces will be excluded from the slither output
	ExcludeInterfaces bool `json:"excludeInterfaces"`
}

type ProjectConfig struct {
	// Name describes the project name.
	Name string `json:"name" description:"The project name"`

	// ContractWhitelist describes the only contracts that should be included
	ContractWhitelist []string `json:"contractWhitelist"`

	// TargetContracts describes the directory that holds the contracts to be fuzzed.
	TargetContracts DirectoryConfig `json:"targetContracts"`

	// TestContracts describes the directory that holds the test contracts
	TestContracts DirectoryConfig `json:"testContracts"`

	// OnChainConfig describes the onchain configuration for the project
	OnChainConfig OnChainConfig `json:"-"`

	// Port describes the port that the API will be running on
	Port int `json:"port"`

	// IncludeInterfaces describes whether interfaces will be included in the slither output
	IncludeInterfaces bool `json:"includeInterfaces"`

	// IncludeAbstract describes whether abstract contracts will be included in the slither output
	IncludeAbstract bool `json:"includeAbstract"`

	// IncludeLibraries describes whether libraries will be included in the slither output
	IncludeLibraries bool `json:"includeLibraries"`

	// SlitherArgs describes the extra arguments to be provided to Slither
	SlitherArgs map[string]any `json:"slitherArgs"`
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

// WriteToFile writes the ProjectConfig to a provided file path in a JSON-serialized format with comments.
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
	if p.Name == "" {
		return errors.New("project configuration must specify project name")
	}

	if p.OnChainConfig.Enabled {
		if p.OnChainConfig.Address == "" {
			return errors.New("contract address must be specified in onchain mode")
		}
		if p.OnChainConfig.ApiKey == "" {
			return errors.New("API key must be specified in onchain mode")
		}
	} else {
		if p.TargetContracts.Dir == "" {
			return errors.New("target must be specified in local mode")
		}
	}

	return nil
}
