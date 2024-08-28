package disc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/cmproxy"
	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/httpserver"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/pm"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/common/pkg/utils"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/fm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/k8sapiclient"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/app/disc/discutil"
	"gerrit.ericsson.se/udm/nrfagent_discovery/app/disc/nfrequester"
	"gerrit.ericsson.se/udm/nrfagent_discovery/app/disc/schema"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/worker"

	"github.com/buger/jsonparser"
)

const (
	httpContentTypeJSON        = "application/json"
	httpHeaderJSONPatchJSON    = "application/json-patch+json"
	httpContentTypeProblemJSON = "application/problem+json"
	httpResponseFormat         = `{"title": "%s"}`
)

const (
	//Log Message for HTTP Request
	requestLogFormat = `{"request":{"sequenceID":"%s", "URL":"%+v", "method":"%s", "description":%s}}`
	//Log Message for HTTP Response
	responseLogFormat = `{"response":{"sequenceID":"%s", "statusCode":%d, "description":%s}}`
)

const (
	//defaultValidityTime is about 100 years
	defaultValidityTime = 876576 * time.Hour
	//defaultTimeDelta is a time delta before timer fired
	defaultTimeDelta         = 5 * time.Second
	defaultTimeDeltaForSlave = 2 * time.Second
)

type DiscAgentRole int

const (
	AgentRoleUnknown DiscAgentRole = 0
	AgentRoleMaster  DiscAgentRole = 1
	AgentRoleSlave   DiscAgentRole = 2
)

var (
	retryTimeDuration = 2 * time.Second

	nfIsReady               = make(chan bool)
	needReadinessCheckMutex sync.Mutex
	needReadinessCheck      bool = true
	agentRoleQuitMonitor         = make(chan bool)

	agentRoleMonitorStarted bool          = false
	agentRole               DiscAgentRole = AgentRoleUnknown
)

var (
	cacheManager  *cache.CacheManager
	workerManager *worker.WorkerManager
)

//Setup invoked in server package to startup Discovery Agent
func Setup() {
	log.Debugf("Setup: Begin to setup...")
	util.PreComplieRegexp()
	httpserver.Routes = append(httpserver.Routes,
		httpserver.Route{
			Name:        "discoveryRequest",
			Method:      "GET",
			Pattern:     "/nrf-discovery-agent/v1/nf-instances",
			HandlerFunc: nfDiscoveryRequestHandler,
		},
		httpserver.Route{
			Name:        "subscriptionID",
			Method:      "GET",
			Pattern:     "/nrf-discovery-agent/v1/subscriptions",
			HandlerFunc: subscriptionRequestHandler,
		},
		httpserver.Route{
			Name:        "roamSubscriptionID",
			Method:      "GET",
			Pattern:     "/nrf-discovery-agent/v1/roam-subscriptions",
			HandlerFunc: roamSubscriptionRequestHandler,
		},
	)
	log.Infof("Discover Agent url : /nrf-discovery-agent/v1/nf-instances")

	cacheSetup()

	podUUID, err := utils.GetUUIDString()
	if err != nil {
		log.Warnf("Get local uuid failure")
		podUUID = "none"
	}
	common.SetSelfUUID(podUUID)
}

func startAgentRoleMonitor(monitorTimer int) {
	log.Debugf("Start agent role monitor, timer: %d second", monitorTimer)
	agentRoleMonitorStarted = true
	ticker := time.NewTicker(time.Second * time.Duration(monitorTimer))
	go func() {
		for {
			select {
			case <-ticker.C:
				curAgentRole := getAgentRole()
				if agentRole == AgentRoleSlave && curAgentRole == AgentRoleMaster {
					log.Infof("Agent role switch from slave to master")
					workerManager.StopAllWorker()
					workerManager.LaunchAllLeftTask()
					agentRole = AgentRoleMaster
				}
			case <-agentRoleQuitMonitor:
				ticker.Stop()
				agentRoleMonitorStarted = false
				log.Debugf("Ticker which is used to monitor agent role stops")
				return
			}
		}
	}()
}

// StopAgentRoleMonitor is to stop agent role monitor
func StopAgentRoleMonitor() {
	if agentRoleMonitorStarted {
		log.Debugf("Stop agent role monitor.")
		agentRoleQuitMonitor <- true
	}

	for {
		if !agentRoleMonitorStarted {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(500))
	}
}

func getReadinessCheckFlag() bool {
	return needReadinessCheck
}
func setReadinessCheckFlag(status bool) {
	needReadinessCheck = status
}

func getAgentRole() DiscAgentRole {
	var curAgentRole DiscAgentRole

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		curAgentRole = AgentRoleMaster
	} else {
		curAgentRole = AgentRoleSlave
	}

	return curAgentRole
}

