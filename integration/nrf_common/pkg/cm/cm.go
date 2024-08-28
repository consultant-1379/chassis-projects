package cm

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"unsafe"

	"gerrit.ericsson.se/udm/common/pkg/cmproxy"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

const (

	// ManagementWorkMode indicates nrf managemnt serivce
	ManagementWorkMode = constvalue.APP_WORKMODE_NRF_MGMT
	// DiscoveryWorkMode indicates nrf discovery service
	DiscoveryWorkMode = constvalue.APP_WORKMODE_NRF_DISC
	// ProvsionWorkMode indicates nrf provision service
	ProvsionWorkMode = constvalue.APP_WORKMODE_NRF_PROV
	// ProvsionWorkMode indicates nrf provision service
	NotificationWorkMode = constvalue.AppWorkmodeNrfNotif
)

// NrfProfileChangeHandler is function type fo callback for nrf profile change handler
type cmConfigChangeHandler func() bool

var (
	cmChangeHandler cmConfigChangeHandler
	//PeerNRFInfoMap to store plmnnrfurl selfurl ipaddress of all site
	PeerNRFInfoMap *TPeerNRFsInfo

	//NRFURLInfo to store  plmnnrfurl selfurl ipaddress
	NRFURLInfo *TNRFURLInfo

	//MyURLInfo to store my url info
	MyURLInfo *TMyURLInfo
)

//TNRFURLInfo to store peer NRF(upper layer, same layer NRF) url info, id and so on
type TNRFURLInfo struct {
	//PeerNrfAPIRootMap PeerNRF address
	PeerNrfAPIRootMap map[int][]string
	//PeerNrfAPIRootInstanceIDMap  PeerNrfAPIRoot mapping with Plmn Profile ID
	PeerNrfAPIRootInstanceIDMap map[string]string
	//PeerNrfPriority PeerNRF address priority
	PeerNrfPriority []int
	//PeerNRFInstanceID to store peer NRF instanceID, if some instance remove from cm, need delete cache entrys that from this instanceid
	PeerNRFInstanceID []string
	//PeerNRFAddressIdentifier to store peer NRF fqdn, ipv4 address, ipv6 address, used to check against http header forwarded.
	PeerNRFAddressIdentifier []string
}

//AtomicSetNRFURLInfo to set nrfinfo that contain plmnurl selfurl
func (n *TNRFURLInfo) AtomicSetNRFURLInfo() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&NRFURLInfo)), unsafe.Pointer(n))
}

//GetNRFURLInfo  to get nrfinfo
func GetNRFURLInfo() *TNRFURLInfo {
	return (*TNRFURLInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&NRFURLInfo))))
}
func (n *TNRFURLInfo) init() {
	n.PeerNrfAPIRootMap = make(map[int][]string)
	n.PeerNrfAPIRootInstanceIDMap = make(map[string]string)
	n.PeerNrfPriority = make([]int, 0)
	n.PeerNRFInstanceID = make([]string, 0)
	n.PeerNRFAddressIdentifier = make([]string, 0)
}

// TMyURLInfo contain all kinds of URI info of mine
type TMyURLInfo struct {
	//MyAddressIdentifier to store my fqdn, ipv4 address, ipv6 address, used to construct http header forwarded.
	MyAddressIdentifier []string
	//MyDiscAPIURI to store discovery service self open uri
	MyDiscAPIURI []string
}

func (n *TMyURLInfo) init() {
	n.MyDiscAPIURI = make([]string, 0)
	n.MyAddressIdentifier = make([]string, 0)
}

//AtomicSetMyURLInfo to set nrfinfo that contain plmnurl selfurl
func (n *TMyURLInfo) AtomicSetMyURLInfo() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&MyURLInfo)), unsafe.Pointer(n))
}

//GetMyURLInfo  to get nrfinfo
func GetMyURLInfo() *TMyURLInfo {
	return (*TMyURLInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&MyURLInfo))))
}

// TPeerNRFInfo to store peer nrf urlinfo, layer ...
type TPeerNRFInfo struct {
	NRFURLInfo *TNRFURLInfo
	Layer      string
	Name       string
}

// String is to print the peernrf info
func (p TPeerNRFInfo) String() string {
	return fmt.Sprintf("next-hop NRF: {Name: %s, Layer: %s}\n", p.Name, p.Layer)
}

// TPeerNRFsInfo contains all peer NRF info
type TPeerNRFsInfo map[string]*TPeerNRFInfo

// AtomicSetPeerNRFInfoMap to set peer nrfinfomap
func (p *TPeerNRFsInfo) AtomicSetPeerNRFInfoMap() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&PeerNRFInfoMap)), unsafe.Pointer(p))
}

