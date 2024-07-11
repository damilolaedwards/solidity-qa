package cmd

import (
	"assistant/config"
	"assistant/logging/colors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// initCmd represents the command provider for init
var initCmd = &cobra.Command{
	Use:           "init",
	Short:         "Initializes a project configuration",
	Long:          `Initializes a project configuration`,
	RunE:          cmdRunInit,
	SilenceUsage:  true,
	SilenceErrors: false,
}

func init() {
	// Add flags to init command
	err := addInitFlags()
	if err != nil {
		cmdLogger.Panic("Failed to initialize the init command", err)
	}

	// Add the init command and its associated flags to the root command
	rootCmd.AddCommand(initCmd)
}

// cmdRunInit executes the init CLI command and updates the project configuration with any flags
func cmdRunInit(cmd *cobra.Command, args []string) error {
	// Check to see if --out flag was used and store the value of --out flag
	outputFlagUsed := cmd.Flags().Changed("out")
	outputPath, err := cmd.Flags().GetString("out")
	if err != nil {
		cmdLogger.Error("Failed to run the init command", err)
		return err
	}

	// If we weren't provided an output path (flag was not used), we use our working directory
	if !outputFlagUsed {
		workingDirectory, err := os.Getwd()
		if err != nil {
			cmdLogger.Error("Failed to run the init command", err)
			return err
		}
		outputPath = filepath.Join(workingDirectory, DefaultProjectConfigFilename)
	}

	// By default, projectConfig will be the default project config for the DefaultCompilationPlatform
	projectConfig, err := config.GetDefaultProjectConfig()
	if err != nil {
		cmdLogger.Error("Failed to run the init command", err)
		return err
	}

	// Update the project configuration given whatever flags were set using the CLI
	err = updateProjectConfigWithInitFlags(cmd, projectConfig)
	if err != nil {
		cmdLogger.Error("Failed to run the init command", err)
		return err
	}

	// Write our project configuration
	err = projectConfig.WriteToFile(outputPath)
	if err != nil {
		cmdLogger.Error("Failed to run the init command", err)
		return err
	}

	// Print a success message
	if absoluteOutputPath, err := filepath.Abs(outputPath); err == nil {
		outputPath = absoluteOutputPath
	}
	cmdLogger.Info("Project configuration successfully output to: ", colors.Bold, outputPath, colors.Reset)
	return nil
}
