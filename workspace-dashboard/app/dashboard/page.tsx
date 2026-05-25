'use client'

import { useEffect, useState } from 'react'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ServiceCard from '@/components/ServiceCard'
import StatusBadge from '@/components/StatusBadge'
import ProtectedPage from '@/components/ProtectedPage'
import { PageSkeleton } from '@/components/LoadingSkeleton'
import { useAuth } from '@/lib/auth-context'
import { getStatus as fetchStatus, ServiceStatus } from '@/lib/api'

const SERVICE_CARDS = [
  { name: 'Gitea', href: '/services/gitea', description: 'Git repositories & issues', icon: '📦', apiName: 'gitea' },
  { name: 'Ollama', href: '/services/ollama', description: 'AI models & management', icon: '🧠', apiName: 'ollama' },
  { name: 'Coolify', href: '/services/coolify', description: 'App deployments & servers', icon: '🚀', apiName: 'coolify' },
  { name: 'AI Chat', href: '/chat', description: 'Chat with local models', icon: '💬', apiName: '' },
  { name: 'OpenCode.ai', href: '/services/opencode', description: 'AI coding assistant', icon: '🤖', apiName: 'opencode' },
  { name: 'Open WebUI', href: '/services/openwebui', description: 'Full AI chat interface', icon: '🌐', apiName: '' },
]

export default function DashboardPage() {
  const { user, authenticated, logout } = useAuth()
  const [services, setServices] = useState<ServiceStatus[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!authenticated) return
    fetchStatus()
      .then(data => setServices(data.services))
      .catch(err => console.error('status fetch failed', err))
      .finally(() => setLoading(false))
  }, [authenticated])

  if (!authenticated || loading) {
    return <PageSkeleton />
  }

  const healthyCount = services.filter(s => s.status === 'healthy').length

  function getStatus(name: string): 'healthy' | 'unhealthy' | 'unknown' {
    const svc = services.find(s => s.name === name)
    if (!svc) return 'unknown'
    return svc.status === 'healthy' ? 'healthy' : 'unhealthy'
  }

  return (
    <ProtectedPage>
      <div className="min-h-screen bg-gray-950 flex">
        <Sidebar />
        <div className="flex-1">
          <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
            <div>
              <h1 className="text-lg font-semibold text-white">Dashboard</h1>
              <p className="text-sm text-gray-400">
                {healthyCount}/{services.length} services healthy
              </p>
            </div>
            <UserMenu username={user || 'admin'} onLogout={logout} />
          </header>

          <main className="p-8">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-10">
              {SERVICE_CARDS.map(card => (
                <ServiceCard
                  key={card.href}
                  name={card.name}
                  href={card.href}
                  description={card.description}
                  status={card.apiName ? getStatus(card.apiName) : 'healthy'}
                  icon={card.icon}
                />
              ))}
            </div>

            <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
              <h2 className="text-lg font-semibold text-white mb-4">Service Status</h2>
              <div className="space-y-3">
                {services.length === 0 && (
                  <p className="text-gray-500 text-sm">No services detected</p>
                )}
                {services.map((svc) => (
                  <div key={svc.name} className="flex items-center justify-between py-2 border-b border-gray-800 last:border-0">
                    <div className="flex items-center gap-3">
                      <span className="text-sm font-medium text-gray-200 capitalize">{svc.name}</span>
                      {svc.error && (
                        <span className="text-xs text-gray-500 truncate max-w-xs">{svc.error}</span>
                      )}
                    </div>
                    <StatusBadge status={svc.status === 'healthy' ? 'healthy' : 'unhealthy'} />
                  </div>
                ))}
              </div>
            </div>
          </main>
        </div>
      </div>
    </ProtectedPage>
  )
}
