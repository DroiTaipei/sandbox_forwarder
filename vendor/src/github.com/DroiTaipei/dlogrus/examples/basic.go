package main

import (
	log "github.com/DroiTaipei/dlogrus"
	"time"
)

func init() {
	logFile := "debug.log"
	logLevel := "debug"
	logVer := "1"
	logTimeFormat := time.RFC3339
	log.Initialize(logFile, logLevel, logVer, logTimeFormat)
	// log.ConnectKafka([]string{"127.0.0.1:9092"}, "test")
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.WithFieldsFatal(log.Fields{
				"omg":    true,
				"err":    err,
				"number": 100,
			}, "The ice breaks!")
		}
	}()

	log.WithFieldsDebug(log.Fields{
		"animal": "walrus",
		"number": 8,
	}, "Started observing beach")

	log.WithFieldsInfo(log.Fields{
		"animal": "walrus",
		"size":   10,
	}, "A group of walrus emerges from the ocean")
	log.WithFieldsError(log.Fields{
		"temperature": -4,
	}, "Temperature changes")

	log.WithFieldsWarn(log.Fields{
		"Alert": "For the world",
	}, "We need warning")
	panic("All is accident")

}
