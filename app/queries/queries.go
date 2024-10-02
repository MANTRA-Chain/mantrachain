package queries

import (
	"fmt"
	"sync"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	tokenfactorytypes "github.com/MANTRA-Chain/mantrachain/x/tokenfactory/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

// stargateResponsePools keeps a whitelist and its deterministic
// response binding for stargate queries.
// CONTRACT: Since results of queries go into blocks, queries being added here should always be
// deterministic or can cause non-determinism in the state machine.
//
// The query is multi-threaded, so we're using a sync.Pool
// to manage the allocation and de-allocation of newly created
// protobuf objects.
var stargateResponsePools = make(map[string]*sync.Pool)

func init() {
	// Whitelist specific queries
	setWhitelistedQuery("/osmosis.tokenfactory.v1beta1.Query/Params", &tokenfactorytypes.QueryParamsResponse{})
	setWhitelistedQuery("/osmosis.tokenfactory.v1beta1.Query/DenomAuthorityMetadata", &tokenfactorytypes.QueryDenomAuthorityMetadataResponse{})
}

// setWhitelistedQuery sets the whitelisted query at the provided path.
// This method also creates a sync.Pool for the provided proto message.
func setWhitelistedQuery(queryPath string, protoResponse proto.Message) {
	stargateResponsePools[queryPath] = &sync.Pool{
		New: func() any {
			return proto.Clone(protoResponse)
		},
	}
}

// IsWhitelistedQuery returns an error if the query is not whitelisted.
func IsWhitelistedQuery(queryPath string) error {
	_, isWhitelisted := stargateResponsePools[queryPath]
	if !isWhitelisted {
		return wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", queryPath)}
	}
	return nil
}

// getWhitelistedQuery returns the whitelisted query at the provided path.
// CONTRACT: Must call returnStargateResponseToPool to avoid pointless allocations.
func getWhitelistedQuery(queryPath string) (proto.Message, error) {
	pool, isWhitelisted := stargateResponsePools[queryPath]
	if !isWhitelisted {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", queryPath)}
	}
	protoMsg, ok := pool.Get().(proto.Message)
	if protoMsg == nil {
		return nil, fmt.Errorf("pool returned nil for query path '%s'", queryPath)
	}
	if !ok {
		return nil, fmt.Errorf("failed to assert type to proto.Message")
	}
	return protoMsg, nil
}

// returnStargateResponseToPool returns the provided proto message to the appropriate pool based on its query path.
func returnStargateResponseToPool(queryPath string, pb proto.Message) {
	pool, exists := stargateResponsePools[queryPath]
	if !exists {
		// Optional: handle the unexpected queryPath
		return
	}
	pool.Put(pb)
}

// StargateQuerier dispatches whitelisted stargate queries.
func StargateQuerier(queryRouter baseapp.GRPCQueryRouter, cdc codec.Codec) func(ctx sdk.Context, request *wasmvmtypes.StargateQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmvmtypes.StargateQuery) ([]byte, error) {
		// Check if the query path is whitelisted.
		if err := IsWhitelistedQuery(request.Path); err != nil {
			return nil, err
		}

		protoResponseType, err := getWhitelistedQuery(request.Path)
		if err != nil {
			return nil, err
		}
		// Ensure the proto message is returned to the pool to prevent leaks.
		defer returnStargateResponseToPool(request.Path, protoResponseType)

		route := queryRouter.Route(request.Path)
		if route == nil {
			return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("no route to query '%s'", request.Path)}
		}

		res, err := route(ctx, &abci.RequestQuery{
			Data: request.Data,
			Path: request.Path,
		})
		if err != nil {
			return nil, err
		}

		if res.Value == nil {
			return nil, fmt.Errorf("response value is nil")
		}

		// Unmarshal the response.
		err = cdc.Unmarshal(res.Value, protoResponseType)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response for query path '%s': %w", request.Path, err)
		}

		// Marshal to JSON for the contract.
		bz, err := cdc.MarshalJSON(protoResponseType)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response to JSON: %w", err)
		}

		return bz, nil
	}
}

// RegisterCustomPlugins registers custom plugins for the wasm module.
func RegisterCustomPlugins(
	queryRouter baseapp.GRPCQueryRouter,
	cdc codec.Codec,
) []wasmkeeper.Option {
	// Register the Stargate querier with the custom queries.
	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Stargate: StargateQuerier(queryRouter, cdc),
	})

	return []wasmkeeper.Option{
		queryPluginOpt,
	}
}
