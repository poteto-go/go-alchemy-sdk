package types

type IAlchemyProvider interface {
	// get  the number of the most recent block.
	GetBlockNumber() (int, error)

	// send raw transaction
	// & return result of response
	Send(method string, params ...string) (string, error)
}
