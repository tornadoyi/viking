package log

import (
	"fmt"
	"go.uber.org/zap/zapcore"
)

func NewEncoder(encoding string, config EncoderConfig) (Encoder, error){
	switch encoding {
	case "console": return zapcore.NewConsoleEncoder(config), nil
	case "json": return zapcore.NewJSONEncoder(config), nil
	}
	return nil, fmt.Errorf("no encoder registered for name %q", encoding)
}
