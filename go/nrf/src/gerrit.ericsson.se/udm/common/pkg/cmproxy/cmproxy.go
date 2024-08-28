package cmproxy

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/buger/jsonparser"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/msgbus"
)

var (
	// PendingInitializeFailed is used by UT/BT
	PendingInitializeFailed = true
)

const (
	// Default eric-cm-mediator URI
	defaultCmMediatorURI = "http://eric-cm-mediator:5003/cm/api/v1.1/"
	// Default retrive timer interval is 5
	defautRetrieveInterval = 5
	// Default subscriptsion leaseSeconds is 3600
	defaultHeartBeatInterval = 2880
)

type cmproxyStatus int

const (
	initializing cmproxyStatus = iota
	runningWithMessageBus
	runningWithoutMessageBus
	shuttingDown
)

func (s cmproxyStatus) String() string {
	status := []string{
		"initializing",
		"running with message bus",
		"running without message bus",
		"shutting down",
	}
	return status[s]
}

var (
	cmMediatorURI string

	status = initializing
	wg     sync.WaitGroup

	cmConfigList     map[string]*cmConfig
	cmConfigListLock sync.Mutex

	cmSubscriptions     map[string]*cmSubscription
	cmSubscriptionsLock sync.Mutex

	retrieveTicker    *time.Ticker
	resubscribeTicker *time.Ticker

	notificationsTopics map[string]bool

	cmHTTPClient *httpclient.HttpClient
	cmMessageBus *msgbus.MessageBus
)

type cmConfig struct {
	subscriptionID string
	jsonPath       string
	topic          string
	handler        callbackHandler
	localETag      string
}

// callback handler should NOT invoke RegisterConf() and DeRegisterConf()
// avoid the dead lock of cmConfigListLock.
type callbackHandler func(event, configName, format string, rawData []byte)

// InitCmproxyLog provide API for CM proxy Logging initialization
func InitCmproxyLog(level log.Level, networkFunc, podIP, serviceID string) {
	// set log level. the value can be ErrorLevel/WarnLevel/InfoLevel/DebugLevel
	log.SetLevel(level)
	// set output
	log.SetOutput(os.Stdout)
	// set network function name, it will be displayed in log output
	log.SetNF(networkFunc)
	// set pod ip, it will be displayed in log output
	log.SetPodIP(podIP)
	// set log format to user-defined json
	// set service ID, it will be displayed in log output
	log.SetServiceID(serviceID)
	// set log format to user-defined json
	log.SetFormatter(&log.JSONFormatter{})
}

// Init provide API for CM Proxy initialization
// parameter cmmSvc: "CMM_SERVICE:PORT"
func Init(cmmURI string) {
	if status != initializing {
		log.Infof("cmproxy has beed initialized")
		return
	}

	if cmmURI == "" {
		log.Infof("cm mediator URI is not provided by application, using default value for cmproxy")
		cmMediatorURI = defaultCmMediatorURI
	} else {
		cmMediatorURI = cmmURI
	}
	log.Infof("cm mediator URI: %s", cmMediatorURI)

	cmConfigList = make(map[string]*cmConfig)
	cmSubscriptions = make(map[string]*cmSubscription)

	notificationsTopics = make(map[string]bool)
}

