import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

const getPair = (pairs: any[], baseCoinDenom: string, quoteCoinDenom: string) => pairs.find((pair: any) => pair.base_coin_denom === baseCoinDenom && pair.quote_coin_denom === quoteCoinDenom)

const existsPair = (pairs: any[], baseCoinDenom: string, quoteCoinDenom: string) => pairs.some((pair: any) => pair.base_coin_denom === baseCoinDenom && pair.quote_coin_denom === quoteCoinDenom)

const notExistsPair = (pairs: any[], baseCoinDenom: string, quoteCoinDenom: string) => pairs.every((pair: any) => pair.base_coin_denom !== baseCoinDenom || pair.quote_coin_denom !== quoteCoinDenom)

const queryPairs = async (client: any) => {
  const res = await client.MantrachainLiquidityV1Beta1.query.queryPairs()
  return res?.data?.pairs || []
}

export const createPairIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, baseCoinDenom: string, quoteCoinDenom: string, numAttempts = 2) => {
  if (notExistsPair(await queryPairs(client), baseCoinDenom, quoteCoinDenom)) {
    const res = await client.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
      value: {
        creator: account,
        baseCoinDenom,
        quoteCoinDenom
      }
    })
    
    if (res.code !== 0) {
      throw new Error(res.rawLog)
    }
  } else {
    return
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryPairs(client),
    async (res) => existsPair(res, baseCoinDenom, quoteCoinDenom),
    numAttempts,
  )
}

export const getPairId = async (client: any, baseCoinDenom: string, quoteCoinDenom: string): Promise<number> => {
  const pairs = await queryPairs(client)
  
  if (notExistsPair(pairs, baseCoinDenom, quoteCoinDenom)) {
    throw new Error(`Pair ${baseCoinDenom}:${quoteCoinDenom} does not exist`)
  }

  return parseInt(getPair(pairs, baseCoinDenom, quoteCoinDenom)['id'])
}

const getPool = (pools: any[], pairId: string) => pools.find((pair: any) => pair.pair_id === pairId)

const existsPool = (pools: any[], pairId: string) => pools.some((pair: any) => pair.pair_id === pairId)

const notExistsPool = (pools: any[], pairId: string) => pools.every((pair: any) => pair.pair_id !== pairId)

const queryPools = async (client: any, pairId: string) => {
  const res = await client.MantrachainLiquidityV1Beta1.query.queryPools({
    pair_id: pairId
  })
  return res?.data?.pools || []
}

export const createPoolIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, pairId: string, baseCoinDenom: string, baseCoinAmount: string, quoteCoinDenom: string, quoteCoinAmount: string, numAttempts = 2) => {
  if (notExistsPool(await queryPools(client, pairId), pairId)) {
    const res = await client.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePool({
      value: {
        creator: account,
        pairId,
        depositCoins: [
          {
            denom: baseCoinDenom,
            amount: baseCoinAmount
          },
          {
            denom: quoteCoinDenom,
            amount: quoteCoinAmount
          }
        ]
      }
    })

    if (res.code !== 0) {
      throw new Error(res.rawLog)
    }
  } else {
    return
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryPools(client, pairId),
    async (res) => existsPool(res, pairId),
    numAttempts,
  )
}

export const getPoolId = async (client: any, pairId: string): Promise<number> => {
  const res = await queryPools(client, pairId)

  if (notExistsPool(res, pairId)) {
    throw new Error(`Pool pair id: ${pairId} does not exist`)
  }

  return parseInt(getPool(res, pairId)['id'])
}