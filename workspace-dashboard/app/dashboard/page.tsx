'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ServiceCard from '@/components/ServiceCard'
import StatusBadge from '@/components/StatusBadge'
import { verify, getStatus, logout, getUser, ServiceStatus } from '@/lib/api'

export default function DashboardPage() {
  const [authenticated, setAuthenticated] = useState(false)
  const [services, setServices] = useState<ServiceStatus[]>([])
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    async function init() {
      const ok = await verify()
      if (!ok) {
        router.push('/login')
        return
      }
      setAuthenticated(true)
      setUser(getUser())
      try {
        const data = await getStatus()
        setServices(data.services)
      } catch (err) {
        console.error('status fetch failed', err)
      }
      setLoading(false)
    }
    init()
  }, [router])

  async function handleLogout() {
    await logout()
    router.push('/login')
  }

  if (!authenticated || loading) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="animate-pulse text-gray-400">Loading...</div>
      </div>
    )
  }

  const healthyCount = services.filter(s => s.status === 'healthy').length

  return (
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
          <UserMenu username={user || 'admin'} onLogout={handleLogout} />
        </header>

        <main className="p-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-10">
            <ServiceCard
              name="Gitea"
              href="/services/gitea"
              description="Git repositories & issues"
              status={services.find(s => s.name === 'gitea')?.status || 'unknown'}
              icon="📦"
            />
            <ServiceCard
              name="Ollama"
              href="/services/ollama"
              description="AI models & chat"
              status={services.find(s => s.name === 'ollama')?.status || 'unknown'}
              icon="🧠"
            />
            <ServiceCard
              name="Coolify"
              href="/services/coolify"
              description="App deployments & servers"
              status={services.find(s => s.name === 'coolify')?.status || 'unknown'}
              icon="🚀"
            />
            <ServiceCard
              name="AI Chat"
              href="/chat"
              description="Chat with local models"
              status="healthy"
              icon="💬"
            />
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
  )
}
