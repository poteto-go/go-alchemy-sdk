![](https://img.shields.io/badge/go-geth-lightblue)

Returns the number of p2p peers as reported by the net_peerCount method.

```go
func PeerCount() (count uint64, err error)
```

```go
func main() {
	...
	alchemy := gas.NewAlchemy(setting)
	res, _ := alchemy.Core.PeerCount()
}
```
