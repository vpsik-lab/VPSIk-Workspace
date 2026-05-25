'use client'

import { useEffect, useState, useRef, useCallback } from 'react'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ProtectedPage from '@/components/ProtectedPage'
import { useAuth } from '@/lib/auth-context'
import { getOllamaModels, chatOllamaStream, runOllamaTask, OllamaModel } from '@/lib/api'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

const TASK_TEMPLATES = [
  { id: 'explain', label: 'Explain', icon: '💡', prompt: 'Explain this in simple terms:\n\n' },
  { id: 'summarize', label: 'Summarize', icon: '📝', prompt: 'Summarize this concisely:\n\n' },
  { id: 'review', label: 'Review Code', icon: '🔍', prompt: 'Review this code for bugs and improvements:\n\n' },
  { id: 'expand', label: 'Expand', icon: '📖', prompt: 'Expand on this with more detail:\n\n' },
]

export default function ChatPage() {
  const { user, authenticated, logout } = useAuth()
  const [models, setModels] = useState<OllamaModel[]>([])
  const [selectedModel, setSelectedModel] = useState('')
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [sending, setSending] = useState(false)
  const [streamingContent, setStreamingContent] = useState('')
  const [showTasks, setShowTasks] = useState(false)
  const abortRef = useRef<AbortController | null>(null)
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!authenticated) return
    getOllamaModels()
      .then(data => {
        setModels(data)
        if (data.length > 0) setSelectedModel(data[0].name)
      })
      .catch(err => console.error('failed to fetch models', err))
  }, [authenticated])

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages, streamingContent])

  useEffect(() => {
    return () => {
      if (abortRef.current) {
        abortRef.current.abort()
        abortRef.current = null
      }
    }
  }, [])

  const handleSend = useCallback(async (content: string) => {
    if (!content.trim() || !selectedModel || sending) return

    const userMsg: Message = { role: 'user', content: content.trim() }
    setMessages(prev => [...prev, userMsg])
    setInput('')
    setSending(true)
    setStreamingContent('')
    setShowTasks(false)

    const apiMessages = [...messages, userMsg].map(m => ({
      role: m.role,
      content: m.content,
    }))

    let fullReply = ''
    abortRef.current = await chatOllamaStream(
      selectedModel,
      apiMessages,
      (chunk) => {
        fullReply += chunk
        setStreamingContent(fullReply)
      },
      () => {
        setMessages(prev => [...prev, { role: 'assistant', content: fullReply }])
        setStreamingContent('')
        setSending(false)
      },
      (err) => {
        setMessages(prev => [...prev, { role: 'assistant', content: `Error: ${err.message}` }])
        setStreamingContent('')
        setSending(false)
      },
    )
  }, [selectedModel, messages, sending])

  async function handleTask(taskId: string) {
    const task = TASK_TEMPLATES.find(t => t.id === taskId)
    if (!task || !input.trim() || !selectedModel || sending) return

    const userMsg: Message = { role: 'user', content: `${task.icon} ${task.label}: ${input.trim()}` }
    setMessages(prev => [...prev, userMsg])
    setInput('')
    setSending(true)
    setStreamingContent('')
    setShowTasks(false)

    try {
      const reply = await runOllamaTask(selectedModel, taskId, input.trim())
      setMessages(prev => [...prev, { role: 'assistant', content: reply }])
    } catch (err: unknown) {
      setMessages(prev => [...prev, { role: 'assistant', content: `Error: ${err instanceof Error ? err.message : String(err)}` }])
    }
    setSending(false)
  }

  function stopStreaming() {
    if (abortRef.current) {
      abortRef.current.abort()
      abortRef.current = null
    }
  }

  return (
    <ProtectedPage>
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
            <UserMenu username={user || 'admin'} onLogout={logout} />
          </header>

          <div className="flex-1 overflow-y-auto p-8 space-y-4">
            {messages.length === 0 && !streamingContent && (
              <div className="text-center text-gray-500 mt-20">
                <p className="text-4xl mb-4">💬</p>
                <p>Start a conversation with your local AI</p>
                <div className="flex gap-2 justify-center mt-4">
                  {TASK_TEMPLATES.map(task => (
                    <span key={task.id} className="px-3 py-1 bg-gray-800 rounded-full text-xs text-gray-400">
                      {task.icon} {task.label}
                    </span>
                  ))}
                </div>
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
            {streamingContent && (
              <div className="flex justify-start">
                <div className="bg-gray-800 border border-gray-700 rounded-xl px-5 py-3 max-w-xl">
                  <p className="text-xs text-gray-500 mb-1">AI</p>
                  <p className="text-sm whitespace-pre-wrap">{streamingContent}</p>
                </div>
              </div>
            )}
            <div ref={bottomRef} />
          </div>

          <footer className="border-t border-gray-800 p-4">
            <form onSubmit={(e) => { e.preventDefault(); handleSend(input) }} className="max-w-3xl mx-auto">
              {showTasks && input.trim() && !sending && (
                <div className="flex gap-2 mb-2">
                  {TASK_TEMPLATES.map(task => (
                    <button
                      key={task.id}
                      type="button"
                      onClick={() => handleTask(task.id)}
                      className="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 text-xs rounded-lg transition"
                    >
                      {task.icon} {task.label}
                    </button>
                  ))}
                </div>
              )}
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => setShowTasks(!showTasks)}
                  className="px-3 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-400 rounded-lg transition text-sm"
                  title="AI Tasks"
                >
                  🎯
                </button>
                <input
                  type="text"
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  placeholder="Type your message..."
                  disabled={sending}
                  className="flex-1 px-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-vpsik-500"
                />
                {sending ? (
                  <button
                    type="button"
                    onClick={stopStreaming}
                    className="px-6 py-2.5 bg-red-600 hover:bg-red-500 text-white rounded-lg transition"
                  >
                    Stop
                  </button>
                ) : (
                  <button
                    type="submit"
                    disabled={!input.trim()}
                    className="px-6 py-2.5 bg-vpsik-600 hover:bg-vpsik-500 disabled:bg-gray-700 text-white rounded-lg transition"
                  >
                    Send
                  </button>
                )}
              </div>
            </form>
          </footer>
        </div>
      </div>
    </ProtectedPage>
  )
}
