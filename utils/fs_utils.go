package utils

import (
	"fmt"
	"os"
	"path/filepath"
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
