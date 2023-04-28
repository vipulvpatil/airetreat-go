package logger

import (
	"fmt"
)

var logMode = "stdout"

func LogToSentry() {
	logMode = "sentry"
}

func LogToStdout() {
	logMode = "stdout"
}

func LogMessageln(a ...any) {
	fmt.Println(a...)
}

func LogMessagef(format string, a ...any) {
	fmt.Printf(format, a...)
}

func LogError(err error) {
	//do nothing on stdout
}
