# Mantrachain Modules Params

## Guard Module

The `x/guard` module contains the following params:

- admin_account
- account_privileges_token_collection_creator
- account_privileges_token_collection_id
- default_privileges

`AdminAccount` is the chain main admin account wallet address.

`AccountPrivilegesTokenCollectionCreator` is the collection creator of the guard soul-bond nft token collection.

`AccountPrivilegesTokenCollectionId` is the collection id of the guard soul-bond nft token collection.

`DefaultPrivileges` is the default privileges for the guard module.

The `guard` module params can be updated with governance.

## Farming Module

The `x/farming` module contains the following params:

- private_plan_creation_fee
- next_epoch_days
- farming_fee_collector
- delayed_staking_gas_fee
- max_num_private_plans

`PrivatePlanCreationFee` specifies the fee for plan creation this fee prevents from spamming and is collected in the community pool.

`NextEpochDays` is the epoch length in number of days it updates internal state called `CurrentEpochDays` that is used to process staking and reward distribution in end blocker.

`FarmingFeeCollector` is the module account address to collect fees within the farming module.

`DelayedStakingGasFee` is used to impose gas fee for the delayed staking.

`MaxNumPrivatePlans` is the maximum number of active private plans.

The `farming` module params can be updated with governance.

## TxFees Module

The `x/txfees` module contains the following params:

- base_denom

`BaseDenom` param, which is the native token(`uaum`) in which the gas fees are paid.

The `txfees` module params can be updated with governance.

## Coin Factory Module

The `x/coinfactory` module contains the following params:

- denom_creation_fee
- denom_creation_gas_consume

`DenomCreationFee` defines the fee to be charged on the creation of a new denom. The fee is drawn from the MsgCreateDenom's sender account, and transferred to the admin's wallet.

`DenomCreationGasConsume` defines the gas cost for creating a new denom. This is intended as a spam deterrence mechanism.

The `coinfactory` module params can be updated with governance.

## Farming Module

The `x/farming` module contains the following params:

- private_plan_creation_fee
- next_epoch_days
- farming_fee_collector
- delayed_staking_gas_fee
- max_num_private_plans

`PrivatePlanCreationFee` specifies the fee for plan creation this fee prevents from spamming and is collected in the community pool.

`NextEpochDays` is the epoch length in number of days it updates internal state called `CurrentEpochDays` that is used to process staking and reward distribution in end blocker.

`FarmingFeeCollector` is the module account address to collect fees within the farming module.

`DelayedStakingGasFee` is used to impose gas fee for the delayed staking.

`MaxNumPrivatePlans` is the maximum number of active private plans.

The `farming` module params can be updated with governance.

## Liquid Farming Module

The `x/liquidfarming` module contains the following params:

- fee_collector
- rewards_auction_duration
- liquid_farms

`FeeCollector` is the module account address to collect fees within the liquid farming module.

The `liquidfarming` module params can be updated with governance.

## Liquidity Module

The `x/liquidity` module contains the following params:

- batch_size
- tick_precision
- fee_collector_address
- dust_collector_address
- min_initial_pool_coin_supply
- pair_creation_fee
- pool_creation_fee
- min_initial_deposit_amount
- max_price_limit_ratio
- max_num_market_making_order_ticks
- max_order_lifespan
- swap_fee_rate
- withdraw_fee_rate
- deposit_extra_gas
- withdraw_extra_gas
- order_extra_gas
- max_num_active_pools_per_pair

`N/A` details about the params can be found atm.

The `liquidity` module params can be updated with governance.

## LPFrarm Module

The `x/lpfarm` module contains the following params:

- private_plan_creation_fee
- fee_collector
- max_num_private_plans
- max_block_duration

The `lpfarm` module params can be updated with governance.

## Market Maker Module

The `x/marketmaker` module contains the following params:

- incentive_budget_address
- deposit_amount
- common
- incentive_pairs

`IncentiveBudgetAddress` is the address containing the funds used to distribute incentives.

`deposit_amount` is the amount of deposit to be applied to the market maker, which is calculated per pair and is refunded when the market maker included or rejected through the `MarketMaker Proposal`.

`Common` is common variables used in market maker scoring system.

`IncentivePairs` include the pairs that are incentive target pairs and the variables used in market maker scoring system.

The `marketmaker` module params can be updated with governance.

## Token Module

The `x/token` module contains the following params:

- valid_nft_collection_id
- nft_collection_default_id
- nft_collection_default_name
- valid_nft_collection_metadata_symbol_min_length
- valid_nft_collection_metadata_symbol_max_length
- valid_nft_collection_metadata_description_max_length
- valid_nft_collection_metadata_name_max_length
- valid_nft_collection_metadata_images_max_count
- valid_nft_collection_metadata_images_type_max_length
- valid_nft_collection_metadata_links_max_count
- valid_nft_collection_metadata_links_type_max_length
- valid_nft_collection_metadata_options_max_count
- valid_nft_collection_metadata_options_type_max_length
- valid_nft_collection_metadata_options_value_max_length
- valid_nft_collection_metadata_options_sub_value_max_length
- valid_nft_id
- valid_nft_metadata_max_count
- valid_nft_metadata_title_max_length
- valid_nft_metadata_description_max_length
- valid_nft_metadata_images_max_count
- valid_nft_metadata_images_type_max_length
- valid_nft_metadata_links_max_count
- valid_nft_metadata_links_type_max_length
- valid_nft_metadata_attributes_max_count
- valid_nft_metadata_attributes_type_max_length
- valid_nft_metadata_attributes_value_max_length
- valid_nft_metadata_attributes_sub_value_max_length
- valid_burn_nft_max_count

The `token` module params contain varios restrictions being validated in the token collection/nft creation and/or update.

The `token` module params can be updated with governance.

## ***Note***: For more information about the Cosmos SDK native modules params, please refer to the [Cosmos SDK docs](https://docs.cosmos.network/v0.47)
