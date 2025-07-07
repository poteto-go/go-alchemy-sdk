# `Ether.EstimateGas`

Returns an estimate of the amount of gas that would be required to submit transaction to the network.

An estimate may not be accurate since there could be another transaction on the network that was not accounted for, but after being mined affects the relevant state.
This is an alias for {@link TransactNamespace.estimateGas}.

- [x] normal case
  - [x] call eth_estimateGas & estimate gas and transform to big.Int
- [x] error case
  - [x] if error occur in marshal json, return `core.ErrFailedToMarshalParameter`
  - [x] if error occur in Send, return internal error
  - [x] if error occur in FromBigHex, return internal error

# `Core.EstimateGas`
