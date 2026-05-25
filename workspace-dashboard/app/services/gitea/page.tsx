'use client'

import { useEffect, useState } from 'react'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ProtectedPage from '@/components/ProtectedPage'
import { ListSkeleton } from '@/components/LoadingSkeleton'
import { useAuth } from '@/lib/auth-context'
import { getGiteaRepos, GiteaRepo, API_BASE, authHeaders } from '@/lib/api'

interface WebhookForm {
  repo: string
  url: string
  secret: string
}

export default function GiteaPage() {
  const { user, authenticated, logout } = useAuth()
  const [repos, setRepos] = useState<GiteaRepo[]>([])
  const [loading, setLoading] = useState(true)
  const [showWebhook, setShowWebhook] = useState(false)
  const [webhook, setWebhook] = useState<WebhookForm>({ repo: '', url: '', secret: '' })
  const [webhookMsg, setWebhookMsg] = useState('')

  useEffect(() => {
    if (!authenticated) return
    getGiteaRepos()
      .then(setRepos)
      .catch(err => console.error('failed to fetch repos', err))
      .finally(() => setLoading(false))
  }, [authenticated])

  async function handleCreateWebhook(e: React.FormEvent) {
    e.preventDefault()
    setWebhookMsg('')
    try {
      const res = await fetch(`${API_BASE}/api/gitea/webhooks`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({
          repo: webhook.repo,
          url: webhook.url,
          secret: webhook.secret,
          events: ['push', 'pull_request'],
        }),
      })
      if (!res.ok) throw new Error('webhook creation failed')
      setWebhookMsg('✅ Webhook created')
      setShowWebhook(false)
    } catch (err: unknown) {
      setWebhookMsg(`❌ ${err instanceof Error ? err.message : String(err)}`)
    }
  }

  return (
    <ProtectedPage>
      <div className="min-h-screen bg-gray-950 flex">
        <Sidebar />
        <div className="flex-1">
          <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
            <h1 className="text-lg font-semibold text-white">Gitea Repositories</h1>
            <UserMenu username={user || 'admin'} onLogout={logout} />
          </header>
          <main className="p-8">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-gray-300 font-medium">Repositories</h2>
              <button
                onClick={() => setShowWebhook(!showWebhook)}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 text-gray-300 text-sm rounded-lg transition"
              >
                {showWebhook ? 'Close' : '+ Webhook'}
              </button>
            </div>

            {showWebhook && (
              <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 mb-6">
                <h3 className="text-white font-medium mb-4">Create Webhook</h3>
                <form onSubmit={handleCreateWebhook} className="space-y-3">
                  <select
                    value={webhook.repo}
                    onChange={(e) => setWebhook({ ...webhook, repo: e.target.value })}
                    className="w-full px-4 py-2 bg-gray-800 border border-gray-700 text-gray-200 rounded-lg text-sm"
                    required
                  >
                    <option value="">Select repository</option>
                    {repos.map(r => (
                      <option key={r.full_name} value={r.full_name}>{r.full_name}</option>
                    ))}
                  </select>
                  <input
                    type="url"
                    value={webhook.url}
                    onChange={(e) => setWebhook({ ...webhook, url: e.target.value })}
                    placeholder="Webhook URL (e.g., http://coolify:8000/api/v1/deploy)"
                    className="w-full px-4 py-2 bg-gray-800 border border-gray-700 text-gray-200 rounded-lg text-sm"
                    required
                  />
                  <input
                    type="text"
                    value={webhook.secret}
                    onChange={(e) => setWebhook({ ...webhook, secret: e.target.value })}
                    placeholder="Secret (optional)"
                    className="w-full px-4 py-2 bg-gray-800 border border-gray-700 text-gray-200 rounded-lg text-sm"
                  />
                  <button
                    type="submit"
                    className="px-4 py-2 bg-vpsik-600 hover:bg-vpsik-500 text-white text-sm rounded-lg transition"
                  >
                    Create Webhook
                  </button>
                  {webhookMsg && (
                    <p className="text-sm text-gray-400">{webhookMsg}</p>
                  )}
                </form>
              </div>
            )}

            {loading ? (
              <ListSkeleton rows={4} />
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
    </ProtectedPage>
  )
}
