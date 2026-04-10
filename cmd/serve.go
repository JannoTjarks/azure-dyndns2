package cmd

import (
	"github.com/JannoTjarks/azure-dyndns2/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultPort = "8080"

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve dyndns2 api",
	Long:  "Starts a webserver that follows the dyndns2 standard to create DNS records in a Azure DNS Zone",
	Run: func(cmd *cobra.Command, args []string) {
		server.Serve(
			viper.GetString(port.name),
			viper.GetString(dnsZoneName.name),
			viper.GetString(dnsResourceGroupName.name),
			viper.GetString(dnsSubscriptionId.name),
		)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP(port.name, port.shorthand, port.value, port.usage)
	serveCmd.Flags().String(dnsZoneName.name, dnsResourceGroupName.value, dnsResourceGroupName.usage)
	serveCmd.Flags().String(dnsResourceGroupName.name, dnsResourceGroupName.value, dnsResourceGroupName.usage)
	serveCmd.Flags().String(dnsSubscriptionId.name, dnsSubscriptionId.value, dnsSubscriptionId.usage)

	serveCmd.PreRun = func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())
		viper.SetDefault(port.name, defaultPort)
		viper.AutomaticEnv()

		viper.BindEnv(port.name, port.envVar)
		viper.BindEnv(dnsZoneName.name, dnsZoneName.envVar)
		viper.BindEnv(dnsResourceGroupName.name, dnsResourceGroupName.envVar)
		viper.BindEnv(dnsSubscriptionId.name, dnsSubscriptionId.envVar)
	}
}
