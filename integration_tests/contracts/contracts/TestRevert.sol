// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.4;

contract TestRevert {
    uint256 state;
    constructor() {
        state = 0;
    }
    function transfer(uint256 value) public payable {
        uint256 minimal = 5 * 10 ** 18;
        state = value;
        if (state < minimal) {
            revert("Not enough tokens to transfer");
        }
    }
    function query() public view returns (uint256) {
        return state;
    }
}
