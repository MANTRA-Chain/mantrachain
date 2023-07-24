import {getGasFee, MantrachainSdk} from '../helpers/sdk'
import { createDenomIfNotExists, genCoinDenom } from "../helpers/coinfactory";
import { getPairId } from '../helpers/liquidity';
import { updateAccountPrivileges, updateCoinRequiredPrivileges } from '../helpers/guard';
import {mintGuardSoulBondNft} from "../helpers/token";
import {Privileges, utils} from "../../../../../mantrachain-sdk/src/utils";

describe('Lpfarm module', () => {
    let sdk: MantrachainSdk

    let baseCoinDenom = 'liquidity' + new Date().getTime().toString();
    // with this we manipulate the time and space
    // of this test environment according to our needs
    // to the infinity and beyond!
    let quoteCoinDenom = 'liquidity' + new Date().getTime().toString() + 1;

    let pairId = 0;

    let poolId;

    beforeAll(async () => {
        sdk = new MantrachainSdk()
        await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)

        await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, baseCoinDenom)
        await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, quoteCoinDenom)

        await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
            value: {
                sender: sdk.adminAddress,
                amount: {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                    amount: "10000000000000000000"
                }
            },
            fee: getGasFee()
        })

        await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
            value: {
                sender: sdk.adminAddress,
                amount: {
                    denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                    amount: "10000000000000000000"
                }
            },
            fee: getGasFee()
        })

        const res = await sdk.clientAdmin.MantrachainGuardV1.query.queryParams();
        const defaultPrivileges = utils.base64ToBytes(
            res.data.params.default_privileges
        );

        const privileges = Privileges.fromBuffer(defaultPrivileges).set(64).toBuffer();

        await sdk.clientAdmin.MantrachainGuardV1.tx.sendMsgUpdateRequiredPrivileges({
            value: {
                creator: sdk.adminAddress,
                index: utils.strToIndex(
                    genCoinDenom(sdk.adminAddress, baseCoinDenom)
                ), // denom of the coin from the `CoinFactory` module represented in bytes
                privileges,
                kind: "coin",
            },
            fee: getGasFee()
        });

        await sdk.clientAdmin.MantrachainGuardV1.tx.sendMsgUpdateRequiredPrivileges({
            value: {
                creator: sdk.adminAddress,
                index: utils.strToIndex(
                    genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                ), // denom of the coin from the `CoinFactory` module represented in bytes
                privileges,
                kind: "coin",
            },
            fee: getGasFee()
        });

        await sdk.clientAdmin.MantrachainTokenV1.tx.sendMsgMintNft({
            value: {
                creator: sdk.adminAddress,
                receiver: sdk.recipientAddress,
                collectionCreator: sdk.adminAddress,
                collectionId: "account_privileges_guard_nft_collection",
                nft: {
                    id: sdk.recipientAddress,
                    title: 'AccountPrivileges',
                    images: [],
                    url: '',
                    description: 'AccountPrivileges',
                    links: [],
                    attributes: [],
                    data: undefined
                },
                strict: true,
                did: false
            },
            fee: getGasFee()
        })

        await sdk.clientAdmin.MantrachainGuardV1.tx.sendMsgUpdateAccountPrivileges({
            value: {
                creator: sdk.adminAddress,
                account: sdk.recipientAddress,
                privileges: privileges
            },
            fee: getGasFee()
        })

        await sdk.clientAdmin.CosmosBankV1Beta1.tx.sendMsgSend({
            value: {
                fromAddress: sdk.adminAddress,
                toAddress: sdk.recipientAddress,
                amount: [
                    {
                        denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                        amount: '100000000000'
                    },
                    {
                        denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                        amount: '100000000000'
                    }
                ]
            },
            fee: getGasFee()
        })


        await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
            value: {
                creator: sdk.adminAddress,
                baseCoinDenom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                quoteCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom)
            },
            fee: getGasFee()
        })

        const allPairs = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPairs();
        const lastPair = allPairs.data.pairs.pop();
        pairId = lastPair.id

        await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePool(
            {
                value: {
                    creator: sdk.adminAddress,
                    pairId: pairId,
                    depositCoins: [
                        {
                            denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                            amount: "1000000"
                        },
                        {
                            denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                            amount: "1000000"
                        }
                    ]
                },
                fee: getGasFee()
            }
        )

        const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPools({
            pair_id: pairId.toString()
        });

        const poolIndex = resp.data.pools.length - 1;
        poolId = resp.data.pools[poolIndex].id;

        await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgDeposit(
            {
                value: {
                    depositor: sdk.recipientAddress,
                    poolId: poolId,
                    depositCoins: [
                        {
                            denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                            amount: "1000000"
                        },
                        {
                            denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                            amount: "1000000"
                        }
                    ]
                },
                fee: getGasFee()
            }
        )
    })

    describe('User', () => {
        test('should throw when trying to create private farming plan', async () => {
            const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPool(poolId);

            const pool = resp.data.pool;

            const res = sdk.clientRecipient.MantrachainLpfarmV1Beta1.tx.sendMsgCreatePrivatePlan({
                value: {
                    creator: sdk.recipientAddress,
                    description: 'money',
                    rewardAllocations: [
                        {
                            pairId: 0,
                            denom: pool.pool_coin_denom,
                            rewardsPerDay: [
                                {
                                    denom: 'uaum',
                                    value: '100000000000'
                                }
                            ]
                        },
                    ],
                    startTime: new Date(),
                    endTime: new Date(new Date().setMonth(new Date().getMonth() + 1))
                },
                fee: getGasFee()
            })

            await expect(res).rejects.toThrow()
        })

        test('should be able to stake its pool coins', async () => {
            const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPool(poolId);
            const pool = resp.data.pool;

            await sdk.clientAdmin.MantrachainLpfarmV1Beta1.tx.sendMsgCreatePrivatePlan({
                value: {
                    creator: sdk.adminAddress,
                    description: 'money',
                    rewardAllocations: [
                        {
                            pairId: 0,
                            denom: pool.pool_coin_denom,
                            rewardsPerDay: [
                                {
                                    denom: 'uaum',
                                    amount: '1000000'
                                }
                            ]
                        },
                    ],
                    startTime: new Date(new Date().setMinutes(new Date().getMinutes() - 10)),
                    endTime: new Date(new Date().setMonth(new Date().getMonth() + 1))
                },
                fee: getGasFee()
            })

            const plan = await sdk.clientAdmin.MantrachainLpfarmV1Beta1.query.queryPlans();
            const lastPlan = plan.data.plans.pop();
            const planFarmingPoolAddress = lastPlan.farming_pool_address;

            await sdk.clientAdmin.CosmosBankV1Beta1.tx.sendMsgSend({
                value: {
                    fromAddress: sdk.adminAddress,
                    toAddress: planFarmingPoolAddress,
                    amount: [{
                        denom: 'uaum',
                        amount: '100000000000'
                    }]
                },
                fee: getGasFee()
            })

            await sdk.clientRecipient.MantrachainLpfarmV1Beta1.tx.sendMsgFarm({
                value: {
                    farmer: sdk.recipientAddress,
                    coin: {
                        denom: pool.pool_coin_denom,
                        amount: '1000000'
                    }
                },
                fee: getGasFee()
            })

            const position = await sdk.clientRecipient.MantrachainLpfarmV1Beta1.query.queryPosition(
                sdk.recipientAddress, pool.pool_coin_denom
            )

            expect(Number(position.data.position.farming_amount)).toBe(1000000)
        })

        test('should be able to claim accumulated staking rewards', async () => {
            const balanceOfRewardCoinBefore = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {denom: 'uaum'}
            )

            const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPool(poolId);
            const pool = resp.data.pool;

            await new Promise((r) => setTimeout(r, 7000));

            await sdk.clientRecipient.MantrachainLpfarmV1Beta1.tx.sendMsgHarvest({
                value: {
                    farmer: sdk.recipientAddress,
                    denom: pool.pool_coin_denom
                },
                fee: getGasFee()
            })


            const balanceOfRewardCoinAfter = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {denom: 'uaum'}
            )

            expect(Number(balanceOfRewardCoinAfter.data.balance.amount)).toBeGreaterThan(Number(balanceOfRewardCoinBefore.data.balance.amount))
        })

        test('should be able to unfarm his pool tokens', async () => {
            const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPool(poolId);
            const pool = resp.data.pool;

            await sdk.clientRecipient.MantrachainLpfarmV1Beta1.tx.sendMsgUnfarm({
                value: {
                    farmer: sdk.recipientAddress,
                    coin: {
                        denom: pool.pool_coin_denom,
                        amount: '1000000'
                    }
                },
                fee: getGasFee()
            })

            const balanceOfPoolCoins = await sdk.clientAdmin.CosmosBankV1Beta1.query.queryBalance(sdk.recipientAddress, {
                denom: pool.pool_coin_denom
            })

            expect(Number(balanceOfPoolCoins.data.balance.amount)).toBeGreaterThan(0)
        })
    })
})