// InitiateRun initialize cmproxy, fmproxy, message bus, and another common services
func InitiateRun() {
	// Electing leader
	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		log.Infof("InitiateRun: act as a leader")
		agentRole = AgentRoleMaster
	} else {
		log.Infof("InitiateRun: act as a slave")
		agentRole = AgentRoleSlave
	}

	cacheManager = cache.Instance()
	workerManager = worker.Instance()

	workModeMonitorTimer := 2
	worker.StartWorkModeMonitor(workModeMonitorTimer)

	agentRoleMonitorTimer := 2
	startAgentRoleMonitor(agentRoleMonitorTimer)

	// Initializ subscriptionInfoList data in ConfigMap
	loadConfigmapStorage()

	// Initialize cmproxy
	cmproxy.Init(os.Getenv("CM_URI_PREFIX"))
	// Register configurations
	cmproxy.RegisterConf(os.Getenv("CM_CONFNAME_PREFIX")+"-nrfagentlog", "ericsson-nrfagentlog:nrfagentlog",
		"nrf-agent-cmproxy", common.CmNrfAgentLogHandler, cmproxy.NtfFormatFull)
	cmproxy.RegisterConf("ericsson-nrf-agent", "ericsson-nrf-agent:nrf-agent",
		"nrf-agent-cmproxy", cmNrfAgentConfHandler, cmproxy.NtfFormatFull)
	// Running
	cmproxy.Run()

	//load schema
	err := schema.LoadDiscoverSchema()
	if err != nil {
		log.Error("InitiateRun: Load notify NF Profile schema failed")
		os.Exit(1)
	}

	// Initialize fmproxy
	fm.Init(consts.AlarmDiscServiceName)

	waitCMReady()

	// Initialize Message Bus
	if err := initMessageBus(); err != nil {
		log.Warnf("InitiateRun: Initial NRF Discovery Agent MessageBus failed")
	}

	// Initialize NRF Monitor
	e := client.InitializeMonitor(client.SeconaryMonitor)
	if e != nil {
		log.Errorf("failed to initialize NRF Server monitor")
	} else {
		log.Infof("initialize NRF Server monitor as secondary monitor")
	}

	// Initialize ConfigMap monitor
	go configmapMonitor(cm.Opts.FileNotifyDir, 5*time.Second, configmapTargetNfProfilesHandler)

	waitDISCAgentReady()

}

func waitCMReady() {
	for {
		var nrfServers structs.NrfServerList
		if !structs.GetNrfServerList(&nrfServers) {
			time.Sleep(retryTimeDuration)
			continue
		}
		if !structs.ValidateNrfServerList(&nrfServers) {
			time.Sleep(retryTimeDuration)
			continue
		}
		if !structs.ValidateStatusNotifIPEndPoint() {
			log.Errorf("waitCMReady: StatusNotifIPEndPoint of NRF Agent Notification is not configured CM service")
			time.Sleep(retryTimeDuration)
			continue
		}
		break
	}
}

func configmapTargetNfProfilesHandler(event, configurationName string, rawData []byte) {
	if strings.Contains(event, "CREATE") || strings.Contains(event, "WRITE") {
		var err error
		var targetNfProfilesData []byte
		targetNfProfilesData, _, _, err = jsonparser.Get(rawData, "targetNfProfiles")
		if err != nil {
			log.Errorf("Failed to run jsonparser.Get() targetNfProfiles, %s", err.Error())
			return
		}
		log.Debugf("targetNfProfiles %s", string(targetNfProfilesData))

		var targetNfProfiles []structs.TargetNfProfile
		err = json.Unmarshal(targetNfProfilesData, &targetNfProfiles)
		if err != nil {
			log.Errorf("Unmarshal targetNfProfiles data failure, Error:%s", err.Error())
			return
		}

		for _, targetNfProfile := range targetNfProfiles {
			var targetNf structs.TargetNf
			targetNf.RequesterNfType = targetNfProfile.RequesterNfType
			targetNf.TargetNfType = targetNfProfile.TargetNfType
			targetNf.TargetServiceNames = targetNfProfile.TargetServiceNames
			targetNf.NotifCondition = targetNfProfile.NotifCondition
			targetNf.SubscriptionValidityTime = targetNfProfile.SubscriptionValidityTime

			_, existed := cacheManager.GetTargetNf(targetNf.RequesterNfType, targetNf.TargetNfType)
			if existed {
				continue
			}

			cacheManager.InitCache(targetNf.RequesterNfType, targetNf.TargetNfType)

			cacheManager.SetTargetNf(targetNf.RequesterNfType, targetNf)
			log.Debugf("%s TargetNf is %+v", targetNf.RequesterNfType, targetNf)
		}
	}

	cmTargetNfProfilesHandler(event, configurationName, cmproxy.NtfFormatFull, rawData)

}

func waitDISCAgentReady() {
	go discInitForReadiness()
	for {
		select {
		case <-nfIsReady:
			log.Infof("NRF Discovery Agent is ready to work")
			return
		default:
			{
				log.Infof("NRF Discovery Agent is not ready, wait for it to be ready")
				time.Sleep(retryTimeDuration)
			}
		}
	}
}

func discInitForReadiness() {
	nfInfos := loopFetchRequesterNfInfo()

	if nfInfos == nil {
		needReadinessCheckMutex.Lock()
		defer needReadinessCheckMutex.Unlock()
		if getReadinessCheckFlag() == true {
			nfIsReady <- true
			log.Debug("start set needReadinessCheck false")
			setReadinessCheckFlag(false)
			log.Debug("finish set needReadinessCheck false")
		}

		return
	}

	for _, nfInfo := range nfInfos {
		go workerManager.PrepareDiscoveryAgent(nfInfo)
	}

	if workerManager.WaitAgentReady() {
		needReadinessCheckMutex.Lock()
		defer needReadinessCheckMutex.Unlock()
		if getReadinessCheckFlag() == true {
			nfIsReady <- true
			//needReadinessCheck used for the firstly deploy or reboot
			log.Debug("start set needReadinessCheck false")
			setReadinessCheckFlag(false)
			log.Debug("finish set needReadinessCheck false")
		}
	}
}

func loopFetchRequesterNfInfo() []structs.NfInfoForRegDisc {
	nfInfos, ok := fetchRequesterNfInfo()
	for !ok {
		log.Warnf("Fetch nfInfo failure and retry to fetch nfInfo from NRF Register Agent %v seconds later", retryTimeDuration)
		time.Sleep(retryTimeDuration)
		nfInfos, ok = fetchRequesterNfInfo()
	}

	if nfInfos == nil {
		log.Warnf("Fetch All nfInfo from NRF Register Agent is empty")
		return nil
	}

	log.Info("Fetch all nfInfo from NRF Register Agent success")
	return nfInfos
}

