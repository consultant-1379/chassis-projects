package cm

import (
	"encoding/json"
	"os"
	"strings"

	"fmt"
	"sort"
	"sync"

	"strconv"

	"gerrit.ericsson.se/udm/common/pkg/cmproxy"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

const (

	// ManagementWorkMode indicates nrf managemnt serivce
	ManagementWorkMode = "nrf_mgmt"
	// DiscoveryWorkMode indicates nrf discovery service
	DiscoveryWorkMode = "nrf_disc"
	// ProvsionWorkMode indicates nrf provision service
	ProvsionWorkMode = "nrf_prov"
	// ProvsionWorkMode indicates nrf provision service
	NotificationWorkMode = "nrf_notify"
)

// NrfProfileChangeHandler is function type fo callback for nrf profile change handler
type cmConfigChangeHandler func() bool

var (
	nrfProfileChangeHandler cmConfigChangeHandler
	plmnNrfChangeHandler    cmConfigChangeHandler
	discLocalCacheHandler   cmConfigChangeHandler
	//PlmnNrfAPIRootMap PlmnNRF address
	PlmnNrfAPIRootMap map[int][]string
	//PlmnNrfAPIRootInstanceIDMap  PlmnNrfAPIRoot mapping with Plmn ProfileID
	PlmnNrfAPIRootInstanceIDMap map[string]string
	//PlmnNrfPriority PlmnNRF address priority
	PlmnNrfPriority []int
	//Mutex for refresh plmnnrf address
	Mutex *sync.RWMutex = new(sync.RWMutex)
	//DiscNRFSelfAPIURI to store discovery service self open uri
	DiscNRFSelfAPIURI []string
	//PlmnNRFInstanceID to store plmn NRF instanceID, if some instance remove from cm, need delete cache entrys that from this instanceid
	PlmnNRFInstanceID []string
)

// RegisterNrfProfileChangeHandler is to register handler of NrfProfile change
func RegisterNrfProfileChangeHandler(callback cmConfigChangeHandler) {
	nrfProfileChangeHandler = callback
}

// RegisterPlmnNrfChangeHandler is to register handler of plmn nrf change
func RegisterPlmnNrfChangeHandler(callback cmConfigChangeHandler) {
	plmnNrfChangeHandler = callback
}

//RegisterDiscLocalCacheHandler is to register handler of discovery cache configuration change
func RegisterDiscLocalCacheHandler(callback cmConfigChangeHandler) {
	discLocalCacheHandler = callback
}

// Setup to register and update configuration of nrf
func Setup() {

	log.Debugf("Initialize CM Proxy function")

	cmmService := os.Getenv("CMM_SERVICE") + "/"
	cmproxy.Init(cmmService)

	configSchemaName := os.Getenv("CMM_CONFIG_SCHEMA_NAME")
	configModuleName := os.Getenv("CMM_CONFIG_MODULE_NAME")
	configTopicName := os.Getenv("CMM_CONFIG_TOPIC")
	cmproxy.RegisterConf(configSchemaName, configModuleName, configTopicName, cmUpdateHandler, cmproxy.NtfFormatFull)

	cmproxy.Run()

	initOtherConf()
}

func cmUpdateHandler(Event, ConfigurationName, format string, RawData []byte) {

	log.Debugf("configuration is updated")

	if format != cmproxy.NtfFormatFull {
		return
	}

	log.Debugf(string(RawData))

	var message EricssonNrfFunction
	err := json.Unmarshal(RawData, &message)
	if err != nil {
		log.Errorf("decode error %v", err)
		return
	}

	updateNrfCommon(&message)

	updateServiceProfile(&message, ServiceName)

	updateNFProfile(&message)

	updateNrfPolicy(&message)

	updateNFServiceLog(&message)

	if strings.EqualFold(ServiceName, ManagementWorkMode) {
		if NrfCommon.Role == "region-nrf" {
			ret := nrfProfileChangeHandler()
			if !ret {
				log.Errorf("Send nrfProfile update to PLMN NRF failed.")
			}
		}
	}

	if strings.EqualFold(ServiceName, DiscoveryWorkMode) {
		ret := discLocalCacheHandler()
		if !ret {
			log.Errorf("update discovery local cache fail")
		}
	}

	if NrfCommon.Role == "region-nrf" {
		name := "nnrf-nfm"
		if strings.EqualFold(ServiceName, DiscoveryWorkMode) {
			name = constvalue.NNRFDISC
		}
		ret := ConstructPlmnNrfAPIRoot(name)
		if !ret {
			log.Errorf(" construct PLMN nrf address failed.")
		}
	}

	if strings.EqualFold(ServiceName, DiscoveryWorkMode) {
		constructDiscNrfSelfAPIURI()
		if NrfCommon.PlmnNrf != nil &&  len(NrfCommon.PlmnNrf.Profile) > 0 {
			for _, nrf := range NrfCommon.PlmnNrf.Profile {
				PlmnNRFInstanceID = append(PlmnNRFInstanceID, nrf.ID)
			}
		}
	}
}

func constructDiscNrfSelfAPIURI() {

	Mutex.Lock()
	DiscNRFSelfAPIURI = make([]string, 0)
	peerScheme := NrfCommon.RemoteDefaultSetting.Scheme
	peerPort := NrfCommon.RemoteDefaultSetting.Port
	DiscNRFSelfAPIURI = append(DiscNRFSelfAPIURI, peerScheme+"://"+NfProfile.Fqdn+":"+strconv.Itoa(peerPort)+"/nnrf-disc/v1/")
	for _, v4 := range NfProfile.Ipv4Address {
		DiscNRFSelfAPIURI = append(DiscNRFSelfAPIURI, peerScheme+"://"+v4+":"+strconv.Itoa(peerPort)+"/nnrf-disc/v1/")
	}

	for _, v6 := range NfProfile.Ipv6Address {
		DiscNRFSelfAPIURI = append(DiscNRFSelfAPIURI, peerScheme+"://"+v6+":"+strconv.Itoa(peerPort)+"/nnrf-disc/v1/")
	}

	if DiscoveryNfServices != nil {
		for _, s := range DiscoveryNfServices {
			for _, v := range s.Version {
				if s.IPEndpoint != nil {
					for _, ss := range s.IPEndpoint {
						DiscNRFSelfAPIURI = append(DiscNRFSelfAPIURI, s.Scheme+"://"+ss.Ipv4Address+":"+strconv.Itoa(ss.Port)+"/nnrf-disc/"+v.APIVersionInURI+"/")
						DiscNRFSelfAPIURI = append(DiscNRFSelfAPIURI, s.Scheme+"://"+ss.Ipv6Address+":"+strconv.Itoa(ss.Port)+"/nnrf-disc/"+v.APIVersionInURI+"/")
						DiscNRFSelfAPIURI = append(DiscNRFSelfAPIURI, s.Scheme+"://"+s.Fqdn+":"+strconv.Itoa(ss.Port)+"/nnrf-disc"+v.APIVersionInURI+"/")
					}
				}
			}

		}
	}

	Mutex.Unlock()

	if len(DiscNRFSelfAPIURI) == 0 {
		log.Warningf("Please configure NRF Discovery API Root configuration in nf-profile and discovery-nf-service")
		return
	}

	for _, v := range DiscNRFSelfAPIURI {
		log.Debugf("NRF Self Root API : %s", v)
	}

}

// ConstructPlmnNrfAPIRoot is to construct apiRoot of Plmn nrf
func ConstructPlmnNrfAPIRoot(serviceName string) bool {

	var plmnNrfHostPortMap map[string]int

	if NrfCommon.PlmnNrf == nil {
		log.Warnf("cm.NrfCommon.PlmnNrf == nil")
		return false
	}

	if len(NrfCommon.PlmnNrf.Profile) == 0 {
		log.Warnf("cm.NrfCommon.PlmnNrf.Profile is empty")
		return false
	}

	Mutex.Lock()
	PlmnNrfAPIRootMap = make(map[int][]string)
	PlmnNrfAPIRootInstanceIDMap = make(map[string]string)
	PlmnNrfPriority = append([]int{})
	for _, nrfProfile := range NrfCommon.PlmnNrf.Profile {
		profileID := nrfProfile.ID
		nrfServices := nrfProfile.Service

		log.Debugf("the nrfServices is %v", nrfServices)
		for _, nrfService := range nrfServices {
			var scheme, fqdn string
			if nrfService.Name == serviceName {
				scheme = nrfService.Scheme
				fqdn = nrfService.Fqdn
				plmnNrfHostPortMap = parsePlmnNrfHostPortByNrfService(nrfService.IPEndpoint, scheme, fqdn)

				if len(plmnNrfHostPortMap) == 0 {
					fqdn = nrfProfile.Fqdn
					plmnNrfHostPortMap = parsePlmnNrfHostPortByNrfProfile(nrfProfile.Ipv4Address, nrfProfile.Ipv6Address, nrfService.IPEndpoint, scheme, fqdn)
				}

				priority := 0
				if nrfService.Priority != nil {
					priority = *(nrfService.Priority)
				} else {
					if nrfProfile.Priority != nil {
						priority = *(nrfProfile.Priority)
					}
				}

				for host, port := range plmnNrfHostPortMap {
					var apiRoot string
					if nrfService.APIPrefix != "" {
						apiRoot = fmt.Sprintf("%s://%s:%d/%s", scheme, host, port, nrfService.APIPrefix)
					} else {
						apiRoot = fmt.Sprintf("%s://%s:%d", scheme, host, port)
					}
					PlmnNrfAPIRootMap[priority] = append(PlmnNrfAPIRootMap[priority], apiRoot)
					PlmnNrfAPIRootInstanceIDMap[apiRoot] = profileID
				}
			}

		}
	}
	for k := range PlmnNrfAPIRootMap {
		PlmnNrfPriority = append(PlmnNrfPriority, k)
	}
	sort.Ints(PlmnNrfPriority)
	Mutex.Unlock()

	for _, priority := range PlmnNrfPriority {
		for _, apiRoot := range PlmnNrfAPIRootMap[priority] {
			log.Debugf("PLMN NRF apiRoot is %s", apiRoot)
		}
	}
	return true
}

func parsePlmnNrfHostPortByNrfService(IPEndpoints []TIPEndpoint, scheme string, fqdn string) map[string]int {
	var port int
	plmnNrfHostPortMap := make(map[string]int)
	port = 0

	for _, ipEndpoint := range IPEndpoints {
		port = 0
		if ipEndpoint.Port > 0 {
			port = ipEndpoint.Port
		}
		ipAddress := ""
		if ipEndpoint.Ipv4Address != "" {
			ipAddress = ipEndpoint.Ipv4Address
		} else if ipEndpoint.Ipv4Address == "" && ipEndpoint.Ipv6Address != "" {
			ipAddress = "[" + ipEndpoint.Ipv6Address + "]"
		}
		if ipAddress != "" {
			if port == 0 {
				if scheme == "https" {
					port = 443
				} else if scheme == "http" {
					port = 80
				}
			}
			plmnNrfHostPortMap[ipAddress] = port
		}
	}

	if len(plmnNrfHostPortMap) == 0 && fqdn != "" {

		if port == 0 {
			if scheme == "https" {
				port = 443
			} else if scheme == "http" {
				port = 80
			}
		}
		plmnNrfHostPortMap[fqdn] = port
	}

	return plmnNrfHostPortMap
}

func parsePlmnNrfHostPortByNrfProfile(IPv4Address, IPv6Address []string, IPEndpoints []TIPEndpoint, scheme string, fqdn string) map[string]int {

	plmnNrfHostPortMap := make(map[string]int)
	port := 0

	for _, ipEndpoint := range IPEndpoints {
		if ipEndpoint.Port > 0 {
			port = ipEndpoint.Port
			break
		}
	}

	if len(IPv4Address) > 0 {
		for _, vIPv4Address := range IPv4Address {
			if vIPv4Address != "" {
				if port == 0 {
					if scheme == "https" {
						port = 443
					} else if scheme == "http" {
						port = 80
					}
				}
				plmnNrfHostPortMap[vIPv4Address] = port
			}
		}
	} else {
		for _, vIPv6Address := range IPv6Address {
			ipAddress := "[" + vIPv6Address + "]"
			if ipAddress != "" {
				if port == 0 {
					if scheme == "https" {
						port = 443
					} else if scheme == "http" {
						port = 80
					}
				}
				plmnNrfHostPortMap[ipAddress] = port
			}
		}
	}

	if len(plmnNrfHostPortMap) == 0 && fqdn != "" {
		if port == 0 {
			if scheme == "https" {
				port = 443
			} else if scheme == "http" {
				port = 80
			}
		}
		plmnNrfHostPortMap[fqdn] = port
	}

	return plmnNrfHostPortMap
}

func updateNrfCommon(message *EricssonNrfFunction) {
	nrfCommon := message.NrfCommon
	if nrfCommon != nil {
		nrfCommon.ParseConf()
	}
}

func updateServiceProfile(message *EricssonNrfFunction, serviceName string) {

	if strings.EqualFold(serviceName, ManagementWorkMode) || strings.EqualFold(serviceName, NotificationWorkMode) {
		serviceProfile := message.ManagementService
		if serviceProfile != nil {
			serviceProfile.ParseConf()
			log.Debugf("update service profile successfully")
			serviceProfile.Show()
		}

	} else if strings.EqualFold(serviceName, DiscoveryWorkMode) {
		serviceProfile := message.DiscoveryService
		if serviceProfile != nil {
			serviceProfile.ParseConf()
			log.Debugf("update service profile successfully")
			serviceProfile.Show()
		}
	} else if strings.EqualFold(serviceName, ProvsionWorkMode) {
		serviceProfile := message.ProvisionService
		if serviceProfile != nil {
			serviceProfile.ParseConf()
			log.Debugf("update provision service profile successfully")
			serviceProfile.Show()
		}
	} else {
		log.Warnf("ignore work mode " + serviceName)
	}
}

func updateNFProfile(message *EricssonNrfFunction) {
	nrfNFProfile := message.NfProfile
	if nrfNFProfile != nil {
		nrfNFProfile.ParseConf()
		log.Debugf("update nf-profile successfully")
		nrfNFProfile.Show()
	}
}

func updateNrfPolicy(message *EricssonNrfFunction) {
	nrfPolicy := message.NrfPolicy
	if nrfPolicy != nil {
		nrfPolicy.ParseConf()
		log.Debugf("update nrf policy successfully")
		nrfPolicy.Show()
	}
}

func updateNFServiceLog(message *EricssonNrfFunction) {

	log.Debugf("update service log conf for " + ServiceName)

	nfServiceLogs := message.NfServiceLogs
	for _, serviceLog := range nfServiceLogs {

		if strings.EqualFold(serviceLog.LogID, ServiceName) {
			updateServiceLog(&serviceLog)
		} else if strings.EqualFold(ServiceName, NotificationWorkMode) && strings.EqualFold(serviceLog.LogID, ManagementWorkMode) {
			updateServiceLog(&serviceLog)
		} else {
			log.Infof("ignore log-id = " + serviceLog.LogID)
		}
	}
}

func updateServiceLog(serviceLog *TNfServiceLog) {
	if serviceLog != nil {
		serviceLog.ParseConf()

		log.Debugf("update service log conf successfully")
		serviceLog.Show()
	}
}

func checkFileIsExist(filename string) bool {

	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist

}

// except from CM mediator, get other configruation from env, and etc
func initOtherConf() {
	SetServiceVersion()
}
