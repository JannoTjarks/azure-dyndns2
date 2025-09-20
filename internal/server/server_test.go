package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
