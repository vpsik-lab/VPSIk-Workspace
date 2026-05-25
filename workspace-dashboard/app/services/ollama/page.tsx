'use client'

import { useEffect, useState } from 'react'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ProtectedPage from '@/components/ProtectedPage'
import { ListSkeleton } from '@/components/LoadingSkeleton'
import { useAuth } from '@/lib/auth-context'
import { getOllamaModels, pullOllamaModel, deleteOllamaModel, OllamaModel } from '@/lib/api'

export default function OllamaPage() {
  const { user, authenticated, logout } = useAuth()
  const [models, setModels] = useState<OllamaModel[]>([])
  const [loading, setLoading] = useState(true)
  const [pullName, setPullName] = useState('')
  const [pulling, setPulling] = useState(false)
  const [error, setError] = useState('')

  async function loadModels() {
    setLoading(true)
    try {
      const data = await getOllamaModels()
      setModels(data)
    } catch (err) {
      console.error('failed to fetch models', err)
    }
    setLoading(false)
  }

  useEffect(() => {
    if (!authenticated) return
    loadModels()
  }, [authenticated])

  async function handlePull(e: React.FormEvent) {
    e.preventDefault()
    if (!pullName.trim() || pulling) return
    setPulling(true)
    setError('')
    try {
      await pullOllamaModel(pullName.trim())
      setPullName('')
      await loadModels()
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : String(err))
    }
    setPulling(false)
  }

  async function handleDelete(modelName: string) {
    if (!confirm(`Delete model "${modelName}"?`)) return
    try {
      await deleteOllamaModel(modelName)
      await loadModels()
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : String(err))
    }
  }

  function formatSize(bytes: number): string {
    const units = ['B', 'KB', 'MB', 'GB']
    let i = 0
    let size = bytes
    while (size >= 1024 && i < units.length - 1) {
      size /= 1024
      i++
    }
    return `${size.toFixed(i === 0 ? 0 : 1)} ${units[i]}`
  }

  return (
    <ProtectedPage>
      <div className="min-h-screen bg-gray-950 flex">
        <Sidebar />
        <div className="flex-1">
          <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
            <h1 className="text-lg font-semibold text-white">Ollama Models</h1>
            <UserMenu username={user || 'admin'} onLogout={logout} />
          </header>
          <main className="p-8">
            <form onSubmit={handlePull} className="mb-8 flex gap-3">
              <input
                type="text"
                value={pullName}
                onChange={(e) => setPullName(e.target.value)}
                placeholder="Pull a model (e.g., llama3.2, mistral)..."
                disabled={pulling}
                className="flex-1 px-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-vpsik-500"
              />
              <button
                type="submit"
                disabled={!pullName.trim() || pulling}
                className="px-6 py-2.5 bg-vpsik-600 hover:bg-vpsik-500 disabled:bg-gray-700 text-white font-medium rounded-lg transition"
              >
                {pulling ? 'Pulling...' : 'Pull'}
              </button>
            </form>

            {error && (
              <div className="mb-6 bg-red-900/50 border border-red-800 text-red-300 px-4 py-2.5 rounded-lg text-sm">
                {error}
              </div>
            )}

            {loading ? (
              <ListSkeleton rows={4} />
            ) : models.length === 0 ? (
              <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
                <p className="text-gray-500">No models found. Pull one above.</p>
              </div>
            ) : (
              <div className="grid gap-4">
                {models.map((model) => (
                  <div key={model.name} className="bg-gray-900 border border-gray-800 rounded-xl p-5 flex items-center justify-between">
                    <div>
                      <h3 className="text-white font-semibold">{model.name}</h3>
                      <p className="text-xs text-gray-500 mt-1">Digest: {model.digest.slice(0, 19)}...</p>
                    </div>
                    <div className="flex items-center gap-4">
                      <span className="text-sm text-gray-400">{formatSize(model.size)}</span>
                      <button
                        onClick={() => handleDelete(model.name)}
                        className="text-red-400 hover:text-red-300 text-sm transition"
                        title="Delete model"
                      >
                        🗑
                      </button>
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
