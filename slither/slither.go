package slither

import (
	"assistant/types"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed slither.py
var pythonScript string

func RunSlitherOnTarget(target string, outputFile *os.File) error {
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
	if _, err := tmpfile.Write([]byte(pythonScript)); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return nil
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Println("Error closing temporary file:", err)
		return nil
	}

	// Arguments to pass to the Python script
	args := []string{target, outputFile.Name()}

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

func RunSlitherOnDir(dirPath string, excludePaths ...string) (*types.SlitherOutput, error) {
	var contracts []types.Contract

	// Create a temporary file to hold the slither output
	tmpfile, err := os.CreateTemp("", "slither-output-*.json")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return nil, nil
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println("Error removing temporary file:", err)
		}
	}(tmpfile.Name()) // Clean up

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path should be excluded
		for _, excludePath := range excludePaths {
			if strings.HasPrefix(path, excludePath) || strings.HasPrefix(path, filepath.Join(dirPath, excludePath)) {
				return filepath.SkipDir
			}
		}

		// If it's a file, run slither on this path
		if !info.IsDir() {
			err = RunSlitherOnTarget(path, tmpfile)
			if err != nil {
				return err
			}

			// Read the contents of the temporary file
			file, err := os.ReadFile(tmpfile.Name())
			if err != nil {
				return err
			}

			// Parse the slither output
			var contractsInFile []types.Contract
			err = json.Unmarshal(file, &contractsInFile)
			if err != nil {
				return err
			}

			// Append the contracts to the result
			contracts = append(contracts, contractsInFile...)

			// Clear the temporary file
			err = os.Truncate(tmpfile.Name(), 0)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.SlitherOutput{
		Contracts: contracts,
	}, nil
}

func GenerateRegexFromPaths(paths []string) string {
	var escapedPaths []string
	for _, path := range paths {
		// Escape special regex characters in the path
		escapedPath := regexp.QuoteMeta(path)
		// Add a pattern to match the path and its subdirectories
		escapedPaths = append(escapedPaths, fmt.Sprintf("%s(/.*)?", escapedPath))
	}
	// Join all the individual patterns with | to create the final regex
	return fmt.Sprintf("^(%s)$", strings.Join(escapedPaths, "|"))
}
