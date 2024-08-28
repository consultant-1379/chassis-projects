package fmproxy

import (
	"time"
)

const (
	// AlarmSchemaVersion is the schema version for alarm
	AlarmSchemaVersion = "0.2"

	// AlarmDefaultExpireTime
	AlarmDefaultExpireTime = 300

	// MinExpireTime is minimum valid expiration time
	MinExpireTime = 10

	// RaiseAlarm is to raise alarm
	RaiseAlarm = true

	// ClearAlarm is to clean alarm
	ClearAlarm = false

	// SeverityClear is Clear alarm
	SeverityClear = "Clear"
	// SeverityWarning is warning alarm
	SeverityWarning = "Warning"
	// SeverityMinor is minor alarm
	SeverityMinor = "Minor"
	// SeverityMajor is major alarm
	SeverityMajor = "Major"
	// SeverityCritical is critical alarm
	SeverityCritical = "Critical"
)

// AlarmInfo struct define
type AlarmInfo struct {
	IsAutoResend          bool           // whether FM proxy resend automaticall after expiration
	FaultName             string         // alarm name
	FaultyResource        string         // it may be pod-name
	Expiration            int            // The expiration time of the fault in seconds: -1 mean no expiration, the default is used if app not provide,the minimum valid value is 10s.
	Severity              string         // Severity may be Clear Warning Minor Major Critical, default value is used if it is empty
	AdditionalInformation AdditionalInfo //it is either AddtionMultiKeyValue or JsonAddtion
	Description           string         // Extra information providing further insight about the alarm
	timestamp             time.Time      // the time when the alarm is raise,it is optional
}

// AdditionalInfo is for additionalInformation of alarm
type AdditionalInfo interface {
	getAddtionJsonStr() string
}

// AddtionMultiKeyValue is to request providing the addtionInfo by pattern key:value
type AddtionMultiKeyValue struct {
	Key   string
	Value string
}

// JsonAddtion is to request providing the addtionInfo as json object'
// e.g. JsonAddtion{JsonStr: `{ "parameters": "FM_GATEWAY" }` }
type JsonAddtion struct {
	JsonStr string
}

func (m *AddtionMultiKeyValue) getAddtionJsonStr() string {
	addtionJsonStr := `{ "`
	addtionJsonStr = addtionJsonStr + m.Key + `": "`
	addtionJsonStr = addtionJsonStr + m.Value + `" }`
	return addtionJsonStr
}

func (m *JsonAddtion) getAddtionJsonStr() string {
	return m.JsonStr
}
