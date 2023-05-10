import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

const queryNftCollection = async (client: any, creator: string, id: string) => {
  const res = await client.MantrachainTokenV1.query.queryNftCollection(creator, id)
  return res?.data?.nft_collection || null
}

const existsNftCollection = async (collection: any) => collection !== null

export const createNftCollectionIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, collection: any, numAttempts = 20) => {
  if (!(await existsNftCollection(await queryNftCollection(client, account, collection.id)))) {
    await client.MantrachainTokenV1.tx.sendMsgCreateNftCollection({
      value: {
        creator: account,
        collection
      }
    })
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => { },
    async () => await existsNftCollection(await queryNftCollection(client, account, collection.id)),
    numAttempts,
  )
}