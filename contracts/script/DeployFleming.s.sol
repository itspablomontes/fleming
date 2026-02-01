// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

import {Script, console} from "forge-std/Script.sol";
import {FlemingAnchor} from "../src/FlemingAnchor.sol";
import {VCRegistry} from "../src/VCRegistry.sol";
import {ZKVerifier} from "../src/ZKVerifier.sol";

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

    event DeploymentComplete(
        address indexed anchor,
        address indexed vcRegistry,
        address indexed zkVerifier,
        uint256 chainId
    );

    // ─── Structs ───────────────────────────────────────────────────────────────

    struct DeploymentConfig {
        address anchorer;
        address issuer;
        uint256 expectedPublicInputCount;
    }

    struct DeploymentResult {
        FlemingAnchor anchor;
        VCRegistry vcRegistry;
        ZKVerifier zkVerifier;
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

        emit DeploymentComplete(
            address(result.anchor),
            address(result.vcRegistry),
            address(result.zkVerifier),
            chainId
        );

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
        if (
            chainId != BASE_SEPOLIA &&
            chainId != BASE_MAINNET &&
            chainId != ANVIL_LOCAL
        ) {
            revert InvalidChain(chainId);
        }

        vm.startBroadcast(deployerKey);
    }

    /// @notice Load deployment configuration from environment
    function _loadConfig()
        internal
        view
        returns (DeploymentConfig memory config)
    {
        // Try to load from env, fallback to deployer address
        try vm.envAddress("ANCHORER_ADDRESS") returns (address addr) {
            config.anchorer = addr;
        } catch {
            config.anchorer = deployer;
            console.log(
                "ANCHORER_ADDRESS not set, using deployer:",
                config.anchorer
            );
        }
        try vm.envAddress("ISSUER_ADDRESS") returns (address addr) {
            config.issuer = addr;
        } catch {
            config.issuer = deployer;
            console.log(
                "ISSUER_ADDRESS not set, using deployer:",
                config.issuer
            );
        }
        try vm.envUint("ZK_PUBLIC_INPUTS") returns (uint256 count) {
            config.expectedPublicInputCount = count;
        } catch {
            config.expectedPublicInputCount = 2; // Default for basic circuits
            console.log(
                "ZK_PUBLIC_INPUTS not set, using default:",
                config.expectedPublicInputCount
            );
        }
        // Validate addresses
        if (config.anchorer == address(0)) revert ZeroAddress("anchorer");
        if (config.issuer == address(0)) revert ZeroAddress("issuer");

        return config;
    }

    /// @notice Deploy all contracts
    function _deploy(
        DeploymentConfig memory config
    ) internal returns (DeploymentResult memory result) {
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

        // 3. Deploy ZKVerifier (Phase C.2)
        console.log("Deploying ZKVerifier...");
        result.zkVerifier = new ZKVerifier(config.expectedPublicInputCount);
        console.log("  Address:", address(result.zkVerifier));
        console.log(
            "  Expected public inputs:",
            config.expectedPublicInputCount
        );
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
        if (address(result.zkVerifier) == address(0)) {
            revert DeploymentFailed("ZKVerifier deployment failed");
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
        console.log(
            "VCREGISTRY_CONTRACT_ADDRESS=%s",
            address(result.vcRegistry)
        );
        console.log(
            "ZKVERIFIER_CONTRACT_ADDRESS=%s",
            address(result.zkVerifier)
        );
        console.log("");
        console.log("Verify contracts:");
        console.log(
            "forge verify-contract %s FlemingAnchor --chain base-sepolia",
            address(result.anchor)
        );
        console.log(
            "forge verify-contract %s VCRegistry --chain base-sepolia",
            address(result.vcRegistry)
        );
        console.log(
            "forge verify-contract %s ZKVerifier --chain base-sepolia",
            address(result.zkVerifier)
        );
        console.log("");
    }

    /// @notice Get human-readable network name
    function _getNetworkName(uint256 id) internal pure returns (string memory) {
        if (id == BASE_SEPOLIA) return "Base Sepolia";
        if (id == BASE_MAINNET) return "Base Mainnet";
        if (id == ANVIL_LOCAL) return "Anvil Local";
        return "Unknown";
    }
}
