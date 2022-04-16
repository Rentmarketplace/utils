package utils_v1_0

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// Log Export global
var Log *logrus.Logger

func Logger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  os.Getenv("LOG_PATH") + "Info-" + time.Now().Format("01-02-2006") + ".log",
		logrus.ErrorLevel: os.Getenv("LOG_PATH") + "Error-" + time.Now().Format("01-02-2006") + ".log",
		logrus.WarnLevel:  os.Getenv("LOG_PATH") + "Warning-" + time.Now().Format("01-02-2006") + ".log",
		logrus.DebugLevel: os.Getenv("LOG_PATH") + "Debug-" + time.Now().Format("01-02-2006") + ".log",
		logrus.TraceLevel: os.Getenv("LOG_PATH") + time.Now().Format("01-02-2006") + ".test.log",
	}

	Log = logrus.New()
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	return Log
}
