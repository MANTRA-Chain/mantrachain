<!-- order: 5 -->

# Events

## Handlers

### MsgCreatePrivatePlan

| Type                                           | Attribute Key        | Attribute Value                                |
|------------------------------------------------|----------------------|------------------------------------------------|
| message                                        | action               | /aumega.lpfarm.v1beta1.Msg/CreatePrivatePlan |
| aumega.lpfarm.v1beta1.EventCreatePrivatePlan | creator              | {planCreatorAddress}                           |
| aumega.lpfarm.v1beta1.EventCreatePrivatePlan | plan_id              | {planId}                                       |
| aumega.lpfarm.v1beta1.EventCreatePrivatePlan | farming_pool_address | {farmingPoolAddress}                           |

### MsgFarm

| Type                              | Attribute Key     | Attribute Value                   |
|-----------------------------------|-------------------|-----------------------------------|
| message                           | action            | /aumega.lpfarm.v1beta1.Msg/Farm |
| aumega.lpfarm.v1beta1.EventFarm | farmer            | {farmerAddress}                   |
| aumega.lpfarm.v1beta1.EventFarm | coin              | {coin}                            |
| aumega.lpfarm.v1beta1.EventFarm | withdrawn_rewards | {withdrawnRewards}                |

### MsgUnfarm

| Type                                | Attribute Key     | Attribute Value                     |
|-------------------------------------|-------------------|-------------------------------------|
| message                             | action            | /aumega.lpfarm.v1beta1.Msg/Unfarm |
| aumega.lpfarm.v1beta1.EventUnfarm | farmer            | {farmerAddress}                     |
| aumega.lpfarm.v1beta1.EventUnfarm | coin              | {coin}                              |
| aumega.lpfarm.v1beta1.EventUnfarm | withdrawn_rewards | {withdrawnRewards}                  |

### MsgHarvest

| Type                                 | Attribute Key     | Attribute Value                      |
|--------------------------------------|-------------------|--------------------------------------|
| message                              | action            | /aumega.lpfarm.v1beta1.Msg/Harvest |
| aumega.lpfarm.v1beta1.EventHarvest | farmer            | {farmerAddress}                      |
| aumega.lpfarm.v1beta1.EventHarvest | denom             | {farmingAssetDenom}                  |
| aumega.lpfarm.v1beta1.EventHarvest | withdrawn_rewards | {withdrawnRewards}                   |