// CheckLayerConsistent to check whether layer is consistent
func (p *TPeerNRFsInfo) CheckLayerConsistent() bool {
	if len(*p) == 0 {
		return true
	}

	lastLayer := "defaultLayer"
	for _, peerNRFInfo := range *p {
		if lastLayer != "defaultLayer" {
			if lastLayer != peerNRFInfo.Layer {
				log.Warnf("Peer NRF %s 's layer is %s, different with other NRF in next-hop", peerNRFInfo.Name, peerNRFInfo.Layer)
				return false
			}
		}
		lastLayer = peerNRFInfo.Layer
	}
	log.Debugf("NRF in next-hop is in the %s with me.", lastLayer)
	return true
}

// GetPeerNRFInfoMap  to get peer nrfinfomap
func GetPeerNRFInfoMap() *TPeerNRFsInfo {
	return (*TPeerNRFsInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&PeerNRFInfoMap))))
}

// RegisterCMChangeHandler is to register handler of discovery overload redirection configuration change
func RegisterCMChangeHandler(callback cmConfigChangeHandler) {
	cmChangeHandler = callback
}

// Setup to register and update configuration of nrf
func Setup() {

	log.Debugf("Initialize CM Proxy function")

	cmmService := os.Getenv("CMM_SERVICE") + "/"
	cmproxy.Init(cmmService)
	initCM()
	configSchemaName := os.Getenv("CMM_CONFIG_SCHEMA_NAME")
	configModuleName := os.Getenv("CMM_CONFIG_MODULE_NAME")
	configTopicName := os.Getenv("CMM_CONFIG_TOPIC")
	cmproxy.RegisterConf(configSchemaName, configModuleName, configTopicName, cmUpdateHandler, cmproxy.NtfFormatFull)

	cmproxy.Run()

	initOtherConf()
}

//proctection to avoid get null point
func initCM() {
	var nrfURLInfo TNRFURLInfo
	nrfURLInfo.init()
	nrfURLInfo.AtomicSetNRFURLInfo()

	peerNRFInfoMap := make(TPeerNRFsInfo)
	peerNRFInfoMap["init"] = &TPeerNRFInfo{
		NRFURLInfo: &nrfURLInfo,
	}
	peerNRFInfoMap.AtomicSetPeerNRFInfoMap()

	var myURLInfo TMyURLInfo
	myURLInfo.init()
	myURLInfo.AtomicSetMyURLInfo()

	var discoveryService TDiscoveryService
	discoveryService.atomicSetDiscService()

	var managementService TManagementService
	managementService.atomicSetMgmtService()

	var nfProfile TNfProfile
	nfProfile.atomicSetNFProfile()

	var plmnListRaw TPlmnListRawData
	plmnListRaw.atomicSetPlmnListRaw()

	var nrfNFServices TNRFNFServices
	nrfNFServices.atomicSetNRFNFServices()

	var nrfCommon TCommon
	nrfCommon.atomicSetCommon()

	var nrfPolicy TNrfPolicy
	nrfPolicy.atomicSetNRFPolicy()

	var provisionService TProvisionService
	provisionService.atomicSetPrivisionService()

	var ingressAddress TProvIngressAddress
	ingressAddress.atomicSetIngressAddress()
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

	var myURLInfo TMyURLInfo
	if strings.EqualFold(ServiceName, ManagementWorkMode) {
		myURLInfo.ConstructMyAddressIdentifier(constvalue.NNRFNFM)
	} else if strings.EqualFold(ServiceName, DiscoveryWorkMode) {
		myURLInfo.ConstructMyAddressIdentifier(constvalue.NNRFDISC)
	}

	if strings.EqualFold(ServiceName, DiscoveryWorkMode) {
		myURLInfo.constructMyDiscAPIURI()
	}
	myURLInfo.AtomicSetMyURLInfo()

	if GetNRFRole() == REGIONNRF {
		peerNRFInfoMap := make(TPeerNRFsInfo)
		if GetNRFCommon().RegionNrf != nil && GetNRFCommon().RegionNrf.NextHop != nil && len(GetNRFCommon().RegionNrf.NextHop.Site) > 0 {

			for _, site := range GetNRFCommon().RegionNrf.NextHop.Site {
				name := constvalue.NNRFNFM
				if strings.EqualFold(ServiceName, DiscoveryWorkMode) {
					name = constvalue.NNRFDISC
				}

				nrfURLInfo := TNRFURLInfo{}
				nrfURLInfo.ConstructPeerNrfAddressIdentifier(name, site.Profile, site.ID)
				ret := nrfURLInfo.ConstructPeerNrfAPIRoot(name, site.Profile)
				if !ret {
					log.Errorf("construct Peer nrf address failed.")
				}

				if name == DiscoveryWorkMode {
					for _, nrfProfile := range site.Profile {
						nrfURLInfo.PeerNRFInstanceID = append(nrfURLInfo.PeerNRFInstanceID, nrfProfile.ID)
					}

				}
				peerNRFInfo := &TPeerNRFInfo{
					Name:       site.ID,
					Layer:      site.Layer,
					NRFURLInfo: &nrfURLInfo,
				}
				log.Debugf(peerNRFInfo.String())
				peerNRFInfoMap[site.ID] = peerNRFInfo
			}
		}
		peerNRFInfoMap.CheckLayerConsistent()
		peerNRFInfoMap.AtomicSetPeerNRFInfoMap()
	}

	if cmChangeHandler != nil {
		ret := cmChangeHandler()
		if !ret {
			log.Errorf("call cmChangeHandler fail")
		}
	}
}

