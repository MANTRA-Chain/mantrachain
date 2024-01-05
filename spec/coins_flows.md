# Aumega Coins Flows

## Tx Fees

There are several flags users can use to indicate how much they are willing to pay in fees for a Tx. The flags are:

`--gas` refers to how much gas, which represents computational resources, Tx consumes. Gas is dependent on the transaction and is not precisely calculated until execution, but can be estimated by providing auto as the value for --gas.

`--gas-adjustment` (optional) can be used to scale gas up in order to avoid underestimating. For example, users can specify their gas adjustment as 1.5 to use 1.5 times the estimated gas.

`--gas-prices` specifies how much the user is willing to pay per unit of gas, which can be one or multiple denominations of tokens. For example, --gas-prices=0.025uaum means the user is willing to pay 0.025uaum per unit of gas.

`--fees` specifies how much in fees the user is willing to pay in total.

`--timeout-height` specifies a block timeout height to prevent the tx from being committed past a certain height.

The ultimate value of the fees paid is equal to the gas multiplied by the gas prices. In other words, fees = ceil(gas * gasPrices). Thus, since fees can be calculated using gas prices and vice versa, the users specify only one of the two.

Later, validators decide whether or not to include the transaction in their block by comparing the given or calculated gas-prices to their local `min-gas-prices`. Tx is rejected if its gas-prices is not high enough, so users are incentivized to pay more.

## ***Note***: All of the gas fees are sent to an ***admin wallet*** and not to the validators

## Gas

In the Cosmos SDK, gas is a special unit that is used to track the consumption of resources during execution. gas is typically consumed whenever read and writes are made to the store, but it can also be consumed if expensive computation needs to be done. It serves two main purposes:

- Make sure blocks are not consuming too many resources and are finalized. This is implemented by default in the Cosmos SDK via the block gas meter.
- Prevent spam and abuse from end-user. To this end, gas consumed during message execution is typically priced, resulting in a fee (`fees = gas * gas-prices`). fees generally have to be paid by the sender of the message. Note that the Cosmos SDK does not enforce gas pricing by default, as there may be other ways to prevent spam (e.g. bandwidth schemes). Still, most applications implement fee mechanisms to prevent spam by using the AnteHandler.

The auth module AnteHandler checks and increments sequence numbers, checks signatures and account numbers, and deducts fees from the first signer of the transaction.

The `GasMeter` tracks how much gas is used during the execution of Tx. The user-provided amount of gas for Tx is known as `GasWanted`. If `GasConsumed`, the amount of gas consumed during execution, ever exceeds `GasWanted`, the execution stops and the changes made to the cached copy of the state are not committed. Otherwise, CheckTx sets `GasUsed` equal to `GasConsumed` and returns it in the result. After calculating the gas and fee values, validator-nodes check that the user-specified gas-prices is greater than their locally defined `min-gas-prices`.

## The Minting Mechanism

The minting mechanism was designed to:

- allow for a flexible inflation rate determined by market demand targeting a particular bonded-stake ratio
- effect a balance between market liquidity and staked supply
  
In order to best determine the appropriate market rate for inflation rewards, a moving change rate is used. The moving change rate mechanism ensures that if the % bonded is either over or under the goal `%-bonded`, the inflation rate will adjust to further incentivize or disincentivize being bonded, respectively. Setting the goal `%-bonded` at less than 100% encourages the network to maintain some non-staked tokens which should help provide some liquidity.

It can be broken down in the following way:

- If the inflation rate is below the goal %-bonded the inflation rate will increase until a maximum value is reached
- If the goal % bonded (67% in Cosmos-Hub) is maintained, then the inflation rate will stay constant
- If the inflation rate is above the goal %-bonded the inflation rate will decrease until a minimum value is reached

Minting parameters are recalculated and inflation paid at the beginning of each block.

The target annual inflation rate is recalculated each block. The inflation is also subject to a rate change (positive or negative) depending on the distance from the desired ratio (67%). The maximum rate change possible is defined to be 13% per year, however the annual inflation is capped as between 7% and 20%.

The minting module contains the following parameters:

