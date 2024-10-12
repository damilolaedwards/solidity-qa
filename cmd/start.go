package cmd

import (
	"assistant/api"
	"assistant/config"
	"assistant/internal/slither"
	"assistant/logging/colors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:           "start",
	Short:         "Parses contracts and spins up frontend server on specified port",
	Long:          "Parses contracts using Slither and spins up frontend server on specified port",
	Args:          cmdValidateStartArgs,
	RunE:          cmdRunStart,
	SilenceUsage:  true,
	SilenceErrors: false,
}

func init() {
	// Add all the flags allowed for the start command
	err := addStartFlags()
	if err != nil {
		cmdLogger.Panic("Failed to initialize the fuzz command", err)
	}

	// Add the start command and its associated flags to the root command
	rootCmd.AddCommand(startCmd)
}

// cmdValidateStartArgs makes sure that there are only one positional argument can be provided to the start command
func cmdValidateStartArgs(cmd *cobra.Command, args []string) error {
	// Make sure we have no positional args
	if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
		err = fmt.Errorf("start can only accept one positional argument, target contracts directory")
		cmdLogger.Error("Failed to validate args to the start command", err)
		return err
	}

	return nil
}

// cmdRunStart runs the start CLI command
func cmdRunStart(cmd *cobra.Command, args []string) error {
	var projectConfig *config.ProjectConfig

	// Check to see if --config flag was used and store the value of --config flag
	configFlagUsed := cmd.Flags().Changed("config")
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		cmdLogger.Error("Failed to run the start command", err)
		return err
	}

	// If --config was not used, look for `assistant.json` in the current work directory
	if !configFlagUsed {
		workingDirectory, err := os.Getwd()
		if err != nil {
			cmdLogger.Error("Failed to run the start command", err)
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
			cmdLogger.Error("Failed to run the start command", err)
			return err
		}
	}

	// Possibility #2: If the --config flag was used, and we couldn't find the file, we'll throw an error
	if configFlagUsed && existenceError != nil {
		cmdLogger.Error("Failed to run the start command", err)
		return existenceError
	}

	// Possibility #3: --config flag was not used and assistant.json was not found, so use the default project config
	if !configFlagUsed && existenceError != nil {
		cmdLogger.Warn(fmt.Sprintf("Unable to find the config file at %v, will use the default project configuration", configPath))

		projectConfig, err = config.GetDefaultProjectConfig()
		if err != nil {
			cmdLogger.Error("Failed to run the start command", err)
			return err
		}
	}

	// Update the project configuration given whatever flags were set using the CLI
	err = updateProjectConfigWithStartFlags(cmd, projectConfig)
	if err != nil {
		cmdLogger.Error("Failed to run the start command", err)
		return err
	}

	// Validate project config
	err = projectConfig.Validate()
	if err != nil {
		cmdLogger.Error("Failed to run the start command", err)
		return err
	}

	// Update target if provided as a positional argument
	if len(args) == 1 {
		projectConfig.TargetContracts.Dir = args[0]
	}

	if projectConfig.OnChainConfig.Enabled {
		cmdLogger.Info("Running Slither on contract at address: ", colors.Green, projectConfig.OnChainConfig.Address, colors.Reset)
	} else {
		if len(projectConfig.ContractWhitelist) > 0 {
			cmdLogger.Info("Running Slither on the target contracts directory: ", colors.Green, projectConfig.TargetContracts.Dir, colors.Reset, ", Excluding paths: ", colors.Red, projectConfig.TargetContracts.ExcludePaths, colors.Reset, ", Selecting contracts: ", colors.Green, projectConfig.ContractWhitelist, colors.Reset, "\n")
		} else {
			cmdLogger.Info("Running Slither on the target contracts directory: ", colors.Green, projectConfig.TargetContracts.Dir, colors.Reset, ", Excluding paths: ", colors.Red, projectConfig.TargetContracts.ExcludePaths, colors.Reset, "\n")
		}
	}

	// Parse contracts
	contracts, contractCodes, err := slither.ParseContracts(projectConfig)
	if err != nil {
		cmdLogger.Error("Failed to run the start command", err)
		return err
	}

	cmdLogger.Info("Successfully ran Slither on target")

	// Start the API
	api.InitializeAPI(contractCodes, contracts).Start(projectConfig)

	return nil
}
