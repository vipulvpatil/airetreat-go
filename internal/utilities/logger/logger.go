package logger

import (
	"fmt"

	"github.com/getsentry/sentry-go"
)

var logMode = "stdout"

func LogUsingSentry(sentryDsn, environment string) error {
	logMode = "sentry"
	return sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDsn,
		TracesSampleRate: 1.0,
		Environment:      environment,
	})
}

func LogToStdout() {
	logMode = "stdout"
}

func getLogger() loggerInterface {
	switch logMode {
	case "stdout":
		return &stdoutLogger{}
	case "sentry":
		return &sentryLogger{}
	}
	return nil
}

func LogMessageln(a ...any) {
	l := getLogger()
	l.logMessageln(a...)
}

func LogMessagef(format string, a ...any) {
	l := getLogger()
	l.logMessagef(format, a...)
}

func LogError(err error) {
	l := getLogger()
	l.logError(err)
}

type loggerInterface interface {
	logMessageln(a ...any)
	logMessagef(format string, a ...any)
	logError(err error)
}

type stdoutLogger struct{}

func (l *stdoutLogger) logMessageln(a ...any) {
	fmt.Println(a...)
}

func (l *stdoutLogger) logMessagef(format string, a ...any) {
	fmt.Printf(format, a...)
}

func (l *stdoutLogger) logError(err error) {
	//do nothing on stdout
}

type sentryLogger struct{}

func (l *sentryLogger) logMessageln(a ...any) {
	sentry.CaptureMessage(fmt.Sprintln(a...))
}

func (l *sentryLogger) logMessagef(format string, a ...any) {
	sentry.CaptureMessage(fmt.Sprintf(format, a...))
}

func (l *sentryLogger) logError(err error) {
	sentry.CaptureException(err)
}
