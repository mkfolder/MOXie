import { ArrowUpRight, ArrowDownLeft, ExternalLink } from 'lucide-react'

import DefaultLayout from '@/layouts/default'

const transactions = [
  { from: '0x1a2B...c3d4', to: '0x5e6F...a7b8', amount: '+2.5 SOL', date: '2 min ago', incoming: true },
  { from: '0x9c0D...e1f2', to: '0x3g4H...i5j6', amount: '-0.8 SOL', date: '15 min ago', incoming: false },
  { from: '0x7k8L...m9n0', to: '0x1o2P...q3r4', amount: '+5.1 SOL', date: '1 hour ago', incoming: true },
  { from: '0x4f5G...h6i7', to: '0x8j9K...l0m1', amount: '-1.2 SOL', date: '3 hours ago', incoming: false },
  { from: '0x2n3O...p4q5', to: '0x6r7S...t8u9', amount: '+0.5 SOL', date: '5 hours ago', incoming: true },
]

const TransactionsPage = () => {
  return (
    <DefaultLayout>
      <div className="flex flex-col gap-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Transactions</h1>
          <p className="text-muted mt-1 text-sm">Every SOL accounted for</p>
        </div>

        <div className="bg-surface rounded-2xl">
          <div className="text-muted hidden grid-cols-4 gap-4 border-b px-6 py-4 text-xs font-medium tracking-wider uppercase md:grid">
            <span>From</span>
            <span>To</span>
            <span>Amount</span>
            <span className="text-right">Date</span>
          </div>
          {transactions.map((tx, i) => (
            <div key={i} className="grid grid-cols-4 items-center gap-4 px-6 py-4">
              <div className="flex items-center gap-3">
                <div
                  className={`flex h-8 w-8 items-center justify-center rounded-full ${tx.incoming ? 'bg-accent/10' : 'bg-danger/10'}`}
                >
                  {tx.incoming ? (
                    <ArrowDownLeft className="text-accent" size={14} />
                  ) : (
                    <ArrowUpRight className="text-danger" size={14} />
                  )}
                </div>
                <span className="font-mono text-sm">{tx.from}</span>
              </div>
              <span className="font-mono text-sm">{tx.to}</span>
              <span className={`text-sm font-semibold ${tx.incoming ? 'text-success' : 'text-danger'}`}>
                {tx.amount}
              </span>
              <div className="flex items-center justify-end gap-2">
                <span className="text-muted text-sm">{tx.date}</span>
                <button className="text-muted hover:text-foreground transition-colors">
                  <ExternalLink size={14} />
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </DefaultLayout>
  )
}

export default TransactionsPage
