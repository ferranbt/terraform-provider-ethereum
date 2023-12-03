// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.4;

contract Inputs {
    address public addr;
    Config config;
    uint256 number;
    bytes32 bytesI;

    struct Config {
        uint256 number;
    }

    constructor(
        address _addr,
        uint256 _number,
        bytes32 _bytesI,
        Config memory _config
    ) {
        addr = _addr;
        bytesI = _bytesI;
        config = _config;
        number = _number;
    }

    function applyFunc(
        address _addr,
        uint256 _number,
        bytes32 _bytesI,
        Config memory _config
    ) public payable {
        addr = _addr;
        bytesI = _bytesI;
        config = _config;
        number = _number;
    }
}

contract WithEvents {
    event One(
        uint64 num
    );

    function applyFunc() public {
        emit One(1);
    }
}

contract Call {
    function multipleOutput() public pure returns (uint64, uint64) {
        return (1, 2);
    }

    function singleOutput() public pure returns (uint64) {
        return 1;
    }
}
