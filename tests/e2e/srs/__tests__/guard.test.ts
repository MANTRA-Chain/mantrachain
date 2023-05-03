import { MantrachainSdk } from '../helpers/sdk'

describe('Guard module', () => {
  let sdk: MantrachainSdk

  beforeAll(async () => {
    sdk = new MantrachainSdk()
    await sdk.init()
  })

  describe('Not Authenticated', () => {
    test('Should throw when create denom with account with no permission', async () => {
      const promise = sdk.client.MantrachainCoinfactoryV1Beta1.tx.sendMsgCreateDenom({
        value: {
          sender: sdk.recipientAddress,
          subdenom: 'testcoin'
        }
      })

      return expect(promise).rejects.toThrow(
        /unauthorized/
      );
    });
  })
})