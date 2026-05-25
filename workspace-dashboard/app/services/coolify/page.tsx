'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import { verify, getCoolifyServers, logout, getUser, CoolifyServer } from '@/lib/api'

export default function CoolifyPage() {
  const [servers, setServers] = useState<CoolifyServer[]>([])
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    async function init() {
      const ok = await verify()
      if (!ok) { router.push('/login'); return }
      setUser(getUser())
      try {
        const data = await getCoolifyServers()
        setServers(data)
      } catch (err) {
        console.error('failed to fetch servers', err)
      }
      setLoading(false)
    }
    init()
  }, [router])

  return (
    <div className="min-h-screen bg-gray-950 flex">
      <Sidebar />
      <div className="flex-1">
        <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
          <h1 className="text-lg font-semibold text-white">Coolify Servers</h1>
          <UserMenu username={user || 'admin'} onLogout={() => { logout(); router.push('/login') }} />
        </header>
        <main className="p-8">
          {loading ? (
            <div className="animate-pulse text-gray-400">Loading...</div>
          ) : servers.length === 0 ? (
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
              <p className="text-gray-500">No servers found</p>
            </div>
          ) : (
            <div className="grid gap-4">
              {servers.map((server) => (
                <div key={server.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5 flex items-center justify-between">
                  <div>
                    <h3 className="text-white font-semibold">{server.name}</h3>
                    <p className="text-sm text-gray-400 mt-1">{server.ip}</p>
                  </div>
                  <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    server.status === 'ready'
                      ? 'bg-green-900/50 text-green-400'
                      : 'bg-yellow-900/50 text-yellow-400'
                  }`}>
                    {server.status}
                  </span>
                </div>
              ))}
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
