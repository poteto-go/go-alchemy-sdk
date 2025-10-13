package internal

import (
	"math"
	"sync"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/utils"
)

type BackoffConfig struct {
	Mode           string  `yaml:"mode"`
	MaxRetries     int     `yaml:"max_retries"`
	InitialDelayMs float64 `yaml:"initial_delay_ms"`
	MaxDelayMs     float64 `yaml:"max_delay_ms"`
}

var DefaultBackoffConfig = BackoffConfig{
	Mode:           "exponential",
	MaxRetries:     1,
	InitialDelayMs: 1000,
	MaxDelayMs:     30000,
}

type BackoffManager struct {
	config    BackoffConfig
	retries   int
	lastDelay float64
	lock      sync.Mutex
}

func NewBackoffManager(config BackoffConfig) *BackoffManager {
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

	currentDelay := b.calcBackOffDelay()

	currentDelay = math.Max(currentDelay, b.config.InitialDelayMs)
	currentDelay = math.Min(currentDelay, b.config.MaxDelayMs)

	b.retries++
	b.lastDelay = currentDelay

	time.Sleep(time.Duration(currentDelay) * time.Millisecond)

	return nil
}

func (b *BackoffManager) calcBackOffDelay() float64 {
	switch b.config.Mode {
	case "exponential":
		return b.calcExponentialBackOffDelay()
	default:
		return b.config.InitialDelayMs
	}
}

func (b *BackoffManager) calcExponentialBackOffDelay() float64 {
	// if first retry return 0
	if b.retries == 0 {
		return float64(0)
	}

	return math.Min(b.lastDelay+(utils.RandomF64(1)-0.5)*b.lastDelay, b.config.MaxDelayMs)
}
