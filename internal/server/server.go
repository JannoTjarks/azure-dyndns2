package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/JannoTjarks/azure-dyndns2/internal/utils"
)

var config serverConfig

type serverConfig struct {
	port              string
	dnsZoneName       string
	resourceGroupName string
	subscriptionId    string
}

func newServerConfig(port string, dnsZoneName string, resourceGroupName string, subscriptionId string) serverConfig {
	config := serverConfig{}
	config.port = port
	config.dnsZoneName = dnsZoneName
	config.resourceGroupName = resourceGroupName
	config.subscriptionId = subscriptionId
	return config
}

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
	err := utils.CreateOrUpdateDynDnsRecord(hostname, myip, config.dnsZoneName, config.resourceGroupName, config.subscriptionId)
	if err != nil {
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusInternalServerError))
		return
	}
	fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
}

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
