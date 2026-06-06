import { User, Settings, Key, Bell, Shield } from 'lucide-react'
import DefaultLayout from '@/layouts/default'

const sections = [
  { label: 'Profile', icon: User, description: 'Manage your personal information' },
  { label: 'API Keys', icon: Key, description: 'Manage your API keys and access tokens' },
  { label: 'Notifications', icon: Bell, description: 'Configure your notification preferences' },
  { label: 'Security', icon: Shield, description: 'Security settings and two-factor authentication' },
]

const ProfilePage = () => {
  return (
    <DefaultLayout>
      <div className="flex flex-col gap-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Profile</h1>
          <p className="text-muted mt-1 text-sm">All about you</p>
        </div>

        <div className="bg-surface rounded-2xl p-6">
          <div className="flex items-center gap-4">
            <div className="bg-accent/10 flex h-16 w-16 items-center justify-center rounded-2xl">
              <User size={28} className="text-accent" />
            </div>
            <div>
              <p className="text-lg font-semibold">johndoe</p>
              <p className="text-muted text-sm">john@example.com</p>
            </div>
          </div>
        </div>

        <div className="grid gap-3">
          {sections.map(({ label, icon: Icon, description }) => (
            <button
              key={label}
              className="bg-surface hover:bg-surface/80 flex cursor-pointer items-center gap-4 rounded-2xl p-5 text-left transition-colors"
            >
              <div className="bg-primary/10 flex h-10 w-10 items-center justify-center rounded-xl">
                <Icon size={20} className="text-primary" />
              </div>
              <div>
                <p className="text-sm font-medium">{label}</p>
                <p className="text-muted mt-0.5 text-xs">{description}</p>
              </div>
              <Settings size={16} className="text-muted ml-auto" />
            </button>
          ))}
        </div>
      </div>
    </DefaultLayout>
  )
}

export default ProfilePage
