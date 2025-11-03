package v7rc0

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

const (
	NameSlot   = 0
	SymbolSlot = 1
)

var WOMContractAddress = map[string][]common.Address{
	"mantra-1":            {common.HexToAddress("0xE3047710EF6cB36Bcf1E58145529778eA7Cb5598")},
	"mantra-dukong-1":     {common.HexToAddress("0x10d26F0491fA11c5853ED7C1f9817b098317DC46")},
	"mantra-canary-net-1": {common.HexToAddress("0x523A024258fc56E4d6d79D4367a98F2548A9f401")},
	"mantra-dryrun-1":     {common.HexToAddress("0x523A024258fc56E4d6d79D4367a98F2548A9f401")},
}

func migrateWOMs(ctx sdk.Context, evmKeeper evmkeeper.Keeper) error {
	addresses := WOMContractAddress[ctx.ChainID()]
	for _, addr := range addresses {
		if err := migrateWOM(ctx, evmKeeper, addr); err != nil {
			return err
		}
	}
	return nil
}

func migrateWOM(ctx sdk.Context, evmKeeper evmkeeper.Keeper, contract common.Address) error {
	// build excluded slots for name symbol fields
	excluded := make(map[common.Hash]struct{})

	// assume 0 and 1 are slots for name and symbol strings.
	for _, slot := range []int{NameSlot, SymbolSlot} {
		slots, err := stringSlots(ctx, evmKeeper, contract, slot)
		if err != nil {
			return err
		}

		for _, s := range slots {
			excluded[s] = struct{}{}
		}
	}

	setStringField(ctx, evmKeeper, contract, NameSlot, "WMANTRA Token")
	setStringField(ctx, evmKeeper, contract, SymbolSlot, "WMANTRA")

	evmKeeper.ForEachStorage(ctx, contract, func(key, value common.Hash) bool {
		if _, ok := excluded[key]; ok {
			return true
		}

		// assume to be balance or allowance slots, multiply value by 4
		val := uint256.NewInt(0).SetBytes(value.Bytes())

		var overflow bool
		val, overflow = val.MulOverflow(val, uint256.NewInt(4))

		if overflow {
			// set to max uint256 when overflow
			value = common.MaxHash
		} else {
			value = common.BigToHash(val.ToBig())
		}

		evmKeeper.SetState(ctx, contract, key, value.Bytes())
		return true
	})
	return nil
}

// stringSlots computes all the storage slots related to a string field located at given slot
func stringSlots(ctx sdk.Context, evmKeeper evmkeeper.Keeper, contract common.Address, slot int) ([]common.Hash, error) {
	var slots []common.Hash

	// get length slot
	lengthSlot := common.BigToHash(big.NewInt(int64(slot)))
	slots = append(slots, lengthSlot)

	// get length value
	lengthValue := evmKeeper.GetState(ctx, contract, lengthSlot)
	lengthBig := new(big.Int).SetBytes(lengthValue.Bytes())
	if !lengthBig.IsUint64() {
		return nil, errors.New("string length exceeds uint64")
	}
	length := lengthBig.Uint64()
	if length == 0 {
		return slots, nil
	}

	if length > 1024 {
		return nil, errors.New("string length too large")
	}

	// padded length to 32 bytes slots
	numDataSlots := (length + 31) / 32

	// compute data slots
	dataStart := new(big.Int).SetBytes(crypto.Keccak256(lengthSlot.Bytes()))
	for i := uint64(0); i < numDataSlots; i++ {
		dataSlot := common.BigToHash(new(big.Int).Add(dataStart, big.NewInt(int64(i))))
		slots = append(slots, dataSlot)
	}

	return slots, nil
}

func setStringField(ctx sdk.Context, evmKeeper evmkeeper.Keeper, contract common.Address, slot int, value string) {
	if len(value) > 32 {
		panic("string length exceeds 32 bytes")
	}
	lengthSlot := common.BigToHash(big.NewInt(int64(slot)))
	length := common.BigToHash(big.NewInt(int64(len(value))))
	evmKeeper.SetState(ctx, contract, lengthSlot, length.Bytes())

	dataSlot := crypto.Keccak256Hash(lengthSlot.Bytes())
	var data common.Hash
	copy(data[:], []byte(value))
	evmKeeper.SetState(ctx, contract, dataSlot, data.Bytes())
}
