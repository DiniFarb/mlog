package mlog

import (
	"fmt"
)

func Trace(message string, args ...interface{}) {
	create()
	if mlogger.level >= Ltrace {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   mlogger.appName,
			Level:     "TRACE",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

func Debug(message string, args ...interface{}) {
	create()
	if mlogger.level >= Ldebug {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   mlogger.appName,
			Level:     "DEBUG",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

func Info(message string, args ...interface{}) {
	create()
	if mlogger.level >= Linfo {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   mlogger.appName,
			Level:     "INFO ",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

func Warn(message string, args ...interface{}) {
	create()
	if mlogger.level >= Lwarn {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   mlogger.appName,
			Level:     "WARN ",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

func Error(message string, args ...interface{}) {
	create()
	if mlogger.level >= Lerror {
		logline := LogLine{
			Timestamp: createTimeStamp(),
			AppName:   mlogger.appName,
			Level:     "ERROR",
			Message:   fmt.Sprintf(message, args...),
		}
		log(logline)
	}
}

// Used for Errors in MLog itself, can be used in AddCustomOutput f.e.
func MLogError(message string, args ...interface{}) {
	logline := LogLine{
		Timestamp: createTimeStamp(),
		AppName:   mlogger.appName,
		Level:     "MLOG ERROR",
		Message:   fmt.Sprintf(message, args...),
	}
	fmt.Println(ApplyFormat(logline))
}

func ApplyFormat(logine LogLine) string {
	formated := ""
	switch mlogger.format {
	case Fjson:
		formated = formatJson(logine)
	case FCustom:
		formated = mlogger.customFormat(logine)
	default:
		formated = formatDefaultText(logine)
	}
	return formated
}

func SetAppName(appName string) {
	create()
	mlogger.appName = appName
}

func SetLevel(level Level) {
	create()
	mlogger.level = level
}

func SetTimeFormat(timeformat string) {
	create()
	mlogger.timeformat = timeformat
}

func SetFormat(format Format) {
	create()
	mlogger.format = format
}

func SetCustomFormat(log func(LogLine) string) {
	create()
	mlogger.format = FCustom
	mlogger.customFormat = log
}

// This will start in memory queue for logs.
// If the out function returns false, the log will be put back in the queue.
func AddCustomOutput(out func(LogLine) bool) {
	create()
	queue := newQueue(out)
	mlogger.outputQueues = append(mlogger.outputQueues, queue)
}
