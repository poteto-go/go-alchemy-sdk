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
