import { MantrachainSdk } from '../helpers/sdk'
import { createNftCollectionIfNotExists } from '../helpers/token'

describe('Token module', () => {
  let sdk: MantrachainSdk

  beforeAll(async () => {
    sdk = new MantrachainSdk()
    await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)

    await createNftCollectionIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, {
      id: "token0",
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

    await createNftCollectionIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, {
      id: "token1",
      name: 'test collection',
      images: [],
      url: '',
      description: '',
      links: [],
      options: [],
      category: 'utility',
      opened: false,
      symbol: 'TEST',
      soulBondedNfts: false,
      restrictedNfts: true,
      data: null
    })
  })

  describe('Soul-bond nfts collection', () => {
    test('Should throw when approve nft from soul-bond nft collection', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgApproveNft({
        value: {
          creator: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token0",
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
          collectionId: "token0",
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
          collectionId: "token0",
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
          collectionId: "token0",
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

  describe('Not Authenticated', () => {
    test('Should throw when create restricted nft collection from non-admin account', async () => {
      const collection = {
        id: "token2",
        name: 'test collection',
        images: [],
        url: '',
        description: '',
        links: [],
        options: [],
        category: 'utility',
        opened: false,
        symbol: 'TEST',
        soulBondedNfts: false,
        restrictedNfts: true,
        data: null
      }

      const promise = sdk.clientValidator.MantrachainTokenV1.tx.sendMsgCreateNftCollection({
        value: {
          creator: sdk.validatorAddress,
          collection
        }
      })

      return expect(promise).rejects.toThrow(
        /guard token: fail/
      )
    })

    test('Should throw when mint nft for restricted nft collection from non-admin account', async () => {
      const nft = {
        id: "0",
        title: 'test nft',
        images: [],
        url: '',
        description: '',
        links: [],
        attributes: [],
        data: null
      }

      const promise = sdk.clientValidator.MantrachainTokenV1.tx.sendMsgMintNft({
        value: {
          creator: sdk.validatorAddress,
          receiver: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nft
        }
      })

      return expect(promise).rejects.toThrow(
        /guard token: fail/
      )
    })

    test('Should throw when mint nfts for restricted nft collection from non-admin account', async () => {
      const nft = {
        id: "0",
        title: 'test nft',
        images: [],
        url: '',
        description: '',
        links: [],
        attributes: [],
        data: null
      }

      const promise = sdk.clientValidator.MantrachainTokenV1.tx.sendMsgMintNfts({
        value: {
          creator: sdk.validatorAddress,
          receiver: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nfts: {
            nfts: [nft]
          }
        }
      })

      return expect(promise).rejects.toThrow(
        /guard token: fail/
      )
    })

    test('Should throw when burn nft for restricted nft collection from non-admin account', async () => {
      const promise = sdk.clientValidator.MantrachainTokenV1.tx.sendMsgBurnNft({
        value: {
          creator: sdk.validatorAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nftId: "0"
        }
      })

      return expect(promise).rejects.toThrow(
        /guard token: fail/
      )
    })

    test('Should throw when burn nfts for restricted nft collection from non-admin account', async () => {
      const promise = sdk.clientValidator.MantrachainTokenV1.tx.sendMsgBurnNfts({
        value: {
          creator: sdk.validatorAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nfts: {
            nftsIds: ["0"]
          }
        }
      })

      return expect(promise).rejects.toThrow(
        /guard token: fail/
      )
    })

    test('Should throw when approve nft for restricted nft collection from non-admin account', async () => {
      const promise = sdk.clientValidator.MantrachainTokenV1.tx.sendMsgApproveNft({
        value: {
          creator: sdk.validatorAddress,
          receiver: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nftId: "0",
          approved: true,
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /guard token: fail/
      )
    })

    test('Should throw when approve nfts for restricted nft collection from non-admin account', async () => {
      const promise = sdk.clientValidator.MantrachainTokenV1.tx.sendMsgApproveNfts({
        value: {
          creator: sdk.validatorAddress,
          receiver: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nfts: {
            nftsIds: ["0"]
          },
          approved: true,
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /guard token: fail/
      )
    })

    test('Should throw when transfer nft for restricted nft collection from non-admin account', async () => {
      const promise = sdk.clientValidator.MantrachainTokenV1.tx.sendMsgTransferNft({
        value: {
          creator: sdk.validatorAddress,
          owner: sdk.adminAddress,
          receiver: sdk.validatorAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nftId: "0",
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /guard token: fail/
      )
    })

    test('Should throw when transfer nfts for restricted nft collection from non-admin account', async () => {
      const promise = sdk.clientValidator.MantrachainTokenV1.tx.sendMsgTransferNfts({
        value: {
          creator: sdk.validatorAddress,
          owner: sdk.adminAddress,
          receiver: sdk.validatorAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nfts: {
            nftsIds: ["0"]
          },
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /guard token: fail/
      )
    })
  })
})