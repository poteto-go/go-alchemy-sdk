package internal

import (
	"math"
	"sync"
	"time"

	"github.com/poteto-go/go-alchemy-sdk/core"
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
	MaxRetries:     5,
	InitialDelayMs: 1000,
	MaxDelayMs:     30000,
}

type IBackoffManager interface {
	Reset()
	Backoff() error
}

type BackoffManager struct {
	config    BackoffConfig
	retries   int
	lastDelay float64
	lock      sync.Mutex
}

func NewBackoffManager(config BackoffConfig) IBackoffManager {
	return &BackoffManager{
		config:    config,
		retries:   0,
		lastDelay: 0,
	}
}

func (b *BackoffManager) Reset() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.retries = 0
	b.lastDelay = 0
}

func (b *BackoffManager) Backoff() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.retries >= b.config.MaxRetries {
		return core.ErrOverMaxRetries
	}

	var currentDelay float64 = 0
	if b.config.Mode == "exponential" {
		currentDelay = b.exponentialBackoff()
	}

	currentDelay = math.Max(currentDelay, b.config.InitialDelayMs)
	currentDelay = math.Min(currentDelay, b.config.MaxDelayMs)

	b.retries++
	b.lastDelay = currentDelay

	time.Sleep(time.Duration(currentDelay) * time.Millisecond)

	return nil
}

func (b *BackoffManager) exponentialBackoff() float64 {
	if b.retries == 0 {
		return float64(0)
	}
	return math.Min(b.lastDelay+(utils.RandomF64(0, 1)-0.5), b.config.MaxDelayMs)
}
