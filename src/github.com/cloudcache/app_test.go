package cloudcache

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	sandAppRes *AppRes
	prodAppRes *AppRes
)

func BeforeTest() {
	sandAppRes = &AppRes{
		AppID: "mfpvmbzh-hZEHSCm90pXoWA2DUfRqB0YlQC4qngk",
		Stage: SandStage,
	}
	prodAppRes = &AppRes{
		AppID: "mfpvmbzhM9iqko8nk2u5klmoBREpPQNElQCYlSQN",
		Stage: ProdStage,
	}
}

func TestGetTable(t *testing.T) {
	assertEqual(t, GetCloudCacheTable("mfpvmbzh-hZEHSCm90pXoWA2DUfRqB0YlQC4qngk"), "TPE_D__R_OI__BAAS_CACHE__mfpvmbzh_sand")
	assertEqual(t, GetCloudCacheTable("mfpvmbzhM9iqko8nk2u5klmoBREpPQNElQCYlSQN"), "TPE_D__R_OI__BAAS_CACHE__mfpvmbzh_prod")
	assertEqual(t, GetCloudCacheTable("Unknown"), "TPE_D__R_OI__BAAS_CACHE_30000")
}

func TestAppResGetTable(t *testing.T) {
	assertEqual(t, sandAppRes.GetCloudCacheTable(), "TPE_D__R_OI__BAAS_CACHE__mfpvmbzh_sand")
	assertEqual(t, prodAppRes.GetCloudCacheTable(), "TPE_D__R_OI__BAAS_CACHE__mfpvmbzh_prod")
}

func TestGenTableSQL(t *testing.T) {
	expected := `CREATE TABLE TPE_D__R_OI__BAAS_CACHE__mfpvmbzh_sand (
   KEY varchar(32) NOT NULL,
   VALUE varchar(1024) NOT NULL,
   PRIMARY KEY (KEY),
   LIMIT PARTITION ROWS 30000
);
PARTITION TABLE TPE_D__R_OI__BAAS_CACHE__mfpvmbzh_sand ON COLUMN KEY;`
	assertEqual(t, sandAppRes.GenCloudCacheTableSQL(), expected)
	expected = `CREATE TABLE TPE_D__R_OI__BAAS_CACHE__mfpvmbzh_prod (
   KEY varchar(32) NOT NULL,
   VALUE varchar(1024) NOT NULL,
   PRIMARY KEY (KEY),
   LIMIT PARTITION ROWS 30000
);
PARTITION TABLE TPE_D__R_OI__BAAS_CACHE__mfpvmbzh_prod ON COLUMN KEY;`
	assertEqual(t, prodAppRes.GenCloudCacheTableSQL(), expected)
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}, message ...string) {
	if expected == actual {
		return
	}
	message = append(message, fmt.Sprintf("\nExpected: %v\nreceived: %v", expected, actual))
	t.Fatal(strings.Join(message, " "))
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
