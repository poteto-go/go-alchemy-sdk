package internal

import (
	"encoding/hex"
	"strings"
)

func DecodeHex(hexStr string) ([]byte, error) {
	trimedHex := strings.TrimPrefix(hexStr, "0x")
	return hex.DecodeString(trimedHex)
}
