package utils

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

func NewAzureDynDnsRecord(name string, ip string) azureDynDnsRecord {
	record := azureDynDnsRecord{}
	record.Name = name
	record.MyIP = ip
	record.TTL = 3600
	return record
}

func NewAzureDnsZone(dnsZone string, resourceGroup string, subscription string) azureDnsZone {
	record := azureDnsZone{}
	record.Name = dnsZone
	record.ResourceGroup = resourceGroup
	record.Subscription = subscription
	return record
}
