package wallet

type WalletStableCoin interface {
	WalletERC20
}

type walletStableCoin struct {
	walletERC20
}
