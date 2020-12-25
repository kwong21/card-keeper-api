package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// APILogger is a wrapper struct around logrus logger
type APILogger struct {
	*logrus.Logger
}

// Fields wraps logrus.Fields, which is a map[string]interface{}
type Fields logrus.Fields

// NewLogger initializes the APILogger
func NewLogger() *APILogger {
	logFile, err := os.OpenFile("card-keeper-api.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println("Can't open log file", err)
	}

	baseLogger := logrus.New()

	baseLogger.Out = logFile

	baseLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	baseLogger.SetReportCaller(true)

	javLogger := &APILogger{baseLogger}

	return javLogger
}

// LogInfo writes Info log statements.
func (l *APILogger) LogInfo(message string) {
	l.Info(message)
}

// LogInfoWithFields writes Info log statements with fields
func (l *APILogger) LogInfoWithFields(f Fields, message string) {
	l.WithFields(logrus.Fields(f)).Info(message)
}

// LogErrorWithFields write Error log statements with fields
func (l *APILogger) LogErrorWithFields(f Fields, message string) {
	l.WithFields(logrus.Fields(f)).Error(message)
}

// LogWarnWithFields write Warn log statements with fields
func (l *APILogger) LogWarnWithFields(f Fields, message string) {
	l.WithFields(logrus.Fields(f)).Warn(message)
}
