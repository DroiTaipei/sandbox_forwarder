package postgres

import (
	"fmt"
	"time"
)

const (
	// DbUsageSticTable is the table for mongo colloection storage usage
	DbUsageSticTable = "db_usage_statistics"
	// AppUsageSticTable is the table for app monthly resource usage
	AppUsageSticTable = "app_usage_statistics"
)

// Mostly table in baas_statistics database is partition table
// With the format as
// Parent: {TableName}
// Childs: {TableName}_2017_01, {TableName}_2017_02, {TableName}_2017_03 ...
// Ideally, Any kind access target is Parent Table
// Use trigger the rearrange the data into Child Table(Partition)
// The Partition Key is date_part('year_month',access_time)
// We found while Insert(Update & Delete are wait to experiement)
// Parent Table will not return anything, including err
// We use gorm as ORM , while inserting , it add `RETURNING ID` at the end of INSERT SQL statement
// But Parent Table will return nothing ......, and gorm will return error
// So we try to access Child Table directly
func getParitionTable(parent string, partitionKey time.Time) string {
	if partitionKey.IsZero() {
		return parent
	}
	return fmt.Sprintf("%s_%s", parent, partitionKey.Format("2006_01"))
}

// DbUsageStic data structure
type DbUsageStic struct {
	ID             string    `json:"ID" gorm:"column:id"`
	AppID          string    `json:"appID" gorm:"column:application_id"`
	CollectionName string    `json:"collectionName"`
	Size           int64     `json:"size"`
	DocumentsCount int64     `json:"documentsCount"`
	AccessTime     time.Time `json:"accessTime"`
}

// TableName return the table name in database for gorm using
func (s DbUsageStic) TableName() string {
	return getParitionTable(DbUsageSticTable, s.AccessTime)
}

// AppUsageStic data structure
type AppUsageStic struct {
	ID           string    `json:"ID" gorm:"column:id"`
	AppID        string    `json:"appID" gorm:"column:application_id"`
	DeveloperID  string    `json:"developerID" gorm:"column:developer_id"`
	ReqCnt       int64     `json:"reqCnt" gorm:"column:request_count"`
	ReqIn        int64     `json:"reqIn" gorm:"column:request_throughput_in"`
	ReqOut       int64     `json:"reqOut" gorm:"column:request_throughput_out"`
	FileCdnOut   int64     `json:"fileCdnOut"  gorm:"column:file_cdn_throughput_out"`
	FileStorSize int64     `json:"fileStorSize"  gorm:"column:file_storage_size"`
	DbStorSize   int64     `json:"dbStorSize"  gorm:"column:db_storage_size"`
	AccessTime   time.Time `json:"accessTime"`
}

// TableName return the table name in database for gorm using
func (s AppUsageStic) TableName() string {
	return getParitionTable(AppUsageSticTable, s.AccessTime)
}
