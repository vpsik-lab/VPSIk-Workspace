"use client"

import { useEffect, type ReactNode } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/lib/auth-context"
import { Loader2 } from "lucide-react"

export default function ProtectedPage({ children }: { children: ReactNode }) {
  const { authenticated, loading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!loading && !authenticated) {
      router.push("/login")
    }
  }, [loading, authenticated, router])

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Loading...</p>
        </div>
      </div>
    )
  }

  if (!authenticated) return null

  return <>{children}</>
}
