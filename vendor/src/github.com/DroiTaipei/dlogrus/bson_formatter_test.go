package dlogrus

import (
	"encoding/binary"
	"testing"

	"github.com/DroiTaipei/mgo/bson"
	"github.com/devopstaku/logrus"
)

func TestAccessLogUnserialize(t *testing.T) {
	Initialize(map[string]string{"standrd_log_version": "1", "access_log_version": "1"})
	formatter := newBSONFormatter("")
	data := map[string]interface{}{"Dct": int64(5), "A": "1", "Aid": "sf8umbzhPSbyjbHi0foaOt9KttUgRv7hlQBYJQAA", "Dt": "voltdb", "Rid": "1f8umbzhObFZyDl9d_M9EyqM1swAJVAkiz-84QEA1f8umbzhObFZyDl9d84QEA1K", "Aidm": "prod", "Dh": "voltdb02", "Dc": "SELECT * FROM oz_topic  ORDER BY topic_id ASC OFFSET 0 LIMIT 10"}
	b, err := formatter.Format(logrus.WithFields(logrus.Fields(data)))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	prefix := b[:PREFIX_LENGTH]
	needle := 0
	if string(prefix[needle:needle+AID_LENGTH]) != "sf8umbzhPSbyjbHi0foaOt9KttUgRv7hlQBYJQAA" {
		t.Fatal("AID is not equal ", err)
	}
	needle += AID_LENGTH
	if string(prefix[needle:needle+RID_LENGTH]) != "1f8umbzhObFZyDl9d_M9EyqM1swAJVAkiz-84QEA1f8umbzhObFZyDl9d84QEA1K" {
		t.Fatal("RID is not equal ", err)
	}
	needle += RID_LENGTH
	if string(prefix[needle:needle+DB_TYPE_LENGTH]) != "voltdb" {
		t.Fatal("DBT is not equal ", err)
	}
	needle += DB_TYPE_LENGTH
	if binary.LittleEndian.Uint32(prefix[needle:needle+QUERY_SPENT_LENGTH]) != uint32(5) {
		t.Fatal("DCT is not equal ", err)
	}
	needle += QUERY_SPENT_LENGTH

	entry := make(map[string]interface{})
	err = bson.Unmarshal(b[PREFIX_LENGTH:], &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	for k, v := range data {
		if entry[k] != v {
			t.Fatal(k, " field not equal")
		}
	}

}
