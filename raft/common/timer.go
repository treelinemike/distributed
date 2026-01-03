package common

import (
	"log"
	"math/rand/v2"
	"time"
)

type (
	Timeout struct {
		Min_ms int `mapstructure:"Min_ms"`
		Max_ms int `mapstructure:"Max_ms"`
	}
	RaftTimer struct {
		Timer   *time.Timer
		Timeout Timeout
	}
)

func (t *RaftTimer) Reset() error {
	randtval := time.Duration(rand.IntN(t.Timeout.Max_ms-t.Timeout.Min_ms)+t.Timeout.Min_ms) * time.Millisecond
	t.Timer.Reset(randtval)
	log.Printf("Setting election timeout at %v", randtval)
	return nil
}

func (t *RaftTimer) Stop() error {
	t.Timer.Stop()
	log.Printf("Stopping election timer")
	return nil
}

func (t *RaftTimer) SetTimeout(to Timeout) error {
	t.Timeout = to
	return nil
}
