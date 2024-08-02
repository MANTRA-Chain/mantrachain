package types

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	// TODO: all hardcoded due to import cycle, see if there is better fix
	whitelistAddresses := []sdk.AccAddress{
		address.Module("coinfactory"),

		address.Module("liquidity", []byte("FeeCollector")),
		address.Module("liquidity", []byte("DustCollector")),
		address.Module("liquidity", []byte("GlobalEscrow")),
		address.Module("liquidity"),

		address.Module("lpfarm", []byte("FeeCollector")),
		address.Module("lpfarm", []byte("RewardsPool")),

		address.Module("farming", []byte("ecosystem_incentive_mm")),
		address.Module("marketmaker", []byte("ClaimableIncentiveReserveAcc")),
		sdk.AccAddress(crypto.AddressHash([]byte("marketmaker"))),
	}
	whitelistTransfersAccAddrs := make([]WhitelistTransfersAccAddrs, len(whitelistAddresses))
	for i, addr := range whitelistAddresses {
		whitelistTransfersAccAddrs[i] = WhitelistTransfersAccAddrs{
			Index:         addr,
			Account:       addr,
			IsWhitelisted: true,
		}
	}
	return &GenesisState{
		AccountPrivilegesList:      []AccountPrivileges{},
		GuardTransferCoins:         nil,
		RequiredPrivilegesList:     []RequiredPrivileges{},
		WhitelistTransfersAccAddrs: whitelistTransfersAccAddrs,
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		AccountPrivilegesList:      []AccountPrivileges{},
		GuardTransferCoins:         nil,
		RequiredPrivilegesList:     []RequiredPrivileges{},
		WhitelistTransfersAccAddrs: []WhitelistTransfersAccAddrs{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: params,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in accountPrivileges
	accountPrivilegesIndexMap := make(map[string]struct{})

	for _, elem := range gs.AccountPrivilegesList {
		index := string(elem.Account)
		if _, ok := accountPrivilegesIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for accountPrivileges")
		}
		accountPrivilegesIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in requiredPrivileges
	requiredPrivilegesIndexMap := make(map[string]struct{})

	for _, elem := range gs.RequiredPrivilegesList {
		var key []byte
		indexBytes := []byte(RequiredPrivilegesStoreKey([]byte(elem.Kind)))
		key = append(key, indexBytes...)
		key = append(key, Placeholder...)
		key = append(key, elem.Index...)
		key = append(key, Placeholder...)

		index := string(key)
		if _, ok := requiredPrivilegesIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for requiredPrivileges")
		}
		requiredPrivilegesIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in whitelistTransfersAccAddrs
	whitelistTransfersAccAddrsIndexMap := make(map[string]struct{})

	for _, elem := range gs.WhitelistTransfersAccAddrs {
		index := string(elem.Index)
		if _, ok := whitelistTransfersAccAddrsIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for whitelistTransfersAccAddrs")
		}
		whitelistTransfersAccAddrsIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
