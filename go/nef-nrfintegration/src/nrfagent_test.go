package main

import (
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddSuffix(t *testing.T) {
	uri := "http://localhost"
	uri = addSuffix(uri)
	if uri != "http://localhost/" {
		t.Error()
	}
	uri = "http://localhost/"
	uri = addSuffix(uri)
	if uri != "http://localhost/" {
		t.Error()
	}
}

func TestSendHB(t *testing.T) {
	testSrvName := "nnef-bdt"
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Error(r.Method)
			t.FailNow()
		}
		if r.URL.EscapedPath() != "/nrf-register-agent/v1/nf-status/NEF/"+testSrvName {
			t.Error(r.URL.EscapedPath())
			t.FailNow()
		}
		w.WriteHeader(http.StatusOK)
	})
	ts := httptest.NewServer(h2c.NewHandler(srv, &http2.Server{}))
	defer ts.Close()

	// Create client for test
	config = &Config{}
	config.nrfAgentUri = ts.URL
	config.heartbeatInterval = 10
	config.nfProfilePath = "ericsson-nef:nef"
	config.nfCMProfileName = "ericsson-nef"
	if err := sendHeartbeat(testSrvName); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestSendHB_Neg(t *testing.T) {
	testSrvName := "nnef-bdt"
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Error(r.Method)
			t.FailNow()
		}
		if r.URL.EscapedPath() != "/nrf-register-agent/v1/nf-status/NEF/"+testSrvName {
			t.Error(r.URL.EscapedPath())
			t.FailNow()
		}
		w.WriteHeader(http.StatusNotImplemented)
	})
	ts := httptest.NewServer(h2c.NewHandler(srv, &http2.Server{}))
	defer ts.Close()

	// Create client for test
	config = &Config{}
	config.nrfAgentUri = ts.URL
	config.heartbeatInterval = 10
	config.nfProfilePath = "ericsson-nef:nef"
	config.nfCMProfileName = "ericsson-nef"
	if err := sendHeartbeat(testSrvName); err == nil {
		t.FailNow()
	} else if err.Error() != "NRFAgent responded with StatusCode = 501 for HB of the service( nnef-bdt )" {
		t.Error(err.Error())
		t.FailNow()
	}
}
