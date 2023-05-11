import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'
import { queryBalance } from './bank'

const queryDenomsFromCreator = async (client: any, account: string) => {
  const res = await client.MantrachainCoinfactoryV1Beta1.query.queryDenomsFromCreator(account)
  return res?.data?.denoms || []
}

const notExistsDenom = (denoms: string[], account: string, subdenom: string) => denoms.every((denom: string) => denom !== genCoinDenom(account, subdenom))

const existsDenom = (denoms: string[], account: string, subdenom: string) => denoms.some((denom: string) => denom === genCoinDenom(account, subdenom))

export const genCoinDenom = (account: string, subdenom: string) => `factory/${account}/${subdenom}`

export const createDenomIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, subdenom: string, numAttempts = 20) => {
  if (notExistsDenom(await queryDenomsFromCreator(client, account), account, subdenom)) {
    await client.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: account,
        subdenom
      }
    })
  } else {
    return
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryDenomsFromCreator(client, account),
    async (res) => existsDenom(res, account, subdenom),
    numAttempts,
  )
}

export const mintCoins = async (sdk: MantrachainSdk, client: any, account: string, subdenom: string, amount: number, minBalance?: number, numAttempts = 20) => {
  const denom = genCoinDenom(account, subdenom)
  const privBalance = await queryBalance(client, account, denom)

  if (!!minBalance && privBalance >= minBalance) {
    return privBalance
  }

  await client.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
    value: {
      sender: account,
      amount: {
        denom,
        amount: amount.toString()
      }
    }
  })

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryBalance(client, account, denom),
    async (balance) => balance === privBalance + amount,
    numAttempts,
  )
}