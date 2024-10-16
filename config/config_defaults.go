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
		ContractWhitelist: []string{},
		TestContracts: DirectoryConfig{
			Dir:          "",
			ExcludePaths: []string{},
		},
		OnChainConfig: OnChainConfig{
			Enabled:           false,
			Address:           "",
			NetworkPrefix:     "mainet",
			ApiKey:            "",
			ExcludeInterfaces: false,
		},
		Port:              8080,
		IncludeAbstract:   false,
		IncludeInterfaces: false,
		IncludeLibraries:  false,
		SlitherArgs:       map[string]any{},
	}

	// Return the project configuration
	return projectConfig, nil
}
