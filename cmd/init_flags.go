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

	// Update fuzz tests directory
	if cmd.Flags().Changed("fuzz-tests-dir") {
		projectConfig.FuzzTests.Dir, err = cmd.Flags().GetString("fuzz-tests-dir")
		if err != nil {
			return err
		}
	}

	// Update unit tests directory
	if cmd.Flags().Changed("unit-tests-dir") {
		projectConfig.UnitTests.Dir, err = cmd.Flags().GetString("unit-tests-dir")
		if err != nil {
			return err
		}
	}

	// Update coverage report file
	if cmd.Flags().Changed("coverage-report-file") {
		projectConfig.CoverageReportFile, err = cmd.Flags().GetString("coverage-report-file")
		if err != nil {
			return err
		}
	}
	return nil
}
