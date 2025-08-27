package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/JannoTjarks/azure-dyndns2/internal/utils"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
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
		ctx := context.Background()

		dnsZone := utils.NewAzureDnsZone(
			dnsZone,
			resourceGroup,
			subscription,
		)

		dynDnsRecord := utils.NewAzureDynDnsRecord(
			strings.Trim(hostname, dnsZone.Name),
			myip,
		)

		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("authentication failure: %v", err)
		}

		clientFactory, err := armdns.NewClientFactory(dnsZone.Subscription, cred, nil)
		if err != nil {
			log.Fatalf("failed to create armdns client: %v", err)
		}

		recordSet := armdns.RecordSet{
			Properties: &armdns.RecordSetProperties{
				ARecords: []*armdns.ARecord{
					{
						IPv4Address: to.Ptr(dynDnsRecord.MyIP),
					}},
				TTL: to.Ptr(dynDnsRecord.TTL),
			},
		}

		recordSetCreateUpdateOptions := armdns.RecordSetsClientCreateOrUpdateOptions{
			IfMatch:     nil,
			IfNoneMatch: nil,
		}

		client := clientFactory.NewRecordSetsClient()

		res, err := client.CreateOrUpdate(
			ctx,
			dnsZone.ResourceGroup,
			dnsZone.Name,
			dynDnsRecord.Name,
			armdns.RecordTypeA,
			recordSet,
			&recordSetCreateUpdateOptions,
		)

		if err != nil {
			log.Fatalf("failed to finish the request: %v", err)
		}

		fmt.Println(*res.Properties.ProvisioningState)
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