func fetchRequesterNfInfo() ([]structs.NfInfoForRegDisc, bool) {
	nrfRegAgentName := os.Getenv("NAME_NRFAGENT_REG")
	nrfRegAentPort := os.Getenv("NAME_NRFAGENT_REG_PORT")
	urlReg := "http://" + nrfRegAgentName + ":" + nrfRegAentPort + "/nrf-register-agent/v1/dumpNfInfo"
	log.Infof("request nfInfo url is %s", urlReg)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDo("h2", "GET", urlReg, hdr, nil)
	if err != nil {
		log.Errorf("Failed to send request to NRF Register Agent, %s", err.Error())
		return nil, false
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("fetch all fqdn failed, StatusCode(%d)", resp.StatusCode)
		return nil, false
	}
	if len(resp.Body) == 0 {
		log.Infof("No NF instance registered in Register Agent")
		return nil, true
	}
	log.Infof("Response body is %+v", string(resp.Body))
	var nfInfos []structs.NfInfoForRegDisc
	err = json.Unmarshal(resp.Body, &nfInfos)
	if err != nil {
		log.Errorf("Unmarshal nfInfoSlice error, Error:%s", err.Error())
		return nil, false
	}
	log.Debugf("fqdn from Register Agent is %+v", nfInfos)
	return nfInfos, true
}

/*
func unsubscribeByNfType(requesterNfType string) {
	if worker.IsKeepCacheMode() {
		log.Info("keep cache mode, not send message to NRF.")
		return
	}
	targetNfs, ok := cacheManager.GetTargetNfs(requesterNfType)
	if !ok {
		log.Errorf("Failed to get targetNfProfiles for nfType[%s], please check configmap status", requesterNfType)
		return
	}

	for _, targetNf := range targetNfs {
		unsubscribeSubscription(requesterNfType, targetNf.TargetNfType)
	}
}
*/
/*
func unsubscribeSubscription(requesterNfType, targetNfType string) {
	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON

	subscriptionIDURLs, ok := cacheManager.GetSubscriptionIDs(requesterNfType, targetNfType)
	for ok && len(subscriptionIDURLs) > 0 {
		subscriptionID := subscriptionIDURLs[0]
		log.Infof("subscriptionID of %s: %+v", requesterNfType, subscriptionID)

		if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
			resp, err := client.HTTPDoToNrfMgmt("h2", "DELETE", subscriptionID, hdr, bytes.NewBuffer([]byte("")))
			if err != nil {
				log.Errorf("failed to send DELETE subscription request to NRF, %s", err.Error())
			} else {
				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
					log.Errorf("Failed to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
				} else {
					log.Infof("Success to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
				}
			}
		} else {
			log.Infof("Slaver discovery agent has no need to send unsubscription request to NRF")
		}

		cacheManager.DelSubscriptionInfo(requesterNfType, targetNfType, subscriptionID)
		cacheManager.DelSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
		cacheManager.UpdateSubscriptionStorage()

		subscriptionIDURLs, ok = cacheManager.GetSubscriptionIDs(requesterNfType, targetNfType)
	}

	roamsubscriptionIDURLs, ok := cacheManager.GetRoamingSubscriptionIDs(requesterNfType, targetNfType)
	for ok && len(roamsubscriptionIDURLs) > 0 {
		roamsubscriptionID := roamsubscriptionIDURLs[0]
		log.Infof("roaming subscriptionID of %s: %+v", requesterNfType, roamsubscriptionID)

		if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
			resp, err := client.HTTPDoToNrfMgmt("h2", "DELETE", roamsubscriptionID, hdr, bytes.NewBuffer([]byte("")))
			if err != nil {
				log.Errorf("failed to send DELETE roaming subscription request to NRF, %s", err.Error())
			} else {
				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
					log.Errorf("Failed to DELETE roaming subscription(%s), StatusCode(%d)", roamsubscriptionID, resp.StatusCode)
				} else {
					log.Infof("Success to DELETE roaming subscription(%s), StatusCode(%d)", roamsubscriptionID, resp.StatusCode)
				}
			}
		} else {
			log.Infof("Slaver discovery agent has no need to send unsubscription request to NRF")
		}

		cacheManager.DelRoamingSubscriptionInfo(requesterNfType, targetNfType, roamsubscriptionID)
		cacheManager.DelRoamingSubscriptionMonitor(requesterNfType, targetNfType, roamsubscriptionID)

		roamsubscriptionIDURLs, ok = cacheManager.GetRoamingSubscriptionIDs(requesterNfType, targetNfType)
	}
}
*/

func handleDiscoveryRequest(targetNf *structs.TargetNf, nfInstanceID string) {
	requestOptions := ""
	if nfInstanceID != "" {
		requestOptions = consts.SearchDataTargetInstID + "=" + nfInstanceID
	}

	fqdn, exists := cacheManager.GetRequesterFqdn(targetNf.RequesterNfType)
	if !exists {
		log.Warnf("%s nf instance was deregistered", targetNf.RequesterNfType)
		return
	}

	query := util.GetDiscoveryRequestURL(targetNf, requestOptions, fqdn)
	if query == "" {
		log.Errorf("Failed to get Discovery request URL")
		return
	}

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		log.Debugf("Master Agent KeepCache Mode status %t", worker.IsKeepCacheMode())
		if worker.IsKeepCacheMode() {
			log.Info("Master Agent is keep cache mode, not send message to NRF.")
			return
		}
		if _, err := handleDiscoveryRequestToNrf(targetNf, query); err != nil {
			log.Errorf("Failed to query %s,%+v from NRF, %s",
				targetNf.TargetNfType, targetNf.TargetServiceNames, err.Error())
		}
	} else {
		if _, err := handleDiscoveryRequestToMaster(targetNf, query); err != nil {
			log.Errorf("Failed to query %s,%+v from master agent, %s",
				targetNf.TargetNfType, targetNf.TargetServiceNames, err.Error())
		}
	}
}

