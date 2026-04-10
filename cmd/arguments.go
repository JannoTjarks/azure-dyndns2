package cmd

type argument struct {
	name      string
	shorthand string
	value     string
	usage     string
	envVar    string
}

var (
	port = argument{
		name:      "port",
		shorthand: "p",
		usage:     "The port which will be used by the webapi",
		envVar:    "AZURE_DYNDNS_PORT",
	}
	hostname = argument{
		name:      "hostname",
		shorthand: "n",
		usage:     "Hostname which will be updated - Must be a Fully Qualified Domain Name (fqdn)",
	}
	myip = argument{
		name:      "myip",
		shorthand: "i",
		usage:     "IP Address which will be set in Azure DNS - Must be a IPv4 Address",
	}
	dnsZoneName = argument{
		name:   "dns-zone-name",
		usage:  "The name of the Azure DNS zone",
		envVar: "AZURE_DYNDNS_DNS_ZONE_NAME",
	}
	dnsResourceGroupName = argument{
		name:   "dns-resource-group-name",
		usage:  "The name of the Resource Group which contains the Azure DNS zone",
		envVar: "AZURE_DYNDNS_DNS_RESOURCE_GROUP_NAME",
	}
	dnsSubscriptionId = argument{
		name:   "dns-subscription-id",
		usage:  "The Subscription Id which contains the Azure DNS zone",
		envVar: "AZURE_DYNDNS_DNS_SUBSCRIPTION_ID",
	}
)
