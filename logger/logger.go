package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

type Logger struct {
	Log *zap.SugaredLogger
}

func NewLogger() *Logger {
	now := time.Now()
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%04d%02d%02d%s", "log", now.Year(), now.Month(), now.Day(), ".log"),
		MaxSize:    2,
		MaxBackups: 100,
		MaxAge:     30,
		Compress:   false,
	})
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), w, zap.NewAtomicLevel())
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	return &Logger{
		Log: logger,
	}
}
