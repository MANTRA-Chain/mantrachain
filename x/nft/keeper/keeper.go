package keeper

import (
	"github.com/LimeChain/mantrachain/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

// Keeper of the nft store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	bk       types.BankKeeper
}

// NewKeeper creates a new nft Keeper instance
func NewKeeper(key storetypes.StoreKey,
	cdc codec.BinaryCodec, ak types.AccountKeeper, bk types.BankKeeper,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: key,
		bk:       bk,
	}
}
