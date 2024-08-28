package client

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/msgbus"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/fm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

var (
	monitorRole  string
	monitorMSB   *msgbus.MessageBus
	monitorTopic = consts.MsgbusTopicNamePrefix + "connections"
	monitorChan  chan bool

	nrfServerConnLock sync.Mutex
	nrfMgmtURLPrefix  string
	nrfDiscURLPrefix  string
	syncNrfStatus     SyncNrfStatus
	nrfStatusLock     sync.Mutex

	defaultRetryTimes = 3

	nrfConnStatus NRFConnStatus
)

type SyncNrfStatus int

const (
	SyncNrfStatusUnknown SyncNrfStatus = 0
	SyncNrfStatusOK      SyncNrfStatus = 1
	SyncNrfStatusNOK     SyncNrfStatus = 2
	SyncNrfStatusFailure SyncNrfStatus = 3
)

type NRFConnStatus int

const (
	NRFConnUnknown NRFConnStatus = 0
	NRFConnNormal  NRFConnStatus = 1
	NRFConnLost    NRFConnStatus = 2
)
const (
	//PrimaryMonitor is primary role
	PrimaryMonitor = "primaryRole"
	//SeconaryMonitor is secondary role
	SeconaryMonitor = "secondaryRole"
)

type monitorMessage struct {
	Role       string `json:"monitorRole"`
	MgmtPrefix string `json:"mgmtUrlPrefix"`
	DiscPrefix string `json:"discUrlPrefix"`
	ForceReq   bool   `json:"forceReq,omitempty"` //ForceReq is for SeconaryMonitor fource request NRF url from PrimaryMonitor
}

func getNrfServerPrefix() (string, string) {
	nrfServerConnLock.Lock()
	defer nrfServerConnLock.Unlock()
	return nrfMgmtURLPrefix, nrfDiscURLPrefix
}

func setNrfServerPrefix(nrfMgmtPrefix, nrfDiscPrefix string) {
	nrfServerConnLock.Lock()
	defer nrfServerConnLock.Unlock()
	nrfMgmtURLPrefix = nrfMgmtPrefix
	nrfDiscURLPrefix = nrfDiscPrefix
}

// ResetNrfServerPrefix trigger re-select behavior
func ResetNrfServerPrefix() {
	setNrfServerPrefix("", "")
	if monitorChan != nil {
		monitorChan <- true
	}
}

// HTTPDoToNrfMgmt HTTPDo to NRF Managemnet
var HTTPDoToNrfMgmt = func(httpv, method, urlSuffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
	resp, err := request2NrfServer(structs.NrfMgmtServiceName, httpv, method, urlSuffix, hdr, body)
	// The code retry to connect to NRF management is in regagent&regproxy.
	// There is no need to retry again in common connections package
	for retryCnt := 0; err != nil && retryCnt < 0; retryCnt++ {
		resp, err = request2NrfServer(structs.NrfMgmtServiceName, httpv, method, urlSuffix, hdr, body)
	}
	if err != nil {
		fm.ConnectionStatus("nrf-mgmt", false)
	} else {
		fm.ConnectionStatus("nrf-mgmt", true)
	}
	return resp, err
}

// HTTPDoToNrfDisc HTTPDo to NRF Discovery
var HTTPDoToNrfDisc = func(httpv, method, urlSuffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
	resp, err := request2NrfServer(structs.NrfDiscServiceName, httpv, method, urlSuffix, hdr, body)
	//retry to send requests
	for retryCnt := 0; err != nil && retryCnt < defaultRetryTimes; retryCnt++ {
		resp, err = request2NrfServer(structs.NrfDiscServiceName, httpv, method, urlSuffix, hdr, body)
	}
	if err != nil {
		fm.ConnectionStatus("nrf-disc", false)
	} else {
		fm.ConnectionStatus("nrf-disc", true)
	}
	return resp, err
}

