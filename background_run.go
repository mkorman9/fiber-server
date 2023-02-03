package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func runInBackground(fn func() error) {
	errorsChannel := make(chan struct{}, 1)
	signalsChannel := make(chan os.Signal)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Stack().
					Err(fmt.Errorf("%v", r)).
					Msg("Panic while running task in background")
			}
		}()

		if err := fn(); err != nil {
			select {
			case errorsChannel <- struct{}{}:
			default:
			}
		}
	}()

	for {
		select {
		case <-errorsChannel:
			return
		case <-signalsChannel:
			return
		}
	}
}
