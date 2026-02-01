// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

import {Script, console} from "forge-std/Script.sol";
import {FlemingAnchor} from "../src/FlemingAnchor.sol";
import {VCRegistry} from "../src/VCRegistry.sol";

/// @title DeployFleming
/// @notice Unified deployment script for all Fleming contracts
/// @dev Supports Base Sepolia (84532) and Base Mainnet (8453)
/// @custom:version 1.0.0
contract DeployFleming is Script {
    // ─── Custom Errors ─────────────────────────────────────────────────────────

    error InvalidChain(uint256 chainId);
    error DeploymentFailed(string reason);
    error MissingEnvironmentVariable(string varName);
    error ZeroAddress(string name);

    // ─── Events ────────────────────────────────────────────────────────────────

    event DeploymentComplete(address indexed anchor, address indexed vcRegistry, uint256 chainId);

    // ─── Structs ───────────────────────────────────────────────────────────────

    struct DeploymentConfig {
        address anchorer;
        address issuer;
    }

    struct DeploymentResult {
        FlemingAnchor anchor;
        VCRegistry vcRegistry;
    }

    // ─── Constants ─────────────────────────────────────────────────────────────

    uint256 constant BASE_SEPOLIA = 84532;
    uint256 constant BASE_MAINNET = 8453;
    uint256 constant ANVIL_LOCAL = 31337;

    // ─── State ─────────────────────────────────────────────────────────────────

    address public deployer;
    uint256 public chainId;

    // ─── External Functions ────────────────────────────────────────────────────

    function run() external returns (DeploymentResult memory result) {
        // Setup
        _setup();

        // Load configuration
        DeploymentConfig memory config = _loadConfig();

        // Deploy contracts
        result = _deploy(config);

        // Verify deployment
        _verifyDeployment(result);

        // Log results
        _logDeployment(result);

        emit DeploymentComplete(address(result.anchor), address(result.vcRegistry), chainId);

        // Write to JSON file (Local Dev Automation)
        _writeDeploymentJSON(result);

        return result;
    }

    // ─── Internal Functions ────────────────────────────────────────────────────

    /// @notice Setup deployment environment
    function _setup() internal {
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");
        deployer = vm.addr(deployerKey);
        chainId = block.chainid;

        console.log("========================================");
        console.log("Fleming Protocol Deployment");
        console.log("========================================");
        console.log("Deployer:", deployer);
        console.log("Chain ID:", chainId);
        console.log("Network:", _getNetworkName(chainId));
        console.log("");

        // Validate chain
        if (chainId != BASE_SEPOLIA && chainId != BASE_MAINNET && chainId != ANVIL_LOCAL) {
            revert InvalidChain(chainId);
        }

        vm.startBroadcast(deployerKey);
    }

    /// @notice Load deployment configuration from environment
    function _loadConfig() internal view returns (DeploymentConfig memory config) {
        // Try to load from env, fallback to deployer address
        try vm.envAddress("ANCHORER_ADDRESS") returns (address addr) {
            config.anchorer = addr;
        } catch {
            config.anchorer = deployer;
            console.log("ANCHORER_ADDRESS not set, using deployer:", config.anchorer);
        }
        try vm.envAddress("ISSUER_ADDRESS") returns (address addr) {
            config.issuer = addr;
        } catch {
            config.issuer = deployer;
            console.log("ISSUER_ADDRESS not set, using deployer:", config.issuer);
        }
        // Validate addresses
        if (config.anchorer == address(0)) revert ZeroAddress("anchorer");
        if (config.issuer == address(0)) revert ZeroAddress("issuer");

        return config;
    }

    /// @notice Deploy all contracts
    function _deploy(DeploymentConfig memory config)
        internal
        returns (DeploymentResult memory result)
    {
        console.log("--- Deployment Phase ---");

        // 1. Deploy FlemingAnchor (Phase B)
        console.log("Deploying FlemingAnchor...");
        result.anchor = new FlemingAnchor(config.anchorer);
        console.log("  Address:", address(result.anchor));
        console.log("  Initial anchorer:", config.anchorer);
        console.log("");

        // 2. Deploy VCRegistry (Phase C.1)
        console.log("Deploying VCRegistry...");
        result.vcRegistry = new VCRegistry(config.issuer);
        console.log("  Address:", address(result.vcRegistry));
        console.log("  Initial issuer:", config.issuer);
        console.log("");

        vm.stopBroadcast();

        return result;
    }

    /// @notice Verify deployment succeeded
    function _verifyDeployment(DeploymentResult memory result) internal pure {
        console.log("--- Verification Phase ---");

        // Verify all contracts deployed
        if (address(result.anchor) == address(0)) {
            revert DeploymentFailed("FlemingAnchor deployment failed");
        }
        if (address(result.vcRegistry) == address(0)) {
            revert DeploymentFailed("VCRegistry deployment failed");
        }

        console.log("All contracts deployed successfully");
        console.log("");
    }

    /// @notice Log deployment results and environment variables
    function _logDeployment(DeploymentResult memory result) internal pure {
        console.log("========================================");
        console.log("DEPLOYMENT COMPLETE");
        console.log("========================================");
        console.log("");
        console.log("Add to your .env:");
        console.log("ANCHOR_CONTRACT_ADDRESS=%s", address(result.anchor));
        console.log("VCREGISTRY_CONTRACT_ADDRESS=%s", address(result.vcRegistry));

        console.log("");
        console.log("Verify contracts:");
        console.log(
            "forge verify-contract %s FlemingAnchor --chain base-sepolia", address(result.anchor)
        );
        console.log(
            "forge verify-contract %s VCRegistry --chain base-sepolia", address(result.vcRegistry)
        );

        console.log("");
    }

    /// @notice Write deployment artifacts to JSON file
    function _writeDeploymentJSON(DeploymentResult memory result) internal {
        string memory json = "deployment_export";
        vm.serializeAddress(json, "anchor", address(result.anchor));
        string memory output = vm.serializeAddress(json, "vcRegistry", address(result.vcRegistry));

        // Local dev expects a stable `deployments.json` path (docker compose copies it to a shared volume
        // so the backend can auto-wire the anchor contract address).
        //
        // For non-local networks, write explicit per-network artifacts so they can be committed as a
        // human/audit source of truth without getting overwritten by local runs.
        string memory fileName = "deployments.json";
        if (chainId == BASE_SEPOLIA) {
            fileName = "deployments.base_sepolia.json";
        } else if (chainId == BASE_MAINNET) {
            fileName = "deployments.base.json";
        }

        vm.writeJson(output, fileName);
        console.log("Wrote deployments to:", fileName);
    }

    /// @notice Get human-readable network name
    function _getNetworkName(uint256 id) internal pure returns (string memory) {
        if (id == BASE_SEPOLIA) return "Base Sepolia";
        if (id == BASE_MAINNET) return "Base Mainnet";
        if (id == ANVIL_LOCAL) return "Anvil Local";
        return "Unknown";
    }
}
