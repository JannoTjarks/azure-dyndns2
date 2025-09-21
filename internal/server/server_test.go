package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAzureDynDnsRecord(t *testing.T) {
	port := "8080"
	dnsZoneName := "wikipedia.de"
	resourceGroupName := "publicdns-rg"
	subscriptionId := "054a27fd-2cd9-42fa-9139-c99cf680fd35"

	want := serverConfig{port, dnsZoneName, resourceGroupName, subscriptionId}

	record := newServerConfig(port, dnsZoneName, resourceGroupName, subscriptionId)

	if record != want {
		t.Errorf(`utils.serverConfig() = %q, want match for %#q, nil`, record, want)
	}
}

func TestVersionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(versionHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	want := fmt.Sprintf("%s\n", `{"version":"dev","commit":"none","date":"unknown","buildby":"unknown"}`)
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), want)
	}
}

func TestVersionHandlerWithExtraPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/version/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(versionHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	want := fmt.Sprintf("%s\n", `404 page not found`)
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), want)
	}
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	want := fmt.Sprintf("%s\n", "Thanks for using azure-dyndns2!")
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), want)
	}
}

func TestRootHandlerWrongPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestRootHandlerWierdPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/../", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestFallbackHandlerDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fallbackHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	want := fmt.Sprintf("%s\n", "HTTP Method is not allowed!")
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), want)
	}
}

func TestFallbackHandlerPost(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fallbackHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	want := fmt.Sprintf("%s\n", "HTTP Method is not allowed!")
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), want)
	}
}
