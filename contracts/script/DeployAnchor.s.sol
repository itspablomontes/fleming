// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Script, console} from "forge-std/Script.sol";
import {FlemingAnchor} from "../src/FlemingAnchor.sol";

/// @title DeployAnchor
/// @notice Deployment script for FlemingAnchor contract
/// @dev Run with: forge script script/DeployAnchor.s.sol --rpc-url $BASE_SEPOLIA_RPC_URL --broadcast --verify
contract DeployAnchor is Script {
    function setUp() public {}

    function run() public {
        uint256 deployerPrivateKey = vm.envUint("ANCHOR_PRIVATE_KEY");

        console.log("Deploying FlemingAnchor...");
        console.log("Deployer:", vm.addr(deployerPrivateKey));

        vm.startBroadcast(deployerPrivateKey);

        FlemingAnchor anchor = new FlemingAnchor(vm.addr(deployerPrivateKey));

        vm.stopBroadcast();

        console.log("FlemingAnchor deployed to:", address(anchor));
        console.log("");
        console.log("Add to your .env:");
        console.log("ANCHOR_CONTRACT_ADDRESS=%s", address(anchor));
    }
}
