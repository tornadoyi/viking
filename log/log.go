

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
	l, err := NewLoggerWithConfig(cfg, opts...)
	if err != nil { return nil, err }
	loggers[name] = l
	return l, nil
}


func SetDefaultLogger(logger *Logger) {
	mutex.Lock()
	defer mutex.Unlock()
	defaultLogger = logger
}


func DPanic(args ...interface{}) { defaultLogger.DPanic(args...) }
func DPanicf(template string, args ...interface{}) { defaultLogger.DPanicf(template, args...) }
func DPanicw(msg string, keysAndValues ...interface{}) { defaultLogger.DPanicw(msg, keysAndValues...) }
func Debug(args ...interface{}) { defaultLogger.Debug(args...) }
func Debugf(template string, args ...interface{}) { defaultLogger.Debugf(template, args...) }
func Debugw(msg string, keysAndValues ...interface{}) { defaultLogger.Debugw(msg, keysAndValues...) }
func Error(args ...interface{}) { defaultLogger.Error(args...) }
func Errorf(template string, args ...interface{}) { defaultLogger.Errorf(template, args...) }
func Errorw(msg string, keysAndValues ...interface{}) { defaultLogger.Errorw(msg, keysAndValues...) }
func Fatal(args ...interface{}) { defaultLogger.Fatal(args...) }
func Fatalf(template string, args ...interface{}) { defaultLogger.Fatalf(template, args...) }
func Fatalw(msg string, keysAndValues ...interface{}) { defaultLogger.Fatalw(msg, keysAndValues...) }
func Info(args ...interface{}) { defaultLogger.Info(args...) }
func Infof(template string, args ...interface{}) { defaultLogger.Infof(template, args...) }
func Infow(msg string, keysAndValues ...interface{}) { defaultLogger.Infow(msg, keysAndValues...) }
func Panic(args ...interface{}) { defaultLogger.Panic(args...) }
func Panicf(template string, args ...interface{}) { defaultLogger.Panicf(template, args...) }
func Panicw(msg string, keysAndValues ...interface{}) { defaultLogger.Panicw(msg, keysAndValues...) }
func Warn(args ...interface{}) { defaultLogger.Warn(args...) }
func Warnf(template string, args ...interface{}) { defaultLogger.Warnf(template, args...) }
func Warnw(msg string, keysAndValues ...interface{}) { defaultLogger.Warnw(msg, keysAndValues...) }
func Sync() error { return defaultLogger.Sync() }






