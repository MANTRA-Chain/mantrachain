import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

export const createDenomIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, subdenom: string) => {
  const queryDenom = (): any => client.MantrachainCoinfactoryV1Beta1.query.queryDenomsFromCreator(account)
  const notExistsDenom = (res: any) => res.data?.denoms?.every((denom: string) => denom !== `factory/${account}/${subdenom}`)
  const existsDenom = (res: any) => res.data?.denoms?.some((denom: string) => denom === `factory/${account}/${subdenom}`)

  if (notExistsDenom(await queryDenom())) {
    await client.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: account,
        subdenom
      }
    })
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryDenom(),
    async (res) => existsDenom(res),
    20,
  );
}