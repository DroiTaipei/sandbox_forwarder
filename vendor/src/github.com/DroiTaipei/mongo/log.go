package mongo

import (
	"time"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
)

const (
	VERSION                    = "1"
	TOPIC_KEY                  = "TOPIC"
	UNKNOWN_VALUE              = "U"
	STANDARD_LOG_VERSION_FIELD = "V"
	ACCESS_LOG_VERSION_FIELD   = "A"
	FUNCTION_FIELD             = "fc"
	FUNCTION_ARGS_FIELD        = "fa"
	POD_NAME_FIELD             = "Pd"
	NODE_NAME_FIELD            = "Nd"
	METHOD_FIELD               = "Md"
	URI_FIELD                  = "Uri"
	DB_TYPE_FIELD              = "Dt"
	DB_HOSTNAME_FIELD          = "Dh"
	DB_COMMAND_FIELD           = "Dc"
	DB_COMMAND_TIME_FIELD      = "Dct"
	DB_COMMAND_TYPE_FIELD      = "Dcc"
	REQUEST_TIME_FIELD         = "Rt"
	DISCARD_FILE_NAME          = "Discard"
	DB_TYPE                    = "mongo"
)

// systemLogHeaders used for error log format
var systemCtx = &droictx.DoneContext{}
var mongoAccTopic = ""
var mongoStdTopic = ""

func init() {
	systemCtx.Set("Aid", "ay8umbzhD9bxb_hwRC7z-RMyw2vFYUzXlQDNDIwA")
	systemCtx.Set("SAid", "ay8umbzhD9bxb_hwRC7z-RMyw2vFYUzXlQDNDIwA")
	systemCtx.Set("URL", "/MongoDao")
	systemCtx.Set("Aidm", "prod")
	systemCtx.Set("SAidm", "prod")
	systemCtx.Set("Rid", "M1111111111111T1111111111")
}

func spentTime(t time.Time) int64 {
	d := time.Since(t)
	// 其實正解應該是 int64(math.Ceil(d.Seconds() * 1e3))
	// 不過不想浪費效能來算....直接無條件進位了！
	// 所以千萬不能用這個來作效能評估，只能拿來計價（喂
	return (d.Nanoseconds() / 1e6) + 1
}

func accessLog(ctx droictx.Context, method, sql string, start time.Time) {
	logger := droipkg.GetLogger().WithMap(ctx.Map()).
		WithField(ACCESS_LOG_VERSION_FIELD, VERSION).
		WithField(DB_COMMAND_FIELD, sql).
		WithField(DB_COMMAND_TYPE_FIELD, method).
		WithField(DB_COMMAND_TIME_FIELD, spentTime(start)).
		WithField(DB_TYPE_FIELD, DB_TYPE)
	if len(mongoAccTopic) > 0 {
		logger = logger.WithField(TOPIC_KEY, mongoAccTopic)
	}

	logger.Info(ACCESS_LOG_VERSION_FIELD)
}

func errLog(ctx droictx.Context, msg string) {
	logger := droipkg.GetLogger().WithMap(ctx.Map()).
		WithField(STANDARD_LOG_VERSION_FIELD, VERSION).
		WithField(DB_TYPE_FIELD, DB_TYPE)
	if len(mongoStdTopic) > 0 {
		logger = logger.WithField(TOPIC_KEY, mongoStdTopic)
	}

	logger.Error(msg)
}

func infoLog(ctx droictx.Context, msg string) {
	logger := droipkg.GetLogger().WithMap(ctx.Map()).
		WithField(STANDARD_LOG_VERSION_FIELD, VERSION).
		WithField(DB_TYPE_FIELD, DB_TYPE)
	if len(mongoStdTopic) > 0 {
		logger = logger.WithField(TOPIC_KEY, mongoStdTopic)
	}
	logger.Info(msg)
}
