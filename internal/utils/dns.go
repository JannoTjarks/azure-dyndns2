package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func CreateOrUpdateDynDnsRecord(hostname string, myip string, dnsZoneName string, resourceGroupName string, subscriptionId string) error {
	ctx := context.Background()

	dnsZone := newAzureDnsZone(
		dnsZoneName,
		resourceGroupName,
		subscriptionId,
	)

	pqdn, _ := strings.CutSuffix(hostname, "."+dnsZone.Name)

	dynDnsRecord := newAzureDynDnsRecord(
		pqdn,
		myip,
	)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Printf("authentication failure: %v", err)
		return err
	}

	armdnsClient, err := armdns.NewClientFactory(dnsZone.Subscription, cred, nil)
	if err != nil {
		fmt.Printf("failed to create armdns client: %v", err)
		return err
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
		fmt.Printf("failed to finish the request: %v", err)
		return err
	}

	fmt.Printf("%s: DNS CreateOrUpdate %s\n", time.Now().Format("2006-01-02T15:04:05Z07:00"), *res.Properties.ProvisioningState)
	return nil
}
