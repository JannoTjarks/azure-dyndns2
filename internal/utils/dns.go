package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
)

type azureDynDnsRecord struct {
	Name string
	MyIP string
	TTL  int64
}

type azureDnsZone struct {
	Name          string
	ResourceGroup string
	Subscription  string
}

func newAzureDynDnsRecord(name string, ip string) azureDynDnsRecord {
	record := azureDynDnsRecord{}
	record.Name = name
	record.MyIP = ip
	record.TTL = 3600
	return record
}

func newAzureDnsZone(dnsZone string, resourceGroup string, subscription string) azureDnsZone {
	record := azureDnsZone{}
	record.Name = dnsZone
	record.ResourceGroup = resourceGroup
	record.Subscription = subscription
	return record
}

func CreateOrUpdateDynDnsRecord(hostname string, myip string, dnsZoneName string, resourceGroupName string, subscriptionId string) {
	ctx := context.Background()

	dnsZone := newAzureDnsZone(
		dnsZoneName,
		resourceGroupName,
		subscriptionId,
	)

	dynDnsRecord := newAzureDynDnsRecord(
		strings.Trim(hostname, dnsZone.Name),
		myip,
	)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("authentication failure: %v", err)
	}

	armdnsClient, err := armdns.NewClientFactory(dnsZone.Subscription, cred, nil)
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

	recordsetClient := armdnsClient.NewRecordSetsClient()

	res, err := recordsetClient.CreateOrUpdate(
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
}
