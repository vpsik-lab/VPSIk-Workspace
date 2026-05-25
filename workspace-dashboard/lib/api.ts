function getApiBase(): string {
  if (typeof window !== 'undefined') {
    return ''
  }
  return 'http://api:8081'
}

export const API_BASE = getApiBase()

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

function fetchOpts(): RequestInit {
  return {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
  }
}

export async function login(username: string, password: string): Promise<string> {
  const res = await fetch(`${API_BASE}/api/auth/login`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'login failed' }))
    throw new Error(err.error || 'login failed')
  }

  const data = await res.json()
  localStorage.setItem('vpsik_user', data.username)
  return data.token
}

export async function verify(): Promise<boolean> {
  const res = await fetch(`${API_BASE}/api/auth/verify`, fetchOpts())
  return res.ok
}

export async function logout(): Promise<void> {
  await fetch(`${API_BASE}/api/auth/logout`, { method: 'POST', credentials: 'include' }).catch(() => {})
  localStorage.removeItem('vpsik_user')
}

export function getUser(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('vpsik_user')
}

export async function getStatus(): Promise<StatusResponse> {
  const res = await fetch(`${API_BASE}/api/status`, fetchOpts())
  if (!res.ok) throw new Error('failed to fetch status')
  return res.json()
}

export async function getGiteaRepos(): Promise<GiteaRepo[]> {
  const res = await fetch(`${API_BASE}/api/gitea/repos`, fetchOpts())
  if (!res.ok) throw new Error('failed to fetch repos')
  return res.json()
}

export interface CoolifyProject {
  id: string
  uuid: string
  name: string
  description: string
}

export interface CoolifyEnvironment {
  id: string
  name: string
  created_at: string
}

export interface CoolifyApplication {
  id: string
  name: string
  uuid: string
  fqdn: string
  status: string
  git_source: string
  repository: string
  git_branch: string
  updated_at: string
}

export async function getCoolifyProjects(): Promise<CoolifyProject[]> {
  const res = await fetch(`${API_BASE}/api/coolify/projects`, fetchOpts())
  if (!res.ok) throw new Error('failed to fetch projects')
  return res.json()
}

export async function getCoolifyEnvironments(projectUUID: string): Promise<CoolifyEnvironment[]> {
  const res = await fetch(`${API_BASE}/api/coolify/projects/${projectUUID}/environments`, fetchOpts())
  if (!res.ok) throw new Error('failed to fetch environments')
  return res.json()
}

export async function getCoolifyApplications(projectUUID: string, envName: string): Promise<CoolifyApplication[]> {
  const res = await fetch(`${API_BASE}/api/coolify/projects/${projectUUID}/environments/${envName}/applications`, fetchOpts())
  if (!res.ok) throw new Error('failed to fetch applications')
  return res.json()
}

export async function deployCoolifyResource(resourceUUID: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/coolify/deploy`, {
    method: 'POST',
    ...fetchOpts(),
    body: JSON.stringify({ resource_uuid: resourceUUID }),
  })
  if (!res.ok) throw new Error('deploy failed')
}

export async function restartCoolifyResource(resourceUUID: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/coolify/restart`, {
    method: 'POST',
    ...fetchOpts(),
    body: JSON.stringify({ resource_uuid: resourceUUID }),
  })
  if (!res.ok) throw new Error('restart failed')
}

export async function getCoolifyDeploymentLogs(deploymentID: string): Promise<string> {
  const res = await fetch(`${API_BASE}/api/coolify/deployments/${deploymentID}/logs`, fetchOpts())
  if (!res.ok) throw new Error('failed to fetch logs')
  const data = await res.json()
  return data.logs
}

// ─── Ollama ──────────────────────────────────────────────────────

export async function getOllamaModels(): Promise<OllamaModel[]> {
  const res = await fetch(`${API_BASE}/api/ollama/models`, fetchOpts())
  if (!res.ok) throw new Error('failed to fetch models')
  return res.json()
}

