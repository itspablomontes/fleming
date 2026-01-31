// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/// @title FlemingAnchor
/// @notice Immutable on-chain anchoring of audit log Merkle roots
/// @dev Once a root is anchored, it cannot be changed or deleted.
/// @custom:security-contact security@fleming.health
/// @custom:version 1.0.0
contract FlemingAnchor is Ownable {
    // ─── Custom Errors ─────────────────────────────────────────────────────────

    /// @notice Root has already been anchored
    /// @param root The root that was already anchored
    error AlreadyAnchored(bytes32 root);

    /// @notice Root cannot be zero bytes
    error ZeroRoot();

    /// @notice New owner cannot be zero address
    error ZeroAddress();

    // ─── Events ────────────────────────────────────────────────────────────────

    /// @notice Emitted when a new Merkle root is anchored
    /// @param root The Merkle root (SHA-256 hash, 32 bytes)
    /// @param timestamp The block timestamp at anchoring time
    /// @param blockNumber The block number at anchoring time
    /// @param anchorer The address that performed the anchoring
    event RootAnchored(
        bytes32 indexed root,
        uint256 timestamp,
        uint256 blockNumber,
        address indexed anchorer
    );

    /// @notice Emitted when the anchorer address is changed
    /// @param previousAnchorer The previous allowed anchorer
    /// @param newAnchorer The new allowed anchorer
    event AnchorerUpdated(
        address indexed previousAnchorer,
        address indexed newAnchorer
    );

    // ─── Storage ───────────────────────────────────────────────────────────────

    /// @notice Mapping from Merkle root to anchor timestamp (0 if not anchored)
    /// @dev Packed slot: timestamp (32 bytes) = 1 storage slot per entry
    mapping(bytes32 => uint256) public anchors;

    /// @notice Total number of anchored roots (for statistics and monitoring)
    /// @dev Uses unchecked increments for gas efficiency
    uint256 public anchorCount;

    /// @notice Address authorized to anchor roots (can be a multisig or contract)
    /// @dev Separate from owner to allow automated anchoring by backend
    address public anchorer;

    // ─── Modifiers ─────────────────────────────────────────────────────────────

    /// @notice Restricts function to authorized anchorer only
    modifier onlyAnchorer() {
        if (msg.sender != anchorer && msg.sender != owner()) {
            revert Unauthorized(msg.sender);
        }
        _;
    }

    /// @notice Thrown when caller is not authorized
    error Unauthorized(address caller);

    // ─── Constructor ───────────────────────────────────────────────────────────

    /// @notice Deploy contract and set initial anchorer
    /// @param initialAnchorer The address initially authorized to anchor roots
    /// @dev If initialAnchorer is address(0), defaults to deployer
    constructor(address initialAnchorer) Ownable(msg.sender) {
        if (initialAnchorer == address(0)) {
            anchorer = msg.sender;
        } else {
            anchorer = initialAnchorer;
        }
        emit AnchorerUpdated(address(0), anchorer);
    }

    // ─── External Functions ────────────────────────────────────────────────────

    /// @notice Anchor a Merkle root on-chain
    /// @dev Reverts if root is zero or already anchored. CEI pattern enforced.
    /// @param root The Merkle root hash (SHA-256, 32 bytes)
    /// @custom:gas-target < 50,000 gas
    function anchor(bytes32 root) external onlyAnchorer {
        // CHECKS
        if (root == bytes32(0)) revert ZeroRoot();
        if (anchors[root] != 0) revert AlreadyAnchored(root);

        // EFFECTS
        anchors[root] = block.timestamp;
        unchecked {
            ++anchorCount;
        }

        // INTERACTIONS (none for this function, but emit event)
        emit RootAnchored(root, block.timestamp, block.number, msg.sender);
    }

    /// @notice Batch anchor multiple Merkle roots in a single transaction
    /// @dev More gas efficient than individual anchors for multiple roots
    /// @param roots Array of Merkle root hashes to anchor
    /// @custom:gas Saves ~15k gas per additional root vs individual calls
    function batchAnchor(bytes32[] calldata roots) external onlyAnchorer {
        uint256 length = roots.length;
        if (length == 0) return;

        uint256 currentTimestamp = block.timestamp;
        uint256 currentBlock = block.number;
        uint256 anchoredCount;

        for (uint256 i = 0; i < length; ++i) {
            bytes32 root = roots[i];

            // Skip if already anchored (don't revert, continue batch)
            if (root == bytes32(0) || anchors[root] != 0) continue;

            anchors[root] = currentTimestamp;
            emit RootAnchored(root, currentTimestamp, currentBlock, msg.sender);

            unchecked {
                ++anchoredCount;
            }
        }

        if (anchoredCount == 0) return;
        unchecked {
            anchorCount += anchoredCount;
        }
    }

    /// @notice Check if a Merkle root has been anchored
    /// @param root The Merkle root to verify
    /// @return timestamp The anchor timestamp (0 if not anchored)
    /// @custom:gas ~2,400 gas (cold) / ~100 gas (warm)
    function getAnchorTimestamp(
        bytes32 root
    ) external view returns (uint256 timestamp) {
        return anchors[root];
    }

    /// @notice Check if a root is anchored (boolean convenience)
    /// @param root The Merkle root to check
    /// @return True if the root has been anchored
    /// @custom:gas ~2,400 gas (cold) / ~100 gas (warm)
    function isAnchored(bytes32 root) external view returns (bool) {
        return anchors[root] != 0;
    }

    /// @notice Get the age of an anchored root in seconds
    /// @param root The Merkle root to check
    /// @return age Seconds since anchoring (0 if not anchored)
    /// @custom:gas ~2,500 gas
    function getAnchorAge(bytes32 root) external view returns (uint256 age) {
        uint256 timestamp = anchors[root];
        if (timestamp == 0) return 0;
        return block.timestamp - timestamp;
    }

    // ─── Admin Functions ───────────────────────────────────────────────────────

    /// @notice Update the authorized anchorer address
    /// @param newAnchorer The new address authorized to anchor roots
    /// @dev Only callable by owner. Set to multisig for production.
    function setAnchorer(address newAnchorer) external onlyOwner {
        if (newAnchorer == address(0)) revert ZeroAddress();

        address previousAnchorer = anchorer;
        anchorer = newAnchorer;

        emit AnchorerUpdated(previousAnchorer, newAnchorer);
    }

    /// @notice Renounce anchorer role (emergency stop for automated anchoring)
    /// @dev Can only be called by owner. After this, only owner can anchor.
    function renounceAnchorer() external onlyOwner {
        address previousAnchorer = anchorer;
        anchorer = address(0);

        emit AnchorerUpdated(previousAnchorer, address(0));
    }

    // ─── View Functions ────────────────────────────────────────────────────────

    /// @notice Get contract metadata for verification
    /// @return version Contract version string
    /// @return chainId Current chain ID
    function getContractInfo()
        external
        view
        returns (string memory version, uint256 chainId)
    {
        return ("1.0.0", block.chainid);
    }
}
