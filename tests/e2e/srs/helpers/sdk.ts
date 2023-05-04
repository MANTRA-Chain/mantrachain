import { Client } from '@mantrachain/sdk'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'
import * as dotenv from 'dotenv'
import { BlockWaiter } from './wait'

dotenv.config()

export class MantrachainSdk {
  clientAdmin: any
  clientRecipient: any
  adminWallet: DirectSecp256k1HdWallet
  adminAddress: string
  recipientWallet: DirectSecp256k1HdWallet
  recipientAddress: string
  blockWaiter: BlockWaiter

  async init(host = 'http://127.0.0.1:1317', rpc = 'http://127.0.0.1:26648', ws = 'ws://127.0.0.1:26657/websocket') {
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
        apiURL: host,
        rpcURL: rpc,
      },
      this.adminWallet
    )

    this.clientRecipient = new Client(
      {
        apiURL: host,
        rpcURL: rpc,
      },
      this.recipientWallet
    )

    this.blockWaiter = new BlockWaiter(ws);
  }
}