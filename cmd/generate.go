package cmd

import (
	"assistant/api"
	"assistant/config"
	"assistant/logging/colors"
	"assistant/slither"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var generateCmd = &cobra.Command{
	Use:           "generate",
	Short:         "Generate invariants for Medusa",
	Long:          `Generate invariants for Medusa`,
	Args:          cmdValidateGenerateArgs,
	RunE:          cmdRunGenerate,
	SilenceUsage:  true,
	SilenceErrors: false,
}

func init() {
	// Add all the flags allowed for the generate command
	err := addGenerateFlags()
	if err != nil {
		cmdLogger.Panic("Failed to initialize the fuzz command", err)
	}

	// Add the generate command and its associated flags to the root command
	rootCmd.AddCommand(generateCmd)
}

// cmdValidateGenerateArgs makes sure that there are no positional arguments provided to the generate command
func cmdValidateGenerateArgs(cmd *cobra.Command, args []string) error {
	// Make sure we have no positional args
	if err := cobra.NoArgs(cmd, args); err != nil {
		err = fmt.Errorf("generate does not accept any positional arguments, only flags and their associated values")
		cmdLogger.Error("Failed to validate args to the generate command", err)
		return err
	}
	return nil
}

// cmdRunGenerate runs the generate CLI command
func cmdRunGenerate(cmd *cobra.Command, args []string) error {
	var projectConfig *config.ProjectConfig

	// Check to see if --config flag was used and store the value of --config flag
	configFlagUsed := cmd.Flags().Changed("config")
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		cmdLogger.Error("Failed to run the generate command", err)
		return err
	}

	// If --config was not used, look for `assistant.json` in the current work directory
	if !configFlagUsed {
		workingDirectory, err := os.Getwd()
		if err != nil {
			cmdLogger.Error("Failed to run the generate command", err)
			return err
		}
		configPath = filepath.Join(workingDirectory, DefaultProjectConfigFilename)
	}

	// Check to see if the file exists at configPath
	_, existenceError := os.Stat(configPath)

	// Possibility #1: File was found
	if existenceError == nil {
		// Try to read the configuration file and throw an error if something goes wrong
		cmdLogger.Info("Reading the configuration file at: ", colors.Bold, configPath, colors.Reset)
		projectConfig, err = config.ReadProjectConfigFromFile(configPath)
		if err != nil {
			cmdLogger.Error("Failed to run the generate command", err)
			return err
		}
	}

	// Possibility #2: If the --config flag was used, and we couldn't find the file, we'll throw an error
	if configFlagUsed && existenceError != nil {
		cmdLogger.Error("Failed to run the generate command", err)
		return existenceError
	}

	// Possibility #3: --config flag was not used and assistant.json was not found, so use the default project config
	if !configFlagUsed && existenceError != nil {
		cmdLogger.Warn(fmt.Sprintf("Unable to find the config file at %v, will use the default project configuration", configPath))

		projectConfig, err = config.GetDefaultProjectConfig()
		if err != nil {
			cmdLogger.Error("Failed to run the generate command", err)
			return err
		}
	}

	// Update the project configuration given whatever flags were set using the CLI
	err = updateProjectConfigWithGenerateFlags(cmd, projectConfig)
	if err != nil {
		cmdLogger.Error("Failed to run the generate command", err)
		return err
	}

	// Validate project config
	err = projectConfig.Validate()
	if err != nil {
		cmdLogger.Error("Failed to run the generate command", err)
		return err
	}

	// Run slither
	cmdLogger.Info("Running Slither on the target contracts directory: ", colors.Green, projectConfig.TargetContracts.Dir, colors.Reset, ", Excluding paths: ", colors.Red, projectConfig.TargetContracts.ExcludePaths, colors.Reset, "\n")
	contracts, contractCodes, err := slither.ParseContracts(projectConfig)
	if err != nil {
		cmdLogger.Error("Failed to run the generate command", err)
		return err
	}
	cmdLogger.Info("Successfully ran Slither on the target contracts directory")

	// Start the API
	api.InitializeAPI(contractCodes, contracts).Start(projectConfig)

	return nil
}
