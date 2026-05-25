'use client'

import { useState } from 'react'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ProtectedPage from '@/components/ProtectedPage'
import { useAuth } from '@/lib/auth-context'
import { openCodeChat } from '@/lib/api'

const SUGGESTIONS = [
  'Explain this codebase architecture',
  'How do I add a new API endpoint?',
  'Find potential bugs in the authentication flow',
  'Suggest improvements to the deployment pipeline',
  'Generate a README for this project',
]

export default function OpenCodePage() {
  const { user, logout } = useAuth()
  const [input, setInput] = useState('')
  const [context, setContext] = useState('')
  const [reply, setReply] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!input.trim() || loading) return

    setLoading(true)
    setError('')
    setReply('')

    try {
      const result = await openCodeChat(input.trim(), context.trim())
      setReply(result)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : String(err))
    }
    setLoading(false)
  }

  return (
    <ProtectedPage>
      <div className="min-h-screen bg-gray-950 flex">
        <Sidebar />
        <div className="flex-1 flex flex-col">
          <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
            <h1 className="text-lg font-semibold text-white">OpenCode.ai</h1>
            <UserMenu username={user || 'admin'} onLogout={logout} />
          </header>

          <main className="flex-1 p-8 overflow-y-auto">
            {!reply && !loading && (
              <div className="mb-8">
                <p className="text-gray-400 text-sm mb-4">Suggestions:</p>
                <div className="flex flex-wrap gap-2">
                  {SUGGESTIONS.map((s, i) => (
                    <button
                      key={i}
                      onClick={() => setInput(s)}
                      className="px-4 py-2 bg-gray-800 hover:bg-gray-700 text-gray-300 text-sm rounded-lg transition"
                    >
                      {s}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {error && (
              <div className="mb-4 bg-red-900/50 border border-red-800 text-red-300 px-4 py-2.5 rounded-lg text-sm">
                {error}
              </div>
            )}

            {loading && (
              <div className="bg-gray-800 border border-gray-700 rounded-xl p-5 mb-4">
                <div className="animate-pulse text-gray-400 text-sm">Thinking...</div>
              </div>
            )}

            {reply && (
              <div className="bg-gray-800 border border-gray-700 rounded-xl p-5 mb-4">
                <p className="text-xs text-gray-500 mb-2">OpenCode.ai</p>
                <p className="text-sm text-gray-200 whitespace-pre-wrap">{reply}</p>
              </div>
            )}
          </main>

          <footer className="border-t border-gray-800 p-4">
            <form onSubmit={handleSubmit} className="max-w-3xl mx-auto space-y-2">
              <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Ask OpenCode.ai about your code..."
                disabled={loading}
                className="w-full px-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-vpsik-500"
              />
              <div className="flex gap-3">
                <input
                  type="text"
                  value={context}
                  onChange={(e) => setContext(e.target.value)}
                  placeholder="Optional context (file path, project name)..."
                  disabled={loading}
                  className="flex-1 px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-gray-100 placeholder-gray-500 text-sm focus:outline-none focus:ring-2 focus:ring-vpsik-500"
                />
                <button
                  type="submit"
                  disabled={!input.trim() || loading}
                  className="px-6 py-2 bg-vpsik-600 hover:bg-vpsik-500 disabled:bg-gray-700 text-white rounded-lg transition"
                >
                  {loading ? 'Sending...' : 'Send'}
                </button>
              </div>
            </form>
          </footer>
        </div>
      </div>
    </ProtectedPage>
  )
}
