
package log

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)



func TestLogCreator(t *testing.T) {
	content := `
level: info
disableCaller: false
disableStacktrace: false
sampling:
  initial: 100
  thereafter: 100
encoding: console
encoderConfig:
  messageKey: msg
  levelKey: lvl
  timeKey: ts
  nameKey: logger
  callerKey: caller
  stacktraceKey: stacktrace
  lineEnding: "\n"
  levelEncoder: capitalColor
  timeEncoder: iso8601
  callerEncoder: short
stdout: true
strerr: false
file:
  filename: %v
  maxsize: 1
  maxage: 0
  localtime: true
  maxbackups: 3
  compress: false
`
	log_path := "test.log"
	content = fmt.Sprintf(content, log_path)
	var cfg *Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil { t.Fatal(err) }
	logger, err := NewLoggerWithConfig(cfg)
	if err != nil { t.Fatal(err) }

	logger.Infow("test", "filed1", 1, "field2", 2)

	os.Remove(log_path)
}