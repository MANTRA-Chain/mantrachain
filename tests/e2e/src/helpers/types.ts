export type ChannelsList = {
  channels: {
    state: string
    ordering: string
    counterparty: {
      port_id: string
      channel_id: string
    }
    connection_hops: string[]
    version: string
    port_id: string
    channel_id: string
  }[]
}