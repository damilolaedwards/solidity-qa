package slither

import (
	"assistant/config"
	"assistant/types"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:embed get_target_contracts.py
var getTargetContractsScript string

func GetContractsDataUsingSlither(projectConfig *config.ProjectConfig, outputFile *os.File) error {
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
	if _, err := tmpfile.Write([]byte(getTargetContractsScript)); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return nil
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Println("Error closing temporary file:", err)
		return nil
	}

	args := []string{"--target", ".", "--out", outputFile.Name(),
		"--contracts-dir", projectConfig.TargetContracts.Dir}

	if projectConfig.TargetContracts.ExcludePaths != nil && len(projectConfig.TargetContracts.ExcludePaths) > 0 {
		args = append(args, "--exclude-contract-paths",
			strings.Join(projectConfig.TargetContracts.ExcludePaths, ","))
	}

	if projectConfig.TestContracts.Dir != "" {
		args = append(args, "--tests-dir", projectConfig.TestContracts.Dir)
		if projectConfig.TestContracts.ExcludePaths != nil && len(projectConfig.TestContracts.ExcludePaths) > 0 {
			args = append(args, "--exclude-test-paths", strings.Join(projectConfig.TestContracts.ExcludePaths, ","))
		}
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

func ParseContracts(projectConfig *config.ProjectConfig) ([]types.Contract, string, error) {
	var slitherOutput types.SlitherOutput

	// Create a temporary file to hold the slither output
	tmpfile, err := os.CreateTemp("", "slither-output-*.json")
	if err != nil {
		return nil, "", fmt.Errorf("error creating temporary file: %v", err)
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
		return nil, "", fmt.Errorf("unable to read directory")
	}

	// Run slither on project
	err = GetContractsDataUsingSlither(projectConfig, tmpfile)
	if err != nil {
		return nil, "", err
	}

	// Read the contents of the temporary file
	file, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return nil, "", err
	}

	// Parse the slither output
	err = json.Unmarshal(file, &slitherOutput)
	if err != nil {
		return nil, "", err
	}

	contractCodes := getContractCodes(slitherOutput.Contracts)

	// Filter slither output
	filteredContracts := filterSlitherOutput(slitherOutput.Contracts, projectConfig.IncludeInterfaces, projectConfig.IncludeAbstract, projectConfig.IncludeLibraries)

	return filteredContracts, contractCodes, nil
}

func getContractCodes(contracts []types.SlitherContract) string {
	var contractCodes strings.Builder
	for _, contract := range contracts {
		contractCodes.WriteString(contract.Code)
		for _, subContract := range contract.InheritedContracts {
			contractCodes.WriteString(getContractCodes([]types.SlitherContract{subContract}))
		}
	}
	return contractCodes.String()
}

func filterSlitherOutput(slitherContracts []types.SlitherContract, includeInterfaces bool, includeAbstract bool, includeLibraries bool) []types.Contract {
	var filteredContracts []types.Contract

	for _, slitherContract := range slitherContracts {
		if !includeInterfaces && slitherContract.IsInterface {
			continue
		}
		if !includeAbstract && slitherContract.IsAbstract {
			continue
		}
		if !includeLibraries && slitherContract.IsLibrary {
			continue
		}

		filteredContracts = append(filteredContracts, types.Contract{
			ID:                 slitherContract.ID,
			Name:               slitherContract.Name,
			Functions:          slitherContract.Functions,
			InheritedContracts: filterSlitherOutput(slitherContract.InheritedContracts, includeInterfaces, includeAbstract, includeLibraries),
			IsAbstract:         slitherContract.IsAbstract,
			IsInterface:        slitherContract.IsAbstract,
			IsLibrary:          slitherContract.IsAbstract,
		})
	}

	return filteredContracts
}
