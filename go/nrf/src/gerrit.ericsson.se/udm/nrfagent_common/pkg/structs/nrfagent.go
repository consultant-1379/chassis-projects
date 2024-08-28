package structs

import (
	"encoding/json"
	"strings"
	"sync"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
)

const (
	//NrfMgmtServiceName is service name of NRF management
	NrfMgmtServiceName = "nnrf-nfm"
	//NrfDiscServiceName is service name of NRF discovery
	NrfDiscServiceName = "nnrf-disc"
)

//TargetNfProfile for NRF Agent
type TargetNfProfile struct {
	RequesterNfType          string          `json:"requesterNfType"`
	TargetNfType             string          `json:"targetNfType"`
	TargetServiceNames       []string        `json:"targetServiceNames"`
	NotifCondition           *NotifCondition `json:"notifCondition,omitempty"`
	SubscriptionValidityTime int             `json:"subscriptionValidityTime"`
	RequesterNfFqdn          string          `json:"requesterNfFqdn,omitempty"`
	CallbackURI              string          `json:"callbackUri,omitempty"`
	TargetPlmn               PlmnID          `json:"targetPlmn,omitempty"`
	SupportedFeatures        string          `json:"supportedFeatures,omitempty"`
	NfSpecificNrfServerList  NrfServerList   `json:"nfSpecificNrfServerList,omitempty"`
}

//NrfAgentConf for NRF Agent
type NrfAgentConf struct {
	DefaultNrfServerList  NrfServerList         `json:"nrf,omitempty"`
	StatusNotifIPEndPoint StatusNotifIPEndPoint `json:"notification-address,omitempty"`
}

//StatusNotifIPEndPoint for NRF Agent
type StatusNotifIPEndPoint struct {
	Ipv4Address string `json:"ipv4-address,omitempty"`
	Ipv6Address string `json:"ipv6-address,omitempty"`
	Transport   string `json:"transport,omitempty"`
	Port        int    `json:"port,omitempty"`
}

//NrfServerList for NRF Agent
type NrfServerList struct {
	NrfServerProfileList []NrfServerProfile `json:"profile,omitempty"`
	Mode                 string             `json:"mode,omitempty"`
}

//NrfServerProfile for NRF Agent
type NrfServerProfile struct {
	NrfServiceEndPoints  []NrfServiceEndPoint `json:"service,omitempty"`
	NrfServerProfileID   string               `json:"id"`
	NrfServerFqdn        string               `json:"fqdn,omitempty"`
	NrfServerIpv4Address []string             `json:"ipv4-address,omitempty"`
	NrfServerIpv6Address []string             `json:"ipv6-address,omitempty"`
	NrfServerPriority    int                  `json:"priority,omitempty"`
	NrfServerCapacity    int                  `json:"capacity,omitempty"`
	NrfServerLocality    string               `json:"locality,omitempty"`
}

//NrfServiceEndPoint for NRF Agent
type NrfServiceEndPoint struct {
	ID                int                    `json:"id"`
	ServiceName       string                 `json:"name"`
	Versions          []NrfServiceVersion    `json:"version"`
	Scheme            string                 `json:"scheme"`
	Fqdn              string                 `json:"fqdn,omitempty"`
	IPEndPoints       []NrfServiceIPEndPoint `json:"ip-endpoint,omitempty"`
	Priority          int                    `json:"priority,omitempty"`
	Capacity          int                    `json:"capacity,omitempty"`
	APIPrefix         string                 `json:"api-prefix,omitempty"`
	SupportedFeatures string                 `json:"supported-features,omitempty"`
}

//NrfServiceVersion definition
type NrfServiceVersion struct {
	APIVersionInUrI string `json:"api-version-in-uri,omitempty"`
	APIFullVersion  string `json:"api-full-version,omitempty"`
	Expiry          string `json:"expiry,omitempty"`
}

//NrfServiceIPEndPoint definition
type NrfServiceIPEndPoint struct {
	ID          int    `json:"id"`
	Ipv4Address string `json:"ipv4-address,omitempty"`
	Ipv6Address string `json:"ipv6-address,omitempty"`
	Transport   string `json:"transport,omitempty"`
	Port        int    `json:"port,omitempty"`
}

/*
 * --------------------------------------------struct list for nf-service-log -------------------------------------------------------------
 */
type CMNfServiceLog struct {
	LogID    string   `json:"log-id,omitempty"`
	Severity string   `json:"severity,omitempty"`
	PodLogs  []PodLog `json:"pod-log"`
}

type PodLog struct {
	PodID    string `json:"pod-id,omitempty"`
	Severity string `json:"severity,omitempty"`
}

var (
	nrfAgentsConfLock       sync.Mutex
	nfProfile               NfProfile
	nrfServerList           NrfServerList
	nrfAgentNotifIPEndpoint *StatusNotifIPEndPoint
	nfServicelogs           []CMNfServiceLog

	targetNfProfiles = make(map[string][]TargetNfProfile)
)

