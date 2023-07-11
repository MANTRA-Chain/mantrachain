import axios from 'axios'
import { ChannelsList } from './types'
import { wait } from './wait'

const BLOCKS_COUNT_BEFORE_START = 10

export const setup = async (host: string) => {
  await waitForHTTP(host)
}

const waitForHTTP = async (
  host = 'http://127.0.0.1:1317',
  path = `blocks/${BLOCKS_COUNT_BEFORE_START}`,
  timeout = 280000,
) => {
  const start = Date.now()
  console.log('Waiting for client to start...')
  while (Date.now() < start + timeout) {
    try {
      const r = await axios.get(`${host}/${path}`, {
        timeout: 1000,
      })
      if (r.status === 200) {
        return
      }
      // eslint-disable-next-line no-empty
    } catch (e) { }
    await wait(10)
  }
  throw new Error('No port opened')
}

export const waitForChannel = async (
  host = 'http://127.0.0.1:1317',
  timeout = 100000,
) => {
  const start = Date.now()

  while (Date.now() < start + timeout) {
    try {
      const r = await axios.get<ChannelsList>(
        `${host}/ibc/core/channel/v1/channels`,
        {
          timeout: 1000,
        },
      )
      if (
        r.data.channels.length > 0 &&
        r.data.channels.every(
          (channel) => channel.counterparty.channel_id !== '',
        )
      ) {
        await wait(20)
        return
      }
      // eslint-disable-next-line no-empty
    } catch (e) { }
    await wait(10)
  }

  throw new Error('No channel opened')
}
