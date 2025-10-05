# azure-dyndns2
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/JannoTjarks/azure-dyndns2.svg)](https://github.com/JannoTjarks/azure-dyndns2)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=JannoTjarks_azure-dyndns2&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=JannoTjarks_azure-dyndns2)
[![GitHub license](https://img.shields.io/github/license/JannoTjarks/azure-dyndns2.svg)](https://github.com/JannoTjarks/azure-dyndns2/blob/master/LICENSE)
[![Latest release](https://img.shields.io/github/v/release/JannoTjarks/azure-dyndns2)](https://github.com/JannoTjarks/azure-dyndns2/releases)


DynDNS (Dynamic DNS) is a service that allows the DNS record of a router or server to be automatically updated when it changes, ensuring stable reachability.

Currently (Fall 2025) Azure has no built-in solution for DynDNS. This project goal is to enable Azure DNS to support DynDNS by adding a small webapi, which can be addressed by the DynDNS Update API, also called the dyndns2 standard.
That’s why this project is called azure-dyndns2.

Using a web API instead of a command line tool allows for greater flexibility and compatibility with various clients, such as simple ISP-managed routers, which typically have only few basic configuration options.

## Usage
### One quick example
```bash
./azure-dyndns2 one-shot --hostname <fqdn> --myip <ip-address> --dns-zone-name <azure-zone-name> --dns-resource-group-name <azure-resource-group-name> --dns-subscription-id <azure-subscription-id>
```

### All you need to know
`azure-dyndns2` has two modes:
| Mode | Description | Status |
| --- | --- | --- |
| one-shot | Sets a A Record in Azure DNS one-time | stable |
| serve | Starts a webserver which accepts http requests that are following the dyndns2 standard. | In testing |

As a base configuration the AzureDNS zone, the Resource Group and the Subscription MUST be specified.  
When the application is used in serve mode, the network port CAN be specified. Port number 8080 is set as the default value.

All configurations can be done via Commandline flags or with Environment variables:
| Configuration | Commandline flag | Environment variable | Default Value |
| --- | --- | --- | --- |
| DNS zone name | dns-zone-name | AZURE_DYNDNS_DNS_ZONE_NAME | |
| Resource Group Name | dns-resource-group-name | AZURE_DYNDNS_DNS_RESOURCE_GROUP_NAME | |
| Subscription ID | dns-subscription-id | AZURE_DYNDNS_DNS_SUBSCRIPTION_ID | |
| Port | port | AZURE_DYNDNS_PORT | 8080 |

Authentication and authorization on your Azure tenant must be performed in order for a DNS record update to be carried out. The Microsoft modules used support the DefaultAzureCredential class for this purpose. 

Currently two methods are tested:
- [AzureCliCredential](https://learn.microsoft.com/en-us/dotnet/api/azure.identity.azureclicredential?view=azure-dotnet)
- [EnvironmentCredential](https://learn.microsoft.com/en-us/dotnet/api/azure.identity.environmentcredential?view=azure-dotnet)

The AzureCliCredential should only be used inside of interactive sessions. For non-interactive invokes the EnvironmentCredential method is recommended. In most cases, a principal with secret is used for this purpose. The following environment variables must be set:
| Variable | Description |
| --- | --- |
| AZURE_TENANT_ID | The Microsoft Entra tenant (directory) ID |
| AZURE_CLIENT_ID | The client (application) ID of an App Registration in the tenant |
| AZURE_CLIENT_SECRET | A client secret that was generated for the App Registration |

### Usage of the binary
#### Run the one-shot mode
```bash
./azure-dyndns2 one-shot --hostname <fqdn> --myip <ip-address> --dns-zone-name <azure-zone-name> --dns-resource-group-name <azure-resource-group-name> --dns-subscription-id <azure-subscription-id>
```

#### Run the serve mode
```bash
./azure-dyndns2 serve --dns-zone-name <azure-zone-name> --dns-resource-group-name <azure-resource-group-name> --dns-subscription-id <azure-subscription-id>
```

### Usage of the Container
A container image is available on the GitHub Container Registry.

The container starts automatically the `serve mode`. If you want to use the one-shot command, you have to specify the command. It is recommended to do all configuration with Environment variables. The examples would work with an .env file in the following schema:
```env
AZURE_DYNDNS_DNS_ZONE_NAME=xyz
AZURE_DYNDNS_DNS_RESOURCE_GROUP_NAME=xyz
AZURE_DYNDNS_DNS_SUBSCRIPTION_ID=xyz

AZURE_TENANT_ID=xyz
AZURE_CLIENT_ID=xyz
AZURE_CLIENT_SECRET=xyz
```
#### Run the one-shot mode
```bash
podman run --env-file .env ghcr.io/jannotjarks/azure-dyndns2:latest /app/azure-dyndns2 one-shot --hostname <fqdn> --myip <ip-address>
```

#### Run the serve mode
```bash
podman run --env-file .env -p 8080:8080 ghcr.io/jannotjarks/azure-dyndns2:latest
```

## The DynDNS Update API
The DynDNS Update API allows the update of an ip address with an WebAPI/REST call. This call is described here: [DynDNS Perform Update](https://help.dyn.com/perform-update.html).

Based on this standard the following minimal schema would allow the update of a dns record:
```bash
# URL Syntax
https://{user}:{password}@{dyndns-service}/nic/update?hostname={hostname}&myip={IP Address}

# Raw HTTP GET Request
GET /nic/update?hostname={hostname}&myip={ipaddress} HTTP/1.0
Host: {dyndns-service}
Authorization: Basic {base64-decoded-authorization}
User-Agent: {client-user-agent}
```
The dyndns update request should always be done with a HTTP GET request; DynDNS seems still to accept HTTP POST requests, but dyn.com mentions, that this behavior can change without notice!

Because of that, this tool will accept right from the beginning only HTTP GET requests on the path /nic/update.

The webserver returns meaningful HTTP Status Codes and also follow the Return Codes of the dyndns2 standard

## Background
When a forced disconnection or Zwangstrennung occurs, resulting in an IP address change, it can disrupt the reachability and configuration of servers or applications that depend on a specific IP address. DynDNS can help maintain access to your network even with dynamic IPs.

#### Zwangstrennung: A Forced Disconnection

A Zwangstrennung, or forced disconnection, is a technical measure implemented by internet service providers (ISPs) to disconnect a user's internet connection and reassign their IP address. This process is typically automated and occurs at regular intervals, such as every 24 hours. In Germany and Austria, forced disconnections are particularly prevalent.

#### Impact on Users

In Germany, for example, many ISPs implement a 24-hour forced disconnection policy, where customers' IP addresses are renewed every 24 hours. This practice is common in Austria as well.

A Zwangstrennung can have implications for users, particularly those who require a stable and continuous internet connection for applications such as online gaming, video streaming, or remote work.
However, for most users, the impact is often minimal, as modern routers and ISP configurations typically schedule Zwangstrennungen during nighttime hours when users are less likely to be actively using their internet connections. This minimizes disruptions and allows users to maintain a seamless online experience during peak usage hours.

#### Positive aspects

Forced disconnections, or Zwangstrennungen, can have a positive effect on users by enhancing their online security and anonymity through regular IP address changes, making it more difficult for hackers and trackers to identify and target their devices.

#### Technical Reasons

The technical reasons for implementing a Zwangstrennung include:

Dynamic IP address management: ISPs use dynamic IP address assignment to manage their IP address pools and ensure efficient use of available addresses. By forcing a disconnection and reassigning IP addresses, ISPs can free up addresses and assign them to other users.
DHCP lease renewal: ISPs use the Dynamic Host Configuration Protocol (DHCP) to assign IP addresses to users. A Zwangstrennung can be triggered by the expiration of a DHCP lease, which requires the user to renew their IP address and potentially receive a new one.
