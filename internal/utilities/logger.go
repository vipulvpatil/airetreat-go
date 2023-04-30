package utilities

import (
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

type LoggerParams struct {
	Mode         string
	SentryParams struct {
		Dsn         string
		Environment string
	}
}

func InitLogger(params LoggerParams) (Logger, func(time.Duration) bool, error) {
	switch params.Mode {
	case "stdout":
		return &stdoutLogger{}, nil, nil
	case "sentry":
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              params.SentryParams.Dsn,
			TracesSampleRate: 1.0,
			Environment:      params.SentryParams.Environment,
		})
		return &sentryLogger{}, sentry.Flush, err
	}
	return nil, nil, errors.New("incorrect logger mode")
}

type Logger interface {
	LogMessageln(a ...any)
	LogMessagef(format string, a ...any)
	LogError(err error)
}

type stdoutLogger struct{}

func (l *stdoutLogger) LogMessageln(a ...any) {
	fmt.Println(a...)
}

func (l *stdoutLogger) LogMessagef(format string, a ...any) {
	fmt.Printf(format, a...)
}

func (l *stdoutLogger) LogError(err error) {
	//do nothing on stdout
}

type sentryLogger struct{}

func (l *sentryLogger) LogMessageln(a ...any) {
	sentry.CaptureMessage(fmt.Sprintln(a...))
}

func (l *sentryLogger) LogMessagef(format string, a ...any) {
	sentry.CaptureMessage(fmt.Sprintf(format, a...))
}

func (l *sentryLogger) LogError(err error) {
	sentry.CaptureException(err)
}
