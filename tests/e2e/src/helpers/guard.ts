import { AumegaSdk } from '../helpers/sdk'
import { getWithAttempts } from './wait'
import { utils, Privileges } from '@mantrachain/sdk'

const queryGuardTransferCoins = async (client: any) => {
  const res = await client.AumegaGuardV1.query.queryGuardTransferCoins()
  return res?.data?.guard_transfer_coins || false
}

const notSetGuardTransferCoins = (expected: boolean, actual: boolean) => expected !== actual

const queryDefaultPrivileges = async (client: any): Promise<Buffer> => {
  const res = await client.AumegaGuardV1.query.queryParams()
  return utils.base64ToBytes(
    res.data.params.default_privileges
  )
}

const queryCoinRequiredPrivileges = async (client: any, denom: string) => {
  try {
    const res = await client.AumegaGuardV1.query.queryRequiredPrivileges(
      utils.strToBase64(denom),
      {
        kind: "coin",
      }
    )

    return !!res?.data?.privileges
      ? utils.base64ToBytes(res.data.privileges)
      : Buffer.from([])
  } catch (e) {
    return Buffer.from([])
  }
}

const queryAccountPrivileges = async (client: any, account: string) => {
  const res = await client.AumegaGuardV1.query.queryAccountPrivileges(account)
  return utils.base64ToBytes(res.data.privileges)
}

export const setGuardTransferCoins = async (sdk: AumegaSdk, client: any, account: string, enabled: boolean, numAttempts = 2) => {
  if (notSetGuardTransferCoins(await queryGuardTransferCoins(client), enabled)) {
    const res = await client.AumegaGuardV1.tx.sendMsgUpdateGuardTransferCoins({
      value: {
        creator: account,
        enabled
      }
    })

    if (res.code !== 0) {
      throw new Error(res.rawLog)
    }
  } else {
    return
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryGuardTransferCoins(client),
    async (res) => !notSetGuardTransferCoins(res, enabled),
    numAttempts,
  )
}

export const updateAccountPrivileges = async (sdk: AumegaSdk, client: any, account: string, receiver: string, setBits?: number[], unsetBits?: number[], numAttempts = 2) => {
  const accountPrivileges: any = await queryAccountPrivileges(client, receiver)
  const defaultPrivileges = await queryDefaultPrivileges(client)
  let newAccountPrivileges: any = Buffer.from([])

  if (!!setBits?.length || !!unsetBits?.length) {
    newAccountPrivileges = Privileges.fromBuffer(accountPrivileges)

    if (!!setBits?.length) {
      setBits.forEach(bit => {
        newAccountPrivileges = newAccountPrivileges.set(bit)
      })
    }

    if (!!unsetBits?.length) {
      unsetBits.forEach(bit => {
        newAccountPrivileges = newAccountPrivileges.unset(bit)
      })
    }

    newAccountPrivileges = newAccountPrivileges.toBuffer()
  } else if (accountPrivileges.equals(defaultPrivileges)) {
    return
  }

  const res = await client.AumegaGuardV1.tx.sendMsgUpdateAccountPrivileges({
    value: {
      creator: account,
      account: receiver,
      privileges: !newAccountPrivileges.length ? null : newAccountPrivileges,
    }
  })

  if (res.code !== 0) {
    throw new Error(res.rawLog)
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryAccountPrivileges(client, receiver),
    async (privileges) => {
      if (!newAccountPrivileges.length) {
        return defaultPrivileges.equals(privileges)
      } else {
        return newAccountPrivileges.equals(privileges)
      }
    },
    numAttempts,
  )
}

export const updateCoinRequiredPrivileges = async (sdk: AumegaSdk, client: any, account: string, denom: string, setBits?: number[], unsetBits?: number[], numAttempts = 2) => {
  const requiredPrivileges: any = await queryCoinRequiredPrivileges(client, denom)
  const hasRequiredPrivileges = !!requiredPrivileges.length
  let newRequiredPrivileges: any = Buffer.from([])

  if (!setBits?.length && !hasRequiredPrivileges) {
    return
  }

  if (!!setBits?.length || !!unsetBits?.length) {
    if (!hasRequiredPrivileges) {
      newRequiredPrivileges = Privileges.fromBuffer(await queryDefaultPrivileges(client))
    } else {
      newRequiredPrivileges = Privileges.fromBuffer(requiredPrivileges)
    }

    if (!!setBits?.length) {
      setBits.forEach(bit => {
        newRequiredPrivileges = newRequiredPrivileges.set(bit)
      })
    }

    if (!!unsetBits?.length) {
      unsetBits.forEach(bit => {
        newRequiredPrivileges = newRequiredPrivileges.unset(bit)
      })
    }

    newRequiredPrivileges = newRequiredPrivileges.toBuffer()
  }

  const res = await client.AumegaGuardV1.tx.sendMsgUpdateRequiredPrivileges({
    value: {
      creator: account,
      index: utils.strToIndex(denom),
      privileges: !newRequiredPrivileges.length ? null : newRequiredPrivileges,
      kind: "coin",
    }
  })

  if (res.code !== 0) {
    throw new Error(res.rawLog)
  }

  return getWithAttempts(
    sdk.blockWaiter,
    async () => await queryCoinRequiredPrivileges(client, denom),
    async (privileges) => newRequiredPrivileges.equals(privileges),
    numAttempts,
  )
}