//UpdateNfProfile update configuration nfProfile
func UpdateNfProfile(data []byte) bool {
	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	var newNfProfile NfProfile
	err := json.Unmarshal(data, &newNfProfile)
	if err != nil {
		log.Debugf("failed to Unmarshal NfProfile, %s\n", err.Error())
		return false
	}
	log.Debugf("%v\n", newNfProfile)

	nfProfile = newNfProfile
	return true
}

//GetNfProfile get configuration nfProfile
func GetNfProfile(c *NfProfile) bool {
	if c == nil {
		return false
	}

	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	*c = nfProfile
	return true
}

//UpdateTargetNfProfilesByName update configuration targetNfProfiles
func UpdateTargetNfProfilesByName(configName string, data []byte) bool {
	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	var newTargetNfProfiles []TargetNfProfile
	err := json.Unmarshal(data, &newTargetNfProfiles)
	if err != nil {
		log.Debugf("failed to Unmarshal TargetNfProfile[%s], %s\n", configName, err.Error())
		return false
	}
	log.Debugf("UpdateTargetNfProfiles[%s]: %v\n", configName, newTargetNfProfiles)

	targetNfProfiles[configName] = newTargetNfProfiles
	return true
}

//CheckTargetNfProfilesByName check configuration targetNfProfiles existed or not
func CheckTargetNfProfilesByName(configName string) bool {
	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	_, existed := targetNfProfiles[configName]
	if !existed {
		log.Infof("%s was not loaded by fsnotify", configName)
		delete(targetNfProfiles, configName)
		return false
	}

	return true
}

//GetTargetNfProfilesByName get configuration targetNfProfiles by name
func GetTargetNfProfilesByName(configName string) []TargetNfProfile {
	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	existed := CheckTargetNfProfilesByName(configName)
	if !existed {
		return nil
	}

	// memory leak, this part of memory can be only recycled by GC
	newTargetNfProfiles := []TargetNfProfile{}

	nfProfiles, _ := targetNfProfiles[configName]
	copy(newTargetNfProfiles, nfProfiles)

	return newTargetNfProfiles
}

//GetTargetNfProfiles get configuration targetNfProfiles
func GetTargetNfProfiles() []TargetNfProfile {
	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	var newTargetNfProfiles []TargetNfProfile
	for _, v := range targetNfProfiles {
		newTargetNfProfiles = append(newTargetNfProfiles, v...)
	}
	return newTargetNfProfiles
}

func updatePodLogLevel(serviceLog CMNfServiceLog) {
	var podID string = cm.PodIp
	var usePodLevel bool = false

	for _, podLog := range serviceLog.PodLogs {

		if podLog.PodID == podID {
			if podLog.Severity != "inherit" {
				log.SetLevel(log.LevelUint(strings.ToUpper(podLog.Severity)))
				log.Infof("Change log level to %s success", podLog.Severity)
				usePodLevel = true
				break
			} else {
				break
			}
		}
	}

	if !usePodLevel {
		log.SetLevel(log.LevelUint(strings.ToUpper(serviceLog.Severity)))
		log.Infof("Change log level to %s success", serviceLog.Severity)
	}
}

//UpdateNfServiceLog update nfServiceLog
func UpdateNfServiceLog(data []byte, serviceName string) bool {
	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	var newNfServicelogs []CMNfServiceLog
	err := json.Unmarshal(data, &newNfServicelogs)
	if err != nil {
		log.Debugf("failed to Unmarshal nfServicelogs, %s\n", err.Error())
		return false
	}

	var serviceLog CMNfServiceLog
	found := false
	for _, serviceLog = range newNfServicelogs {
		if serviceLog.LogID == serviceName {
			found = true

			updatePodLogLevel(serviceLog)
			//log.SetLevel(log.Level(cm.Opts.LogLevel))
			//log.SetLevel(log.Level(serviceLog.Severity))
			//log.SetLevel(log.Level(serviceLog.Severity))
			//log.SetOutput(os.Stdout)

			//log.SetServiceID(serviceName)
			//log.SetNF("nrf-agent")
			//log.SetPodIP(newNfServiceLog.PodLogs.)
			//log.SetFormatter(&log.JSONFormatter{})
			break
		}
	}

	return found
}

//GetNfServiceLog get configuration nfServiceLog
func GetNfServiceLog(c *CMNfServiceLog, serviceName string) bool {
	if c == nil {
		return false
	}

	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	found := false
	for _, ns := range nfServicelogs {
		if ns.LogID == serviceName {
			*c = ns
			found = true
			break
		}
	}
	return found
}

