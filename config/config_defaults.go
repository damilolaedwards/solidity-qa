package config

// GetDefaultProjectConfig obtains a default configuration for a project. It populates a default compilation config
// based on the provided platform, or a nil one if an empty string is provided.
func GetDefaultProjectConfig() (*ProjectConfig, error) {
	// Create a project configuration
	projectConfig := &ProjectConfig{
		Name: "example",
		TargetContracts: DirectoryConfig{
			Dir:          "",
			ExcludePaths: []string{},
		},
		TestContracts: DirectoryConfig{
			Dir:          "",
			ExcludePaths: []string{},
		},
		OnChainConfig: OnChainConfig{
			Enabled:       false,
			Address:       "",
			NetworkPrefix: "mainnet",
			ApiKey:        "",
		},
		Port:              8080,
		IncludeAbstract:   false,
		IncludeInterfaces: false,
		IncludeLibraries:  false,
	}

	// Return the project configuration
	return projectConfig, nil
}
