import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

const existsPair = (res: any, baseCoinDenom: string, quoteCoinDenom: string) => res.data?.pairs?.some((pair: any) => pair.base_coin_denom === baseCoinDenom && pair.quote_coin_denom === quoteCoinDenom)

const notExistsPair = (res: any, baseCoinDenom: string, quoteCoinDenom: string) => res.data?.pairs?.every((pair: any) => pair.base_coin_denom !== baseCoinDenom && pair.quote_coin_denom !== quoteCoinDenom)

const queryPairs = (client: any, baseCoinDenom): any => client.MantrachainLiquidityV1Beta1.query.queryPairs({
  denoms: baseCoinDenom // This doesn't work properly, it should be [baseCoinDenom, quoteCoinDenom], 
  // but the chain parse the request as ["", baseCoinDenom, quoteCoinDenom] if we pass an array
})

const getPair = (res: any, baseCoinDenom: string, quoteCoinDenom: string) => res.data?.pairs?.find((pair: any) => pair.base_coin_denom === baseCoinDenom && pair.quote_coin_denom === quoteCoinDenom)


export const createPairIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, baseCoinDenom: string, quoteCoinDenom: string) => {
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
    20,
  )
}

export const getPairId = async (client: any, baseCoinDenom: string, quoteCoinDenom: string): Promise<number> => {
  const res = await queryPairs(client, baseCoinDenom)

  if (notExistsPair(res, baseCoinDenom, quoteCoinDenom)) {
    throw new Error(`Pair ${baseCoinDenom}:${quoteCoinDenom} does not exist`)
  }

  return parseInt(getPair(res, baseCoinDenom, quoteCoinDenom)['id'])
}