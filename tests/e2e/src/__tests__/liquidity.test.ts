import { MantrachainSdk } from "../helpers/sdk"
import { createDenomIfNotExists, genCoinDenom } from "../helpers/coinfactory"
import { mintGuardSoulBondNft } from '../helpers/token'
import { updateCoinRequiredPrivileges, updateAccountPrivileges } from '../helpers/guard'

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

  let baseCoinDenom = 'atom' + new Date().getTime().toString()
  // with this we manipulate the time and space
  // of this test environment according to our needs
  // to the infinity and beyond!
  let quoteCoinDenom = 'osmo' + new Date().getTime().toString() + 1

  let pairId = 0
  let poolId = 0

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
      }
    })

    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom: genCoinDenom(sdk.adminAddress, quoteCoinDenom),
          amount: "1000000000000000000"
        }
      }
    })
  })

  describe('Admin', () => {
    test('should be able to create pair for existing denoms', async () => {
      await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
        value: {
          creator: sdk.adminAddress,
          baseCoinDenom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
          quoteCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom)
        }
      })

      const allPairs = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPairs()
      const lastPair = allPairs.data.pairs.pop()
      pairId = lastPair.id

      expect(Number(lastPair.id)).toBeGreaterThan(0)
    })

    test('should throw when trying to create already existing pair for existing denoms', async () => {
      await expect(sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
        value: {
          creator: sdk.adminAddress,
          baseCoinDenom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
          quoteCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom)
        }
      })).rejects.toThrow(/pair already exists/)
    })

    test('should be able to create pool for existing pair', async () => {
      await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePool(
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
          }
        }
      )

      const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPools({
        pair_id: String(pairId)
      })

      const lastPool = resp.data.pools.pop()

      poolId = lastPool.id

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
          }
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
      const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPool(poolId)
      const pool = resp.data.pool

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
          }
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
          }
        }
      )

      const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPools({
        pair_id: pairId.toString()
      })

      const lastPool = resp.data.pools.pop()

      expect(lastPool.type).toBe('POOL_TYPE_RANGED')
      expect(lastPool.creator).toBe(sdk.adminAddress)
    })
  })

  describe('User', () => {
    test('should throw when trying to create pairs', async () => {
      await expect(sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
        value: {
          creator: sdk.recipientAddress,
          baseCoinDenom: genCoinDenom(sdk.adminAddress, baseCoinDenom),
          quoteCoinDenom: genCoinDenom(sdk.adminAddress, quoteCoinDenom)
        }
      })).rejects.toThrow(/unauthorized/)
    })

    test('should throw when trying to create pools', async () => {
      const res = sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePool(
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
          }
        }
      )

      await expect(res).rejects.toThrow(/unauthorized/)
    })

    test('should throw when trying to create rangedPools', async () => {
      const res = sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreateRangedPool(
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
          }
        }
      )

      await expect(res).rejects.toThrow(/unauthorized/)
    })

    test('should be able to deposit liquidity to existing pool', async () => {
      await updateCoinRequiredPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, genCoinDenom(sdk.adminAddress, baseCoinDenom), [64])
      await updateCoinRequiredPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, genCoinDenom(sdk.adminAddress, quoteCoinDenom), [64])
      await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress, [64])
      await mintGuardSoulBondNft(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress)

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
        }
      })

      const balanceOfBaseCoinsBefore = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
        sdk.recipientAddress, {
        denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
      }
      )

      await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgDeposit({
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
        }
      }
      )

      const balanceOfBaseCoinsAfter = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
        sdk.recipientAddress, {
        denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
      })

      expect(Number(balanceOfBaseCoinsBefore.data.balance.amount)).toBeGreaterThan(Number(balanceOfBaseCoinsAfter.data.balance.amount))
    })

    test('should be able to withdraw liquidity from existing pool after deposit', async () => {
      const resp = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.query.queryPool(poolId)

      const pool = resp.data.pool

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
          }
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
        }
      })

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
            orderLifespan: undefined
          }
        }
      )

      const balanceOfBaseCoinsAfter = await sdk.clientRecipient.CosmosBankV1Beta1.query.queryBalance(
        sdk.recipientAddress, {
        denom: genCoinDenom(sdk.adminAddress, baseCoinDenom)
      }
      )

      expect(Number(balanceOfBaseCoinsBefore.data.balance.amount)).toBeGreaterThan(Number(balanceOfBaseCoinsAfter.data.balance.amount))
    })

    test('should be able to cancel order', async () => {

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
          price: '1500000000000000000',
          amount: '1000000',
          orderLifespan: { seconds: 120 }
        }
      }
      )
      const orders = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.query.queryOrders(pairId.toString())

      // const order = orders.data.orders.find(o => o.type == 'ORDER_TYPE_MM')
      const order = orders.data.orders.pop()

      await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCancelOrder({
        value: {
          orderer: sdk.recipientAddress,
          pairId: pairId,
          orderId: order.id
        }
      })

      await new Promise((r) => setTimeout(r, 7000))

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
            orderLifespan: { seconds: 120 }
          }
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
            orderLifespan: { seconds: 120 }
          }
        }
      )

      await sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCancelAllOrders({
        value: {
          orderer: sdk.recipientAddress
        }
      })

      await new Promise((r) => setTimeout(r, 7000))

      const ordersAfter = await sdk.clientRecipient.MantrachainLiquidityV1Beta1.query.queryOrders(pairId.toString())

      expect(ordersAfter.data.orders.length).toBe(0)
    })
  })
})