package fmproxy

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

var (
	serviceName string
	//faultMap is to store alarm information
	faultMap      = make(map[string]*AlarmInfo)
	faultMapMutex = new(sync.Mutex)
)

// Init provide API for FM Proxy initialization
// parameter kafkaConn: "FM_SERVICE:PORT"
// parameter alarmServiceName: alarm modeling file name  e.g."nrf-alarm.json", nrf-alarm is serviecname
func Init(kafkaConn string, alarmServiceName string) error {
	if alarmServiceName == "" {
		return errors.New("serviceName of alarm is empty")
	} else {
		serviceName = alarmServiceName
	}

	if kafkaConn == "" {
		log.Infof("kafka connection is not provided by application, using default value for fmproxy")
		kafkaConnection = DefaultKafka
	} else {
		kafkaConnection = kafkaConn
	}
	log.Infof("kafka connection: %s", kafkaConnection)

	return initMsgbus()
}

// SendAlarm function send alarm message to kafka via sync producer
// parameter isRaise: RaiseAlarm(true) is to raise alarm, ClearAlarm(false) is to clean alarm
func SendAlarm(alarmPara *AlarmInfo, isRaise bool) error {
	if alarmPara.FaultName == "" {
		log.Error("faultName of alarm is empty")
		return errors.New("faultName of alarm is empty")
	}

	if alarmPara.FaultyResource == "" {
		log.Error("faultyResource of alarm is empty")
		return errors.New("faultyResource of alarm is empty")
	}

	alarmPara.timestamp = time.Now()
	if RaiseAlarm != isRaise {
		alarmPara.Severity = "Clear"
	}

	faultMapMutex.Lock()
	defer faultMapMutex.Unlock()

	alarmKey := alarmPara.FaultName + alarmPara.FaultyResource
	_, found := faultMap[alarmKey]

	log.Debugf("key: %v, found: %v", alarmKey, found)

	//clear alarm
	if found && !isRaise {
		err := sendMsg(alarmPara)
		delete(faultMap, alarmKey)
		return err
	}

	if isRaise {
		// update alarm
		if !found {
			err := sendMsg(alarmPara)
			if err != nil {
				return err
			}
			go autoResendAlarm(alarmKey, alarmPara.Expiration)

		} else {

			if faultMap[alarmKey].Severity != alarmPara.Severity {
				err := sendMsg(alarmPara)
				if err != nil {
					return err
				}
			}
		}
		faultMap[alarmKey] = alarmPara
	}

	return nil
}

func structureAlarmMsg(alarmPara *AlarmInfo) string {
	alarmMsg := "{ "
	alarmMsg = alarmMsg + `"version": "` + AlarmSchemaVersion + `", `
	alarmMsg = alarmMsg + `"faultName": "` + alarmPara.FaultName + `", `
	alarmMsg = alarmMsg + `"serviceName": "` + serviceName + `", `
	alarmMsg = alarmMsg + `"faultyResource": "` + alarmPara.FaultyResource + `", `
	t := time.Now()
	n := t.UnixNano() / 1000
	alarmMsg = alarmMsg + `"eventTime": ` + strconv.FormatInt(n, 10)

	if alarmPara.Severity != "" {
		alarmMsg = alarmMsg + `, "severity": "` + alarmPara.Severity + `"`
	}

	if alarmPara.Description != "" {
		alarmMsg = alarmMsg + `, "description": "` + alarmPara.Description + `"`
	}

	switch {
	case alarmPara.Expiration >= MinExpireTime:
		alarmMsg = alarmMsg + `, "expiration": ` + strconv.Itoa(alarmPara.Expiration)
	case alarmPara.Expiration < MinExpireTime && alarmPara.Expiration > 0:
		alarmMsg = alarmMsg + `, "expiration": ` + strconv.Itoa(MinExpireTime)
	case alarmPara.Expiration < 0:
		alarmMsg = alarmMsg + `, "expiration": 0`
	default:
	}

	if nil != alarmPara.AdditionalInformation {
		additionalInformation := alarmPara.AdditionalInformation.getAddtionJsonStr()
		if additionalInformation != "" {
			alarmMsg = alarmMsg + `, "additionalInformation": ` + additionalInformation
		}
	}

	alarmMsg = alarmMsg + ` }`
	return alarmMsg
}

func autoResendAlarm(alarmKey string, expireTime int) {
	var expiration int
	if expireTime >= MinExpireTime {
		expiration = expireTime
	} else {
		expiration = MinExpireTime
	}

	interval := expiration/2 + 1
	autoResendTicker := time.NewTicker(time.Second * (time.Duration)(interval))
	if autoResendTicker != nil {
		alarmPara, found := faultMap[alarmKey]
		defer autoResendTicker.Stop()
		for found {
			select {
			case <-autoResendTicker.C:
				faultMapMutex.Lock()
				alarmPara, found = faultMap[alarmKey]
				if found {
					if !alarmPara.IsAutoResend {
						if isMsgExpired(alarmPara, expiration) {
							delete(faultMap, alarmKey)
							faultMapMutex.Unlock()
							break
						}
					}
					err := sendMsg(alarmPara)
					if err != nil {
						log.Errorf("failed to resend alarm: %s", err.Error())
					}
				}
				faultMapMutex.Unlock()
			default:
				time.Sleep(time.Second)
			}
		}
	} else {
		log.Errorf("failed to create ticker for resending alarm")
	}
}

func isMsgExpired(alarmPara *AlarmInfo, expiration int) bool {
	currentTime := time.Now().UnixNano() / 1000000
	timestamp := alarmPara.timestamp.UnixNano() / 1000000
	expireTime := timestamp + (int64(expiration) * 1000)
	return currentTime > expireTime
}
