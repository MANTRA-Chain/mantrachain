package wasmbinding

import (
	"fmt"
	"sync"

	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"

	"github.com/cosmos/gogoproto/proto"

	tokenfactorytypes "github.com/MANTRA-Chain/mantrachain/x/tokenfactory/types"
)

// stargateResponsePools keeps whitelist and its deterministic
// response binding for stargate queries.
// CONTRACT: since results of queries go into blocks, queries being added here should always be
// deterministic or can cause non-determinism in the state machine.
//
// The query is multi-threaded so we're using a sync.Pool
// to manage the allocation and de-allocation of newly created
// pb objects.
var stargateResponsePools = make(map[string]*sync.Pool)

// Note: When adding a migration here, we should also add it to the Async ICQ params in the upgrade.
// In the future we may want to find a better way to keep these in sync

//nolint:staticcheck
func init() {

	// tokenfactory
	setWhitelistedQuery("/osmosis.tokenfactory.v1beta1.Query/Params", &tokenfactorytypes.QueryParamsResponse{})
	setWhitelistedQuery("/osmosis.tokenfactory.v1beta1.Query/DenomAuthorityMetadata", &tokenfactorytypes.QueryDenomAuthorityMetadataResponse{})
	// Does not include denoms_from_creator, TBD if this is the index we want contracts to use instead of admin

}

// IsWhitelistedQuery returns if the query is not whitelisted.
func IsWhitelistedQuery(queryPath string) error {
	_, isWhitelisted := stargateResponsePools[queryPath]
	if !isWhitelisted {
		return wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", queryPath)}
	}
	return nil
}

// getWhitelistedQuery returns the whitelisted query at the provided path.
// If the query does not exist, or it was setup wrong by the chain, this returns an error.
// CONTRACT: must call returnStargateResponseToPool in order to avoid pointless allocs.
func getWhitelistedQuery(queryPath string) (proto.Message, error) {
	protoResponseAny, isWhitelisted := stargateResponsePools[queryPath]
	if !isWhitelisted {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", queryPath)}
	}
	protoMarshaler, ok := protoResponseAny.Get().(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to assert type to proto.Messager")
	}
	return protoMarshaler, nil
}

type protoTypeG[T any] interface {
	*T
	proto.Message
}

// setWhitelistedQuery sets the whitelisted query at the provided path.
// This method also creates a sync.Pool for the provided protoMarshaler.
// We use generics so we can properly instantiate an object that the
// queryPath expects as a response.
func setWhitelistedQuery[T any, PT protoTypeG[T]](queryPath string, _ PT) {
	stargateResponsePools[queryPath] = &sync.Pool{
		New: func() any {
			return PT(new(T))
		},
	}
}

// returnStargateResponseToPool returns the provided protoMarshaler to the appropriate pool based on it's query path.
func returnStargateResponseToPool(queryPath string, pb proto.Message) {
	stargateResponsePools[queryPath].Put(pb)
}

func GetStargateWhitelistedPaths() (keys []string) {
	// Iterate over the map and collect the keys
	keys = make([]string, 0, len(stargateResponsePools))
	for k := range stargateResponsePools {
		keys = append(keys, k)
	}
	return keys
}
