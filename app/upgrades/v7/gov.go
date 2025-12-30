package v7

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func migrateGov(ctx sdk.Context, govKeeper govkeeper.Keeper) error {
	govParams, err := govKeeper.Params.Get(ctx)
	if err != nil {
		return err
	}

	// migrate min_deposit, min_initial_deposit, expedited_min_deposit
	govParams.MinDeposit = convertCoinsToNewDenom(govParams.MinDeposit)
	govParams.ExpeditedMinDeposit = convertCoinsToNewDenom(govParams.ExpeditedMinDeposit)

	if err := govKeeper.Params.Set(ctx, govParams); err != nil {
		return err
	}

	// migrate deposits on active proposals
	iter, err := govKeeper.Proposals.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		kv, err := iter.KeyValue()
		if err != nil {
			return err
		}

		proposal := kv.Value

		if proposal.Status == govv1.StatusDepositPeriod || proposal.Status == govv1.StatusVotingPeriod {
			deposits, err := govKeeper.GetDeposits(ctx, proposal.Id)
			if err != nil {
				return err
			}

			for _, deposit := range deposits {
				deposit.Amount = convertCoinsToNewDenom(deposit.Amount)
				if err = govKeeper.SetDeposit(ctx, *deposit); err != nil {
					return err
				}
			}

			proposal.TotalDeposit = convertCoinsToNewDenom(proposal.TotalDeposit)
			if err = govKeeper.SetProposal(ctx, proposal); err != nil {
				return err
			}
		}
	}

	return nil
}
