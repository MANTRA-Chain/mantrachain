import { MantrachainSdk } from '../helpers/sdk'
import { createDenomIfNotExists, genCoinDenom, mintCoins } from '../helpers/coinfactory'
import { createPairIfNotExists, getPairId } from '../helpers/liquidity'
import { createNftCollectionIfNotExists } from '../helpers/token'
import { setGuardTransferCoins, updateCoinRequiredPrivileges, updateAccountPrivileges } from '../helpers/guard'
import { queryBalance, sendCoins } from '../helpers/bank'

let sdk: MantrachainSdk

describe('Guard module', () => {
  beforeAll(async () => {
    sdk = new MantrachainSdk()
    await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)
  })

  describe('Not Authenticated', () => {
    beforeAll(async () => {
      await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, 'guard0')
      await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, 'guard1')
      await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, 'guard2')
      await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, 'guard3')

      await createPairIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, genCoinDenom(sdk.adminAddress, 'guard0'), genCoinDenom(sdk.adminAddress, 'guard1'))

      await createNftCollectionIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, {
        id: 'guard0',
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

    test('Should throw when update account privileges with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateAccountPrivileges({
        value: {
          creator: sdk.recipientAddress,
          account: sdk.recipientAddress,
          privileges: new Uint8Array(32),
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when update account privileges batch with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateAccountPrivilegesBatch({
        value: {
          creator: sdk.recipientAddress,
          accountsPrivileges: {
            accounts: [sdk.recipientAddress],
            privileges: [new Uint8Array(32)]
          }
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when update account privileges grouped batch with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateAccountPrivilegesGroupedBatch({
        value: {
          creator: sdk.recipientAddress,
          accountsPrivilegesGrouped: {
            accounts: [{
              accounts: [sdk.recipientAddress]
            }],
            privileges: [new Uint8Array(32)]
          }
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when update required privileges with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateRequiredPrivileges({
        value: {
          creator: sdk.recipientAddress,
          index: new Uint8Array(1),
          privileges: new Uint8Array(32),
          kind: 'coin'
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when update required privileges batch with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateRequiredPrivilegesBatch({
        value: {
          creator: sdk.recipientAddress,
          requiredPrivileges: {
            indexes: [new Uint8Array(1)],
            privileges: [new Uint8Array(32)]
          },
          kind: 'coin'
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when update required privileges grouped batch with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateRequiredPrivilegesGroupedBatch({
        value: {
          creator: sdk.recipientAddress,
          requiredPrivilegesGrouped: {
            indexes: [{
              indexes: [new Uint8Array(1)]
            }],
            privileges: [new Uint8Array(32)]
          },
          kind: 'coin'
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when update guard transfer coins with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateGuardTransferCoins({
        value: {
          creator: sdk.recipientAddress,
          enabled: true
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when update locked with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateLocked({
        value: {
          creator: sdk.recipientAddress,
          index: new Uint8Array(1),
          locked: true,
          kind: 'coin'
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when update authz grants with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateAuthzGenericGrantRevokeBatch({
        value: {
          creator: sdk.recipientAddress,
          grantee: sdk.adminAddress,
          authzGrantRevokeMsgsTypes: {
            msgs: [
              { typeUrl: '/mantrachain.coinfactory.v1beta1.MsgCreateDenom', grant: true }
            ],
          }
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when create denom with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
        value: {
          sender: sdk.recipientAddress,
          subdenom: 'guard4'
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when force transfer with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
        value: {
          sender: sdk.recipientAddress,
          amount: {
            denom: 'uaum',
            amount: '1000000000000000000'
          },
          transferFromAddress: sdk.recipientAddress,
          transferToAddress: sdk.adminAddress,
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when create pair with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
        value: {
          creator: sdk.recipientAddress,
          baseCoinDenom: genCoinDenom(sdk.adminAddress, 'guard2'),
          quoteCoinDenom: genCoinDenom(sdk.adminAddress, 'guard3')
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when create pool with account with no permission', async () => {
      const pairId = await getPairId(sdk.clientRecipient, genCoinDenom(sdk.adminAddress, 'guard0'), genCoinDenom(sdk.adminAddress, 'guard1'))

      const promise = sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePool({
        value: {
          creator: sdk.recipientAddress,
          pairId,
          depositCoins: [{
            denom: genCoinDenom(sdk.adminAddress, 'guard0'),
            amount: '1000000000000000000'
          }, {
            denom: genCoinDenom(sdk.adminAddress, 'guard1'),
            amount: '1000000000000000000'
          }]
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when create ranged pool with account with no permission', async () => {
      const pairId = await getPairId(sdk.clientRecipient, genCoinDenom(sdk.adminAddress, 'guard0'), genCoinDenom(sdk.adminAddress, 'guard1'))

      const promise = sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreateRangedPool({
        value: {
          creator: sdk.recipientAddress,
          pairId,
          depositCoins: [{
            denom: genCoinDenom(sdk.adminAddress, 'guard0'),
            amount: '1000000000000000000'
          }, {
            denom: genCoinDenom(sdk.adminAddress, 'guard1'),
            amount: '1000000000000000000'
          }],
          minPrice: '1000000000',
          maxPrice: '10000000000',
          initialPrice: '5000000000'
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when create private plan with account with no permission', async () => {
      const pairId = await getPairId(sdk.clientRecipient, genCoinDenom(sdk.adminAddress, 'guard0'), genCoinDenom(sdk.adminAddress, 'guard1'))

      const startTime = new Date()
      const endTime = new Date()
      endTime.setMinutes(endTime.getMinutes() + 30)

      const promise = sdk.clientRecipient.MantrachainLpfarmV1Beta1.tx.sendMsgCreatePrivatePlan({
        value: {
          creator: sdk.recipientAddress,
          description: 'test plan',
          rewardAllocations: [{
            pairId,
            rewardsPerDay: [{
              denom: genCoinDenom(sdk.adminAddress, 'guard0'),
              amount: '1000'
            }]
          }],
          startTime,
          endTime,
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when create restricted nft collection with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgCreateNftCollection({
        value: {
          creator: sdk.recipientAddress,
          collection: {
            id: 'guard0',
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
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when mint nft for restricted nft collection with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgMintNft({
        value: {
          creator: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: 'guard0',
          nft: {
            id: '0',
            title: 'test nft',
            images: [],
            url: '',
            description: '',
            links: [],
            attributes: [],
            data: null
          }
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when mint nfts for restricted nft collection with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgMintNfts({
        value: {
          creator: sdk.recipientAddress,
          receiver: sdk.recipientAddress,
          strict: false,
          collectionCreator: sdk.adminAddress,
          collectionId: 'guard0',
          nfts: {
            nfts: [{
              id: '0',
              title: 'test nft',
              images: [],
              url: '',
              description: '',
              links: [],
              attributes: [],
              data: null
            }]
          }
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when burn nft from restricted nft collection with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgBurnNft({
        value: {
          creator: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: 'guard0',
          nftId: '0',
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when burn nfts from restricted nft collection with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgBurnNfts({
        value: {
          creator: sdk.recipientAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: 'guard0',
          nfts: {
            nftsIds: ['0']
          },
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when approve nft from restricted nft collection with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgApproveNft({
        value: {
          creator: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: 'guard0',
          nftId: '0',
          approved: true,
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when approve nfts from restricted nft collection with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgApproveNfts({
        value: {
          creator: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: 'guard0',
          nfts: {
            nftsIds: ['0']
          },
          approved: true,
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when transfer nft from restricted nft collection with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgTransferNft({
        value: {
          creator: sdk.recipientAddress,
          owner: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: 'guard0',
          nftId: '0',
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when transfer nfts from restricted nft collection with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainTokenV1.tx.sendMsgTransferNfts({
        value: {
          creator: sdk.recipientAddress,
          owner: sdk.recipientAddress,
          receiver: sdk.adminAddress,
          collectionCreator: sdk.adminAddress,
          collectionId: 'guard0',
          nfts: {
            nftsIds: ['0']
          },
          strict: true
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })
  })

  describe('Transfer coins', () => {
    beforeAll(async () => {
      await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, 'guard5')
      await mintCoins(sdk, sdk.clientAdmin, sdk.adminAddress, 'guard5', 1000000, 1000000)
    })

    test('Should transfer from/to account without guard soul-bond nft when guard transfer coins is false', async () => {
      await setGuardTransferCoins(sdk, sdk.clientAdmin, sdk.adminAddress, false)

      const denom = genCoinDenom(sdk.adminAddress, 'guard5')
      const amount = 1000
      const privBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

      await updateCoinRequiredPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, denom)
      await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress)

      await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress, denom, amount)
      const currBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

      expect(privBalance + amount).toEqual(currBalance);

      await sendCoins(sdk, sdk.clientRecipient, sdk.recipientAddress, sdk.adminAddress, denom, currBalance)
      const latestBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

      expect(latestBalance).toEqual(0);
    })

    test('Should transfer from admin account to account without guard soul-bond nft when guard transfer coins is true', async () => {
      await setGuardTransferCoins(sdk, sdk.clientAdmin, sdk.adminAddress, true)

      const denom = genCoinDenom(sdk.adminAddress, 'guard5')
      const amount = 1000
      const privBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

      await updateCoinRequiredPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, denom)
      await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress)

      await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress, denom, amount)
      const currBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

      expect(privBalance + amount).toEqual(currBalance);
    })

    test('Should throw when transfer from account without guard soul-bond nft when guard transfer coins is true', async () => {
      await setGuardTransferCoins(sdk, sdk.clientAdmin, sdk.adminAddress, true)

      const denom = genCoinDenom(sdk.adminAddress, 'guard5')
      const amount = 1000

      await updateCoinRequiredPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, denom)
      await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress)

      await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress, denom, amount, amount)
      const currBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

      const promise = sendCoins(sdk, sdk.clientRecipient, sdk.recipientAddress, sdk.adminAddress, denom, currBalance)

      return expect(promise).rejects.toThrow(
        /missing soul bond nft/
      )
    })

    test('Should throw when transfer coin without required privileges', async () => {
      await setGuardTransferCoins(sdk, sdk.clientAdmin, sdk.adminAddress, true)

      const denom = genCoinDenom(sdk.adminAddress, 'guard5')
      const amount = 1000

      await updateCoinRequiredPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, denom)
      await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress)

      await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress, denom, amount, amount)
      const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

      const promise = sendCoins(sdk, sdk.clientValidator, sdk.validatorAddress, sdk.adminAddress, denom, currBalance)

      return expect(promise).rejects.toThrow(
        /coin required privileges not found/
      )
    })

    test('Should throw when transfer coin from account without account privileges', async () => {
      await setGuardTransferCoins(sdk, sdk.clientAdmin, sdk.adminAddress, true)

      const denom = genCoinDenom(sdk.adminAddress, 'guard5')
      const amount = 1000

      await updateCoinRequiredPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, denom, [64])
      await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress)
      await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress)

      await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress, denom, amount, amount)
      const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

      const promise = sendCoins(sdk, sdk.clientValidator, sdk.validatorAddress, sdk.recipientAddress, denom, currBalance)

      return expect(promise).rejects.toThrow(
        /insufficient privileges/
      )
    })

    test('Should throw when transfer coin to account without soul-bond nft', async () => {
      await setGuardTransferCoins(sdk, sdk.clientAdmin, sdk.adminAddress, true)

      const denom = genCoinDenom(sdk.adminAddress, 'guard5')
      const amount = 1000

      await updateCoinRequiredPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, denom, [64])
      await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress, [64])
      await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress)

      await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress, denom, amount, amount)
      const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

      const promise = sendCoins(sdk, sdk.clientValidator, sdk.validatorAddress, sdk.recipientAddress, denom, currBalance)

      return expect(promise).rejects.toThrow(
        /missing soul bond nft/
      )
    })

    // TODO: test('Should transfer from and to account with guard soul-bond nft and privileges when guard transfer coins is true', async () => {
    // })
  })
})