'use client'

import { useEffect, useState } from 'react'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ProtectedPage from '@/components/ProtectedPage'
import { ListSkeleton } from '@/components/LoadingSkeleton'
import { useAuth } from '@/lib/auth-context'
import { getStatus, ServiceStatus } from '@/lib/api'

const GRAFANA_URL = process.env.NEXT_PUBLIC_GRAFANA_URL || 'http://localhost:3002'
const PROMETHEUS_URL = process.env.NEXT_PUBLIC_PROMETHEUS_URL || 'http://localhost:9090'

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    ok: 'bg-green-900/50 text-green-400',
    healthy: 'bg-green-900/50 text-green-400',
    running: 'bg-green-900/50 text-green-400',
    error: 'bg-red-900/50 text-red-400',
    down: 'bg-red-900/50 text-red-400',
    degraded: 'bg-yellow-900/50 text-yellow-400',
    unknown: 'bg-gray-800 text-gray-400',
  }
  const cls = colors[status?.toLowerCase()] || colors.unknown
  return (
    <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${cls}`}>
      {status || 'unknown'}
    </span>
  )
}

export default function MonitoringPage() {
  const { user, authenticated, logout } = useAuth()
  const [services, setServices] = useState<ServiceStatus[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!authenticated) return
    getStatus()
      .then(data => setServices(data.services || []))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [authenticated])

  return (
    <ProtectedPage>
      <div className="min-h-screen bg-gray-950 flex">
        <Sidebar />
        <div className="flex-1">
          <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
            <h1 className="text-lg font-semibold text-white">Monitoring</h1>
            <UserMenu username={user || 'admin'} onLogout={logout} />
          </header>
          <main className="p-8 space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <a
                href={GRAFANA_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="bg-gray-900 border border-gray-800 rounded-xl p-6 hover:border-vpsik-600/50 transition group"
              >
                <h2 className="text-white font-semibold text-lg group-hover:text-vpsik-400">Grafana</h2>
                <p className="text-gray-400 text-sm mt-1">Dashboards, alerts, and metrics visualization</p>
                <span className="text-vpsik-400 text-xs mt-2 inline-block">Open Grafana →</span>
              </a>
              <a
                href={PROMETHEUS_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="bg-gray-900 border border-gray-800 rounded-xl p-6 hover:border-vpsik-600/50 transition group"
              >
                <h2 className="text-white font-semibold text-lg group-hover:text-vpsik-400">Prometheus</h2>
                <p className="text-gray-400 text-sm mt-1">Time-series metrics and alerting rules</p>
                <span className="text-vpsik-400 text-xs mt-2 inline-block">Open Prometheus →</span>
              </a>
            </div>

            <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <h2 className="text-white font-semibold mb-4">Service Health</h2>
              {loading ? (
                <ListSkeleton rows={5} />
              ) : error ? (
                <p className="text-red-400 text-sm">{error}</p>
              ) : services.length === 0 ? (
                <p className="text-gray-500 text-sm">No services reported</p>
              ) : (
                <div className="space-y-2">
                  {services.map(svc => (
                    <div key={svc.name} className="flex items-center justify-between py-2 border-b border-gray-800 last:border-0">
                      <span className="text-gray-200 text-sm">{svc.name}</span>
                      <div className="flex items-center gap-2">
                        <StatusBadge status={svc.status} />
                        {svc.error && <span className="text-red-400 text-xs">{svc.error}</span>}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </main>
        </div>
      </div>
    </ProtectedPage>
  )
}
