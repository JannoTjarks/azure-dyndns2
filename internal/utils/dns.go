package utils

import (
	"context"
	"errors"
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

func extractPQDN(hostname string, dnsZoneName string) (string, error) {
	pqdn, found := strings.CutSuffix(hostname, "."+dnsZoneName)
	if !found {
		return "", errors.New("it is not possible to extract the pqdn (partially-qualified domain name) from the hostname! the hostname must be a fqdn")
	}

	return pqdn, nil
}

func CreateOrUpdateDynDnsRecord(hostname string, myip string, dnsZoneName string, resourceGroupName string, subscriptionId string) (string, error) {
	ctx := context.Background()

	dnsZone := newAzureDnsZone(
		dnsZoneName,
		resourceGroupName,
		subscriptionId,
	)

	pqdn, err := extractPQDN(hostname, dnsZoneName)
	if err != nil {
		fmt.Printf("failed to create the request: %v\n", err)
		return "", err
	}

	dynDnsRecord := newAzureDynDnsRecord(
		pqdn,
		myip,
	)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Printf("authentication failure: %v\n", err)
		return "", err
	}

	armdnsClient, err := armdns.NewClientFactory(dnsZone.Subscription, cred, nil)
	if err != nil {
		fmt.Printf("failed to create armdns client: %v\n", err)
		return "", err
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
		fmt.Printf("failed to finish the request: %v\n", err)
		return "", err
	}

	provState := fmt.Sprintf("%s: DNS CreateOrUpdate %s\n", time.Now().Format("2006-01-02T15:04:05Z07:00"), *res.Properties.ProvisioningState)
	return provState, nil
}
