package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type DirectoryConfig struct {
	// Dir describes the path to the directory
	Dir string `json:"directory" description:"The directory path relative to the project root"`

	// ExcludePaths describes the paths that should be excluded when parsing the directory
	ExcludePaths []string `json:"excludePaths" description:"Paths that should be excluded when parsing the directory"`
}

type OnChainConfig struct {
	// Enabled describes whether an onchain contract is to be used
	Enabled bool `json:"-"`

	// Address describes the address of the onchain contract
	Address string `json:"-"`

	// ApiKey describes the API key to be used for the network
	ApiKey string `json:"-"`

	// NetworkPrefix describes the network prefix of the onchain contract
	NetworkPrefix string `json:"-"`

	// ExcludeInterfaces describes whether interfaces will be excluded from the slither output
	ExcludeInterfaces bool `json:"-"`
}

type ProjectConfig struct {
	// Name describes the project name.
	Name string `json:"name" description:"The project name"`

	// TargetContracts describes the directory that holds the contracts to be fuzzed.
	TargetContracts DirectoryConfig `json:"targetContracts" description:"The directory that holds the contracts to be fuzzed"`

	// TestContracts describes the directory that holds the test contracts
	TestContracts DirectoryConfig `json:"testContracts" description:"The directory that holds the test contracts"`

	// OnChainConfig describes the onchain configuration for the project
	OnChainConfig OnChainConfig `json:"-"`

	// Port describes the port that the API will be running on
	Port int `json:"port" description:"The port that the API will be running on"`

	// IncludeInterfaces describes whether interfaces will be included in the slither output
	IncludeInterfaces bool `json:"includeInterfaces" description:"Whether interfaces will be included in the slither output"`

	// IncludeAbstract describes whether abstract contracts will be included in the slither output
	IncludeAbstract bool `json:"includeAbstract" description:"Whether abstract contracts will be included in the slither output"`

	// IncludeLibraries describes whether libraries will be included in the slither output
	IncludeLibraries bool `json:"includeLibraries" description:"Whether libraries will be included in the slither output"`
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

// customMarshal is a helper function to marshal a struct with inline comments
func customMarshal(v interface{}) ([]byte, error) {
	var result strings.Builder
	result.WriteString("{\n")

	t := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		value = value.Elem()
	}

	fields := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
		description := field.Tag.Get("description")
		fieldValue := value.Field(i).Interface()

		// Skip fields with json:"-" tag
		if jsonTag == "-" {
			continue
		}

		if jsonTag != "" {
			var fieldStr string
			if field.Type.Kind() == reflect.Struct {
				// For nested structs, recursively call customMarshal
				nestedJSON, err := customMarshal(fieldValue)
				if err != nil {
					return nil, err
				}
				fieldStr = fmt.Sprintf("\t%q: %s // %s", jsonTag, string(nestedJSON), description)
			} else {
				fieldJSON, err := json.Marshal(fieldValue)
				if err != nil {
					return nil, err
				}
				fieldStr = fmt.Sprintf("\t%q: %s // %s", jsonTag, string(fieldJSON), description)
			}
			fields = append(fields, fieldStr)
		}
	}

	result.WriteString(strings.Join(fields, ",\n"))
	result.WriteString("\n}")

	return []byte(result.String()), nil
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
