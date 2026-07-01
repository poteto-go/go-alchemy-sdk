package ether

import (
	"context"

	"github.com/ethereum/go-ethereum"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/poteto-go/go-alchemy-sdk/constant"
)

// TODO: backoff or
func (ether *Ether) Subscribe(ctx context.Context, channel any, params ...any) (ethereum.Subscription, error) {
	if !ether.isWebSocket() {
		return nil, constant.ErrUnsupportedNotWebsocketProvider
	}

	if err := ether.SetEthClient(); err != nil {
		return nil, err
	}
	defer ether.Close()

	c, ok := ether.Client().(*ethclient.Client)
	if !ok {
		return nil, constant.ErrUnSupportSimulatedMethod
	}

	return c.Client().EthSubscribe(ctx, channel, params...)
}

func (ether *Ether) SubscribeNewHead(ctx context.Context, headerChan chan<- *gethTypes.Header) (ethereum.Subscription, error) {
	if !ether.isWebSocket() {
		return nil, constant.ErrUnsupportedNotWebsocketProvider
	}

	if err := ether.SetEthClient(); err != nil {
		return nil, err
	}
	defer ether.Close()

	c, ok := ether.Client().(*ethclient.Client)
	if !ok {
		return nil, constant.ErrUnSupportSimulatedMethod
	}

	return c.SubscribeNewHead(ctx, headerChan)
}

func (ether *Ether) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, logChan chan<- gethTypes.Log) (ethereum.Subscription, error) {
	if !ether.isWebSocket() {
		return nil, constant.ErrUnsupportedNotWebsocketProvider
	}

	if err := ether.SetEthClient(); err != nil {
		return nil, err
	}
	defer ether.Close()

	c, ok := ether.Client().(*ethclient.Client)
	if !ok {
		return nil, constant.ErrUnSupportSimulatedMethod
	}

	return c.SubscribeFilterLogs(ctx, query, logChan)
}

func (ether *Ether) SubscribeTxReceipts(ctx context.Context, q *ethereum.TransactionReceiptsQuery, receiptsChan chan<- []*gethTypes.Receipt) (ethereum.Subscription, error) {
	if !ether.isWebSocket() {
		return nil, constant.ErrUnsupportedNotWebsocketProvider
	}

	if err := ether.SetEthClient(); err != nil {
		return nil, err
	}
	defer ether.Close()

	c, ok := ether.Client().(*ethclient.Client)
	if !ok {
		return nil, constant.ErrUnSupportSimulatedMethod
	}

	return c.SubscribeTransactionReceipts(ctx, q, receiptsChan)
}