//ConstructPeerNrfAddressIdentifier is to construct peer nrf address identifier
func (n *TNRFURLInfo) ConstructPeerNrfAddressIdentifier(serviceName string, peerNRFProfile []TNrfProfile, peerName string) {
	if GetNRFRole() != REGIONNRF {
		log.Warningf("Current this function only suits for RegionNRF Role")
		return
	}

	n.PeerNRFAddressIdentifier = make([]string, 0)
	for _, nrfProfile := range peerNRFProfile {
		nrfServices := nrfProfile.Service
		for _, nrfService := range nrfServices {
			if nrfService.Name == serviceName {
				for _, ipEndpoint := range nrfService.IPEndpoint {
					if ipEndpoint.Ipv4Address != "" {
						n.PeerNRFAddressIdentifier = append(n.PeerNRFAddressIdentifier, ipEndpoint.Ipv4Address)
					}
					if ipEndpoint.Ipv6Address != "" {
						n.PeerNRFAddressIdentifier = append(n.PeerNRFAddressIdentifier, ipEndpoint.Ipv6Address)
					}
				}

				if nrfService.Fqdn != "" {
					n.PeerNRFAddressIdentifier = append(n.PeerNRFAddressIdentifier, nrfService.Fqdn)
				}
			}
		}
		if len(nrfProfile.Ipv4Address) > 0 {
			n.PeerNRFAddressIdentifier = append(n.PeerNRFAddressIdentifier, nrfProfile.Ipv4Address...)
		}
		if len(nrfProfile.Ipv6Address) > 0 {
			n.PeerNRFAddressIdentifier = append(n.PeerNRFAddressIdentifier, nrfProfile.Ipv6Address...)
		}
		if nrfProfile.Fqdn != "" {
			n.PeerNRFAddressIdentifier = append(n.PeerNRFAddressIdentifier, nrfProfile.Fqdn)
		}
	}

	if len(n.PeerNRFAddressIdentifier) == 0 {
		log.Warningf("Please configure Address identifier for next-hop peer NRF %s, such as fqdn, ipv4 address, ipv6 address", peerName)
	}
	log.Debugf("Current NRF Role :%s, peer NRF %s's Address Identifier : %v", GetNRFRole(), peerName, n.PeerNRFAddressIdentifier)
}

//ConstructMyAddressIdentifier is to construct plmn nrf address identifier
func (n *TMyURLInfo) ConstructMyAddressIdentifier(serviceName string) {
	n.MyAddressIdentifier = make([]string, 0)
	nfProfile := GetNRFNFProfile()

	for _, nrfService := range nfProfile.Service {
		if nrfService.Name == serviceName {
			for _, ipEndpoint := range nrfService.IPEndpoint {
				if ipEndpoint.Ipv4Address != "" {
					n.MyAddressIdentifier = append(n.MyAddressIdentifier, ipEndpoint.Ipv4Address)
				}
				if ipEndpoint.Ipv6Address != "" {
					n.MyAddressIdentifier = append(n.MyAddressIdentifier, ipEndpoint.Ipv6Address)
				}
			}

			if nrfService.Fqdn != "" {
				n.MyAddressIdentifier = append(n.MyAddressIdentifier, nrfService.Fqdn)
			}
		}
	}

	if len(nfProfile.Ipv4Address) != 0 {
		n.MyAddressIdentifier = append(n.MyAddressIdentifier, nfProfile.Ipv4Address...)
	}
	if len(nfProfile.Ipv6Address) != 0 {
		n.MyAddressIdentifier = append(n.MyAddressIdentifier, nfProfile.Ipv6Address...)
	}

	if nfProfile.Fqdn != "" {
		n.MyAddressIdentifier = append(n.MyAddressIdentifier, nfProfile.Fqdn)
	}

	if len(n.MyAddressIdentifier) == 0 {
		log.Warningf("Please configure NRF Address identifier in nfprofile to identify itself, such as fqdn, ipv4 address, ipv6 address")
	}

	log.Debugf("Current NRF Role :%s, My Address Identifier : %v", GetNRFRole(), n.MyAddressIdentifier)
}

