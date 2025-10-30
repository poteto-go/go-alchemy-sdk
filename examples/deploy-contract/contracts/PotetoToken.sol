// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "./ERC20.sol";

contract PotetoToken is ERC20 {
    constructor() ERC20("PotetoToken", "PTT", 1)
    {
        _mint(msg.sender, 100 * 10 ** uint256(decimals));
    }
}