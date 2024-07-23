package cmd

import (
	"assistant/config"
	"fmt"

	"github.com/spf13/cobra"
)

// addGenerateFlags adds the various flags for the generate command
func addGenerateFlags() error {
	// Get the default project config and throw an error if we cant
	defaultConfig, err := config.GetDefaultProjectConfig()
	if err != nil {
		return err
	}

	// Prevent alphabetical sorting of usage message
	generateCmd.Flags().SortFlags = false

	// Config file
	generateCmd.Flags().String("config", "", "path to config file")

	// Flags
	generateCmd.Flags().String("out", "",
		fmt.Sprintf("path to output directory (unless a config file is provided, default is %q)", defaultConfig.Out))

	return nil
}

// updateProjectConfigWithGenerateFlags will update the given projectConfig with any CLI arguments that were provided to the generate command
func updateProjectConfigWithGenerateFlags(cmd *cobra.Command, projectConfig *config.ProjectConfig) error {
	var err error

	// Update output path
	if cmd.Flags().Changed("out") {
		projectConfig.Out, err = cmd.Flags().GetString("out")
		if err != nil {
			return err
		}
	}

	// Update target contracts directory
	if cmd.Flags().Changed("target-contracts-dir") {
		projectConfig.TargetContracts.Dir, err = cmd.Flags().GetString("target-contracts-dir")
		if err != nil {
			return err
		}
	}

	return nil
}
