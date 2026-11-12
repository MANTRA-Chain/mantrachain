package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		// evmKeeper clears EIP-7702 delegations from sanctioned accounts. It
		// may be nil, in which case nothing is cleared.
		evmKeeper types.EVMKeeper

		// the address capable of executing a MsgUpdateParams message.
		// Typically, this should be the x/gov module account.
		authority string

		Schema collections.Schema
		Params collections.Item[types.Params]

		// BlacklistAccounts is a set of accounts that are blacklisted from using the module.
		BlacklistAccounts collections.KeySet[string]
		// Sanctioners maps a sanctioner address to its entry.
		Sanctioners collections.Map[string, types.Sanctioner]
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	evmKeeper types.EVMKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
		evmKeeper:    evmKeeper,

		Params:            collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		BlacklistAccounts: collections.NewKeySet(sb, types.PrefixBlackListAccounts, "blacklist_accounts", collections.StringKey),
		Sanctioners:       collections.NewMap(sb, types.PrefixSanctioners, "sanctioners", collections.StringKey, codec.CollValue[types.Sanctioner](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
