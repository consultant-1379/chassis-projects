package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthzHandler(t *testing.T) {
	srv := http.NewServeMux()
	srv.HandleFunc("/healthz", HealthzHandler)
	ts := httptest.NewServer(srv)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(err)
		t.FailNow()
	}
}

func TestRunProbe(t *testing.T) {
	go runProbe("0")
	time.Sleep(200 * time.Millisecond)
}
