# `GetTokenBalances`

## `Ether`

send with alchemy_getTplemBalances with params

- [x] normal case
  - [x] call with alchemy_getTokenBalances and params & return result
  - [x] if not params provided, call with alchemy_getTokenBalances & return result
- [x] error case
  - [x] if error occur in send, return internal error

## `Core.GetTokenBalances(address)`

- [x] call `ether.GetTokenBalances` with just address & return result
- [x] call `ether.GetTokenBalances` with just address & return internal error

## FIX

response type is TokenBalanceResponse

- [x] normal
- [x] if response includes error
- [x] error on unmarshal

## `Core.GetTokenBalances(address, contracts)`

- [ ] call `ether.GetTokenBalances` with address & contracts & return filtered result
- [ ] call `ether.GetTokenBalances` with address & contracts & return internal error
