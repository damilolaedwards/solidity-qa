package cmd

import (
	"assistant/config"
	"github.com/spf13/cobra"
)

// addGenerateFlags adds the various flags for the generate command
func addGenerateFlags() error {
	// Prevent alphabetical sorting of usage message
	generateCmd.Flags().SortFlags = false

	// Config file
	generateCmd.Flags().String("config", "", "path to config file")

	// Flags
	generateCmd.Flags().String("target-contracts-dir", "", "Target contracts directory/file")

	return nil
}

// updateProjectConfigWithGenerateFlags will update the given projectConfig with any CLI arguments that were provided to the generate command
func updateProjectConfigWithGenerateFlags(cmd *cobra.Command, projectConfig *config.ProjectConfig) error {
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
