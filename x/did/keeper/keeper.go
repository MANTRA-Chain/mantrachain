package keeper

import (
	"fmt"
	"strings"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/did/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		guardKeeper types.GuardKeeper
	}
)

// UnmarshalFn is a generic function to unmarshal bytes
type UnmarshalFn func(value []byte) (interface{}, bool)

// MarshalFn is a generic function to marshal bytes
type MarshalFn func(value interface{}) []byte

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

	guardKeeper types.GuardKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,

		guardKeeper: guardKeeper,
	}
}

func SetGuardKeeper(k *Keeper, guardKeeper types.GuardKeeper) {
	k.guardKeeper = guardKeeper
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetDidDocument(ctx sdk.Context, key []byte, document types.DidDocument) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	documentBs, err := document.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(append(types.DidDocumentKey, key...), documentBs)
}

func (k Keeper) HasDidDocument(ctx sdk.Context, key []byte) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(append(types.DidDocumentKey, key...))
}

func (k Keeper) DeleteDidDocument(ctx sdk.Context, key []byte) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(append(types.DidDocumentKey, key...))
}

func (k Keeper) GetDidDocument(ctx sdk.Context, key []byte) (types.DidDocument, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	documentBs := store.Get(append(types.DidDocumentKey, key...))
	document := types.DidDocument{}
	if documentBs == nil {
		return document, false
	}
	if err := document.Unmarshal(documentBs); err != nil {
		panic(err)
	}
	return document, true
}

func (k Keeper) SetDidMetadata(ctx sdk.Context, key []byte, meta types.DidMetadata) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	metaBs, err := meta.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(append(types.DidMetadataKey, key...), metaBs)
}

func (k Keeper) GetDidMetadata(ctx sdk.Context, key []byte) (types.DidMetadata, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	metaBs := store.Get(append(types.DidMetadataKey, key...))
	meta := types.DidMetadata{}
	if metaBs == nil {
		return meta, false
	}
	if err := meta.Unmarshal(metaBs); err != nil {
		panic(err)
	}
	return meta, true
}

// ResolveDid returning the did document and associated metadata
func (k Keeper) ResolveDid(ctx sdk.Context, did types.DID) (doc types.DidDocument, meta types.DidMetadata, err error) {
	if strings.HasPrefix(did.String(), types.DidKeyPrefix) {
		doc, meta, err = types.ResolveAccountDID(did.String(), ctx.ChainID())
		return
	}
	doc, found := k.GetDidDocument(ctx, []byte(did.String()))
	if !found {
		err = types.ErrDidDocumentNotFound
		return
	}
	meta, _ = k.GetDidMetadata(ctx, []byte(did.String()))
	return
}

func (k Keeper) Marshal(value interface{}) (bytes []byte) {
	switch value := value.(type) {
	case types.DidDocument:
		bytes = k.cdc.MustMarshal(&value)
	case types.DidMetadata:
		bytes = k.cdc.MustMarshal(&value)
	}
	return
}

// GetAllDidDocumentsWithCondition retrieve a list of
// did document by some arbitrary criteria. The selector filter has access
// to both the did and its metadata
func (k Keeper) GetAllDidDocumentsWithCondition(
	ctx sdk.Context,
	key []byte,
	didSelector func(did types.DidDocument) bool,
) (didDocs []types.DidDocument) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, key)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		document := types.DidDocument{}
		if err := document.Unmarshal(iterator.Value()); err != nil {
			panic(err)
		}
		if didSelector(document) {
			didDocs = append(didDocs, document)
		}
	}

	return didDocs
}

// GetAllDidDocuments returns all the DidDocuments
func (k Keeper) GetAllDidDocuments(ctx sdk.Context) []types.DidDocument {
	return k.GetAllDidDocumentsWithCondition(
		ctx,
		types.DidDocumentKey,
		func(did types.DidDocument) bool { return true },
	)
}

// GetDidDocumentsByPubKey retrieve a did document using a pubkey associated to the DID
// TODO: this function is used only in the issuer module ante handler !
func (k Keeper) GetDidDocumentsByPubKey(ctx sdk.Context, pubkey cryptotypes.PubKey) (dids []types.DidDocument) {
	dids = k.GetAllDidDocumentsWithCondition(
		ctx,
		types.DidDocumentKey,
		func(did types.DidDocument) bool {
			return did.HasPublicKey(pubkey)
		},
	)
	// compute the key did

	// generate the address
	addr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), pubkey.Address())
	if err != nil {
		return
	}
	doc, _, err := types.ResolveAccountDID(types.NewKeyDID(addr).String(), ctx.ChainID())
	if err != nil {
		return
	}
	dids = append(dids, doc)
	return
}
