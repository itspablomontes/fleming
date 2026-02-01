// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

import {Test, console} from "forge-std/Test.sol";
import {VCRegistry} from "../../src/VCRegistry.sol";

/// @title VCRegistryTest
/// @notice Comprehensive test suite for VCRegistry contract
/// @dev Tests cover: issuance, revocation, bitmap operations, batch operations
contract VCRegistryTest is Test {
    VCRegistry registry;
    address owner;
    address issuer;
    address unauthorized;

    // Test data
    bytes32 constant HASH_1 = keccak256("vc_content_1");
    bytes32 constant HASH_2 = keccak256("vc_content_2");
    bytes32 constant HASH_3 = keccak256("vc_content_3");
    bytes32 constant ZERO_HASH = bytes32(0);

    event VCIssued(
        uint256 indexed vcId,
        bytes32 indexed vcHash,
        address indexed issuer,
        uint256 timestamp
    );
    event VCRevoked(
        uint256 indexed vcId,
        address indexed revoker,
        uint256 timestamp
    );
    event IssuerUpdated(
        address indexed previousIssuer,
        address indexed newIssuer
    );

    function setUp() public {
        owner = makeAddr("owner");
        issuer = makeAddr("issuer");
        unauthorized = makeAddr("unauthorized");

        vm.prank(owner);
        registry = new VCRegistry(issuer);
    }

    // ============ UNIT TESTS ============

    function test_ConstructorSetsIssuer() public view {
        assertEq(registry.issuer(), issuer);
        assertEq(registry.owner(), owner);
        assertEq(registry.nextVCId(), 1);
    }

    function test_ConstructorDefaultsToOwner() public {
        vm.prank(owner);
        VCRegistry defaultRegistry = new VCRegistry(address(0));
        assertEq(defaultRegistry.issuer(), owner);
    }

    function test_IssueCreatesVC() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        assertEq(vcId, 1);
        assertEq(registry.vcHashes(vcId), HASH_1);
        assertEq(registry.totalIssued(), 1);
        assertGt(registry.issuanceTimestamps(vcId), 0);
    }

    function test_IssueOwnerCanAlsoIssue() public {
        vm.prank(owner);
        uint256 vcId = registry.issue(HASH_1);

        assertEq(vcId, 1);
    }

    function test_UnauthorizedCannotIssue() public {
        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(
                VCRegistry.Unauthorized.selector,
                unauthorized
            )
        );
        registry.issue(HASH_1);
    }

    function test_IssueIncrementsId() public {
        vm.startPrank(issuer);
        uint256 id1 = registry.issue(HASH_1);
        uint256 id2 = registry.issue(HASH_2);
        uint256 id3 = registry.issue(HASH_3);
        vm.stopPrank();

        assertEq(id1, 1);
        assertEq(id2, 2);
        assertEq(id3, 3);
        assertEq(registry.nextVCId(), 4);
    }

    // ============ REVOCATION TESTS ============

    function test_RevokeMarksAsRevoked() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        vm.prank(owner);
        registry.revoke(vcId);

        assertTrue(registry.isRevoked(vcId));
        assertEq(registry.totalRevoked(), 1);
    }

    function test_NonOwnerCannotRevoke() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        vm.prank(issuer);
        vm.expectRevert();
        registry.revoke(vcId);
    }

    function test_CannotRevokeInvalidId() public {
        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(VCRegistry.InvalidVCId.selector, 999)
        );
        registry.revoke(999);
    }

    function test_CannotRevokeZeroId() public {
        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(VCRegistry.InvalidVCId.selector, 0)
        );
        registry.revoke(0);
    }

    function test_CannotRevokeAlreadyRevoked() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        vm.startPrank(owner);
        registry.revoke(vcId);

        vm.expectRevert(
            abi.encodeWithSelector(VCRegistry.AlreadyRevoked.selector, vcId)
        );
        registry.revoke(vcId);
        vm.stopPrank();
    }

    function test_CannotRevokeUnissued() public {
        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(VCRegistry.InvalidVCId.selector, 1)
        );
        registry.revoke(1);
    }

    // ============ IS REVOKED TESTS ============

    function test_IsRevokedReturnsFalseForUnissued() public view {
        assertFalse(registry.isRevoked(1));
    }

    function test_IsRevokedReturnsFalseForZero() public view {
        assertFalse(registry.isRevoked(0));
    }

    function test_IsRevokedReturnsTrueAfterRevoke() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        assertFalse(registry.isRevoked(vcId));

        vm.prank(owner);
        registry.revoke(vcId);

        assertTrue(registry.isRevoked(vcId));
    }

    // ============ VERIFY VC TESTS ============

    function test_VerifyVCSuccess() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        (bool valid, bool revoked, uint256 issuedAt) = registry.verifyVC(
            vcId,
            HASH_1
        );

        assertTrue(valid);
        assertFalse(revoked);
        assertGt(issuedAt, 0);
    }

    function test_VerifyVCFailsWithWrongHash() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        (bool valid, bool revoked, uint256 issuedAt) = registry.verifyVC(
            vcId,
            HASH_2
        );

        assertFalse(valid);
        assertFalse(revoked);
        assertGt(issuedAt, 0);
    }

    function test_VerifyVCFailsWhenRevoked() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        vm.prank(owner);
        registry.revoke(vcId);

        (bool valid, bool revoked, ) = registry.verifyVC(vcId, HASH_1);

        assertFalse(valid);
        assertTrue(revoked);
    }

    function test_VerifyVCInvalidId() public view {
        (bool valid, bool revoked, uint256 issuedAt) = registry.verifyVC(
            999,
            HASH_1
        );

        assertFalse(valid);
        assertFalse(revoked);
        assertEq(issuedAt, 0);
    }

    // ============ BATCH OPERATIONS ============

    function test_BatchIssueMultipleVCs() public {
        bytes32[] memory hashes = new bytes32[](3);
        hashes[0] = HASH_1;
        hashes[1] = HASH_2;
        hashes[2] = HASH_3;

        vm.prank(issuer);
        uint256[] memory ids = registry.batchIssue(hashes);

        assertEq(ids.length, 3);
        assertEq(ids[0], 1);
        assertEq(ids[1], 2);
        assertEq(ids[2], 3);
        assertEq(registry.totalIssued(), 3);
    }

    function test_BatchIssueSkipsZeroHashes() public {
        bytes32[] memory hashes = new bytes32[](3);
        hashes[0] = HASH_1;
        hashes[1] = ZERO_HASH;
        hashes[2] = HASH_2;

        vm.prank(issuer);
        uint256[] memory ids = registry.batchIssue(hashes);

        // ID 2 was skipped (zero hash)
        assertEq(ids[0], 1);
        assertEq(ids[1], 0); // Skipped
        assertEq(ids[2], 2);
    }

    function test_BatchRevokeMultipleVCs() public {
        vm.startPrank(issuer);
        uint256 id1 = registry.issue(HASH_1);
        uint256 id2 = registry.issue(HASH_2);
        uint256 id3 = registry.issue(HASH_3);
        vm.stopPrank();

        uint256[] memory ids = new uint256[](3);
        ids[0] = id1;
        ids[1] = id2;
        ids[2] = id3;

        vm.prank(owner);
        registry.batchRevoke(ids);

        assertTrue(registry.isRevoked(id1));
        assertTrue(registry.isRevoked(id2));
        assertTrue(registry.isRevoked(id3));
        assertEq(registry.totalRevoked(), 3);
    }

    function test_BatchRevokeSkipsInvalid() public {
        vm.prank(issuer);
        uint256 id1 = registry.issue(HASH_1);

        uint256[] memory ids = new uint256[](3);
        ids[0] = id1;
        ids[1] = 999; // Invalid, should skip
        ids[2] = 0; // Zero, should skip

        vm.prank(owner);
        registry.batchRevoke(ids);

        assertTrue(registry.isRevoked(id1));
        assertEq(registry.totalRevoked(), 1);
    }

    function test_BatchRevokeSkipsAlreadyRevoked() public {
        vm.startPrank(issuer);
        uint256 id1 = registry.issue(HASH_1);
        uint256 id2 = registry.issue(HASH_2);
        vm.stopPrank();

        vm.prank(owner);
        registry.revoke(id1);

        uint256[] memory ids = new uint256[](2);
        ids[0] = id1; // Already revoked
        ids[1] = id2;

        vm.prank(owner);
        registry.batchRevoke(ids);

        assertEq(registry.totalRevoked(), 2); // Both revoked
    }

    // ============ ADMIN TESTS ============

    function test_OwnerCanUpdateIssuer() public {
        address newIssuer = makeAddr("newIssuer");

        vm.prank(owner);
        registry.setIssuer(newIssuer);

        assertEq(registry.issuer(), newIssuer);
    }

    function test_NonOwnerCannotUpdateIssuer() public {
        address newIssuer = makeAddr("newIssuer");

        vm.prank(issuer);
        vm.expectRevert();
        registry.setIssuer(newIssuer);
    }

    // ============ EVENT TESTS ============

    function test_Emit_VCIssued() public {
        vm.prank(issuer);

        vm.expectEmit(true, true, true, true);
        emit VCIssued(1, HASH_1, issuer, block.timestamp);

        registry.issue(HASH_1);
    }

    function test_Emit_VCRevoked() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        vm.prank(owner);
        vm.expectEmit(true, true, false, true);
        emit VCRevoked(vcId, owner, block.timestamp);

        registry.revoke(vcId);
    }

    // ============ FUZZ TESTS ============

    function testFuzz_IssueAnyNonZeroHash(bytes32 hash) public {
        vm.assume(hash != bytes32(0));

        vm.prank(issuer);
        uint256 vcId = registry.issue(hash);

        assertGt(vcId, 0);
        assertEq(registry.vcHashes(vcId), hash);
    }

    function testFuzz_RevokeAndCheck(
        uint8 numIssues,
        uint8 revokeIndex
    ) public {
        vm.assume(numIssues > 0 && numIssues <= 100);
        vm.assume(revokeIndex < numIssues);

        vm.startPrank(issuer);
        for (uint256 i = 0; i < numIssues; ++i) {
            registry.issue(keccak256(abi.encodePacked(i)));
        }
        vm.stopPrank();

        uint256 targetId = uint256(revokeIndex) + 1;

        vm.prank(owner);
        registry.revoke(targetId);

        assertTrue(registry.isRevoked(targetId));
        assertEq(registry.totalRevoked(), 1);
    }

    // ============ GAS TESTS ============

    function testGas_Issue() public {
        vm.prank(issuer);

        uint256 gasBefore = gasleft();
        registry.issue(HASH_1);
        uint256 gasUsed = gasBefore - gasleft();

        console.log("Gas used for issue():", gasUsed);
        assertLt(gasUsed, 70000, "Issue should use less than 70k gas");
    }

    function testGas_Revoke() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        vm.prank(owner);
        uint256 gasBefore = gasleft();
        registry.revoke(vcId);
        uint256 gasUsed = gasBefore - gasleft();

        console.log("Gas used for revoke():", gasUsed);
        assertLt(gasUsed, 55000, "Revoke should use less than 55k gas");
    }

    function testGas_IsRevoked() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        uint256 gasBefore = gasleft();
        registry.isRevoked(vcId);
        uint256 gasUsed = gasBefore - gasleft();

        console.log("Gas used for isRevoked() (unrevoked):", gasUsed);
        assertLt(gasUsed, 7000, "isRevoked should use less than 7k gas (cold)");
    }

    function testGas_VerifyVC() public {
        vm.prank(issuer);
        uint256 vcId = registry.issue(HASH_1);

        uint256 gasBefore = gasleft();
        registry.verifyVC(vcId, HASH_1);
        uint256 gasUsed = gasBefore - gasleft();

        console.log("Gas used for verifyVC():", gasUsed);
        assertLt(gasUsed, 7000, "verifyVC should use less than 7k gas (cold)");
    }

    function testGas_BatchIssue() public {
        bytes32[] memory hashes = new bytes32[](10);
        for (uint256 i = 0; i < 10; ++i) {
            hashes[i] = keccak256(abi.encodePacked(i));
        }

        vm.prank(issuer);

        uint256 gasBefore = gasleft();
        registry.batchIssue(hashes);
        uint256 gasUsed = gasBefore - gasleft();

        console.log("Gas used for batchIssue(10):", gasUsed);
        console.log("Gas per VC:", gasUsed / 10);
        assertLt(
            gasUsed,
            500000,
            "Batch issue of 10 should use less than 500k gas"
        );
    }
}
