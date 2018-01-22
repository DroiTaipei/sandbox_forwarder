package api

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"gopkg.in/ory-am/dockertest.v3"
	"os"
	"testing"
	"util/config"
)

var (
	pool     *dockertest.Pool
	resource *dockertest.Resource

	MockMongoHost = "localhost:7379"

	dockerOpt = dockertest.RunOptions{
		Repository:   "mongo",
		Tag:          "latest",
		Env:          []string{},
		ExposedPorts: []string{"7379/tcp"},
	}
)

func BeforeTest() {
	confFilePath := "../../../conf.d/dev.toml"
	fmt.Println("[UnitTest] Read config file from " + confFilePath)

	// Populate config
	config.LoadConfig(confFilePath)

	// err = mongo.Initialize(cfg.GetMgoDBInfo(), "_Id", dlogrus.StandardLogger())
	// if err != nil {
	// 	return droipkg.Wrap(err, "Mongo initailize failed")
	// }
	// fmt.Printf("mongo config %+v\n", cfg.GetMgoDBInfo())
	// defer mongo.Close()
}

func DockerTestInit() (err error) {
	pool, err = dockertest.NewPool("")
	if err != nil {
		return err
	}

	// pulls an image, creates a container based on it and runs it
	resource, err = pool.RunWithOptions(&dockerOpt)
	if err != nil {
		return err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		broker := sarama.NewBroker(MockKafkaHost)
		kafkaErr := broker.Open(nil)
		if kafkaErr != nil {
			return kafkaErr
		}
		producer, err := sarama.NewSyncProducer([]string{MockKafkaHost}, nil)
		if err != nil {
			return err
		}
		defer producer.Close()
		defer broker.Close()
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func TestGetFullURI(t *testing.T) {
	URI := getFullURI("123456", []byte("where={}&limit=2"))
	fmt.Println(URI)
	if URI != "123456?where={}&limit=2" {
		t.Fail()
	}
}

func AfterTest() {

}

func TestMain(m *testing.M) {
	BeforeTest()
	retCode := m.Run()
	AfterTest()
	os.Exit(retCode)
}
