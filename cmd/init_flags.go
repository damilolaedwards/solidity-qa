package cmd

import (
	"assistant/config"
	"github.com/spf13/cobra"
)

// addInitFlags adds the various flags for the init command
func addInitFlags() error {
	// Output path for configuration
	initCmd.Flags().String("out", "", "output path for the new project configuration file")
	initCmd.Flags().String("compilations-dir", "", "compilations directory path")
	initCmd.Flags().String("coverage-report-file", "", "coverage report file")
	initCmd.Flags().String("unit-tests-dir", "", "compilations directory path")

	return nil
}

// updateProjectConfigWithInitFlags will update the given projectConfig with any CLI arguments that were provided to the init command
func updateProjectConfigWithInitFlags(cmd *cobra.Command, projectConfig *config.ProjectConfig) error {
	var err error

	// Update target contracts directory
	if cmd.Flags().Changed("target-contracts-dir") {
		projectConfig.TargetContracts.Dir, err = cmd.Flags().GetString("target-contracts-dir")
		if err != nil {
			return err
		}
	}

	return nil
}
