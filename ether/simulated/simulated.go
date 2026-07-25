/*
Package simulated wires go-ethereum's in-process simulated backend into the SDK.

It is kept out of the core packages on purpose: ethclient/simulated links a full
geth node (core/vm, eth, node, pebble, the p2p stack, ...) which roughly triples
the binary of an SDK user who never touches a simulated chain. Import this
package only when you actually want a simulated chain.
*/
package simulated

import (
	"errors"

	gethSimulated "github.com/ethereum/go-ethereum/ethclient/simulated"

	"github.com/poteto-go/go-alchemy-sdk/ether"
	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/namespace"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

/*
NewSimulatedApi returns an EtherApi backed by geth's simulated backend.

The backend is owned by the caller: the SDK never closes or re-creates it.
*/
func NewSimulatedApi(backend *gethSimulated.Backend) types.EtherApi {
	return ether.NewSimulatedEtherApi(backend, backend.Client())
}

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
	alchemy := simulated.NewSimulatedAlchemy(sim)
*/
func NewSimulatedAlchemy(backend *gethSimulated.Backend) (SimulatedAlchemy, error) {
	if backend == nil {
		return SimulatedAlchemy{}, errors.New("no connected simulated backend")
	}

	alchemyProvider := gas.NewAlchemyProvider(gas.AlchemyConfig{})
	eth := NewSimulatedApi(backend)
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
