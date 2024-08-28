package fm

import (
	"os"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/fmproxy"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

//AlarmMapTable static alarm info
type AlarmMapTable map[string]fmproxy.AlarmInfo

var alarmInfoMap = AlarmMapTable{
	"nrf-mgmt": fmproxy.AlarmInfo{
		IsAutoResend: true,
		FaultName:    "nrfagentNrfmngtConnectionFailure",
		Expiration:   600,
	},
	"nrf-disc": fmproxy.AlarmInfo{
		IsAutoResend: true,
		FaultName:    "nrfagentNrfdiscoverConnectionFailure",
		Expiration:   600,
	},
	"noAvailableDestination": fmproxy.AlarmInfo{
		IsAutoResend: true,
		FaultName:    "nrfagentNoAvailableDestination",
		Expiration:   600,
	},
}

var (
	alarmRaised           = make(map[string]bool)
	unavailableNfTypeList = make(map[string]string)
)

var sendAlarm = func(alarmPara *fmproxy.AlarmInfo, isRaise bool) error {
	return fmproxy.SendAlarm(alarmPara, isRaise)
}

// Init Init fmproxy module
func Init(serviceName string) {
	err := fmproxy.Init(os.Getenv("MESSAGE_BUS_KAFKA"), serviceName)
	if err != nil {
		log.Errorf("fmproxy init fail, error message is: %s", err.Error())
	}
	log.Infof("fmproxy init success")
}

// ConnectionStatus : Report the connection status to fmproxy to trigger the ConnectionFailure error/clear alarm
// parameter target : The connection target, e.g nrf-mgm, the func will use it as the key get the alarm info from table "alarmInfoMap"
// parameter available : The connection status, the func will generate  error or clear alarm according to the status
func ConnectionStatus(target string, available bool) {
	raised := alarmRaised[target]
	if (!available && raised) ||
		(available && !raised) {
		return
	}

	alarm, ok := alarmInfoMap[target]
	if !ok {
		log.Errorf("No alarm definition for target connection: %s", target)
		return
	}
	alarm.FaultyResource = os.Getenv("POD_NAME")

	err := sendAlarm(&alarm, !available)
	if err != nil {
		log.Errorf("send alarm to fmproxy fail, %s", err.Error())
		return
	}
	if !available {
		log.Debugf("send alarm %s to fmproxy done", alarm.FaultName)
	} else {
		log.Debugf("clear alarm %s to fmproxy done", alarm.FaultName)
	}
	alarmRaised[target] = !raised
}

// DestinationStatus : Report the destination NF status to fmproxy to trigger the
// parameter target : The destination NF Type (with filter?)
// parameter available : The connection status, the func will generate  error or clear alarm according to the status
func DestinationStatus(target string, available bool, requesterNf, additionalInfo string) {
	alarm, ok := alarmInfoMap[target]
	if !ok {
		log.Errorf("No alarm definition for target connection: %s", target)
		return
	}
	alarm.FaultyResource = os.Getenv("POD_NAME")

	jsonAddtionStrVal := unavailableNfTypeList[requesterNf]
	if additionalInfo != "" {
		if !available {
			if !strings.Contains(jsonAddtionStrVal, additionalInfo) {
				jsonAddtionStrVal = jsonAddtionStrVal + additionalInfo + ","
			} else {
				return
			}
		} else {
			if strings.Contains(jsonAddtionStrVal, additionalInfo) {
				jsonAddtionStrVal = strings.Replace(jsonAddtionStrVal, additionalInfo+",", "", 1)
			} else {
				return
			}
		}
		unavailableNfTypeList[requesterNf] = jsonAddtionStrVal
	} else {
		unavailableNfTypeList[requesterNf] = ""
	}

	var jsonAddition string
	for _, addition := range unavailableNfTypeList {
		if addition != "" {
			jsonAddition = jsonAddition + addition
		}
	}

	// remove last ',' character
	jsonAddition = strings.TrimSuffix(jsonAddition, ",")

	var err error
	if len(jsonAddition) != 0 {
		alarm.AdditionalInformation = &fmproxy.AddtionMultiKeyValue{Key: "targetNfType", Value: jsonAddition}
		err = sendAlarm(&alarm, true)
	} else {
		alarm.AdditionalInformation = nil
		err = sendAlarm(&alarm, false)
	}
	if err != nil {
		log.Errorf("send alarm to fmproxy fail, %s", err.Error())
		return
	}
	log.Infof("send alarm %s (%s) to fmproxy done", alarm.FaultName, jsonAddition)
}
