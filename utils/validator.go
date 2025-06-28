package utils

import "github.com/poteto-go/go-alchemy-sdk/core"

// https://docs.ethers.org/v5/api/providers/types/
// just support latest | earliest | pending for now
func ValidateBlockTag(blockTag string) error {
	if blockTag == "latest" {
		return nil
	}

	if blockTag == "earliest" {
		return nil
	}

	if blockTag == "pending" {
		return nil
	}

	return core.ErrInvalidBlockTag
}
