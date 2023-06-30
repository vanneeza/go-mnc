package log

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	logDir := filepath.Join("utils", "log")
	logPath := filepath.Join(logDir, "activity.log")

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		Logger.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			writers := []io.Writer{
				file,
				os.Stdout,
			}
			Logger.SetOutput(io.MultiWriter(writers...))
		} else {
			Logger.SetOutput(os.Stdout)
		}
	}

	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&logrus.JSONFormatter{})
}

func Info(message string) {
	Logger.Info(message)
}

func Error(message string, err error) {
	Logger.WithError(err).Error(message)
}
