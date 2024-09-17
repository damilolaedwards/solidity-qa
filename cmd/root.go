package cmd

import (
	"assistant/logging"
	"os"

	"github.com/rs/zerolog"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "assistant",
	Short:         "This is an AI-powered smart contract auditing assistant",
	Long:          `This application is a CLI tool that leverages on Slither and AI to assist smart contract auditors.`,
	SilenceUsage:  true,
	SilenceErrors: false,
}

// cmdLogger is the logger that will be used for the cmd package
var cmdLogger = logging.NewLogger(zerolog.InfoLevel)

// Execute provides an exportable function to invoke the CLI. Returns an error if one was encountered.
func Execute() error {
	// Add stdout as an unstructured, colorized output stream for the command logger
	cmdLogger.AddWriter(os.Stdout, logging.UNSTRUCTURED, true)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	return rootCmd.Execute()
}
