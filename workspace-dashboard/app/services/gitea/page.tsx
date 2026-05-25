'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import { verify, getGiteaRepos, logout, getUser, GiteaRepo } from '@/lib/api'

export default function GiteaPage() {
  const [repos, setRepos] = useState<GiteaRepo[]>([])
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    async function init() {
      const ok = await verify()
      if (!ok) { router.push('/login'); return }
      setUser(getUser())
      try {
        const data = await getGiteaRepos()
        setRepos(data)
      } catch (err) {
        console.error('failed to fetch repos', err)
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
          <h1 className="text-lg font-semibold text-white">Gitea Repositories</h1>
          <UserMenu username={user || 'admin'} onLogout={() => { logout(); router.push('/login') }} />
        </header>
        <main className="p-8">
          {loading ? (
            <div className="animate-pulse text-gray-400">Loading...</div>
          ) : repos.length === 0 ? (
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
              <p className="text-gray-500">No repositories found</p>
            </div>
          ) : (
            <div className="grid gap-4">
              {repos.map((repo) => (
                <div key={repo.full_name} className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                  <div className="flex items-start justify-between">
                    <div>
                      <h3 className="text-white font-semibold">{repo.full_name}</h3>
                      <p className="text-sm text-gray-400 mt-1">{repo.description || 'No description'}</p>
                    </div>
                    <div className="flex items-center gap-3 text-sm text-gray-500">
                      {repo.language && (
                        <span className="px-2 py-0.5 bg-gray-800 rounded text-xs">{repo.language}</span>
                      )}
                      <span>★ {repo.stars_count}</span>
                      <span>⑂ {repo.forks_count}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
