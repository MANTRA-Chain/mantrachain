import { AumegaSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'
import { queryBalance } from './bank'

export const queryDenomsFromCreator = async (client: any, account: string) => {
  const res = await client.AumegaCoinfactoryV1Beta1.query.queryDenomsFromCreator(account)
  return res?.data?.denoms || []
}

export const queryAdmin = async (client: any, denom: string) => {
  const res = await client.AumegaCoinfactoryV1Beta1.query.queryDenomAuthorityMetadata(denom)
  return res?.data?.authority_metadata?.admin || null
}

const notExistsDenom = (denoms: string[], account: string, subdenom: string) => denoms.every((denom: string) => denom !== genCoinDenom(account, subdenom))

export const existsDenom = (denoms: string[], account: string, subdenom: string) => denoms.some((denom: string) => denom === genCoinDenom(account, subdenom))

export const genCoinDenom = (account: string, subdenom: string) => `factory/${account}/${subdenom}`

export const createDenomIfNotExists = async (sdk: AumegaSdk, client: any, account: string, subdenom: string, numAttempts = 2) => {
  if (notExistsDenom(await queryDenomsFromCreator(client, account), account, subdenom)) {
    const res = await client.AumegaCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: account,
        subdenom
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
    async () => await queryDenomsFromCreator(client, account),
    async (res) => existsDenom(res, account, subdenom),
    numAttempts,
  )
}

export const mintCoins = async (sdk: AumegaSdk, client: any, account: string, subdenom: string, amount: number, minBalance?: number, numAttempts = 2) => {
  const denom = genCoinDenom(account, subdenom)
  const privBalance = await queryBalance(client, account, denom)

  if (!!minBalance && privBalance >= minBalance) {
    return privBalance
  }

  const res = await client.AumegaCoinfactoryV1Beta1.tx.sendMsgMint({
    value: {
      sender: account,
      amount: {
        denom,
        amount: amount.toString()
      }
    }
  })

  if (res.code !== 0) {
    throw new Error(res.rawLog)
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryBalance(client, account, denom),
    async (balance) => balance === privBalance + amount,
    numAttempts,
  )
}