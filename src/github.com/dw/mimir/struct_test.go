package mimir

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Do somethings before all test cases
func BeforeTest() {

}

func TestGenBody(t *testing.T) {
	assert := assert.New(t)
	fa := &FormArgs{}
	assert.Equal(0, len(fa.GenBody()))

	fa = &FormArgs{
		AppID:   "xxx",
		Channel: "yyyy",
	}
	expected := []byte("app_id=xxx&channel=yyyy")
	assert.Equal(expected, fa.GenBody())
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
