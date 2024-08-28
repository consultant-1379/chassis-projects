package multisite

import (
	//"fmt"
	//"sort"
	"time"
	//"strings"
	"strconv"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/fm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/multisite/witness"
)

// Monitor struct for multi-site
type Monitor struct {
	isLeaderSite bool
	isActiveSite bool
	weight       float32
}

var monitor *Monitor

// GetMonitor for multi-site
func GetMonitor() *Monitor {
	if monitor == nil {
		monitor = &Monitor{
			isLeaderSite: false,
			isActiveSite: false,
			weight:       1.0,
		}
	}
	return monitor
}

// Run for Monitor
func (m *Monitor) Run() {
	go func() {
		for {
			m.updateSiteStatus()
			time.Sleep(time.Second * time.Duration(configmap.MultisiteMonitorTime))
		}
	}()
}

func (m *Monitor) updateSiteStatus() {
	allSites := getAllSites()
	activeSites, downSites := splitSites(allSites)

	handleDownSites(downSites)

	allNum := len(allSites)
	activeNum := len(activeSites)

	if allNum == 0 || activeNum == 0 {
		m.isActiveSite = false
		m.isLeaderSite = false
		m.weight = 0
		return
	}

	activeLeaderID, activeLeaderWeight := getLeader(activeSites)

	m.weight = float32(activeNum) / float32(allNum)

	if m.weight == 0.5 && activeLeaderWeight == 0.5 {
		// only half sites active, check witness
		m.isActiveSite = witness.GetManager().IsWitnessActive()
	} else {
		m.isActiveSite = (m.weight > 0.5) && (m.weight == activeLeaderWeight)
	}

	if m.isActiveSite {
		instanceID := cm.NfProfile.InstanceID
		m.isLeaderSite = (instanceID == activeLeaderID)
	} else {
		m.isLeaderSite = false
	}

	log.Debug(*m)
}

// IsLeaderSite to show if site is leader
func (m *Monitor) IsLeaderSite() bool {
	return m.isActiveSite && m.isLeaderSite
}

// IsActiveSite to show if site is active
func (m *Monitor) IsActiveSite() bool {
	log.Debugf("ServiceName: %v", cm.ServiceName)
	log.Debug(cm.NrfCommon.GeoRed)
	if cm.ServiceName == cm.ManagementWorkMode && cm.NrfCommon.GeoRed.KeepManagementService {
		return true
	}

	if cm.ServiceName == cm.DiscoveryWorkMode && cm.NrfCommon.GeoRed.KeepDiscoveryService {
		return true
	}

	if cm.ServiceName == cm.NotificationWorkMode && cm.NrfCommon.GeoRed.KeepManagementService {
		return true
	}

	if cm.ServiceName == cm.ProvsionWorkMode && cm.NrfCommon.GeoRed.KeepManagementService {
		return true
	}

	return m.isActiveSite
}

func getAllSites() []StatusInfo {
	var sites []StatusInfo
	oql := "SELECT DISTINCT * FROM /" + region

	result, err := dbmgmt.GetByOQL(region, oql)
	if err != nil || len(result) == 0 {
		log.Warning("Cannot get site info")
		return sites
	}

	for i := range result {
		site, err := decodeMultiSiteInfo(result[i])
		if err == nil {
			sites = append(sites, site)
		}
	}
	log.Debug("All sites from DB:")
	log.Debug(sites)

	return sites
}

func splitSites(sites []StatusInfo) (activeSites, downSites []StatusInfo) {
	unixtime := time.Now().UnixNano() / 1000000
	expireTime := unixtime - (int64(configmap.MultisiteExpireTime) * 1000)
	for i := range sites {
		timestamp, err := strconv.ParseInt(sites[i].LastUpdateTime, 10, 64)
		if err != nil {
			continue
		}

		if timestamp > expireTime {
			activeSites = append(activeSites, sites[i])
		} else {
			downSites = append(downSites, sites[i])
		}
	}
	log.Debug("Active sites from DB:")
	log.Debug(activeSites)
	log.Debug("Down sites from DB:")
	log.Debug(downSites)
	return activeSites, downSites
}

// leader with largest weight and minimal InstanceID
func getLeader(activeSites []StatusInfo) (string, float32) {
	leaderID := activeSites[0].InstanceID
	leaderWeight := activeSites[0].Weight
	for i := range activeSites {
		if activeSites[i].Weight > leaderWeight {
			leaderWeight = activeSites[i].Weight
			leaderID = activeSites[i].InstanceID
		} else if activeSites[i].Weight == leaderWeight && activeSites[i].InstanceID < leaderID {
			leaderID = activeSites[i].InstanceID
		}

		if cm.ServiceName == cm.ManagementWorkMode {
			fm.ClearNRFReplicationFailureAlarm(activeSites[i].Fqdn)
		}
	}
	return leaderID, leaderWeight
}

func handleDownSites(downSites []StatusInfo) {
	for i := range downSites {
		// skip self
		if downSites[i].InstanceID == cm.NfProfile.InstanceID {
			continue
		}

		if cm.ServiceName == cm.ManagementWorkMode {
			fm.RaiseNRFReplicationFailureAlarm(downSites[i].InstanceID, downSites[i].Fqdn)
		}
	}
}