func handleDiscoveryRequestToMaster(targetNf *structs.TargetNf, query string) (*httpclient.HttpRespData, error) {
	// Waiting 1s for master agent fetching the profile from NRF
	// avoid to forward the request repeatly
	time.Sleep(time.Second)

	log.Infof("Try to fetch %s %+v profile from master NRF Discovery Agent",
		targetNf.TargetNfType, targetNf.TargetServiceNames)

	masterURL := util.GetLeaderDiscURL()
	if masterURL == "" {
		log.Errorf("Failed to get master URL")
		return nil, errors.New("NRF Discovery master agent URI not found")
	}
	log.Debugf("Master discovery agent url \"%s\"", masterURL+query)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDo("h2", "GET", masterURL+query, hdr, nil)
	if err != nil {
		log.Errorf("Failed to send discovery request to master discovery agent, %s", err.Error())
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode <= http.StatusUnavailableForLegalReasons {
		log.Infof("NRF Master response code is: %+v", resp.StatusCode)
		cacheManager.SetCacheStatus(targetNf.RequesterNfType, targetNf.TargetNfType, true)
		return nil, errors.New(string(resp.Body))
	}

	err = handleDiscoveryResponse(targetNf, resp)
	if err != nil {
		fmRaiseNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
		return resp, err
	}
	fmClearNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
	return resp, nil
}

func handleDiscoveryRequestToNrf(targetNf *structs.TargetNf, query string) (*httpclient.HttpRespData, error) {
	pm.Inc(consts.NrfDiscoveryRequestsTotal)

	log.Infof("Try to fetch %s %+v profile from NRF Discovery", targetNf.TargetNfType, targetNf.TargetServiceNames)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDoToNrfDisc("h2", "GET", query, hdr, nil)
	if err != nil {
		log.Errorf("failed to send discovery request to NRF, %s", err.Error())
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode <= http.StatusUnavailableForLegalReasons {
		log.Infof("NRF response code is: %+v", resp.StatusCode)
		cacheManager.SetCacheStatus(targetNf.RequesterNfType, targetNf.TargetNfType, true)
		return nil, errors.New(string(resp.Body))
	}

	err = handleDiscoveryResponse(targetNf, resp)
	if err != nil {
		fmRaiseNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
		return nil, err
	}
	fmClearNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
	util.PushMessageToMSB(targetNf.RequesterNfType, targetNf.TargetNfType, "", consts.NFEventDiscResult, resp.Body)
	return resp, nil
}

func handleDiscoveryResponse(targetNf *structs.TargetNf, resp *httpclient.HttpRespData) error {
	pmNrfDiscoveryResponses(resp.StatusCode)

	//log.Infof("NRF response: %+v", resp.SimpleString())
	if resp.StatusCode != http.StatusOK {
		log.Errorf("NRF response: %+v", string(resp.Body))
		return errors.New(string(resp.Body))
	}
	cacheManager.SetCacheStatus(targetNf.RequesterNfType, targetNf.TargetNfType, true)
	var err error
	resp.Body, err = common.ConvertIpv6ToIpv4InSearchResult(resp.Body, cm.IsEnableConvertIpv6ToIpv4())
	if err != nil {
		log.Errorf("Failed to convert Ipv6Address to Ipv4Address in NF profile, %s", err.Error())
		return err
	}

	log.Infof("NRF response: %+v", string(resp.Body))
	nfInstances, validityPeriod, ok := cache.SpliteSeachResult(resp.Body)
	if !ok {
		return fmt.Errorf("invalid SearchResult in NRF response, %+v", string(resp.Body))
	}
	for _, nfProfile := range nfInstances {
		cacheManager.CachedWithTTL(targetNf.RequesterNfType, targetNf.TargetNfType, nfProfile, validityPeriod, false)
	}

	return nil
}

func handleDiscoveryFailure(rw http.ResponseWriter, req *http.Request, logcontent *log.LogStruct, statusCode int, body string) {
	log.Debugf(consts.REQUEST_LOG_FORMAT, logcontent.SequenceId, req.URL, req.Method, logcontent.RequestDescription)
	DiscoveryResponseHander(rw, req, statusCode, body)
	log.Errorf(consts.RESPONSE_LOG_FORMAT, logcontent.SequenceId, statusCode, logcontent.ResponseDescription)
}

func handleDiscoverySuccess(rw http.ResponseWriter, req *http.Request, logcontent *log.LogStruct, statusCode int, validityPeriod int64, body string) {
	log.Debugf(consts.REQUEST_LOG_FORMAT, logcontent.SequenceId, req.URL, req.Method, logcontent.RequestDescription)
	//pm.Inc(constvalue.NfDiscoverySuccessTotal)
	if validityPeriod == 0 {
		rw.Header().Set("Cache-Control", consts.SearchDataCacheControlNoCache)
	} else {
		rw.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(validityPeriod, 10))
	}
	DiscoveryResponseHander(rw, req, statusCode, body)
	log.Debugf(consts.RESPONSE_LOG_FORMAT, logcontent.SequenceId, statusCode, logcontent.ResponseDescription)
}

// DiscoveryResponseHander handle response for  discovery
func DiscoveryResponseHander(rw http.ResponseWriter, req *http.Request, statuscode int, body string) {
	pm.Inc(consts.NfDiscoveryResponsesTotal)

	rw.WriteHeader(statuscode)
	if statuscode != http.StatusOK {
		rw.Header().Set("Content-Type", httpContentTypeProblemJSON)
		if body != "" {
			_, err := rw.Write([]byte(body))
			if err != nil {
				log.Warnf("%v", err)
			}
		}
	} else {
		rw.Header().Set("Content-Type", httpContentTypeJSON)
		//rw.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(validityPeriod, 10))
		if body != "" {
			_, err := rw.Write([]byte(body))
			if err != nil {
				log.Warnf("%v", err)
			}
		}
	}
}

func nfDiscoveryRequestHandler(rw http.ResponseWriter, req *http.Request) {
	startedAt := time.Now()
	defer func() {
		pm.Observe(float64(time.Since(startedAt))/float64(time.Second), consts.NfRequestDuration, consts.NfDiscovery)
	}()

	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}
	pm.Inc(consts.NfDiscoveryRequestsTotal)

	log.Debugf("NF Discovery Request comes, remote[%s] query[%s]", req.RemoteAddr, req.URL.RawQuery)
	var problemDetails *problemdetails.ProblemDetails
	oriqueryForm, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		errorInfo := fmt.Sprintf("Parse URL err: %s", err.Error())
		problemDetails = &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.RequestDescription = fmt.Sprintf(`{"service-names":%v, "target-nf-type":"%s", "requester-nf-type":"%s"}`, "[]", "", "")
		logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, http.StatusBadRequest, problemDetails.ToString())
		return
	}
	for k, v := range oriqueryForm {
		log.Debugf("oriqueryForm Key: %s", k)
		for i, value := range v {
			log.Debugf("oriqueryForm Value[%d]: %s", i, value)
		}
	}
	var queryForm nfrequester.SearchParameterData
	queryForm.InitMember(oriqueryForm)
	problem := queryForm.ValidateNRFDiscovery()
	if problem != nil {
		logcontent.RequestDescription = fmt.Sprintf(`{"service-names":%v, "target-nf-type":"%s", "requester-nf-type":"%s"}`, "[]", "", "")
		logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, problem.ToString())
		handleDiscoveryFailure(rw, req, logcontent, http.StatusBadRequest, problem.ToString())
		return
	}

	requesterNfType := queryForm.FetchRequesterNfTypeParameter()
	targetNfType := queryForm.FetchTargetNfTypeParameter()

	targetNf, exists := cacheManager.GetTargetNf(requesterNfType, targetNfType)
	if !exists {
		errorInfo := fmt.Sprintf("Do not deploy targetNfProfile configmap for requesterNfType:%s,targetNfType:%s", requesterNfType, targetNfType)
		problemDetails = &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.RequestDescription = fmt.Sprintf(`{"service-names":%v, "target-nf-type":"%s", "requester-nf-type":"%s"}`, "[]", "", "")
		logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, http.StatusForbidden, problemDetails.ToString())
		return
	}

	_, exists = cacheManager.GetRequesterFqdn(requesterNfType)
	if !exists {
		errorInfo := fmt.Sprintf("RequesterNfType:%s have not any NF instance online", requesterNfType)
		problemDetails = &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.RequestDescription = fmt.Sprintf(`{"service-names":%v, "target-nf-type":"%s", "requester-nf-type":"%s"}`, "[]", "", "")
		logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, http.StatusForbidden, problemDetails.ToString())
		return
	}

	serviceNames := queryForm.FetchServcieNamesParameter()
	err = discutil.ServiceNameScopeVerify(serviceNames, &targetNf)
	if err != nil {
		errorInfo := err.Error()
		problemDetails = &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.RequestDescription = fmt.Sprintf(`{"service-names":%v, "target-nf-type":"%s", "requester-nf-type":"%s"}`, "[]", "", "")
		logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, http.StatusForbidden, problemDetails.ToString())
		return
	}

	plmns := queryForm.FetchPlmnsParameter()
	err = discutil.RejectVerify(oriqueryForm, plmns)
	if err != nil {
		errorInfo := err.Error()
		problemDetails = &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.RequestDescription = fmt.Sprintf(`{"service-names":%v, "target-nf-type":"%s", "requester-nf-type":"%s"}`, "[]", "", "")
		logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, http.StatusForbidden, problemDetails.ToString())
		return
	}

	var ok bool
	ok = discutil.ForwardVerify(oriqueryForm)
	if ok {
		if keepCacheModeSendResponse(rw, req) {
			return
		}
		rest := proxy(rw, req) // just forward
		if rest {
			log.Info("NRF-Agent-Disc proxy to NRF and response to NF success")
		} else {
			log.Info("NRF-Agent-Disc proxy to NRF and response to NF failure")
		}

		return
	}

	searchParameter := cache.SearchParameter{}
	queryForm.CacheSearchParameterInjection(&searchParameter)

	var content []byte
	var isRoamQuery bool
	ok = discutil.RoamingVerify(oriqueryForm)
	if ok {
		isRoamQuery = true
		log.Infof("Nrf-agent-disc will search in roaming cache")
		content, ok = cacheManager.SearchRoamingCache(requesterNfType, targetNfType, &searchParameter, worker.IsKeepCacheMode())
	} else {
		isRoamQuery = false
		log.Infof("Nrf-agent-disc will search in home cache")
		content, ok = cacheManager.Search(requesterNfType, targetNfType, &searchParameter, worker.IsKeepCacheMode())
	}

	if !ok {
		if keepCacheModeSendResponse(rw, req) {
			return
		}
		log.Info("NF profile not found in cache, forward to NRF Discovery")
		respRawData, ok := syncQueryData(rw, req, &targetNf) // sync data and then apply filter
		if !ok {
			log.Info("Query to NRF failure or NRF response failure")
			return
		}

		validityPeriod, err := util.GetValidityPeriod(respRawData)
		if err != nil {
			errorInfo := err.Error()
			problemDetails := &problemdetails.ProblemDetails{
				Title: errorInfo,
			}
			logcontent.ResponseDescription = errorInfo
			handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
			return
		}

		var respData []byte
		respData, err = applyFilter(respRawData, &searchParameter)
		if err != nil {
			errorInfo := err.Error()
			problemDetails := &problemdetails.ProblemDetails{
				Title: errorInfo,
			}
			logcontent.ResponseDescription = errorInfo
			handleDiscoveryFailure(rw, req, logcontent, http.StatusNotFound, problemDetails.ToString())
			return
		}

		//allowCached := discutil.ResponseAllowCache()
		if !isRoamQuery {
			go cacheMessage(requesterNfType, targetNfType, respRawData, false)
			go util.PushMessageToMSB(requesterNfType, targetNfType, "", consts.NFEventDiscResult, respRawData)
		} else {
			roamTargetPlmnID := queryForm.FetchRoamPlmnIDParameter()
			roamRespRawData := discutil.BuildRoamMessage(requesterNfType, targetNfType, respRawData, roamTargetPlmnID)

			go cacheMessage(requesterNfType, targetNfType, roamRespRawData, true)
			go discutil.SubscribeRoamNfProfiles(respRawData, requesterNfType, targetNfType, &roamTargetPlmnID, validityPeriod)
			go util.PushMessageToMSB(requesterNfType, targetNfType, "", consts.NFEventDiscResult, roamRespRawData)
		}

		handleDiscoverySuccess(rw, req, logcontent, http.StatusOK, validityPeriod, string(respData))

		return
	}

	log.Info("NRF-Agent-Disc search in cache success")
	validityPeriod, err := util.GetValidityPeriod(content)
	if err != nil {
		errorInfo := err.Error()
		problemDetails := &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.ResponseDescription = errorInfo
		handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
		return
	}

	handleDiscoverySuccess(rw, req, logcontent, http.StatusOK, validityPeriod, string(content))
}

