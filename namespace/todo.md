# `Core.GetStorageAt`

Return the value of the provided position at the provided address, at the provided block in `Bytes32` format. For inspecting solidity code.

- [x] 正常系
  - [x] call eth_getStorageAt & return provided block
- [x] エラーケース
  - [x] 内部エラーが発生したとき、そのエラーを返す.
  - [x] 無効な blockTag でエラーを返す
