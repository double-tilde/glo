package logger

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

// TODO: Return errors, do not rely on other packages inside pkg

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

type logEntry struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
	File  string `json:"file"`
	Line  int    `json:"line"`
	Info  string `json:"info"`
}

const resetColor = "\033[0m"

type contextKey string

const (
	keyCallerFile contextKey = "caller_file"
	keyCallerLine contextKey = "caller_line"
)

// Clogger - colourful log

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

func (l *Clogger) writeJSON(entry logEntry) error {
	enc := json.NewEncoder(l.jsonLogFile)
	enc.SetIndent("", "  ")

	return enc.Encode(entry)
}

func getCallerInfo(level int) (int, string, context.Context, string, string, string) {
	_, file, line, _ := runtime.Caller(3)
	shortFile := shortFilePath(file)

	ctx := context.WithValue(context.Background(), keyCallerFile, shortFile)
	ctx = context.WithValue(ctx, keyCallerLine, line)

	color := levelColors[level]
	levelLabel := levelNames[level]
	timestamp := time.Now().Format("15:04:05")

	return line, shortFile, ctx, color, levelLabel, timestamp
}

func (l *Clogger) log(level int, msg string) {
	line, shortFile, ctx, color, levelLabel, timestamp := getCallerInfo(level)

	fmt.Printf("%s[%s]%s %s %s:%d %s\n",
		color, levelLabel, resetColor, timestamp, shortFile, line, msg)

	file, _ := ctx.Value(keyCallerFile).(string)
	lineNum, _ := ctx.Value(keyCallerLine).(int)

	entry := logEntry{
		Time:  time.Now().Format(time.RFC3339),
		Level: levelLabel,
		Msg:   msg,
		File:  file,
		Line:  lineNum,
		Info:  "",
	}

	_ = l.writeJSON(entry)
}

func (l *Clogger) logError(level int, msg string, err error) {
	line, shortFile, ctx, color, levelLabel, timestamp := getCallerInfo(level)

	fmt.Printf("%s[%s]%s %s %s:%d %s - %v\n",
		color, levelLabel, resetColor, timestamp, shortFile, line, msg, err)

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	file, _ := ctx.Value(keyCallerFile).(string)
	lineNum, _ := ctx.Value(keyCallerLine).(int)

	entry := logEntry{
		Time:  time.Now().Format(time.RFC3339),
		Level: levelLabel,
		Msg:   msg,
		File:  file,
		Line:  lineNum,
		Info:  errMsg,
	}

	_ = l.writeJSON(entry)

	if level == LevelFatal {
		os.Exit(1)
	}
}

func (l *Clogger) Trace(msg string) { l.log(LevelTrace, msg) }
func (l *Clogger) Debug(msg string) { l.log(LevelDebug, msg) }
func (l *Clogger) Info(msg string)  { l.log(LevelInfo, msg) }

func (l *Clogger) Warn(msg string, err error)  { l.logError(LevelWarn, msg, err) }
func (l *Clogger) Error(msg string, err error) { l.logError(LevelError, msg, err) }
func (l *Clogger) Fatal(msg string, err error) { l.logError(LevelFatal, msg, err) }

func (l *Clogger) Close() error {
	return l.jsonLogFile.Close()
}
