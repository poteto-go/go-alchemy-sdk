package namespace

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/poteto-go/go-alchemy-sdk/types"
)

type IWS interface {
	// Subscribe subscribes to the specified event and returns a subscription ID.
	Subscribe(ctx context.Context, channel any, params ...any) (ethereum.Subscription, error)

	// subscribe new block
	SubscribeNewHead(ctx context.Context, headerChan chan<- *gethTypes.Header) (ethereum.Subscription, error)

	// subscribe event logs by filter query
	SubscribeLogs(ctx context.Context, query ethereum.FilterQuery, logChan chan<- gethTypes.Log) (ethereum.Subscription, error)

	// subscribe event logs of a specific contract address
	SubscribeContractLogs(ctx context.Context, contractAddress common.Address, logChan chan<- gethTypes.Log) (ethereum.Subscription, error)

	// subscribe transaction receipts
	SubscribeTxReceipts(ctx context.Context, query *ethereum.TransactionReceiptsQuery, receiptsChan chan<- []*gethTypes.Receipt) (ethereum.Subscription, error)
}

func NewWSNamespace(ether types.EtherApi) IWS {
	return &WS{
		ether: ether,
	}
}

type WS struct {
	ether types.EtherApi
}

func (w *WS) Subscribe(ctx context.Context, channel any, params ...any) (ethereum.Subscription, error) {
	return w.ether.Subscribe(ctx, channel, params...)
}

func (w *WS) SubscribeNewHead(ctx context.Context, headerChan chan<- *gethTypes.Header) (ethereum.Subscription, error) {
	return w.ether.SubscribeNewHead(ctx, headerChan)
}

func (w *WS) SubscribeLogs(ctx context.Context, query ethereum.FilterQuery, logChan chan<- gethTypes.Log) (ethereum.Subscription, error) {
	return w.ether.SubscribeFilterLogs(ctx, query, logChan)
}

func (w *WS) SubscribeContractLogs(ctx context.Context, contractAddress common.Address, logChan chan<- gethTypes.Log) (ethereum.Subscription, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	return w.ether.SubscribeFilterLogs(ctx, query, logChan)
}

func (w *WS) SubscribeTxReceipts(ctx context.Context, query *ethereum.TransactionReceiptsQuery, receiptsChan chan<- []*gethTypes.Receipt) (ethereum.Subscription, error) {
	return w.ether.SubscribeTxReceipts(ctx, query, receiptsChan)
}
