package cmd

import "github.com/spf13/cobra"

var integrationsCmd = &cobra.Command{
	Use:   "integrations",
	Short: "interact with rocket.chat integrations",
}

func init() {
	newIntegrationCmd.PersistentFlags().String("name", "", "name for new integration")
	newIntegrationCmd.PersistentFlags().String("user", "", "username used by integration")

	integrationsCmd.AddCommand(newIntegrationCmd)
	newIntegrationCmd.AddCommand(newIncomingCmd)
}

var newIntegrationCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new integration",
}

var newIncomingCmd = &cobra.Command{
	Use:   "incoming",
	Short: "create an inbound integration",
	RunE:  newIncoming,
}

func newIncoming(cmd *cobra.Command, args []string) error {
	return nil
}
