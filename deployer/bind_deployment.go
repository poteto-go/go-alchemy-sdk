package deployer

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
)

func BindDeploymentMetadata(base *bind.MetaData, args ...any) error {
	if base == nil {
		return errors.New("unexpected nil base or initial supplier")
	}

	if len(args) == 0 {
		return nil
	}

	parsedABI, err := abi.JSON(strings.NewReader(base.ABI))
	if err != nil {
		return err
	}

	// "" = constructor
	input, err := parsedABI.Pack("", args...)
	if err != nil {
		return err
	}

	binStr := strings.TrimPrefix(base.Bin, "0x")
	binBytes, err := hex.DecodeString(binStr)
	if err != nil {
		return err
	}

	fullBin := append(binBytes, input...)
	base.Bin = "0x" + hex.EncodeToString(fullBin)
	return nil
}
