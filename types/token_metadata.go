package types

type TokenMetadataResponse struct {
	/* The token's name */
	Name string `json:"name"`

	/* The token's symbol */
	Symbol string `json:"symbol"`

	/* The number of decimals of the token */
	Decimals int `json:"decimals"`

	/* url link to th token's logo */
	Logo string `json:"logo"`
}
