import { AumegaSdk } from '../helpers/sdk'
import { existsDenom, queryDenomsFromCreator, genCoinDenom } from '../helpers/coinfactory'
import { queryBalance, sendCoins } from '../helpers/bank'
import { updateCoinRequiredPrivileges, updateAccountPrivileges } from '../helpers/guard'

describe('Coinfactory module', () => {
  let sdk: AumegaSdk
  const subdenom = `cf${new Date().getTime()}`

  beforeAll(async () => {
    sdk = new AumegaSdk()
    await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)

    // Adds permission and allows for the recipient to be able to create denoms
    await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress, [100])
  })
  
  afterAll(async () => {
    // Clear the permissions for the recipient
    await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.recipientAddress)
  })

  test('Should create denom', async () => {
    const res = await sdk.clientAdmin.AumegaCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.adminAddress,
        subdenom
      }
    })

    expect(res.code).toBe(0)
    expect(existsDenom(await queryDenomsFromCreator(sdk.clientAdmin, sdk.adminAddress), sdk.adminAddress, subdenom)).toBeTruthy()
  })

  test('Should mint coins', async () => {
    const amount = 1000
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    const res = await sdk.clientAdmin.AumegaCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    const currBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should burn coins', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    const res = await sdk.clientAdmin.AumegaCoinfactoryV1Beta1.tx.sendMsgBurn({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    const currBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance).toEqual(currBalance + amount)
  })

  test('Should force transfer coins', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    expect(await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)).toEqual(amount)

    await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress, denom, amount)

    const privBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

    const res = await sdk.clientAdmin.AumegaCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        },
        transferFromAddress: sdk.validatorAddress,
        transferToAddress: sdk.recipientAddress
      }
    })

    const currBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should change admin', async () => {
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const res = await sdk.clientAdmin.AumegaCoinfactoryV1Beta1.tx.sendMsgChangeAdmin({
      value: {
        sender: sdk.adminAddress,
        denom,
        newAdmin: sdk.validatorAddress,
      }
    })

    // TODO: fix queryAdmin issue when adding coin queries for denoms with slashes(/)
    // const admin = await queryAdmin(sdk.clientValidator, denom)

    expect(res.code).toBe(0)
    // expect(admin).toEqual(sdk.validatorAddress)
  })

  test('Should mint coins with the new admin', async () => {
    const amount = 1000
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    const res = await sdk.clientValidator.AumegaCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.validatorAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should burn coins with the new admin', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    const res = await sdk.clientValidator.AumegaCoinfactoryV1Beta1.tx.sendMsgBurn({
      value: {
        sender: sdk.validatorAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance).toEqual(currBalance + amount)
  })

  test('Should return error when force transfer coins with the new admin w/o permission', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    await expect(sdk.clientValidator.AumegaCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
      value: {
        sender: sdk.validatorAddress,
        amount: {
          denom,
          amount: amount.toString()
        },
        transferFromAddress: sdk.recipientAddress,
        transferToAddress: sdk.validatorAddress
      }
    })).rejects.toThrow(
      /unauthorized/
    )
  })

  test('Should force transfer coins when the admin is changed admin but with account with permissions', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    expect(await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)).toEqual(amount)

    const privBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    const res = await sdk.clientAdmin.AumegaCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        },
        transferFromAddress: sdk.recipientAddress,
        transferToAddress: sdk.validatorAddress
      }
    })

    const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should return error when mint coins with not the admin', async () => {
    const amount = 1000
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    await expect(sdk.clientAdmin.AumegaCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })).rejects.toThrow(
      /unauthorized/
    )
  })

  test('Should return error when burn coins with not the admin', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    await expect(sdk.clientAdmin.AumegaCoinfactoryV1Beta1.tx.sendMsgBurn({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })).rejects.toThrow(
      /unauthorized/
    )
  })

  test('Should return error when create denom with account with no permission', async () => {
    await expect(sdk.clientValidator.AumegaCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.validatorAddress,
        subdenom: 'cf0'
      }
    })).rejects.toThrow(
      /unauthorized/
    )
  })

  test('Should create denom with account with permission', async () => {
    const res = await sdk.clientRecipient.AumegaCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.recipientAddress,
        subdenom
      }
    })

    expect(res.code).toBe(0)
    expect(existsDenom(await queryDenomsFromCreator(sdk.clientRecipient, sdk.recipientAddress), sdk.recipientAddress, subdenom)).toBeTruthy()
  })

  test('Should return error when create existing denom', async () => {
    await expect(sdk.clientAdmin.AumegaCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.adminAddress,
        subdenom
      }
    })).rejects.toThrow(
      /attempting to create a denom that already exists/
    )
  })
})