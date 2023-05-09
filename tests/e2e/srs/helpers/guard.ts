import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

const queryGuardTransferCoins = (client: any): any => client.MantrachainGuardV1.query.queryGuardTransferCoins()

const notSetGuardTransferCoins = (res: any, guardTransferCoins: boolean) => res.data?.guard_transfer_coins !== guardTransferCoins

export const setGuardTransferCoins = async (sdk: MantrachainSdk, client: any, account: string, enabled: boolean, numAttempts = 20) => {
  if (notSetGuardTransferCoins(await queryGuardTransferCoins(client), enabled)) {
    await client.MantrachainGuardV1.tx.sendMsgUpdateGuardTransferCoins({
      value: {
        creator: account,
        enabled
      }
    })
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryGuardTransferCoins(client),
    async (res) => !notSetGuardTransferCoins(res, enabled),
    numAttempts,
  )
}