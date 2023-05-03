import { Client } from '@mantrachain/sdk'
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing"
import * as dotenv from 'dotenv'

dotenv.config()

export class MantrachainSdk {
  client: any
  recipientWallet: DirectSecp256k1HdWallet
  recipientAddress: string

  async init() {
    this.recipientWallet = await DirectSecp256k1HdWallet.fromMnemonic(process.env.RECIPIENT_MNEMONIC!, {
      prefix: "mantra",
    })

    this.recipientAddress = (await this.recipientWallet.getAccounts())[0].address

    this.client = new Client(
      {
        apiURL: "http://127.0.0.1:1317",
        rpcURL: "http://127.0.0.1:26648",
      },
      this.recipientWallet
    )
  }
}