export interface HeaderBalances {
  sol: number
  usd: string
}

export interface ApiStatus {
  online: boolean
  helius_configured: boolean
  address_set: boolean
  webhook_set: boolean
  label: string
}
