package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func blockGoroutine() {
	shutdownSignalsChannel := make(chan os.Signal)
	signal.Notify(shutdownSignalsChannel, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case s := <-shutdownSignalsChannel:
			log.Info().Msgf("Unblocking goroutine due to a signal: %v", s)
			return
		}
	}
}
