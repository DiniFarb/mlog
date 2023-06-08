package mlog

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var mlogger *logger
var onceLogger sync.Once

type logger struct {
	timeformat   string
	appName      string
	level        Level
	format       Format
	customFormat func(LogLine) string
	outputQueues []*logQueue
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

func createTimeStamp() string {
	return time.Now().Format(mlogger.timeformat)
}

func formatJson(logine LogLine) string {
	l := strings.TrimSpace(logine.Level)
	return fmt.Sprintf("{\"timestamp\":\"%s\",\"appname\":\"%s\",\"level\":\"%s\",\"message\":\"%s\"}", logine.Timestamp, logine.AppName, l, logine.Message)
}

func formatDefaultText(logine LogLine) string {
	return fmt.Sprintf("[%s] %s | %s | %s", logine.AppName, logine.Timestamp, logine.Level, logine.Message)
}

func create() {
	onceLogger.Do(func() {
		mlogger = &logger{
			timeformat:   "2006-01-02 15:04:05.000",
			appName:      "MLOG",
			level:        Ltrace,
			customFormat: formatDefaultText,
			format:       Ftext,
			outputQueues: make([]*logQueue, 0),
		}
	})
}

func log(logline LogLine) {
	for _, queue := range mlogger.outputQueues {
		maxSize := int64(len(queue.queue)) >= queue.maxQueueSize
		for maxSize {
			time.Sleep(10 * time.Millisecond)
			maxSize = int64(len(queue.queue)) >= queue.maxQueueSize
		}
		queue.enqueue(logline)
	}
	fmt.Println(ApplyFormat(logline))
}
