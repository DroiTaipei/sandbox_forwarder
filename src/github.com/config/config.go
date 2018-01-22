package config

import (
	"errors"
	gcfg "github.com/DroiTaipei/go-config"
	"sort"
	"strings"
	"time"
)

type Config struct {
	*gcfg.Config
	Opt *Options
}

// GetUniqSubKeys - Using prefix to get the array with number suffix
func (cfgs *Config) GetUniqSubKeys(prefix string) []string {
	kv, _ := cfgs.Settings()
	keys := []string{}
	for k := range kv {
		if strings.HasPrefix(k, prefix) {
			tokens := strings.Split(strings.TrimPrefix(k, prefix), ".")
			subKey := tokens[0] // I only interest in the first level subkeys
			dup := false
			for _, key := range keys {
				if subKey == key {
					dup = true
					break
				}
			}
			if !dup {
				keys = append(keys, subKey)
			}
		}
	}
	// Add sort to keep output stable
	sort.Strings(keys)
	return keys
}

// Duration get the time.Duration
// Follow the convention of go-config
func (cfgs *Config) Duration(key string, alt time.Duration) (ret time.Duration) {
	ret = alt
	raw, err := cfgs.String(key)
	if err != nil {
		return
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return
	}
	ret = d
	return
}

type Validater func(map[string]string) error

type Options struct {
	Validate Validater
	EnvMap   map[string]string
}

func NewOptions() *Options {
	return &Options{
		EnvMap: make(map[string]string),
	}
}

func keyToEnv(in string) string {
	return strings.ToUpper(strings.Replace(in, ".", "_", 1))
}

func defaultEnvMap(p gcfg.Provider) (ret map[string]string, err error) {
	ret = make(map[string]string)
	s, err := p.Load()
	if err != nil {
		return
	}
	for k := range s {
		ret[keyToEnv(k)] = k
	}
	return
}

func NewStaticConfig(args map[string]string) (ret *Config) {
	ret, _ = NewConfig(gcfg.NewOnceLoader(gcfg.NewStatic(args)), nil)
	return ret
}

func NewConfig(p gcfg.Provider, opt *Options) (ret *Config, err error) {
	if p == nil {
		err = errors.New("Nil gcfg.Provider")
		return
	}
	if opt == nil {
		opt = NewOptions()
	}
	envMap, err := defaultEnvMap(p)
	if err != nil {
		return
	}

	for k, v := range opt.EnvMap {
		envMap[k] = v
	}
	opt.EnvMap = envMap
	env := gcfg.NewEnvironment(envMap)
	cfgs := gcfg.NewConfig([]gcfg.Provider{p, env})
	cfgs.Validate = opt.Validate
	err = cfgs.Load()
	if err != nil {
		return
	}
	ret = &Config{
		cfgs,
		opt,
	}
	return
}
func LoadConfig(configFile string, opt *Options) (ret *Config, err error) {
	oLoader := gcfg.NewOnceLoader(gcfg.NewTOMLFile(configFile))
	return NewConfig(oLoader, opt)
}
