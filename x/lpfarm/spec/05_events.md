<!-- order: 5 -->

# Events

## Handlers

### MsgCreatePrivatePlan

| Type                                           | Attribute Key        | Attribute Value                                |
|------------------------------------------------|----------------------|------------------------------------------------|
| message                                        | action               | /mantrachain.lpfarm.v1beta1.Msg/CreatePrivatePlan |
| mantrachain.lpfarm.v1beta1.EventCreatePrivatePlan | creator              | {planCreatorAddress}                           |
| mantrachain.lpfarm.v1beta1.EventCreatePrivatePlan | plan_id              | {planId}                                       |
| mantrachain.lpfarm.v1beta1.EventCreatePrivatePlan | farming_pool_address | {farmingPoolAddress}                           |

### MsgFarm

| Type                              | Attribute Key     | Attribute Value                   |
|-----------------------------------|-------------------|-----------------------------------|
| message                           | action            | /mantrachain.lpfarm.v1beta1.Msg/Farm |
| mantrachain.lpfarm.v1beta1.EventFarm | farmer            | {farmerAddress}                   |
| mantrachain.lpfarm.v1beta1.EventFarm | coin              | {coin}                            |
| mantrachain.lpfarm.v1beta1.EventFarm | withdrawn_rewards | {withdrawnRewards}                |

### MsgUnfarm

| Type                                | Attribute Key     | Attribute Value                     |
|-------------------------------------|-------------------|-------------------------------------|
| message                             | action            | /mantrachain.lpfarm.v1beta1.Msg/Unfarm |
| mantrachain.lpfarm.v1beta1.EventUnfarm | farmer            | {farmerAddress}                     |
| mantrachain.lpfarm.v1beta1.EventUnfarm | coin              | {coin}                              |
| mantrachain.lpfarm.v1beta1.EventUnfarm | withdrawn_rewards | {withdrawnRewards}                  |

### MsgHarvest

| Type                                 | Attribute Key     | Attribute Value                      |
|--------------------------------------|-------------------|--------------------------------------|
| message                              | action            | /mantrachain.lpfarm.v1beta1.Msg/Harvest |
| mantrachain.lpfarm.v1beta1.EventHarvest | farmer            | {farmerAddress}                      |
| mantrachain.lpfarm.v1beta1.EventHarvest | denom             | {farmingAssetDenom}                  |
| mantrachain.lpfarm.v1beta1.EventHarvest | withdrawn_rewards | {withdrawnRewards}                   |
