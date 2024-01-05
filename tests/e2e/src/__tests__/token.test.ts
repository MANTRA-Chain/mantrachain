import { AumegaSdk } from '../helpers/sdk'
import { createNftCollectionIfNotExists } from '../helpers/token'

describe('Token module', () => {
  let sdk: AumegaSdk

  beforeAll(async () => {
    sdk = new AumegaSdk()
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
    test('Should return error when approve nft from soul-bond nft collection', async () => {
      await expect(sdk.clientRecipient.AumegaTokenV1.tx.sendMsgApproveNft({
        value: {
          creator: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token0",
          nftId: "0",
          approved: true,
          strict: true
        }
      })).rejects.toThrow(
        /operation disabled/
      )
    })

    test('Should return error when approve nfts from soul-bond nft collection', async () => {
      await expect(sdk.clientRecipient.AumegaTokenV1.tx.sendMsgApproveNfts({
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
      })).rejects.toThrow(
        /operation disabled/
      )
    })

    test('Should return error when transfer nft from soul-bond nft collection', async () => {
      await expect(sdk.clientRecipient.AumegaTokenV1.tx.sendMsgTransferNft({
        value: {
          creator: sdk.recipientAddress,
          owner: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token0",
          nftId: "0",
          strict: true
        }
      })).rejects.toThrow(
        /operation disabled/
      )
    })

    test('Should return error when transfer nfts from soul-bond nft collection', async () => {
      await expect(sdk.clientRecipient.AumegaTokenV1.tx.sendMsgTransferNfts({
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
      })).rejects.toThrow(
        /operation disabled/
      )
    })
  })

  describe('Not Authenticated', () => {
    test('Should return error when create restricted nft collection from non-admin account', async () => {
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

      await expect(sdk.clientValidator.AumegaTokenV1.tx.sendMsgCreateNftCollection({
        value: {
          creator: sdk.validatorAddress,
          collection
        }
      })).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should return error when mint nft for restricted nft collection from non-admin account', async () => {
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

      await expect(sdk.clientValidator.AumegaTokenV1.tx.sendMsgMintNft({
        value: {
          creator: sdk.validatorAddress,
          receiver: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nft
        }
      })).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should return error when mint nfts for restricted nft collection from non-admin account', async () => {
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

      await expect(sdk.clientValidator.AumegaTokenV1.tx.sendMsgMintNfts({
        value: {
          creator: sdk.validatorAddress,
          receiver: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nfts: {
            nfts: [nft]
          }
        }
      })).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should return error when burn nft for restricted nft collection from non-admin account', async () => {
      await expect(sdk.clientValidator.AumegaTokenV1.tx.sendMsgBurnNft({
        value: {
          creator: sdk.validatorAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nftId: "0"
        }
      })).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should return error when burn nfts for restricted nft collection from non-admin account', async () => {
      await expect(sdk.clientValidator.AumegaTokenV1.tx.sendMsgBurnNfts({
        value: {
          creator: sdk.validatorAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nfts: {
            nftsIds: ["0"]
          }
        }
      })).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should return error when approve nft for restricted nft collection from non-admin account', async () => {
      await expect(sdk.clientValidator.AumegaTokenV1.tx.sendMsgApproveNft({
        value: {
          creator: sdk.validatorAddress,
          receiver: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nftId: "0",
          approved: true,
          strict: true
        }
      })).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should return error when approve nfts for restricted nft collection from non-admin account', async () => {
      await expect(sdk.clientValidator.AumegaTokenV1.tx.sendMsgApproveNfts({
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
      })).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should return error when transfer nft for restricted nft collection from non-admin account', async () => {
      await expect(sdk.clientValidator.AumegaTokenV1.tx.sendMsgTransferNft({
        value: {
          creator: sdk.validatorAddress,
          owner: sdk.adminAddress,
          receiver: sdk.validatorAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: "token1",
          nftId: "0",
          strict: true
        }
      })).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should return error when transfer nfts for restricted nft collection from non-admin account', async () => {
      await expect(sdk.clientValidator.AumegaTokenV1.tx.sendMsgTransferNfts({
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
      })).rejects.toThrow(
        /unauthorized/
      )
    })
  })
})