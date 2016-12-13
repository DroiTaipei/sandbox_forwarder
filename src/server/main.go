package main

import (
	"api"
	"flag"
	"fmt"
	"github.com/DroiTaipei/dlogrus"
	"github.com/DroiTaipei/droipkg"
	"github.com/DroiTaipei/mongo"
	"github.com/valyala/fasthttp"
	// "github.com/valyala/fasthttp/reuseport"
	stdlog "log"
	"os"
	"runtime"
	"util/config"
)

type options struct {
	version bool
	prof    bool
}

var opts options

func init() {
	flag.BoolVar(&opts.version, "build", false, "GoLang build version.")
	flag.BoolVar(&opts.prof, "prof", false, "GoLang profiling function.")
}

func main() {
	var cfgFile string

	flag.StringVar(&cfgFile, "config", "./conf.d/current.toml", "Path to Config File")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] <command> \n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if opts.version {
		fmt.Fprintf(os.Stderr, "%s\n", runtime.Version())
	}
	if opts.prof {
		ActivateProfile()
	}
	stdlog.Fatal(run(cfgFile))
}

func run(cfgFilePath string) (err error) {
	cfg, err := config.LoadConfig(cfgFilePath)
	if err != nil {
		return droipkg.Wrap(err, "Load Config Failed")
	}

	err = dlogrus.Initialize(cfg.LogConfigs())
	if err != nil {
		return droipkg.Wrap(err, "Get Log Config Failed")
	}

	if cfg.GetKafkaEnabled() {
		stdlog.Println("trying to connect Kafka")

		ks, alt, slt, err := cfg.GetKafkaInfos()
		if err != nil {
			return droipkg.Wrap(err, "Get Kafka config failed")
		}
		err = dlogrus.ConnectKafka(ks, alt, slt)
		if err != nil {
			return droipkg.Wrap(err, "Connected Kafka failed")
		}
		stdlog.Println("Kafka Connected")
	}

	err = mongo.Initialize(cfg.GetMgoDBInfo(), "_Id", dlogrus.StandardLogger())
	if err != nil {
		return droipkg.Wrap(err, "Mongo initailize failed")
	}
	fmt.Printf("mongo config %+v\n", cfg.GetMgoDBInfo())
	defer mongo.Close()

	droipkg.SetLogger(dlogrus.StandardLogger())

	api_port, forwarder_port := cfg.GetAPIPort()
	timeout := cfg.GetTimeout()

	apiRouter := api.ApiRegist(timeout)
	forwarderRouter := api.ForwarderRegist(timeout)

	go func() {
		bind_api := fmt.Sprintf(":%d", api_port)

		stdlog.Println("API server start at port ", api_port)

		err := fasthttp.ListenAndServe(bind_api, apiRouter.Handler)
		if err != nil {
			stdlog.Fatalf("API server crash with error %v", err.Error())
		}
	}()

	bind_forwarder := fmt.Sprintf(":%d", forwarder_port)
	fasthttp.ListenAndServe(bind_forwarder, forwarderRouter.Handler)

	return nil
}
