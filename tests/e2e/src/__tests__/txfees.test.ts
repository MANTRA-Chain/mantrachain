import { MantrachainSdk } from '../helpers/sdk'
import { createDenomIfNotExists, genCoinDenom } from "../helpers/coinfactory"
import { createPairIfNotExists, createPoolIfNotExists, getPairId, getPoolId } from "../helpers/liquidity"
import { queryBalance, sendCoins } from '../helpers/bank'

describe('Txfees module', () => {
  let sdk: MantrachainSdk

  let testDenom = 'txfees' + new Date().getTime().toString()
  let gasFeesDenom = 'txfees' + new Date().getTime().toString() + 1

  let pairId = 0

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
    await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, gasFeesDenom)

    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom: genCoinDenom(sdk.adminAddress, gasFeesDenom),
          amount: "1000000000000000000"
        }
      }
    })

    await createPairIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom), "uaum")
    pairId = await getPairId(sdk.clientAdmin, genCoinDenom(sdk.adminAddress, gasFeesDenom), "uaum")
    await createPoolIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, String(pairId), genCoinDenom(sdk.adminAddress, gasFeesDenom), "100000000", "uaum", "10000000")

    // To set the last price of the pair
    const res = await sdk.clientAdmin.MantrachainLiquidityV1Beta1.tx.sendMsgLimitOrder({
      value: {
        orderer: sdk.adminAddress,
        pairId: pairId,
        direction: 1,
        offerCoin: {
          denom: "uaum",
          amount: "1000000"
        },
        demandCoinDenom: genCoinDenom(sdk.adminAddress, gasFeesDenom),
        price: '140000000000000000',
        amount: '1000000',
        orderLifespan: undefined
      }
    })

    expect(res.code).toBe(0)

    await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom), 10000)
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
          denom: genCoinDenom(sdk.adminAddress, gasFeesDenom),
          amount: "1000"
        }],
        gas: "100000"
      }
    })).rejects.toThrow(
      /invalid fee denom/
    )
  })

  test('Should pay gas fees with non-native token', async () => {
    let res = await sdk.clientAdmin.MantrachainTxfeesV1.tx.sendMsgCreateFeeToken({
      value: {
        creator: sdk.adminAddress,
        denom: genCoinDenom(sdk.adminAddress, gasFeesDenom),
        pairId: pairId.toString(),
      }
    })

    expect(res.code).toBe(0)

    const currNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, "uaum")
    const currNonNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom))

    res = await sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgCreateNftCollection({
      value: {
        creator: sdk.recipientAddress,
        collection: {
          id: 'txfees' + new Date().getTime().toString(),
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
          denom: genCoinDenom(sdk.adminAddress, gasFeesDenom),
          amount: "1000"
        }],
        gas: "200000"
      }
    })

    expect(res.code).toBe(0)

    const latestNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, "uaum")
    const latestNonNativeBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, genCoinDenom(sdk.adminAddress, gasFeesDenom))

    expect(currNativeBalance).toEqual(latestNativeBalance);
    expect(currNonNativeBalance).toBeGreaterThan(latestNonNativeBalance);
  })
})
