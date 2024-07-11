package config

// GetDefaultProjectConfig obtains a default configuration for a project. It populates a default compilation config
// based on the provided platform, or a nil one if an empty string is provided.
func GetDefaultProjectConfig() (*ProjectConfig, error) {
	// Create a project configuration
	projectConfig := &ProjectConfig{
		Out: "invariants.txt",
		TargetContracts: DirectoryConfig{
			Dir:          "",
			ExcludePaths: []string{},
		},
		FuzzTests: DirectoryConfig{
			Dir:          "",
			ExcludePaths: []string{},
		},
		UnitTests: DirectoryConfig{
			Dir:          "",
			ExcludePaths: []string{},
		},
		CoverageReportFile: "corpus/coverage_report.html",
	}

	// Return the project configuration
	return projectConfig, nil
}
