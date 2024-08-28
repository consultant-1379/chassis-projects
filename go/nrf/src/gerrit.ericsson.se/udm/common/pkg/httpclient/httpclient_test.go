package httpclient

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
	"gerrit.ericsson.se/udm/common/pkg/utils"
)

const hostPort string = "40012"

func healthCheck_Handler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Println("server ", req.URL.String())
}

func TestGetHttp1WithoutTLS(t *testing.T) {
	default_http := InitHttpClient(
		Connections(5),
		Trace(true),
		Timeout(7*time.Second),
		KeepAlive(true),
	)

	h := httpserver.InitHTTPServer(
		httpserver.Trace(true),
		httpserver.HostPort("", hostPort),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),

		httpserver.PathFunc("/health-check", "GET,PUT,DELETE,POST", healthCheck_Handler),
	)

	h.Run()
	time.Sleep(1 * time.Second)

	resp, err := default_http.HttpDoJsonBody("GET", "http://127.0.0.1:" + hostPort + "/health-check", nil)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	fmt.Printf("resp:  %v, err %v\n", resp.SimpleString(), err)
	fmt.Printf("resp:  %v, err %v\n", resp.String(), err)

	if resp.StatusCode != 200 || resp.Protocol != "HTTP/1.1" {
		t.Fatalf("resp:  %v, err %v", resp.SimpleString(), err)
	}
	h.Stop()
	time.Sleep(1 * time.Second)
}

func TestGetHttp2WithTLS(t *testing.T) {
	default_https := InitHttpClient(
		Connections(5),
		Trace(true),
		Timeout(7*time.Second),
		KeepAlive(true),
		HTTP2(true),
	)

	tlsconfig, err := utils.GenTlsConfig(false, "server.crt", "server.key", nil)
	if err != nil {
		t.Fatalf("Can not get tls config")
	}
	h := httpserver.InitHTTPServer(
		httpserver.Trace(true),
		httpserver.HostPort("", hostPort),
		httpserver.HTTP2(true),
		httpserver.TLSConfig(tlsconfig),
		httpserver.MaxConcurrentStreams(100),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),

		httpserver.PathFunc("/health-check", "GET,PUT,DELETE,POST", healthCheck_Handler),
	)

	h.Run()
	time.Sleep(1 * time.Second)

	resp, err := default_https.HttpDoJsonBody("GET", "https://127.0.0.1:" + hostPort + "/health-check", nil)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if resp.StatusCode != 200 || resp.Protocol != "HTTP/2.0" {
		t.Fatalf("resp:  %v, err %v", resp.SimpleString(), err)
	}
	h.Stop()
	time.Sleep(1 * time.Second)
}

func TestGetHttp1WithTLS(t *testing.T) {
	default_https := InitHttpClient(
		Connections(5),
		Trace(true),
		Timeout(7*time.Second),
		KeepAlive(true),
		HTTP2(true),
	)

	tlsconfig, err := utils.GenTlsConfig(false, "server.crt", "server.key", nil)
	if err != nil {
		t.Fatalf("Can not get tls config")
	}
	h := httpserver.InitHTTPServer(
		httpserver.Trace(true),
		httpserver.HostPort("", hostPort),
		//httpserver.HTTP2(true),
		httpserver.TLSConfig(tlsconfig),
		httpserver.MaxConcurrentStreams(100),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),

		httpserver.PathFunc("/health-check", "GET,PUT,DELETE,POST", healthCheck_Handler),
	)

	h.Run()
	time.Sleep(1 * time.Second)

	resp, err := default_https.HttpDoJsonBody("GET", "https://127.0.0.1:40012/health-check", nil)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if resp.StatusCode != 200 || resp.Protocol != "HTTP/1.1" {
		t.Fatalf("resp:  %v, err %v", resp.SimpleString(), err)
	}
	h.Stop()
	time.Sleep(1 * time.Second)
}
