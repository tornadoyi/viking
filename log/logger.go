package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sort"
	"time"
)

type Logger		= zap.SugaredLogger

func NewLoggerWithConfig(cfg *Config, opts ...Option) (*Logger, error) {
	// encoder
	encoder, err := NewEncoder(cfg.Encoding, cfg.EncoderConfig)
	if err != nil { return nil, err}

	// level
	if cfg.Level == (Level{}) { return nil, fmt.Errorf("missing Level") }

	// writers
	ws := make([]WriteSyncer, 0)
	if cfg.Stdout { ws = append(ws, os.Stdout) }
	if cfg.File != nil { ws = append(ws, AddSync(cfg.File)) }
	writer := NewMultiWriteSyncer(ws...)

	// options
	if cfg.Stderr { opts = append(opts, zap.ErrorOutput(os.Stderr)) }
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(ErrorLevel))
	}
	if cfg.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(cfg.Sampling.Initial), int(cfg.Sampling.Thereafter))
		}))
	}
	if len(cfg.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	log := zap.New(
		NewCore(encoder, writer, cfg.Level),
		opts...
	)
	return log.Sugar(), nil
}









