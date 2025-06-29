## `Ether.GetCode`

### Summary

Returns the contract code of the provided address at the block.
If there is no contract deployed, the result is 0x.

### Behavior

- [x] normal case

  - [x] if exist, return code hex string
    - [x] temp impl
    - [x] actual impl
    - [x] check call with method eth_getCode
  - [x] if not exist, return 0x
    - What does ether return if not contract exists?
      - EX:)
        - {"jsonrpc":"2.0","id":1,"result":"0x"}
        - then do nothing

- [x] error case
  - [x] if invalid BlockTag, throw error
  - [x] if connected error, throw error

## `Core.GetCode`

### Summary

Returns the contract code of the provided address at the block.
If there is no contract deployed, the result is 0x.

### Behavior

- [x] normal case
  - [x] call `Ether.GetCode`
    - return val check
- [x] error case
  - [x] if ether has error, return err

## `Core.IsContractAddress`

### Summary

Checks if the provided address is a smart contract.

### Behavior

- [x] common check
  - blockTag == "latest" on `Core.GetCode`
- [x] if has valid code hexString (!="0x") return true
- [x] if invalid code hexString (=="0x") return false
- [x] if error occur in `Core.GetCode`, return false
