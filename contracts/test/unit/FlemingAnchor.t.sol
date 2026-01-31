// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

import {Test, console} from "forge-std/Test.sol";
import {FlemingAnchor} from "../../src/FlemingAnchor.sol";

/// @title FlemingAnchorTest
/// @notice Comprehensive test suite for FlemingAnchor contract
/// @dev Tests cover: unit, fuzz, gas, and edge cases
contract FlemingAnchorTest is Test {
    FlemingAnchor anchor;
    address owner;
    address anchorer;
    address unauthorized;

    // Test data
    bytes32 constant ROOT_1 = keccak256("audit_batch_1");
    bytes32 constant ROOT_2 = keccak256("audit_batch_2");
    bytes32 constant ROOT_3 = keccak256("audit_batch_3");
    bytes32 constant ZERO_ROOT = bytes32(0);

    event RootAnchored(
        bytes32 indexed root,
        uint256 timestamp,
        uint256 blockNumber,
        address indexed anchorer
    );
    event AnchorerUpdated(
        address indexed previousAnchorer,
        address indexed newAnchorer
    );

    function setUp() public {
        owner = makeAddr("owner");
        anchorer = makeAddr("anchorer");
        unauthorized = makeAddr("unauthorized");

        vm.prank(owner);
        anchor = new FlemingAnchor(anchorer);
    }

    // ============ UNIT TESTS ============

    function test_ConstructorSetsAnchorer() public view {
        assertEq(anchor.anchorer(), anchorer);
        assertEq(anchor.owner(), owner);
    }

    function test_ConstructorDefaultsToOwner() public {
        vm.prank(owner);
        FlemingAnchor defaultAnchor = new FlemingAnchor(address(0));
        assertEq(defaultAnchor.anchorer(), owner);
    }

    function test_AnchorStoresTimestamp() public {
        uint256 expectedTimestamp = block.timestamp;

        vm.prank(anchorer);
        anchor.anchor(ROOT_1);

        assertEq(anchor.anchors(ROOT_1), expectedTimestamp);
        assertEq(anchor.anchorCount(), 1);
    }

    function test_AnchorOwnerCanAlsoAnchor() public {
        vm.prank(owner);
        anchor.anchor(ROOT_1);

        assertTrue(anchor.isAnchored(ROOT_1));
    }

    function test_UnauthorizedCannotAnchor() public {
        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(
                FlemingAnchor.Unauthorized.selector,
                unauthorized
            )
        );
        anchor.anchor(ROOT_1);
    }

    function test_IsAnchoredReturnsTrueAfterAnchor() public {
        assertFalse(anchor.isAnchored(ROOT_1));

        vm.prank(anchorer);
        anchor.anchor(ROOT_1);

        assertTrue(anchor.isAnchored(ROOT_1));
    }

    function test_GetAnchorTimestampReturnsZeroForUnanchored() public view {
        assertEq(anchor.getAnchorTimestamp(ROOT_1), 0);
    }

    function test_GetAnchorAgeReturnsZeroForUnanchored() public view {
        assertEq(anchor.getAnchorAge(ROOT_1), 0);
    }

    function test_GetAnchorAgeCalculatesCorrectly() public {
        vm.prank(anchorer);
        anchor.anchor(ROOT_1);

        uint256 ageBeforeWarp = anchor.getAnchorAge(ROOT_1);
        assertEq(ageBeforeWarp, 0);

        vm.warp(block.timestamp + 1 hours);
        uint256 ageAfterWarp = anchor.getAnchorAge(ROOT_1);
        assertEq(ageAfterWarp, 1 hours);
    }

    function test_MultipleRootsCanBeAnchored() public {
        vm.startPrank(anchorer);
        anchor.anchor(ROOT_1);
        anchor.anchor(ROOT_2);
        anchor.anchor(ROOT_3);
        vm.stopPrank();

        assertTrue(anchor.isAnchored(ROOT_1));
        assertTrue(anchor.isAnchored(ROOT_2));
        assertTrue(anchor.isAnchored(ROOT_3));
        assertEq(anchor.anchorCount(), 3);
    }

    // ============ REVERT TESTS ============

    function test_RevertIf_AnchorZeroRoot() public {
        vm.prank(anchorer);
        vm.expectRevert(FlemingAnchor.ZeroRoot.selector);
        anchor.anchor(ZERO_ROOT);
    }

    function test_RevertIf_AnchorDuplicate() public {
        vm.prank(anchorer);
        anchor.anchor(ROOT_1);

        vm.prank(anchorer);
        vm.expectRevert(
            abi.encodeWithSelector(
                FlemingAnchor.AlreadyAnchored.selector,
                ROOT_1
            )
        );
        anchor.anchor(ROOT_1);
    }

    // ============ EVENT TESTS ============

    function test_Emit_RootAnchored() public {
        vm.prank(anchorer);

        vm.expectEmit(true, true, false, true);
        emit RootAnchored(ROOT_1, block.timestamp, block.number, anchorer);

        anchor.anchor(ROOT_1);
    }

    // ============ BATCH TESTS ============

    function test_BatchAnchorMultipleRoots() public {
        bytes32[] memory roots = new bytes32[](3);
        roots[0] = ROOT_1;
        roots[1] = ROOT_2;
        roots[2] = ROOT_3;

        vm.prank(anchorer);
        anchor.batchAnchor(roots);

        assertTrue(anchor.isAnchored(ROOT_1));
        assertTrue(anchor.isAnchored(ROOT_2));
        assertTrue(anchor.isAnchored(ROOT_3));
        assertEq(anchor.anchorCount(), 3);
    }

    function test_BatchAnchorSkipsInvalidRoots() public {
        bytes32[] memory roots = new bytes32[](3);
        roots[0] = ROOT_1;
        roots[1] = ZERO_ROOT; // Should be skipped
        roots[2] = ROOT_2;

        vm.prank(anchorer);
        anchor.batchAnchor(roots);

        assertTrue(anchor.isAnchored(ROOT_1));
        assertFalse(anchor.isAnchored(ZERO_ROOT));
        assertTrue(anchor.isAnchored(ROOT_2));
    }

    function test_BatchAnchorSkipsAlreadyAnchored() public {
        vm.prank(anchorer);
        anchor.anchor(ROOT_1);

        bytes32[] memory roots = new bytes32[](2);
        roots[0] = ROOT_1; // Already anchored, should skip
        roots[1] = ROOT_2;

        vm.prank(anchorer);
        anchor.batchAnchor(roots);

        assertEq(anchor.anchorCount(), 2); // Only ROOT_2 added
    }

    function test_BatchAnchorEmptyArray() public {
        bytes32[] memory roots = new bytes32[](0);

        vm.prank(anchorer);
        anchor.batchAnchor(roots); // Should not revert

        assertEq(anchor.anchorCount(), 0);
    }

    // ============ ADMIN TESTS ============

    function test_OwnerCanUpdateAnchorer() public {
        address newAnchorer = makeAddr("newAnchorer");

        vm.prank(owner);
        anchor.setAnchorer(newAnchorer);

        assertEq(anchor.anchorer(), newAnchorer);
    }

    function test_NonOwnerCannotUpdateAnchorer() public {
        address newAnchorer = makeAddr("newAnchorer");

        vm.prank(anchorer);
        vm.expectRevert();
        anchor.setAnchorer(newAnchorer);
    }

    function test_CannotSetAnchorerToZero() public {
        vm.prank(owner);
        vm.expectRevert(FlemingAnchor.ZeroAddress.selector);
        anchor.setAnchorer(address(0));
    }

    function test_OwnerCanRenounceAnchorer() public {
        vm.prank(owner);
        anchor.renounceAnchorer();

        assertEq(anchor.anchorer(), address(0));
    }

    function test_AfterRenounceOnlyOwnerCanAnchor() public {
        vm.prank(owner);
        anchor.renounceAnchorer();

        // Anchorer can no longer anchor
        vm.prank(anchorer);
        vm.expectRevert(
            abi.encodeWithSelector(
                FlemingAnchor.Unauthorized.selector,
                anchorer
            )
        );
        anchor.anchor(ROOT_1);

        // But owner still can
        vm.prank(owner);
        anchor.anchor(ROOT_1);
        assertTrue(anchor.isAnchored(ROOT_1));
    }

    // ============ FUZZ TESTS ============

    function testFuzz_AnchorAnyNonZeroRoot(bytes32 root) public {
        vm.assume(root != bytes32(0));

        vm.prank(anchorer);
        anchor.anchor(root);

        assertEq(anchor.anchors(root), block.timestamp);
        assertTrue(anchor.isAnchored(root));
    }

    function testFuzz_CannotAnchorSameRootTwice(bytes32 root) public {
        vm.assume(root != bytes32(0));

        vm.prank(anchorer);
        anchor.anchor(root);

        vm.prank(anchorer);
        vm.expectRevert();
        anchor.anchor(root);
    }

    function testFuzz_BatchAnchor(uint8 count) public {
        vm.assume(count > 0 && count <= 100);

        bytes32[] memory roots = new bytes32[](count);
        for (uint256 i = 0; i < count; ++i) {
            roots[i] = keccak256(abi.encodePacked(i));
        }

        vm.prank(anchorer);
        anchor.batchAnchor(roots);

        assertEq(anchor.anchorCount(), count);
    }

    // ============ GAS TESTS ============

    function testGas_Anchor() public {
        vm.prank(anchorer);

        uint256 gasBefore = gasleft();
        anchor.anchor(ROOT_1);
        uint256 gasUsed = gasBefore - gasleft();

        console.log("Gas used for anchor():", gasUsed);
        assertLt(gasUsed, 55000, "Anchor should use less than 55k gas");
    }

    function testGas_BatchAnchor() public {
        bytes32[] memory roots = new bytes32[](10);
        for (uint256 i = 0; i < 10; ++i) {
            roots[i] = keccak256(abi.encodePacked(i));
        }

        vm.prank(anchorer);

        uint256 gasBefore = gasleft();
        anchor.batchAnchor(roots);
        uint256 gasUsed = gasBefore - gasleft();

        console.log("Gas used for batchAnchor(10):", gasUsed);
        console.log("Gas per root:", gasUsed / 10);
        assertLt(
            gasUsed,
            400000,
            "Batch anchor of 10 should use less than 400k gas"
        );
    }

    function testGas_Verify() public {
        vm.prank(anchorer);
        anchor.anchor(ROOT_1);

        uint256 gasBefore = gasleft();
        anchor.getAnchorTimestamp(ROOT_1);
        uint256 gasUsed = gasBefore - gasleft();

        console.log("Gas used for getAnchorTimestamp():", gasUsed);
        assertLt(gasUsed, 3000, "Verify should use less than 3k gas (warm)");
    }

    // ============ VIEW FUNCTION TESTS ============

    function test_GetContractInfo() public view {
        (string memory version, uint256 chainId) = anchor.getContractInfo();
        assertEq(version, "1.0.0");
        assertEq(chainId, block.chainid);
    }

    // ============ TIME WARP TESTS ============

    function test_TimestampIsImmutable() public {
        vm.prank(anchorer);
        anchor.anchor(ROOT_1);
        uint256 originalTimestamp = anchor.anchors(ROOT_1);

        vm.warp(block.timestamp + 365 days);

        // Cannot re-anchor
        vm.prank(anchorer);
        vm.expectRevert();
        anchor.anchor(ROOT_1);

        // Timestamp unchanged
        assertEq(anchor.anchors(ROOT_1), originalTimestamp);
    }
}
