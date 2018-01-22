package config

import (
	"github.com/DroiTaipei/droipkg"
	gcfg "github.com/DroiTaipei/go-config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func BeforeTest() {

}

func TestKeyToEnv(t *testing.T) {
	source := "log.kafka_enabled"
	target := "LOG_KAFKA_ENABLED"
	assert.Equal(t, target, keyToEnv(source))
}

func TestDefaultEnvMap(t *testing.T) {
	args := map[string]string{
		"a.b_c_d": "NoUsedValueInThisCase",
	}
	expected := map[string]string{
		"A_B_C_D": "a.b_c_d",
	}
	p := gcfg.NewOnceLoader(gcfg.NewStatic(args))
	result, _ := defaultEnvMap(p)
	assert.EqualValues(t, expected, result)
}

func TestEnvMapOveride(t *testing.T) {
	os.Setenv("A_B_C_D", "123")
	args := map[string]string{
		"a.b_c_d": "NoUsedValueInThisCase",
	}
	expected := map[string]string{
		"A_B_C_D": "a_b.c_d",
	}
	p := gcfg.NewOnceLoader(gcfg.NewStatic(args))
	opt := &Options{
		EnvMap: expected,
	}
	c, err := NewConfig(p, opt)
	assert.Nil(t, err)
	r, err := c.String("a.b_c_d")
	assert.Nil(t, err)
	assert.Equal(t, "NoUsedValueInThisCase", r)
	r, err = c.String("a_b.c_d")
	assert.Nil(t, err)
	assert.Equal(t, "123", r)
	assert.EqualValues(t, expected, c.Opt.EnvMap)
}

func TestEnvMapMerge(t *testing.T) {
	args := map[string]string{
		"a.b_c_d": "NoUsedValueInThisCase",
	}
	append := map[string]string{
		"B_C_D_E": "b.c_d_E",
	}
	expected := map[string]string{
		"A_B_C_D": "a.b_c_d",
		"B_C_D_E": "b.c_d_E",
	}
	p := gcfg.NewOnceLoader(gcfg.NewStatic(args))
	opt := &Options{
		EnvMap: append,
	}
	c, err := NewConfig(p, opt)
	assert.Nil(t, err)
	assert.EqualValues(t, expected, c.Opt.EnvMap)
}

func TestGetUniqSubKeys(t *testing.T) {
	args := map[string]string{
		"a.b.c.1": "NoUsedValueInThisCase",
		"a.b.c.2": "NoUsedValueInThisCase",
		"a.b.c.3": "NoUsedValueInThisCase",
		"a.b.c_4": "NoUsedValueInThisCase",
	}

	expected := []string{
		"1",
		"2",
		"3",
	}

	p := gcfg.NewOnceLoader(gcfg.NewStatic(args))
	c, err := NewConfig(p, nil)
	assert.Nil(t, err)
	assert.EqualValues(t, expected, c.GetUniqSubKeys("a.b.c."))
}

func TestNewConfig(t *testing.T) {
	var p gcfg.Provider
	c, err := NewConfig(p, nil)
	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func TestValidate(t *testing.T) {
	args := map[string]string{
		"a.b_c_d": "NoUsedValueInThisCase",
	}
	validate := func(in map[string]string) error {
		if _, ok := in["a.b"]; !ok {
			return droipkg.NewError("a.b not found")
		}
		return nil
	}
	p := gcfg.NewOnceLoader(gcfg.NewStatic(args))
	opt := &Options{
		Validate: validate,
	}
	_, err := NewConfig(p, opt)
	assert.NotNil(t, err)
}

func TestNewOptions(t *testing.T) {
	opt := NewOptions()
	assert.NotNil(t, opt)
	assert.NotNil(t, opt.EnvMap)
	assert.Nil(t, opt.Validate)
}

func TestDuration(t *testing.T) {
	args := map[string]string{
		"test.duration": "57s",
	}

	expected := 57 * time.Second
	c := NewStaticConfig(args)
	assert.Equal(t, expected, c.Duration("test.duration", 0*time.Second))
	expected = 21 * time.Minute
	assert.Equal(t, expected, c.Duration("not.exists", expected))
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