// InitializeMonitor initialize NRF server monitor
func InitializeMonitor(role string) error {
	if role != PrimaryMonitor &&
		role != SeconaryMonitor {
		return errors.New("unknow NRF Server monitor role: " + role)
	}
	monitorRole = role
	{
		nrfStatusLock.Lock()
		defer nrfStatusLock.Unlock()
		syncNrfStatus = SyncNrfStatusUnknown
	}
	nrfConnStatus = NRFConnUnknown

	monitorMSB = msgbus.NewMessageBus(os.Getenv("MESSAGE_BUS_KAFKA"))
	if monitorMSB == nil {
		return errors.New("failed to initialize message bus for NRF Server monitor")
	}
	log.Infof("initialize message bus for NRF Server monitor done")

	err := monitorMSB.ConsumeTopic(monitorTopic, monitorMessageHandler)
	if err != nil {
		log.Errorf("failed to ConsumeTopic %s for NRF Server monitor", monitorTopic)
		return err
	}

	go func() {
		monitorLoop(time.Duration(cm.GetHTTP2Timeout()) * time.Second)
	}()

	return nil
}

func monitorMessageHandler(msg []byte) {
	log.Infof("monitorMessageHandler: %s", string(msg))

	var mm monitorMessage
	err := json.Unmarshal(msg, &mm)
	if err != nil {
		log.Errorf("invalid message from message bus, %s", err.Error())
		return
	}
	if mm.Role == monitorRole {
		return
	}

	switch monitorRole {
	case PrimaryMonitor:
		primaryMonitorMessageHandler(&mm)
	case SeconaryMonitor:
		secondaryMonitorMessageHandler(&mm)
	default:
		log.Errorf("unknow NRF Server monitor role: %s", monitorRole)
	}
}

func primaryMonitorMessageHandler(mm *monitorMessage) {
	var err error

	mgmtPrefix, discPrefix := getNrfServerPrefix()
	if mgmtPrefix == "" ||
		discPrefix == "" {
		ResetNrfServerPrefix()
		sendResponseForForceReq(mm)
		return
	}

	_, mgmtValid := hb2NrfServer(mgmtPrefix, structs.NrfMgmtServiceName)
	if !mgmtValid {
		log.Errorf("primaryMonitorMessageHandler: lost connection to NRF Management.")
		ResetNrfServerPrefix()
		sendResponseForForceReq(mm)
		return
	}
	_, discValid := hb2NrfServer(discPrefix, structs.NrfDiscServiceName)
	if !discValid {
		log.Errorf("primaryMonitorMessageHandler: lost connection to NRF Discovery, %v", discValid)
		ResetNrfServerPrefix()
		sendResponseForForceReq(mm)
		return
	}

	respMm := &monitorMessage{
		Role:       monitorRole,
		MgmtPrefix: mgmtPrefix,
		DiscPrefix: discPrefix,
	}
	err = sendMonitorMessage(respMm)
	if err != nil {
		log.Errorf("%s", err.Error())
		nrfStatusLock.Lock()
		defer nrfStatusLock.Unlock()
		syncNrfStatus = SyncNrfStatusFailure
	} else {
		nrfStatusLock.Lock()
		defer nrfStatusLock.Unlock()
		syncNrfStatus = SyncNrfStatusOK
	}
}

func secondaryMonitorMessageHandler(mm *monitorMessage) {
	if mm.Role != PrimaryMonitor {
		return
	}

	if mm.MgmtPrefix == "" || mm.DiscPrefix == "" {
		//NRFConnUnknown status will not change to NRFConnLost directly, only NRFConnNormal may change to NRFConnLost
		if nrfConnStatus == NRFConnNormal {
			log.Warningf("NRF status set to unavilable, NRF Management URL: %s, NRF Discovery URL: %s.", mm.MgmtPrefix, mm.DiscPrefix)
			nrfConnStatus = NRFConnLost
		}
		return
	}

	setNrfServerPrefix(mm.MgmtPrefix, mm.DiscPrefix)
	nrfConnStatus = NRFConnNormal
	log.Infof("NRF status set to avilable, NRF Management URL: %s, NRF Discovery URL: %s", mm.MgmtPrefix, mm.DiscPrefix)
}

