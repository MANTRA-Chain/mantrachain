# Mantrachain Modules Params

## Guard Module

The `x/guard` module contains the following params:

- admin_account
- account_privileges_token_collection_creator
- account_privileges_token_collection_id
- default_privileges
- base_denom

`AdminAccount` is the chain main admin account wallet address.

`AccountPrivilegesTokenCollectionCreator` is the collection creator of the guard soul-bond nft token collection.

`AccountPrivilegesTokenCollectionId` is the collection id of the guard soul-bond nft token collection.

`DefaultPrivileges` is the default privileges for the guard module.

`BaseDenom` is the native token(`uom`) in which the gas fees are paid.

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

`BaseDenom` param, which is the native token(`uom`) in which the gas fees are paid.

The `txfees` module params can be updated with governance.

## Coin Factory Module

The `x/coinfactory` module contains the following params:

- denom_creation_fee
- denom_creation_gas_consume

`DenomCreationFee` defines the fee to be charged on the creation of a new denom. The fee is drawn from the MsgCreateDenom's sender account, and transferred to the admin's wallet.

`DenomCreationGasConsume` defines the gas cost for creating a new denom. This is intended as a spam deterrence mechanism.

The `coinfactory` module params can be updated with governance.

## Liquid Farming Module

The `x/liquidfarming` module contains the following params:

- fee_collector
- rewards_auction_duration
- liquid_farms

`FeeCollector` is the module account address to collect fees within the liquid farming module.

`RewardsAuctionDuration` is the duration of the rewards auction.

`LiquidFarms` include the farming variables used in liquid farming.

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

`BatchSize` is how many orders to be matched by a tick.

`TickPrecision` is the number of blocks per tick

`FeeCollectorAddress` is the module account address to collect fees within the liquid module.

`DustCollectorAddress` is the module account address to collect dust within the liquid module.

`MinInitialPoolCoinSupply` is the minimum initial pool coin supply.

`PairCreationFee` is the fee to be charged on the creation of a new pair. The fee is drawn from the MsgCreatePair's sender account, and transferred to the fee collector's wallet.

`PoolCreationFee` is the fee to be charged on the creation of a new pool. The fee is drawn from the MsgCreatePool's sender account, and transferred to the fee collector's wallet.

`MinInitialDepositAmount` is the minimum initial deposit amount when providing a liquidity.

`MaxPriceLimitRatio` is the maximum price limit ratio. Default: 10%.

`MaxNumMarketMakingOrderTicks` is the maximum number of market making order ticks.

`MaxOrderLifespan` is the maximum order lifespan.

`SwapFeeRate` is the swap fee rate. Unused.

`WithdrawFeeRate` is the withdraw fee rate.

`DepositExtraGas` is the deposit extra gas.

`WithdrawExtraGas` is the withdraw extra gas.

`OrderExtraGas` is the order extra gas.

`MaxNumActivePoolsPerPair` is the maximum number of active pools per pair.

The `liquidity` module params can be updated with governance.

## LPFrarm Module

The `x/lpfarm` module contains the following params:

- private_plan_creation_fee
- fee_collector
- max_num_private_plans
- max_block_duration

`PrivatePlanCreationFee` specifies the fee for plan creation this fee prevents from spamming.

`FeeCollector` is the module account address to collect fees within the lpfarm module.

`MaxNumPrivatePlans` is the maximum number of active private plans.

`MaxBlockDuration` is the maximum block duration for the private plan.

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

`ValidNFTCollectionId` is a regex which validates the nft collection id.

`NFTCollectionDefaultId` is the default nft collection id.

`NFTCollectionDefaultName` is the default nft collection name.

`ValidNFTCollectionMetadataSymbolMinLength` is the valid nft collection metadata symbol min length.

`ValidNFTCollectionMetadataSymbolMaxLength` is the valid nft collection metadata symbol max length.

`ValidNFTCollectionMetadataDescriptionMaxLength` is the valid nft collection metadata description max length.

`ValidNFTCollectionMetadataNameMaxLength` is the valid nft collection metadata name max length.

`ValidNFTCollectionMetadataImagesMaxCount` is the valid nft collection metadata images max count.

`ValidNFTCollectionMetadataImagesTypeMaxLength` is the valid nft collection metadata images type max length.

`ValidNFTCollectionMetadataLinksMaxCount` is the valid nft collection metadata links max count.

`ValidNFTCollectionMetadataLinksTypeMaxLength` is the valid nft collection metadata links type max length.

`ValidNFTCollectionMetadataOptionsMaxCount` is the valid nft collection metadata options max count.

