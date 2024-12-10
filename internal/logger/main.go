package logger

import (
	"fmt"
	"os"
)

type Level int

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

var level Level = LevelInfo

func SetLevel(l Level) {
	level = l
}

func isDebug() bool {
	return level == LevelDebug
}

func SetLevelDebug(active bool) {
	if active {
		level = LevelDebug
	}
}

func Infof(message string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, message+"\n", args...)
}

func Info(message string) {
	fmt.Fprint(os.Stdout, message)
}

func Debugf(message string, args ...interface{}) {
	if isDebug() {
		fmt.Fprintf(os.Stdout, message+"\n", args...)
	}
}

func Debug(message string) {
	if isDebug() {
		fmt.Fprint(os.Stdout, message)
	}
}

func Errorf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, message+"\n", args...)
}

func Error(message string) {
	fmt.Fprint(os.Stdout, message)
}
