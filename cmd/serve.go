package cmd

import (
	"github.com/JannoTjarks/azure-dyndns2/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve dyndns2 api",
	Long:  "Starts a webserver that follows the dyndns2 standard to create DNS records in a Azure DNS Zone",
	Run: func(cmd *cobra.Command, args []string) {
		server.Serve(port, dnsZoneName, resourceGroupName, subscriptionId)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "The port which will be used by the webapi")
	serveCmd.Flags().StringVarP(&dnsZoneName, "dns-zone-name", "z", "", "The name of the Azure DNS zone")
	serveCmd.Flags().StringVarP(&resourceGroupName, "dns-resource-group-name", "r", "", "The name of the Resource Group which contains the Azure DNS zone")
	serveCmd.Flags().StringVarP(&subscriptionId, "dns-subscription-id", "s", "", "The Subscription Id which contains the Azure DNS zone")
}
