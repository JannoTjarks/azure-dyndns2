package cmd

import (
	"fmt"
	"log"

	"github.com/JannoTjarks/azure-dyndns2/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var oneshotCmd = &cobra.Command{
	Use:     "one-shot",
	Aliases: []string{"oneshot"},
	Short:   "Set a DNS Record direct from the CLI",
	Long:    "Allows you to set a DNS Record direct from the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		provState, err := utils.CreateOrUpdateDynDnsRecord(
			viper.GetString("hostname"),
			viper.GetString("myip"),
			viper.GetString("dns-zone-name"),
			viper.GetString("dns-resource-group-name"),
			viper.GetString("dns-subscription-id"),
		)
		if err != nil {
			log.Fatal("An error occured...")
		}

		fmt.Print(provState)
	},
}

func init() {
	rootCmd.AddCommand(oneshotCmd)

	oneshotCmd.Flags().StringP("hostname", "n", "", "Hostname which will be updated - Must be a Fully Qualified Domain Name (fqdn)")
	oneshotCmd.Flags().StringP("myip", "i", "", "IP Address which will be set in Azure DNS - Must be a IPv4 Address")
	oneshotCmd.Flags().String("dns-zone-name", "", "The name of the Azure DNS zone")
	oneshotCmd.Flags().String("dns-resource-group-name", "", "The name of the Resource Group which contains the Azure DNS zone")
	oneshotCmd.Flags().String("dns-subscription-id", "", "The Subscription Id which contains the Azure DNS zone")


	viper.BindPFlags(oneshotCmd.Flags())
	viper.AutomaticEnv()

	viper.BindEnv("dns-zone-name", "AZURE_DYNDNS_DNS_ZONE_NAME")
	viper.BindEnv("dns-resource-group-name", "AZURE_DYNDNS_DNS_RESOURCE_GROUP_NAME")
	viper.BindEnv("dns-subscription-id", "AZURE_DYNDNS_DNS_SUBSCRIPTION_ID")
}
