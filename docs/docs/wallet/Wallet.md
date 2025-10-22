---
sidebar_position: 1
---

Wallet class inherits Signer and can sign transactions and messages using.

```go
func New(privateKeyStr string) (w Wallet, err error)
```

You can create Wallet w/ your privateKey string.

```go
func main() {
  w, _ := wallet.New("privateKey")
}
```
