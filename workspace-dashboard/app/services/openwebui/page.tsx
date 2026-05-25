'use client'

import Sidebar from '@/components/Sidebar'
import UserMenu from '@/components/UserMenu'
import ProtectedPage from '@/components/ProtectedPage'
import { useAuth } from '@/lib/auth-context'

export default function OpenWebUIPage() {
  const { user, logout } = useAuth()
  const openWebUIUrl = process.env.NEXT_PUBLIC_OPENWEBUI_URL || 'http://localhost:3001'

  return (
    <ProtectedPage>
      <div className="min-h-screen bg-gray-950 flex">
        <Sidebar />
        <div className="flex-1 flex flex-col">
          <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-8">
            <h1 className="text-lg font-semibold text-white">Open WebUI</h1>
            <UserMenu username={user || 'admin'} onLogout={logout} />
          </header>
          <main className="flex-1 p-4">
            <iframe
              src={openWebUIUrl}
              className="w-full h-full rounded-xl border border-gray-800 bg-white"
              title="Open WebUI"
              sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
            />
          </main>
        </div>
      </div>
    </ProtectedPage>
  )
}
