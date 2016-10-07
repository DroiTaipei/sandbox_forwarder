package config

import (
	"errors"
	"fmt"
	"github.com/DroiTaipei/dlogrus"
	"github.com/DroiTaipei/mongo"

	gconfig "github.com/zpatrick/go-config"
	"strings"
	"time"
)

const (
	EMPTY        = ""
	API_PORT     = 8080
	API_POD_NAME = "db-api"

	MGO_PORT         = "7379"
	MGO_MAX_CONN     = 250
	MGO_MAX_IDLE     = 240
	MGO_USER         = "bass"
	MGO_DATABASE     = "baas"
	MGO_AUTHDATABASE = "admin"
	MGO_TIMEOUT      = 300
	MGO_SHARDUSER    = "sharder"

	LOG_FILE_NAME = "db_api.log"
	LOG_VERSION   = "1"
	LOG_LEVEL     = "debug"
	LOG_FORMATTER = "json"

	KAFKA_ENABLED              = false
	FLUSH_FREQUENCY            = 500
	MAX_CONNECTIONS            = 5
	MAX_CONNECTIONS_PER_BROKER = 5
	BATCH_SIZE                 = 16384
	MAX_REQUESTS               = 10
	SEND_ROUTINES              = 10
	RECEIVE_ROUTINES           = 10
	REQUIRED_ACKS              = 1
	MAX_BATCH_SIZE             = 16384
	QUEUE_LENGTH               = 1024
	ENQUEUE_TIMEOUT            = 1000
)

type Config struct {
	*gconfig.Config
}

