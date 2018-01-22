package droipkg

import (
	"github.com/DroiTaipei/logrus"
	"io/ioutil"
	"os"
)

type Logger interface {
	WithError(err error) *logrus.Entry
	WithField(key string, value interface{}) *logrus.Entry
	WithMap(input map[string]interface{}) *logrus.Entry
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Println(args ...interface{})
}

type DroiLogger struct {
	*logrus.Logger
}

func (dl *DroiLogger) WithMap(input map[string]interface{}) *logrus.Entry {
	return dl.WithFields(logrus.Fields(input))
}

func init() {
	// Dummmy for Stupid Usage
	l := logrus.New()
	l.Out = os.Stdout
	l.Level = logrus.DebugLevel
	SetLogger(&DroiLogger{l})
}

var stdLogger Logger

func SetLogger(l Logger) {
	stdLogger = l
}

func GetDiscardLogger() Logger {
	// Dummmy for Stupid Usage
	l := logrus.New()
	l.Out = ioutil.Discard
	l.Level = logrus.PanicLevel
	return &DroiLogger{l}
}

func GetLogger() Logger {
	return stdLogger
}
