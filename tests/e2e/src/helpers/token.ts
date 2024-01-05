import { AumegaSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'

const queryNftCollection = async (client: any, creator: string, id: string) => {
  try {
    const res = await client.AumegaTokenV1.query.queryNftCollection(creator, id)
    return res?.data?.nft_collection || null
  } catch (e) {
    return null
  }
}

const queryNft = async (client: any, collectionCreator: string, collectionId: string, id: string) => {
  try {
    const res = await client.AumegaTokenV1.query.queryNft(collectionCreator, collectionId, id)
    return res?.data?.nft || null
  } catch (e) {
    return null
  }
}

const existsNftCollection = (collection: any) => collection !== null

const existsNft = (nft: any) => nft !== null

const queryGuardNftTokenCollectionParams = async (client: any): Promise<any> => {
  const res = await client.AumegaGuardV1.query.queryParams()
  return {
    collectionCreator: res.data.params.account_privileges_token_collection_creator,
    collectionId: res.data.params.account_privileges_token_collection_id
  }
}

export const createNftCollectionIfNotExists = async (sdk: AumegaSdk, client: any, account: string, collection: any, numAttempts = 2) => {
  if (!existsNftCollection(await queryNftCollection(client, account, collection.id))) {
    const res = await client.AumegaTokenV1.tx.sendMsgCreateNftCollection({
      value: {
        creator: account,
        collection
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
    async () => { },
    async () => await existsNftCollection(await queryNftCollection(client, account, collection.id)),
    numAttempts,
  )
}

export const mintGuardSoulBondNft = async (sdk: AumegaSdk, client: any, account: string, receiver: string, numAttempts = 2) => {
  const guardCollectionParams = await queryGuardNftTokenCollectionParams(client)
  if (!existsNft(await queryNft(client, guardCollectionParams.collectionCreator, guardCollectionParams.collectionId, receiver))) {
    const res = await client.AumegaTokenV1.tx.sendMsgMintNft({
      value: {
        creator: account,
        receiver,
        collectionCreator: guardCollectionParams.collectionCreator,
        collectionId: guardCollectionParams.collectionId,
        nft: {
          id: receiver,
          title: 'AccountPrivileges',
          images: [],
          url: '',
          description: 'AccountPrivileges',
          links: [],
          attributes: [],
          data: null
        },
        strict: true,
        did: true,
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
    async () => await queryNft(client, guardCollectionParams.collectionCreator, guardCollectionParams.collectionId, receiver),
    async (nft) => existsNft(nft),
    numAttempts,
  )
}

export const burnGuardSoulBondNft = async (sdk: AumegaSdk, client: any, account: string, id: string, numAttempts = 2) => {
  const guardCollectionParams = await queryGuardNftTokenCollectionParams(client)
  if (existsNft(await queryNft(client, guardCollectionParams.collectionCreator, guardCollectionParams.collectionId, id))) {
    const res = await client.AumegaTokenV1.tx.sendMsgBurnNft({
      value: {
        creator: account,
        collectionCreator: guardCollectionParams.collectionCreator,
        collectionId: guardCollectionParams.collectionId,
        nftId: id
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
    async () => await queryNft(client, guardCollectionParams.collectionCreator, guardCollectionParams.collectionId, id),
    async (nft) => !existsNft(nft),
    numAttempts,
  )
}
