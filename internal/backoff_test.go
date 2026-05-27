package internal

import (
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

var backoffConfigTest = types.BackoffConfig{
	Mode:           "exponential",
	MaxRetries:     3,
	InitialDelayMs: 10,
	MaxDelayMs:     30,
}

func TestNewBackoffManager(t *testing.T) {
	// Arrange
	config := types.BackoffConfig{
		Mode:           "exponential",
		MaxRetries:     1,
		InitialDelayMs: 10,
		MaxDelayMs:     10,
	}

	// Act
	manager := NewBackoffManager(config)

	// Assert
	assert.NotNil(t, manager)
	assert.Equal(t, config.Mode, manager.config.Mode)
	assert.Equal(t, config.MaxRetries, manager.config.MaxRetries)
	assert.Equal(t, config.InitialDelayMs, manager.config.InitialDelayMs)
	assert.Equal(t, config.MaxDelayMs, manager.config.MaxDelayMs)
	assert.Equal(t, 0, manager.retries)
	assert.Equal(t, float64(0), manager.lastDelay)
}

func TestBackoffManager_Backoff(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("exponential backoff", func(t *testing.T) {
			// Arrange
			manager := NewBackoffManager(backoffConfigTest)

			for i := 1; i <= manager.config.MaxRetries; i++ {
				// Act
				err := manager.Backoff()

				// Assert
				assert.Nil(t, err)
				assert.Equal(t, i, manager.retries)
				assert.Condition(
					t,
					func() (success bool) {
						return manager.lastDelay >= manager.config.InitialDelayMs && manager.lastDelay <= manager.config.MaxDelayMs
					},
				)
			}
		})
	})

	t.Run("error case:", func(t *testing.T) {
		t.Run("max over retries -> constant.ErrOverMaxRetries", func(t *testing.T) {
			// Arrange
			manager := NewBackoffManager(backoffConfigTest)
			manager.config.MaxRetries = 0

			// Act
			err := manager.Backoff()

			// Assert
			assert.ErrorIs(t, err, constant.ErrOverMaxRetries)
		})

	})
}

func TestBackoffManager_calculateExponentialBackOffDelay(t *testing.T) {
	t.Run("first retry returns 0", func(t *testing.T) {
		// Arrange
		manager := NewBackoffManager(backoffConfigTest)

		// Act
		delay := manager.calculateExponentialBackOffDelay()

		// Assert
		assert.Equal(t, float64(0), delay)
	})

	t.Run("grows exponentially (doubles) without jitter", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock RandomF64 to return 0.5 so jitter term is 0
		patches.ApplyFunc(utils.RandomF64, func(max float64) float64 {
			return 0.5
		})

		// Arrange
		config := types.BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     10,
			InitialDelayMs: 10,
			MaxDelayMs:     100000,
		}
		manager := NewBackoffManager(config)
		manager.retries = 1
		manager.lastDelay = 10

		// Act
		delay := manager.calculateExponentialBackOffDelay()

		// Assert: base is doubled (10 * 2 = 20), jitter is 0 -> 20
		assert.Equal(t, float64(20), delay)
	})

	t.Run("jitter at lower bound returns base * 0.5", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock RandomF64 to return 0 -> jitter is -0.5 * base
		patches.ApplyFunc(utils.RandomF64, func(max float64) float64 {
			return 0
		})

		// Arrange
		config := types.BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     10,
			InitialDelayMs: 10,
			MaxDelayMs:     100000,
		}
		manager := NewBackoffManager(config)
		manager.retries = 1
		manager.lastDelay = 10

		// Act
		delay := manager.calculateExponentialBackOffDelay()

		// Assert: base = 20, jitter = -10, result = 10
		assert.Equal(t, float64(10), delay)
	})

	t.Run("capped at MaxDelayMs", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock RandomF64 returning 0.5 so jitter is 0
		patches.ApplyFunc(utils.RandomF64, func(max float64) float64 {
			return 0.5
		})

		// Arrange
		config := types.BackoffConfig{
			Mode:           "exponential",
			MaxRetries:     10,
			InitialDelayMs: 10,
			MaxDelayMs:     30,
		}
		manager := NewBackoffManager(config)
		manager.retries = 1
		manager.lastDelay = 20

		// Act
		delay := manager.calculateExponentialBackOffDelay()

		// Assert: base = 40, capped at 30
		assert.Equal(t, float64(30), delay)
	})
}
