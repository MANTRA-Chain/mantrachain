package distrclaim

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/cosmos/evm/precompiles/common"
	erc20types "github.com/cosmos/evm/x/erc20/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func (p *Precompile) runClaimRewardsAndConvertCoin(
	ctx sdk.Context,
	input []byte,
	stateDB vm.StateDB,
	contract *vm.Contract,
) ([]byte, error) {
	var in ClaimRewardsAndConvertCoinCall
	if _, err := (&in).Decode(input); err != nil {
		return nil, err
	}

	delegatorAddr := in.Delegator
	maxRetrieve := in.MaxRetrieve
	denom := in.Denom

	if delegatorAddr == (common.Address{}) {
		return nil, fmt.Errorf(cmn.ErrInvalidDelegator, delegatorAddr)
	}
	if denom == "" {
		return nil, fmt.Errorf(cmn.ErrInvalidType, "denom", "string", denom)
	}

	out, err := p.claimRewardsAndConvertCoin(ctx, delegatorAddr, maxRetrieve, denom, stateDB, contract)
	if err != nil {
		return nil, err
	}

	return out.Encode()
}

func (p *Precompile) claimRewardsAndConvertCoin(
	ctx sdk.Context,
	delegatorAddr common.Address,
	maxRetrieve uint32,
	denom string,
	_ vm.StateDB,
	contract *vm.Contract,
) (*ClaimRewardsAndConvertCoinReturn, error) {

	maxVals, err := p.stakingKeeper.MaxValidators(ctx)
	if err != nil {
		return nil, err
	}
	if maxRetrieve > maxVals {
		return nil, fmt.Errorf(
			"maxRetrieve (%d) parameter exceeds the maximum number of validators (%d)",
			maxRetrieve,
			maxVals,
		)
	}

	msgSender := contract.Caller()
	if msgSender != delegatorAddr {
		return nil, fmt.Errorf(cmn.ErrRequesterIsNotMsgSender, msgSender.String(), delegatorAddr.String())
	}

	delegatorBech32, err := p.addrCdc.BytesToString(delegatorAddr.Bytes())
	if err != nil {
		return nil, err
	}

	// Ensure rewards are paid to the delegator account itself, so the subsequent
	// conversion burns from the same account.
	_, err = p.distributionMsgServer.SetWithdrawAddress(ctx, &distrtypes.MsgSetWithdrawAddress{
		DelegatorAddress: delegatorBech32,
		WithdrawAddress:  delegatorBech32,
	})
	if err != nil {
		return nil, err
	}

	res, err := p.stakingKeeper.GetDelegatorValidators(ctx, delegatorAddr.Bytes(), maxRetrieve)
	if err != nil {
		return nil, err
	}

	totalCoins := sdk.Coins{}
	for _, validator := range res.Validators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			return nil, err
		}

		coins, err := p.distributionKeeper.WithdrawDelegationRewards(ctx, delegatorAddr.Bytes(), valAddr)
		if err != nil {
			return nil, err
		}

		totalCoins = totalCoins.Add(coins...)
	}

	converted := big.NewInt(0)
	amount := totalCoins.AmountOf(denom)
	if amount.IsPositive() {
		amtBig := amount.BigInt()
		_, err := p.erc20MsgServer.ConvertCoin(ctx, &erc20types.MsgConvertCoin{
			Coin: sdk.Coin{
				Denom:  denom,
				Amount: amount,
			},
			Receiver: delegatorAddr.Hex(),
			Sender:   delegatorBech32,
		})
		if err != nil {
			return nil, err
		}
		converted = amtBig
	}

	return &ClaimRewardsAndConvertCoinReturn{Amount: converted}, nil
}