func proxy(rw http.ResponseWriter, req *http.Request) bool {
	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}
	//var problemDetails *problemdetails.ProblemDetails

	rawQuery := req.URL.RawQuery
	log.Infof("rawQuery=%s", rawQuery)

	//hdr = req.Header
	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	query := "nf-instances?" + rawQuery
	resp, err := client.HTTPDoToNrfDisc("h2", "GET", query, hdr, nil)
	if err != nil {
		errorInfo := err.Error()
		problemDetails := &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("Discovery request to NRF fail, err:%s", errorInfo),
		}
		logcontent.ResponseDescription = fmt.Sprintf(`"Discovery request to NRF fail. %s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
		return false
	}

	//pm.Inc(consts.NrfDiscoveryRequestsTotal)
	//pmNrfDiscoveryResponses(resp.StatusCode)

	log.Debugf("NRF response: %s", resp.SimpleString())

	respData := resp.Body
	if resp.StatusCode != http.StatusOK {
		errorInfo := string(respData)
		problemDetails := &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.ResponseDescription = fmt.Sprintf(`"NRF response body: %s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, resp.StatusCode, problemDetails.ToString())
		return false
	}
	/*
		validityPeriod, err := util.GetValidityPeriod(respData)
		if err != nil {
			errorInfo := err.Error()
			problemDetails := &problemdetails.ProblemDetails{
				Title: errorInfo,
			}
			logcontent.ResponseDescription = errorInfo
			handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
			return false
		}
	*/
	profiles, _, _, err := jsonparser.Get(respData, "nfInstances")
	if err != nil {
		log.Errorf("Get nfInstances from NRF reponse fail, err:%s\n", err.Error())
		errorInfo := err.Error()
		problemDetails := &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.ResponseDescription = errorInfo
		handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
		return false
	}

	noCacheResp := fmt.Sprintf("{\"validityPeriod\":%d,\"nfInstances\":%s}", 0, string(profiles))

	logcontent.ResponseDescription = "NRF response success"
	handleDiscoverySuccess(rw, req, logcontent, http.StatusOK, 0, noCacheResp)

	return true
}

