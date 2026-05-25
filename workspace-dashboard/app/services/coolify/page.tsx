'use client'

import { useEffect, useState } from 'react'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ProtectedPage from '@/components/ProtectedPage'
import { ListSkeleton } from '@/components/LoadingSkeleton'
import { useAuth } from '@/lib/auth-context'
import {
  getCoolifyProjects, getCoolifyEnvironments, getCoolifyApplications,
  deployCoolifyResource, restartCoolifyResource,
  CoolifyProject, CoolifyEnvironment, CoolifyApplication,
} from '@/lib/api'

export default function CoolifyPage() {
  const { user, authenticated, logout } = useAuth()
  const [projects, setProjects] = useState<CoolifyProject[]>([])
  const [selectedProject, setSelectedProject] = useState<string>('')
  const [environments, setEnvironments] = useState<CoolifyEnvironment[]>([])
  const [selectedEnv, setSelectedEnv] = useState<string>('')
  const [apps, setApps] = useState<CoolifyApplication[]>([])
  const [deploying, setDeploying] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!authenticated) return
    getCoolifyProjects()
      .then(data => {
        setProjects(data)
        if (data.length > 0) {
          setSelectedProject(data[0].uuid)
        }
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [authenticated])

  useEffect(() => {
    if (!selectedProject) return
    getCoolifyEnvironments(selectedProject)
      .then(data => {
        setEnvironments(data)
        if (data.length > 0) {
          setSelectedEnv(data[0].name)
        }
      })
      .catch(console.error)
  }, [selectedProject])

  useEffect(() => {
    if (!selectedProject || !selectedEnv) return
    getCoolifyApplications(selectedProject, selectedEnv)
      .then(setApps)
      .catch(console.error)
  }, [selectedProject, selectedEnv])

  async function handleDeploy(uuid: string) {
    setDeploying(uuid)
    try {
      await deployCoolifyResource(uuid)
    } catch (err) {
      console.error('deploy failed', err)
    }
    setDeploying(null)
  }

  async function handleRestart(uuid: string) {
    setDeploying(uuid)
    try {
      await restartCoolifyResource(uuid)
    } catch (err) {
      console.error('restart failed', err)
    }
    setDeploying(null)
  }

  function getStatusColor(status: string) {
    switch (status?.toLowerCase()) {
      case 'running': case 'ready': return 'bg-green-900/50 text-green-400'
      case 'exited': case 'stopped': return 'bg-red-900/50 text-red-400'
      default: return 'bg-yellow-900/50 text-yellow-400'
    }
  }

  return (
    <ProtectedPage>
      <div className="min-h-screen bg-gray-950 flex">
        <Sidebar />
        <div className="flex-1">
          <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
            <h1 className="text-lg font-semibold text-white">Coolify — Deployment Platform</h1>
            <UserMenu username={user || 'admin'} onLogout={logout} />
          </header>
          <main className="p-8">
            <div className="flex gap-4 mb-6">
              <select
                value={selectedProject}
                onChange={(e) => setSelectedProject(e.target.value)}
                className="bg-gray-800 border border-gray-700 text-gray-200 rounded-lg px-4 py-2 text-sm"
              >
                {projects.map(p => (
                  <option key={p.uuid} value={p.uuid}>{p.name}</option>
                ))}
              </select>
              <select
                value={selectedEnv}
                onChange={(e) => setSelectedEnv(e.target.value)}
                className="bg-gray-800 border border-gray-700 text-gray-200 rounded-lg px-4 py-2 text-sm"
              >
                {environments.map(e => (
                  <option key={e.name} value={e.name}>{e.name}</option>
                ))}
              </select>
            </div>

            {loading ? (
              <ListSkeleton rows={4} />
            ) : apps.length === 0 ? (
              <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
                <p className="text-gray-500">No applications found in this project/environment</p>
              </div>
            ) : (
              <div className="space-y-4">
                {apps.map(app => (
                  <div key={app.uuid} className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                    <div className="flex items-center justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-3">
                          <h3 className="text-white font-semibold">{app.name}</h3>
                          <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(app.status)}`}>
                            {app.status || 'unknown'}
                          </span>
                        </div>
                        <div className="flex gap-4 mt-2 text-sm text-gray-400">
                          {app.fqdn && <span>🌐 {app.fqdn}</span>}
                          {app.repository && <span>📦 {app.repository}</span>}
                          {app.git_branch && <span>🌿 {app.git_branch}</span>}
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <button
                          onClick={() => handleDeploy(app.uuid)}
                          disabled={deploying === app.uuid}
                          className="px-4 py-2 bg-vpsik-600 hover:bg-vpsik-500 disabled:bg-gray-700 text-white text-sm rounded-lg transition"
                        >
                          {deploying === app.uuid ? '...' : 'Deploy'}
                        </button>
                        <button
                          onClick={() => handleRestart(app.uuid)}
                          disabled={deploying === app.uuid}
                          className="px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 text-gray-200 text-sm rounded-lg transition"
                        >
                          Restart
                        </button>
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
