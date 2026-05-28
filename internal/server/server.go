package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/JannoTjarks/azure-dyndns2/docs"
	"github.com/JannoTjarks/azure-dyndns2/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger/v2"
)

var OpenApiFiles embed.FS

var (
	config ServerConfig

	ipUpdatesProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "azure-dyndns2_processed_ip_updates_total",
		Help: "The total number of processed ip updates via DynDNS Update API since start of the service.",
	})
	ipUpdatesProcessedSucceeded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "azure-dyndns2_processed_ip_updates_succeeded",
		Help: "The total number of succeeded ip updates via DynDNS Update API since start of the service.",
	})
)

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

func ipUpdateHandler(w http.ResponseWriter, req *http.Request) {
	ipUpdatesProcessedTotal.Inc()
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
	ipUpdatesProcessedSucceeded.Inc()
	fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
}

func configHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/config" {
		http.NotFound(w, req)
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusNotFound))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonBytes, _ := json.Marshal(config)

	fmt.Fprintf(w, "%s\n", string(jsonBytes))
	fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
}

func versionHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/version" {
		http.NotFound(w, req)
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusNotFound))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "%s\n", utils.GenerateVersionJson())
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

func openapiHandler(w http.ResponseWriter, req *http.Request) {
	switch file := req.PathValue("file"); file {
	case "openapi.json":
		fileContent, err := OpenApiFiles.ReadFile("docs/openapi.json")
		if err != nil {
			http.NotFound(w, req)
			fmt.Println(formatCommonLog(*req, time.Now(), http.StatusNotFound))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", string(fileContent))
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
	case "openapi.yaml":
		fileContent, err := OpenApiFiles.ReadFile("docs/openapi.yaml")
		if err != nil {
			http.NotFound(w, req)
			fmt.Println(formatCommonLog(*req, time.Now(), http.StatusNotFound))
		}

		w.Header().Set("Content-Type", "application/yaml")
		fmt.Fprintf(w, "%s\n", string(fileContent))
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
	case "docs":
		fileContent, err := OpenApiFiles.ReadFile("docs/index.html")
		if err != nil {
			http.NotFound(w, req)
			fmt.Println(formatCommonLog(*req, time.Now(), http.StatusNotFound))
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "%s\n", string(fileContent))
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusOK))
	default:
		http.NotFound(w, req)
		fmt.Println(formatCommonLog(*req, time.Now(), http.StatusNotFound))
	}
}

func Serve(port string, dnsZoneName string, resourceGroupName string, subscriptionId string) {
	config = newServerConfig(port, dnsZoneName, resourceGroupName, subscriptionId)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /nic/update", ipUpdateHandler)
	mux.HandleFunc("GET /config", configHandler)
	mux.HandleFunc("GET /version", versionHandler)
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	mux.HandleFunc("GET /api/{file}", openapiHandler)
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("/", fallbackHandler)

	mux.Handle("GET /metrics", promhttp.Handler())

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
