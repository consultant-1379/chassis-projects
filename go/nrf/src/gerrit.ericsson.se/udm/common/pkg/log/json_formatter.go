package log

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type logBody struct {
	Timestamp       string      `json:"timestamp"`
	Level           string      `json:"level"`
	NetworkFunction string      `json:"networkFunction"`
	ServiceId       string      `json:"serviceId"`
	PodIP           string      `json:"podIp"`
	Message         interface{} `json:"message"`
	StackInfo       string      `json:"stackinfo"`
}

var (
	serviceID       string
	networkFunction string
	podIP           string
)

const defaultTimestampFormat = time.RFC3339

// SetServiceID set service name
func SetServiceID(id string) {
	serviceID = id
}

// SetNF set network function name
func SetNF(nf string) {
	networkFunction = nf
}

// SetPodIP set pod ip
func SetPodIP(ip string) {
	podIP = ip
}

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime: "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg: "@message",
	//    },
	// }
	// FieldMap FieldMap
}

func levelString(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "DEBUG"
	case logrus.InfoLevel:
		return "INFO"
	case logrus.WarnLevel:
		return "WARN"
	case logrus.ErrorLevel:
		return "ERROR"
	case logrus.FatalLevel:
		return "FATAL"
	case logrus.PanicLevel:
		return "PANIC"
	}

	return "UNKNOWN"
}

// LevelToString return string format of log level
func LevelToString(level Level) string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	}

	return "UNKNOWN"
}

// LevelUint return int log level from string log level
func LevelUint(level string) Level {
	switch level {
	case "DEBUG":
		return DebugLevel
	case "INFO", "NOTICE":
		return InfoLevel
	case "WARNING":
		return WarnLevel
	case "ERROR", "CRITICAL", "ALERT", "EMERGENCY":
		return ErrorLevel
	}

	return WarnLevel
}

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logBodyIns := &logBody{
		Level:           levelString(entry.Level),
		NetworkFunction: networkFunction,
		ServiceId:       serviceID,
		PodIP:           podIP,
		//Message:         entry.Message,
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	if !f.DisableTimestamp {
		logBodyIns.Timestamp = entry.Time.Format(timestampFormat)
	}

	if _, ok := entry.Data["isjson"]; ok {
		logBodyIns.Message = json.RawMessage(entry.Message)
	} else {
		logBodyIns.Message = entry.Message
	}

	if si, ok := entry.Data["stackinfo"]; ok {
		logBodyIns.StackInfo = si.(string)
	}

	serialized, err := json.Marshal(logBodyIns)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
