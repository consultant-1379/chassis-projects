package leaderelection

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/buger/jsonparser"
	"gerrit.ericsson.se/udm/nrf_common/pkg/multisite"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
)

// IsLeader is used to judge whether caller is leader
func IsLeader() bool {
	hostName, _ := os.Hostname()
	resp, err := http.Get("http://localhost:4040")
	if err != nil {
		return false
	}

	body, _ := ioutil.ReadAll(resp.Body)
	name, _ := jsonparser.GetString(body, "name")
	defer func() { _ = resp.Body.Close() }()
	if name == hostName {
		if configmap.MultisiteEnabled {
			return multisite.GetMonitor().IsLeaderSite()
		}
		return true
	}
	return false
}
