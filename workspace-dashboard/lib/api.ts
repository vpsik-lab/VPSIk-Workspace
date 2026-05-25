const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'

export interface ServiceStatus {
  name: string
  status: string
  error?: string
}

export interface StatusResponse {
  services: ServiceStatus[]
  timestamp: string
}

export interface GiteaRepo {
  name: string
  full_name: string
  description: string
  private: boolean
  language: string
  stars_count: number
  forks_count: number
  updated_at: string
}

export interface OllamaModel {
  name: string
  size: number
  digest: string
  modified_at: string
}

export interface CoolifyServer {
  id: number
  name: string
  ip: string
  status: string
}

function getToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('vpsik_token')
}

function authHeaders(): Record<string, string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  const token = getToken()
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return headers
}

export async function login(username: string, password: string): Promise<string> {
  const res = await fetch(`${API_BASE}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'login failed' }))
    throw new Error(err.error || 'login failed')
  }

  const data = await res.json()
  localStorage.setItem('vpsik_token', data.token)
  localStorage.setItem('vpsik_user', data.username)
  return data.token
}

export async function verify(): Promise<boolean> {
  const res = await fetch(`${API_BASE}/api/auth/verify`, {
    headers: authHeaders(),
  })
  return res.ok
}

export async function logout(): Promise<void> {
  localStorage.removeItem('vpsik_token')
  localStorage.removeItem('vpsik_user')
}

export function getUser(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('vpsik_user')
}

export async function getStatus(): Promise<StatusResponse> {
  const res = await fetch(`${API_BASE}/api/status`, {
    headers: authHeaders(),
  })
  if (!res.ok) throw new Error('failed to fetch status')
  return res.json()
}

export async function getGiteaRepos(): Promise<GiteaRepo[]> {
  const res = await fetch(`${API_BASE}/api/gitea/repos`, {
    headers: authHeaders(),
  })
  if (!res.ok) throw new Error('failed to fetch repos')
  return res.json()
}

export async function getOllamaModels(): Promise<OllamaModel[]> {
  const res = await fetch(`${API_BASE}/api/ollama/models`, {
    headers: authHeaders(),
  })
  if (!res.ok) throw new Error('failed to fetch models')
  return res.json()
}

export async function getCoolifyServers(): Promise<CoolifyServer[]> {
  const res = await fetch(`${API_BASE}/api/coolify/servers`, {
    headers: authHeaders(),
  })
  if (!res.ok) throw new Error('failed to fetch servers')
  return res.json()
}

export async function chatOllama(model: string, messages: { role: string; content: string }[]): Promise<string> {
  const res = await fetch(`${API_BASE}/api/ollama/chat`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ model, messages }),
  })
  if (!res.ok) throw new Error('chat failed')
  const data = await res.json()
  return data.reply
}
