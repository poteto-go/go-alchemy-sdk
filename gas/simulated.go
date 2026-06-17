package gas

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

// simulated alchemy connect to simulated backend
type SimulatedAlchemy struct {
	Core       namespace.ICore
	Transact   namespace.ITransact
	Nft        namespace.INft
	ERC1155    namespace.IErc1155
	ERC20      namespace.IERC20
	StableCoin namespace.IStableCoin
	Debug      namespace.IDebug
	provider   types.IAlchemyProvider
}

/*
With Geth's simulatedBackend, you can connect to a simulated blockchain node without launching a chain.
This enables the execution of high-speed tests.

! simulated backend doesn't support un-geth supported method

	sim := simulated.NewBackend(
		types.GenesisAlloc{ addr: {Balance: big.NewInt(...)} },
		options...,
	)
	defer sim.Close()
	alchemy := gas.NewSimulatedAlchemy(sim)
*/
func NewSimulatedAlchemy(backend *simulated.Backend) (SimulatedAlchemy, error) {
	if backend == nil {
		return SimulatedAlchemy{}, errors.New("no connected simulated backend")
	}

	alchemyProvider := NewAlchemyProvider(AlchemyConfig{})
	eth := ether.NewSimulatedApi(backend)
	alchemyProvider.SetEth(eth)
	coreNamespace := namespace.NewCore(eth)
	transactNamespace := namespace.NewTransactNamespace(eth)
	nftNamespace := namespace.NewNftNamespace(eth)
	erc1155Namespace := namespace.NewErc1155Namespace(eth)
	erc20Namespace := namespace.NewERC20Namespace(eth)
	stableCoinNamespace := namespace.NewStableCoinNamespace(eth)
	debugNamespace := namespace.NewSimulatedDebugNamespace(eth)
	return SimulatedAlchemy{
		Core:       coreNamespace,
		Transact:   transactNamespace,
		Nft:        nftNamespace,
		ERC1155:    erc1155Namespace,
		ERC20:      erc20Namespace,
		StableCoin: stableCoinNamespace,
		Debug:      debugNamespace,
		provider:   alchemyProvider,
	}, nil
}

func (gas *SimulatedAlchemy) GetProvider() types.IAlchemyProvider {
	return gas.provider
}
