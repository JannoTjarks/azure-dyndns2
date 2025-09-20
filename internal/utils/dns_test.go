package utils

import "testing"

func TestNewAzureDynDnsRecord(t *testing.T) {
	name := "hello.world"
	myip := "127.0.0.2"
	want := azureDynDnsRecord{Name: name, MyIP: myip, TTL: 3600}

	record := newAzureDynDnsRecord(name, myip)

	if record != want {
		t.Errorf(`utils.newAzureDynDnsRecord() = %q, want match for %#q, nil`, record, want)
	}
}

func TestNewAzureDnsZone(t *testing.T) {
	zone := "world"
	rg := "dns-world-rg"
	sub := "e09dfa34-17d9-4bba-8583-fa63b45a9b2a"
	want := azureDnsZone{Name: zone, ResourceGroup: rg, Subscription: sub}

	record := newAzureDnsZone(zone, rg, sub)

	if record != want {
		t.Errorf(`utils.newAzureDnsZone() = %q, want match for %#q, nil`, record, want)
	}
}

func TestExtractPQDN(t *testing.T) {
	hostname := "en.wikipedia.org"
	dnsZoneName := "wikipedia.org"
	want := "en"

	pqdn, err := extractPQDN(hostname, dnsZoneName)
	if pqdn != want {
		t.Errorf(`utils.extractPQDN() = %q, want match for %#q, nil`, pqdn, want)
	}
	if err != nil {
		t.Errorf(`utils.extractPQDN() = invoke of function throws an error`)
	}
}

func TestExtractPQDNFaulty(t *testing.T) {
	hostname := "en.wikipedia.org"
	dnsZoneName := "google.de"

	_, err := extractPQDN(hostname, dnsZoneName)
	if err == nil {
		t.Errorf(`utils.extractPQDN() = Error should be shrown, because extracting of PQDN was not successful`)
	}
}
