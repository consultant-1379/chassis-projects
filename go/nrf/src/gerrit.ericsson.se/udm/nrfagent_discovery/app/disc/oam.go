package disc

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/cmproxy"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/pm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/fm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"

	"github.com/buger/jsonparser"
	"github.com/fsnotify/fsnotify"
)

func cmNrfAgentConfHandler(Event, ConfigurationName, format string, RawData []byte) {
	log.Infof("cmNrfAgentConfHandler: %s, %s, %s", Event, format, string(RawData))
	if format != cmproxy.NtfFormatFull {
		log.Warnf("notification format:%s is not recommended", format)
		return
	}

	nrfProfilesData, _, _, err := jsonparser.Get(RawData, "nrf")
	if err != nil {
		log.Errorf("failed to get nrf in %s, %s", ConfigurationName, err.Error())
	} else {
		log.Infof("NRF Profile: %s", string(nrfProfilesData))
		structs.UpdateNrfServerList(nrfProfilesData)
		client.ResetNrfServerPrefix()
	}
	//	nrfAgentCommonFuncData, _, _, err := jsonparser.Get(RawData, "nrf-agent-common")
	//	if err != nil {
	//		log.Errorf("failed to get nrf-agent-common in %s, %s", ConfigurationName, err.Error())
	//	} else {
	//		log.Infof("common function: %s", string(nrfAgentCommonFuncData))
	//		//		UpdateCommonFunction(nrfAgentCommonFuncData)
	//	}
	statusNotifIPEndpoint, _, _, err := jsonparser.Get(RawData, "notification-address")
	if err != nil {
		log.Errorf("failed to get notification-address in %s, %s", ConfigurationName, err.Error())
	} else {
		log.Infof("notification-address: %s", string(statusNotifIPEndpoint))
		structs.UpdateStatusNotifIPEndPoint(statusNotifIPEndpoint)
	}
}

type callbackHandler func(event, fileName string, rawData []byte)

func configmapMonitor(notifyDir string, sec time.Duration, handler callbackHandler) {
	absNotifyDir, _ := filepath.Abs(notifyDir)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("fsnotify failed to create fsWatcher %s, %s", absNotifyDir, err.Error())
		return
	}
	defer func() {
		if watcher != nil {
			_ = watcher.Close()
		}
	}()

	fsTicker := time.NewTicker(sec)
	if fsTicker == nil {
		log.Errorf("fsnotify failed to create timer")
	}
	defer func() {
		if fsTicker != nil {
			fsTicker.Stop()
		}
	}()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Debugf("fsnotify get event: %+v", event)
				go configmapMonitorEventHandler(event, handler)

			case e := <-watcher.Errors:
				log.Errorf("fsnotify error: %s", e.Error())

			case <-fsTicker.C:
				go configmapMonitorTimerHandler(absNotifyDir, handler)
			}
		}
	}()

	err = watcher.Add(absNotifyDir)
	if err != nil {
		log.Errorf("fsnotify add %s to fsWatcher error: %s", absNotifyDir, err.Error())
		return
	}
	log.Debugf("fsnotify add %s to fsWatcher done", absNotifyDir)

	<-done
}

func configmapMonitorEventHandler(event fsnotify.Event, handler callbackHandler) {
	if event.Op&fsnotify.Write == fsnotify.Write ||
		event.Op&fsnotify.Create == fsnotify.Create {
		absConfName, _ := filepath.Abs(event.Name)
		rawData, err := ioutil.ReadFile(event.Name)
		if err != nil {
			log.Errorf("fsnotify failed to load file %s, %s", absConfName, err.Error())
			return
		}
		if len(rawData) == 0 {
			log.Infof("fsnotify %s is empty", absConfName)
			return
		}
		if handler != nil {
			handler(event.Op.String(), absConfName, rawData)
		}
	}
}

func configmapMonitorTimerHandler(notifyDir string, handler callbackHandler) {
	files, _ := ioutil.ReadDir(notifyDir)
	for _, f := range files {
		absConfName := notifyDir + "/" + f.Name()
		c := structs.CheckTargetNfProfilesByName(absConfName)
		if !c {
			rawData, err := ioutil.ReadFile(absConfName)
			if err != nil {
				log.Errorf("fsnotify failed to load file %s", err.Error())
				continue
			}
			if len(rawData) == 0 {
				log.Infof("fsnotify %s is empty", absConfName)
				continue
			}
			if handler != nil {
				handler(fsnotify.Create.String(), absConfName, rawData)
			}
		}
	}
}

func cmTargetNfProfilesHandler(Event, ConfigurationName, format string, RawData []byte) {
	log.Infof("cmTargetNfProfilesHandler: %s, %s, %s", Event, format, string(RawData))
	if format != cmproxy.NtfFormatFull {
		log.Warnf("notification format:%s is not recommended", format)
		return
	}
	if len(RawData) == 0 {
		return
	}

	var rawData []byte
	var err error
	rawData, _, _, err = jsonparser.Get(RawData, "targetNfProfiles")
	if err != nil {
		log.Errorf("failed to run jsonparser.Get() targetNfProfiles, %s", err.Error())
		return
	}
	structs.UpdateTargetNfProfilesByName(ConfigurationName, rawData)
	CreateNfTypeTopic()
}

func pmNrfDiscoveryResponses(code int) {
	if code >= 200 && code <= 299 {
		pm.Inc(consts.NrfDiscoveryResponses2xx)
	} else if code >= 300 && code <= 399 {
		pm.Inc(consts.NrfDiscoveryResponses3xx)
	} else if code >= 400 && code <= 499 {
		pm.Inc(consts.NrfDiscoveryResponses4xx)
	} else if code >= 500 && code <= 599 {
		pm.Inc(consts.NrfDiscoveryResponses5xx)
	}
	pm.Inc(consts.NrfDiscoveryResponsesTotal)
}

func fmRaiseNoAvailableDestination(requesterNfType, targetNfType string) {
	var searchParameter cache.SearchParameter
	searchParameter.SetRequesterNfType(requesterNfType)
	searchParameter.SetTargetNfType(targetNfType)
	if _, hit := cache.Instance().Search(requesterNfType, targetNfType, &searchParameter, false); !hit {
		fm.DestinationStatus("noAvailableDestination", false, requesterNfType, targetNfType)
	}
}

func fmClearNoAvailableDestination(requesterNfType, targetNfType string) {
	fm.DestinationStatus("noAvailableDestination", true, requesterNfType, targetNfType)
}
