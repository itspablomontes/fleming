// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {Bitmap} from "./libraries/Bitmap.sol";

/// @title VCRegistry
/// @notice Verifiable Credential registry with efficient revocation
/// @dev Stores VC hashes and uses bitmap for revocation status.
contract VCRegistry is Ownable {
    using Bitmap for uint256[];

    // ─── Custom Errors ─────────────────────────────────────────────────────────

    /// @notice VC ID is invalid (not yet issued)
    /// @param vcId The invalid VC ID
    error InvalidVCId(uint256 vcId);

    /// @notice VC has already been revoked
    /// @param vcId The VC ID that was already revoked
    error AlreadyRevoked(uint256 vcId);

    /// @notice VC hash cannot be zero
    error ZeroHash();

    /// @notice Address cannot be zero
    error ZeroAddress();

    /// @notice Caller is not authorized for this operation
    /// @param caller The unauthorized caller
    error Unauthorized(address caller);

    // ─── Events ────────────────────────────────────────────────────────────────

    /// @notice Emitted when a new VC is issued
    /// @param vcId Unique identifier for the VC
    /// @param vcHash Hash of the VC content (for verification)
    /// @param issuer Address that issued the VC
    /// @param timestamp Block timestamp of issuance
    event VCIssued(
        uint256 indexed vcId,
        bytes32 indexed vcHash,
        address indexed issuer,
        uint256 timestamp
    );

    /// @notice Emitted when a VC is revoked
    /// @param vcId The VC ID that was revoked
    /// @param revoker Address that performed the revocation
    /// @param timestamp Block timestamp of revocation
    event VCRevoked(
        uint256 indexed vcId,
        address indexed revoker,
        uint256 timestamp
    );

    /// @notice Emitted when the issuer role is updated
    /// @param previousIssuer The previous authorized issuer
    /// @param newIssuer The new authorized issuer
    event IssuerUpdated(
        address indexed previousIssuer,
        address indexed newIssuer
    );

    // ─── Storage ───────────────────────────────────────────────────────────────

    /// @notice Bitmap for revocation status (1 bit per VC)
    /// @dev Index = vcId, bit = 1 if revoked. Supports 2^248 VCs (practically unlimited)
    uint256[] public revocationBitmap;

    /// @notice Mapping from VC ID to VC content hash
    /// @dev Used to verify VC integrity off-chain
    mapping(uint256 => bytes32) public vcHashes;

    /// @notice Mapping from VC ID to issuance timestamp
    mapping(uint256 => uint256) public issuanceTimestamps;

    /// @notice Next available VC ID (auto-incrementing)
    /// @dev Starts at 1, so vcId 0 is invalid
    uint256 public nextVCId;

    /// @notice Total number of revoked VCs
    uint256 public totalRevoked;

    /// @notice Address authorized to issue new VCs
    address public issuer;

    // ─── Modifiers ─────────────────────────────────────────────────────────────

    /// @notice Restricts function to authorized issuer only
    modifier onlyIssuer() {
        if (msg.sender != issuer && msg.sender != owner()) {
            revert Unauthorized(msg.sender);
        }
        _;
    }

    // ─── Constructor ───────────────────────────────────────────────────────────

    /// @notice Deploy contract and set initial issuer
    /// @param initialIssuer The address initially authorized to issue VCs
    constructor(address initialIssuer) Ownable(msg.sender) {
        nextVCId = 1; // Start at 1 so 0 is always invalid
        // Pre-allocate the first slot to avoid a costly first-time push on revoke.
        revocationBitmap.push(0);

        if (initialIssuer == address(0)) {
            issuer = msg.sender;
        } else {
            issuer = initialIssuer;
        }

        emit IssuerUpdated(address(0), issuer);
    }

    // ─── External Functions ────────────────────────────────────────────────────

    /// @notice Issue a new Verifiable Credential
    /// @param vcHash Hash of the VC content (keccak256 of SD-JWT or similar)
    /// @return vcId The unique identifier assigned to this VC
    function issue(bytes32 vcHash) external onlyIssuer returns (uint256 vcId) {
        if (vcHash == bytes32(0)) revert ZeroHash();

        vcId = nextVCId;
        unchecked {
            ++nextVCId;
        }

        vcHashes[vcId] = vcHash;
        issuanceTimestamps[vcId] = block.timestamp;

        emit VCIssued(vcId, vcHash, msg.sender, block.timestamp);

        return vcId;
    }

    /// @notice Batch issue multiple VCs in a single transaction
    /// @param hashes Array of VC content hashes
    /// @return vcIds Array of assigned VC IDs
    function batchIssue(
        bytes32[] calldata hashes
    ) external onlyIssuer returns (uint256[] memory vcIds) {
        uint256 length = hashes.length;
        if (length == 0) return new uint256[](0);

        vcIds = new uint256[](length);
        uint256 currentId = nextVCId;
        uint256 currentTimestamp = block.timestamp;

        for (uint256 i = 0; i < length; ++i) {
            bytes32 hash = hashes[i];
            if (hash == bytes32(0)) continue; // Skip invalid hashes

            uint256 vcId = currentId;
            unchecked {
                ++currentId;
            }

            vcHashes[vcId] = hash;
            issuanceTimestamps[vcId] = currentTimestamp;
            vcIds[i] = vcId;

            emit VCIssued(vcId, hash, msg.sender, currentTimestamp);
        }

        nextVCId = currentId;

        return vcIds;
    }

    /// @notice Revoke a Verifiable Credential
    /// @param vcId The VC ID to revoke
    /// @dev Only callable by owner. Irreversible operation.
    function revoke(uint256 vcId) external onlyOwner {
        _revoke(vcId);
    }

    /// @notice Batch revoke multiple VCs
    /// @param vcIds Array of VC IDs to revoke
    function batchRevoke(uint256[] calldata vcIds) external onlyOwner {
        uint256 length = vcIds.length;
        if (length == 0) return;

        for (uint256 i = 0; i < length; ++i) {
            uint256 vcId = vcIds[i];

            // Skip if invalid or already revoked
            if (vcId == 0 || vcId >= nextVCId) continue;
            if (revocationBitmap.get(vcId)) continue;

            _revoke(vcId);
        }
    }

    /// @notice Check if a VC has been revoked
    /// @param vcId The VC ID to check
    /// @return True if the VC has been revoked
    function isRevoked(uint256 vcId) external view returns (bool) {
        if (vcId == 0) return false;
        return revocationBitmap.get(vcId);
    }

    /// @notice Total number of issued VCs
    /// @dev Derived from `nextVCId` since IDs are contiguous.
    function totalIssued() external view returns (uint256) {
        unchecked {
            return nextVCId - 1;
        }
    }

    /// @notice Verify a VC's status and hash
    /// @param vcId The VC ID to verify
    /// @param expectedHash The expected VC content hash
    /// @return valid True if VC exists, hash matches, and not revoked
    /// @return revoked True if the VC has been revoked
    /// @return issuedAt The timestamp when the VC was issued (0 if not found)
    function verifyVC(
        uint256 vcId,
        bytes32 expectedHash
    ) external view returns (bool valid, bool revoked, uint256 issuedAt) {
        if (vcId == 0 || vcId >= nextVCId) {
            return (false, false, 0);
        }

        issuedAt = issuanceTimestamps[vcId];
        if (issuedAt == 0) {
            return (false, false, 0);
        }

        revoked = revocationBitmap.get(vcId);
        valid = !revoked && vcHashes[vcId] == expectedHash;

        return (valid, revoked, issuedAt);
    }

    /// @notice Get the revocation count up to a specific VC ID
    /// @param upToVCId Count revocations up to and including this ID
    /// @return count Number of revoked VCs in the range
    /// @dev Gas intensive for large ranges - prefer off-chain indexing
    function getRevocationCount(
        uint256 upToVCId
    ) external view returns (uint256 count) {
        if (upToVCId == 0 || upToVCId >= nextVCId) return 0;
        return revocationBitmap.countSetBits(upToVCId);
    }

    // ─── Admin Functions ───────────────────────────────────────────────────────

    /// @notice Update the authorized issuer address
    /// @param newIssuer The new address authorized to issue VCs
    function setIssuer(address newIssuer) external onlyOwner {
        if (newIssuer == address(0)) revert ZeroAddress();

        address previousIssuer = issuer;
        issuer = newIssuer;

        emit IssuerUpdated(previousIssuer, newIssuer);
    }

    // ─── Internal Functions ────────────────────────────────────────────────────

    /// @notice Internal revoke function (used by revoke and batchRevoke)
    /// @param vcId The VC ID to revoke
    function _revoke(uint256 vcId) internal {
        if (vcId == 0 || vcId >= nextVCId) revert InvalidVCId(vcId);
        if (revocationBitmap.get(vcId)) revert AlreadyRevoked(vcId);

        revocationBitmap.set(vcId);

        unchecked {
            ++totalRevoked;
        }

        emit VCRevoked(vcId, msg.sender, block.timestamp);
    }
}
