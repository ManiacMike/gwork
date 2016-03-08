package gwork

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

type LogLevel uint8

// defined log levels
const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelNotice
	LogLevelWarning
	LogLevelError
)

type logData struct {
	level   LogLevel
	message string
}

func (this *logData) LevelString() (ret string) {
	switch this.level {
	case LogLevelDebug:
		ret = "DEBUG"
	case LogLevelInfo:
		ret = "INFO"
	case LogLevelNotice:
		ret = "NOTICE"
	case LogLevelWarning:
		ret = "WARNING"
	case LogLevelError:
		ret = "ERROR"
	}
	return
}

var (
	// a channel use to transfer log
	logChannel chan *logData
	logger     *log.Logger

	// a buffer use to write logs
	buf            bytes.Buffer
	currBufferSize uint16 = 0
)

// Start the log server
func logStart() {
	logChannel = make(chan *logData, conf.LogQueueSize)
	logger = log.New(&buf, "", log.LstdFlags)
	go func() {
		for {
			select {
			case aLog := <-logChannel:
				writeLog(aLog)
			}
		}
	}()
}

func Logf(level LogLevel, format string, a ...interface{}) {
	if level < conf.LogLevel {
		return
	}
	message := fmt.Sprintf(format, a...)
	aLog := &logData{level, message}
	logChannel <- aLog
}

func Log(level LogLevel, a ...interface{}) {
	if level < conf.LogLevel {
		return
	}
	message := fmt.Sprint(a...)
	aLog := &logData{level, message}
	logChannel <- aLog
}

// Write logs to buffer. Flush logs to stdout if buffer size reach the configured size.
func writeLog(aLog *logData) {
	logger.Printf("[%s] %s\n", aLog.LevelString(), aLog.message)
	if aLog.level == LogLevelError {
		defer func() { os.Exit(1) }()
		logger.Printf("SHUTDOWN\n")
	}
	currBufferSize++
	if currBufferSize >= conf.LogBufferSize {
		FlushLog()
	}
}

func FlushLog() {
	log.Print(&buf)
	buf.Reset()
	currBufferSize = 0
}
