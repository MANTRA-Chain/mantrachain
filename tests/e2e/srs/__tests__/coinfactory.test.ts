import { MantrachainSdk } from '../helpers/sdk'
import { existsDenom, queryDenomsFromCreator, genCoinDenom, queryAdmin } from '../helpers/coinfactory'
import { queryBalance, sendCoins } from '../helpers/bank'

describe('Coinfactory module', () => {
  let sdk: MantrachainSdk
  const subdenom = `cf${new Date().getTime()}`

  beforeAll(async () => {
    sdk = new MantrachainSdk()
    await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)
  })

  test('Should create denom', async () => {
    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.adminAddress,
        subdenom
      }
    })

    expect(existsDenom(await queryDenomsFromCreator(sdk.clientAdmin, sdk.adminAddress), sdk.adminAddress, subdenom)).toBeTruthy()
  })

  test('Should mint coins', async () => {
    const amount = 1000
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    const currBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should burn coins', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgBurn({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    const currBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    expect(privBalance).toEqual(currBalance + amount)
  })

  test('Should force transfer coins', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    expect(await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)).toEqual(amount)

    await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress, denom, amount)

    const privBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
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

    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should change admin', async () => {
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgChangeAdmin({
      value: {
        sender: sdk.adminAddress,
        denom,
        newAdmin: sdk.validatorAddress,
      }
    })

    // TODO: fix queryAdmin issue when adding coin queries for denoms with slashes(/)
    // const admin = await queryAdmin(sdk.clientValidator, denom)

    // expect(admin).toEqual(sdk.validatorAddress)
  })

  test('Should mint coins with the new admin', async () => {
    const amount = 1000
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    await sdk.clientValidator.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.validatorAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should burn coins with the new admin', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    await sdk.clientValidator.MantrachainCoinfactoryV1Beta1.tx.sendMsgBurn({
      value: {
        sender: sdk.validatorAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    expect(privBalance).toEqual(currBalance + amount)
  })

  test('Should throw when force transfer coins with the new admin w/o permission', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const promise = sdk.clientValidator.MantrachainCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
      value: {
        sender: sdk.validatorAddress,
        amount: {
          denom,
          amount: amount.toString()
        },
        transferFromAddress: sdk.recipientAddress,
        transferToAddress: sdk.validatorAddress
      }
    })

    return expect(promise).rejects.toThrow(
      /unauthorized/
    )
  })

  test('Should force transfer coins when the admin is changed admin but with account with permissions', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    expect(await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)).toEqual(amount)

    const privBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
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

    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should throw when mint coins with not the admin', async () => {
    const amount = 1000
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const promise = sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    return expect(promise).rejects.toThrow(
      /unauthorized/
    )
  })

  test('Should throw when burn coins with not the admin', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const promise = sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgBurn({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      }
    })

    return expect(promise).rejects.toThrow(
      /unauthorized/
    )
  })

  test('Should throw when create denom with account with no permission', async () => {
    const promise = sdk.clientValidator.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.validatorAddress,
        subdenom: 'cf0'
      }
    })

    return expect(promise).rejects.toThrow(
      /unauthorized/
    )
  })

  test('Should throw when create existing denom', async () => {
    const promise = sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.adminAddress,
        subdenom
      }
    })

    return expect(promise).rejects.toThrow(
      /attempting to create a denom that already exists/
    )
  })
})