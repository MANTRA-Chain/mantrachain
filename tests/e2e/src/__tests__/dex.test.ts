import {getGasFee, MantrachainSdk} from "../helpers/sdk";
import {createDenomIfNotExists, genCoinDenom} from "../helpers/coinfactory";
import {updateAccountPrivileges, updateCoinRequiredPrivileges} from "../helpers/guard";
import {Privileges, utils} from "../../../../../mantrachain-sdk";

/** OrderDirection enumerates order directions. */
export enum OrderDirection {
    /** ORDER_DIRECTION_UNSPECIFIED - ORDER_DIRECTION_UNSPECIFIED specifies unknown order direction */
    ORDER_DIRECTION_UNSPECIFIED = 0,
    /** ORDER_DIRECTION_BUY - ORDER_DIRECTION_BUY specifies buy(swap quote coin to base coin) order direction */
    ORDER_DIRECTION_BUY = 1,
    /** ORDER_DIRECTION_SELL - ORDER_DIRECTION_SELL specifies sell(swap base coin to quote coin) order direction */
    ORDER_DIRECTION_SELL = 2,
    UNRECOGNIZED = -1,
}

describe('Liquidity module', () => {
    let sdk: MantrachainSdk

    let baseCoinDenom = 'atom' + new Date().getTime().toString();
    // with this we manipulate the time and space
    // of this test environment according to our needs
    // to the infinity and beyond!
    let quoteCoinDenom = 'osmo' + new Date().getTime().toString() + 1;

    let pairId = 0;
    let poolId = 0;

    let basicPoolId;
    let rangedPoolId;

    let basicPool;
    let rangedPool;

    let genericUsers: {
        address: string,
        wallet,
        client
    }[] = [];

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
                    amount: "1000000000000000000"
                }
            },
            fee: getGasFee()
        })

        await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
            value: {
                sender: sdk.adminAddress,
                amount: {
                    denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                    amount: "1000000000000000000"
                }
            },
            fee: getGasFee()
        })

        // await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
        //     value: {
        //         creator: sdk.adminAddress,
        //         baseCoinDenom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
        //         quoteCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom)
        //     },
        //     fee: getGasFee()
        // })
        //
        // const allPairs = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPairs();
        // const lastPair = allPairs.data.pairs.pop();
        // pairId = lastPair.id
        //
        //
        // await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePool(
        //     {
        //         value: {
        //             creator: sdk.adminAddress,
        //             pairId: pairId,
        //             depositCoins: [
        //                 {
        //                     denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
        //                     amount: "10000000"
        //                 },
        //                 {
        //                     denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
        //                     amount: "10000000"
        //                 }
        //             ]
        //         },
        //         fee: {
        //             amount: [{ denom: "uaum", amount: "42" }],
        //             gas: "210000"
        //         }
        //     }
        // )
        //
        // await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreateRangedPool(
        //     {
        //         value: {
        //             creator: sdk.adminAddress,
        //             pairId: pairId,
        //             depositCoins: [
        //                 {
        //                     denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
        //                     amount: "10000000"
        //                 },
        //                 {
        //                     denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
        //                     amount: "10000000"
        //                 }
        //             ],
        //             minPrice: '800000000000000000',
        //             maxPrice: '1200000000000000000',
        //             initialPrice: '1000000000000000000',
        //         },
        //         fee: {
        //             amount: [{denom: "uaum", amount: "40"}],
        //             gas: "200000"
        //         }
        //     }
        // )
        //
        // const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPools({
        //     pair_id: pairId.toString()
        // });
        //
        //
        // rangedPool = resp.data.pools.pop()
        // basicPool = resp.data.pools.pop()
        //
        // rangedPoolId = rangedPool.id
        // basicPoolId = basicPool.id

        // for (let i = 1; i <= 3; i++) {
        //     const wallet = await DirectSecp256k1HdWallet.generate(24, {prefix: 'mantra'});
        //     const user = {
        //         address: (await wallet.getAccounts())[0].address,
        //         wallet: wallet,
        //         client: new Client({
        //                 apiURL: process.env.API_URL,
        //                 rpcURL: process.env.RPC_URL,
        //             },
        //             wallet
        //         )
        //     }
        //
        //     await sdk.clientAdmin.CosmosBankV1Beta1.tx.sendMsgSend({
        //         value: {
        //             fromAddress: sdk.adminAddress,
        //             toAddress: user.address,
        //             amount: [
        //                 {
        //                     denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
        //                     amount: '10000000000000'
        //                 },
        //                 {
        //                     denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
        //                     amount: '10000000000000'
        //                 },
        //                 {
        //                     denom: 'uaum',
        //                     amount: '100000000'
        //                 }
        //             ]
        //         },
        //         fee: getGasFee()
        //     })
        //
        //     genericUsers.push(user);
        // }
    })

    describe('Admin', () => {
        test('should be able to create pair for existing denoms', async () => {

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

            expect(Number(lastPair.id)).toBeGreaterThan(0);
        });

        test('should throw when trying to create already existing pair for existing denoms', async () => {

            const res = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
                value: {
                    creator: sdk.adminAddress,
                    baseCoinDenom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                    quoteCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom)
                },
                fee: getGasFee()
            })

            expect(res.code).not.toBe(0)
            expect(res.rawLog).toMatch(/pair already exists/);
        });

        test('should be able to create pool for existing pair', async () => {

            const a = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePool(
                {
                    value: {
                        creator: sdk.adminAddress,
                        pairId: pairId,
                        depositCoins: [
                            {
                                denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                                amount: "1000000000"
                            },
                            {
                                denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                                amount: "1000000000"
                            }
                        ]
                    },
                    fee: getGasFee()
                }
            )

            const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPools({
                pair_id: String(pairId)
            });

            const lastPool = resp.data.pools.pop();

            poolId = lastPool.id;

            expect(lastPool.creator).toBe(sdk.adminAddress)
        })

        test('should be able to deposit liquidity to existing pool', async () => {
            const balanceOfBaseCoinsBefore = await sdk.clientAdmin.CosmosBankV1Beta1.query.queryBalance(
                sdk.adminAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )

            await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgDeposit(
                {
                    value: {
                        depositor: sdk.adminAddress,
                        poolId: poolId,
                        depositCoins: [
                            {
                                denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                                amount: "1000000000"
                            },
                            {
                                denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                                amount: "1000000000"
                            }
                        ]
                    },
                    fee: getGasFee()
                }
            )

            const balanceOfBaseCoinsAfter = await sdk.clientAdmin.CosmosBankV1Beta1.query.queryBalance(
                sdk.adminAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )
            expect(Number(balanceOfBaseCoinsBefore.data.balance.amount)).toBeGreaterThan(Number(balanceOfBaseCoinsAfter.data.balance.amount))
        })

        test('should be able to withdraw liquidity from existing pool', async () => {
            const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPool(poolId);
            const pool = resp.data.pool;

            const balanceOfBaseCoinsBefore = await sdk.clientAdmin.CosmosBankV1Beta1.query.queryBalance(
                sdk.adminAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )

            await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgWithdraw(
                {
                    value: {
                        withdrawer: sdk.adminAddress,
                        poolId: poolId,
                        poolCoin: {
                            denom: pool.pool_coin_denom,
                            amount: "5000000"
                        }
                    },
                    fee: getGasFee()
                }
            )

            const balanceOfBaseCoinsAfter = await sdk.clientAdmin.CosmosBankV1Beta1.query.queryBalance(
                sdk.adminAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )
            expect(Number(balanceOfBaseCoinsBefore.data.balance.amount)).toBeLessThan(Number(balanceOfBaseCoinsAfter.data.balance.amount))
        })

        test('should be able to create rangedPool for existing pair', async () => {
            await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreateRangedPool(
                {
                    value: {
                        creator: sdk.adminAddress,
                        pairId: pairId,
                        depositCoins: [
                            {
                                denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                                amount: "1000000000"
                            },
                            {
                                denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                                amount: "1000000000"
                            }
                        ],
                        minPrice: '1400000000000000000',
                        maxPrice: '1600000000000000000',
                        initialPrice: '1500000000000000000',
                    },
                    fee: getGasFee()
                }
            )

            const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPools({
                pair_id: pairId.toString()
            });

            const lastPool = resp.data.pools.pop();

            expect(lastPool.type).toBe('POOL_TYPE_RANGED')
            expect(lastPool.creator).toBe(sdk.adminAddress)
        })
    })

    describe('User', () => {
        test('should throw when trying to create pairs', async () => {
            const res = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
                value: {
                    creator: sdk.recipientAddress,
                    baseCoinDenom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                    quoteCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom)
                },
                fee: getGasFee()
            })

            expect(res.code).not.toBe(0)
        })

        test('should throw when trying to create pools', async () => {

            const res = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePool(
                {
                    value: {
                        creator: sdk.recipientAddress,
                        pairId: pairId,
                        depositCoins: [
                            {
                                denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                                amount: "100000000"
                            },
                            {
                                denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                                amount: "100000000"
                            }
                        ]
                    },
                    fee: getGasFee()
                }
            )

            expect(res.code).not.toBe(0)
        })

        test('should throw when trying to create rangedPools', async () => {
            const res = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreateRangedPool(
                {
                    value: {
                        creator: sdk.recipientAddress,
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
                        ],
                        minPrice: '1400000000000000000',
                        maxPrice: '1600000000000000000',
                        initialPrice: '1500000000000000000',
                    },
                    fee: getGasFee()
                }
            )

            expect(res.code).not.toBe(0)
        })

        test('should be able to deposit liquidity to existing pool', async () => {

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

            // await updateCoinRequiredPrivileges(
            //     sdk,
            //     sdk.clientAdmin,
            //     sdk.adminAddress,
            //     genCoinDenom(sdk.adminAddress, baseCoinDenom),
            //     [1, 1]
            // )
            // await updateCoinRequiredPrivileges(
            //     sdk,
            //     sdk.clientAdmin,
            //     sdk.adminAddress,
            //     genCoinDenom(sdk.adminAddress, quoteCoinDenom),
            //     [1, 1]
            // )

            const bob = await sdk.clientAdmin.MantrachainTokenV1.tx.sendMsgMintNft({
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

            const prv = await sdk.clientAdmin.MantrachainGuardV1.tx.sendMsgUpdateAccountPrivileges({
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
                            amount: '10000000000000'
                        },
                        {
                            denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                            amount: '10000000000000'
                        }
                    ]
                },
                fee: getGasFee()
            })

            const balanceOfBaseCoinsBefore = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )

            const a = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgDeposit({
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

            const balanceOfBaseCoinsAfter = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )
            expect(Number(balanceOfBaseCoinsBefore.data.balance.amount)).toBeGreaterThan(Number(balanceOfBaseCoinsAfter.data.balance.amount))
        })

        test('should be able to withdraw liquidity from existing pool after deposit', async () => {
            const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPool(poolId);

            const pool = resp.data.pool;

            const balanceOfBaseCoinsBefore = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {
                    denom: pool.pool_coin_denom
                }
            )

            await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgWithdraw(
                {
                    value: {
                        withdrawer: sdk.recipientAddress,
                        poolId: poolId,
                        poolCoin: {
                            denom: pool.pool_coin_denom,
                            amount: "200000000"
                        }
                    },
                    fee: getGasFee()
                }
            )

            const balanceOfBaseCoinsAfter = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {
                    denom: pool.pool_coin_denom
                }
            )

            expect(Number(balanceOfBaseCoinsBefore.data.balance.amount)).toBeGreaterThan(Number(balanceOfBaseCoinsAfter.data.balance.amount))
        })

        test('should be able to make limit order', async () => {
            const balanceOfBaseCoinsBefore = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )

            await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgLimitOrder({
                    value: {
                        orderer: sdk.recipientAddress,
                        pairId: pairId,
                        direction: OrderDirection.ORDER_DIRECTION_SELL,
                        offerCoin: {
                            denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                            amount: "1000000"
                        },
                        demandCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                        price: '1400000000000000000',
                        amount: '1000000',
                        orderLifespan: undefined
                    },
                    fee: getGasFee()
                }
            )

            const balanceOfBaseCoinsAfter = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )

            expect(Number(balanceOfBaseCoinsBefore.data.balance.amount)).toBeGreaterThan(Number(balanceOfBaseCoinsAfter.data.balance.amount))
        })

        test('should be able to make market order', async () => {
            const balanceOfBaseCoinsBefore = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )

            const z = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgMarketOrder(
                {
                    value: {
                        orderer: sdk.recipientAddress,
                        pairId: pairId,
                        direction: OrderDirection.ORDER_DIRECTION_SELL,
                        offerCoin: {
                            denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                            amount: "1000000"
                        },
                        demandCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                        amount: '1000000',
                        orderLifespan: undefined
                    },
                    fee: getGasFee()
                }
            )

            const balanceOfBaseCoinsAfter = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
                sdk.recipientAddress, {
                    denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
                }
            )

            expect(Number(balanceOfBaseCoinsBefore.data.balance.amount)).toBeGreaterThan(Number(balanceOfBaseCoinsAfter.data.balance.amount))
        })

        // test('should be able to make MM order', async () => {
        //     const balanceOfBaseCoinsBefore = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
        //         sdk.recipientAddress, {
        //             denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
        //         }
        //     )
        //
        //     const a = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgMMOrder(
        //         {
        //             value: {
        //                 orderer: sdk.recipientAddress,
        //                 pairId: Number(pairId),
        //                 maxSellPrice: '1500000000000000000',
        //                 minSellPrice: '1400000000000000000',
        //                 sellAmount: '1000000',
        //                 maxBuyPrice:'1500000000000000000',
        //                 minBuyPrice:'1400000000000000000',
        //                 buyAmount: '1000000'
        //             },
        //             fee: {
        //                 amount: [{ denom: "uaum", amount: "10000000" }],
        //                 gas: "50000000000"
        //             }
        //         }
        //     )
        //     console.log(a)
        //     const balanceOfBaseCoinsAfter = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
        //         sdk.recipientAddress, {
        //             denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
        //         }
        //     )
        //
        //     expect(Number(balanceOfBaseCoinsBefore.data.balance.amount)).toBeGreaterThan(Number(balanceOfBaseCoinsAfter.data.balance.amount))
        // })

        test('should be able to cancel order', async () => {

            const z = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgLimitOrder({
                    value: {
                        orderer: sdk.recipientAddress,
                        pairId: pairId,
                        direction: OrderDirection.ORDER_DIRECTION_SELL,
                        offerCoin: {
                            denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                            amount: "1000000"
                        },
                        demandCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                        price: '1500000000000000000',
                        amount: '1000000',
                        orderLifespan: {seconds: 120}
                    },
                    fee: getGasFee()
                }
            )
            console.log(z)
            const orders = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.query.queryOrders(pairId.toString())
            console.log(orders)
            // const order = orders.data.orders.find(o => o.type == 'ORDER_TYPE_MM')
            const order = orders.data.orders.pop();

            await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCancelOrder({
                value: {
                    orderer: sdk.recipientAddress,
                    pairId: pairId,
                    orderId: order.id
                },
                fee: getGasFee()
            })

            await new Promise((r) => setTimeout(r, 7000));

            const ordersAfter = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.query.queryOrders(pairId.toString())
            expect(ordersAfter.data.orders.length).toBe(0)
        })

        test('should be able to cancel all orders', async () => {
            await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgLimitOrder(
                {
                    value: {
                        orderer: sdk.recipientAddress,
                        pairId: pairId,
                        direction: OrderDirection.ORDER_DIRECTION_SELL,
                        offerCoin: {
                            denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                            amount: "1000000"
                        },
                        demandCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                        price: '1500000000000000000',
                        amount: '1000000',
                        orderLifespan: {seconds: 120}
                    },
                    fee: getGasFee()
                }
            )

            await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgMarketOrder(
                {
                    value: {
                        orderer: sdk.recipientAddress,
                        pairId: pairId,
                        direction: OrderDirection.ORDER_DIRECTION_SELL,
                        offerCoin: {
                            denom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
                            amount: "1000000"
                        },
                        demandCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
                        amount: '1000000',
                        orderLifespan: {seconds: 120}
                    },
                    fee: getGasFee()
                }
            )

            await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCancelAllOrders({
                value: {
                    orderer: sdk.recipientAddress
                },
                fee: getGasFee()
            })

            await new Promise((r) => setTimeout(r, 7000));

            const ordersAfter = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.query.queryOrders(pairId.toString())

            expect(ordersAfter.data.orders.length).toBe(0)
        })
    })
})