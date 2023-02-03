package main

import (
	"os"
	"os/signal"
	"syscall"
)

type signalLock struct {
	ch chan struct{}
}

func newSignalLock() *signalLock {
	return &signalLock{
		ch: make(chan struct{}),
	}
}

func (s *signalLock) Unblock() {
	select {
	case s.ch <- struct{}{}:
	default:
	}
}

func (s *signalLock) Block() {
	signalsChannel := make(chan os.Signal)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-s.ch:
			return
		case <-signalsChannel:
			return
		}
	}
}
