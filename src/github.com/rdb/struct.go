package rdb

import (
	"reflect"
	"strings"
	"time"
)

const (
	// unknown stage or status
	unknown = "Unknown"

	// StageFlagSand Sandbox Value in database
	StageFlagSand = 1
	// StageFlagProd Production Value in database
	StageFlagProd = 2

	// StatusFlagValid Value in database
	StatusFlagValid = 1
	// StatusFlagDeleted Value in database
	StatusFlagDeleted = 2
	// StatusFlagSuspended Value in database
	StatusFlagSuspended = 3
)

var (
	appStage = map[int]string{
		StageFlagSand: "Sandbox",
		StageFlagProd: "Production",
	}
	appStatus = map[int]string{
		StatusFlagValid:     "Valid",
		StatusFlagDeleted:   "Deleted",
		StatusFlagSuspended: "Suspended",
	}
)

// Application data structure
type Application struct {
	ID              string    `json:"ID" gorm:"primary_key"`
	Name            string    `json:"name"`
	PackageName     string    `json:"packageName"`
	DeveloperID     string    `json:"developerID"`
	URL             string    `json:"URL"`
	Description     string    `json:"description"`
	Icon            string    `json:"icon"`
	ClientKey       string    `json:"clientKey"`
	RestAPIKey      string    `json:"restApiKey"`
	CloudCodeKey    string    `json:"cloudCodeKey"`
	MasterKey       string    `json:"masterKey"`
	MasterKeyCount  int       `json:"masterKeyCount"`
	SecretKey       string    `json:"secretKey"`
	Preference      string    `json:"preference"`
	AuthPublicData  string    `json:"authPublicData"`
	AuthPrivateData string    `json:"authPrivateData"`
	CreatedTime     time.Time `gorm:"column:creation_time" json:"createdTime"`
	Status          int       `json:"status"`
	RunningVer      string    `json:"runningVer"`
	WebConfig       string    `json:"webConfig"`
	StageFlag       int       `json:"stageFlag"`
	QiniuAccount    string    `json:"qiniuAccount"`
	CRC32AppID      int       `gorm:"column:crc32_app_id"`
	ModRID          int       `gorm:"column:mod_r_id"`
}

// UpdateApplication - update application record on restricted fields. used for CloudOps API now
type UpdateApplication struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Preference  string `json:"preference,omitempty"`
	Status      int    `json:"status,omitempty"`
}

// TableName return the table name in database for gorm using
func (s Application) TableName() string {
	return "baas.application"
}

// FieldMap return the map between json field name and db field name
func (s *Application) FieldMap() map[string]string {
	return map[string]string{
		"name":            "name",
		"packageName":     "package_name",
		"developerID":     "developer_id",
		"URL":             "url",
		"description":     "description",
		"icon":            "icon",
		"clientKey":       "client_key",
		"restApiKey":      "rest_api_key",
		"cloudCodeKey":    "cloud_code_key",
		"masterKey":       "master_key",
		"masterKeyCount":  "master_key_count",
		"secretKey":       "secret_key",
		"preference":      "preference",
		"authPublicData":  "auth_public_data",
		"authPrivateData": "auth_private_data",
		"createdTime":     "creation_time",
		"status":          "status",
		"runningVer":      "running_ver",
		"webConfig":       "web_config",
		"stageFlag":       "stage_flag",
		"qiniuAccount":    "qiniu_account",
	}
}

// GetStage return the stage of app
func (s *Application) GetStage() string {
	return GetStageName(s.StageFlag)
}

// GetStatus return the status of app
func (s *Application) GetStatus() string {
	return GetStatusName(s.Status)
}

// GetStageName return the stage string representation
func GetStageName(stage int) string {
	if v, ok := appStage[stage]; ok {
		return v
	}
	return unknown
}

// GetStatusName return the status string representation
func GetStatusName(status int) string {
	if v, ok := appStatus[status]; ok {
		return v
	}
	return unknown
}

