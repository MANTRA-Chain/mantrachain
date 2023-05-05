import { MantrachainSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

const queryNftCollection = (client: any, creator: string, id: string): any => client.MantrachainTokenV1.query.queryNftCollection(creator, id)

const existsNftCollection = async (client: any, creator: string, id: string) => {
  try {
    const res = await queryNftCollection(client, creator, id)
    return res?.data?.['nft_collection']?.['id'] === id && res?.data?.['nft_collection']?.['creator'] === creator
  } catch (e) {
    return false
  }
}

export const createNftCollectionIfNotExists = async (sdk: MantrachainSdk, client: any, account: string, collection: any, numAttempts = 20) => {
  if (!(await existsNftCollection(client, account, collection.id))) {
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
    async () => await existsNftCollection(client, account, collection.id),
    numAttempts,
  )
}