var sendMonitorMessage = func(msg *monitorMessage) error {
	if msg == nil {
		return errors.New("message should not be nil")
	}
	msgBuf, err := json.Marshal(*msg)
	if err != nil {
		return errors.New("failed to marshal NRF monitor message, " + err.Error())
	}
	if monitorMSB == nil {
		return errors.New("message bus was not initialized")
	}
	err = monitorMSB.SendMessage(monitorTopic, string(msgBuf))
	if err != nil {
		return errors.New("failed to send message to message bus, " + err.Error())
	}
	log.Debugf("sendMonitorMessage success, msgBuf: %s", string(msgBuf))
	return nil
}

func monitorLoop(sec time.Duration) {
	monitorT := time.NewTicker(sec)
	if monitorT == nil {
		log.Errorf("failed to create timer for NRF Servers monitor")
		return
	}
	defer monitorT.Stop()

	monitorChan = make(chan bool)
	if monitorChan == nil {
		log.Errorf("failed to create channel for NRF Servers monitor")
		return
	}
	defer close(monitorChan)

	//Set forceReq as true for the first monitor message of SeconaryMonitor
	monitorHandler(true)
	for {
		select {
		case <-monitorT.C:
			monitorHandler(false)
		case <-monitorChan:
			monitorHandler(false)
		}
	}
}

func monitorHandler(forceReq bool) {
	mgmtPrefix, discPrefix := getNrfServerPrefix()
	if mgmtPrefix != "" && discPrefix != "" {
		//PrimaryMonitor check NRF connection every period, and reconnect NRF for next period if connection lost
		if monitorRole == PrimaryMonitor {
			if syncNrfStatus != SyncNrfStatusFailure {
				validateNrfServerConn(mgmtPrefix, discPrefix)
				return
			}
		} else {
			return
		}

	}

	switch monitorRole {
	case PrimaryMonitor:
		selectNrfServer()
	case SeconaryMonitor:
		mm := &monitorMessage{
			Role:       monitorRole,
			MgmtPrefix: "",
			DiscPrefix: "",
			ForceReq:   forceReq,
		}
		err := sendMonitorMessage(mm)
		if err != nil {
			log.Errorf("%s", err.Error())
		}
	default:
		log.Errorf("unknow NRF Server monitor role: %s", monitorRole)
	}
}

