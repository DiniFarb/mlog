package mlog

import (
	"fmt"
	"sync"
	"time"
)

type Logger struct {
	timeformat string
	appName    string
	level      Level
	format     Format
	custom     func(LogLine)
}

type LogLine struct {
	Timestamp string `json:"timestamp"`
	AppName   string `json:"appname"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type Format int
type Level int

const (
	_ Format = iota
	Ftext
	Fjson
	FCustom
)

const (
	_ Level = iota
	Lerror
	Lwarn
	Linfo
	Ldebug
	Ltrace
)

var logger *Logger
var once sync.Once

func Trace(message string, args ...interface{}) {
	create()
	if logger.level >= Ltrace {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   logger.appName,
			Level:     "TRACE",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

func Debug(message string, args ...interface{}) {
	create()
	if logger.level >= Ldebug {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   logger.appName,
			Level:     "DEBUG",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

func Info(message string, args ...interface{}) {
	create()
	if logger.level >= Linfo {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   logger.appName,
			Level:     "INFO",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

func Warn(message string, args ...interface{}) {
	create()
	if logger.level >= Lwarn {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   logger.appName,
			Level:     "WARN",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

func Error(message string, args ...interface{}) {
	create()
	if logger.level >= Lerror {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   logger.appName,
			Level:     "ERROR",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

func createTimeStamp() string {
	return time.Now().Format(logger.timeformat)
}

func log(logine LogLine) {
	switch logger.format {
	case Fjson:
		logJson(logine)
	case FCustom:
		logger.custom(logine)
	default:
		logText(logine)
	}
}

func logJson(logine LogLine) {
	fmt.Printf("{\"timestamp\":\"%s\",\"appname\":\"%s\",\"level\":\"%s\",\"message\":\"%s\"}\n", logine.Timestamp, logine.AppName, logine.Level, logine.Message)
}

func logText(logine LogLine) {
	fmt.Printf("[%s] %s | %s | %s\n", logine.AppName, logine.Timestamp, logine.Level, logine.Message)
}

func create() {
	once.Do(func() {
		logger = &Logger{
			timeformat: "2006-01-02 15:04:05.000",
			appName:    "MLOG",
			level:      Ltrace,
			custom:     logText,
			format:     Ftext,
		}
	})
}

func SetAppName(appName string) {
	create()
	logger.appName = appName
}

func SetLevel(level Level) {
	create()
	logger.level = level
}

func SetTimeFormat(timeformat string) {
	create()
	logger.timeformat = timeformat
}

func SetFormat(format Format) {
	create()
	logger.format = format
}

func SetCustomFormat(log func(LogLine)) {
	create()
	logger.format = FCustom
	logger.custom = log
}
