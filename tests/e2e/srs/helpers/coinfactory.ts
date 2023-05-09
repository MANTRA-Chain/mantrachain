import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'
import { queryBalance } from './bank'

const queryDenom = (client: any, account: string): any => client.MantrachainCoinfactoryV1Beta1.query.queryDenomsFromCreator(account)

const notExistsDenom = (res: any, account: string, subdenom: string) => res.data?.denoms?.every((denom: string) => denom !== getCoinDenom(account, subdenom))

const existsDenom = (res: any, account: string, subdenom: string) => res.data?.denoms?.some((denom: string) => denom === getCoinDenom(account, subdenom))

export const getCoinDenom = (account: string, subdenom: string) => `factory/${account}/${subdenom}`

export const createDenomIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, subdenom: string, numAttempts = 20) => {
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
    numAttempts,
  )
}

export const mintCoins = async (sdk: MantrachainSdk, client: any, account: string, subdenom: string, amount: string, minBalance?: string, numAttempts = 20) => {
  const denom = getCoinDenom(account, subdenom)
  const privBalance = await queryBalance(client, account, denom)

  if (!!minBalance && parseInt(privBalance?.data?.balance?.amount) >= parseInt(minBalance)) {
    return privBalance
  }

  await client.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
    value: {
      sender: account,
      amount: {
        denom,
        amount
      }
    }
  })

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryBalance(client, account, denom),
    async (res) => parseInt(res?.data?.balance?.amount) === parseInt(privBalance?.data?.balance?.amount) + parseInt(amount),
    numAttempts,
  )
}