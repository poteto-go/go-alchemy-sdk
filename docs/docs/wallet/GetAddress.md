---
sidebar_position: 2
---

get address of wallet

```go
func GetAddress() (address string)
```

```go
func main() {
	w, _ := wallet.New("<privateKey>")
	address := w.GetAddress()
	fmt.Println(address)
}
```
