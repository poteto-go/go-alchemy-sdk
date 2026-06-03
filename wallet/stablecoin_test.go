package wallet

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestWallet_StableCoin(t *testing.T) {
	t.Run("returns WalletStableCoin", func(t *testing.T) {
		w := createConnectedWallet()

		sc := w.StableCoin()

		assert.NotNil(t, sc)
	})
}

func TestWallet_StableCoin_MintNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	toAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can mint stablecoin", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().MintNoWait(contractAddress, toAddress, big.NewInt(100), nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on mint", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().MintNoWait(contractAddress, toAddress, big.NewInt(100), nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().MintNoWait(contractAddress, toAddress, big.NewInt(100), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Mint(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	toAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can mint stablecoin and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().Mint(context.Background(), contractAddress, toAddress, big.NewInt(100), nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Mint(context.Background(), contractAddress, toAddress, big.NewInt(100), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_BurnNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	expectedHash := common.HexToHash("0x123")

	t.Run("can burn stablecoin", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().BurnNoWait(contractAddress, big.NewInt(50), nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on burn", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().BurnNoWait(contractAddress, big.NewInt(50), nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().BurnNoWait(contractAddress, big.NewInt(50), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Burn(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	expectedHash := common.HexToHash("0x123")

	t.Run("can burn stablecoin and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().Burn(context.Background(), contractAddress, big.NewInt(50), nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Burn(context.Background(), contractAddress, big.NewInt(50), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_BlacklistNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	targetAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can blacklist address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().BlacklistNoWait(contractAddress, targetAddress, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on blacklist", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().BlacklistNoWait(contractAddress, targetAddress, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().BlacklistNoWait(contractAddress, targetAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Blacklist(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	targetAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can blacklist address and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().Blacklist(context.Background(), contractAddress, targetAddress, nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Blacklist(context.Background(), contractAddress, targetAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_UnBlacklistNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	targetAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can unBlacklist address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().UnBlacklistNoWait(contractAddress, targetAddress, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on unBlacklist", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().UnBlacklistNoWait(contractAddress, targetAddress, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().UnBlacklistNoWait(contractAddress, targetAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_UnBlacklist(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	targetAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can unBlacklist address and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().UnBlacklist(context.Background(), contractAddress, targetAddress, nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().UnBlacklist(context.Background(), contractAddress, targetAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_IsBlacklisted(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	targetAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"

	t.Run("returns true when address is blacklisted", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)
		expected[31] = 1

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().IsBlacklisted(contractAddress, targetAddress)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("returns false when address is not blacklisted", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().IsBlacklisted(contractAddress, targetAddress)

		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().IsBlacklisted(contractAddress, targetAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Currency(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"

	t.Run("returns currency string", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := "USD"

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return utils.EncodeABIString(expected), nil
			},
		)

		result, err := w.StableCoin().Currency(contractAddress)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Currency(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Version(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"

	t.Run("returns version string", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := "1"

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return utils.EncodeABIString(expected), nil
			},
		)

		result, err := w.StableCoin().Version(contractAddress)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Version(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_PauseNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	expectedHash := common.HexToHash("0x123")

	t.Run("can pause contract", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().PauseNoWait(contractAddress, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on pause", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().PauseNoWait(contractAddress, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().PauseNoWait(contractAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Pause(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	expectedHash := common.HexToHash("0x123")

	t.Run("can pause contract and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().Pause(context.Background(), contractAddress, nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Pause(context.Background(), contractAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_UnpauseNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	expectedHash := common.HexToHash("0x123")

	t.Run("can unpause contract", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().UnpauseNoWait(contractAddress, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on unpause", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().UnpauseNoWait(contractAddress, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().UnpauseNoWait(contractAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Unpause(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	expectedHash := common.HexToHash("0x123")

	t.Run("can unpause contract and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().Unpause(context.Background(), contractAddress, nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Unpause(context.Background(), contractAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Paused(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"

	t.Run("returns true when contract is paused", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)
		expected[31] = 1

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().Paused(contractAddress)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("returns false when contract is not paused", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().Paused(contractAddress)

		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Paused(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_MasterMinter(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	masterMinterAddress := common.HexToAddress("0xE25583099BA105D9ec0A67f5Ae86D90e50036425")

	t.Run("returns master minter address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)
		copy(expected[12:], masterMinterAddress.Bytes())

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().MasterMinter(contractAddress)

		assert.Nil(t, err)
		assert.Equal(t, masterMinterAddress, result)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().MasterMinter(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Pauser(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	pauserAddress := common.HexToAddress("0xE25583099BA105D9ec0A67f5Ae86D90e50036425")

	t.Run("returns pauser address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)
		copy(expected[12:], pauserAddress.Bytes())

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().Pauser(contractAddress)

		assert.Nil(t, err)
		assert.Equal(t, pauserAddress, result)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Pauser(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Blacklister(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	blacklisterAddress := common.HexToAddress("0xE25583099BA105D9ec0A67f5Ae86D90e50036425")

	t.Run("returns blacklister address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)
		copy(expected[12:], blacklisterAddress.Bytes())

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().Blacklister(contractAddress)

		assert.Nil(t, err)
		assert.Equal(t, blacklisterAddress, result)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Blacklister(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_Owner(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	ownerAddress := common.HexToAddress("0xE25583099BA105D9ec0A67f5Ae86D90e50036425")

	t.Run("returns owner address", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)
		copy(expected[12:], ownerAddress.Bytes())

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().Owner(contractAddress)

		assert.Nil(t, err)
		assert.Equal(t, ownerAddress, result)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().Owner(contractAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_TransferOwnershipNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	newOwnerAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can transfer ownership", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().TransferOwnershipNoWait(contractAddress, newOwnerAddress, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on transfer ownership", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().TransferOwnershipNoWait(contractAddress, newOwnerAddress, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().TransferOwnershipNoWait(contractAddress, newOwnerAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_ConfigureMinterNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	minterAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can configure minter", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().ConfigureMinterNoWait(contractAddress, minterAddress, big.NewInt(1_000_000), nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on configure minter", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().ConfigureMinterNoWait(contractAddress, minterAddress, big.NewInt(1_000_000), nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().ConfigureMinterNoWait(contractAddress, minterAddress, big.NewInt(1_000_000), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_ConfigureMinter(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	minterAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can configure minter and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().ConfigureMinter(context.Background(), contractAddress, minterAddress, big.NewInt(1_000_000), nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().ConfigureMinter(context.Background(), contractAddress, minterAddress, big.NewInt(1_000_000), nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_RemoveMinterNoWait(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	minterAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can remove minter", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)

		hash, err := w.StableCoin().RemoveMinterNoWait(contractAddress, minterAddress, nil)

		assert.Nil(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("handle error on remove minter", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return common.Hash{}, errors.New("error")
			},
		)

		_, err := w.StableCoin().RemoveMinterNoWait(contractAddress, minterAddress, nil)

		assert.Error(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().RemoveMinterNoWait(contractAddress, minterAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_RemoveMinter(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	minterAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can remove minter and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().RemoveMinter(context.Background(), contractAddress, minterAddress, nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().RemoveMinter(context.Background(), contractAddress, minterAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_IsMinter(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	minterAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"

	t.Run("returns true when address is a minter", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)
		expected[31] = 1

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().IsMinter(contractAddress, minterAddress)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("returns false when address is not a minter", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().IsMinter(contractAddress, minterAddress)

		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().IsMinter(contractAddress, minterAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_MinterAllowance(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	minterAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"

	t.Run("returns minter allowance", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()
		expected := make([]byte, 32)
		expected[31] = 100

		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"CallContract",
			func(_ *ether.Ether, _ ethereum.CallMsg, _ string) ([]byte, error) {
				return expected, nil
			},
		)

		result, err := w.StableCoin().MinterAllowance(contractAddress, minterAddress)

		assert.Nil(t, err)
		assert.Equal(t, int64(100), result.Int64())
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().MinterAllowance(contractAddress, minterAddress)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}

func TestWallet_StableCoin_TransferOwnership(t *testing.T) {
	contractAddress := "0x1234567890123456789012345678901234567890"
	newOwnerAddress := "0xE25583099BA105D9ec0A67f5Ae86D90e50036425"
	expectedHash := common.HexToHash("0x123")

	t.Run("can transfer ownership and wait", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		w := createConnectedWallet()

		patches.ApplyMethod(
			reflect.TypeOf(w),
			"SendTransaction",
			func(_ *wallet, _ types.TransactionRequest) (common.Hash, error) {
				return expectedHash, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(w.provider.Eth()),
			"WaitMined",
			func(_ *ether.Ether, _ context.Context, _ common.Hash) (*gethTypes.Receipt, error) {
				return &gethTypes.Receipt{}, nil
			},
		)

		_, err := w.StableCoin().TransferOwnership(context.Background(), contractAddress, newOwnerAddress, nil)

		assert.Nil(t, err)
	})

	t.Run("error w/o connect wallet", func(t *testing.T) {
		w, _ := New(testPrivHex)

		_, err := w.StableCoin().TransferOwnership(context.Background(), contractAddress, newOwnerAddress, nil)

		assert.ErrorIs(t, err, constant.ErrWalletIsNotConnected)
	})
}
