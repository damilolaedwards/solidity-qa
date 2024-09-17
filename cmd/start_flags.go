package cmd

import (
	"assistant/config"

	"github.com/spf13/cobra"
)

// addStartFlags adds the various flags for the start command
func addStartFlags() error {
	// Prevent alphabetical sorting of usage message
	startCmd.Flags().SortFlags = false

	// Flags
	startCmd.Flags().String("config", "", "path to config file")
	startCmd.Flags().String("name", "", "project name")
	startCmd.Flags().Bool("onchain", false, "used to specify that an onchain contract will be used rather than a local project")
	startCmd.Flags().String("address", "", "address of contract to be analyzed")
	startCmd.Flags().String("network-prefix", "", "network prefix of contract to be analyzed")
	startCmd.Flags().String("api-key", "", "network API key")
	startCmd.Flags().Bool("exclude-interfaces", false, "used to specify that interfaces should be excluded from the analysis")

	return nil
}

// updateProjectConfigWithStartFlags will update the given projectConfig with any CLI arguments that were provided to the start command
func updateProjectConfigWithStartFlags(cmd *cobra.Command, projectConfig *config.ProjectConfig) error {
	var err error

	if cmd.Flags().Changed("name") {
		projectConfig.Name, err = cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
	}

	// Update etherscan config if onchain is set to true
	if cmd.Flags().Changed("onchain") {
		projectConfig.OnChainConfig.Enabled, err = cmd.Flags().GetBool("onchain")
		if err != nil {
			return err
		}

		projectConfig.OnChainConfig.NetworkPrefix, err = cmd.Flags().GetString("network-prefix")
		if err != nil {
			return err
		}

		projectConfig.OnChainConfig.Address, err = cmd.Flags().GetString("address")
		if err != nil {
			return err
		}

		projectConfig.OnChainConfig.ApiKey, err = cmd.Flags().GetString("api-key")
		if err != nil {
			return err
		}

		projectConfig.OnChainConfig.ExcludeInterfaces, err = cmd.Flags().GetBool("exclude-interfaces")
		if err != nil {
			return err
		}
	}

	return nil
}
