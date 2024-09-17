package cmd

import (
	"assistant/config"

	"github.com/spf13/cobra"
)

// addInitFlags adds the various flags for the init command
func addInitFlags() error {
	// Output path for configuration
	initCmd.Flags().String("out", "", "output path for the new project configuration file")
	initCmd.Flags().String("name", "", "name of the project")
	initCmd.Flags().String("port", "8080", "port for the API")
	initCmd.Flags().String("target", "", "directory containing the contracts to be fuzzed")
	initCmd.Flags().String("test-dir", "", "directory containing the test contracts")

	return nil
}

// updateProjectConfigWithInitFlags will update the given projectConfig with any CLI arguments that were provided to the init command
func updateProjectConfigWithInitFlags(cmd *cobra.Command, projectConfig *config.ProjectConfig) error {
	var err error

	// Update name
	if cmd.Flags().Changed("name") {
		projectConfig.Name, err = cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
	}

	// Update port
	if cmd.Flags().Changed("port") {
		projectConfig.Port, err = cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
	}

	// Update target contracts directory
	if cmd.Flags().Changed("target") {
		projectConfig.TargetContracts.Dir, err = cmd.Flags().GetString("target")
		if err != nil {
			return err
		}
	}

	// Update test contracts directory
	if cmd.Flags().Changed("test-dir") {
		projectConfig.TestContracts.Dir, err = cmd.Flags().GetString("test-dir")
		if err != nil {
			return err
		}
	}

	return nil
}
