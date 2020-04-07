package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DebugLevel = zap.DebugLevel
	InfoLevel = zap.InfoLevel
	WarnLevel = zap.WarnLevel
	ErrorLevel = zap.ErrorLevel
	DPanicLevel = zap.DPanicLevel
	PanicLevel = zap.PanicLevel
	FatalLevel = zap.FatalLevel
)


type Level							= zap.AtomicLevel
type SamplingConfig					= zap.SamplingConfig
type Option							= zap.Option

type EncoderConfig					= zapcore.EncoderConfig
type Encoder						= zapcore.Encoder
type WriteSyncer					= zapcore.WriteSyncer

type RotateLogger					= lumberjack.Logger


var AddSync							= zapcore.AddSync
var NewMultiWriteSyncer				= zapcore.NewMultiWriteSyncer
var NewCore							= zapcore.NewCore