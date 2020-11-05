package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var DefaultLogger = NewLogger()

func NewLogger() *logrus.Logger {
	return &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}

func SetLevel(lvl logrus.Level) {
	DefaultLogger.SetLevel(lvl)
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}