func (cfgs *Config) GetUniqSubKeys(prefix string) []string {
	kv, _ := cfgs.Settings()
	keys := []string{}
	for k, _ := range kv {
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
	return keys
}

func LoadConfig(configFile string) (ret *Config, err error) {
	tomlProvider := gconfig.NewTOMLFile(configFile)

	mappings := map[string]string{
		"KAFKA_STD_TOPIC":                  "kafka.standard_log_topic",
		"KAFKA_ACC_TOPIC":                  "kafka.access_log_topic",
		"KAFKA_HOSTS":                      "kafka.hosts",
		"KAFKA_FlUSH_FREQUENCY":            "kafka.flush_frequency",
		"KAFKA_MAX_CONNECTIONS":            "kafka.max_connections",
		"KAFKA_MAX_CONNECTIONS_PER_BROKER": "kafka.max_connections_per_broker",
		"KAFKA_BATCH_SIZE":                 "kafka.batch_size",
		"KAFKA_MAX_REQUESTS":               "kafka.max_requests",
		"KAFKA_SEND_ROUTINES":              "kafka.send_routines",
		"KAFKA_RECEIVE_ROUTINES":           "kafka.receive_routines",
		"KAFKA_REQUIRED_ACKS":              "kafka.required_acks",
		"KAFKA_QUEUE_LENGTH":               "kafka.queue_length",
		"KAFKA_ENQUEUE_TIMEOUT":            "kafka.enqueue_timeout",

		"LOG_LEVEL":         "log.level",
		"LOG_KAFKA_ENABLED": "log.kafka_enabled",
		"LOG_FILE_NAME":     "log.file_name",
		"LOG_FORMATTER":     "log.formatter",
		"API_PORT":          "api.port",
	}

	env := gconfig.NewEnvironment(mappings)
	cfgs := gconfig.NewConfig([]gconfig.Provider{tomlProvider, env})
	err = cfgs.Load()
	if err != nil {
		return
	}
	ret = &Config{cfgs}
	return
}

func (cfgs *Config) GetAPIPort() (ret int) {
	ret, _ = cfgs.IntOr("api.port", API_PORT)
	return
}

func (cfgs *Config) GetKafkaEnabled() (ret bool) {
	ret, _ = cfgs.BoolOr("log.kafka_enabled", KAFKA_ENABLED)
	return
}

func (cfgs *Config) LogConfigs() (fileName, level, formatter, stdLogVer, accessLogVer string) {
	fileName, _ = cfgs.StringOr("log.file_name", LOG_FILE_NAME)
	stdLogVer, _ = cfgs.StringOr("log.standard_log_version", LOG_VERSION)
	accessLogVer, _ = cfgs.StringOr("log.access_log_version", LOG_VERSION)
	level, _ = cfgs.StringOr("log.level", LOG_LEVEL)
	formatter, _ = cfgs.StringOr("log.formatter", LOG_FORMATTER)
	return
}

func (cfgs *Config) GetKafkaInfos() (ks dlogrus.KafkaSetting, accessLogTopic, standardLogTopic string, err error) {
	accessLogTopic, err = cfgs.String("kafka.access_log_topic")
	if err != nil {
		return
	}
	standardLogTopic, err = cfgs.String("kafka.standard_log_topic")
	if err != nil {
		return
	}
	queueLength, _ := cfgs.IntOr("kafka.queue_length", QUEUE_LENGTH)
	enqueueTimeout, _ := cfgs.IntOr("kafka.enqueue_timeout", ENQUEUE_TIMEOUT)

	var rawHost string
	rawHost, err = cfgs.String("kafka.hosts")

	fmt.Println("raw host : ", rawHost, " error : ", err)
	if err != nil {
		return
	}
	var freq, mc, mcpb, bs, mr, sr, rr, ra int
	freq, _ = cfgs.IntOr("kafka.flush_frequency", FLUSH_FREQUENCY)
	mc, _ = cfgs.IntOr("kafka.max_connections", MAX_CONNECTIONS)
	mcpb, _ = cfgs.IntOr("kafka.max_connections_per_broker", MAX_CONNECTIONS_PER_BROKER)
	bs, _ = cfgs.IntOr("kafka.batch_size", BATCH_SIZE)
	mr, _ = cfgs.IntOr("kafka.max_requests", MAX_REQUESTS)
	sr, _ = cfgs.IntOr("kafka.send_routines", SEND_ROUTINES)
	rr, _ = cfgs.IntOr("kafka.receive_routines", RECEIVE_ROUTINES)
	ra, _ = cfgs.IntOr("kafka.required_acks", REQUIRED_ACKS)

	if bs > MAX_BATCH_SIZE {
		err = errors.New("Batch size upper bound is 16384")
		return
	}
	if ra > 1 || ra < -1 {
		err = errors.New("Require Acks legal value -1, 0, 1")
		return
	}
	ks = dlogrus.KafkaSetting{
		Hosts:                   strings.Split(rawHost, ";"),
		Linger:                  time.Duration(freq) * time.Millisecond,
		MaxConnections:          mc,
		MaxConnectionsPerBroker: mcpb,
		BatchSize:               bs,
		MaxRequests:             mr,
		SendRoutines:            sr,
		ReceiveRoutines:         rr,
		RequiredAcks:            ra,
		LocalQueueLength:        queueLength,
		EnqueueTimeout:          time.Duration(enqueueTimeout) * time.Millisecond,
	}
	return
}

func (cfgs *Config) GetMgoDBInfo() *mongo.DBInfo {
	t := cfgs.GetMgoDBInfos()
	fmt.Printf("config %+v\n", t[0])
	return t[0]
}

func (cfgs *Config) GetMgoDBInfos() (ret []*mongo.DBInfo) {
	prefix := "database.mgo.instance."
	subkeys := cfgs.GetUniqSubKeys(prefix)

	fmt.Printf("%+v\n", subkeys)

	var key string
	var mgoDirect bool
	var mgoMaxConn, mgoMaxIdle, mgoDefaultTimeout int
	var mgoDefaultPort, mgoDefaultUser, mgoDefaultPassword, mgoDefaultDatabase, mgoDefaultAuthDatabase, mgoShardUser, mgoShardPassword string
	var name string

	name, _ = cfgs.StringOr("database.mgo.default.name", "")
	mgoMaxConn, _ = cfgs.IntOr("database.mgo.default.max_conn", MGO_MAX_CONN)
	mgoMaxIdle, _ = cfgs.IntOr("database.mgo.default.max_idle", MGO_MAX_IDLE)
	mgoDefaultUser, _ = cfgs.StringOr("database.mgo.default.user", MGO_USER)
	mgoDefaultPassword, _ = cfgs.StringOr("database.mgo.default.password", "")
	mgoDefaultDatabase, _ = cfgs.StringOr("database.mgo.default.database", MGO_DATABASE)
	mgoDefaultAuthDatabase, _ = cfgs.StringOr("database.mgo.default.authdatabase", MGO_AUTHDATABASE)
	mgoDefaultTimeout, _ = cfgs.IntOr("database.mgo.default.timeout", MGO_TIMEOUT)
	mgoShardUser, _ = cfgs.StringOr("database.mgo.default.sharduser", MGO_SHARDUSER)
	mgoShardPassword, _ = cfgs.StringOr("database.mgo.default.shardpassword", "")
	mgoDirect, _ = cfgs.Bool("database.mgo.default.direct")
	mgoAddrs := []string{}
	b := len(subkeys)
	for i := 0; i < b; i++ {
		key = prefix + subkeys[i]
		enabled, err := cfgs.Bool(fmt.Sprintf("%s.enabled", key))

		if err == nil && enabled {
			host, _ := cfgs.StringOr(fmt.Sprintf("%s.host", key), "")
			port, _ := cfgs.StringOr(fmt.Sprintf("%s.port", key), mgoDefaultPort)
			mgoAddrs = append(mgoAddrs, host+":"+port)
		}
	}
	ret = append(ret, mongo.NewDBInfo(name, mgoAddrs, mgoDefaultUser,
		mgoDefaultPassword, mgoDefaultDatabase, mgoDefaultAuthDatabase, mgoDefaultTimeout, mgoMaxConn, mgoMaxIdle, mgoShardUser, mgoShardPassword, mgoDirect))
	return
}