`ValidNFTCollectionMetadataOptionsTypeMaxLength` is the valid nft collection metadata options type max length.

`ValidNFTCollectionMetadataOptionsValueMaxLength` is the valid nft collection metadata options value max length.

`ValidNFTCollectionMetadataOptionsSubValueMaxLength` is the valid nft collection metadata options sub value max length.

`ValidNFTId` is a regex which validates the nft id.

`ValidNFTMetadataMaxCount` is the max count of the nfts that can be minted in a single batch transaction.

`ValidNFTMetadataTitleMaxLength` is the valid nft metadata title max length.

`ValidNFTMetadataDescriptionMaxLength` is the valid nft metadata description max length.

`ValidNFTMetadataImagesMaxCount` is the valid nft metadata images max count.

`ValidNFTMetadataImagesTypeMaxLength` is the valid nft metadata images type max length.

`ValidNFTMetadataLinksMaxCount` is the valid nft metadata links max count.

`ValidNFTMetadataLinksTypeMaxLength` is the valid nft metadata links type max length.

`ValidNFTMetadataAttributesMaxCount` is the valid nft metadata attributes max count.

`ValidNFTMetadataAttributesTypeMaxLength` is the valid nft metadata attributes type max length.

`ValidNFTMetadataAttributesValueMaxLength` is the valid nft metadata attributes value max length.

`ValidNFTMetadataAttributesSubValueMaxLength` is the valid nft metadata attributes sub value max length.

`ValidBurnNFTMaxCount` is the max count of the nfts that can be burnt in a single batch transaction.

The `token` module params can be updated with governance.

# Cosmos SDK native modules


### Auth Module

The `x/auth` contains the following parameters:

- `max_memo_characters` - max tx memo length // 256
- `sig_verify_cost_ed25519` - difficulty of signature verification // 590
- `sig_verify_cost_secp256k1` - difficulty of signature verification - 1000
- `tx_sig_limit` - transaction signature limit // 7
- `tx_size_cost_per_byte` - transaction size cost per byte // 10

### Bank Module

The `x/bank` contains the following parameters:

- `default_send_enabled` - is the send transactions enabled by default / true

### Crisis

- `constant_fee` -  // 1000uom

### Distribution

- `community_tax` - the tax that gets allocated to the community pool // 2%
- `withdra_addr_enabled` - probably related to mantrachaind tx distribution set-withdraw-addr // true

### gov

- `max_deposit_period` - maximum deposit period 172800s
- `min_deposit` - minimum deposit // 10000000uom
- `burn_proposal_deposit_prevote` - burn the prevote deposits // false
- `burn_vote_quorum` -  burn the // false
- `burn_vote_veto`- burn the // true
- `min_initial_deposit_ratio` - minimum initial deposit // "0.000000000000000000"
- `quorum` - // "0.334000000000000000"
- `threshold` - // "0.500000000000000000"
- `veto_threshold` - // "0.334000000000000000"
- `voting_period` - // 172800s
- `quorum` - // "0.334000000000000000"
- `threshold` - // "0.500000000000000000"
- `veto_threshold` - // "0.334000000000000000"
- `voting_period` - // 172800s


### Mint

- `mint_denom` - which denom is minted // "uom"
- `inflation_rate_change` - the inflation rate change over year // "0.130000000000000000"
- `inflation_max` - max annual inflation // "0.200000000000000000"
- `inflation_min` - min annual inflation // "0.070000000000000000"
- `goal_bonded` - % goal bonded tokens // "0.670000000000000000"
- `blocks_per_year` - blocks per year // "6311520"

### Slashing

- `signed_blocks_window` - blocks per tracked window // "100"
- `min_signed_per_window` -  minimum singed blocks per window // "0.500000000000000000"
- `downtime_jail_duration` - how many seconds downtime the validator is allowed to have // "600000000000"
- `slash_fraction_double_sign` - how much validator is slashed for double signing // "0.050000000000000000"
- `slash_fraction_downtime` - how much is validator slashed for downtime // "0.010000000000000000"

### Staking

- `unbonding_time` - time duration of unbonding // "259200000000000"
- `max_validators` - maximum number of validators // 100
- `key_max_entries` - max entries for either unbonding delegation or redelegation (per pair/trio).// 7
- `historical_entries` - number historical entries to persist in store. //3
- `bond_denom string` - bondable coin denomination // "uom"
- `min_commission_rate` - chain-wide minimum commission rate that a validator can charge their delegators. // "0.000000000000000000"



## ***Note***: For more information about the Cosmos SDK native modules params, please refer to the [Cosmos SDK docs](https://docs.cosmos.network/v0.47)
