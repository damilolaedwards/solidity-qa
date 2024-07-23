package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ReadDirectoryContents(dirPath string, excludePaths ...string) (string, error) {
	var result strings.Builder

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path should be excluded
		for _, excludePath := range excludePaths {
			if strings.HasPrefix(path, excludePath) || strings.HasPrefix(path, filepath.Join(dirPath, excludePath)) {
				return filepath.SkipDir
			}
		}

		// If it's a file, read and append its contents
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			result.WriteString(fmt.Sprintf("--- Start of file: %s ---\n", path))
			result.Write(content)
			result.WriteString(fmt.Sprintf("\n--- End of file: %s ---\n\n", path))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return result.String(), nil
}

// GenerateRegexFromPaths generates a regex pattern to match the given paths and their subdirectories
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
