import { MantrachainSdk } from '../helpers/sdk'
import { createNftCollectionIfNotExists } from '../helpers/token'

describe('Token module', () => {
  let sdk: MantrachainSdk

  beforeAll(async () => {
    sdk = new MantrachainSdk()
    await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)

    await createNftCollectionIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, {
      id: "token-0",
      name: 'test collection',
      images: [],
      url: '',
      description: '',
      links: [],
      options: [],
      category: 'utility',
      opened: false,
      symbol: 'TEST',
      soulBondedNfts: true,
      restrictedNfts: true,
      data: null
    })
  })

  describe('Not Authenticated', () => {
    test('Should throw when approve nft from soul-bond nft collection', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgApproveNft({
        value: {
          creator: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token-0",
          nftId: "0",
          approved: true,
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /operation disabled/
      )
    })

    test('Should throw when approve nfts from soul-bond nft collection', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgApproveNfts({
        value: {
          creator: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token-0",
          nfts: {
            nftsIds: ["0"]
          },
          approved: true,
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /operation disabled/
      )
    })

    test('Should throw when transfer nft from soul-bond nft collection', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgTransferNft({
        value: {
          creator: sdk.recipientAddress,
          owner: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token-0",
          nftId: "0",
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /operation disabled/
      )
    })

    test('Should throw when transfer nfts from soul-bond nft collection', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgTransferNfts({
        value: {
          creator: sdk.recipientAddress,
          owner: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token-0",
          nfts: {
            nftsIds: ["0"]
          },
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /operation disabled/
      )
    })
  })
})