func getBaseNrfURLs(ep *structs.NrfServiceEndPoint, nrfServerProfile *structs.NrfServerProfile, isNFLevel bool) []string {
	var nrfURLs []string

	if !structs.ValidateNrfServiceEndPoint(ep, nrfServerProfile, isNFLevel) {
		if isNFLevel == false {
			log.Warnf("NRF Service level mandatary fields scheme, apiPrefix, versions, fqdn or ipEndPoint may be not implement in %+v", ep)
		} else {
			log.Warnf("NRF NF level ipv4address, scheme, apiPrefix, versions, fqdn or port in ipEndPoint may be not implement in %+v", *nrfServerProfile)
		}
		return nil
	}
	if isNFLevel != true {
		if ep.IPEndPoints != nil {
			for _, v := range ep.IPEndPoints {
				ipAddress := v.Ipv4Address

				if cm.IsEnableIpv6() {
					if v4Addr := common.ConvertIpv6ToIpv4(v.Ipv6Address); v4Addr != "" {
						ipAddress = v4Addr
					}
				}

				if ipAddress == "0.0.0.0" ||
					v.Port == 0 {
					continue
				}

				if ep.APIPrefix != "" {
					reqURL := ep.Scheme + "://" +
						ipAddress + ":" +
						strconv.Itoa(v.Port) + "/" +
						ep.APIPrefix + "/" +
						ep.ServiceName + "/" +
						ep.Versions[0].APIVersionInUrI + "/"
					nrfURLs = append(nrfURLs, reqURL)
				} else {
					reqURL := ep.Scheme + "://" +
						ipAddress + ":" +
						strconv.Itoa(v.Port) + "/" +
						ep.ServiceName + "/" +
						ep.Versions[0].APIVersionInUrI + "/"
					nrfURLs = append(nrfURLs, reqURL)
				}

			}
		}
		if ep.Fqdn != "" {
			if ep.APIPrefix != "" {
				reqURL := ep.Scheme + "://" +
					ep.Fqdn + "/" +
					ep.APIPrefix + "/" +
					ep.ServiceName + "/" +
					ep.Versions[0].APIVersionInUrI + "/"
				nrfURLs = append(nrfURLs, reqURL)
			} else {
				reqURL := ep.Scheme + "://" +
					ep.Fqdn + "/" +
					ep.ServiceName + "/" +
					ep.Versions[0].APIVersionInUrI + "/"
				nrfURLs = append(nrfURLs, reqURL)
			}

		}
	} else {
		//get NF level address
		if len(nrfServerProfile.NrfServerIpv4Address) == 0 && nrfServerProfile.NrfServerFqdn == "" {
			log.Warnf("nrf NF level IPv4address or FQDN  may be not implement in %+v", *nrfServerProfile)
			return nil
		}
		for _, ipAddress := range nrfServerProfile.NrfServerIpv4Address {
			if ep.IPEndPoints != nil && ipAddress != "0.0.0.0" {
				for _, v := range ep.IPEndPoints {
					if v.Port == 0 {
						continue
					}
					if ep.APIPrefix != "" {
						reqURL := ep.Scheme + "://" +
							ipAddress + ":" +
							strconv.Itoa(v.Port) + "/" +
							ep.APIPrefix + "/" +
							ep.ServiceName + "/" +
							ep.Versions[0].APIVersionInUrI + "/"
						nrfURLs = append(nrfURLs, reqURL)
					} else {
						reqURL := ep.Scheme + "://" +
							ipAddress + ":" +
							strconv.Itoa(v.Port) + "/" +
							ep.ServiceName + "/" +
							ep.Versions[0].APIVersionInUrI + "/"
						nrfURLs = append(nrfURLs, reqURL)
					}
				}
			}
		}

		if nrfServerProfile.NrfServerFqdn != "" {
			if ep.APIPrefix != "" {
				reqURL := ep.Scheme + "://" +
					nrfServerProfile.NrfServerFqdn + "/" +
					ep.APIPrefix + "/" +
					ep.ServiceName + "/" +
					ep.Versions[0].APIVersionInUrI + "/"
				nrfURLs = append(nrfURLs, reqURL)
			} else {
				reqURL := ep.Scheme + "://" +
					nrfServerProfile.NrfServerFqdn + "/" +
					ep.ServiceName + "/" +
					ep.Versions[0].APIVersionInUrI + "/"
				nrfURLs = append(nrfURLs, reqURL)
			}
		}
	}

	return nrfURLs
}

type nrfServerSlice []structs.NrfServerProfile

func (s nrfServerSlice) Len() int           { return len(s) }
func (s nrfServerSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s nrfServerSlice) Less(i, j int) bool { return s[i].NrfServerPriority < s[j].NrfServerPriority }