// Run provide API for CM Proxy execution
// Create tickers to retrieve configurations and subscribe subscriptions
func Run() {
	// Init Http Client for cmproxy
	var (
		defaultConnections = 2
		defaultTimeout     = 10 * time.Second
		defaultKeepAlive   = true
	)
	cmHTTPClient = httpclient.InitHttpClient(
		httpclient.Connections(defaultConnections),
		httpclient.Timeout(defaultTimeout),
		httpclient.KeepAlive(defaultKeepAlive),
	)
	if cmHTTPClient == nil {
		panic("failed to initialize HTTP Client for cmproxy")
	}

	// Initialize message bus for cmproxy
	isMessagebusReady := initMessagebus()
	// Get configurations and post subscriptions
	if isConfigurationsReady := initConfigurations(2 * time.Second); isConfigurationsReady {
		if cmMessageBus != nil && isMessagebusReady {
			status = runningWithMessageBus
		} else {
			status = runningWithoutMessageBus
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		getConfigurationsTimer(defautRetrieveInterval * time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		heartBeatTimer(defaultHeartBeatInterval * time.Second)
	}()

	log.Debugf("initialize cmproxy done")
}

// Stop provide API for CM Proxy shutdown
// Stop tickers and wait for go routines done
func Stop() {
	closeMessageBus()

	status = shuttingDown
	wg.Wait()
}

func initConfigurations(sec time.Duration) bool {
	ready := false
	for !ready && PendingInitializeFailed {
		if ready = getConfigurations(); ready {
			if ready = putSubscriptions(); !ready {
				time.Sleep(sec)
			}
		} else {
			time.Sleep(sec)
		}
	}
	return ready
}

func getConfigurationsTimer(sec time.Duration) {
	retrieveTicker = time.NewTicker(sec)
	if retrieveTicker != nil {
		defer retrieveTicker.Stop()
		for status != shuttingDown {
			select {
			case <-retrieveTicker.C:
				getConfigurations()
			default:
				time.Sleep(time.Second)
			}
		}
	} else {
		log.Errorf("failed to create retrieve timer")
	}
}

func heartBeatTimer(sec time.Duration) {
	resubscribeTicker = time.NewTicker(sec) // Half of leaseSeconds
	if resubscribeTicker != nil {
		defer resubscribeTicker.Stop()
		for status != shuttingDown {
			select {
			case <-resubscribeTicker.C:
				putSubscriptions()
			default:
				time.Sleep(time.Second)
			}
		}
	} else {
		log.Errorf("failed to create heartbeat timer")
	}
}

func cmConfigCreated(configName, jsonPath, topic, id string, callback callbackHandler) bool {
	cmConfigListLock.Lock()
	defer cmConfigListLock.Unlock()

	if cmConfigList == nil {
		return false
	}

	_, existed := cmConfigList[configName]
	if existed {
		log.Infof("%s was registered", configName)
		return false
	}

	cmConfigList[configName] =
		&cmConfig{
			subscriptionID: id,
			jsonPath:       jsonPath,
			topic:          topic,
			handler:        callback,
			localETag:      "",
		}
	notificationsTopics[topic] = true

	return true
}

func cmConfigUpdated(event, configName, format, baseETag, configETag string, rawData []byte) {
	switch event {
	case EventConfigCreated:
		eventConfigCreatedHandler(configName, configETag)
	case EventConfigUpdated:
		eventConfigUpdatedHandler(configName, format, baseETag, configETag, rawData)
	case EventConfigDeleted:
		eventConfigDeletedHandler(configName)
	}
}

func cmConfigDeleted(configName string) {
	cmConfigListLock.Lock()
	defer cmConfigListLock.Unlock()

	delete(cmConfigList, configName)
}

func eventConfigCreatedHandler(configName, configETag string) {
	cmConfigListLock.Lock()
	defer cmConfigListLock.Unlock()

	config, existed := cmConfigList[configName]
	if !existed {
		log.Infof("%s was not registered, ignore this notification", configName)
		delete(cmConfigList, configName)
		return
	}
	if config == nil {
		return
	}

	//Detail content of configuration is not in "configCreated" message.
	//So, clear local etage and trigger HTTP Get behavior, then configurations will
	//be retrieved by retrieve timer, and callback will be invoked.
	log.Infof("%s is created in cm mediator serivce", configName)
	config.localETag = ""
}

func eventConfigUpdatedHandler(configName, format, baseETag, configETag string, rawData []byte) {
	cmConfigListLock.Lock()
	defer cmConfigListLock.Unlock()

	var data []byte
	var err error

	config, existed := cmConfigList[configName]
	if !existed {
		log.Infof("%s was not registered, ignore this notification", configName)
		delete(cmConfigList, configName)
		return
	}
	if config == nil {
		return
	}

	if config.localETag != configETag {
		switch format {
		case NtfFormatPatch:
			if config.localETag != baseETag {
				//trigger timer to getConfiguration()
				log.Infof("failed to update Etag for %s, localETag in notification is not the same with baseETag", configName)
				config.localETag = ""
				return
			}
			data = rawData
		case NtfFormatFull:
			// remove root JSON key
			data, _, _, err = jsonparser.Get(rawData, config.jsonPath)
			if err != nil {
				log.Infof("failed to parse jsonPath %s for %s", config.jsonPath, configName)
				return
			}
		}

		// Store latest etag
		config.localETag = configETag
		// Invoke callback handler to notify App Service
		config.handler(EventConfigUpdated, configName, format, data)
	}
}

func eventConfigDeletedHandler(configName string) {
	cmConfigListLock.Lock()
	defer cmConfigListLock.Unlock()

	config, existed := cmConfigList[configName]
	if !existed {
		log.Infof("%s was not registered, ignore this notification", configName)
		delete(cmConfigList, configName)
		return
	}
	if config == nil {
		return
	}

	log.Infof("%s is removed from cm mediator serivce", configName)
	config.handler(EventConfigDeleted, configName, NtfFormatFull, nil)
}

func getCmmURL(useCases, id string) string {
	url := cmMediatorURI + useCases
	if id != "" {
		url = url + "/" + id
	}
	return url
}

// RegisterConf provide API for APP register Configurations to CM Proxy
func RegisterConf(configName, jsonPath, topic string, callback callbackHandler, format string) {
	subscriptionID := configName + "_" + topic + "_sub"
	if !cmConfigCreated(configName, jsonPath, topic, subscriptionID, callback) {
		return
	}

	s := newSubscription(subscriptionID, configName)
	s.Callback = "kafka:" + topic
	s.UpdateNotificationFormat = format
	if !cmSubscriptionCreated(subscriptionID, s) {
		log.Infof("register configuration %s to cmproxy, but post subscription failed", configName)
		return
	}

	log.Infof("register configuration %s to cmproxy", configName)
}

// DeRegisterConf provide API for APP de-registers Configurations to CM Proxy
func DeRegisterConf(configName, topic string) {
	subscriptionID := configName + "_" + topic + "_sub"
	cmConfigDeleted(configName)
	cmSubscriptionDeleted(subscriptionID)

	deleteSubscription(subscriptionID)
	log.Infof("deregister configuration %s from cmproxy done", configName)
}

func getConfigurations() bool {
	cmConfigListLock.Lock()
	defer cmConfigListLock.Unlock()

	ready := true
	for configName, config := range cmConfigList {
		if config.localETag != "" && status == runningWithMessageBus {
			continue
		}

		resp, err := getConfiguration(configName)
		if err != nil {
			ready = false
		} else {
			var configuration cmConfiguration
			err := json.Unmarshal(resp.Body, &configuration)
			if err != nil {
				log.Errorf("failed to Unmarshal configuration %s, %s, %s.", configName, string(resp.Body), err.Error())
				return false
			}

			cmConfigListLock.Unlock()
			// Notify App Service;
			// event: event, format: full
			cmConfigUpdated(EventConfigUpdated, configName, NtfFormatFull, "", resp.Etag, configuration.Data)
			cmConfigListLock.Lock()
		}
	}
	return ready
}

func putSubscriptions() bool {
	cmSubscriptionsLock.Lock()
	defer cmSubscriptionsLock.Unlock()

	ready := true
	for subscriptionID, s := range cmSubscriptions {
		err := putSubscription(subscriptionID, s)
		if err != nil {
			ready = false
		}
	}
	return ready
}
