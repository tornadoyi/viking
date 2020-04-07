

package log

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"sync"
)


var (
	mutex sync.RWMutex
	loggers = make(map[string]*Logger)
	defaultLogger = (*Logger)(nil)
	_ = createDefaultLogger()
)


func GetLogger(name string) (*Logger, bool){
	mutex.RLock()
	defer mutex.RUnlock()
	l, ok := loggers[name]
	return l, ok
}


func createDefaultLogger() *Logger{
	const defaultLoggerConfig = `
level: debug
disableCaller: true
disableStacktrace: true
encoding: console
encoderConfig:
  messageKey: msg
  levelKey: lvl
  timeKey: ts
  lineEnding: "\n"
  levelEncoder: capitalColor
  timeEncoder: iso8601
stdout: true
strerr: true
`
	var cfg *Config
	if err := yaml.Unmarshal([]byte(defaultLoggerConfig), &cfg); err != nil { panic(err) }
	l, err := NewLoggerWithConfig(cfg)
	if err != nil { panic(err) }
	SetDefaultLogger(l)
	return l
}

func CreateLogger(name string, cfg *Config, opts... Option) (*Logger, error){
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := loggers[name]; ok { return nil, fmt.Errorf("Repeated logger %v", name)}
	l, err := NewLoggerWithConfig(cfg)
	if err != nil { return nil, err }
	loggers[name] = l
	return l, nil
}


func SetDefaultLogger(logger *Logger) {
	mutex.Lock()
	defer mutex.Unlock()

	// set defaylt logger
	defaultLogger = logger

	// export function
	DPanic = defaultLogger.DPanic
	DPanicf = defaultLogger.DPanicf
	DPanicw = defaultLogger.DPanicw
	Debug = defaultLogger.Debug
	Debugf = defaultLogger.Debugf
	Debugw = defaultLogger.Debugw
	Error = defaultLogger.Error
	Errorf = defaultLogger.Errorf
	Errorw = defaultLogger.Errorw
	Fatal = defaultLogger.Fatal
	Fatalf = defaultLogger.Fatalf
	Fatalw = defaultLogger.Fatalw
	Info = defaultLogger.Info
	Infof = defaultLogger.Infof
	Infow = defaultLogger.Infow
	Panic = defaultLogger.Panic
	Panicf = defaultLogger.Panicf
	Panicw = defaultLogger.Panicw
	Warn = defaultLogger.Warn
	Warnf = defaultLogger.Warnf
	Warnw = defaultLogger.Warnw
}


var DPanic (func(args ...interface{}))
var DPanicf func(template string, args ...interface{})
var DPanicw func(msg string, keysAndValues ...interface{})
var Debug func(args ...interface{})
var Debugf func(template string, args ...interface{})
var Debugw func(msg string, keysAndValues ...interface{})
var Error func(args ...interface{})
var Errorf func(template string, args ...interface{})
var Errorw func(msg string, keysAndValues ...interface{})
var Fatal func(args ...interface{})
var Fatalf func(template string, args ...interface{})
var Fatalw func(msg string, keysAndValues ...interface{})
var Info func(args ...interface{})
var Infof func(template string, args ...interface{})
var Infow func(msg string, keysAndValues ...interface{})
var Panic func(args ...interface{})
var Panicf func(template string, args ...interface{})
var Panicw func(msg string, keysAndValues ...interface{})
var Warn func(args ...interface{})
var Warnf func(template string, args ...interface{})
var Warnw func(msg string, keysAndValues ...interface{})







