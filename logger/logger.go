package logger

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

func InitLogger() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   "log",
		MaxSize:    100, // 100M for a single file
		MaxBackups: 60,  // Maximum 60 log files
		MaxAge:     1,   // 1day
		Compress:   false,
		LocalTime:  true,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(lumberJackLogger), zap.DebugLevel),
	)

	logger = zap.New(core)

	logger.Debug("Init Logger OK")
}

func Sync() {
	logger.Sync()
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
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		logger.Debug("http requst",
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.String("remote address", c.Request.RemoteAddr),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("duration", duration),
		)
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error("[recovery for broken pipe]",
						zap.String("path", c.Request.URL.Path),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
