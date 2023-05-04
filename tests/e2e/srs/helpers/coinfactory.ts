import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

export const createDenomIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, denom: string) => {
  const queryDenom = (): any => client.MantrachainCoinfactoryV1Beta1.query.queryDenomsFromCreator(account)
  const notExistsDenom = (res: any) => res.data?.denoms?.every((denom: string) => !denom.includes(denom))
  const existsDenom = (res: any) => res.data?.denoms?.some((denom: string) => denom.includes(denom))

  const res = await queryDenom()

  if (notExistsDenom(res)) {
    await client.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: account,
        subdenom: denom
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