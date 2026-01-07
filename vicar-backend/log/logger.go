package log

import (
	"io"
	"time"

	"github.com/fatih/color"
)

var minLogLevel = INFO
var logWriter *io.Writer
var isDebug = false

func SetDebug(debug bool) {
	isDebug = debug
	if debug {
		minLogLevel = DEBUG
	}
}

func IsDebug() bool {
	return isDebug
}

// SetLogWriter sets the writer to write logs to.
func SetLogWriter(writer io.Writer) {
	logWriter = &writer
}

func Log(area LogArea, level Level, emote, message string, args ...any) {
	if level < minLogLevel {
		return
	}

	time := time.Now().Format("2006-01-02 15:04:05")
	areaString := area.ColoredString()
	levelString := level.ColoredString(emote)
	msgString := color.New(color.FgWhite).Sprintf(message, args...)

	msg := time + " - " + areaString + levelString + msgString

	if logWriter != nil {
		(*logWriter).Write([]byte(msg + "\n"))
	} else {
		println(msg)
	}
}

func Debug(area LogArea, emote, message string, args ...interface{}) {
	Log(area, DEBUG, emote, message, args...)
}

func Info(area LogArea, emote, message string, args ...interface{}) {
	Log(area, INFO, emote, message, args...)
}

func Success(area LogArea, emote, message string, args ...interface{}) {
	Log(area, SUCCESS, emote, message, args...)
}

func Warning(area LogArea, emote, message string, args ...interface{}) {
	Log(area, WARNING, emote, message, args...)
}

func Error(area LogArea, emote, message string, args ...interface{}) {
	Log(area, ERROR, emote, message, args...)
}
