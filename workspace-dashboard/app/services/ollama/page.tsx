'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import { verify, getOllamaModels, logout, getUser, OllamaModel } from '@/lib/api'

export default function OllamaPage() {
  const [models, setModels] = useState<OllamaModel[]>([])
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    async function init() {
      const ok = await verify()
      if (!ok) { router.push('/login'); return }
      setUser(getUser())
      try {
        const data = await getOllamaModels()
        setModels(data)
      } catch (err) {
        console.error('failed to fetch models', err)
      }
      setLoading(false)
    }
    init()
  }, [router])

  function formatSize(bytes: number): string {
    if (bytes >= 1e9) return `${(bytes / 1e9).toFixed(1)} GB`
    if (bytes >= 1e6) return `${(bytes / 1e6).toFixed(0)} MB`
    return `${(bytes / 1e3).toFixed(0)} KB`
  }

  return (
    <div className="min-h-screen bg-gray-950 flex">
      <Sidebar />
      <div className="flex-1">
        <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
          <h1 className="text-lg font-semibold text-white">Ollama Models</h1>
          <UserMenu username={user || 'admin'} onLogout={() => { logout(); router.push('/login') }} />
        </header>
        <main className="p-8">
          {loading ? (
            <div className="animate-pulse text-gray-400">Loading...</div>
          ) : models.length === 0 ? (
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
              <p className="text-gray-500">No models found</p>
            </div>
          ) : (
            <div className="grid gap-4">
              {models.map((model) => (
                <div key={model.name} className="bg-gray-900 border border-gray-800 rounded-xl p-5 flex items-center justify-between">
                  <div>
                    <h3 className="text-white font-semibold">{model.name}</h3>
                    <p className="text-xs text-gray-500 mt-1">Digest: {model.digest.slice(0, 19)}...</p>
                  </div>
                  <span className="text-sm text-gray-400">{formatSize(model.size)}</span>
                </div>
              ))}
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
