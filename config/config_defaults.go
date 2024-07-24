package config

// GetDefaultProjectConfig obtains a default configuration for a project. It populates a default compilation config
// based on the provided platform, or a nil one if an empty string is provided.
func GetDefaultProjectConfig() (*ProjectConfig, error) {
	// Create a project configuration
	projectConfig := &ProjectConfig{
		TargetContracts: DirectoryConfig{
			Dir:          "",
			ExcludePaths: []string{},
		},
	}

	// Return the project configuration
	return projectConfig, nil
}
