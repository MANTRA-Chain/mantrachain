import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

export const createPairIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, baseCoinDenom: string, quoteCoinDenom: string) => {
  const queryPairs = (): any => client.MantrachainLiquidityV1Beta1.query.queryPairs({
    denoms: baseCoinDenom
  })
  const notExistsPair = (res: any) => res.data?.pairs?.every((pair: any) => pair.base_coin_denom !== baseCoinDenom && pair.quote_coin_denom !== quoteCoinDenom)
  const existsPair = (res: any) => res.data?.pairs?.some((pair: any) => pair.base_coin_denom === baseCoinDenom && pair.quote_coin_denom === quoteCoinDenom)

  if (notExistsPair(await queryPairs())) {
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
    async () => await queryPairs(),
    async (res) => existsPair(res),
    20,
  );
}