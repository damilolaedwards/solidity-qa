package cmd

import (
	"assistant/config"
	"assistant/llm"
	"assistant/logging/colors"
	"assistant/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
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

	targetContracts, err := utils.ReadDirectoryContents(projectConfig.TargetContracts.Dir, projectConfig.TargetContracts.ExcludePaths...)
	if err != nil {
		cmdLogger.Error("Failed to run the generate command", err)
		return err
	}

	var fuzzTests string
	if projectConfig.FuzzTests.Dir != "" {
		fuzzTests, err = utils.ReadDirectoryContents(projectConfig.FuzzTests.Dir, projectConfig.FuzzTests.ExcludePaths...)
		if err != nil {
			cmdLogger.Error("Failed to run the generate command", err)
			return err
		}
	}

	var unitTests string
	if projectConfig.UnitTests.Dir != "" {
		unitTests, err = utils.ReadDirectoryContents(projectConfig.UnitTests.Dir, projectConfig.UnitTests.ExcludePaths...)
		if err != nil {
			cmdLogger.Error("Failed to run the generate command", err)
			return err
		}
	}

	var parsedCoverageReport utils.CoverageReport
	if projectConfig.CoverageReportFile != "" {
		_, err := os.Stat(projectConfig.CoverageReportFile)
		if err != nil {
			cmdLogger.Error("Failed to run the generate command", err)
			return err
		}

		coverageReport, err := os.ReadFile(projectConfig.CoverageReportFile)
		if err != nil {
			cmdLogger.Error("Failed to run the generate command", err)
			return err
		}

		parsedCoverageReport, err = utils.ParseCoverageReportHTML(string(coverageReport))
		if err != nil {
			cmdLogger.Error("Failed to run the generate command", err)
			return err
		}

		// Exclude reports of files we do not need
		includePaths := []string{projectConfig.TargetContracts.Dir}
		excludePaths := projectConfig.TargetContracts.ExcludePaths

		if projectConfig.FuzzTests.Dir != "" {
			includePaths = append(includePaths, projectConfig.FuzzTests.Dir)
			excludePaths = append(excludePaths, projectConfig.FuzzTests.ExcludePaths...)
		}

		utils.FilterCoverageFiles(&parsedCoverageReport, includePaths, excludePaths)
	}

	// Construct the messages for the LLM
	messages := append(llm.TrainingPrompts(), llm.Message{
		Role:    "user",
		Content: llm.GenerateInvariantsPrompt(targetContracts, fuzzTests, unitTests, fmt.Sprintf("%v", parsedCoverageReport)),
	})

	// Make request to LLM to generate invariants
	invariants, err := llm.AskGPT4Turbo(messages)
	if err != nil {
		cmdLogger.Error("Failed to run the generate command", err)
		return err
	}

	invariants = fmt.Sprintf("=============== Invariants generated at %v ===============\n\n%v\n\n", time.Now().String(), invariants)

	// Open the out file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(projectConfig.Out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			cmdLogger.Error("Error closing "+projectConfig.Out, err)
		}
	}(file)

	// Write to the file
	_, err = file.WriteString(invariants)
	if err != nil {
		return err
	}

	return nil
}
