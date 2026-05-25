'use client'

import { useEffect, useState, useRef } from 'react'
import { useRouter } from 'next/navigation'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import { verify, getOllamaModels, chatOllama, logout, getUser, OllamaModel } from '@/lib/api'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

export default function ChatPage() {
  const [models, setModels] = useState<OllamaModel[]>([])
  const [selectedModel, setSelectedModel] = useState('')
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [sending, setSending] = useState(false)
  const [user, setUser] = useState<string | null>(null)
  const router = useRouter()
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    async function init() {
      const ok = await verify()
      if (!ok) { router.push('/login'); return }
      setUser(getUser())
      try {
        const data = await getOllamaModels()
        setModels(data)
        if (data.length > 0) setSelectedModel(data[0].name)
      } catch (err) {
        console.error('failed to fetch models', err)
      }
    }
    init()
  }, [router])

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  async function handleSend(e: React.FormEvent) {
    e.preventDefault()
    if (!input.trim() || !selectedModel || sending) return

    const userMsg: Message = { role: 'user', content: input.trim() }
    setMessages(prev => [...prev, userMsg])
    setInput('')
    setSending(true)

    try {
      const apiMessages = messages.concat(userMsg).map(m => ({
        role: m.role,
        content: m.content,
      }))
      const reply = await chatOllama(selectedModel, apiMessages)
      setMessages(prev => [...prev, { role: 'assistant', content: reply }])
    } catch (err: any) {
      setMessages(prev => [...prev, { role: 'assistant', content: `Error: ${err.message}` }])
    }
    setSending(false)
  }

  return (
    <div className="min-h-screen bg-gray-950 flex">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
          <div className="flex items-center gap-4">
            <h1 className="text-lg font-semibold text-white">AI Chat</h1>
            <select
              value={selectedModel}
              onChange={(e) => setSelectedModel(e.target.value)}
              className="bg-gray-800 border border-gray-700 text-gray-200 text-sm rounded-lg px-3 py-1.5"
            >
              {models.map(m => (
                <option key={m.name} value={m.name}>{m.name}</option>
              ))}
            </select>
          </div>
          <UserMenu username={user || 'admin'} onLogout={() => { logout(); router.push('/login') }} />
        </header>

        <div className="flex-1 overflow-y-auto p-8 space-y-4">
          {messages.length === 0 && (
            <div className="text-center text-gray-500 mt-20">
              <p className="text-4xl mb-4">💬</p>
              <p>Start a conversation with your local AI</p>
            </div>
          )}
          {messages.map((msg, i) => (
            <div key={i} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
              <div className={`max-w-xl rounded-xl px-5 py-3 ${
                msg.role === 'user'
                  ? 'bg-vpsik-600/20 text-gray-100 border border-vpsik-600/30'
                  : 'bg-gray-800 text-gray-200 border border-gray-700'
              }`}>
                <p className="text-xs text-gray-500 mb-1">
                  {msg.role === 'user' ? 'You' : 'AI'}
                </p>
                <p className="text-sm whitespace-pre-wrap">{msg.content}</p>
              </div>
            </div>
          ))}
          {sending && (
            <div className="flex justify-start">
              <div className="bg-gray-800 border border-gray-700 rounded-xl px-5 py-3">
                <div className="animate-pulse text-gray-400 text-sm">Thinking...</div>
              </div>
            </div>
          )}
          <div ref={bottomRef} />
        </div>

        <footer className="border-t border-gray-800 p-4">
          <form onSubmit={handleSend} className="max-w-3xl mx-auto flex gap-3">
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Type your message..."
              disabled={sending}
              className="flex-1 px-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-vpsik-500"
            />
            <button
              type="submit"
              disabled={sending || !input.trim()}
              className="px-6 py-2.5 bg-vpsik-600 hover:bg-vpsik-500 disabled:bg-gray-700 text-white rounded-lg transition"
            >
              Send
            </button>
          </form>
        </footer>
      </div>
    </div>
  )
}
