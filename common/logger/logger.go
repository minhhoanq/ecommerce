package logger

import (
	"log"

	"github.com/minhhoanq/ecommerce/common/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Interface interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	DPanic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

type Level = zapcore.Level

type wrapLogger struct {
	l *zap.Logger
}

var _ Interface = (*wrapLogger)(nil)

const (
	InfoLevel   Level = zap.InfoLevel   //0, default level
	WarnLevel   Level = zap.WarnLevel   //1
	ErrorLevel  Level = zap.ErrorLevel  //2
	DPanicLevel Level = zap.DPanicLevel //3, used in development log
	//PanicLevel logs a message, then panic
	PanicLevel Level = zap.PanicLevel //4
	//FatalLevel logs a message, then calls os.Exit(1)
	FatalLevel Level = zap.FatalLevel
	DebugLevel Level = zap.DebugLevel
)

type Field = zap.Field

var (
	Skip       = zap.Skip
	Binary     = zap.Binary
	Bool       = zap.Bool
	ByteString = zap.ByteString
	Complex128 = zap.Complex128

	Complex64 = zap.Complex64

	Float64 = zap.Float64
	Float32 = zap.Float32
	Int     = zap.Int
	Int64   = zap.Int64

	Int32     = zap.Int32
	Int16     = zap.Int16
	Int8      = zap.Int8
	String    = zap.String
	Uint      = zap.Uint
	Uint64    = zap.Uint64
	Uint32    = zap.Uint32
	Uint16    = zap.Uint16
	Uint8     = zap.Uint8
	Uintptr   = zap.Uintptr
	Reflect   = zap.Reflect
	Namespace = zap.Namespace
	Stringer  = zap.Stringer
	Time      = zap.Time
	Stack     = zap.Stack
	Duration  = zap.Duration
	Any       = zap.Any
)

// NewWrapLogger initializes a new logger
func NewWrapLogger(level Level, isDevelopment bool) *wrapLogger {
	var cfg zap.Config

	// Configure logger based on environment
	if isDevelopment {
		cfg = zap.NewDevelopmentConfig()

		cfg.EncoderConfig = zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Log màu sắc
			EncodeTime:     zapcore.ISO8601TimeEncoder,       // Format thời gian chuẩn ISO8601
			EncodeCaller:   zapcore.ShortCallerEncoder,       // Rút gọn đường dẫn file
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}

		// CAdjust the output to ConsoleEncoder (instead of JSONEncoder)
		cfg.Encoding = "console"

		cfg.Level = zap.NewAtomicLevelAt(level)
	} else {
		cfg = zap.NewProductionConfig()
	}

	// Set log level
	cfg.Level = zap.NewAtomicLevelAt(level)

	// Build the logger
	logger, err := cfg.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// Wrap the logger
	return &wrapLogger{l: logger}
}

var w *wrapLogger

func init() {
	w = NewWrapLogger(DebugLevel, false)
}

func GetLogLevel(l string) Level {
	var zapLogLvl zapcore.Level
	err := zapLogLvl.Set(l)
	if err != nil {
		log.Println("Cannot parse loglevel, err:", err.Error())
		zapLogLvl = zap.ErrorLevel
	}
	return zapLogLvl
}

func Setup(enviroment string, l string) {
	w = NewWrapLogger(GetLogLevel(l), enviroment == constants.LocalEnvName)
}

func With(fields ...Field) {
	w.l.With(fields...)
}

func (l *wrapLogger) Debug(msg string, fields ...Field) {
	w.l.Debug(msg, fields...)
}

func (l *wrapLogger) Info(msg string, fields ...Field) {
	w.l.Info(msg, fields...)
}

func (l *wrapLogger) Warn(msg string, fields ...Field) {
	w.l.Warn(msg, fields...)
}

func (l *wrapLogger) Error(msg string, fields ...Field) {
	w.l.Error(msg, fields...)
}

func (l *wrapLogger) DPanic(msg string, fields ...Field) {
	w.l.DPanic(msg, fields...)
}

func (l *wrapLogger) Panic(msg string, fields ...Field) {
	w.l.Panic(msg, fields...)
}

func (l *wrapLogger) Fatal(msg string, fields ...Field) {
	w.l.Fatal(msg, fields...)
}