//UpdateNrfServerList update data in nrf-agent configuration
func UpdateNrfServerList(data []byte) bool {
	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	var newNrfServerList NrfServerList
	err := json.Unmarshal(data, &newNrfServerList)
	if err != nil {
		log.Errorf("failed to Unmarshal NrfServerList, %s\n", err.Error())
		return false
	}
	log.Debugf("%v\n", newNrfServerList)

	if !ValidateNrfServerList(&newNrfServerList) {
		log.Errorf("failed to validate NrfServerList, ignore this notification form cm")
		return false
	}

	nrfServerList = newNrfServerList
	return true
}

//GetNrfServerList get configuration nrfServerList by serviceName
func GetNrfServerList(c *NrfServerList) bool {
	if c == nil {
		return false
	}

	nrfAgentsConfLock.Lock()
	defer nrfAgentsConfLock.Unlock()

	*c = nrfServerList
	return true
}

//ValidateNrfServerList check validation of NrfServerList
func ValidateNrfServerList(c *NrfServerList) bool {
	nrfServerCnt := 0
	mgmtValidated := false
	discValidated := false

	for _, sp := range c.NrfServerProfileList {
		for _, ep := range sp.NrfServiceEndPoints {
			validated := ValidateNrfServiceEndPoint(&ep, &sp, true)

			switch ep.ServiceName {
			case NrfMgmtServiceName:
				mgmtValidated = validated
			case NrfDiscServiceName:
				discValidated = validated
			}
		}

		if mgmtValidated && discValidated {
			nrfServerCnt++
			continue
		} else {
			mgmtValidated = false
			discValidated = false
		}
	}

	if nrfServerCnt > 0 {
		return true
	} else {
		log.Errorf("no NrfServerProfile is completed configured in cm")
		return false
	}
}

//ValidateNrfServiceEndPoint check mandatory parameters in NrfServiceEndPoint
func ValidateNrfServiceEndPoint(ep *NrfServiceEndPoint, np *NrfServerProfile, isNFLevel bool) bool {
	validated := false

	if ep.Scheme == "" ||
		ep.ServiceName == "" ||
		len(ep.Versions) < 1 ||
		ep.Versions[0].APIVersionInUrI == "" ||
		(len(ep.IPEndPoints) < 1 &&
			ep.Fqdn == "") {
		return false
	}

	if ep.IPEndPoints != nil {
		for _, ipEp := range ep.IPEndPoints {
			if ipEp.Ipv4Address == "0.0.0.0" ||
				ipEp.Port == 0 {
				continue
			}

			validated = true
			break

		}
	}

	if ep.Fqdn != "" {
		validated = true
	}

	//Only be used when service level address can not be connectted
	if validated != true && isNFLevel == true && (len(np.NrfServerIpv4Address) > 0 || np.NrfServerFqdn != "") {
		for _, ipv4Addr := range np.NrfServerIpv4Address {
			if ipv4Addr != "0.0.0.0" {
				break
			}
		}
		if ep.IPEndPoints != nil {
			for _, ipEp := range ep.IPEndPoints {
				if ipEp.Port == 0 {
					continue
				}
				validated = true
				break
			}
		}
	}

	return validated
}

//UpdateStatusNotifIPEndPoint update data in nrf-agent configuration
func UpdateStatusNotifIPEndPoint(data []byte) bool {
	var newStatusNotifIPEndPoint StatusNotifIPEndPoint
	err := json.Unmarshal(data, &newStatusNotifIPEndPoint)
	if err != nil {
		log.Errorf("Failed to Unmarshal newStatusNotifIPEndPoint, %s\n", err.Error())
		return false
	}
	log.Debugf("UpdateStatusNotifIPEndPoint: newStatusNotifIPEndPoint is %+v\n", newStatusNotifIPEndPoint)

	nrfAgentNotifIPEndpoint = &newStatusNotifIPEndPoint
	return true
}

//GetStatusNotifIPEndPoint get configuration statusNotifIpEndpoint
func GetStatusNotifIPEndPoint() (*StatusNotifIPEndPoint, bool) {
	if nrfAgentNotifIPEndpoint == nil {
		return nil, false
	}

	return nrfAgentNotifIPEndpoint, true
}

//ValidateStatusNotifIPEndPointValue check mandatory parameters in nrfAgentNotifIPEndpoint
func ValidateStatusNotifIPEndPointValue(ep *StatusNotifIPEndPoint) bool {
	if ep.Ipv4Address == "" || ep.Port == 0 {
		return false
	}
	return true
}

//ValidateStatusNotifIPEndPoint for checking statusNotifIpEndpoint
func ValidateStatusNotifIPEndPoint() bool {
	_, ok := GetStatusNotifIPEndPoint()
	if !ok {
		return false
	} else {
		validated := ValidateStatusNotifIPEndPointValue(nrfAgentNotifIPEndpoint)
		return validated
	}

}