|Key|Type|Example|
|-|-|-|
|MintDenom|string|"uaum"|
|InflationRateChange|string (dec)|"0.130000000000000000"|
|InflationMax|string (dec)|"0.200000000000000000"|
|InflationMin|string (dec)|"0.070000000000000000"|
|BlocksPerYear|string (uint64)|"6311520"|

Setting `InflationMin`, `InflationMax` and `InflationRateChange` to 0 will disable inflation.

A user can query and interact with the mint module using the `CLI`, `gRPC` and/or `REST`.

### REST

A user can query the mint module using REST endpoints.

#### Annual Provisions

```js
/cosmos/mint/v1beta1/annual_provisions
```

Example Output:

```js
{
  "annualProvisions": "1432452520532626265712995618"
}
```

#### Inflation

```js
/cosmos/mint/v1beta1/inflation
```

Example Output:

```js
{
  "inflation": "130197115720711261"
}
```

#### Params

```js
/cosmos/mint/v1beta1/params
```

Example Output:

```js
{
  "params": {
    "mintDenom": "uaum",
    "inflationRateChange": "130000000000000000",
    "inflationMax": "200000000000000000",
    "inflationMin": "70000000000000000",
    "goalBonded": "670000000000000000",
    "blocksPerYear": "6311520"
  }
}
```

Those parameters can be updated through a governance proposal.

