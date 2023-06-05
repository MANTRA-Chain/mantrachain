import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'
import { getGasFee } from './sdk'

const existsPair = (pairs: any[], baseCoinDenom: string, quoteCoinDenom: string) => pairs.some((pair: any) => pair.base_coin_denom === baseCoinDenom && pair.quote_coin_denom === quoteCoinDenom)

const notExistsPair = (pairs: any[], baseCoinDenom: string, quoteCoinDenom: string) => pairs.every((pair: any) => pair.base_coin_denom !== baseCoinDenom && pair.quote_coin_denom !== quoteCoinDenom)

const queryPairs = async (client: any, baseCoinDenom: string, quoteCoinDenom: string) => {
  const res = await client.MantrachainLiquidityV1Beta1.query.queryPairsByDenoms({
    denom1: baseCoinDenom,
    denom2: quoteCoinDenom,
  })
  return res?.data?.pairs || []
}

const getPair = (pairs: any[], baseCoinDenom: string, quoteCoinDenom: string) => pairs.find((pair: any) => pair.base_coin_denom === baseCoinDenom && pair.quote_coin_denom === quoteCoinDenom)

export const createPairIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, baseCoinDenom: string, quoteCoinDenom: string, numAttempts = 2) => {
  if (notExistsPair(await queryPairs(client, baseCoinDenom, quoteCoinDenom), baseCoinDenom, quoteCoinDenom)) {
    const res = await client.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
      value: {
        creator: account,
        baseCoinDenom,
        quoteCoinDenom
      },
      fee: getGasFee()
    })

    if (res.code !== 0) {
      throw new Error(res.rawLog)
    }
  } else {
    return
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryPairs(client, baseCoinDenom, quoteCoinDenom),
    async (res) => existsPair(res, baseCoinDenom, quoteCoinDenom),
    numAttempts,
  )
}

export const getPairId = async (client: any, baseCoinDenom: string, quoteCoinDenom: string): Promise<number> => {
  const res = await queryPairs(client, baseCoinDenom, quoteCoinDenom)

  if (notExistsPair(res, baseCoinDenom, quoteCoinDenom)) {
    throw new Error(`Pair ${baseCoinDenom}:${quoteCoinDenom} does not exist`)
  }

  return parseInt(getPair(res, baseCoinDenom, quoteCoinDenom)['id'])
}