func syncQueryData(rw http.ResponseWriter, req *http.Request, targetNf *structs.TargetNf) ([]byte, bool) {
	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}
	//var problemDetails *problemdetails.ProblemDetails

	//1 refactor requester
	v := req.URL.Query()
	v.Del(consts.SearchDataRequesterNfType)
	v.Del(consts.SearchDataTargetNfType)
	v.Del(consts.SearchDataServiceName)
	v.Del(consts.SearchDataRequesterNFInstFQDN)
	// Remove the options impact NF services
	v.Del(consts.SearchDataSupportedFeatures)

	fqdn, exists := cacheManager.GetRequesterFqdn(targetNf.RequesterNfType)
	if !exists {
		log.Warnf("%s nf instance was deregistered", targetNf.RequesterNfType)
		return nil, false
	}

	query := util.GetDiscoveryRequestURL(targetNf, v.Encode(), fqdn)
	if query == "" {
		errorInfo := "NRF Discovery URI not found"
		problemDetails := &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("%s", errorInfo),
		}
		logcontent.ResponseDescription = fmt.Sprintf(`"Get NRF Discovery URI failed. %s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
		return nil, false
	}

	//2 query to NRF
	pm.Inc(consts.NrfDiscoveryRequestsTotal)

	log.Infof("Try to fetch %s %v profile from NRF Discovery", targetNf.TargetNfType, targetNf.TargetServiceNames)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDoToNrfDisc("h2", "GET", query, hdr, nil)
	if err != nil {
		errorInfo := err.Error()
		problemDetails := &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("Discovery request to NRF fail, err:%s", errorInfo),
		}
		logcontent.ResponseDescription = fmt.Sprintf(`"Discovery request to NRF fail. %s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
		return nil, false
	}

	//3 handler response, handleDiscoveryResponse start////////////
	pmNrfDiscoveryResponses(resp.StatusCode)

	log.Infof("NRF response: %s", resp.SimpleString())

	respBody := resp.Body
	if resp.StatusCode != http.StatusOK {
		errorInfo := string(respBody)
		problemDetails := &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.ResponseDescription = fmt.Sprintf(`"NRF response body: %s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, resp.StatusCode, problemDetails.ToString())
		return nil, false
	}

	//why convert ipv6 to ipv4 before return
	newBody, err := common.ConvertIpv6ToIpv4InSearchResult(respBody, cm.IsEnableConvertIpv6ToIpv4())
	if err != nil {
		errorInfo := fmt.Sprintf("Convert nfProfile ipv6 to ipv4 failed, %s", err.Error())
		problemDetails := &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.ResponseDescription = errorInfo
		handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
		return nil, false
	}

	return newBody, true

	/*
		if err != nil {
			fmRaiseNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
			return nil, err
		}
		fmClearNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)

		/////push to messageBus/////
		pushMessageToMSB(targetNf.RequesterNfType, targetNf.TargetNfType, "", consts.NFEventDiscResult, resp.Body)
	*/
}

func cacheMessage(requesterNfType string, targetNfType string, content []byte, isRoam bool) {
	validityPeriod, err := util.GetValidityPeriod(content)
	if err != nil {
		log.Errorf("Get validityPeriod failed, err:%s", err.Error())
		return
	}

	nfInstances, err := util.GetNfInstances(content)
	if err != nil {
		log.Errorf("Get nfInstances failed, err:%s", err.Error())
		return
	}

	for _, nfProfile := range nfInstances {
		cacheManager.CachedWithTTL(requesterNfType, targetNfType, nfProfile, uint(validityPeriod), isRoam)
	}
}

/*
func cacheRoamMessage(requesterNfType string, targetNfType string, content []byte, roamTargetPlmnID structs.PlmnID) {
	validityPeriod, err := util.GetValidityPeriod(content)
	if err != nil {
		log.Errorf("Get validityPeriod failed, err:%s", err.Error())
		return
	}

	nfInstances, err := util.GetNfInstances(content)
	if err != nil {
		log.Errorf("Get nfInstances failed, err:%s", err.Error())
		return
	}

	for _, nfProfile := range nfInstances {
		nfInstanceID := util.GetNfInstanceID(nfProfile)
		if nfInstanceID == "" {
			log.Error("nfProfile less nfInstanceID")
			continue
		}

		ok := cacheManager.ProbeRoam(requesterNfType, targetNfType, nfInstanceID)
		if ok {
			oldProfile := cacheManager.GetProfileByID(requesterNfType, targetNfType, nfInstanceID, true)
			if oldProfile == nil {
				continue
			}
			newNfProfile := discutil.AppendPlmnInfo(nfProfile, oldProfile, roamTargetPlmnID)
			cacheManager.CachedWithTTL(requesterNfType, targetNfType, newNfProfile, uint(validityPeriod), true)
		} else {
			missPlmn := discutil.PlmnMissProber(nfProfile)
			if missPlmn {
				newNfProfile := discutil.AddPlmnInfo(nfProfile, roamTargetPlmnID)
				cacheManager.CachedWithTTL(requesterNfType, targetNfType, newNfProfile, uint(validityPeriod), true)
			} else {
				cacheManager.CachedWithTTL(requesterNfType, targetNfType, nfProfile, uint(validityPeriod), true)
			}
		}
	}
}
*/
func applyFilter(rawContent []byte, searchParameter *cache.SearchParameter) ([]byte, error) {
	if !searchParameter.SearchServiceName() &&
		!searchParameter.SearchSupportedFeatures() {
		return rawContent, nil
	}

	var searchResult structs.SearchResult
	if err := json.Unmarshal(rawContent, &searchResult); err != nil {
		log.Errorf("applyFilter: failed to unmarshal SearchResult %+v, %s", string(rawContent), err.Error())
		return nil, fmt.Errorf("NRF response an invalid SearchResult, and requested NF profile not found in NRF")
	}
	//If NRF response 200 with no NfInstances, agent transfer it to NF.
	if len(searchResult.NfInstances) == 0 {
		log.Debugf("applyFilter: no NfInstances for NRF response.")
		return rawContent, nil
	}
	for i := 0; i < len(searchResult.NfInstances); {
		if !cache.NfServiceFilter(&searchResult.NfInstances[i], searchParameter) {
			searchResult.NfInstances = append(searchResult.NfInstances[:i], searchResult.NfInstances[i+1:]...)
		} else {
			i++
		}
	}
	if len(searchResult.NfInstances) == 0 {
		return nil, fmt.Errorf("requested NF profile not found in NRF")
	}
	content, err := json.Marshal(searchResult)
	if err != nil {
		log.Errorf("fixSearchResultFromNrf: failed to marshal SearchResult %+v, %s", searchResult, err.Error())
		return nil, fmt.Errorf("NRF response an invalid SearchResult, and requested NF profile not found in NRF")
	}
	return content, nil
}

//when start, master reponse the ttl delay 3 seconds compare with master, master delay 5 seconds compare with MGMT
func subscriptionRequestHandler(rw http.ResponseWriter, req *http.Request) {
	requesterNfType := req.FormValue(consts.SearchDataRequesterNfType)
	targetNfType := req.FormValue(consts.SearchDataTargetNfType)
	serviceName := req.FormValue(consts.SearchDataServiceName)

	log.Infof("ENTRY FROM Slave node(%s): query subscriptionInfo of requesterNfType:%s,targetNfType:%s,serviceNames:%s", req.RemoteAddr, requesterNfType, targetNfType, serviceName)

	//masterCacheStatus will wait for cache build during new deployment, cache is ok, then response the subscriptionInfo
	masterCacheStatus := cacheManager.GetCacheStatus(requesterNfType, targetNfType)
	if masterCacheStatus == false {
		log.Infof("Master cache is not ready for nfType[%s], no subsciptionInfo in cache", requesterNfType)
		rw.Header().Set("Content-Type", httpContentTypeJSON)
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	subscriptionInfo, ok := cacheManager.GetServiceSubscriptionInfo(requesterNfType, targetNfType, serviceName)
	if !ok {
		log.Infof("Master cache no such subscriptionInfo for requesterNfTpe[%s], targetNfType[%s], serviceName[%s]", requesterNfType, targetNfType, serviceName)
		rw.Header().Set("Content-Type", httpContentTypeJSON)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	validityTime := subscriptionInfo.ValidityTime.Add(-defaultTimeDeltaForSlave)
	subscriptionInfo.ValidityTime = validityTime

	subscriptionInfoData, err := json.Marshal(subscriptionInfo)
	if err != nil {
		log.Errorf("Marshal subscriptionInfo fail, err:%s", err.Error())
		rw.Header().Set("Content-Type", httpContentTypeJSON)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", httpContentTypeJSON)
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(subscriptionInfoData)
	if err != nil {
		log.Warnf("%v", err)
	}

	log.Infof("NRF-Discovery-Agent active response subscriptionInfo:%v", subscriptionInfo)
}

func roamSubscriptionRequestHandler(rw http.ResponseWriter, req *http.Request) {
	requesterNfType := req.FormValue(consts.SearchDataRequesterNfType)
	targetNfType := req.FormValue(consts.SearchDataTargetNfType)
	nfInstanceID := req.FormValue(consts.SearchDataTargetInstID)

	log.Infof("ENTRY FROM Slave node(%s): query subscriptionInfo of requesterNfType:%s,targetNfType:%s,nfInstanceID:%s", req.RemoteAddr, requesterNfType, targetNfType, nfInstanceID)

	//masterCacheStatus will wait for cache build during new deployment, cache is ok, then response the subscriptionInfo
	masterCacheStatus := cacheManager.GetCacheStatus(requesterNfType, targetNfType)
	if masterCacheStatus == false {
		log.Infof("Master cache is not ready for nfType[%s], no subsciptionInfo in cache", requesterNfType)
		rw.Header().Set("Content-Type", httpContentTypeJSON)
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	subscriptionInfo, ok := cacheManager.GetNfProfileSubscriptionInfo(requesterNfType, targetNfType, nfInstanceID)
	if !ok {
		log.Infof("Master cache no such subscriptionInfo for requesterNfTpe[%s], targetNfType[%s], nfInstanceID[%s]", requesterNfType, targetNfType, nfInstanceID)
		rw.Header().Set("Content-Type", httpContentTypeJSON)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	validityTime := subscriptionInfo.ValidityTime.Add(-defaultTimeDeltaForSlave)
	subscriptionInfo.ValidityTime = validityTime

	subscriptionInfoData, err := json.Marshal(subscriptionInfo)
	if err != nil {
		log.Errorf("Marshal subscriptionInfo fail, err:%s", err.Error())
		rw.Header().Set("Content-Type", httpContentTypeJSON)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", httpContentTypeJSON)
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(subscriptionInfoData)
	if err != nil {
		log.Warnf("%v", err)
	}

	log.Infof("NRF-Discovery-Agent active response subscriptionInfo:%v", subscriptionInfo)
}

func close(req *http.Request) {
	err := req.Body.Close()
	if err != nil {
		log.Error("close http request Body failure")
	}
}

func loadConfigmapStorage() bool {
	subscriptionInfoData, err := k8sapiclient.GetConfigMapData(consts.ConfigMapStorage, consts.ConfigMapKeySubsInfo)
	if err != nil {
		log.Errorf("Load subscriptionInfo from configmap:%s fail, err:%s", consts.ConfigMapStorage, err.Error())
		return false
	}
	if len(subscriptionInfoData) == 0 {
		return true
	}

	rest := cacheManager.SubscriptionInfoProvision(subscriptionInfoData)

	return rest
}

func keepCacheModeSendResponse(rw http.ResponseWriter, req *http.Request) bool {
	if req == nil {
		log.Errorf("Requester is nil")
		return false
	}
	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}

	if worker.IsKeepCacheMode() {
		errorInfo := "Keep cache mode, not found in agent local cache."
		problemDetails := &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("%s", errorInfo),
		}
		logcontent.ResponseDescription = fmt.Sprintf("%s", errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, http.StatusNotFound, problemDetails.ToString())
		return true
	}
	return false
}
