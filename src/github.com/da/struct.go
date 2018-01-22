package da

import (
	"encoding/json"
	"github.com/DroiTaipei/mgo/bson"
	"reflect"
)

var M01Map map[string]string

func init() {
	var sM01 M01
	M01Map = genFieldMap(reflect.TypeOf(sM01))
}

type StoredCounts struct {
	Catalog      string         `json:"catalog"`
	StoredCounts int            `json:"storedCounts"`
	Stored       map[string]int `json:"stored,omitempty"`
}

type SyncBeat struct {
	Counts  []StoredCounts `json:"counts"`
	Version string         `json:"version"`
	Os      int            `json:"os"`
}

type SendCounts struct {
	Catalog    string         `json:"catalog"`
	SendCounts *int           `json:"sendCounts,omitempty"`
	Send       map[string]int `json:"send,omitempty"`
}

type BeatResponse struct {
	Handle        []SendCounts `json:"handle"`
	IPv4          string       `json:"ipv4"`
	IPv6          string       `json:"ipv6"`
	NextPeriodSec int          `json:"nextPeriodSec"`
	Version       string       `json:"version"`
}

type CatalogRecords struct {
	Catalog string          `json:"catalog"`
	Records json.RawMessage `json:"records"`
}

type DataRequest struct {
	Datas   []CatalogRecords `json:"datas"`
	IPv4    string           `json:"ipv4"`
	IPv6    string           `json:"ipv6"`
	Version string           `json:"version"`
	Os      int              `json:"os"`
}

func ignore(in string) bool {
	return len(in) == 0 || in == "-"
}

func genFieldMap(s reflect.Type) map[string]string {
	ret := make(map[string]string)

	b := s.NumField()
	for i := 0; i < b; i++ {
		field := s.Field(i)
		k := field.Tag.Get("json")
		v := field.Tag.Get("bson")
		if !ignore(k) && !ignore(v) {
			ret[k] = v
		}
	}
	return ret
}

// Ref : https://docs.google.com/spreadsheets/d/1JRvp2auivK3XP49CWJhAQYcnmSIolQVy4OzkKBweSa0/edit#gid=2114744419

type M01 struct {
	Imsi         string `json:"s01" bson:"imsi"`
	SimMsgCenter string `json:"s02" bson:"sim_msg_center"`
	SimIccID     string `json:"s03" bson:"sim_icc_id"`
	SimPhoneNum  string `json:"s04" bson:"sim_phone_num"`
	SimCellID    string `json:"s05" bson:"sim_cell_id"`
	Imei         string `json:"p01" bson:"imei"`
	Lcd          string `json:"p02" bson:"lcd"`
	Model        string `json:"p03" bson:"model"`
	Logo         string `json:"p04" bson:"logo"`
	Plf          string `json:"p05" bson:"plf"`
	Adv          string `json:"p06" bson:"adv"`
	Lang         string `json:"p07" bson:"lang"`
	FreemeVer    string `json:"p08" bson:"freeme_ver"`
	AndroidID    string `json:"p09" bson:"android_id"`
	RAMBytes     int64  `json:"p10" bson:"ram_bytes"`
	ROMBytes     int64  `json:"p11" bson:"rom_bytes"`
	Breakout     bool   `json:"p12" bson:"breakout"`
	DalvikArt    string `json:"p13" bson:"dalvik_art"`
	NetT         string `json:"p14" bson:"net_t"`
	PkgName      string `json:"a01" bson:"pkg_name"`
	AppName      string `json:"a02" bson:"app_name"`
	AppID        string `json:"a03" bson:"app_id"`
	Ch           string `json:"a04" bson:"ch"`
	Ver          string `json:"a05" bson:"ver"`
	IsSystemApp  bool   `json:"a06" bson:"is_system_app"`
	DeviceID     string `json:"did" bson:"did"`
	AppSign      string `json:"a07" bson:"app_sign"`
	CoreSDKVer   string `json:"a08" bson:"coresdk_ver"`
	ApmVer       int64  `json:"a09" bson:"apm_ver"`
	Stamp        int64  `json:"st" bson:"st"`
	IsOpenByUser bool   `json:"u01" bson:"is_open_by_user"`
	DataDate     string `bson:"pt"`
	Os           int    `bson:"os"`
}

func (m *M01) GenBsonM() (ret bson.M) {
	ret = bson.M{
		"imsi":            m.Imsi,
		"sim_msg_center":  m.SimMsgCenter,
		"sim_icc_id":      m.SimIccID,
		"sim_phone_num":   m.SimPhoneNum,
		"sim_cell_id":     m.SimCellID,
		"imei":            m.Imei,
		"lcd":             m.Lcd,
		"model":           m.Model,
		"logo":            m.Logo,
		"plf":             m.Plf,
		"adv":             m.Adv,
		"lang":            m.Lang,
		"freeme_ver":      m.FreemeVer,
		"android_id":      m.AndroidID,
		"ram_bytes":       m.RAMBytes,
		"rom_bytes":       m.ROMBytes,
		"breakout":        m.Breakout,
		"dalvik_art":      m.DalvikArt,
		"net_t":           m.NetT,
		"pkg_name":        m.PkgName,
		"app_name":        m.AppName,
		"app_id":          m.AppID,
		"ch":              m.Ch,
		"ver":             m.Ver,
		"is_system_app":   m.IsSystemApp,
		"did":             m.DeviceID,
		"app_sign":        m.AppSign,
		"coresdk_ver":     m.CoreSDKVer,
		"apm_ver":         m.ApmVer,
		"st":              m.Stamp,
		"is_open_by_user": m.IsOpenByUser,
		"pt":              m.DataDate,
		"os":              m.Os,
	}
	return
}

type M02 struct {
	LastInSt    int64  `json:"lst" bson:"last_in_st"`
	LastOutSt   int64  `json:"let" bson:"last_out_st"`
	Stamp       int64  `json:"st" bson:"st"`
	SimCellID   string `json:"s05" bson:"sim_cell_id"`
	AppID       string `json:"a03" bson:"app_id"`
	Ch          string `json:"a04" bson:"ch"`
	Ver         string `json:"a05" bson:"ver"`
	IsSystemApp bool   `json:"a06" bson:"is_system_app"`
	DeviceID    string `json:"did" bson:"did"`
	NetT        string `json:"p14" bson:"net_t"`
	DataDate    string `bson:"pt"`
	Os          int    `bson:"os"`
}

func (m *M02) GenBsonM() (ret bson.M) {
	ret = bson.M{
		"last_in_st":    m.LastInSt,
		"last_out_st":   m.LastOutSt,
		"st":            m.Stamp,
		"sim_cell_id":   m.SimCellID,
		"app_id":        m.AppID,
		"ch":            m.Ch,
		"ver":           m.Ver,
		"is_system_app": m.IsSystemApp,
		"did":           m.DeviceID,
		"net_t":         m.NetT,
		"pt":            m.DataDate,
		"os":            m.Os,
	}
	return
}
