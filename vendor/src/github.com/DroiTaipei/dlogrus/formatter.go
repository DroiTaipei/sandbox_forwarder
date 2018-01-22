package dlogrus

import (
	"encoding/json"
	"fmt"

	"github.com/DroiTaipei/logrus"
)

const (
	DEFAULT_TIMESTAMP_FORMAT = "2006-01-02 15:04:05.000000"
	BUILT_IN_FIELD_NUM       = 6
)

type (
	Formatter struct{}
	byKey     [][2]string
)

type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

var restrictedTags []string

func init() {
	restrictedTags = []string{
		"V",   // Version
		"M",   // Message
		"T",   // Time
		"L",   // Label
		"R",   // Role
		"Aid", // AppId
		"Rid", // Request Id
		"Pd",  // Pod Name
		"Nd",  // Node Idetifier
	}
}

func newJSONFormatter(timestampFormat string) *JSONFormatter {
	if timestampFormat == "" {
		timestampFormat = DEFAULT_TIMESTAMP_FORMAT
	}
	return &JSONFormatter{TimestampFormat: timestampFormat}
}

func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+BUILT_IN_FIELD_NUM)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/devopstaku/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	data[POD_NAME_FIELD] = settings[POD_NAME_FIELD]
	data[NODE_NAME_FIELD] = settings[NODE_NAME_FIELD]
	data[TIME_FIELD] = entry.Time.Format(f.TimestampFormat)
	data[MESSAGE_FIELD] = entry.Message
	data[LEVEL_FIELD] = entry.Level.String()

	if _, ok := data[ACCESS_LOG_VERSION_FIELD]; ok {
		data[ACCESS_LOG_VERSION_FIELD] = settings[ACCESS_LOG_VERSION_FIELD]
		entry.Data[ACCESS_LOG_VERSION_FIELD] = settings[ACCESS_LOG_VERSION_FIELD]
	} else {
		data[STANDARD_LOG_VERSION_FIELD] = settings[STANDARD_LOG_VERSION_FIELD]
		entry.Data[STANDARD_LOG_VERSION_FIELD] = settings[STANDARD_LOG_VERSION_FIELD]
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
