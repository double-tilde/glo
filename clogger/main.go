package clogger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// clogger - colourful log

const (
	LevelTrace = -8
	LevelDebug = -4
	LevelInfo  = 0
	LevelWarn  = 4
	LevelError = 8
	LevelFatal = 12
)

var levelNames = map[int]string{
	LevelTrace: "TRACE",
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

var levelColors = map[int]string{
	LevelTrace: "\033[36m", // Cyan
	LevelDebug: "\033[34m", // Blue
	LevelInfo:  "\033[32m", // Green
	LevelWarn:  "\033[33m", // Yellow
	LevelError: "\033[31m", // Red
	LevelFatal: "\033[35m", // Magenta
}

const resetColor = "\033[0m"

type contextKey string

const (
	keyCallerFile contextKey = "caller_file"
	keyCallerLine contextKey = "caller_line"
)

type Clogger struct {
	jsonLogFile *os.File
}

func New(filePath string) *Clogger {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("failed to open log file: ", err)
	}

	return &Clogger{
		jsonLogFile: f,
	}
}

func (l *Clogger) log(level int, msg string, err error) {
	_, file, line, _ := runtime.Caller(2)
	shortFile := shortFilePath(file)

	ctx := context.WithValue(context.Background(), keyCallerFile, shortFile)
	ctx = context.WithValue(ctx, keyCallerLine, line)

	color := levelColors[level]
	levelLabel := levelNames[level]
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s[%s]%s %s %s:%d %s - details: %s\n",
		color, levelLabel, resetColor, timestamp, shortFile, line, msg, err)

	entry := map[string]any{
		"time":    time.Now().Format(time.RFC3339),
		"level":   levelLabel,
		"msg":     msg,
		"file":    ctx.Value(keyCallerFile),
		"line":    ctx.Value(keyCallerLine),
		"details": err,
	}
	_ = l.writeJSON(entry)

	if level == LevelFatal {
		os.Exit(1)
	}
}

func shortFilePath(fullPath string) string {
	wd, err := os.Getwd()
	if err != nil {
		return fullPath
	}
	rel, err := filepath.Rel(wd, fullPath)
	if err != nil {
		return fullPath
	}
	return rel
}

func (l *Clogger) writeJSON(entry map[string]any) error {
	enc := json.NewEncoder(l.jsonLogFile)
	return enc.Encode(entry)
}

func (l *Clogger) Trace(msg string, err error) { l.log(LevelTrace, msg, err) }
func (l *Clogger) Debug(msg string, err error)  { l.log(LevelDebug, msg, err) }
func (l *Clogger) Info(msg string, err error)   { l.log(LevelInfo, msg, err) }
func (l *Clogger) Warn(msg string, err error)   { l.log(LevelWarn, msg, err) }
func (l *Clogger) Error(msg string, err error)  { l.log(LevelError, msg, err) }
func (l *Clogger) Fatal(msg string, err error)  { l.log(LevelFatal, msg, err) }

func (l *Clogger) Close() error {
	return l.jsonLogFile.Close()
}
