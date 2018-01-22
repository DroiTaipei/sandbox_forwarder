package rdb

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Do somethings before all test cases
func BeforeTest() {

}

func TestJsonEncoding(t *testing.T) {
	expected := []byte(`{"ID":{"$like":"ABC%"},"status":1}`)
	s := QueryAppPrefixPayload{
		ID: LikeCriterion{
			Like: "ABC%",
		},
		Status: 1,
	}
	result, _ := json.Marshal(s)
	assert.Equal(t, expected, result)
}

func TestGetStage(t *testing.T) {
	s := Application{
		StageFlag: 0,
	}
	assert.Equal(t, "Unknown", s.GetStage())
	s.StageFlag = 1
	assert.Equal(t, "Sandbox", s.GetStage())
	s.StageFlag = 2
	assert.Equal(t, "Production", s.GetStage())

}

func TestGetStatus(t *testing.T) {
	s := Application{
		Status: 0,
	}
	assert.Equal(t, "Unknown", s.GetStatus())
	s.Status = 1
	assert.Equal(t, "Valid", s.GetStatus())
	s.Status = 2
	assert.Equal(t, "Deleted", s.GetStatus())
	s.Status = 3
	assert.Equal(t, "Suspended", s.GetStatus())
}

// Because there is a special cases for upload file data structure
func TestUploadFileFieldMap(t *testing.T) {
	s := UploadFile{}
	result := map[string]string{
		"appID":            "app_id",
		"fileID":           "fid_raw",
		"path":             "path",
		"type":             "type",
		"size":             "size",
		"modifyTime":       "modify_time",
		"createdTime":      "created_time",
		"fileDescription":  "file_description",
		"MD5":              "md5",
		"CDN":              "cdn",
		"CDNMap":           "cdn_map",
		"status":           "status",
		"statusUpdateTime": "status_update_time",
		"objectID":         "object_id",
		"callerAppID":      "caller_app_id",
		"scanDescription":  "scan_description",
		"scanStatus":       "scan_status",
		"scanTime":         "scan_time",
	}
	assert.EqualValues(t, result, s.FieldMap())
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
