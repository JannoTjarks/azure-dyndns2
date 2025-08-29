package cmd

import (
	"github.com/JannoTjarks/azure-dyndns2/internal/utils"

	"github.com/spf13/cobra"
)

var (
	hostname      string
	myip          string
	dnsZone       string
	resourceGroup string
	subscription  string
)

var oneshotCmd = &cobra.Command{
	Use:     "one-shot",
	Aliases: []string{"oneshot"},
	Short:   "Set a DNS Record direct from the CLI",
	Long:    "Allows you to set a DNS Record direct from the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CreateOrUpdateDynDnsRecord(hostname, myip, dnsZone, resourceGroup, subscription)
	},
}

func init() {
	rootCmd.AddCommand(oneshotCmd)

	oneshotCmd.Flags().StringVarP(&hostname, "hostname", "n", "", "Hostname which will be updated - Must be a Fully Qualified Domain Name (fqdn)")
	oneshotCmd.Flags().StringVarP(&myip, "myip", "i", "", "IP Adress which will be set in Azure DNS - Must be a IPv4 Address")

	oneshotCmd.Flags().StringVarP(&dnsZone, "dns-zone", "z", "", "The name of the Azure DNS zone")
	oneshotCmd.Flags().StringVarP(&resourceGroup, "dns-resource-group", "r", "", "The name of the Resource Group which contains the Azure DNS zone")
	oneshotCmd.Flags().StringVarP(&subscription, "dns-subscription", "s", "", "The Subscription Id which contains the Azure DNS zone")
}
