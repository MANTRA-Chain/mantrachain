import { Client } from '@mantrachain/sdk'
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing"
import * as dotenv from 'dotenv'

dotenv.config()

export class MantrachainSdk {
  clientAdmin: any
  clientRecipient: any
  adminWallet: DirectSecp256k1HdWallet
  adminAddress: string
  recipientWallet: DirectSecp256k1HdWallet
  recipientAddress: string

  async init() {
    this.adminWallet = await DirectSecp256k1HdWallet.fromMnemonic(process.env.ADMIN_MNEMONIC!, {
      prefix: "mantra",
    })
    this.recipientWallet = await DirectSecp256k1HdWallet.fromMnemonic(process.env.RECIPIENT_MNEMONIC!, {
      prefix: "mantra",
    })

    this.adminAddress = (await this.adminWallet.getAccounts())[0].address
    this.recipientAddress = (await this.recipientWallet.getAccounts())[0].address

    this.clientAdmin = new Client(
      {
        apiURL: "http://127.0.0.1:1317",
        rpcURL: "http://127.0.0.1:26648",
      },
      this.adminWallet
    )

    this.clientRecipient = new Client(
      {
        apiURL: "http://127.0.0.1:1317",
        rpcURL: "http://127.0.0.1:26648",
      },
      this.recipientWallet
    )
  }
}