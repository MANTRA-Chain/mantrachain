import { MantrachainSdk } from '../helpers/sdk'
import { createDenomIfNotExists, getCoinDenom } from '../helpers/coinfactory'
import { createPairIfNotExists, getPairId } from '../helpers/liquidity'

describe('Guard module', () => {
  let sdk: MantrachainSdk

  beforeAll(async () => {
    sdk = new MantrachainSdk()
    await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)

    await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, 'testcoin1')
    await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, 'testcoin2')
    await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, 'testcoin3')
    await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, 'testcoin4')

    await createPairIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, getCoinDenom(sdk.adminAddress, 'testcoin3'), getCoinDenom(sdk.adminAddress, 'testcoin4'))
  })

  describe('Not Authenticated', () => {
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
          subdenom: 'testcoin'
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
          baseCoinDenom: getCoinDenom(sdk.adminAddress, 'testcoin1'),
          quoteCoinDenom: getCoinDenom(sdk.adminAddress, 'testcoin2')
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when create pool with account with no permission', async () => {
      const pairId = await getPairId(sdk.clientRecipient, getCoinDenom(sdk.adminAddress, 'testcoin3'), getCoinDenom(sdk.adminAddress, 'testcoin4'))

      const promise = sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePool({
        value: {
          creator: sdk.recipientAddress,
          pairId,
          depositCoins: [{
            denom: getCoinDenom(sdk.adminAddress, 'testcoin3'),
            amount: '1000000000000000000'
          }, {
            denom: getCoinDenom(sdk.adminAddress, 'testcoin4'),
            amount: '1000000000000000000'
          }]
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      )
    })

    test('Should throw when create ranged pool with account with no permission', async () => {
      const pairId = await getPairId(sdk.clientRecipient, getCoinDenom(sdk.adminAddress, 'testcoin3'), getCoinDenom(sdk.adminAddress, 'testcoin4'))

      const promise = sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreateRangedPool({
        value: {
          creator: sdk.recipientAddress,
          pairId,
          depositCoins: [{
            denom: getCoinDenom(sdk.adminAddress, 'testcoin3'),
            amount: '1000000000000000000'
          }, {
            denom: getCoinDenom(sdk.adminAddress, 'testcoin4'),
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
      const pairId = await getPairId(sdk.clientRecipient, getCoinDenom(sdk.adminAddress, 'testcoin3'), getCoinDenom(sdk.adminAddress, 'testcoin4'))

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
              denom: getCoinDenom(sdk.adminAddress, 'testcoin3'),
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
  })
})