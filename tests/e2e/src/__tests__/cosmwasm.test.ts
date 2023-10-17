import { MantrachainSdk } from '../helpers/sdk'
import { createDenomIfNotExists, genCoinDenom, mintCoins } from '../helpers/coinfactory'
import { setGuardTransferCoins, updateCoinRequiredPrivileges, updateAccountPrivileges } from '../helpers/guard'
import { burnGuardSoulBondNft, mintGuardSoulBondNft } from '../helpers/token'
import { queryBalance, sendCoins } from '../helpers/bank'

describe('Cosm Wasm module', () => {
  let sdk: MantrachainSdk
  const amount = 1000
  const subdenom = 'cosmwasm0'
  let contractAddress = ''

  beforeAll(async () => {
    sdk = new MantrachainSdk()
    await sdk.init(process.env.API_URL, process.env.RPC_URL, process.env.WS_URL)

    const denom = genCoinDenom(sdk.adminAddress, subdenom)
    await createDenomIfNotExists(sdk, sdk.clientAdmin, sdk.adminAddress, subdenom)
    await mintCoins(sdk, sdk.clientAdmin, sdk.adminAddress, subdenom, amount, amount)
    await setGuardTransferCoins(sdk, sdk.clientAdmin, sdk.adminAddress, true)
    await updateCoinRequiredPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, denom, [64])
    await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress, [64])
    await sendCoins(sdk, sdk.clientAdmin, sdk.adminAddress, sdk.validatorAddress, denom, amount)
    contractAddress = (await sdk.clientAdmin.CosmwasmWasmV1.query.queryContractsByCode("1"))["data"]["contracts"][0]
  })

  test('Should return error when transfer coin to contract without soul-bond nft', async () => {
    const denom = genCoinDenom(sdk.adminAddress, subdenom)
    await burnGuardSoulBondNft(sdk, sdk.clientAdmin, sdk.adminAddress, contractAddress)

    await expect(sdk.clientValidator.CosmosBankV1Beta1.tx.sendMsgSend({
      value: {
        fromAddress: sdk.validatorAddress,
        toAddress: contractAddress,
        amount: [{
          denom,
          amount: amount.toString()
        }]
      }
    })).rejects.toThrow(
      /missing soul bond nft/
    )
  })
  
  test('Should return error when transfer coin to contract without account privileges', async () => {
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    await mintGuardSoulBondNft(sdk, sdk.clientAdmin, sdk.adminAddress, contractAddress)
    await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, contractAddress)

    await expect(sdk.clientValidator.CosmosBankV1Beta1.tx.sendMsgSend({
      value: {
        fromAddress: sdk.validatorAddress,
        toAddress: contractAddress,
        amount: [{
          denom,
          amount: amount.toString()
        }]
      }
    })).rejects.toThrow(
      /insufficient privileges/
    )
  })

  test('Should transfer coin to a contract with guard soul-bond nft and privileges', async () => {
    const denom = genCoinDenom(sdk.adminAddress, subdenom)

    await updateAccountPrivileges(sdk, sdk.clientAdmin, sdk.adminAddress, contractAddress, [64])
    const currBalance = await queryBalance(sdk.clientAdmin, contractAddress, denom)

    await sendCoins(sdk, sdk.clientValidator, sdk.validatorAddress, contractAddress, denom, amount)
    const latestBalance = await queryBalance(sdk.clientAdmin, contractAddress, denom)

    expect(latestBalance).toEqual(currBalance + amount);
  })
})
