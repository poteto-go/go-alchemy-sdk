package types

type Signature struct {
	V uint8
	R [32]byte
	S [32]byte
}