export async function chatOllamaStream(
  model: string,
  messages: { role: string; content: string }[],
  onChunk: (content: string) => void,
  onDone: () => void,
  onError: (err: Error) => void,
): Promise<AbortController> {
  const controller = new AbortController()

  fetch(`${API_BASE}/api/ollama/chat`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', Accept: 'text/event-stream' },
    body: JSON.stringify({ model, messages, stream: true }),
    signal: controller.signal,
  }).then(async (res) => {
    if (!res.ok) {
      onError(new Error('stream chat failed'))
      return
    }

    const reader = res.body?.getReader()
    if (!reader) {
      onError(new Error('no reader available'))
      return
    }

    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const jsonStr = line.slice(6)
          if (jsonStr === '[DONE]') {
            onDone()
            return
          }
          try {
            const data = JSON.parse(jsonStr)
            if (data.content) {
              onChunk(data.content)
            }
            if (data.done) {
              onDone()
              return
            }
          } catch {
            // skip malformed lines
          }
        }
      }
    }
    onDone()
  }).catch((err) => {
    if (err.name !== 'AbortError') {
      onError(err)
    }
  })

  return controller
}

export async function pullOllamaModel(model: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/ollama/pull`, {
    method: 'POST',
    ...fetchOpts(),
    body: JSON.stringify({ model }),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'pull failed' }))
    throw new Error(err.error || 'pull failed')
  }
}

export async function deleteOllamaModel(model: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/ollama/models/${encodeURIComponent(model)}`, {
    method: 'DELETE',
    ...fetchOpts(),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'delete failed' }))
    throw new Error(err.error || 'delete failed')
  }
}

export async function runOllamaTask(model: string, task: string, content: string): Promise<string> {
  const res = await fetch(`${API_BASE}/api/ollama/task`, {
    method: 'POST',
    ...fetchOpts(),
    body: JSON.stringify({ model, task, content }),
  })
  if (!res.ok) throw new Error('task failed')
  const data = await res.json()
  return data.reply
}

// ─── OpenCode.ai ─────────────────────────────────────────────────

export async function openCodeChat(message: string, context?: string, repoPath?: string): Promise<string> {
  const res = await fetch(`${API_BASE}/api/opencode/chat`, {
    method: 'POST',
    ...fetchOpts(),
    body: JSON.stringify({ message, context, repo_path: repoPath }),
  })
  if (!res.ok) throw new Error('opencode chat failed')
  const data = await res.json()
  return data.reply
}

// ─── Restic / Backup ─────────────────────────────────────────────

export interface ResticSnapshot {
  id: string
  time: string
  hostname: string
  tags: string[]
  paths: string[]
  short_id: string
}

export interface ResticBackupStats {
  files_new: number
  files_changed: number
  files_unmodified: number
  dir_new: number
  dir_changed: number
  total_bytes: number
}

export async function getResticSnapshots(): Promise<ResticSnapshot[]> {
  const res = await fetch(`${API_BASE}/api/restic/snapshots`, fetchOpts())
  if (!res.ok) throw new Error('failed to fetch snapshots')
  return res.json()
}

export async function createResticBackup(paths: string[], tags?: string[]): Promise<ResticBackupStats> {
  const res = await fetch(`${API_BASE}/api/restic/backup`, {
    method: 'POST',
    ...fetchOpts(),
    body: JSON.stringify({ paths, tags }),
  })
  if (!res.ok) throw new Error('backup failed')
  return res.json()
}

export async function restoreResticSnapshot(snapshotID: string, target: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/restic/restore`, {
    method: 'POST',
    ...fetchOpts(),
    body: JSON.stringify({ snapshot_id: snapshotID, target }),
  })
  if (!res.ok) throw new Error('restore failed')
}

export async function forgetResticSnapshots(keepLast: number, tags?: string[]): Promise<void> {
  const res = await fetch(`${API_BASE}/api/restic/forget`, {
    method: 'POST',
    ...fetchOpts(),
    body: JSON.stringify({ keep_last: keepLast, tags }),
  })
  if (!res.ok) throw new Error('forget failed')
}

export async function checkResticRepo(): Promise<void> {
  const res = await fetch(`${API_BASE}/api/restic/check`, fetchOpts())
  if (!res.ok) throw new Error('restic check failed')
}

export async function getResticStats(): Promise<Record<string, unknown>> {
  const res = await fetch(`${API_BASE}/api/restic/stats`, fetchOpts())
  if (!res.ok) throw new Error('failed to fetch restic stats')
  return res.json()
}
