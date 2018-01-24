package config

import (
	"errors"
	"fmt"
	"github.com/DroiTaipei/dlogrus"
	"github.com/DroiTaipei/mongo"

	// gconfig "github.com/zpatrick/go-config"
	"github.com/DroiTaipei/droipkg/config"
	"strings"
	"time"
)

const (
	EMPTY             = ""
	API_PORT          = 8080
	FORWARDER_PORT    = 8099
	FORWARDER_TIMEOUT = 190
	API_POD_NAME      = "db-api"
	// API_MODE          = "release"
	// API_VERSION       = "v1"
	// API_TIMEOUT       = "30"

	MGO_PORT         = "7379"
	MGO_MAX_CONN     = 250
	MGO_MAX_IDLE     = 240
	MGO_USER         = "bass"
	MGO_DATABASE     = "baas"
	MGO_AUTHDATABASE = "admin"
	MGO_TIMEOUT      = 300
	MGO_SHARDUSER    = "sharder"
	MGO_SECONDARY    = true

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

	// Jaeger default config
	JAEGER_HOST           = "tpe-jaeger-agent"
	JAEGER_PORT           = 6831
	JAEGER_SAMPLE_RATE    = 0.5
	JAEGER_QUEUE_SIZE     = 65536
	JAEGER_FLUSH_INTERVAL = 60 * time.Second
)

type Config struct {
	*config.Config
}

type TracerOption struct {
	Ip                  string        `toml:"Ip"`
	Port                int           `toml:"Port"`
	SampleRate          float64       `toml:"SampleRate"`
	QueueSize           int           `toml:"QueueSize"`
	BufferFlushInterval time.Duration `toml:"BufferFlushInterval"`
}

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
	return keys
}

func LoadConfig(configFile string) (ret *Config, err error) {
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
		"API_PORT":          "api.api_port",
		"FORWARDER_PORT":    "api.forwarder_port",
		"FORWARDER_TIMEOUT": "api.timeout",
	}

	opts := &config.Options{
		EnvMap: mappings,
	}
	cfgs, err := config.LoadConfig(configFile, opts)
	if err != nil {
		panic(err)
	}
	ret = &Config{cfgs}

	return
}

func (cfgs *Config) GetAPIPort() (api_port, forwarder_port int) {
	api_port, _ = cfgs.IntOr("api.api_port", API_PORT)
	forwarder_port, _ = cfgs.IntOr("api.forwarder_port", FORWARDER_PORT)
	return
}

func (cfgs *Config) GetTimeout() (timeout int) {
	timeout, _ = cfgs.IntOr("api.timeout", FORWARDER_TIMEOUT)
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
		MetadataRetries:         10,
		MetadataBackoff:         time.Duration(10) * time.Second,
		MetadataTTL:             time.Duration(300) * time.Second,
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
	var mgoDirect, mgoSecondary bool
	var mgoMaxConn, mgoDefaultTimeout int
	var mgoDefaultPort, mgoDefaultUser, mgoDefaultPassword, mgoDefaultAuthDatabase string
	var name string

	name, _ = cfgs.StringOr("database.mgo.default.name", "")
	mgoMaxConn, _ = cfgs.IntOr("database.mgo.default.max_conn", MGO_MAX_CONN)
	// mgoMaxIdle, _ = cfgs.IntOr("database.mgo.default.max_idle", MGO_MAX_IDLE)
	mgoDefaultUser, _ = cfgs.StringOr("database.mgo.default.user", MGO_USER)
	mgoDefaultPassword, _ = cfgs.StringOr("database.mgo.default.password", "")
	// mgoDefaultDatabase, _ = cfgs.StringOr("database.mgo.default.database", MGO_DATABASE)
	mgoDefaultAuthDatabase, _ = cfgs.StringOr("database.mgo.default.authdatabase", MGO_AUTHDATABASE)
	mgoDefaultTimeout, _ = cfgs.IntOr("database.mgo.default.timeout", MGO_TIMEOUT)
	// mgoShardUser, _ = cfgs.StringOr("database.mgo.default.sharduser", MGO_SHARDUSER)
	// mgoShardPassword, _ = cfgs.StringOr("database.mgo.default.shardpassword", "")
	mgoDirect, _ = cfgs.Bool("database.mgo.default.direct")
	mgoSecondary, _ = cfgs.BoolOr("database.mgo.default.secondary", MGO_SECONDARY)
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
	// ret = append(ret, mongo.NewDBInfo(name, mgoAddrs, mgoDefaultUser,
	// mgoDefaultPassword, mgoDefaultDatabase, mgoDefaultAuthDatabase, mgoDefaultTimeout, mgoMaxConn, mgoMaxIdle, mgoShardUser, mgoShardPassword, mgoDirect))
	ret = append(ret, mongo.NewDBInfo(name, mgoAddrs, mgoDefaultUser,
		mgoDefaultPassword, mgoDefaultAuthDatabase, mgoDefaultTimeout, mgoMaxConn, mgoDirect, mgoSecondary))
	return
}

func (cfgs *Config) GetMgoDBName() (mgoDatabase string) {
	mgoDatabase, _ = cfgs.StringOr("database.mgo.default.database", MGO_DATABASE)
	return
}

func (cfgs *Config) GetJaegerConfig() *TracerOption {
	host, _ := cfgs.StringOr("jaeger.host", JAEGER_HOST)
	port, _ := cfgs.IntOr("jaeger.port", JAEGER_PORT)
	sampleRate, _ := cfgs.FloatOr("jaeger.sample_rate", JAEGER_SAMPLE_RATE)
	qSize, _ := cfgs.IntOr("jaeger.queue_size", JAEGER_QUEUE_SIZE)
	buffer := cfgs.Duration("jaeger.flush_interval", JAEGER_FLUSH_INTERVAL)
	return &TracerOption{Ip: host, Port: port, SampleRate: sampleRate, QueueSize: qSize, BufferFlushInterval: buffer}
}

// func (cfgs *Config) GetApiConfigs() (ret map[string]string) {
// 	ret["api_port"], _ = cfgs.StringOr("api.api_port", API_PORT)
// 	ret["forwarder_port"], _ = cfgs.StringOr("api.forwarder_port", FORWARDER_PORT)
// 	ret["mode"], _ = cfgs.StringOr("api.mode", API_MODE)
// 	ret["version"], _ = cfgs.StringOr("api.version", API_VERSION)
// 	ret["readTimeout"], _ = cfgs.StringOr("api.readTimeout", API_TIMEOUT)
// 	ret["writeTimeout"], _ = cfgs.StringOr("api.writeTimeout", API_TIMEOUT)
// 	return ret
// }
