package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/JannoTjarks/azure-dyndns2/docs"
	"github.com/JannoTjarks/azure-dyndns2/internal/utils"
	"github.com/swaggo/http-swagger/v2"
)

var config ServerConfig

type ServerConfig struct {
	Port              string `json:"port"`
	DnsZoneName       string `json:"dnsZoneName"`
	ResourceGroupName string `json:"resourceGroupName"`
	SubscriptionId    string `json:"subscriptionId"`
}

func newServerConfig(port string, dnsZoneName string, resourceGroupName string, subscriptionId string) ServerConfig {
	config := ServerConfig{}
	config.Port = port
	config.DnsZoneName = dnsZoneName
	config.ResourceGroupName = resourceGroupName
	config.SubscriptionId = subscriptionId
	return config
}

// @Summary		Update IP address of a DNS Record
// @Description	Set a DNS A Record by using the DynDNS Update API, also called the dyndns2 standard.
// @Produce		plain
// @Param			hostname	query		string		true	"Hostname which will be updated - Must be a Fully Qualified Domain Name"
// @Param			myip	query		string		false	"IP Address which will be set in Azure DNS - Must be a IPv4 Address"
// @Success		200		{string}	string
// @Failure		400		{string}	string
// @Failure		500		{string}	string
// @Router		/nic/update [get]
func ipUpdateHandler(w http.ResponseWriter, req *http.Request) {
	if !req.URL.Query().Has("hostname") {
		http.Error(w, "You have to specify the http query parameter 'hostname'", http.StatusBadRequest)
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusBadRequest))
		return
	}

	var hostname string
	var myip string

	hostname = req.URL.Query().Get("hostname")
	if req.URL.Query().Has("myip") {
		myip = req.URL.Query().Get("myip")
	} else {
		myip = req.RemoteAddr
	}
	provState, err := utils.CreateOrUpdateDynDnsRecord(hostname, myip, config.DnsZoneName, config.ResourceGroupName, config.SubscriptionId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusInternalServerError))
		return
	}

	fmt.Fprintf(w, "%s", provState)
	fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
}

// @Summary		Returns the current configuration of used azure-dyndns2 instance
// @Description	Returns the current configuration of used azure-dyndns2 instance as json object.
// @Produce		json
// @Success		200		{string}	string
// @Failure		404		{string}	string
// @Router		/config [get]
func configHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/config" {
		http.NotFound(w, req)
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusNotFound))
		return
	}

	jsonBytes, _ := json.Marshal(config)

	fmt.Fprintf(w, "%s\n", string(jsonBytes))
	fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
}

// @Summary		Returns the current version of azure-dyndns2
// @Description	Returns the current version of azure-dyndns2 as json object.
// @Produce		json
// @Success		200		{string}	string
// @Failure		404		{string}	string
// @Router		/version [get]
func versionHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/version" {
		http.NotFound(w, req)
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusNotFound))
		return
	}

	fmt.Fprintf(w, "%s\n", utils.GenerateVersionJson())
	fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
}

// @Summary		Returns a small welcome message
// @Description	Returns a small welcome message. This welcome page can be used as a very simple healthcheck.
// @Produce		plain
// @Success		200		{string}	string
// @Failure		404		{string}	string
// @Router		/ [get]
func rootHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusNotFound))
		return
	}

	fmt.Fprintf(w, "Thanks for using azure-dyndns2!\n")
	fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
}

func fallbackHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET")
	http.Error(w, "HTTP Method is not allowed!", http.StatusMethodNotAllowed)
	fmt.Println(formatCommonLog(*req, time.Now(), http.StatusMethodNotAllowed))
}

func Serve(port string, dnsZoneName string, resourceGroupName string, subscriptionId string) {
	config = newServerConfig(port, dnsZoneName, resourceGroupName, subscriptionId)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /nic/update", ipUpdateHandler)
	mux.HandleFunc("GET /config", configHandler)
	mux.HandleFunc("GET /version", versionHandler)
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("/", fallbackHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on %s\n", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal(err)
	}
}
