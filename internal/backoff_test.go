package internal

import (
	"testing"

	"github.com/poteto-go/go-alchemy-sdk/core"
	"github.com/stretchr/testify/assert"
)

var backoffConfigTest = BackoffConfig{
	Mode:           "exponential",
	MaxRetries:     3,
	InitialDelayMs: 10,
	MaxDelayMs:     30,
}

func TestNewBackoffManager(t *testing.T) {
	// Arrange
	config := BackoffConfig{
		Mode:           "exponential",
		MaxRetries:     1,
		InitialDelayMs: 10,
		MaxDelayMs:     10,
	}

	// Act
	manager := NewBackoffManager(config).(*BackoffManager)

	// Assert
	assert.NotNil(t, manager)
	assert.Equal(t, config.Mode, manager.config.Mode)
	assert.Equal(t, config.MaxRetries, manager.config.MaxRetries)
	assert.Equal(t, config.InitialDelayMs, manager.config.InitialDelayMs)
	assert.Equal(t, config.MaxDelayMs, manager.config.MaxDelayMs)
	assert.Equal(t, 0, manager.retries)
	assert.Equal(t, float64(0), manager.lastDelay)
}

func TestBackoffManager_Reset(t *testing.T) {
	// Arrange
	manager := NewBackoffManager(backoffConfigTest).(*BackoffManager)

	// Act
	manager.Reset()

	// Assert
	assert.Equal(t, 0, manager.retries)
	assert.Equal(t, float64(0), manager.lastDelay)
}

func TestBackoffManager_Backoff(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("exponential backoff", func(t *testing.T) {
			// Arrange
			manager := NewBackoffManager(backoffConfigTest).(*BackoffManager)

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
		t.Run("max over retries -> core.ErrOverMaxRetries", func(t *testing.T) {
			// Arrange
			manager := NewBackoffManager(backoffConfigTest).(*BackoffManager)
			manager.config.MaxRetries = 0

			// Act
			err := manager.Backoff()

			// Assert
			assert.ErrorIs(t, err, core.ErrOverMaxRetries)
		})

	})
}
