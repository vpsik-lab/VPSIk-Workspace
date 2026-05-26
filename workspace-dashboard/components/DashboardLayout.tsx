"use client"

import type { ReactNode } from "react"
import Sidebar from "@/components/Sidebar"
import UserMenu from "@/components/UserMenu"
import ProtectedPage from "@/components/ProtectedPage"
import { FadeIn, StaggerContainer } from "@/components/motion-wrapper"
import { useAuth } from "@/lib/auth-context"
import { ArrowUpRight, Activity } from "lucide-react"
import { Badge } from "@/components/ui/badge"

interface DashboardLayoutProps {
  title: string
  subtitle?: string
  children: ReactNode
  actions?: ReactNode
  showUpgrade?: boolean
  animated?: boolean
}

export default function DashboardLayout({
  title,
  subtitle,
  children,
  actions,
  showUpgrade = true,
  animated = true,
}: DashboardLayoutProps) {
  const { user, authenticated, logout } = useAuth()

  const content = (
    <div className="min-h-screen bg-background flex">
      <Sidebar />
      <div className="flex-1 flex flex-col min-w-0">
        <FadeIn>
          <header className="h-16 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 flex items-center justify-between px-4 md:px-8 sticky top-0 z-30">
            <div className="flex items-center gap-3">
              <Activity className="h-5 w-5 text-primary hidden lg:block" />
              <h1 className="text-lg font-semibold">{title}</h1>
              {subtitle && (
                <Badge variant="secondary" className="hidden sm:inline-flex">
                  {subtitle}
                </Badge>
              )}
            </div>
            <div className="flex items-center gap-3">
              {actions}
              {showUpgrade && (
                <a
                  href="https://vpsik.com/pro"
                  target="_blank"
                  className="hidden sm:inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-500 hover:to-blue-500 text-white transition-all duration-200"
                >
                  Upgrade to Pro
                  <ArrowUpRight className="h-3 w-3" />
                </a>
              )}
              <UserMenu username={user || "admin"} onLogout={logout} />
            </div>
          </header>
        </FadeIn>

        {animated ? (
          <StaggerContainer className="flex-1 p-4 md:p-8">
            {children}
          </StaggerContainer>
        ) : (
          <main className="flex-1 p-4 md:p-8">{children}</main>
        )}
      </div>
    </div>
  )

  return <ProtectedPage>{content}</ProtectedPage>
}
