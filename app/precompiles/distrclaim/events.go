package distrclaim

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const claimRewardsAndConvertCoinEventSig = "ClaimRewardsAndConvertCoin(address,string,uint256)"

var claimRewardsAndConvertCoinEventID = crypto.Keccak256Hash([]byte(claimRewardsAndConvertCoinEventSig))

func topicAddress(addr common.Address) common.Hash {
	return common.BytesToHash(common.LeftPadBytes(addr.Bytes(), 32))
}

func (p *Precompile) emitClaimRewardsAndConvertCoinEvent(
	ctx sdk.Context,
	evm *vm.EVM,
	delegator common.Address,
	denom string,
	amount *big.Int,
) error {
	if evm == nil || evm.StateDB == nil {
		return nil
	}

	stringT, err := abi.NewType("string", "", nil)
	if err != nil {
		return err
	}
	uint256T, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return err
	}

	arguments := abi.Arguments{{Type: stringT}, {Type: uint256T}}
	packed, err := arguments.Pack(denom, amount)
	if err != nil {
		return err
	}

	evm.StateDB.AddLog(&types.Log{
		Address:     p.ContractAddress,
		Topics:      []common.Hash{claimRewardsAndConvertCoinEventID, topicAddress(delegator)},
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}
