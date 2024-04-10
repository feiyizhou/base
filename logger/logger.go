package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConf struct {
	Level      string `json:"level" mapstructure:"level"`
	Filename   string `json:"filename" mapstructure:"filename"`
	MaxSize    int    `json:"maxSize" mapstructure:"maxSize"`
	MaxAge     int    `json:"maxAge" mapstructure:"maxAge"`
	MaxBackups int    `json:"maxBackups" mapstructure:"maxBackups"`
	Compress   bool   `json:"compress" mapstructure:"compress"`
}

// InitLogger 初始化Logger
func InitLogger(cfg LogConf) error {
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress)
	encoder := getEncoder()
	var level = new(zapcore.Level)
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return err
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	syncFile := zapcore.AddSync(lumberJackLogger)
	return zapcore.AddSync(syncFile)
}

func Debug(args ...interface{}) {
	zap.L().Sugar().Debugln(args)
}

func Info(args ...interface{}) {
	zap.L().Sugar().Infoln(args)
}

func Warn(args ...interface{}) {
	zap.L().Sugar().Warnln(args)
}

func Error(args ...interface{}) {
	zap.L().Sugar().Errorln(args)
}

func Debugf(template string, args ...interface{}) {
	zap.L().Sugar().Debugf(template, args)
}

func Infof(template string, args ...interface{}) {
	zap.L().Sugar().Infof(template, args)
}

func Warnf(template string, args ...interface{}) {
	zap.L().Sugar().Warnf(template, args)
}

func Errorf(template string, args ...interface{}) {
	zap.L().Sugar().Errorf(template, args)
}
