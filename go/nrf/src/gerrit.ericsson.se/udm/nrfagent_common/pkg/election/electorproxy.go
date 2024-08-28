package election

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"

	"github.com/buger/jsonparser"
)

var (
	leaderElectorURL = "http://localhost:4040"
)

// GetLeader is used to leader get the hostname of leader
var GetLeader = func() string {
	resp, err := http.Get(leaderElectorURL)
	if err != nil {
		return ""
	}
	body, _ := ioutil.ReadAll(resp.Body)
	name, _ := jsonparser.GetString(body, "name")
	defer func() { _ = resp.Body.Close() }()

	return name
}

// IsLeader is used to judge whether caller is leader
var IsLeader = func(identity string) bool {
	selfID := identity
	leaderID := GetLeader()

	if selfID == "" ||
		leaderID == "" {
		return false
	}
	if selfID != leaderID {
		return false
	}
	return true
}

// IsLeaderAlive is used to check whether leaderID is alive or not
func IsLeaderAlive(leaderID, probePort, probeURL string) bool {
	if leaderID == "" {
		log.Errorf("IsLeaderAlive: failed to get leader form leader-elector")
		return false
	}
	if leaderID == os.Getenv("POD_IP") {
		return true
	}

	leaderDiscReadinessProbe := "http://" + leaderID + ":" + probePort + probeURL
	resp, err := http.Get(leaderDiscReadinessProbe)
	if err != nil ||
		resp.StatusCode < http.StatusOK ||
		resp.StatusCode > http.StatusBadRequest {
		log.Infof("IsLeaderAlive: leader(%s) is not ready, waiting for a few seconds", leaderID)
		return false
	}
	return true
}

// IsActiveLeader is used to check whether leaderID is alive or not
var IsActiveLeader = func(probePort, probeURL string) bool {
	selfID := os.Getenv("POD_IP")
	leaderID := GetLeader()
	alive := IsLeaderAlive(leaderID, probePort, probeURL)
	if !alive {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		// Via leader election one of these three pods is selected as the new leader,
		// and you should see the leader failover to a different pod.
		// Because pods in Kubernetes have a grace period before termination,
		// this may take 30-40 seconds.
		timer := time.NewTimer(40 * time.Second)
		defer timer.Stop()

		for !alive {
			select {
			case <-ticker.C:
				leaderID = GetLeader()
				if alive = IsLeaderAlive(leaderID, probePort, probeURL); !alive {
					continue
				}
			case <-timer.C:
				log.Errorf("IsActiveLeader: failed to get leader form leader-elector, try to act as leader")
				return true
			}
			log.Debugf("IsActiveLeader: leader(%s) is ready", leaderID)
		}
	}
	if selfID != leaderID {
		return false
	}
	return true
}
