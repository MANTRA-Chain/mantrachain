import { AumegaSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

export const queryBalance = async (client: any, account: string, denom: string) => {
  const res = await client.CosmosBankV1Beta1.query.queryBalance(account, { denom })
  return !!res?.data?.balance?.amount ? parseInt(res.data.balance.amount) : 0
}

export const sendCoins = async (sdk: AumegaSdk, client: any, fromAddress: string, toAddress: string, denom: string, amount: number, minBalance?: number, numAttempts = 2) => {
  const privBalance = await queryBalance(client, toAddress, denom)

  if (!!minBalance && privBalance >= minBalance) {
    return privBalance
  }

  const res = await client.CosmosBankV1Beta1.tx.sendMsgSend({
    value: {
      fromAddress,
      toAddress,
      amount: [{
        denom,
        amount: amount.toString()
      }]
    }
  })

  if (res.code !== 0) {
    throw new Error(res.rawLog)
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryBalance(client, toAddress, denom),
    async (balance) => balance === privBalance + amount,
    numAttempts,
  )
}