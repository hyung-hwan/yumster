package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Log_Formatter struct {
}

func (f *Log_Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var (
		now time.Time
		msg string
	)
	now = time.Now()
	msg = fmt.Sprintf("%s [%d] %s\n",
		now.Format("2006-01-02 15:04:05 -0700"),
		os.Getpid(), entry.Message)
	return []byte(msg), nil
}

type Logger struct {
	logrus.Logger
}

const LogDebugLevel = logrus.DebugLevel
const LogInfoLevel = logrus.InfoLevel
const LogWarnLevel = logrus.WarnLevel
const LogErrorLevel = logrus.ErrorLevel

func NewLogger() *Logger {
	var (
		logger *Logger
	)

	logger = &Logger{}

	// set default configurations
	//log.SetOutput(io.MultiWriter(os.Stderr))
	logger.SetOutput(os.Stderr)
	logger.SetFormatter(&Log_Formatter{})
	logger.SetLevel(logrus.InfoLevel)

	return logger
}

func (logger *Logger) Print(msg ...interface{}) {
	logger.Info(msg...)
}

/*
func (logger *Logger) Debugx(msg string) {
	var id string = obj.GetID()
	if id == "" {
		logger.Debug(msg)
	} else {
		logger.Debugf("%s: %s", id, msg)
	}
}

func (logger *Logger) Infox(obj ISvc, msg string) {
	var id string = obj.GetID()
	if id == "" {
		logger.Info(msg)
	} else {
		logger.Infof("%s: %s", id, msg)
	}
}

func (logger *Logger) Warnx(obj ISvc, msg string) {
	var id string = obj.GetID()
	if id == "" {
		logger.Warn(msg)
	} else {
		logger.Warnf("%s: %s", id, msg)
	}
}

func (logger *Logger) Errorx(obj ISvc, msg string) {
	var id string = obj.GetID()
	if id == "" {
		logger.Error(msg)
	} else {
		logger.Error("%s: %s", id, msg)
	}
}

func (logger *Logger) Debugfx(obj ISvc, fmt string, a ...interface{}) {
	var id string = obj.GetID()
	if id == "" {
		logger.Debugf(fmt, a...)
	} else {
		logger.Debugf(("%s: " + fmt), append([]interface{}{id}, a...)...)
	}
}

func (logger *Logger) Infofx(obj ISvc, fmt string, a ...interface{}) {
	var id string = obj.GetID()
	if id == "" {
		logger.Infof(fmt, a...)
	} else {
		logger.Infof(("%s: " + fmt), append([]interface{}{id}, a...)...)
	}
}

func (logger *Logger) Warnfx(obj ISvc, fmt string, a ...interface{}) {
	var id string = obj.GetID()
	if id == "" {
		logger.Warnf(fmt, a...)
	} else {
		logger.Warnf(("%s: " + fmt), append([]interface{}{id}, a...)...)
	}
}

func (logger *Logger) Errorfx(obj ISvc, fmt string, a ...interface{}) {
	var id string = obj.GetID()
	if id == "" {
		logger.Errorf(fmt, a...)
	} else {
		logger.Errorf(("%s: " + fmt), append([]interface{}{id}, a...)...)
	}
}
*/
