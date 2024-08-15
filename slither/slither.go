package slither

import (
	"assistant/config"
	"assistant/types"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed get_contracts_data.py
var getContractsDataScript string

//go:embed get_target_contracts.py
var getTargetContractsScript string

func GetContractsDataUsingSlither(target string, outputFile *os.File, includeInterfaces bool, includeLibraries bool, includeAbstract bool) error {
	// Run the command
	// Create a temporary file to hold the Python script
	tmpfile, err := os.CreateTemp("", "script-*.py")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return nil
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println("Error removing temporary file:", err)
		}
	}(tmpfile.Name()) // Clean up

	// Write the Python script to the temporary file
	if _, err := tmpfile.Write([]byte(getContractsDataScript)); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return nil
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Println("Error closing temporary file:", err)
		return nil
	}

	args := []string{"--target", target, "--out", outputFile.Name()}

	if includeInterfaces {
		args = append(args, "--include-interfaces", "true")
	}
	if includeLibraries {
		args = append(args, "--include-libraries", "true")
	}
	if includeAbstract {
		args = append(args, "--include-abstract", "true")
	}

	// Prepare the command
	cmd := exec.Command("python3", append([]string{tmpfile.Name()}, args...)...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error running slither: %v\n", err)
		fmt.Printf("stderr: %s\n", output)
		return err
	}

	// Print out slither output
	fmt.Println(string(output))

	return nil
}

func GetContractsData(projectConfig *config.ProjectConfig) (*types.SlitherOutput, error) {
	var contracts []types.Contract

	// Create a temporary file to hold the slither output
	tmpfile, err := os.CreateTemp("", "slither-output-*.json")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println("Error removing temporary file:", err)
		}
	}(tmpfile.Name()) // Clean up

	// Check if provided directory is a directory
	info, err := os.Stat(projectConfig.TargetContracts.Dir)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("unable to read directory")
	}

	// Run slither on project
	err = GetContractsDataUsingSlither(".", tmpfile, projectConfig.IncludeInterfaces, projectConfig.IncludeLibraries, projectConfig.IncludeAbstract)
	if err != nil {
		return nil, err
	}

	// Read the contents of the temporary file
	file, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return nil, err
	}

	// Parse the slither output
	err = json.Unmarshal(file, &contracts)
	if err != nil {
		return nil, err
	}

	// Filter slither output
	contracts = filterSlitherOutput(contracts, projectConfig.TargetContracts.Dir, projectConfig.TargetContracts.ExcludePaths)

	return &types.SlitherOutput{
		Contracts: contracts,
	}, nil
}

func GetTargetContracts(projectConfig *config.ProjectConfig) (string, error) {
	// Create a temporary file to hold the Python script
	tmpfile, err := os.CreateTemp("", "script-*.py")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return "", nil
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println("Error removing temporary file:", err)
		}
	}(tmpfile.Name()) // Clean up

	// Write the Python script to the temporary file
	if _, err := tmpfile.Write([]byte(getTargetContractsScript)); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return "", nil
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Println("Error closing temporary file:", err)
		return "", nil
	}

	args := []string{"--target", projectConfig.TargetContracts.Dir}

	// Prepare the command
	cmd := exec.Command("python3", append([]string{tmpfile.Name()}, args...)...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error running slither: %v\n", err)
		fmt.Printf("stderr: %s\n", output)
		return "", err
	}

	return string(output), nil
}

func filterSlitherOutput(contracts []types.Contract, targetDir string, excludePaths []string) []types.Contract {
	var filteredContracts []types.Contract

	for _, contract := range contracts {
		// Normalize paths
		contractPath := filepath.Clean(contract.Path)
		targetDir = filepath.Clean(targetDir)

		// Check if contract path is under target directory
		if !strings.HasPrefix(contractPath, targetDir) {
			continue
		}

		// Check if contract path or its parent directory is in exclude paths
		excluded := false
		for _, excludePath := range excludePaths {
			excludePath = filepath.Clean(excludePath)
			if contractPath == excludePath || strings.HasPrefix(contractPath, excludePath+string(filepath.Separator)) || strings.HasPrefix(contractPath, filepath.Join(targetDir, excludePath)+string(filepath.Separator)) {
				excluded = true
				break
			}

			// Check parent directories
			parentDir := filepath.Dir(contractPath)
			for parentDir != "." && parentDir != string(filepath.Separator) {
				if parentDir == excludePath {
					excluded = true
					break
				}
				parentDir = filepath.Dir(parentDir)
			}

			if excluded {
				break
			}
		}

		if !excluded {
			filteredContracts = append(filteredContracts, contract)
		}
	}

	return filteredContracts
}
