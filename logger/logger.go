package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   "log",
		MaxSize:    100, // 100M for a single file
		MaxBackups: 60,  // Maximum 60 log files
		MaxAge:     1,   // 1day
		Compress:   false,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(lumberJackLogger), zap.DebugLevel),
	)

	logger = zap.New(core)

	logger.Debug("Init Logger OK")
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Debug(msg)
}

func Debugln(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	Debug(msg)
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Info(msg)
}

func Infoln(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	Info(msg)
}

func Error(err error, msg string) {
	logger.Error(msg, zap.Error(err))
}

func Errorf(err error, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Error(err, msg)
}

func Errorln(err error, v ...interface{}) {
	msg := fmt.Sprintln(v...)
	Error(err, msg)
}

func Fatal(err error, msg string) {
	logger.Fatal(msg, zap.Error(err))
}

func Fatalf(err error, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Fatal(err, msg)
}

func Fatalln(err error, v ...interface{}) {
	msg := fmt.Sprintln(v...)
	Fatal(err, msg)
}

func Panic(err error, msg string) {
	logger.Panic(msg, zap.Error(err))
}

func Panicf(err error, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Panic(err, msg)
}

func Panicln(err error, v ...interface{}) {
	msg := fmt.Sprintln(v...)
	Panic(err, msg)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		logger.Debug("Http request",
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("remote address", c.Request.RemoteAddr),
			zap.String("method", c.Request.Method),
			zap.String("host", c.Request.Host),
			zap.String("uri", c.Request.RequestURI),
		)
	}
}
