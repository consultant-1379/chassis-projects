package witness

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
	witnessID   string
	witnessFQDN string
}

var manager *Manager

//GetManager is to get manager object
func GetManager() *Manager {
	if manager == nil {
		manager = &Manager{
			witnessID:   "",
			witnessFQDN: "",
		}
	}
	return manager
}

func (m *Manager) IsWitnessActive() bool {
	witness := getWitnessFromDB()
	if witness != nil && !isWitnessExpired(witness) {
		return true
	}
	return false
}

func getWitnessFromDB() *WitnessInfo {
	result, err := dbmgmt.GetByKey(region, "witness")
        if err != nil || result == "" {
		log.Warning("Cannot get Witness")
		return nil
	}

	witness, err := decodeWitnessInfo(result)
	if err != nil {
		return nil
	}

	return &witness
}

func isWitnessExpired(witness *WitnessInfo) bool {
	unixtime := time.Now().UnixNano() / 1000000
	expireTime := unixtime - (int64(configmap.MultisiteExpireTime) * 1000)

	timestamp, err := strconv.ParseInt(witness.LastUpdateTime, 10, 64)
	if err != nil {
		return false
	}

	return timestamp < expireTime
}

func updateWitnessInDB(instanceId, fqdn string) {
	unixtime := time.Now().UnixNano() / 1000000
	timestamp := strconv.FormatInt(unixtime, 10)

	var info WitnessInfo
	info.InstanceID = instanceId
	info.Fqdn = fqdn
	info.LastUpdateTime = timestamp
	data, err := encodeWitnessInfo(info)
	if err != nil {
		return
	}

	err = dbmgmt.Insert(region, "witness", data)
        if err != nil {
                log.Error(err)
        }
}

/*
func clearWitnessInDB() {
	cmd := "remove --region=/" + region + " --all"
        _, err := kvdbclient.GetInstance().SendGFSHCommand(cmd)
        if err != nil {
                log.Warning("send GFSH command fail by kvdbclient")
        }
}
*/

func removeWitnessInDB(instanceId string) {
	var keys []string
	keys = append(keys, "withness")
	err := dbmgmt.Remove(region, keys)
        if err != nil {
                log.Error(err)
        }
}

func (m *Manager) IsWitness(instanceId, fqdn string) bool {
	switch cm.NrfCommon.GeoRed.WitnessNF.IdentityType {
	case "nf-instance-id":
		if cm.NrfCommon.GeoRed.WitnessNF.IdentityValue == instanceId {
			m.witnessID = instanceId
			m.witnessFQDN = fqdn
			return true
		}
	case "nf-fqdn":
		if cm.NrfCommon.GeoRed.WitnessNF.IdentityValue == fqdn {
			m.witnessID = instanceId
			m.witnessFQDN = fqdn
			return true
		}
	default:
		log.Warning("Unsupported WitnessKeyName")
	}
	return false
}

func (m *Manager) isWitness(instanceId string) bool {
	return m.witnessID == instanceId
}

// DeregisterWitness for witness deregister & heartbeat failure
func (m *Manager) DeregisterWitness(nfInstanceId string) {
	if !m.isWitness(nfInstanceId) {
		return
	}

	removeWitnessInDB(nfInstanceId)
	m.witnessID = ""
	m.witnessFQDN = ""
}

// UpdateWitness for witness heartbeat update
func (m *Manager) UpdateWitness(nfInstanceId string) {
	if !m.isWitness(nfInstanceId) {
		return
	}

	updateWitnessInDB(m.witnessID, m.witnessFQDN)
}
