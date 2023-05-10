import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

const existsPair = (pairs: any[], baseCoinDenom: string, quoteCoinDenom: string) => pairs.some((pair: any) => pair.base_coin_denom === baseCoinDenom && pair.quote_coin_denom === quoteCoinDenom)

const notExistsPair = (pairs: any[], baseCoinDenom: string, quoteCoinDenom: string) => pairs.every((pair: any) => pair.base_coin_denom !== baseCoinDenom && pair.quote_coin_denom !== quoteCoinDenom)

const queryPairs = async (client: any, baseCoinDenom: string) => {
  const res = await client.MantrachainLiquidityV1Beta1.query.queryPairs({
    denoms: baseCoinDenom // This doesn't work properly, it should be [baseCoinDenom, quoteCoinDenom], 
    // but the chain parse the denoms query as ["", baseCoinDenom, quoteCoinDenom] if we pass an array
  })
  return res?.data?.pairs || []
}

const getPair = (pairs: any[], baseCoinDenom: string, quoteCoinDenom: string) => pairs.find((pair: any) => pair.base_coin_denom === baseCoinDenom && pair.quote_coin_denom === quoteCoinDenom)

export const createPairIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, baseCoinDenom: string, quoteCoinDenom: string, numAttempts = 20) => {
  const res = await queryPairs(client, baseCoinDenom)

  if (notExistsPair(res, baseCoinDenom, quoteCoinDenom)) {
    await client.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
      value: {
        creator: account,
        baseCoinDenom,
        quoteCoinDenom
      }
    })
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryPairs(client, baseCoinDenom),
    async (res) => existsPair(res, baseCoinDenom, quoteCoinDenom),
    numAttempts,
  )
}

export const getPairId = async (client: any, baseCoinDenom: string, quoteCoinDenom: string): Promise<number> => {
  const res = await queryPairs(client, baseCoinDenom)

  if (notExistsPair(res, baseCoinDenom, quoteCoinDenom)) {
    throw new Error(`Pair ${baseCoinDenom}:${quoteCoinDenom} does not exist`)
  }

  return parseInt(getPair(res, baseCoinDenom, quoteCoinDenom)['id'])
}