func selectNrfServer() {
	var nrfServers structs.NrfServerList
	if !structs.GetNrfServerList(&nrfServers) {
		log.Errorf("failed to get NrfServerList")
		return
	}

	if len(nrfServers.NrfServerProfileList) == 0 {
		log.Errorf("failed to get NrfServerProfileList")
		return
	}

	sort.Sort(nrfServerSlice(nrfServers.NrfServerProfileList))

	var currentSyncStatus = SyncNrfStatusNOK
	//fetch service level address
	for _, nrfServerProfile := range nrfServers.NrfServerProfileList {
		log.Infof("try to connect NRF Server %s under service level", nrfServerProfile.NrfServerProfileID)
		currentSyncStatus = connectNrfServer(&nrfServerProfile, false)
		if currentSyncStatus != SyncNrfStatusNOK {
			break
		}
	}

	//fetch NF level address when service level address can not connect
	if currentSyncStatus == SyncNrfStatusNOK {
		for _, nrfServerProfile := range nrfServers.NrfServerProfileList {
			log.Infof("try to connect NRF Server %s under NF level", nrfServerProfile.NrfServerProfileID)
			currentSyncStatus = connectNrfServer(&nrfServerProfile, true)
			if currentSyncStatus != SyncNrfStatusNOK {
				break
			}
		}
	}

	//Send NRF connection NOK status to SeconaryMonitor when status change from OK/Unknown to NOK.
	if monitorRole == PrimaryMonitor &&
		syncNrfStatus != SyncNrfStatusNOK &&
		currentSyncStatus == SyncNrfStatusNOK {

		mm := &monitorMessage{
			Role:       monitorRole,
			MgmtPrefix: "",
			DiscPrefix: "",
		}
		err := sendMonitorMessage(mm)
		if err != nil {
			log.Errorf("%s", err.Error())
			nrfStatusLock.Lock()
			defer nrfStatusLock.Unlock()
			syncNrfStatus = SyncNrfStatusFailure
			return
		}
	}
	nrfStatusLock.Lock()
	defer nrfStatusLock.Unlock()
	syncNrfStatus = currentSyncStatus
}

func connectNrfServer(nrfServerProfile *structs.NrfServerProfile, isNFLevel bool) SyncNrfStatus {
	var mgmtPrefix, discPrefix string
	var connected SyncNrfStatus = SyncNrfStatusNOK

	for _, nrfServerEndPoint := range nrfServerProfile.NrfServiceEndPoints {
		switch nrfServerEndPoint.ServiceName {
		case structs.NrfMgmtServiceName:
			if mgmtPrefix == "" {
				mgmtPrefix = connectNrfServerEndPoint(&nrfServerEndPoint, structs.NrfMgmtServiceName, nrfServerProfile, isNFLevel)
			}
		case structs.NrfDiscServiceName:
			if discPrefix == "" {
				discPrefix = connectNrfServerEndPoint(&nrfServerEndPoint, structs.NrfDiscServiceName, nrfServerProfile, isNFLevel)
			}
		default:
			log.Errorf("invalid serivceName %s", nrfServerEndPoint.ServiceName)
		}

		//update NRF server prefix
		if mgmtPrefix == "" ||
			discPrefix == "" {
			continue
		}
		setNrfServerPrefix(mgmtPrefix, discPrefix)
		connected = SyncNrfStatusOK

		//update NRF server prefix to secondary monitor
		mm := &monitorMessage{
			Role:       monitorRole,
			MgmtPrefix: mgmtPrefix,
			DiscPrefix: discPrefix,
		}
		err := sendMonitorMessage(mm)
		if err != nil {
			log.Errorf("%s", err.Error())
			connected = SyncNrfStatusFailure
		}
		break
	}

	return connected
}

func connectNrfServerEndPoint(nrfServerEndPoint *structs.NrfServiceEndPoint, serviceName string, nrfServerProfile *structs.NrfServerProfile, isNFLevel bool) string {
	baseNrfURL := ""
	baseNrfURLs := getBaseNrfURLs(nrfServerEndPoint, nrfServerProfile, isNFLevel)

	if len(baseNrfURLs) == 0 {
		if isNFLevel == false {
			log.Errorf("failed to connect NRF %s, related Service Level address info: %v", nrfServerEndPoint.ServiceName, baseNrfURLs)
		} else {
			log.Errorf("failed to connect NRF %s, related NF Level address info: %v", nrfServerEndPoint.ServiceName, baseNrfURLs)
		}
		return baseNrfURL
	}

	for _, nrfURL := range baseNrfURLs {
		_, valid := hb2NrfServer(nrfURL, serviceName)
		if !valid {
			continue
		}
		baseNrfURL = nrfURL
		log.Infof("Valid NRF %s: %s", nrfServerEndPoint.ServiceName, baseNrfURL)
		break
	}

	return baseNrfURL
}

