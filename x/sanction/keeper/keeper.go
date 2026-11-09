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
		cdc           codec.BinaryCodec
		storeService  store.KVStoreService
		logger        log.Logger
		accountKeeper types.AccountKeeper

		// the address capable of executing a MsgUpdateParams message.
		// Typically, this should be the x/gov module account.
		authority string

		// protectedAddrs can never be sanctioned by a sanctioner: the
		// authority and the app's module accounts, including ones that don't
		// exist in x/auth yet.
		protectedAddrs map[string]struct{}

		Schema collections.Schema
		Params collections.Item[types.Params]

		// BlacklistAccounts is a set of accounts that are blacklisted from using the module.
		BlacklistAccounts collections.KeySet[string]
		// Sanctioners maps a sanctioner address to its entry.
		Sanctioners collections.Map[string, types.Sanctioner]
		// SanctionRecords maps a blacklisted address to who added it, why and
		// when. Entries that predate sanctioners have no record.
		SanctionRecords collections.Map[string, types.SanctionRecord]
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	accountKeeper types.AccountKeeper,
	moduleAddrs []string,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	protected := make(map[string]struct{}, len(moduleAddrs)+1)
	protected[authority] = struct{}{}
	for _, addr := range moduleAddrs {
		protected[addr] = struct{}{}
	}

	k := Keeper{
		cdc:            cdc,
		storeService:   storeService,
		authority:      authority,
		logger:         logger,
		accountKeeper:  accountKeeper,
		protectedAddrs: protected,

		Params:            collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		BlacklistAccounts: collections.NewKeySet(sb, types.PrefixBlackListAccounts, "blacklist_accounts", collections.StringKey),
		Sanctioners:       collections.NewMap(sb, types.PrefixSanctioners, "sanctioners", collections.StringKey, codec.CollValue[types.Sanctioner](cdc)),
		SanctionRecords:   collections.NewMap(sb, types.PrefixSanctionRecords, "sanction_records", collections.StringKey, codec.CollValue[types.SanctionRecord](cdc)),
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
