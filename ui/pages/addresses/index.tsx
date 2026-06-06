import { Wallet, Plus, Copy, ExternalLink } from 'lucide-react'
import DefaultLayout from '@/layouts/default'

const addresses = [
  { label: 'Primary Wallet', address: '8x2K...3mNp', balance: '42.5 SOL' },
  { label: 'Treasury', address: '9y4R...7qVs', balance: '125.0 SOL' },
  { label: 'Hot Wallet', address: '3a1B...6cDf', balance: '8.2 SOL' },
]

const AddressesPage = () => {
  return (
    <DefaultLayout>
      <div className="flex flex-col gap-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">Addresses</h1>
            <p className="text-muted mt-1 text-sm">Manage your Solana wallet addresses</p>
          </div>
          <button className="bg-accent text-accent-foreground hover:bg-accent/90 flex items-center gap-2 rounded-xl px-4 py-2 text-sm font-medium transition-colors">
            <Plus size={16} />
            Add Address
          </button>
        </div>

        <div className="grid gap-4">
          {addresses.map(({ label, address, balance }) => (
            <div key={label} className="bg-surface flex items-center justify-between rounded-2xl p-5">
              <div className="flex items-center gap-4">
                <div className="bg-accent/10 flex h-10 w-10 items-center justify-center rounded-xl">
                  <Wallet size={20} className="text-accent" />
                </div>
                <div>
                  <p className="text-sm font-medium">{label}</p>
                  <p className="text-muted mt-0.5 font-mono text-xs">{address}</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <span className="text-sm font-semibold">{balance}</span>
                <button className="text-muted hover:text-foreground transition-colors">
                  <Copy size={16} />
                </button>
                <button className="text-muted hover:text-foreground transition-colors">
                  <ExternalLink size={16} />
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </DefaultLayout>
  )
}

export default AddressesPage
