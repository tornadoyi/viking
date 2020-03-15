package log

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	l4g "viking/log/log4go"
)

var (
	loggers = make(map[string]Logger)
	defaultLogger = createDefaultLogger()
	mutex sync.Mutex
)


func GetLogger(name string) Logger{
	defer mutex.Unlock()
	mutex.Lock()
	log, ok := loggers[name]
	if !ok { return nil}
	return log
}

func SetDefaultLogger(log Logger) {
	defaultLogger = log
}


func createDefaultLogger() Logger{
	log := CreateLogger("__default__")
	lw := l4g.NewConsoleLogWriter()
	lw.SetFormat("[%D %T] [%L] %M")
	log.AddFilter("stdout", l4g.FINEST, lw)
	return log
}

func CreateLogger(name string) Logger{
	defer mutex.Unlock()
	mutex.Lock()

	log := make(Logger)
	if _, ok := loggers[name]; ok { panic(errors.New(fmt.Sprintf("Repeated logger name %v", name))) }
	loggers[name] = log

	// destructor
	runtime.SetFinalizer(&log, func (log *Logger){
		log.Close()
	})
	return log
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

func ParseYamlConfig(config string) (*l4g.YamlLogConfig, error) { return l4g.ParseYamlConfig(config) }