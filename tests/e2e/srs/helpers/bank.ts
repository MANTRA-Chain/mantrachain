import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

export const queryBalance = (client: any, account: string, denom: string): any => client.CosmosBankV1Beta1.query.queryBalance(account, { denom })

export const sendCoins = async (sdk: MantrachainSdk, client: any, fromAddress: string, toAddress: string, denom: string, amount: string, numAttempts = 20) => {
  const privBalance = await queryBalance(client, toAddress, denom)

  await client.CosmosBankV1Beta1.tx.sendMsgSend({
    value: {
      fromAddress,
      toAddress,
      amount: [{
        denom,
        amount
      }]
    }
  })

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryBalance(client, toAddress, denom),
    async (res) => parseInt(res?.data?.balance?.amount) === parseInt(privBalance?.data?.balance?.amount) + parseInt(amount),
    numAttempts,
  )
}