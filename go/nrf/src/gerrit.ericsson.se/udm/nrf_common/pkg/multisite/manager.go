package multisite

import (
	"strconv"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
)

// Manager for Manager struct
type Manager struct {
}

var manager *Manager

//GetManager is to get manager object
func GetManager() *Manager {
	if manager == nil {
		manager = &Manager{}
	}
	return manager
}

//Run is to starting manager
func (m *Manager) Run() {
	log.Debug("starting multisite manager")
	go func() {
		for {
			registerSite()
			time.Sleep(time.Second * time.Duration(configmap.MultisiteHeartbeatTime))
		}
	}()
}

func registerSite() {
	instanceID := cm.NfProfile.InstanceID
	unixtime := time.Now().UnixNano() / 1000000
	timestamp := strconv.FormatInt(unixtime, 10)

	var info StatusInfo
	info.InstanceID = instanceID
	info.LastUpdateTime = timestamp
	info.Fqdn = cm.NfProfile.Fqdn
	info.Weight = GetMonitor().weight
	data, err := encodeMultiSiteInfo(info)
	if err != nil {
		return
	}

	err = dbmgmt.Insert(region, instanceID, data)
        if err != nil {
		log.Error(err)
        }
}
