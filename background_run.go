package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func runInBackground(fn func()) {
	doneChannel := make(chan struct{}, 1)
	signalsChannel := make(chan os.Signal)

	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(signalsChannel)
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Stack().
					Err(fmt.Errorf("%v", r)).
					Msg("Panic while running task in background")
			}
		}()

		fn()
		doneChannel <- struct{}{}
	}()

	for {
		select {
		case <-doneChannel:
			return
		case <-signalsChannel:
			return
		}
	}
}
