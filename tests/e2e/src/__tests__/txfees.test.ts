import { MantrachainSdk } from '../helpers/sdk'
import { createDenomIfNotExists, genCoinDenom } from "../helpers/coinfactory"
import { createPairIfNotExists, createPoolIfNotExists, getPairId, getPoolId } from "../helpers/liquidity"
import { queryBalance, sendCoins } from '../helpers/bank'

describe('Txfees module', () => {
  let sdk: MantrachainSdk

  let testDenom = 'txfees' + new Date().getTime().toString()
  let gasFeesDenom1 = 'txfees' + new Date().getTime().toString() + 1
  let gasFeesDenom2 = 'txfees' + new Date().getTime().toString() + 2

  let pairId1 = 0
  let pairId2 = 0

  beforeAll(async () => {
    sdk = new MantrachainSdk()
    await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)

    await sdk.clientAdmin.MantrachainGuardV1.tx.sendMsgUpdateGuardTransferCoins({
      value: {
        creator: sdk.adminAddress,
        enabled: false
      }
    })

    await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, testDenom)
    await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, gasFeesDenom1)
    await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, gasFeesDenom2)

    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom: genCoinDenom(sdk.adminAddress, gasFeesDenom1),
          amount: "1000000000000000000"
        }
      }
    })

    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom: genCoinDenom(sdk.adminAddress, gasFeesDenom2),
          amount: "1000000000000000000"
        }
      }
    })

    await createPairIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom1), "uaum")
    pairId1 = await getPairId(sdk.clientAdmin, genCoinDenom(sdk.adminAddress, gasFeesDenom1), "uaum")
    await createPoolIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, String(pairId1), genCoinDenom(sdk.adminAddress, gasFeesDenom1), "100000000", "uaum", "10000000")

    await createPairIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, "uaum", genCoinDenom(sdk.adminAddress, gasFeesDenom2))
    pairId2 = await getPairId(sdk.clientAdmin, "uaum", genCoinDenom(sdk.adminAddress, gasFeesDenom2))
    await createPoolIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, String(pairId2), genCoinDenom(sdk.adminAddress, gasFeesDenom2), "100000000", "uaum", "10000000")

    // To set the last price of the pair
    let res = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgLimitOrder({
      value: {
        orderer: sdk.adminAddress,
        pairId: pairId1,
        direction: 1,
        offerCoin: {
          denom: "uaum",
          amount: "1000000"
        },
        demandCoinDenom: genCoinDenom(sdk.adminAddress, gasFeesDenom1),
        price: '140000000000000000',
        amount: '1000000',
        orderLifespan: undefined
      }
    })

    expect(res.code).toBe(0)

    // To set the last price of the pair
    res = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgLimitOrder({
      value: {
        orderer: sdk.adminAddress,
        pairId: pairId2,
        direction: 2,
        offerCoin: {
          denom: "uaum",
          amount: "1000000"
        },
        demandCoinDenom: genCoinDenom(sdk.adminAddress, gasFeesDenom2),
        price: '140000000000000000',
        amount: '1000000',
        orderLifespan: undefined
      }
    })

    expect(res.code).toBe(0)

    await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom1), 10000)
    await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom2), 10000)
  })

  afterAll(async () => {
    await sdk.clientAdmin.MantrachainGuardV1.tx.sendMsgUpdateGuardTransferCoins({
      value: {
        creator: sdk.adminAddress,
        enabled: true
      }
    })
  })

  test('Should return error when try to pay gas fees with non-native token', async () => {
    await expect(sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom: genCoinDenom(sdk.adminAddress, testDenom),
          amount: "1000"
        }
      },
      fee: {
        amount: [{
          denom: genCoinDenom(sdk.adminAddress, gasFeesDenom1),
          amount: "1000"
        }],
        gas: "100000"
      }
    })).rejects.toThrow(
      /invalid fee denom/
    )
  })

  test('Should pay gas fees with non-native token for pair: non native coin/native coin', async () => {
    let res = await sdk.clientAdmin.MantrachainTxfeesV1.tx.sendMsgCreateFeeToken({
      value: {
        creator: sdk.adminAddress,
        denom: genCoinDenom(sdk.adminAddress, gasFeesDenom1),
        pairId: pairId1.toString(),
      }
    })

    expect(res.code).toBe(0)

    const currNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, "uaum")
    const currNonNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom1))

    res = await sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgCreateNftCollection({
      value: {
        creator: sdk.recipientAddress,
        collection: {
          id: 'txfees0' + new Date().getTime().toString(),
          name: 'txfees test collection',
          images: [],
          url: '',
          description: '',
          links: [],
          options: [],
          category: 'utility',
          opened: false,
          symbol: 'TEST',
          soulBondedNfts: false,
          restrictedNfts: false,
          data: null
        },
      },
      fee: {
        amount: [{
          denom: genCoinDenom(sdk.adminAddress, gasFeesDenom1),
          amount: "1000"
        }],
        gas: "200000"
      }
    })

    expect(res.code).toBe(0)

    const latestNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, "uaum")
    const latestNonNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom1))

    expect(currNativeBalance).toEqual(latestNativeBalance);
    expect(currNonNativeBalance).toBeGreaterThan(latestNonNativeBalance);
  })

  test('Should pay gas fees with non-native token for pair: native coin/non native coin', async () => {
    let res = await sdk.clientAdmin.MantrachainTxfeesV1.tx.sendMsgCreateFeeToken({
      value: {
        creator: sdk.adminAddress,
        denom: genCoinDenom(sdk.adminAddress, gasFeesDenom2),
        pairId: pairId2.toString(),
      }
    })

    expect(res.code).toBe(0)

    const currNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, "uaum")
    const currNonNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom2))

    res = await sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgCreateNftCollection({
      value: {
        creator: sdk.recipientAddress,
        collection: {
          id: 'txfees1' + new Date().getTime().toString(),
          name: 'txfees test collection',
          images: [],
          url: '',
          description: '',
          links: [],
          options: [],
          category: 'utility',
          opened: false,
          symbol: 'TEST',
          soulBondedNfts: false,
          restrictedNfts: false,
          data: null
        },
      },
      fee: {
        amount: [{
          denom: genCoinDenom(sdk.adminAddress, gasFeesDenom2),
          amount: "1000"
        }],
        gas: "200000"
      }
    })

    expect(res.code).toBe(0)

    const latestNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, "uaum")
    const latestNonNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom2))

    expect(currNativeBalance).toEqual(latestNativeBalance);
    expect(currNonNativeBalance).toBeGreaterThan(latestNonNativeBalance);
  })
})
