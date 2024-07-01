package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/errors"

	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
)

// hash concatenates and hashes two byte slices.
func hash(left, right []byte) ([]byte, error) {
	hasher := sha256.New()
	if len(left) != 32 || len(right) != 32 {
		return nil, errors.Wrap(types.ErrInvalidMerklePath, "invalid merkle path length")
	}
	_, err := hasher.Write(left)
	if err != nil {
		return nil, err
	}
	_, err = hasher.Write(right)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

// verifyMerklePath verifies the Merkle path for a given leaf and returns true if the path is valid.
func verifyMerklePath(leafHash []byte, merklePath [][]byte, rootHash []byte, leafIndex uint64) (bool, error) {
	var err error
	if len(merklePath) == 0 {
		return false, errors.Wrap(types.ErrInvalidMerklePath, "empty merkle path")
	}

	calculatedHash := leafHash
	for i, pathHash := range merklePath {
		// Determine the direction (left or right) for this level using the ith bit of leafIndex
		direction := (leafIndex >> i) & 1

		if direction == 0 { // If the bit is 0, the current node is a left child
			calculatedHash, err = hash(calculatedHash, pathHash)
		} else { // If the bit is 1, the current node is a right child
			calculatedHash, err = hash(pathHash, calculatedHash)
		}
		if err != nil {
			return false, err
		}
	}

	// Compare the calculated root hash with the given root hash
	return hex.EncodeToString(calculatedHash) == hex.EncodeToString(rootHash), nil
}

func pathToChunks(path []byte) ([][]byte, error) {
	chunkSize := 32
	if len(path)%chunkSize != 0 {
		return nil, errors.Wrap(types.ErrInvalidMerklePath, "invalid merkle path length")
	}
	merklePath := make([][]byte, 0)
	for i := 0; i < len(path); i += chunkSize {
		end := i + chunkSize

		// Make sure not to exceed the slice bounds
		if end > len(path) {
			return nil, errors.Wrap(types.ErrInvalidMerklePath, "invalid merkle path length")
		}

		merklePath = append(merklePath, path[i:end])
	}

	return merklePath, nil
}

func genLeafHash(creator string, amt string) ([]byte, error) {
	leafStr := fmt.Sprintf("%s-%s", creator, amt)
	hasher := sha256.New()
	hasher.Write([]byte(leafStr))
	return hasher.Sum(nil), nil
}
