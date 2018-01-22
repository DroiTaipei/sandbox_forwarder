package v0

import (
	"fmt"
	"strings"
	"time"
)

type ToDBNamer func(string) string

type Application struct {
	Id              string `gorm:"primary_key"`
	Name            string
	PackageName     string
	DeveloperId     string
	Url             string
	Description     string
	Icon            string
	ClientKey       string
	RestApiKey      string
	CloudCodeKey    string
	MasterKey       string
	MasterKeyCount  int
	SecretKey       string
	Preference      string
	AuthPublicData  string
	AuthPrivateData string
	CreationTime    time.Time
	Status          int
	RunningVer      string
	WebConfig       string
	StageFlag       int
	QiniuAccount    string
	CRC32AppID      int `gorm:"column:crc32_app_id"`
	ModRID          int `gorm:"column:mod_r_id"`
}

func (s Application) TableName() string {
	return "baas.application"
}

func (s *Application) ParseOrder(conv ToDBNamer, input string) string {
	t := strings.Split(input, ",")
	var fieldName, direction string
	var r []string
	b := len(t)
	for i := 0; i < b; i++ {
		direction = "ASC"
		fieldName = t[i]

		if strings.Index(fieldName, "-") == 0 {
			direction = "DESC"
			fieldName = fieldName[1:]
		}
		if s.checkValidField(fieldName) {
			r = append(r, fmt.Sprintf("%s %s", conv(fieldName), direction))
		}

	}

	return strings.Join(r, ", ")
}

func (*Application) checkValidField(input string) bool {
	validFields := []string{
		"Id",
		"Name",
		"PackageName",
		"DeveloperId",
		"Url",
		"Description",
		"Icon",
		"ClientKey",
		"WebKey",
		"NetKey",
		"Preference",
	}
	b := len(validFields)
	for i := 0; i < b; i++ {
		if input == validFields[i] {
			return true
		}
	}
	return len(input) > 0
}

type Developer struct {
	Id                string `gorm:"primary_key"`
	RegisterationDate time.Time
	Activated         bool
	CreatingAppCount  int
	DeveloperToken    string
}

func (s Developer) TableName() string {
	return "baas.developer"
}

func (s *Developer) ParseOrder(conv ToDBNamer, input string) string {
	t := strings.Split(input, ",")
	var fieldName, direction string
	var r []string
	b := len(t)
	for i := 0; i < b; i++ {
		direction = "ASC"
		fieldName = t[i]

		if strings.Index(fieldName, "-") == 0 {
			direction = "DESC"
			fieldName = fieldName[1:]
		}
		if s.checkValidField(fieldName) {
			r = append(r, fmt.Sprintf("%s %s", conv(fieldName), direction))
		}

	}

	return strings.Join(r, ", ")
}

// FIXME : For quickly allow JOIN , mofify this
func (*Developer) checkValidField(input string) bool {
	validFields := []string{
		"Id",
		"Name",
		"Email",
		"CreateTime",
		"ModifyTime",
	}
	b := len(validFields)
	for i := 0; i < b; i++ {
		if input == validFields[i] {
			return true
		}
	}
	return len(input) > 0
}

//Team is about schema of baas.team
type Team struct {
	ApplicationId string
	DeveloperId   string
}

func (s Team) TableName() string {
	return "baas.team"
}

func (s *Team) ParseOrder(conv ToDBNamer, input string) string {
	t := strings.Split(input, ",")
	var fieldName, direction string
	var r []string
	b := len(t)
	for i := 0; i < b; i++ {
		direction = "ASC"
		fieldName = t[i]

		if strings.Index(fieldName, "-") == 0 {
			direction = "DESC"
			fieldName = fieldName[1:]
		}
		if s.checkValidField(fieldName) {
			r = append(r, fmt.Sprintf("%s %s", conv(fieldName), direction))
		}

	}

	return strings.Join(r, ", ")
}

// FIXME : For quickly allow JOIN , mofify this
func (*Team) checkValidField(input string) bool {
	validFields := []string{
		"ApplicationId",
		"DeveloperId",
	}
	b := len(validFields)
	for i := 0; i < b; i++ {
		if input == validFields[i] {
			return true
		}
	}
	return len(input) > 0
}

type File struct {
	AppID            string `gorm:"column:app_id"`
	FIDRaw           int64  `gorm:"column:fid_raw"`
	FID              string `gorm:"column:fid"`
	UID              string `gorm:"column:uid"`
	Path             string
	Type             string
	Size             int64
	ModifyTime       time.Time
	CreatedTime      time.Time
	FileDescription  string
	MD5              string `gorm:"column:md5"`
	CDN              string `gorm:"column:cdn"`
	CDNMap           string `gorm:"column:cdn_map"`
	DfsGroup         int
	DfsPath          string
	Status           int64
	StatusUpdateTime int64
	CRC32AppID       int `gorm:"column:crc32_app_id" json:"-"`
	ModRID           int `gorm:"column:mod_r_id" json:"-"`
}

func (s File) TableName() string {
	return "baas.upload_file_mod50"
}

func (s File) FieldMap() map[string]string {
	return map[string]string{
		"AppID":            "app_id",
		"FID":              "fid",
		"UID":              "uid",
		"Path":             "path",
		"Type":             "type",
		"Size":             "size",
		"ModifyTime":       "modify_time",
		"CreatedTime":      "created_time",
		"FileDescription":  "file_description",
		"MD5":              "md5",
		"CDN":              "cdn",
		"CDNMap":           "cdn_map",
		"DfsGroup":         "dfs_group",
		"DfsPath":          "dfs_path",
		"Status":           "status",
		"StatusUpdateTime": "status_update_time",
	}
}
