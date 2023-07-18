import { Client } from '@mantrachain/sdk'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'
import * as dotenv from 'dotenv'
import { BlockWaiter } from './wait'

dotenv.config()

export class MantrachainSdk {
  clientValidator: any
  clientRecipient: any
  clientAdmin: any
  validatorWallet: DirectSecp256k1HdWallet
  validatorAddress: string
  recipientWallet: DirectSecp256k1HdWallet
  recipientAddress: string
  adminWallet: DirectSecp256k1HdWallet
  adminAddress: string
  blockWaiter: BlockWaiter

  async init(host = 'http://127.0.0.1:1317', rpc = 'http://127.0.0.1:26657', ws = 'ws://127.0.0.1:26657') {
    this.validatorWallet = await DirectSecp256k1HdWallet.fromMnemonic(process.env.VALIDATOR_MNEMONIC!, {
      prefix: "mantra",
    })
    this.recipientWallet = await DirectSecp256k1HdWallet.fromMnemonic(process.env.RECIPIENT_MNEMONIC!, {
      prefix: "mantra",
    })
    this.adminWallet = await DirectSecp256k1HdWallet.fromMnemonic(process.env.ADMIN_MNEMONIC!, {
      prefix: "mantra",
    })

    this.validatorAddress = (await this.validatorWallet.getAccounts())[0].address
    this.recipientAddress = (await this.recipientWallet.getAccounts())[0].address
    this.adminAddress = (await this.adminWallet.getAccounts())[0].address

    this.clientValidator = new Client(
      {
        apiURL: host,
        rpcURL: rpc,
      },
      this.validatorWallet
    )
    this.clientRecipient = new Client(
      {
        apiURL: host,
        rpcURL: rpc,
      },
      this.recipientWallet
    )
    this.clientAdmin = new Client(
      {
        apiURL: host,
        rpcURL: rpc,
      },
      this.adminWallet
    )

    this.blockWaiter = new BlockWaiter(ws);
  }
}

export const getGasFee = () => ({
  amount: [{ denom: "uaum", amount: "20" }],
  gas: "200000"
})