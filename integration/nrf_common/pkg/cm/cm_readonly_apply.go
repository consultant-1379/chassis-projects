package cm

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/clusterleader"
	"gerrit.ericsson.se/udm/common/pkg/cmproxy"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

const (
	leaderRetryTime   = 3
	followerRetryTime = 10
)

var (
	configName string
)

// TPatchItem is to construct PATCH data to CM
type TPatchItem struct {
	Value interface{} `json:"value,omitempty"`
	From  string      `json:"from,omitempty"`
	Op    string      `json:"op"`
	Path  string      `json:"path"`
}

// ApplyReadOnlyData apply read-only data to CM after NRF initial configuration created
func ApplyReadOnlyData() {
	configName = os.Getenv("CMM_CONFIG_SCHEMA_NAME")
	if configName == "" {
		log.Errorf("environment variable CMM_CONFIG_SCHEMA_NAME is empty, quit ApplyReadOnlyData")
		return
	}

	leaderTryTime := 0
	followerTryTime := 0
	for {
		if !isApplyNeed() {
			break
		}
		if clusterleader.IsLeader() {
			leaderTryTime++
			if leaderTryTime > (leaderRetryTime + 1) {
				log.Warningf("reach to the retry time %d, stop retry", leaderRetryTime)
				break
			}

			log.Infof("to apply read-only data to CM")
			if err := applyReadOnlyData(); err != nil {
				log.Warningf("failed to apply read-only data to CM, %s. Try again ...", err.Error())
			} else {
				log.Infof("apply read-only data to CM successfully")
				log.Infof("sleep 3 seconds to wait for read-only data to take effect")
				time.Sleep(time.Second * 3)
				break
			}
		} else {
			followerTryTime++
			if followerTryTime > (followerRetryTime + 1) {
				break
			}
		}

		time.Sleep(time.Second * 1)
	}
}

func isApplyNeed() bool {
	nfProfile := GetNRFNFProfile()
	if nfProfile != nil && len(nfProfile.Service) > 0 {
		for _, service := range nfProfile.Service {
			if service.Version == nil || len(service.Version) == 0 {
				return true
			}
		}
	}

	return false
}

func applyReadOnlyData() error {
	readOnlyData := constructReadonlyDataPATCH()
	if readOnlyData == nil {
		return fmt.Errorf("construct PATCH data error")
	}

	err := sendReadonlyDataPatchToCM(readOnlyData)
	if err != nil {
		return fmt.Errorf("send PATCH data to CM falied")
	}

	return nil
}

func constructReadonlyDataPATCH() []byte {
	patchData := make([]TPatchItem, 0)

	// nf-profile.type
	patchData = append(patchData, TPatchItem{
		Op:    "add",
		Path:  "/ericsson-nrf:nrf/nf-profile/type",
		Value: "nrf",
	})

	// nf-profile.service-persistence
	patchData = append(patchData, TPatchItem{
		Op:    "add",
		Path:  "/ericsson-nrf:nrf/nf-profile/service-persistence",
		Value: false,
	})

	// nf-profile.service[index1].version[index2].api-version-in-uri
	nfProfile := GetNRFNFProfile()
	if nfProfile != nil {
		for serviceIndex, service := range nfProfile.Service {
			if service.Version == nil {
				patchData = append(patchData, TPatchItem{
					Op:   "add",
					Path: fmt.Sprintf("/ericsson-nrf:nrf/nf-profile/service/%d/version", serviceIndex),
					Value: []TVersion{
						TVersion{
							APIFullVersion:  "1.R15.1.1",
							APIVersionInURI: "v1",
							Expiry:          "2020-07-06T02:54:32Z",
						},
					},
				})
			} else if len(service.Version) == 0 {
				patchData = append(patchData, TPatchItem{
					Op:   "add",
					Path: fmt.Sprintf("/ericsson-nrf:nrf/nf-profile/service/%d/version/0", serviceIndex),
					Value: TVersion{
						APIFullVersion:  "1.R15.1.1",
						APIVersionInURI: "v1",
						Expiry:          "2020-07-06T02:54:32Z",
					},
				})
			}
		}
	}

	data, err := json.Marshal(patchData)
	if err != nil {
		log.Errorf("Marshal error, %s", err.Error())
		return nil
	}

	return data
}

func sendReadonlyDataPatchToCM(body []byte) error {
	_, err := cmproxy.UpdateConfiguration(configName, body)

	return err
}
