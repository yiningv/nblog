package log

import (
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var (
	logger *zap.Logger
	once   sync.Once
)

type ZapConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool

	Level         int8
	IsDev         bool
	StdLogEnable  bool
	FileLogEnable bool
}

// 初始化zap log
func InitLogByConfig(conf *ZapConfig) {
	once.Do(func() {
		newZapLogger(conf)
	})
}

// 初始化zap log
func InitLogByToml(tomlPath string) {
	conf := new(ZapConfig)
	if _, err := toml.DecodeFile(tomlPath, conf); err != nil {
		panic(err)
	}
	once.Do(func() {
		newZapLogger(conf)
	})
}

// 创建zap的logger实例
func newZapLogger(c *ZapConfig) {
	encoderConf := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if c.IsDev {
		encoder = zapcore.NewConsoleEncoder(encoderConf)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConf)
	}

	var writeSyncers []zapcore.WriteSyncer
	if c.StdLogEnable {
		writeSyncers = append(writeSyncers, zapcore.Lock(os.Stdout))
	}

	if c.FileLogEnable {
		hook := &lumberjack.Logger{
			Filename:   c.Filename,
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			LocalTime:  c.LocalTime,
			Compress:   c.Compress,
		}
		writeSyncers = append(writeSyncers, zapcore.AddSync(hook))
	}

	writeSyncer := zapcore.NewMultiWriteSyncer(writeSyncers...)

	atomicLevel := zap.NewAtomicLevelAt(zapcore.Level(c.Level))

	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		atomicLevel,
	)

	caller := zap.AddCaller()
	skip := zap.AddCallerSkip(1)
	logger = zap.New(core, caller, skip)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zapcore.Field) {
	logger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanic(msg string, fields ...zapcore.Field) {
	logger.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zapcore.Field) {
	logger.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zapcore.Field) {
	logger.Fatal(msg, fields...)
}
