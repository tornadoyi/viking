package log

import (
	l4g "github.com/jeanphorn/log4go"
	"sync"
)

var (
	loggers = make(map[string]Logger)
	defaultLogger Logger
	mutex sync.Mutex
)

func GetLogger(name string) Logger{
	defer mutex.Unlock()
	mutex.Lock()
	log, ok := loggers[name]
	if !ok {
		log = make(Logger)
		loggers[name] = log
	}
	return log
}

func SetDefault(log Logger) {
	defaultLogger = log
}



// export
type Logger = l4g.Logger

func Finest(arg0 interface{}, args ...interface{}) { defaultLogger.Finest(arg0, args...) }

func Fine(arg0 interface{}, args ...interface{}) { defaultLogger.Fine(arg0, args...) }

func Debug(arg0 interface{}, args ...interface{}) { defaultLogger.Debug(arg0, args...) }

func Trace(arg0 interface{}, args ...interface{}) { defaultLogger.Trace(arg0, args...) }

func Info(arg0 interface{}, args ...interface{}) { defaultLogger.Info(arg0, args...) }

func Warn(arg0 interface{}, args ...interface{}) { defaultLogger.Warn(arg0, args...) }

func Error(arg0 interface{}, args ...interface{}) { defaultLogger.Error(arg0, args...) }

func Critical(arg0 interface{}, args ...interface{}) { defaultLogger.Critical(arg0, args...) }