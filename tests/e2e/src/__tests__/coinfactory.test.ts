import { MantrachainSdk, getGasFee } from '../helpers/sdk'
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
    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.adminAddress,
        subdenom
      },
      fee: getGasFee()
    })

    expect(res.code).toBe(0)
    expect(existsDenom(await queryDenomsFromCreator(sdk.clientAdmin, sdk.adminAddress), sdk.adminAddress, subdenom)).toBeTruthy()
  })

  test('Should mint coins', async () => {
    const amount = 1000
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      },
      fee: getGasFee()
    })

    const currBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should burn coins', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientAdmin, sdk.adminAddress, denom)

    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgBurn({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      },
      fee: getGasFee()
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

    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        },
        transferFromAddress: sdk.validatorAddress,
        transferToAddress: sdk.recipientAddress
      },
      fee: getGasFee()
    })

    const currBalance = await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should change admin', async () => {
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgChangeAdmin({
      value: {
        sender: sdk.adminAddress,
        denom,
        newAdmin: sdk.validatorAddress,
      },
      fee: getGasFee()
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

    const res = await sdk.clientValidator.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.validatorAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      },
      fee: getGasFee()
    })

    const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should burn coins with the new admin', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const privBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    const res = await sdk.clientValidator.MantrachainCoinfactoryV1Beta1.tx.sendMsgBurn({
      value: {
        sender: sdk.validatorAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      },
      fee: getGasFee()
    })

    const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance).toEqual(currBalance + amount)
  })

  test('Should return error when force transfer coins with the new admin w/o permission', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const res = await sdk.clientValidator.MantrachainCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
      value: {
        sender: sdk.validatorAddress,
        amount: {
          denom,
          amount: amount.toString()
        },
        transferFromAddress: sdk.recipientAddress,
        transferToAddress: sdk.validatorAddress
      },
      fee: getGasFee()
    })

    expect(res.code).not.toBe(0)
    expect(res.rawLog).toMatch(
      /unauthorized/
    )
  })

  test('Should force transfer coins when the admin is changed admin but with account with permissions', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    expect(await queryBalance(sdk.clientRecipient, sdk.recipientAddress, denom)).toEqual(amount)

    const privBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgForceTransfer({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        },
        transferFromAddress: sdk.recipientAddress,
        transferToAddress: sdk.validatorAddress
      },
      fee: getGasFee()
    })

    const currBalance = await queryBalance(sdk.clientValidator, sdk.validatorAddress, denom)

    expect(res.code).toBe(0)
    expect(privBalance + amount).toEqual(currBalance)
  })

  test('Should return error when mint coins with not the admin', async () => {
    const amount = 1000
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgMint({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      },
      fee: getGasFee()
    })

    expect(res.code).not.toBe(0)
    expect(res.rawLog).toMatch(
      /unauthorized/
    )
  })

  test('Should return error when burn coins with not the admin', async () => {
    const amount = 500
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgBurn({
      value: {
        sender: sdk.adminAddress,
        amount: {
          denom,
          amount: amount.toString()
        }
      },
      fee: getGasFee()
    })

    expect(res.code).not.toBe(0)
    expect(res.rawLog).toMatch(
      /unauthorized/
    )
  })

  test('Should return error when create denom with account with no permission', async () => {
    const res = await sdk.clientValidator.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.validatorAddress,
        subdenom: 'cf0'
      },
      fee: getGasFee()
    })

    expect(res.code).not.toBe(0)
    expect(res.rawLog).toMatch(
      /unauthorized/
    )
  })

  test('Should return error when create existing denom', async () => {
    const res = await sdk.clientAdmin.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
      value: {
        sender: sdk.adminAddress,
        subdenom
      },
      fee: getGasFee()
    })

    expect(res.code).not.toBe(0)
    expect(res.rawLog).toMatch(
      /attempting to create a denom that already exists/
    )
  })
})