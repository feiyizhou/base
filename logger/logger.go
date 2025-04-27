package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConf struct {
	Level      string `json:"level" mapstructure:"level"`
	Filename   string `json:"filename" mapstructure:"filename"`
	MaxSize    int    `json:"maxSize" mapstructure:"maxSize"`
	MaxAge     int    `json:"maxAge" mapstructure:"maxAge"`
	MaxBackups int    `json:"maxBackups" mapstructure:"maxBackups"`
	Compress   bool   `json:"compress" mapstructure:"compress"`
}

// InitLogger init logger
func InitLogger(cfg LogConf) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		parseLogLevel(cfg.Level),
	)
	logger := zap.New(core, zap.AddStacktrace(zap.PanicLevel))
	zap.ReplaceGlobals(logger)
}

func parseLogLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug", "DEBUG":
		return zapcore.DebugLevel
	case "info", "INFO":
		return zapcore.InfoLevel
	case "warn", "WARN":
		return zapcore.WarnLevel
	case "error", "ERROR":
		return zapcore.ErrorLevel
	case "dpanic", "DPANIC":
		return zapcore.DPanicLevel
	case "panic", "PANIC":
		return zapcore.PanicLevel
	case "fatal", "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Debug(args ...any) {
	zap.L().Sugar().Debug(args...)
}

func Info(args ...any) {
	zap.L().Sugar().Infoln(args...)
}

func Warn(args ...any) {
	zap.L().Sugar().Warnln(args...)
}

func Error(args ...any) {
	zap.L().Sugar().Errorln(args...)
}

func Debugf(template string, args ...any) {
	zap.L().Sugar().Debugf(template, args...)
}

func Infof(template string, args ...any) {
	zap.L().Sugar().Infof(template, args...)
}

func Warnf(template string, args ...any) {
	zap.L().Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...any) {
	zap.L().Sugar().Errorf(template, args...)
}
