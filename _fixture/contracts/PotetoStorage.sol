// SPDX-License-Identifier: GPL-3.0

pragma solidity >0.7.0 < 0.9.0;
/**
* @title Storage
* @dev store or retrieve a variable value
*/

contract PotetoStorage {

	uint256 value;

	event Stored(address indexed sender, uint256 value);

	function store(uint256 number) public{
		value = number;
		emit Stored(msg.sender, number);
	}

	function retrieve() public view returns (uint256){
		return value;
	}
}
