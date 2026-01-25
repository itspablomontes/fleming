package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var (
	ErrEmptyLeaves     = errors.New("audit: merkle tree requires at least one leaf")
	ErrLeafNotFound    = errors.New("audit: entry hash not found in merkle tree")
	ErrInvalidHash     = errors.New("audit: invalid hex hash")
	ErrInvalidTreeRoot = errors.New("audit: merkle tree root is empty")
)

type ProofStep struct {
	Hash   string
	IsLeft bool
}

type Proof struct {
	EntryHash string
	Steps     []ProofStep
}

type MerkleTree struct {
	Leaves []string
	Levels [][]string
	Root   string
}

func BuildMerkleTree(entries []Entry) (*MerkleTree, error) {
	if len(entries) == 0 {
		return nil, ErrEmptyLeaves
	}

	leaves := make([]string, 0, len(entries))
	for _, entry := range entries {
		hash := entry.Hash
		if hash == "" {
			hash = entry.ComputeHash()
		}
		if !isValidHexHash(hash) {
			return nil, ErrInvalidHash
		}
		leaves = append(leaves, hash)
	}

	levels, err := buildLevels(leaves)
	if err != nil {
		return nil, err
	}

	root := levels[len(levels)-1][0]
	return &MerkleTree{
		Leaves: leaves,
		Levels: levels,
		Root:   root,
	}, nil
}

func ComputeRoot(leaves []string) (string, error) {
	if len(leaves) == 0 {
		return "", ErrEmptyLeaves
	}

	levels, err := buildLevels(leaves)
	if err != nil {
		return "", err
	}

	return levels[len(levels)-1][0], nil
}

func GenerateProof(tree *MerkleTree, entryHash string) (*Proof, error) {
	if tree == nil || len(tree.Levels) == 0 {
		return nil, ErrInvalidTreeRoot
	}
	if !isValidHexHash(entryHash) {
		return nil, ErrInvalidHash
	}

	index := -1
	for i, leaf := range tree.Leaves {
		if leaf == entryHash {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, ErrLeafNotFound
	}

	steps := make([]ProofStep, 0, len(tree.Levels)-1)
	for levelIndex := 0; levelIndex < len(tree.Levels)-1; levelIndex++ {
		level := tree.Levels[levelIndex]
		isRight := index%2 == 1
		var siblingIndex int
		if isRight {
			siblingIndex = index - 1
		} else {
			siblingIndex = index + 1
		}

		if siblingIndex >= len(level) {
			siblingIndex = index
		}

		steps = append(steps, ProofStep{
			Hash:   level[siblingIndex],
			IsLeft: isRight,
		})

		index = index / 2
	}

	return &Proof{
		EntryHash: entryHash,
		Steps:     steps,
	}, nil
}

func VerifyProof(root string, entryHash string, proof *Proof) bool {
	if proof == nil || root == "" || !isValidHexHash(entryHash) {
		return false
	}

	hash := entryHash
	for _, step := range proof.Steps {
		if !isValidHexHash(step.Hash) {
			return false
		}

		var err error
		if step.IsLeft {
			hash, err = hashPair(step.Hash, hash)
		} else {
			hash, err = hashPair(hash, step.Hash)
		}
		if err != nil {
			return false
		}
	}

	return hash == root
}

func buildLevels(leaves []string) ([][]string, error) {
	if len(leaves) == 0 {
		return nil, ErrEmptyLeaves
	}

	level := make([]string, len(leaves))
	copy(level, leaves)
	for _, leaf := range level {
		if !isValidHexHash(leaf) {
			return nil, ErrInvalidHash
		}
	}

	levels := [][]string{level}
	for len(level) > 1 {
		if len(level)%2 == 1 {
			level = append(level, level[len(level)-1])
		}

		next := make([]string, 0, len(level)/2)
		for i := 0; i < len(level); i += 2 {
			parent, err := hashPair(level[i], level[i+1])
			if err != nil {
				return nil, err
			}
			next = append(next, parent)
		}
		levels = append(levels, next)
		level = next
	}

	return levels, nil
}

func hashPair(left, right string) (string, error) {
	leftBytes, err := hex.DecodeString(left)
	if err != nil {
		return "", ErrInvalidHash
	}
	rightBytes, err := hex.DecodeString(right)
	if err != nil {
		return "", ErrInvalidHash
	}
	sum := sha256.Sum256(append(leftBytes, rightBytes...))
	return hex.EncodeToString(sum[:]), nil
}

func isValidHexHash(value string) bool {
	if len(value) == 0 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
