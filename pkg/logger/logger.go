package logger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

type Logger struct {
	env          string
	requestID    string
	minLevel     LogLevel
}

func New(env string) *Logger {
	minLevel := LevelInfo
	if env == "development" {
		minLevel = LevelDebug
	}
	
	return &Logger{
		env:      env,
		minLevel: minLevel,
	}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		env:       l.env,
		requestID: requestID,
		minLevel:  l.minLevel,
	}
}

func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
	}
	return levels[level] >= levels[l.minLevel]
}

func (l *Logger) format(level LogLevel, msg string, fields map[string]interface{}) string {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] [%s] [%s]", timestamp, level, l.env))
	
	if l.requestID != "" {
		sb.WriteString(fmt.Sprintf(" [req:%s]", l.requestID))
	}
	
	sb.WriteString(fmt.Sprintf(" %s", msg))
	
	if len(fields) > 0 {
		sb.WriteString(" |")
		for key, value := range fields {
			sb.WriteString(fmt.Sprintf(" %s=%v", key, value))
		}
	}
	
	return sb.String()
}

func (l *Logger) log(level LogLevel, msg string, fields map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}
	log.Println(l.format(level, msg, fields))
}

func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	fieldMap := map[string]interface{}{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}
	l.log(LevelDebug, msg, fieldMap)
}

func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	fieldMap := map[string]interface{}{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}
	l.log(LevelInfo, msg, fieldMap)
}

func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	fieldMap := map[string]interface{}{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}
	l.log(LevelWarn, msg, fieldMap)
}

func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	fieldMap := map[string]interface{}{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}
	
	if err != nil {
		fieldMap["error"] = err.Error()
		// Add stack trace for errors in development
		if l.env == "development" {
			fieldMap["stack"] = getStackTrace()
		}
	}
	
	l.log(LevelError, msg, fieldMap)
}

func getStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
