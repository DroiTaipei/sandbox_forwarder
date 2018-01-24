package dlogrus

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/DroiTaipei/droipkg"
	"github.com/DroiTaipei/logrus"
)

const (
	UNKNOWN_VALUE                = "U"
	STANDARD_LOG_VERSION_FIELD   = "V"
	ACCESS_LOG_VERSION_FIELD     = "A"
	TIME_FIELD                   = "T"
	MESSAGE_FIELD                = "M"
	LEVEL_FIELD                  = "L"
	POD_NAME_FIELD               = "Pd"
	NODE_NAME_FIELD              = "Nd"
	DISCARD_FILE_NAME            = "Discard"
	OS_STDOUT                    = "os.Stdout"
	TEXT_FORMATTER_CODE          = "text"
	JSON_FORMATTER_CODE          = "json"
	BSON_FORMATTER_CODE          = "bson"
	STANDARD_LOG_VERSION_DEFAULT = "1"
	ACCESS_LOG_VERSION_DEFAULT   = "1"
)

type Fields logrus.Fields

var settings map[string]string
var logFd *os.File
var logger *logrus.Logger
var hook *kafkaHook

func init() {
	logrus.SetFormatter(newJSONFormatter(""))
	// Output to stderr instead of stdout, could also be a file.
	logrus.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
	logger = logrus.StandardLogger()
}

func setOutput(fileName string) (err error) {

	switch fileName {
	case DISCARD_FILE_NAME:
		logger.Out = ioutil.Discard
	case OS_STDOUT:
		logger.Out = os.Stdout
	default:
		logFd, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		logger.Out = logFd
	}
	return
}

func setLevel(lvl string) (err error) {
	level, err := logrus.ParseLevel(lvl)
	if err == nil {
		logger.Level = level
	}
	return
}

func setFormatter(formatterName string) (err error) {
	switch formatterName {
	case JSON_FORMATTER_CODE:
		logger.Formatter = newJSONFormatter("")
	case BSON_FORMATTER_CODE:
		logger.Formatter = newBSONFormatter("")
	case TEXT_FORMATTER_CODE:
		logger.Formatter = new(logrus.TextFormatter)
	default:
		err = errors.New("There is no valid formatter")
	}
	return
}

func Initialize(fileName, level, formatter, stdLogVer, accessLogVer string) (err error) {
	if len(fileName) > 0 {
		err = setOutput(fileName)
		if err != nil {
			return
		}
	}
	if len(level) > 0 {
		err = setLevel(level)
		if err != nil {
			return
		}
	}
	if len(formatter) > 0 {
		err = setFormatter(formatter)
		if err != nil {
			return
		}
	}

	settings = make(map[string]string)
	settings[STANDARD_LOG_VERSION_FIELD] = STANDARD_LOG_VERSION_DEFAULT
	if len(stdLogVer) > 0 {
		settings[STANDARD_LOG_VERSION_FIELD] = stdLogVer
	}

	settings[ACCESS_LOG_VERSION_FIELD] = ACCESS_LOG_VERSION_DEFAULT

	if len(accessLogVer) > 0 {
		settings[ACCESS_LOG_VERSION_FIELD] = accessLogVer
	}

	pd, err := GetPodname()
	if err != nil {
		pd = UNKNOWN_VALUE
	}
	settings[POD_NAME_FIELD] = pd

	nd, err := GetIP()
	if err != nil {
		nd = UNKNOWN_VALUE
	}
	settings[NODE_NAME_FIELD] = nd
	return
}

func StandardLogger() *droipkg.DroiLogger {
	return &droipkg.DroiLogger{logger}
}

// Sugar Setting
func SetDevelopMode() (err error) {
	err = setFormatter(TEXT_FORMATTER_CODE)
	if err != nil {
		return
	}
	// Output to stderr instead of stdout, could also be a file.
	err = setOutput(OS_STDOUT)
	if err != nil {
		return
	}
	err = setLevel("debug")
	logger.Hooks = make(logrus.LevelHooks)
	return
}

func Close() {
	if logFd != nil {
		logFd.Close()
	}
	if hook != nil {
		hook.Close()
	}
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func WithFieldsDebug(fields Fields, msg string) {
	logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func WithFieldsInfo(fields Fields, msg string) {
	logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func WithFieldsWarn(fields Fields, msg string) {
	logger.WithFields(logrus.Fields(fields)).Warn(msg)
}

func WithFieldsError(fields Fields, msg string) {
	logger.WithFields(logrus.Fields(fields)).Error(msg)
}

func WithFieldsFatal(fields Fields, msg string) {
	logger.WithFields(logrus.Fields(fields)).Fatal(msg)
}
