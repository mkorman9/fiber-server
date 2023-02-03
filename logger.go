package main

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"time"
)

func setupLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	zerolog.TimestampFieldName = "time"
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = true
	zerolog.ErrorStackMarshaler = stackTraceMarshaller

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    !isatty.IsTerminal(os.Stdout.Fd()),
		TimeFormat: "2006-01-02 15:04:05",
	})
}

func stackTraceMarshaller(_ error) interface{} {
	var stackTrace []map[string]string

	for i := 3; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)

		stackTrace = append(stackTrace, map[string]string{
			"src":  fmt.Sprintf("%v:%v", file, line),
			"func": fn.Name(),
		})
	}

	return stackTrace
}
