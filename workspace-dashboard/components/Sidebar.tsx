'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

const navItems = [
  { href: '/dashboard', label: 'Dashboard', icon: '◉' },
  { href: '/chat', label: 'AI Chat', icon: '💬' },
  { href: '/services/gitea', label: 'Gitea', icon: '📦' },
  { href: '/services/ollama', label: 'Ollama', icon: '🧠' },
  { href: '/services/coolify', label: 'Coolify', icon: '🚀' },
  { href: '/services/opencode', label: 'OpenCode.ai', icon: '🤖' },
  { href: '/services/openwebui', label: 'Open WebUI', icon: '🌐' },
  { href: '/services/codeserver', label: 'Code Server', icon: '⌨️' },
  { href: '/services/mattermost', label: 'Mattermost', icon: '💬' },
  { href: '/services/outline', label: 'Outline', icon: '📝' },
  { href: '/services/plane', label: 'Plane', icon: '📋' },
  { href: '/services/backup', label: 'Backup', icon: '💾' },
  { href: '/services/monitoring', label: 'Monitoring', icon: '📊' },
]

export default function Sidebar() {
  const pathname = usePathname()

  return (
    <aside className="w-60 bg-gray-900 border-r border-gray-800 flex flex-col">
      <div className="h-16 flex items-center px-6 border-b border-gray-800">
        <Link href="/dashboard" className="flex items-center gap-2">
          <span className="text-xl font-bold text-white">VPSIk</span>
          <span className="text-xs text-gray-500">Workspace</span>
        </Link>
      </div>

      <nav className="flex-1 py-4 px-3 space-y-1 overflow-y-auto">
        {navItems.map((item) => {
          const active = pathname === item.href || pathname.startsWith(item.href + '/')
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition ${
                active
                  ? 'bg-vpsik-600/20 text-vpsik-400 font-medium'
                  : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800'
              }`}
            >
              <span className="text-base">{item.icon}</span>
              {item.label}
            </Link>
          )
        })}
      </nav>

      <div className="px-3 py-3 border-t border-gray-800">
        <p className="text-[10px] text-gray-600 uppercase tracking-wider">Workspace v0.1</p>
      </div>
    </aside>
  )
}
