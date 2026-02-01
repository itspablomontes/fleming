// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

/// @title Bitmap
/// @notice Gas-efficient bitmap operations for storing boolean flags
/// @dev Stores 256 flags per uint256 slot, achieving ~98% storage efficiency
///      compared to mapping(uint256 => bool)
library Bitmap {
    /// @notice Thrown when attempting to set an index that exceeds bitmap capacity
    /// @param index The requested index
    /// @param maxCapacity The maximum supported index
    error IndexOutOfBounds(uint256 index, uint256 maxCapacity);

    /// @notice Maximum index supported (2^256 / 256 - 1, practically unlimited)
    uint256 constant MAX_INDEX = type(uint256).max >> 8;

    /// @notice Set a bit at the given index to true
    /// @param bitmap The storage array to modify
    /// @param index The bit index to set
    function set(uint256[] storage bitmap, uint256 index) internal {
        if (index > MAX_INDEX) revert IndexOutOfBounds(index, MAX_INDEX);

        uint256 slot = index >> 8; // index / 256
        uint256 offset = index & 0xFF; // index % 256

        // Ensure the slot exists
        if (slot >= bitmap.length) {
            // Extend the array with zeros
            uint256 slotsNeeded = slot - bitmap.length + 1;
            for (uint256 i = 0; i < slotsNeeded; ++i) {
                bitmap.push();
            }
        }

        // Set the bit using OR
        bitmap[slot] |= (uint256(1) << offset);
    }

    /// @notice Get the value of a bit at the given index
    /// @param bitmap The storage array to read from
    /// @param index The bit index to read
    /// @return True if the bit is set, false otherwise
    function get(
        uint256[] storage bitmap,
        uint256 index
    ) internal view returns (bool) {
        if (index > MAX_INDEX) return false;

        uint256 slot = index >> 8;

        // If slot doesn't exist, bit is not set
        if (slot >= bitmap.length) return false;

        uint256 offset = index & 0xFF;

        // Check if bit is set using AND
        return (bitmap[slot] >> offset) & 1 == 1;
    }

    /// @notice Unset (clear) a bit at the given index
    /// @param bitmap The storage array to modify
    /// @param index The bit index to clear
    /// @dev Note: Clearing doesn't shrink the array, just sets bit to 0
    function unset(uint256[] storage bitmap, uint256 index) internal {
        if (index > MAX_INDEX) revert IndexOutOfBounds(index, MAX_INDEX);

        uint256 slot = index >> 8;

        // If slot doesn't exist, nothing to unset
        if (slot >= bitmap.length) return;

        uint256 offset = index & 0xFF;

        // Clear the bit using AND with inverted mask
        bitmap[slot] &= ~(uint256(1) << offset);
    }

    /// @notice Count the number of set bits (population count) up to a given index
    /// @param bitmap The storage array to read from
    /// @param upToIndex Count bits up to and including this index
    /// @return count The number of set bits
    /// @dev Gas intensive for large bitmaps - use off-chain indexing instead
    function countSetBits(
        uint256[] storage bitmap,
        uint256 upToIndex
    ) internal view returns (uint256 count) {
        if (bitmap.length == 0) return 0;

        uint256 lastSlot = upToIndex >> 8;
        uint256 lastOffset = upToIndex & 0xFF;

        // Cap at actual bitmap length
        if (lastSlot >= bitmap.length) {
            lastSlot = bitmap.length - 1;
            lastOffset = 255;
        }

        // Count full slots
        for (uint256 slot = 0; slot < lastSlot; ++slot) {
            count += _popcount(bitmap[slot]);
        }

        // Count partial last slot
        if (lastOffset < 255) {
            uint256 mask = (uint256(1) << (lastOffset + 1)) - 1;
            count += _popcount(bitmap[lastSlot] & mask);
        } else {
            count += _popcount(bitmap[lastSlot]);
        }

        return count;
    }

    /// @notice Population count (number of 1 bits) in a uint256
    /// @param x The value to count
    /// @return The number of 1 bits
    /// @dev Uses the SWAR (SIMD Within A Register) algorithm
    function _popcount(uint256 x) internal pure returns (uint256) {
        x =
            (x &
                0x5555555555555555555555555555555555555555555555555555555555555555) +
            ((x >> 1) &
                0x5555555555555555555555555555555555555555555555555555555555555555);
        x =
            (x &
                0x3333333333333333333333333333333333333333333333333333333333333333) +
            ((x >> 2) &
                0x3333333333333333333333333333333333333333333333333333333333333333);
        x =
            (x &
                0x0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F) +
            ((x >> 4) &
                0x0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F0F);
        x =
            (x &
                0x00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF) +
            ((x >> 8) &
                0x00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF);
        x =
            (x &
                0x0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF) +
            ((x >> 16) &
                0x0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF);
        x =
            (x &
                0x00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF) +
            ((x >> 32) &
                0x00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF);
        x =
            (x &
                0x0000000000000000FFFFFFFFFFFFFFFF0000000000000000FFFFFFFFFFFFFFFF) +
            ((x >> 64) &
                0x0000000000000000FFFFFFFFFFFFFFFF0000000000000000FFFFFFFFFFFFFFFF);
        x =
            (x &
                0x00000000000000000000000000000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF) +
            ((x >> 128) &
                0x00000000000000000000000000000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF);
        return x;
    }
}
