import { DollarSign, ShoppingCart, Users, Wallet, TrendingUp, TrendingDown, Plus, List, UserCog, ArrowUpRight, ArrowDownLeft, ChevronRight } from 'lucide-react'
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts'
import DefaultLayout from '@/layouts/default'

const stats = [
  {
    label: 'Total Revenue',
    value: '$12,345',
    change: '+12.3%',
    icon: DollarSign,
    accent: 'text-success',
    bgAccent: 'bg-success/10',
  },
  {
    label: 'Total Orders',
    value: '1,234',
    change: '+8.1%',
    icon: ShoppingCart,
    accent: 'text-warning',
    bgAccent: 'bg-warning/10',
  },
  {
    label: 'Active Addresses',
    value: '5,678',
    icon: Users,
    accent: 'text-accent',
    bgAccent: 'bg-accent/10',
  },
  {
    label: 'Avg Transaction Value',
    value: '$89.50',
    change: '-2.1%',
    icon: Wallet,
    accent: 'text-danger',
    bgAccent: 'bg-danger/10',
  },
]

const chartData = [
  { day: 'Mon', amount: 1200 },
  { day: 'Tue', amount: 1800 },
  { day: 'Wed', amount: 1400 },
  { day: 'Thu', amount: 2200 },
  { day: 'Fri', amount: 1900 },
  { day: 'Sat', amount: 2600 },
  { day: 'Sun', amount: 1500 },
]

const IndexPage = () => {
  return (
    <DefaultLayout>
      <div className="flex flex-col gap-6">
        <div className="grid grid-cols-4 gap-4">
          {stats.map(({ label, value, change, icon: Icon, accent, bgAccent }) => (
            <div key={label} className="bg-surface flex flex-col gap-3 rounded-2xl p-5">
              <div className="flex items-center gap-3">
                <div className={`${bgAccent} flex h-10 w-10 items-center justify-center rounded-xl`}>
                  <Icon size={20} className={accent} />
                </div>
                <span className="text-muted text-sm">{label}</span>
              </div>
              <div className="flex items-baseline justify-between">
                <span className="text-3xl font-semibold tracking-tight">{value}</span>
                {change && (
                  <span
                    className={`flex items-center gap-0.5 rounded-full px-2 py-0.5 text-xs font-medium ${change.startsWith('+') ? 'bg-success/10 text-success' : 'bg-danger/10 text-danger'}`}
                  >
                    {change.startsWith('+') ? (
                      <TrendingUp size={12} />
                    ) : (
                      <TrendingDown size={12} />
                    )}
                    {change}
                  </span>
                )}
              </div>
            </div>
          ))}
        </div>

        <div className="bg-surface rounded-2xl p-6">
          <h2 className="mb-1 text-base font-semibold">Payment Activity</h2>
          <p className="text-muted mb-6 text-sm">Daily payment volume for the last 7 days</p>
          <ResponsiveContainer width="100%" height={280} className="outline-none">
            <LineChart data={chartData} margin={{ top: 10, right: 10, bottom: 10, left: 10 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="var(--foreground)" strokeOpacity={0.1} vertical={false} />
              <XAxis
                dataKey="day"
                axisLine={false}
                tickLine={false}
                tick={{ fill: 'currentColor', fontSize: 13, opacity: 0.5 }}
                dy={10}
              />
              <YAxis
                axisLine={false}
                tickLine={false}
                tick={{ fill: 'currentColor', fontSize: 13, opacity: 0.5 }}
                dx={-10}
                tickFormatter={(v) => `$${v}`}
              />
              <Tooltip
                cursor={{ stroke: 'var(--foreground)', strokeOpacity: 0.1, strokeDasharray: '3 3', strokeWidth: 1 }}
                contentStyle={{
                  background: 'var(--surface)',
                  border: '1px solid var(--separator)',
                  borderRadius: '12px',
                  fontSize: '13px',
                }}
                labelStyle={{ fontWeight: 600 }}
                formatter={(value) => value ? [`$${value}`, 'Volume'] : ['$0', 'Volume']}
              />
              <Line
                type="monotone"
                dataKey="amount"
                stroke="var(--accent)"
                strokeWidth={2}
                dot={{ fill: 'var(--surface)', stroke: 'var(--accent)', strokeWidth: 2, r: 5 }}
                activeDot={{ fill: 'var(--accent)', r: 7, stroke: 'var(--accent)', strokeWidth: 0 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="bg-surface rounded-2xl p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-base font-semibold">Recent Transactions</h2>
              <button className="text-accent text-xs font-medium hover:underline">View all</button>
            </div>
            <div className="flex flex-col">
              {[
                { from: '0x1a2B...c3d4', to: '0x5e6F...a7b8', amount: '+2.5 SOL', time: '2 min ago', incoming: true },
                { from: '0x9c0D...e1f2', to: '0x3g4H...i5j6', amount: '-0.8 SOL', time: '15 min ago', incoming: false },
                { from: '0x7k8L...m9n0', to: '0x1o2P...q3r4', amount: '+5.1 SOL', time: '1 hour ago', incoming: true },
              ].map((tx, i) => (
                <div key={i}>
                  {i > 0 && <div className="border-separator mx-1 border-t" />}
                  <div className="flex items-center gap-3 px-1 py-3">
                    <div className={`relative flex h-9 w-9 items-center justify-center rounded-full ${tx.incoming ? 'bg-accent/10' : 'bg-danger/10'}`}>
                      {tx.incoming ? (
                        <ArrowDownLeft size={15} className="text-accent" />
                      ) : (
                        <ArrowUpRight size={15} className="text-danger" />
                      )}
                    </div>
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm font-medium">{tx.from} → {tx.to}</p>
                      <p className="text-muted text-xs">{tx.time}</p>
                    </div>
                    <span className={`whitespace-nowrap text-sm font-semibold ${tx.incoming ? 'text-success' : 'text-danger'}`}>{tx.amount}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-surface rounded-2xl p-6">
            <h2 className="mb-4 text-base font-semibold">Quick Actions</h2>
            <div className="flex flex-col gap-2">
              {[
                { title: 'Add Address', subtitle: 'Set up an active address', icon: Plus, accent: 'text-accent', bg: 'bg-accent/10' },
                { title: 'View Orders', subtitle: 'Check all your payment orders', icon: List, accent: 'text-warning', bg: 'bg-warning/10' },
                { title: 'Update Profile', subtitle: 'Configure your settings', icon: UserCog, accent: 'text-success', bg: 'bg-success/10' },
              ].map(({ title, subtitle, icon: Icon, accent, bg }) => (
                <button
                  key={title}
                  className="group flex cursor-pointer items-center gap-3 rounded-xl px-3 py-3 text-left transition-colors hover:bg-background"
                >
                  <div className={`${bg} flex h-10 w-10 items-center justify-center rounded-xl`}>
                    <Icon size={20} className={accent} />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium">{title}</p>
                    <p className="text-muted truncate text-xs">{subtitle}</p>
                  </div>
                  <ChevronRight size={16} className="text-muted opacity-0 transition-opacity group-hover:opacity-100" />
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>
    </DefaultLayout>
  )
}

export default IndexPage
