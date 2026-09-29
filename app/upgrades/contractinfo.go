package upgrades

import (
	"fmt"
	"strings"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/encoding/protowire"
)

// ContractInfo field numbers in the linked wasmd; the retracted v0.61.0-v0.61.4
// had them swapped (CosmWasm/wasmd#2390).
const (
	contractInfoExtensionField  protowire.Number = 7
	contractInfoIBC2PortIDField protowire.Number = 8
)

// MigrateContractInfo rewrites ContractInfo records written by the retracted
// wasmd versions, which put ibc2_port_id at field 7 where this build expects
// extension. Records that already decode are skipped, so it is idempotent; a
// record that cannot be repaired aborts before anything is written.
func MigrateContractInfo(ctx sdk.Context, wasmStoreKey *storetypes.KVStoreKey) (int, error) {
	store := prefix.NewStore(ctx.KVStore(wasmStoreKey), wasmtypes.ContractKeyPrefix)

	fixed := map[string][]byte{}
	iter := store.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		if new(wasmtypes.ContractInfo).Unmarshal(iter.Value()) == nil {
			continue
		}
		bz, err := repairContractInfo(iter.Value())
		if err != nil {
			return 0, fmt.Errorf("contract %s: %w", sdk.AccAddress(iter.Key()), err)
		}
		fixed[string(iter.Key())] = bz
	}

	// Write after the scan so a failure leaves the store untouched. Every key
	// already exists, so map order cannot change the resulting tree.
	for key, bz := range fixed {
		store.Set([]byte(key), bz)
	}
	return len(fixed), nil
}

// repairContractInfo swaps fields 7 and 8 and checks the result decodes with an
// IBC2 port, so a record broken for any other reason is not rewritten.
func repairContractInfo(bz []byte) ([]byte, error) {
	fixed, err := swapProtoFields(bz, contractInfoExtensionField, contractInfoIBC2PortIDField)
	if err != nil {
		return nil, err
	}
	var info wasmtypes.ContractInfo
	if err := info.Unmarshal(fixed); err != nil {
		return nil, fmt.Errorf("still undecodable after swapping fields 7 and 8: %w", err)
	}
	if !strings.HasPrefix(info.IBC2PortID, wasmkeeper.PortIDPrefixV2) {
		return nil, fmt.Errorf("field 7 held %q, not an IBC2 port", info.IBC2PortID)
	}
	return fixed, nil
}

// swapProtoFields exchanges field numbers a and b, copying every value verbatim.
func swapProtoFields(bz []byte, a, b protowire.Number) ([]byte, error) {
	out := make([]byte, 0, len(bz))
	for len(bz) > 0 {
		num, typ, n := protowire.ConsumeTag(bz)
		if n < 0 {
			return nil, protowire.ParseError(n)
		}
		bz = bz[n:]

		// A group's end tag lives inside the value and would not be renumbered.
		if typ == protowire.StartGroupType {
			return nil, fmt.Errorf("field %d: groups are not supported", num)
		}

		n = protowire.ConsumeFieldValue(num, typ, bz)
		if n < 0 {
			return nil, protowire.ParseError(n)
		}
		value := bz[:n]
		bz = bz[n:]

		switch num {
		case a:
			num = b
		case b:
			num = a
		}
		out = protowire.AppendTag(out, num, typ)
		out = append(out, value...)
	}
	return out, nil
}
