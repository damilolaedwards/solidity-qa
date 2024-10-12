package slither

import (
	"assistant/config"
	"assistant/types"
	"assistant/utils"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//go:embed parse_contracts.py
var parseContractsScript string

func runSlitherOnLocal(targetDir string, targetExcludePaths []string, testDir string, testExcludePaths []string, slitherArgs map[string]any) (*types.SlitherOutput, error) {
	// Create a temporary file to hold the Python script
	scriptFile, err := os.CreateTemp("", "script-*.py")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("Error removing temporary file:", err)
		}
	}(scriptFile.Name()) // Clean up

	// Write the Python script to the temporary file
	if _, err := scriptFile.Write([]byte(parseContractsScript)); err != nil {
		return nil, fmt.Errorf("error writing to temporary file: %v", err)
	}
	if err := scriptFile.Close(); err != nil {
		return nil, fmt.Errorf("error closing temporary file: %v", err)
	}

	// Create a temporary file to hold the slither output
	outputFile, err := os.CreateTemp("", "slither-output-*.json")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("Error removing temporary file:", err)
		}
	}(outputFile.Name()) // Clean up

	args := []string{"--target", ".", "--out", outputFile.Name(),
		"--contracts-dir", targetDir}

	if len(slitherArgs) > 0 {
		args = append(args, "--slither-args", utils.MapToDictString(slitherArgs))
	}

	if len(targetExcludePaths) > 0 {
		args = append(args, "--exclude-contract-paths",
			strings.Join(targetExcludePaths, ","))
	}

	if testDir != "" {
		args = append(args, "--tests-dir", testDir)
		if len(testExcludePaths) > 0 {
			args = append(args, "--exclude-test-paths", strings.Join(testExcludePaths, ","))
		}
	}

	// Prepare the command
	cmd := exec.Command("python3", append([]string{scriptFile.Name()}, args...)...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("stderr: %s\n", output)
		return nil, fmt.Errorf("error running slither: %v\n", err)
	}

	// Read file contents
	fileContents, err := os.ReadFile(outputFile.Name())
	if err != nil {
		return nil, fmt.Errorf("error reading slither output file: %v", err)
	}

	// Print out slither output
	fmt.Println(string(output))

	var slitherOutput types.SlitherOutput
	err = json.Unmarshal(fileContents, &slitherOutput)
	if err != nil {
		return nil, err
	}

	return &slitherOutput, nil
}

func runSlitherOnchain(address string, networkPrefix string, apiKey string, slitherArgs map[string]any) (*types.SlitherOutput, error) {
	// Create a temporary file to hold the Python script
	scriptFile, err := os.CreateTemp("", "script-*.py")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("Error removing temporary file:", err)
		}
	}(scriptFile.Name()) // Clean up

	// Write the Python script to the temporary file
	if _, err := scriptFile.Write([]byte(parseContractsScript)); err != nil {
		return nil, fmt.Errorf("error writing to temporary file: %v", err)
	}
	if err := scriptFile.Close(); err != nil {
		return nil, fmt.Errorf("error closing temporary file: %v", err)
	}

	// Create a temporary file to hold the slither output
	outputFile, err := os.CreateTemp("", "slither-output-*.json")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("Error removing temporary file:", err)
		}
	}(outputFile.Name()) // Clean up

	args := []string{"--target", address, "--out", outputFile.Name(), "--onchain",
		"--network-prefix", networkPrefix, "--api-key", apiKey}

	if len(slitherArgs) > 0 {
		args = append(args, "--slither-args", utils.MapToDictString(slitherArgs))
	}

	// Prepare the command
	cmd := exec.Command("python3", append([]string{scriptFile.Name()}, args...)...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("stderr: %s\n", output)
		return nil, fmt.Errorf("error running slither: %v\n", err)
	}

	// Print out output
	fmt.Println(string(output))

	// Read file contents
	fileContents, err := os.ReadFile(outputFile.Name())
	if err != nil {
		return nil, fmt.Errorf("error reading slither output file: %v", err)
	}

	var slitherOutput types.SlitherOutput
	err = json.Unmarshal(fileContents, &slitherOutput)
	if err != nil {
		return nil, err
	}

	return &slitherOutput, nil
}

func ParseContracts(projectConfig *config.ProjectConfig) ([]types.Contract, string, error) {
	var slitherOutput *types.SlitherOutput
	var err error

	if projectConfig.OnChainConfig.Enabled {
		slitherOutput, err = runSlitherOnchain(projectConfig.OnChainConfig.Address, projectConfig.OnChainConfig.NetworkPrefix, projectConfig.OnChainConfig.ApiKey, projectConfig.SlitherArgs)
		if err != nil {
			return nil, "", err
		}
	} else {
		// Check if provided directory is a directory
		info, err := os.Stat(projectConfig.TargetContracts.Dir)
		if err != nil || !info.IsDir() {
			return nil, "", fmt.Errorf("unable to read directory")
		}

		slitherOutput, err = runSlitherOnLocal(projectConfig.TargetContracts.Dir, projectConfig.TargetContracts.ExcludePaths, projectConfig.TestContracts.Dir, projectConfig.TestContracts.ExcludePaths, projectConfig.SlitherArgs)
		if err != nil {
			return nil, "", err
		}
	}

	var filteredContracts []types.Contract
	if projectConfig.OnChainConfig.Enabled {
		filteredContracts = filterSlitherOutput(slitherOutput.Contracts, projectConfig.ContractWhitelist, !projectConfig.OnChainConfig.ExcludeInterfaces, true, true)
	} else {
		filteredContracts = filterSlitherOutput(slitherOutput.Contracts, projectConfig.ContractWhitelist, projectConfig.IncludeInterfaces, projectConfig.IncludeAbstract, projectConfig.IncludeLibraries)
	}

	contractCodes := getContractCodes(slitherOutput.Contracts)

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

func filterSlitherOutput(slitherContracts []types.SlitherContract, whitelist []string, includeInterfaces bool, includeAbstract bool, includeLibraries bool) []types.Contract {
	var filteredContracts []types.Contract
	var whitelistMap map[string]bool

	if len(whitelist) > 0 {
		whitelistMap = make(map[string]bool)
		for _, s := range whitelist {
			whitelistMap[s] = true
		}
	}

	for _, slitherContract := range slitherContracts {
		// Skip contracts not in whitelist
		if len(whitelist) > 0 && !whitelistMap[slitherContract.Name] {
			continue
		}

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
			InheritedContracts: filterSlitherOutput(slitherContract.InheritedContracts, []string{}, includeInterfaces, includeAbstract, includeLibraries),
			IsAbstract:         slitherContract.IsAbstract,
			IsInterface:        slitherContract.IsInterface,
			IsLibrary:          slitherContract.IsLibrary,
		})
	}

	return filteredContracts
}
