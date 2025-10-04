package cmd

import (
	"github.com/JannoTjarks/azure-dyndns2/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve dyndns2 api",
	Long:  "Starts a webserver that follows the dyndns2 standard to create DNS records in a Azure DNS Zone",
	Run: func(cmd *cobra.Command, args []string) {
		server.Serve(
			viper.GetString("port"),
			viper.GetString("dns-zone-name"),
			viper.GetString("dns-resource-group-name"),
			viper.GetString("dns-subscription-id"),
		)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("port", "p", "", "The port which will be used by the webapi")
	serveCmd.Flags().String("dns-zone-name", "", "The name of the Azure DNS zone")
	serveCmd.Flags().String("dns-resource-group-name", "", "The name of the Resource Group which contains the Azure DNS zone")
	serveCmd.Flags().String("dns-subscription-id", "", "The Subscription Id which contains the Azure DNS zone")

	viper.BindPFlags(serveCmd.Flags())
	viper.SetDefault("port", "8080")
	viper.AutomaticEnv()

	viper.BindEnv("port", "AZURE_DYNDNS_PORT")
	viper.BindEnv("dns-zone-name", "AZURE_DYNDNS_DNS_ZONE_NAME")
	viper.BindEnv("dns-resource-group-name", "AZURE_DYNDNS_DNS_RESOURCE_GROUP_NAME")
	viper.BindEnv("dns-subscription-id", "AZURE_DYNDNS_DNS_SUBSCRIPTION_ID")
}
