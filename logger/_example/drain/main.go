package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/blend/go-sdk/logger"
)

func main() {
	log := logger.NewFromEnv()

	log.Listen(logger.Info, "randomly_slow", func(e logger.Event) {
		if rand.Float64() < 0.1 {
			println("randomly slow start")
			time.Sleep(500 * time.Millisecond)
			println("randomly slow stop")
		}
	})

	infoSignal := time.Tick(100 * time.Millisecond)

	done := time.After(10 * time.Second)

	for {
		select {
		case <-infoSignal:
			log.Infof("this is an info event")
		case <-done:
			println("draining")
			log.Drain()
			println("exiting")
			os.Exit(0)
		}
	}
}
