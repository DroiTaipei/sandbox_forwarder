package dlogrus

import (
	"encoding/binary"
	"fmt"
	"github.com/DroiTaipei/mgo/bson"
	"github.com/devopstaku/logrus"
)

const (
	AID_LENGTH         = 40
	RID_LENGTH         = 64
	DB_TYPE_LENGTH     = 6
	QUERY_SPENT_LENGTH = 4
	TIMEFORMAT_LENGTH  = 26
	XPK_LENGTH         = 40
	DCC_LENGTH         = 20
	PREFIX_LENGTH      = AID_LENGTH + RID_LENGTH + DB_TYPE_LENGTH +
		QUERY_SPENT_LENGTH + TIMEFORMAT_LENGTH + XPK_LENGTH + DCC_LENGTH
	LOG_KEY_APP_ID        = "Aid"
	LOG_KEY_REQ_ID        = "Rid"
	LOG_KEY_PLATFORM_KEY  = "XPk"
	DB_TYPE_FIELD         = "Dt"
	DB_HOSTNAME_FIELD     = "Dh"
	DB_COMMAND_TIME_FIELD = "Dct"
	DB_COMMAND_TYPE_FIELD = "Dcc"
)

type BSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

func newBSONFormatter(timestampFormat string) *BSONFormatter {
	if timestampFormat == "" {
		timestampFormat = DEFAULT_TIMESTAMP_FORMAT
	}
	return &BSONFormatter{TimestampFormat: timestampFormat}
}

func (f *BSONFormatter) prefix(d logrus.Fields) ([]byte, error) {
	ret := make([]byte, PREFIX_LENGTH)
	needle := 0
	if v, ok := d[LOG_KEY_APP_ID].(string); ok {
		copy(ret[needle:], []byte(v))
	}
	needle += AID_LENGTH
	if v, ok := d[LOG_KEY_REQ_ID].(string); ok {
		copy(ret[needle:], []byte(v))
	}
	needle += RID_LENGTH
	if v, ok := d[LOG_KEY_PLATFORM_KEY].(string); ok {
		copy(ret[needle:], []byte(v))
	}
	needle += XPK_LENGTH
	if v, ok := d[DB_TYPE_FIELD].(string); ok {
		copy(ret[needle:], []byte(v))
	}
	needle += DB_TYPE_LENGTH
	if v, ok := d[DB_COMMAND_TYPE_FIELD].(string); ok {
		copy(ret[needle:], []byte(v))
	}
	needle += DCC_LENGTH
	if v, ok := d[DB_COMMAND_TIME_FIELD].(int64); ok {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v))
		copy(ret[needle:], b)
	}
	needle += QUERY_SPENT_LENGTH
	if v, ok := d[TIME_FIELD].(string); ok {
		copy(ret[needle:], []byte(v))
	}
	needle += TIMEFORMAT_LENGTH
	return ret, nil
}

func (f *BSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
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
	_, isAccessLog := data[ACCESS_LOG_VERSION_FIELD]

	data[TIME_FIELD] = entry.Time.Format(f.TimestampFormat)
	data[MESSAGE_FIELD] = entry.Message
	data[LEVEL_FIELD] = entry.Level.String()
	data[POD_NAME_FIELD] = settings[POD_NAME_FIELD]
	data[NODE_NAME_FIELD] = settings[NODE_NAME_FIELD]

	if isAccessLog {
		data[ACCESS_LOG_VERSION_FIELD] = settings[ACCESS_LOG_VERSION_FIELD]
		entry.Data[ACCESS_LOG_VERSION_FIELD] = settings[ACCESS_LOG_VERSION_FIELD]
	} else {
		data[STANDARD_LOG_VERSION_FIELD] = settings[STANDARD_LOG_VERSION_FIELD]
		entry.Data[STANDARD_LOG_VERSION_FIELD] = settings[STANDARD_LOG_VERSION_FIELD]
	}

	serialized, err := bson.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to BSON, %v", err)
	}

	if isAccessLog {
		prefix, err := f.prefix(data)
		if err != nil {
			return nil, err
		}
		return append(prefix, serialized...), nil
	}

	return serialized, nil
}
