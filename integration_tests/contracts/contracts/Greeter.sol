// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.4;

contract Greeter {
    string public greeting;

    event ChangeGreeting(address from, string value);

    constructor() {
        greeting = "Hello";
    }

    function setGreeting(string memory _greeting) public {
        greeting = _greeting;
        emit ChangeGreeting(msg.sender, _greeting);
    }

    function greet() public view returns (string memory) {
        return greeting;
    }
}
