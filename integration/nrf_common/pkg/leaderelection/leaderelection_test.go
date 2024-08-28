package leaderelection

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func echoHostName(w http.ResponseWriter, r *http.Request) {
	hostName, _ := os.Hostname()
	body := fmt.Sprintf(`{"name":"%s"}`, hostName)
	fmt.Printf("echoHostName")
	fmt.Fprintf(w, body)
}

func echoOther(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf(`{"name":"%s"}`, "nrf-management-other")
	fmt.Printf("echoOther")
	fmt.Fprintf(w, body)
}

func setupLeader(leaderName string) *httpserver.HttpServer {

	h := httpserver.InitHTTPServer(
		httpserver.HostPort("localhost", "4040"),
		httpserver.ReadTimeout(constvalue.HTTP_SERVER_READ_TIMEOUT),
		httpserver.WriteTimeout(constvalue.HTTP_SERVER_WRITE_TIMEOUT),
	)

	hostName, _ := os.Hostname()
	if leaderName == hostName {
		httpserver.PathFunc("/", "GET", echoHostName)(h)
	} else {
		httpserver.PathFunc("/", "GET", echoOther)(h)
	}

	return h

}

func TestIsLeader(t *testing.T) {

	hostName, _ := os.Hostname()
	h := setupLeader(hostName)
	h.Run()
	time.Sleep(time.Second * 1)

	//default_http := httpclient.InitHttpClient()

	if IsLeader() != true {
		t.Fatalf("Should be leader, but Not")
	}
	h.Stop()

	h = setupLeader("otherName")
	go h.Run()
	time.Sleep(time.Second * 1)
	if IsLeader() != false {
		t.Fatalf("Should NOT be leader!")
	}
	h.Stop()

}
