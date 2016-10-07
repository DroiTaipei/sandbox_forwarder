package dlogrus

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/devopstaku/logrus"
	"github.com/elodina/siesta"
	sp "github.com/elodina/siesta-producer"
)

const (
	MAX_RETRY         = 20
	STANDARD_RETRY    = 3
	KAFKA_NO_RESPONSE = 0
)

var defaultLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.DebugLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
}

type KafkaSetting struct {
	Hosts                   []string
	Linger                  time.Duration
	MaxConnections          int
	MaxConnectionsPerBroker int
	BatchSize               int
	MaxRequests             int
	SendRoutines            int
	ReceiveRoutines         int
	RequiredAcks            int
	LocalQueueLength        int
	EnqueueTimeout          time.Duration
}

type kafkaHook struct {
	Closed         bool
	RequireAck     bool
	accLogTopic    string
	stdLogTopic    string
	enqueueTimeout time.Duration
	levels         []logrus.Level
	producer       *sp.KafkaProducer
	sendQueue      chan *sp.ProducerRecord
	metadatas      chan (<-chan *sp.RecordMetadata)
}

func NewKafkaSetting() KafkaSetting {
	return KafkaSetting{
		Linger:                  100 * time.Millisecond,
		MaxConnections:          5,
		MaxConnectionsPerBroker: 5,
		BatchSize:               16384,
		MaxRequests:             10,
		SendRoutines:            10,
		ReceiveRoutines:         10,
		RequiredAcks:            1,
		LocalQueueLength:        1024,
		EnqueueTimeout:          1000 * time.Millisecond,
	}
}

func ConnectKafka(ks KafkaSetting, accessLogTopic, standardLogTopic string) error {
	hook, err := NewHook(ks, accessLogTopic, standardLogTopic)
	if err == nil {
		logrus.AddHook(hook)
	}
	return err
}

func kafkaConnect(ks KafkaSetting) (*sp.KafkaProducer, error) {
	config := siesta.NewConnectorConfig()
	config.BrokerList = ks.Hosts
	config.MaxConnections = ks.MaxConnections
	config.MaxConnectionsPerBroker = ks.MaxConnectionsPerBroker
	config.MetadataRetries = STANDARD_RETRY
	connector, err := siesta.NewDefaultConnector(config)
	if err != nil {
		return nil, err
	}
	producerConfig := sp.NewProducerConfig()
	producerConfig.BatchSize = ks.BatchSize
	producerConfig.ClientID = "XD"
	producerConfig.MaxRequests = ks.MaxRequests
	producerConfig.SendRoutines = ks.SendRoutines
	producerConfig.ReceiveRoutines = ks.ReceiveRoutines
	// 0 = NoResponse, 1 = WaitForLocal, -1 = WaitForAll, default is 1
	producerConfig.RequiredAcks = ks.RequiredAcks
	producerConfig.AckTimeoutMs = 30000
	// Flush Frequency
	producerConfig.Linger = ks.Linger
	producerConfig.RetryBackoff = 100 * time.Millisecond
	producerConfig.Retries = STANDARD_RETRY
	p := sp.NewKafkaProducer(producerConfig, sp.ByteSerializer, sp.ByteSerializer, connector)
	return p, nil
}

func NewHook(ks KafkaSetting, accessLogTopic, standardLogTopic string) (*kafkaHook, error) {
	h := kafkaHook{
		RequireAck:     ks.RequiredAcks != KAFKA_NO_RESPONSE,
		enqueueTimeout: ks.EnqueueTimeout,
		accLogTopic:    accessLogTopic,
		stdLogTopic:    standardLogTopic,
		levels:         defaultLevels,
	}

	producer, err := kafkaConnect(ks)
	if err != nil {
		return nil, err
	}
	h.producer = producer
	h.sendQueue = make(chan *sp.ProducerRecord, ks.LocalQueueLength)
	go h.sendHandle()
	if h.RequireAck {
		h.metadatas = make(chan (<-chan *sp.RecordMetadata), ks.LocalQueueLength)
		go h.responseHandle()
	}
	return &h, nil
}

func (hook *kafkaHook) sendHandle() {
	var msg *sp.ProducerRecord
	for msg = range hook.sendQueue {
		if hook.RequireAck {
			hook.metadatas <- hook.producer.Send(msg)
		} else {
			hook.producer.Send(msg)
		}
	}
}

func (hook *kafkaHook) responseHandle() {
	for metadataChan := range hook.metadatas {
		select {
		case metadata := <-metadataChan:
			if metadata.Error != siesta.ErrNoError {
				println(metadata.Error.Error())
			}
		default:
			if hook.Closed {
				return
			}
		}
	}
}

func (hook *kafkaHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	topic := hook.stdLogTopic
	if _, ok := entry.Data[ACCESS_LOG_VERSION_FIELD]; ok {
		topic = hook.accLogTopic
	}

	select {
	case hook.sendQueue <- &sp.ProducerRecord{Topic: topic, Value: line}:
		return nil
	case <-time.After(hook.enqueueTimeout):
		return errors.New("Enqueue Timeout Drop")
	}
}

func (hook *kafkaHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *kafkaHook) SetLevels(levels []logrus.Level) {
	hook.levels = levels
}

func (hook *kafkaHook) Close() {
	if hook.producer != nil {
		hook.producer.Close()
	}
	if hook.sendQueue != nil {
		close(hook.sendQueue)
	}
	if hook.metadatas != nil {
		close(hook.metadatas)
	}
	hook.Closed = true
}
