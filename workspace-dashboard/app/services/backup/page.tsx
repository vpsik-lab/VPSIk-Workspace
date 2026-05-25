'use client'

import { useEffect, useState } from 'react'
import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ProtectedPage from '@/components/ProtectedPage'
import { ListSkeleton } from '@/components/LoadingSkeleton'
import { useAuth } from '@/lib/auth-context'
import {
  getResticSnapshots, createResticBackup, restoreResticSnapshot,
  forgetResticSnapshots, checkResticRepo, getResticStats,
  ResticSnapshot, ResticBackupStats,
} from '@/lib/api'

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

export default function BackupPage() {
  const { user, authenticated, logout } = useAuth()
  const [snapshots, setSnapshots] = useState<ResticSnapshot[]>([])
  const [loading, setLoading] = useState(true)
  const [running, setRunning] = useState<string | null>(null)
  const [backupPaths, setBackupPaths] = useState('/data')
  const [backupTags, setBackupTags] = useState('')
  const [restoreTarget, setRestoreTarget] = useState('/restore')
  const [forgetKeep, setForgetKeep] = useState(7)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  function showMsg(type: 'success' | 'error', text: string) {
    setMessage({ type, text })
    setTimeout(() => setMessage(null), 5000)
  }

  async function loadSnapshots() {
    try {
      const data = await getResticSnapshots()
      setSnapshots(data)
    } catch (err: unknown) {
      const e = err instanceof Error ? err.message : 'failed'
      showMsg('error', e)
    }
  }

  useEffect(() => {
    if (!authenticated) return
    loadSnapshots().finally(() => setLoading(false))
  }, [authenticated])

  async function handleBackup() {
    setRunning('backup')
    try {
      const paths = backupPaths.split(',').map(s => s.trim()).filter(Boolean)
      const tags = backupTags.split(',').map(s => s.trim()).filter(Boolean)
      await createResticBackup(paths, tags.length > 0 ? tags : undefined)
      showMsg('success', 'Backup completed')
      loadSnapshots()
    } catch (err: unknown) {
      const e = err instanceof Error ? err.message : 'failed'
      showMsg('error', e)
    }
    setRunning(null)
  }

  async function handleRestore(snapshotID: string) {
    setRunning(snapshotID)
    try {
      await restoreResticSnapshot(snapshotID, restoreTarget)
      showMsg('success', `Restore of ${snapshotID} started`)
    } catch (err: unknown) {
      const e = err instanceof Error ? err.message : 'failed'
      showMsg('error', e)
    }
    setRunning(null)
  }

  async function handleForget() {
    setRunning('forget')
    try {
      await forgetResticSnapshots(forgetKeep)
      showMsg('success', 'Old snapshots pruned')
      loadSnapshots()
    } catch (err: unknown) {
      const e = err instanceof Error ? err.message : 'failed'
      showMsg('error', e)
    }
    setRunning(null)
  }

  async function handleCheck() {
    setRunning('check')
    try {
      await checkResticRepo()
      showMsg('success', 'Repository check passed')
    } catch (err: unknown) {
      const e = err instanceof Error ? err.message : 'failed'
      showMsg('error', e)
    }
    setRunning(null)
  }

  return (
    <ProtectedPage>
      <div className="min-h-screen bg-gray-950 flex">
        <Sidebar />
        <div className="flex-1">
          <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
            <h1 className="text-lg font-semibold text-white">Backup & Recovery</h1>
            <UserMenu username={user || 'admin'} onLogout={logout} />
          </header>
          <main className="p-8 space-y-6">
            {message && (
              <div className={`px-4 py-3 rounded-lg text-sm ${
                message.type === 'success' ? 'bg-green-900/50 text-green-400' : 'bg-red-900/50 text-red-400'
              }`}>
                {message.text}
              </div>
            )}

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                <h2 className="text-white font-semibold mb-4">Create Backup</h2>
                <div className="space-y-3">
                  <div>
                    <label className="text-xs text-gray-500 uppercase">Paths (comma separated)</label>
                    <input
                      value={backupPaths}
                      onChange={e => setBackupPaths(e.target.value)}
                      className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-200 mt-1"
                    />
                  </div>
                  <div>
                    <label className="text-xs text-gray-500 uppercase">Tags (optional)</label>
                    <input
                      value={backupTags}
                      onChange={e => setBackupTags(e.target.value)}
                      placeholder="daily,production"
                      className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-200 mt-1"
                    />
                  </div>
                  <button
                    onClick={handleBackup}
                    disabled={running === 'backup'}
                    className="w-full px-4 py-2 bg-vpsik-600 hover:bg-vpsik-500 disabled:bg-gray-700 text-white text-sm rounded-lg transition"
                  >
                    {running === 'backup' ? 'Running...' : 'Start Backup'}
                  </button>
                </div>
              </div>

              <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                <h2 className="text-white font-semibold mb-4">Prune Snapshots</h2>
                <div className="space-y-3">
                  <div>
                    <label className="text-xs text-gray-500 uppercase">Keep Last</label>
                    <input
                      type="number"
                      value={forgetKeep}
                      onChange={e => setForgetKeep(Number(e.target.value))}
                      min={1}
                      className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-200 mt-1"
                    />
                  </div>
                  <button
                    onClick={handleForget}
                    disabled={running === 'forget'}
                    className="w-full px-4 py-2 bg-yellow-700 hover:bg-yellow-600 disabled:bg-gray-700 text-white text-sm rounded-lg transition"
                  >
                    {running === 'forget' ? 'Pruning...' : 'Prune Old Snapshots'}
                  </button>
                </div>
              </div>

              <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                <h2 className="text-white font-semibold mb-4">Maintenance</h2>
                <div className="space-y-3">
                  <button
                    onClick={handleCheck}
                    disabled={running === 'check'}
                    className="w-full px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 text-gray-200 text-sm rounded-lg transition"
                  >
                    {running === 'check' ? 'Checking...' : 'Check Repository'}
                  </button>
                  <div>
                    <label className="text-xs text-gray-500 uppercase">Restore Target</label>
                    <input
                      value={restoreTarget}
                      onChange={e => setRestoreTarget(e.target.value)}
                      className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-200 mt-1"
                    />
                  </div>
                </div>
              </div>
            </div>

            <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <h2 className="text-white font-semibold mb-4">Snapshots ({snapshots.length})</h2>
              {loading ? (
                <ListSkeleton rows={3} />
              ) : snapshots.length === 0 ? (
                <p className="text-gray-500 text-sm">No snapshots found. Create your first backup above.</p>
              ) : (
                <div className="space-y-3">
                  {snapshots.map(snap => (
                    <div key={snap.id} className="bg-gray-800/50 border border-gray-700 rounded-lg p-4 flex items-center justify-between">
                      <div>
                        <div className="flex items-center gap-2">
                          <span className="text-white font-mono text-sm">{snap.short_id}</span>
                          <span className="text-gray-400 text-xs">{new Date(snap.time).toLocaleString()}</span>
                          {snap.tags?.map(t => (
                            <span key={t} className="px-2 py-0.5 bg-vpsik-900/50 text-vpsik-400 rounded text-xs">{t}</span>
                          ))}
                        </div>
                        <div className="text-gray-500 text-xs mt-1">
                          {snap.hostname} — {snap.paths?.join(', ')}
                        </div>
                      </div>
                      <button
                        onClick={() => handleRestore(snap.short_id)}
                        disabled={running === snap.short_id}
                        className="px-3 py-1.5 bg-vpsik-600 hover:bg-vpsik-500 disabled:bg-gray-700 text-white text-xs rounded-lg transition"
                      >
                        {running === snap.short_id ? '...' : 'Restore'}
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </main>
        </div>
      </div>
    </ProtectedPage>
  )
}
