package audit

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	loger  *lumberjack.Logger
	logger *logrus.Logger
)

func init() {
	loger = &lumberjack.Logger{
		Filename:   "logs/audit.log",
		MaxSize:    50, // megabytes
		MaxBackups: 300,
		MaxAge:     90, // days
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			loger.Rotate()
		}
	}()
	logger = logrus.New()
	logger.SetOutput(loger)
	logger.SetFormatter(&logrus.JSONFormatter{})
}

func Logger() *logrus.Logger {
	return logger
}
