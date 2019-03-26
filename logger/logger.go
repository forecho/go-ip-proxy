package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"sync"
	"time"
)

var singleton *zap.Logger
var once sync.Once

/**
 * 初始化日志组件
 * logPath 日志路径
 * loglevel 日志级别
 * enableDatePrefix 是否启用日志前缀
 */
func Config(logPath string, loglevel string, enableDatePrefix bool) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   newLogPath(logPath, enableDatePrefix), // ⽇志⽂件路径
		MaxSize:    1024,                                  // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 3,                                     // 最多保留3个备份
		MaxAge:     7,                                     // 文件最多保存多少天
		Compress:   true,                                  // 是否压缩 disabled by default
	}
	w := zapcore.AddSync(&hook)
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)

	once.Do(func() {
		singleton = zap.New(core)
	})
	singleton.Info("DefaultLogger init success")
	return singleton
}

// Debug logs a debug message with the given fields
func Debug(message string, fields ...zap.Field) {
	singleton.Debug(message, fields...)
}

func Debugf(message string, params ...interface{}) {
	Debug(fmt.Sprintf(message, params))
}

// Info logs a debug message with the given fields
func Info(message string, fields ...zap.Field) {
	singleton.Info(message, fields...)
}
func Infof(message string, params ...interface{}) {
	Info(fmt.Sprintf(message, params))
}

// Warn logs a debug message with the given fields
func Warn(message string, fields ...zap.Field) {
	singleton.Warn(message, fields...)
}

func Warnf(message string, params ...interface{}) {
	Warn(fmt.Sprintf(message, params))
}
// Error logs a debug message with the given fields
func Error(message string, fields ...zap.Field) {
	singleton.Error(message, fields...)
}
// Error logs a debug message with the given fields
func Errorf(message string, params ...interface{}) {
	Error(fmt.Sprintf(message, params))
}

// Fatal logs a message than calls os.Exit(1)
func Fatal(message string, fields ...zap.Field) {
	singleton.Fatal(message, fields...)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func newLogPath(logPath string, enableDatePrefix bool) string {
	if !enableDatePrefix {
		return logPath
	}
	dir, file := filepath.Split(logPath)
	newFile := fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), file)
	return dir + newFile
}
