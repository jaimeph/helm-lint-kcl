package logger

import "fmt"

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

func SetLevelDebug(active bool) {
	if active {
		level = LevelDebug
	}
}

func Infof(message string, args ...interface{}) {
	fmt.Printf(" INFO "+message+"\n", args...)
}

func Info(message string) {
	fmt.Println(" INFO " + message)
}

func Debugf(message string, args ...interface{}) {
	if level >= LevelDebug {
		fmt.Printf("DEBUG "+message+"\n", args...)
	}
}

func Debug(message string) {
	if level >= LevelDebug {
		fmt.Println("DEBUG " + message)
	}
}

func Errorf(message string, args ...interface{}) {
	fmt.Printf("ERROR "+message+"\n", args...)
}

func Error(message string) {
	fmt.Println("ERROR " + message)
}