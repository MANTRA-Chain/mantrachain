import { MantrachainSdk } from '../helpers/sdk'

describe('Guard module', () => {
  let sdk: MantrachainSdk

  beforeAll(async () => {
    const promises: Promise<any>[] = []
    sdk = new MantrachainSdk()
    await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)

    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.query.queryDenomsFromCreator(sdk.adminAddress)

    if (res.data?.denoms?.every((denom: string) => !denom.includes('testcoin1'))) {
      promises.push(sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
        value: {
          sender: sdk.adminAddress,
          subdenom: 'testcoin1'
        }
      }))
    }

    if (res.data?.denoms?.every((denom: string) => !denom.includes('testcoin2'))) {
      promises.push(sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
        value: {
          sender: sdk.adminAddress,
          subdenom: 'testcoin2'
        }
      }))
    }

    if (!!promises.length) {
      promises.push(sdk.blockWaiter.waitBlocks(1))
    }

    await Promise.all(promises)
  })

  describe('Not Authenticated', () => {
    test('Should throw when update account privileges with account with no permission', async () => {
      const promise = sdk.clientRecipient.MantrachainGuardV1.tx.sendMsgUpdateAccountPrivileges({
        value: {
          creator: sdk.recipientAddress,
          account: sdk.recipientAddress,
          privileges: Buffer.from(new Uint8Array(32)),
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
          privileges: Buffer.from(new Uint8Array(32)),
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
  })

  test('Should throw when create pair with account with no permission', async () => {
    const promise = sdk.clientRecipient.MantrachainLiquidityV1Beta1.tx.sendMsgCreatePair({
      value: {
        creator: sdk.recipientAddress,
        baseCoinDenom: `factory/${sdk.adminAddress}/testcon1`,
        quoteCoinDenom: `factory/${sdk.adminAddress}/testcon2`
      }
    })

    return expect(promise).rejects.toThrow(
      /unauthorized/
    )
  })
})