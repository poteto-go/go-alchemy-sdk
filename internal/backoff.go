package internal

import (
	"math"
	"sync"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

type BackoffManager struct {
	config    types.BackoffConfig
	retries   int
	lastDelay float64
	lock      sync.Mutex
}

func NewBackoffManager(config types.BackoffConfig) *BackoffManager {
	return &BackoffManager{
		config:    config,
		retries:   0,
		lastDelay: 0,
	}
}

func (b *BackoffManager) Backoff() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.retries >= b.config.MaxRetries {
		return constant.ErrOverMaxRetries
	}

	currentDelay := b.calculateBackOffDelay()

	currentDelay = math.Max(currentDelay, b.config.InitialDelayMs)
	currentDelay = math.Min(currentDelay, b.config.MaxDelayMs)

	b.retries++
	b.lastDelay = currentDelay

	time.Sleep(time.Duration(currentDelay) * time.Millisecond)

	return nil
}

func (b *BackoffManager) calculateBackOffDelay() float64 {
	switch b.config.Mode {
	case "exponential":
		return b.calculateExponentialBackOffDelay()
	default:
		return b.config.InitialDelayMs
	}
}

func (b *BackoffManager) calculateExponentialBackOffDelay() float64 {
	// if first retry return 0
	if b.retries == 0 {
		return float64(0)
	}

	return math.Min(b.lastDelay+(utils.RandomF64(1)-0.5)*b.lastDelay, b.config.MaxDelayMs)
}