func (n *TMyURLInfo) constructMyDiscAPIURI() {

	n.MyDiscAPIURI = make([]string, 0)
	peerScheme := constvalue.RemoteDefaultScheme
	peerPort := constvalue.RemoteDefaultPort
	nfProfile := GetNRFNFProfile()
	n.MyDiscAPIURI = append(n.MyDiscAPIURI, peerScheme+"://"+nfProfile.Fqdn+":"+strconv.Itoa(peerPort)+"/nnrf-disc/v1/")
	for _, v4 := range NfProfile.Ipv4Address {
		n.MyDiscAPIURI = append(n.MyDiscAPIURI, peerScheme+"://"+v4+":"+strconv.Itoa(peerPort)+"/nnrf-disc/v1/")
	}

	for _, v6 := range NfProfile.Ipv6Address {
		n.MyDiscAPIURI = append(n.MyDiscAPIURI, peerScheme+"://["+v6+"]:"+strconv.Itoa(peerPort)+"/nnrf-disc/v1/")
	}

	if GetNRFNFServices().DiscoveryNfServices != nil {
		for _, s := range GetNRFNFServices().DiscoveryNfServices {
			for _, v := range s.Version {
				if s.IPEndpoint != nil {
					for _, ss := range s.IPEndpoint {
						n.MyDiscAPIURI = append(n.MyDiscAPIURI, s.Scheme+"://"+ss.Ipv4Address+":"+strconv.Itoa(ss.Port)+"/nnrf-disc/"+v.APIVersionInURI+"/")
						n.MyDiscAPIURI = append(n.MyDiscAPIURI, s.Scheme+"://["+ss.Ipv6Address+"]:"+strconv.Itoa(ss.Port)+"/nnrf-disc/"+v.APIVersionInURI+"/")
						n.MyDiscAPIURI = append(n.MyDiscAPIURI, s.Scheme+"://"+s.Fqdn+":"+strconv.Itoa(ss.Port)+"/nnrf-disc"+v.APIVersionInURI+"/")
					}
				}
			}

		}
	}

	if len(n.MyDiscAPIURI) == 0 {
		log.Warningf("Please configure NRF Discovery API Root configuration in nf-profile and discovery-nf-service")
		return
	}

	for _, v := range n.MyDiscAPIURI {
		log.Debugf("NRF Self Root API : %s", v)
	}

}

// ConstructPeerNrfAPIRoot is to construct apiRoot of Plmn nrf
func (n *TNRFURLInfo) ConstructPeerNrfAPIRoot(serviceName string, peerNRFProfile []TNrfProfile) bool {

	var plmnNrfHostPortMap map[string]int

	n.PeerNrfAPIRootMap = make(map[int][]string)
	n.PeerNrfAPIRootInstanceIDMap = make(map[string]string)
	n.PeerNrfPriority = append([]int{})
	for _, nrfProfile := range peerNRFProfile {
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
					n.PeerNrfAPIRootMap[priority] = append(n.PeerNrfAPIRootMap[priority], apiRoot)
					n.PeerNrfAPIRootInstanceIDMap[apiRoot] = profileID
				}
			}

		}
	}
	for k := range n.PeerNrfAPIRootMap {
		n.PeerNrfPriority = append(n.PeerNrfPriority, k)
	}
	sort.Ints(n.PeerNrfPriority)

	for _, priority := range n.PeerNrfPriority {
		for _, apiRoot := range n.PeerNrfAPIRootMap[priority] {
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
			log.Debugf("update mgmt service profile successfully")
			serviceProfile.Show()
		} else {
			SetDefaultForManagementService()
			log.Debugf("management service use default configuration ")
		}

	} else if strings.EqualFold(serviceName, DiscoveryWorkMode) {
		serviceProfile := message.DiscoveryService
		if serviceProfile != nil {
			serviceProfile.ParseConf()
			log.Debugf("update disc service profile successfully")
			serviceProfile.Show()
		} else {
			SetDefaultForDiscoveryService()
			log.Debugf("discovery service use default configuration ")
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
