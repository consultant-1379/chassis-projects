package log

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// Fields map information
type Fields map[string]interface{}
type fieldKey string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[fieldKey]string

// Default key names for the default fields
const (
	FieldKeyMsg       = "msg"
	FieldKeyLevel     = "level"
	FieldKeyTime      = "time"
	FieldKeyNF        = "networkFunction"
	FieldKeyServiceID = "serviceId"
)

var (
	serviceID       string
	networkFunction string
)

const defaultTimestampFormat = time.RFC3339

// SetServiceID set service ID
func SetServiceID(id string) {
	serviceID = id
}

// SetNF set NF
func SetNF(nf string) {
	networkFunction = nf
}

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
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
	FieldMap FieldMap
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

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	if !f.DisableTimestamp {
		data[f.FieldMap.resolve(FieldKeyTime)] = entry.Time.Format(timestampFormat)
	}
	data[f.FieldMap.resolve(FieldKeyMsg)] = entry.Message
	//data[f.FieldMap.resolve(FieldKeyLevel)] = entry.Level.String()
	data[f.FieldMap.resolve(FieldKeyLevel)] = levelString(entry.Level)

	data[f.FieldMap.resolve(FieldKeyNF)] = networkFunction
	data[f.FieldMap.resolve(FieldKeyServiceID)] = serviceID

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
