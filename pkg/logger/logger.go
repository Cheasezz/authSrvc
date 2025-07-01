package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(message any, args ...any)
	Info(message string, args ...any)
	Error(message any, args ...any)
	Fatal(message any, args ...any)
}

type lg struct {
	logger *logrus.Logger
}

func New(level string) Logger {
	var l logrus.Level

	switch strings.ToLower(level) {
	case "error":
		l = logrus.ErrorLevel
	case "warn":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.SetLevel(l)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	return &lg{
		logger: logger,
	}
}

func (l *lg) Debug(message any, args ...any) {
	l.logger.Debugf(toFormat(message), args...)
}

func (l *lg) Info(message string, args ...any) {
	l.logger.Infof(message, args...)
}

func (l *lg) Error(message any, args ...any) {
	l.logger.Errorf(toFormat(message), args...)
}

func (l *lg) Fatal(message any, args ...any) {
	l.logger.Fatalf(toFormat(message), args...)
}

func toFormat(m any) string {
	switch v := m.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return "%v"
	}
}
