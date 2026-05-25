'use client'

import { useState } from 'react'

interface UserMenuProps {
  username: string
  onLogout: () => void
}

export default function UserMenu({ username, onLogout }: UserMenuProps) {
  const [open, setOpen] = useState(false)

  return (
    <div className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-2 px-3 py-1.5 bg-gray-800 hover:bg-gray-700 rounded-lg transition text-sm"
      >
        <div className="w-6 h-6 rounded-full bg-vpsik-600 flex items-center justify-center text-xs font-medium text-white">
          {username[0].toUpperCase()}
        </div>
        <span className="text-gray-200">{username}</span>
      </button>

      {open && (
        <>
          <div className="fixed inset-0 z-10" onClick={() => setOpen(false)} />
          <div className="absolute right-0 mt-2 w-48 bg-gray-800 border border-gray-700 rounded-xl shadow-xl z-20 py-1">
            <div className="px-4 py-2 text-xs text-gray-500 border-b border-gray-700">
              Signed in as <span className="text-gray-300 font-medium">{username}</span>
            </div>
            <button
              onClick={onLogout}
              className="w-full text-left px-4 py-2.5 text-sm text-red-400 hover:bg-gray-700 transition"
            >
              Sign out
            </button>
          </div>
        </>
      )}
    </div>
  )
}