// QueryAppPrefixPayload - payload sent to Get Sandbox/Production App with App Prefix
type QueryAppPrefixPayload struct {
	ID     LikeCriterion `json:"ID"`
	Status int           `json:"status"`
}

// LikeCriterion - For using fuzzy search
type LikeCriterion struct {
	Like string `json:"$like"`
}

// UploadFile data structure
type UploadFile struct {
	AppID            string    `gorm:"column:app_id" json:"appID"`
	FIDRaw           int64     `gorm:"column:fid_raw" json:"-"`
	FID              string    `gorm:"column:fid" json:"fileID"`
	Path             string    `gorm:"column:path" json:"path"`
	Type             string    `gorm:"column:type" json:"type"`
	Size             int64     `gorm:"column:size" json:"size"`
	ModifyTime       time.Time `gorm:"column:modify_time" json:"modifyTime"`
	CreatedTime      time.Time `gorm:"column:created_time" json:"createdTime"`
	FileDescription  string    `gorm:"column:file_description" json:"fileDescription"`
	MD5              string    `gorm:"column:md5" json:"MD5"`
	CDN              string    `gorm:"column:cdn" json:"CDN"`
	CDNMap           string    `gorm:"column:cdn_map" json:"CDNMap"`
	Status           int64     `gorm:"column:status" json:"status"`
	StatusUpdateTime int64     `gorm:"column:status_update_time" json:"statusUpdateTime"`
	ObjectID         string    `gorm:"column:object_id" json:"objectID"`
	CRC32AppID       int       `gorm:"column:crc32_app_id" json:"-"`
	ModRID           int       `gorm:"column:mod_r_id" json:"-"`
	CallerAppID      string    `gorm:"column:caller_app_id" json:"callerAppID"`
	ScanStatus       int       `gorm:"column:scan_status" json:"scanStatus"`
	ScanTime         time.Time `gorm:"column:scan_time" json:"scanTime"`
	ScanDescription  string    `gorm:"column:scan_description" json:"scanDescription"`
}

// UpdateFile data structure
type UpdateFile struct {
	ModifyTime       string  `gorm:"column:modify_time" json:"modifyTime,omitempty"`
	StatusUpdateTime int64   `gorm:"column:status_update_time" json:"statusUpdateTime,omitempty"`
	ObjectID         *string `gorm:"column:object_id" json:"objectID,omitempty"`
	CDN              *string `gorm:"column:cdn" json:"CDN,omitempty"`
	CDNMap           *string `gorm:"column:cdn_map" json:"CDNMap,omitempty"`
	Status           int     `gorm:"column:status" json:"status,omitempty"`
	ScanStatus       int     `gorm:"column:scan_status" json:"scanStatus,omitempty"`
	ScanTime         string  `gorm:"column:scan_time" json:"scanTime,omitempty"`
	ScanDescription  string  `gorm:"column:scan_description" json:"scanDescription,omitempty"`
}

// TableName return the table name in database for gorm using
func (s UploadFile) TableName() string {
	return "baas.upload_file_mod50"
}

func parseTagSetting(tags reflect.StructTag) map[string]string {
	setting := map[string]string{}
	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}

func ignore(in string) bool {
	return len(in) == 0 || in == "-"
}

func genFieldMap(s reflect.Type) map[string]string {
	ret := make(map[string]string)

	b := s.NumField()
	for i := 0; i < b; i++ {
		field := s.Field(i)
		gSetting := parseTagSetting(field.Tag)
		k := field.Tag.Get("json")
		v := gSetting["COLUMN"]
		if !ignore(k) && !ignore(v) {
			ret[k] = v
		}
	}
	return ret
}

// FieldMap return the map between json field name and db field name
func (s UploadFile) FieldMap() map[string]string {
	ret := genFieldMap(reflect.TypeOf(s))
	// Special Cases
	ret["fileID"] = "fid_raw"
	return ret
}

// BulkCreate data structure
type BulkCreateFile struct {
	Files []UploadFile `json:"files"`
}

// BulkUpdate data structure
type BulkUpdateFile struct {
	Files []map[string]interface{} `json:"files"`
}
