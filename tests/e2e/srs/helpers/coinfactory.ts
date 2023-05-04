import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

const queryDenom = (client: any, account: string): any => client.MantrachainCoinfactoryV1Beta1.query.queryDenomsFromCreator(account)

const notExistsDenom = (res: any, account: string, subdenom: string) => res.data?.denoms?.every((denom: string) => denom !== `factory/${account}/${subdenom}`)

const existsDenom = (res: any, account: string, subdenom: string) => res.data?.denoms?.some((denom: string) => denom === `factory/${account}/${subdenom}`)

export const createDenomIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, subdenom: string) => {
  if (notExistsDenom(await queryDenom(client, account), account, subdenom)) {
    await client.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: account,
        subdenom
      }
    })
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryDenom(client, account),
    async (res) => existsDenom(res, account, subdenom),
    20,
  )
}