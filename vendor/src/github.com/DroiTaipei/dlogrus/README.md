# dlogrus

The wrapper of logrus for Droi


## Function & Usage

```
err := log.Initialize(fileName, level, formatter, stdLogVer, accessLogVer)
if err != nil {
    panic(err)
}
defer log.Close()

kafkaEnabled := true
if err == nil && kafkaEnabled {
    ks := log.NewKafkaSetting()
    ks.Hosts = []string{"127.0.0.1:9092"}
    accessLogTopic := "testDebugLevel"
    standardLogTopic := "testDebugLevel"
    err = log.ConnectKafka(ks, accessLogTopic, standardLogTopic)
    if err != nil {
        panic(err)
    }
    println("Kafka Connected")
}
// 使用 DevelopMode （開發模式
// Output 會在 os.Stdout
// Level 會設為 Debug
// Formatter 會用 Text
// Kafka 也會 Disabled
// 方便開發者 Debug
log.SetDevelopMode()

```

### Example: Basic

```
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
    log.Initialize(logFile, logLevel, logTimeFormat, logVer, logVer)
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

```