For more information on the mint mechanism, see the [Cosmos SDK docs](https://docs.cosmos.network/v0.47/modules/mint).

## The [F1 Distribution Mechanism](https://github.com/cosmos/cosmos-sdk/blob/main/docs/spec/fee_distribution/f1_fee_distr.pdf)

The distribution mechanism was designed to:

- reward validators for their service to the network
- reward delegators for their service to the network
- incentivize validators to maximize their up-time
- incentivize delegators to delegate to validators with high up-time
- incentivize validators to maximize their self-bonded stake
- incentivize delegators to delegate to validators with high self-bonded stake

This simple distribution mechanism describes a functional way to passively distribute rewards between validators and delegators. Note that this mechanism does not distribute funds in as precisely as active reward distribution mechanisms and will therefore be upgraded in the future.

The mechanism operates as follows. Collected rewards are pooled globally and divided out passively to validators and delegators. Each validator has the opportunity to charge commission to the delegators on the rewards collected on behalf of the delegators. Fees are collected directly into a global reward pool and validator proposer-reward pool. Due to the nature of passive accounting, whenever changes to parameters which affect the rate of reward distribution occurs, withdrawal of rewards must also occur.

- Whenever withdrawing, one must withdraw the maximum amount they are entitled to, leaving nothing in the pool.
- Whenever bonding, unbonding, or re-delegating tokens to an existing account, a full withdrawal of the rewards must occur (as the rules for lazy accounting change).
- Whenever a validator chooses to change the commission on rewards, all accumulated commission rewards must be simultaneously withdrawn.

As a part of the lazy computations, each delegator holds an accumulation term specific to each validator which is used to estimate what their approximate fair portion of tokens held in the global fee pool is owed to them.

```js
entitlement = delegator-accumulation / all-delegators-accumulation
```

Rewards are calculated per period. The period is updated each time a validator's delegation changes, for example, when the validator receives a new delegation. The rewards for a single validator can then be calculated by taking the total rewards for the period before the delegation started, minus the current total rewards.

Under the circumstance that there was constant and equal flow of incoming reward tokens every block, this distribution mechanism would be equal to the active distribution (distribute individually to all delegators each block). However, this is unrealistic so deviations from the active distribution will occur based on fluctuations of incoming reward tokens as well as timing of reward withdrawal by other delegators.

In F1 fee distribution, the rewards a delegator receives are calculated when their delegation is withdrawn. This calculation must read the terms of the summation of rewards divided by the share of tokens from the period which they ended when they delegated, and the final period that was created for the withdrawal.

Additionally, as slashes change the amount of tokens a delegation will have (but we calculate this lazily, only when a delegator un-delegates), we must calculate rewards in separate periods before / after any slashes which occurred in between when a delegator delegated and when they withdrew their rewards. Thus slashes, like delegations, reference the period which was ended by the slash event.

All stored historical rewards records for periods which are no longer referenced by any delegations or any slashes can thus be safely removed, as they will never be read (future delegations and future slashes will always reference future periods). This is implemented by tracking a ReferenceCount along with each historical reward storage entry. Each time a new object (delegation or slash) is created which might need to reference the historical record, the reference count is incremented. Each time one object which previously needed to reference the historical record is deleted, the reference count is decremented. If the reference count hits zero, the historical record is deleted.

Validator distribution information for the relevant validator is updated each time:

01. delegation amount to a validator is updated,
02. any delegator withdraws from a validator, or
03. the validator withdraws its commission.

Each delegation distribution only needs to record the height at which it last withdrew fees. Because a delegation must withdraw fees each time it's properties change (aka bonded tokens etc.) its properties will remain constant and the delegator's accumulation factor can be calculated passively knowing only the height of the last withdrawal and its current properties.

The distribution module contains the following parameters:

|Key|Type|Example|
|-|-|-|
|CommunityTax|string (dec)|"0.020000000000000000"|
|BaseProposerReward|string (dec)|"0.010000000000000000"|
|BonusProposerReward|string (dec)|"0.040000000000000000"|
|WithdrawAddrEnabled|bool|true|

Note: The `base_proposer_reward` and `bonus_proposer_reward` fields are deprecated and are no longer used in the `x/distribution` module's reward mechanism.

At each `BeginBlock`, all fees received in the previous block are transferred to the distribution ModuleAccount account. When a delegator or validator withdraws their rewards, they are taken out of the ModuleAccount. During begin block, the different claims on the fees collected are updated as follows:

- The reserve community tax is charged.
- The remainder is distributed proportionally by voting power to all bonded validators

Let `fees` be the total fees collected in the previous block, including inflationary rewards to the stake. All fees are collected in a specific module account during the block. During `BeginBlock`, they are sent to the "distribution" `ModuleAccount`. No other sending of tokens occurs. Instead, the rewards each account is entitled to are stored, and withdrawals can be triggered.

### Reward to the Community Pool

The community pool gets `community_tax * fees`, plus any remaining dust after validators get their rewards that are always rounded down to the nearest integer value.

### Reward To the Validators

The proposer receives no extra rewards. All fees are distributed among all the bonded validators, including the proposer, in proportion to their consensus power.

```js
powFrac = validator power / total bonded validator power
voteMul = 1 - community_tax
```

All validators receive `fees * voteMul * powFrac`.

### Rewards to Delegators

Each validator's rewards are distributed to its delegators. The validator also has a self-delegation that is treated like a regular delegation in distribution calculations.

The validator sets a commission rate. The commission rate is flexible, but each validator sets a maximum rate and a maximum daily increase. These maximums cannot be exceeded and protect delegators from sudden increases of validator commission rates to prevent validators from taking all of the rewards.

### Example Distribution

For this example distribution, the underlying consensus engine selects block proposers in proportion to their power relative to the entire bonded power.

All validators are equally performant at including pre-commits in their proposed blocks. Then hold `(pre_commits included) / (total bonded validator power)` constant so that the amortized block reward for the validator is `( validator power / total bonded power) * (1 - community tax rate)` of the total rewards. Consequently, the reward for a single delegator is:

```js
(delegator proportion of the validator power / validator power) * (validator power / total bonded power)
  * (1 - community tax rate) * (1 - validator commission rate)
= (delegator proportion of the validator power / total bonded power) * (1 -
community tax rate) * (1 - validator commission rate)
```

For more information on the distribution mechanism, see the [Cosmos SDK docs](https://docs.cosmos.network/v0.47/modules/distribution).

## Slashing

`x/evidence` is an implementation of a Cosmos SDK module, that allows for the submission and handling of arbitrary evidence of misbehavior such as equivocation and counterfactual signing.

The evidence module differs from standard evidence handling which typically expects the underlying consensus engine, e.g. Tendermint, to automatically submit evidence when it is discovered by allowing clients and foreign chains to submit more complex evidence directly.

The Cosmos SDK handles two types of evidence inside the ABCI BeginBlock:

- DuplicateVoteEvidence
- LightClientAttackEvidence

The slashing, jailing, and tombstoning calls are delegated through the `x/slashing` module that emits informative events and finally delegates calls to the `x/staking` module.

When a `Validator` is slashed, the following occurs:

- The total `slashAmount` is calculated as the `slashFactor` (a chain parameter) * `TokensFromConsensusPower`, the total number of tokens bonded to the validator at the time of the infraction.
- Every unbonding delegation and pseudo-unbonding redelegation such that the infraction occured before the unbonding or redelegation began from the validator are slashed by the `slashFactor` percentage of the `initialBalance`.
- Each amount slashed from redelegations and unbonding delegations is subtracted from the total slash amount.
- The `remaingSlashAmount` is then slashed from the validator's tokens in the `BondedPool` or `NonBondedPool` depending on the validator's status. This reduces the total supply of tokens.

For more information on the evidence mechanism, see the [Cosmos SDK docs](https://docs.cosmos.network/v0.47/modules/evidence).

## Crisis

The crisis module halts the blockchain under the circumstance that a blockchain invariant is broken. Invariants can be registered with the application during the application initialization process.

Due to the anticipated large gas cost requirement to verify an invariant (and potential to exceed the maximum allowable block gas limit) a constant fee is used instead of the standard gas consumption method. The constant fee is intended to be larger than the anticipated gas cost of running the invariant with the standard gas consumption method.

The `ConstantFee` param is stored in the module params state, it can be updated with governance.

For more information on the crisis mechanism, see the [Cosmos SDK docs](https://docs.cosmos.network/v0.47/modules/crisis).

## Staking

The module enables Cosmos SDK-based blockchain to support an advanced Proof-of-Stake (PoS) system. In this system, holders of the native staking token of the chain can become validators and can delegate tokens to validators, ultimately determining the effective validator set for the system.

### Params

The `staking` module params can be updated with governance.

|Key|Type|Example|
|-|-|-|
|UnbondingTime|string (duration)|"1814400000000000"|
|MaxValidators|string (uint64)|"100"|
|MaxEntries|string (uint64)|"7"|
|HistoricalEntries|string (uint64)|"10000"|
|BondDenom|string|"uaum"|
|MinCommisionRate|string (uint64)|"0.010000000000000000"|

- unbonding_time is the time duration of unbonding.
- max_validators is the maximum number of validators.
- max_entries is the max entries for either unbonding delegation or redelegation (per pair/trio).
- historical_entries is the number of historical entries to persist in store.
- bond_denom is the bondable coin denomination.
- min_commission_rate is the chain-wide minimum commission rate that a validator can charge their delegators.
  
### Validators

Validators can have one of three statuses

- Unbonded: The validator is not in the active set. They cannot sign blocks and do not earn rewards. They can receive delegations.
- Bonded: Once the validator receives sufficient bonded tokens they automatically join the active set during `EndBlock` and their status is updated to Bonded. They are signing blocks and receiving rewards. They can receive further delegations. They can be slashed for misbehavior. Delegators to this validator who unbond their delegation must wait the duration of the UnbondingTime, a chain-specific param, during which time they are still slashable for offences of the source validator if those offences were committed during the period of time that the tokens were bonded.
- Unbonding: When a validator leaves the active set, either by choice or due to slashing, jailing or tombstoning, an unbonding of all their delegations begins. All delegations must then wait the `UnbondingTime` before their tokens are moved to their accounts from the `BondedPool`.

#### Delegator Shares

When one `Delegates` tokens to a `Validator` they are issued a number of delegator shares based on a dynamic exchange rate, calculated as follows from the total number of tokens delegated to the validator and the number of shares issued so far:

```js
Shares per Token = validator.TotalShares() / validator.Tokens()
```

Only the number of shares received is stored on the `DelegationEntry`. When a delegator then `Undelegates`, the token amount they receive is calculated from the number of shares they currently hold and the inverse exchange rate:

```js
Tokens per Share = validator.Tokens() / validatorShares()
```

These `Shares` are simply an accounting mechanism. They are not a fungible asset. The reason for this mechanism is to simplify the accounting around slashing. Rather than iteratively slashing the tokens of every delegation entry, instead the Validators total bonded tokens can be slashed, effectively reducing the value of each issued delegator share.

#### UnbondingDelegation

`Shares` in a `Delegation` can be unbonded, but they must for some time exist as an `UnbondingDelegation`, where shares can be reduced if Byzantine behavior is detected.

#### Redelegation

The bonded tokens worth of a `Delegation` may be instantly redelegated from a source validator to a different validator (destination validator). However when this occurs they must be tracked in a Redelegation object, whereby their shares can be slashed if their tokens have contributed to a Byzantine fault committed by the source validator.

#### State Transitions

State transitions in validators are performed on every `EndBlock` in order to check for changes in the active `ValidatorSet`.

A validator can be `Unbonded`, `Unbonding` or `Bonded`. `Unbonded` and `Unbonding` are collectively called `Not Bonded`. A validator can move directly between all the states, except for from `Bonded` to `Unbonded`.

#### Jail/Unjail

when a validator is jailed it is effectively removed from the CometBFT set.

### How Shares are calculated

At any given point in time, each validator has a number of tokens, `T`, and has a number of shares issued, `S`. Each delegator, `i`, holds a number of shares, `S_i`. The number of tokens is the sum of all tokens delegated to the validator, plus the rewards, minus the slashes.

The delegator is entitled to a portion of the underlying tokens proportional to their proportion of shares. So delegator `i` is entitled to `T * S_i / S` of the validator's tokens.

When a delegator delegates new tokens to the validator, they receive a number of shares proportional to their contribution. So when delegator `j` delegates `T_j` tokens, they receive `S_j = S * T_j / T` shares. The total number of tokens is now `T + T_j`, and the total number of shares is `S + S_j`. js proportion of the shares is the same as their proportion of the total tokens contributed: `(S + S_j) / S = (T + T_j) / T`.

A special case is the initial delegation, when `T = 0` and `S = 0`, so `T_j / T` is undefined. For the initial delegation, delegator j who delegates `T_j` tokens receive `S_j = T_j` shares. So a validator that hasn't received any rewards and has not been slashed will have `T = S`.

When a validator is slashed, the number of tokens is reduced by the slash amount, `T' = T - T_slash`. The number of shares is also reduced by the same proportion, `S' = S - S_slash`. The number of shares held by each delegator is reduced by the same proportion, `S_i' = S_i - S_slash * S_i / S`. The proportion of shares held by each delegator is the same as before the slash: `S_i' / S' = S_i / S`.

When a validator receives rewards, the number of tokens is increased by the reward amount, `T' = T + T_reward`. The number of shares is also increased by the same proportion, `S' = S + S_reward`. The number of shares held by each delegator is increased by the same proportion, `S_i' = S_i + S_reward * S_i / S`. The proportion of shares held by each delegator is the same as before the reward: `S_i' / S' = S_i / S`.

For more information on the staking mechanism, see the [Cosmos SDK docs](https://docs.cosmos.network/v0.47/modules/staking).

## Aumega Additional Modules

### TxFees Module

The `x/txfees` module is used to:

- send the gas fees to an admin wallet
- set a fee token so the gas fees can be paid in a different token than the native token

The `x/txfees` module params contain the `base_denom` param, which is the token in which the gas fees are paid.

In order to add additional tokens, the `admin` should make a `CreateFeeToken` transaction.
Adding additional gas fees tokens is not possible through governance. The `admin` can also update/remove a token by making a `UpdateFeeToken`/`RemoveFeeToken` transaction. The `create` and `update` transactions require a `pairId` param which correspond to a pair including the `base denom` and the corresponding `fee denom`.

### Coinfactory Module

The `x/coinfactory` module is used to:

- create a new token
- mint tokens
- burn tokens
- seize tokens

In order to create a new token, the `admin` should make a `CreateDenom` transaction. The `admin` can also mint/burn/seize tokens by making a `Mint`/`Burn`/`ForceTransfer` transaction. The `mint` and `burn` transactions require a custom token created through the module, while force transfer works for any token.
