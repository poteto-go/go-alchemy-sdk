---
sidebar_position: 1
---

# Deploy Contract

## Deploy Contract to EthSepolia

This is tutorial for deploy contract to public chain.

:::warning

- It does not work on non-Ethernet compatible networks.

:::

You can get full-code here:

- [link](https://github.com/poteto-go/go-alchemy-sdk/tree/main/examples/deploy-contract)

### 1. Setup

```bash
mkdir deploy-contract
cd deploy-contract
go mod init example.com/example/deploy-contract

mkdir contracts
mkdir build
mkdir abi
```

for example:

```
deploy-contract
├── build
├── abi
├── contracts
└── go.mod
```

### 2. Write Contract by Solidity

```sol
// SPDX-License-Identifier: GPL-3.0

pragma solidity >0.7.0 < 0.9.0;
/**
* @title Storage
* @dev store or retrieve a variable value
*/

contract Storage {

	uint256 value;

	function store(uint256 number) public{
		value = number;
	}

	function retrieve() public view returns (uint256){
		return value;
	}
}
```

### 3. Compile Contract

```bash
solc --combined-json abi,bin contracts/Storage.sol > build/Storage.abi
```

### 4. Generate Abi

```bash
abigen --v2 --combined-json build/Storage.abi --pkg abi --type Storage --out abi/Storage.go
```

You get golang code.

```
deploy-contract
├── abi
│   └── Storage.go
├── build
│   └── Storage.abi
├── contracts
│   └── Storage.sol
├── go.mod
└── go.sum
```

And install dependency.

```bash
go mod tidy
```

### 5. Write Code to Deploy Contract

```bash
go get -u github.com/poteto-go/go-alchemy-sdk
```

You can easily deploy contract by metaData.

```go
package main

import (
	"fmt"

	"github.com/poteto-go/go-alchemy-sdk/gas"
	"github.com/poteto-go/go-alchemy-sdk/types"
	"github.com/poteto-go/go-alchemy-sdk/wallet"
	"github.com/poteto-go/tutorial-deploy-contract/abi"
)

func main() {
	setting := gas.AlchemySetting{
		ApiKey:  "<apiKey>",
		Network: types.EthSepolia,
	}

	alchemy := gas.NewAlchemy(setting)

	w, err := wallet.New("<privateKey>")
	if err != nil {
		panic(err)
	}
	w.Connect(alchemy.GetProvider())

	address, err := w.DeployContract(&abi.StorageMetaData)
	if err != nil {
		panic(err)
	}
	fmt.Println(address)
}
```
