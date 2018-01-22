package da

import (
	"encoding/json"
	"github.com/DroiTaipei/mgo/bson"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Do somethings before all test cases
func BeforeTest() {

}

func TestStoredCountsJsonEncoding(t *testing.T) {
	raw := []byte(`{"catalog":"m01","storedCounts":199}`)
	expected := StoredCounts{
		Catalog:      "m01",
		StoredCounts: 199,
		Stored:       nil,
	}
	s := StoredCounts{}
	json.Unmarshal(raw, &s)
	assert.Equal(t, expected, s)
}

func TestStoredJsonDecoding(t *testing.T) {
	raw := []byte(`{"catalog":"m05","stored":{"buy-food-eid":7,"rent-book-eid":9}}`)
	expected := StoredCounts{
		Catalog:      "m05",
		StoredCounts: 0,
		Stored: map[string]int{
			"buy-food-eid":  7,
			"rent-book-eid": 9,
		},
	}
	s := StoredCounts{}
	json.Unmarshal(raw, &s)
	assert.Equal(t, expected, s)
}

func TestSyncBeatJsonEncoding(t *testing.T) {
	expected := []byte(`{"counts":[{"catalog":"m01","storedCounts":199}],"version":"v2.6.978","os":1}`)
	s := SyncBeat{
		Counts: []StoredCounts{
			{
				Catalog:      "m01",
				StoredCounts: 199,
			},
		},
		Version: "v2.6.978",
		Os:      1,
	}
	result, _ := json.Marshal(s)
	assert.Equal(t, expected, result)
}

func TestSendCountsJsonEncoding(t *testing.T) {
	expected := []byte(`{"catalog":"m01","sendCounts":108}`)
	value := 108
	s := SendCounts{
		Catalog:    "m01",
		SendCounts: &value,
	}
	result, _ := json.Marshal(s)
	assert.Equal(t, expected, result)
}

func TestSendJsonEncoding(t *testing.T) {
	expected := []byte(`{"catalog":"m05","send":{"buy-food-eid":7,"rent-book-eid":9}}`)
	s := SendCounts{
		Catalog: "m05",
		Send: map[string]int{
			"buy-food-eid":  7,
			"rent-book-eid": 9,
		},
	}
	result, _ := json.Marshal(s)
	assert.Equal(t, expected, result)
}

func TestBeatResponseJsonEncoding(t *testing.T) {
	expected := []byte(`{"handle":[{"catalog":"m01","sendCounts":108}],"ipv4":"192.168.1.1","ipv6":"::1/128","nextPeriodSec":321,"version":"v2.6.978"}`)
	value := 108

	s := BeatResponse{
		Handle: []SendCounts{
			{
				Catalog:    "m01",
				SendCounts: &value,
			},
		},
		IPv4:          "192.168.1.1",
		IPv6:          "::1/128",
		NextPeriodSec: 321,
		Version:       "v2.6.978",
	}
	result, _ := json.Marshal(s)
	assert.Equal(t, expected, result)
}

func TestCatalogRecordsJsonEncoding(t *testing.T) {
	expected := []byte(`{"catalog":"m01","records":[{"s01":"XD"}]}`)
	s := CatalogRecords{
		Catalog: "m01",
		Records: json.RawMessage(`[{"s01":"XD"}]`),
	}
	result, _ := json.Marshal(s)
	assert.Equal(t, expected, result)
}

func TestDataRequestJsonEncoding(t *testing.T) {
	expected := []byte(`{"datas":[{"catalog":"m01","records":[{"s01":"XD"}]}],"ipv4":"192.168.1.1","ipv6":"::1/128","version":"v2.6.978","os":1}`)
	s := DataRequest{
		Datas: []CatalogRecords{
			{
				Catalog: "m01",
				Records: json.RawMessage(`[{"s01":"XD"}]`),
			},
		},
		IPv4:    "192.168.1.1",
		IPv6:    "::1/128",
		Version: "v2.6.978",
		Os:      1,
	}
	result, _ := json.Marshal(s)
	assert.Equal(t, expected, result)
}

func TestM01Map(t *testing.T) {
	expected := map[string]string{
		"s01": "imsi",
		"s02": "sim_msg_center",
		"s03": "sim_icc_id",
		"s04": "sim_phone_num",
		"s05": "sim_cell_id",
		"p01": "imei",
		"p02": "lcd",
		"p03": "model",
		"p04": "logo",
		"p05": "plf",
		"p06": "adv",
		"p07": "lang",
		"p08": "freeme_ver",
		"p09": "android_id",
		"p10": "ram_bytes",
		"p11": "rom_bytes",
		"p12": "breakout",
		"p13": "dalvik_art",
		"p14": "net_t",
		"a01": "pkg_name",
		"a02": "app_name",
		"a03": "app_id",
		"a04": "ch",
		"a05": "ver",
		"a06": "is_system_app",
		"did": "did",
		"a07": "app_sign",
		"a08": "coresdk_ver",
		"a09": "apm_ver",
		"st":  "st",
		"u01": "is_open_by_user",
	}
	assert.Equal(t, expected, M01Map)
}

func TestM01GenBsonM(t *testing.T) {
	expected := bson.M{
		"imsi":            "A",
		"sim_msg_center":  "B",
		"sim_icc_id":      "C",
		"sim_phone_num":   "D",
		"sim_cell_id":     "E",
		"imei":            "F",
		"lcd":             "G",
		"model":           "H",
		"logo":            "I",
		"plf":             "J",
		"adv":             "K",
		"lang":            "L",
		"freeme_ver":      "M",
		"android_id":      "N",
		"ram_bytes":       int64(16),
		"rom_bytes":       int64(4),
		"breakout":        true,
		"dalvik_art":      "O",
		"net_t":           "P",
		"pkg_name":        "Q",
		"app_name":        "R",
		"app_id":          "S",
		"ch":              "T",
		"ver":             "U",
		"is_system_app":   false,
		"did":             "V",
		"app_sign":        "W",
		"coresdk_ver":     "X",
		"apm_ver":         int64(10),
		"st":              int64(1257894000),
		"is_open_by_user": true,
		"pt":              "2017-09-20",
		"os":              1,
	}
	m := &M01{
		Imsi:         "A",
		SimMsgCenter: "B",
		SimIccID:     "C",
		SimPhoneNum:  "D",
		SimCellID:    "E",
		Imei:         "F",
		Lcd:          "G",
		Model:        "H",
		Logo:         "I",
		Plf:          "J",
		Adv:          "K",
		Lang:         "L",
		FreemeVer:    "M",
		AndroidID:    "N",
		RAMBytes:     int64(16),
		ROMBytes:     int64(4),
		Breakout:     true,
		DalvikArt:    "O",
		NetT:         "P",
		PkgName:      "Q",
		AppName:      "R",
		AppID:        "S",
		Ch:           "T",
		Ver:          "U",
		IsSystemApp:  false,
		DeviceID:     "V",
		AppSign:      "W",
		CoreSDKVer:   "X",
		ApmVer:       int64(10),
		Stamp:        1257894000,
		IsOpenByUser: true,
		DataDate:     "2017-09-20",
		Os:           1,
	}
	assert.Equal(t, expected, m.GenBsonM())
}

func TestM02GenBsonM(t *testing.T) {
	expected := bson.M{
		"last_in_st":    int64(1257891000),
		"last_out_st":   int64(1257893000),
		"st":            int64(1257894000),
		"sim_cell_id":   "A",
		"app_id":        "B",
		"ch":            "C",
		"ver":           "D",
		"is_system_app": true,
		"did":           "E",
		"net_t":         "F",
		"pt":            "2017-09-20",
		"os":            1,
	}
	m := &M02{
		LastInSt:    1257891000,
		LastOutSt:   1257893000,
		Stamp:       1257894000,
		SimCellID:   "A",
		AppID:       "B",
		Ch:          "C",
		Ver:         "D",
		IsSystemApp: true,
		DeviceID:    "E",
		NetT:        "F",
		DataDate:    "2017-09-20",
		Os:          1,
	}
	assert.Equal(t, expected, m.GenBsonM())
}

// Do somethings after all test cases
func AfterTest() {

}

func TestMain(m *testing.M) {
	BeforeTest()
	retCode := m.Run()
	AfterTest()
	os.Exit(retCode)
}
