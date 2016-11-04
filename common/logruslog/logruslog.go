package logruslog

import (
	"github.com/gaobrian/open-falcon-backend/common/vipercfg"
	log "github.com/Sirupsen/logrus"
	"runtime"
)

func logLevel(l string) log.Level {
	switch l {
	case "debug", "Debug", "DEBUG":
		return log.DebugLevel
	case "info", "Info", "INFO":
		return log.InfoLevel
	case "warn", "Warn", "WARN":
		return log.WarnLevel
	case "error", "Error", "ERROR":
		return log.ErrorLevel
	case "fatal", "Fatal", "FATAL":
		return log.FatalLevel
	case "panic", "Panic", "PANIC":
		return log.PanicLevel
	default:
		return log.InfoLevel

	}
}

func Init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	logLevelStr := vipercfg.Config().GetString("logLevel")
	logLevelStr = "debug"
	log.SetLevel(logLevel(logLevelStr))

	if log.GetLevel() == log.DebugLevel{
		hook := new(CallerHook)
		log.AddHook(hook)
	}
}


type CallerHook struct {
}

func (hook *CallerHook) Fire(entry *log.Entry) error {
	skipFrames := 6
	if len(entry.Data) > 0 {
		skipFrames = 6
	}
	_, fn, line, _ := runtime.Caller(skipFrames)
	entry.Data["file"] = fn
	entry.Data["line"] = line
	return nil
}

func (hook *CallerHook) Levels() []log.Level {
	return log.AllLevels
}