var hb2NrfServer = func(url string, serviceName string) (*httpclient.HttpRespData, bool) {
	hdr := make(map[string]string)
	hdr["Content-Type"] = "application/json"
	var verifyUrl string
	switch serviceName {
	case structs.NrfMgmtServiceName:
		verifyUrl = url + "nf-instances?limit=1"
	case structs.NrfDiscServiceName:
		verifyUrl = url + "nf-instances?requester-nf-type=NRF&target-nf-type=NRF&target-nf-instance-id=test"
	default:
		log.Errorf("hb2NrfServer: invalid serivceName %s", serviceName)
		return nil, false
	}

	respData, err := HTTPDo("h2", "GET", verifyUrl, hdr, nil)
	if err != nil {
		log.Errorf("hb2NrfServer connect to NRF %s error: %s", serviceName, err.Error())
		return nil, false
	}
	if respData == nil {
		log.Errorf("hb2NrfServer connect to NRF %s, respData is nil", serviceName)
		return nil, false
	}
	if respData.StatusCode < 600 && respData.StatusCode >= 500 {
		log.Errorf("hb2NrfServer connect to NRF %s, StatusCode %v", serviceName, respData.StatusCode)
		return respData, false
	}
	return respData, true
}

func request2NrfServer(serviceName string, httpv, method, urlSuffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
	var url string

	mgmtPrefix, discPrefix := getNrfServerPrefix()
	if mgmtPrefix == "" ||
		discPrefix == "" {
		ResetNrfServerPrefix()
		return nil, errors.New("no available NRF Server was connected")
	}
	switch serviceName {
	case structs.NrfMgmtServiceName:
		url = mgmtPrefix + urlSuffix
	case structs.NrfDiscServiceName:
		url = discPrefix + urlSuffix
	}

	resp, err := HTTPDo(httpv, method, url, hdr, body)
	if err != nil {
		log.Errorf("failed to send request to NRF Server, %s", err.Error())
		ResetNrfServerPrefix()
	}

	return resp, err
}

func validateNrfServerConn(mgmtPrefix string, discPrefix string) bool {
	_, mgmtValid := hb2NrfServer(mgmtPrefix, structs.NrfMgmtServiceName)
	if !mgmtValid {
		log.Errorf("validateNrfServerConn: lost connection to NRF Management.")
		setNrfServerPrefix("", "")
		return false
	}
	_, discValid := hb2NrfServer(discPrefix, structs.NrfDiscServiceName)
	if !discValid {
		log.Errorf("validateNrfServerConn: lost connection to NRF Discovery, %v", discValid)
		setNrfServerPrefix("", "")
		return false
	}
	return true
}

//GetNRFConnStatus is for NRF connection status
func GetNRFConnStatus() NRFConnStatus {
	return nrfConnStatus
}

//sendResponseForForceReq is for PrimaryMonitor send dummy response message for forceReq
func sendResponseForForceReq(mm *monitorMessage) bool {
	if mm == nil {
		return false
	}
	if mm.ForceReq {
		respMm := &monitorMessage{
			Role:       monitorRole,
			MgmtPrefix: "",
			DiscPrefix: "",
		}
		err := sendMonitorMessage(respMm)
		if err != nil {
			log.Errorf("%s", err.Error())
			nrfStatusLock.Lock()
			defer nrfStatusLock.Unlock()
			syncNrfStatus = SyncNrfStatusFailure
		} else {
			nrfStatusLock.Lock()
			defer nrfStatusLock.Unlock()
			syncNrfStatus = SyncNrfStatusNOK
		}
		return true
	}
	return false

}
