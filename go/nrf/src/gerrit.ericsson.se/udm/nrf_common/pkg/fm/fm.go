package fm

import (
	"fmt"
	"os"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/fmproxy"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

const (
	managementService = "eric-nrf-nnrf-nfm"
	discoveryService  = "eric-nrf-nnrf-disc"

	faultyResourceNotification            = "/nrfe:nrf/non_modeled/nfmanagement/notification"
	faultyResourceManagement              = "/nrfe:nrf/non_modeled/nfmanagement"
	faultyResourceManagementLocal         = "/nrfe:nrf/non_modeled/local[service-name=nnrf-nfm]"
	faultyResourceDiscoveryLocal          = "/nrfe:nrf/non_modeled/local[service-name=nnrf-disc]"
	faultyResourceManagementWithRemote    = "/nrfe:nrf/non_modeled/local[service-name=nnrf-nfm]/remote[%s]"
	faultyResourceDiscoveryWithRemote     = "/nrfe:nrf/non_modeled/local[service-name=nnrf-disc]/remote[%s]"
	alarmNameNrfManagementServiceOverload = "nrfMngtServiceOverloaded"
)

var serviceName string

func Init() {
	kafkaConnectInfo := os.Getenv("MESSAGE_BUS_KAFKA")
	fmFaultMappingFile := os.Getenv("FM_FAULT_MAPPING_FILE")
	if fmFaultMappingFile == "" {
		fmFaultMappingFile = managementService + ".json"
	}
	serviceName = strings.Split(fmFaultMappingFile, ".")[0]

	err := fmproxy.Init(kafkaConnectInfo, serviceName)
	if err != nil {
		log.Errorf("init FM Proxy fail")
	}
}

// RaiseNRFConnectionFailureAlarm send a alarm
// additionalKey: key of additional information
// additionalInfo: value of additional information
func RaiseNRFConnectionFailureAlarm(additionalKey, additionalInfo, remoteInfo string) {
	log.Debugf("sending NRF connnection failure alarm...")

	alarm := buildNRFConnectionFailureAlarm(remoteInfo)
	if alarm == nil {
		return
	}
	jsonAddtion := &fmproxy.AddtionMultiKeyValue{
		Key:   additionalKey,
		Value: additionalInfo,
	}
	alarm.Expiration = 600
	alarm.AdditionalInformation = jsonAddtion
	err := fmproxy.SendAlarm(alarm, fmproxy.RaiseAlarm)
	if err != nil {
		log.Errorf("sending alarm fail")
	}
}

// ClearNRFConnectionFailureAlarm clear NRFConnectionFailureAlarm
func ClearNRFConnectionFailureAlarm(remoteInfo string) {
	log.Debugf("clearing NRF connnection failure alarm...")

	alarm := buildNRFConnectionFailureAlarm(remoteInfo)
	if alarm == nil {
		return
	}
	err := fmproxy.SendAlarm(alarm, fmproxy.ClearAlarm)
	if err != nil {
		log.Errorf("clearing alarm fail")
	}
}

func buildNRFConnectionFailureAlarm(remoteInfo string) *fmproxy.AlarmInfo {
	faulName := ""
	faultResource := ""
	switch serviceName {
	case managementService:
		faulName = "nrfMngtNrfConnectionFailure"
		faultResource = fmt.Sprintf(faultyResourceManagementWithRemote, remoteInfo)
	case discoveryService:
		faulName = "nrfDiscNrfConnectionFailure"
		faultResource = fmt.Sprintf(faultyResourceDiscoveryWithRemote, remoteInfo)
	default:
		log.Errorf("service not supported")
		return nil
	}

	return &fmproxy.AlarmInfo{
		IsAutoResend:   false,
		FaultName:      faulName,
		FaultyResource: faultResource,
	}
}

// RaiseNRFReplicationFailureAlarm send a alarm
func RaiseNRFReplicationFailureAlarm(instanceID, fqdn string) {
	log.Debugf("sending NRF Replication failure alarm...")

	info := fmt.Sprintf("{InstanceID:%s, FQDN:%s}", instanceID, fqdn)
	jsonAddtion := &fmproxy.AddtionMultiKeyValue{
		Key:   "RemoteInfo",
		Value: info,
	}

	remote := "FQDN:" + fqdn
	fr := fmt.Sprintf(faultyResourceManagementWithRemote, remote)
	alarm := &fmproxy.AlarmInfo{
		IsAutoResend:          false,
		FaultName:             "nrfDataReplicationConnectionFailure",
		FaultyResource:        fr,
		Expiration:            600,
		AdditionalInformation: jsonAddtion,
	}

	err := fmproxy.SendAlarm(alarm, fmproxy.RaiseAlarm)
	if err != nil {
		log.Errorf("sending alarm fail")
	}
}

// ClearNRFReplicationFailureAlarm clear NRFReplicationFailureAlarm
func ClearNRFReplicationFailureAlarm(fqdn string) {
	log.Debugf("clearing NRF Replication failure alarm...")

	remote := "FQDN:" + fqdn
	fr := fmt.Sprintf(faultyResourceManagementWithRemote, remote)
	alarm := &fmproxy.AlarmInfo{
		IsAutoResend:   false,
		FaultName:      "nrfDataReplicationConnectionFailure",
		FaultyResource: fr,
	}

	err := fmproxy.SendAlarm(alarm, fmproxy.ClearAlarm)
	if err != nil {
		log.Errorf("sending alarm fail")
	}
}

// SendNRFNotificationOverloadAlarm send/clear a alarm
func SendNRFNotificationOverloadAlarm(isRaise bool, severity string) {
	if isRaise {
		log.Debugf("sending NRF notification overload alarm...")
	} else {
		log.Debugf("clear NRF notification overload alarm...")
	}

	alarm := &fmproxy.AlarmInfo{
		IsAutoResend:   true,
		FaultName:      alarmNameNrfManagementServiceOverload,
		FaultyResource: faultyResourceNotification,
		Expiration:     60,
		Severity:       severity,
		Description:    "NRF Management service overload occurs during sending notifications to NF Service Consumers.",
	}

	err := fmproxy.SendAlarm(alarm, isRaise)
	if err != nil {
		log.Errorf("sending alarm fail")
	}
}

// SendNRFManagementOverloadAlarm send/clear a alarm
func SendNRFManagementOverloadAlarm(isRaise bool) {
	if isRaise {
		log.Debugf("System enters overload status, send NRF management service overload alarm...")
	} else {
		log.Debugf("System leaves overload status, clear NRF management service overload alarm...")
	}

	alarm := &fmproxy.AlarmInfo{
		IsAutoResend:   true,
		FaultName:      alarmNameNrfManagementServiceOverload,
		FaultyResource: faultyResourceManagement,
		Expiration:     60,
		Description:    "NRF Management service overload occurs during handling requests from NF Service Consumers.",
	}

	err := fmproxy.SendAlarm(alarm, isRaise)
	if err != nil {
		log.Errorf("sending alarm fail")
	}
}

// RaiseNRFDatabaseConnectionFailureAlarm raise NRFDatabaseConnectionFailureAlarm
func RaiseNRFDatabaseConnectionFailureAlarm() {
	log.Debugf("Sending NRF Database Connnection Failure Alarm...")
	sendNRFDatabaseConnectionFailureAlarm(fmproxy.RaiseAlarm)
}

// ClearNRFDatabaseConnectionFailureAlarm clear NRFDatabaseConnectionFailureAlarm
func ClearNRFDatabaseConnectionFailureAlarm() {
	log.Debugf("Clearing NRF Database Connnection Failure Alarm...")
	sendNRFDatabaseConnectionFailureAlarm(fmproxy.ClearAlarm)
}

func sendNRFDatabaseConnectionFailureAlarm(isRaise bool) {
	faulName := ""
	faultResource := ""
	switch serviceName {
	case managementService:
		faulName = "nrfMngtDatabaseConnectionFailure"
		faultResource = faultyResourceManagementLocal
	case discoveryService:
		faulName = "nrfDiscDatabaseConnectionFailure"
		faultResource = faultyResourceDiscoveryLocal
	default:
		log.Errorf("service not supported")
		return
	}

	alarm := &fmproxy.AlarmInfo{
		IsAutoResend:   false,
		FaultName:      faulName,
		FaultyResource: faultResource,
	}
	if isRaise {
		alarm.Expiration = 600
	}
	err := fmproxy.SendAlarm(alarm, isRaise)
	if err != nil {
		log.Errorf("sending alarm fail")
	}
}
