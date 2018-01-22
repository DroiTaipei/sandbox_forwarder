package cloudcache

import (
	"fmt"
)

const (
	//DefaultTable is the default CloudCache data store table in VoltDB.
	DefaultTable = "TPE_D__R_OI__BAAS_CACHE_30000"
	//SandStage means represnt AppID is in sandbox stage.
	SandStage = "sand"
	//ProdStage means represnt AppID is in production stage
	ProdStage = "prod"
)

//AppRes is the struct for contain the resource constraints for App.
type AppRes struct {
	AppID string
	Stage string
}

var serviceAppRess map[string]*AppRes

//GetRowLimit is the getter for CloudCache RowLimit
func (ar *AppRes) GetRowLimit() int {
	return 30000
}

//GetCloudCacheTable to get the specified table for AppID
func (ar *AppRes) GetCloudCacheTable() string {
	tmpl := "TPE_D__R_OI__BAAS_CACHE__%s_%s"
	return fmt.Sprintf(tmpl, ar.AppID[0:8], ar.Stage)
}

//GenCloudCacheTableSQL to generate the SQL for create specified table for AppID
func (ar *AppRes) GenCloudCacheTableSQL() string {
	tbl := ar.GetCloudCacheTable()
	return fmt.Sprintf(`CREATE TABLE %s (
   KEY varchar(32) NOT NULL,
   VALUE varchar(1024) NOT NULL,
   PRIMARY KEY (KEY),
   LIMIT PARTITION ROWS %d
);
PARTITION TABLE %s ON COLUMN KEY;`, tbl, ar.GetRowLimit(), tbl)
}

func init() {
	serviceAppRess = map[string]*AppRes{
		// 天氣服務
		"mfpvmbzh-hZEHSCm90pXoWA2DUfRqB0YlQC4qngk": {
			AppID: "mfpvmbzh-hZEHSCm90pXoWA2DUfRqB0YlQC4qngk",
			Stage: SandStage,
		},
		"mfpvmbzhM9iqko8nk2u5klmoBREpPQNElQCYlSQN": {
			AppID: "mfpvmbzhM9iqko8nk2u5klmoBREpPQNElQCYlSQN",
			Stage: ProdStage,
		},
		// 微創海外信息流
		"of2umbzhzQunGpafjz0_rW25FgMCsvIHlQCAdCwJ": {
			AppID: "of2umbzhzQunGpafjz0_rW25FgMCsvIHlQCAdCwJ",
			Stage: SandStage,
		},
		"of2umbzhoAcZ8_DnHF95NP-1SLhCdahYlQDIVxMA": {
			AppID: "of2umbzhoAcZ8_DnHF95NP-1SLhCdahYlQDIVxMA",
			Stage: ProdStage,
		},
	}
}

//GetCloudCacheTable is exported function for externally getting Table
func GetCloudCacheTable(appID string) string {
	AppRes, ok := serviceAppRess[appID]
	if !ok {
		return DefaultTable
	}
	return AppRes.GetCloudCacheTable()
}
