"use client"

import { useState } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import {
  LayoutDashboard,
  MessageSquare,
  GitBranch,
  Brain,
  Rocket,
  Code2,
  Globe,
  Terminal,
  Users,
  FileText,
  ClipboardList,
  HardDrive,
  BarChart3,
  BookOpen,
  ChevronLeft,
  Menu,
  X,
} from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { motion, AnimatePresence } from "framer-motion"

const navItems = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/chat", label: "AI Chat", icon: MessageSquare },
  { href: "/services/gitea", label: "Gitea", icon: GitBranch },
  { href: "/services/ollama", label: "Ollama", icon: Brain },
  { href: "/services/coolify", label: "Coolify", icon: Rocket },
  { href: "/services/opencode", label: "OpenCode.ai", icon: Code2 },
  { href: "/services/openwebui", label: "Open WebUI", icon: Globe },
  { href: "/services/codeserver", label: "Code Server", icon: Terminal },
  { href: "/services/mattermost", label: "Mattermost", icon: Users },
  { href: "/services/outline", label: "Outline", icon: FileText },
  { href: "/services/plane", label: "Plane", icon: ClipboardList },
  { href: "/services/backup", label: "Backup", icon: HardDrive },
  { href: "/services/monitoring", label: "Monitoring", icon: BarChart3 },
  { href: "/docs", label: "Docs", icon: BookOpen },
]

const sidebarVariants = {
  open: { x: 0, transition: { type: "spring" as const, stiffness: 300, damping: 30 } },
  closed: { x: "-100%", transition: { type: "spring" as const, stiffness: 300, damping: 30 } },
}

export default function Sidebar() {
  const pathname = usePathname()
  const [collapsed, setCollapsed] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)

  return (
    <>
      {/* Mobile hamburger */}
      <button
        onClick={() => setMobileOpen(true)}
        className="fixed top-4 left-4 z-50 lg:hidden flex items-center justify-center w-10 h-10 rounded-lg bg-background border shadow-sm"
        aria-label="Open menu"
      >
        <Menu className="h-5 w-5" />
      </button>

      {/* Mobile overlay */}
      <AnimatePresence>
        {mobileOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-40 bg-black/50 lg:hidden"
            onClick={() => setMobileOpen(false)}
          />
        )}
      </AnimatePresence>

      {/* Mobile sidebar */}
      <AnimatePresence>
        {mobileOpen && (
          <motion.aside
            initial="closed"
            animate="open"
            exit="closed"
            variants={sidebarVariants}
            className="fixed inset-y-0 left-0 z-50 w-64 bg-background border-r shadow-2xl lg:hidden flex flex-col"
          >
            <SidebarContent
              pathname={pathname}
              onNavigate={() => setMobileOpen(false)}
              mobile
            />
          </motion.aside>
        )}
      </AnimatePresence>

      {/* Desktop sidebar */}
      <aside
        className={cn(
          "hidden lg:flex flex-col bg-background border-r transition-all duration-300",
          collapsed ? "w-16" : "w-60"
        )}
      >
        <div className="h-16 flex items-center px-4 border-b">
          <Link href="/dashboard" className="flex items-center gap-2 overflow-hidden">
            <div className="flex-shrink-0 w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
              <span className="text-xs font-bold text-primary-foreground">W</span>
            </div>
            <span
              className={cn(
                "font-semibold text-foreground transition-opacity duration-300",
                collapsed && "opacity-0 w-0"
              )}
            >
              WorkSpace OS
            </span>
          </Link>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setCollapsed(!collapsed)}
            className="ml-auto flex-shrink-0 h-8 w-8"
          >
            <ChevronLeft className={cn("h-4 w-4 transition-transform", collapsed && "rotate-180")} />
          </Button>
        </div>

        <nav className="flex-1 py-4 px-2 space-y-1 overflow-y-auto">
          {navItems.map((item) => {
            const active =
              pathname === item.href || pathname.startsWith(item.href + "/")
            const Icon = item.icon
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-200 group relative",
                  active
                    ? "bg-primary/10 text-primary font-medium"
                    : "text-muted-foreground hover:text-foreground hover:bg-accent"
                )}
              >
                <Icon className="h-4 w-4 flex-shrink-0" />
                <span
                  className={cn(
                    "transition-opacity duration-300",
                    collapsed && "opacity-0 w-0 overflow-hidden"
                  )}
                >
                  {item.label}
                </span>
                {active && (
                  <span className="absolute left-0 top-1/2 -translate-y-1/2 w-0.5 h-5 bg-primary rounded-full" />
                )}
              </Link>
            )
          })}
        </nav>

        <div className="px-3 py-3 border-t">
          <p
            className={cn(
              "text-[10px] text-muted-foreground uppercase tracking-wider transition-opacity",
              collapsed && "opacity-0"
            )}
          >
            WorkSpace OS v0.1
          </p>
        </div>
      </aside>
    </>
  )
}

function SidebarContent({
  pathname,
  onNavigate,
  mobile,
}: {
  pathname: string
  onNavigate: () => void
  mobile?: boolean
}) {
  return (
    <>
      <div className="h-16 flex items-center justify-between px-4 border-b">
        <Link
          href="/dashboard"
          className="flex items-center gap-2"
          onClick={onNavigate}
        >
          <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
            <span className="text-xs font-bold text-primary-foreground">W</span>
          </div>
          <span className="font-semibold text-foreground">WorkSpace OS</span>
        </Link>
        <Button variant="ghost" size="icon" onClick={onNavigate} className="h-8 w-8">
          <X className="h-4 w-4" />
        </Button>
      </div>

      <nav className="flex-1 py-4 px-3 space-y-1 overflow-y-auto">
        {navItems.map((item) => {
          const active = pathname === item.href || pathname.startsWith(item.href + "/")
          const Icon = item.icon
          return (
            <Link
              key={item.href}
              href={item.href}
              onClick={onNavigate}
              className={cn(
                "flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-200",
                active
                  ? "bg-primary/10 text-primary font-medium"
                  : "text-muted-foreground hover:text-foreground hover:bg-accent"
              )}
            >
              <Icon className="h-4 w-4 flex-shrink-0" />
              <span>{item.label}</span>
            </Link>
          )
        })}
      </nav>

      {mobile && (
        <div className="px-3 py-3 border-t">
          <p className="text-[10px] text-muted-foreground uppercase tracking-wider">
            WorkSpace OS v0.1
          </p>
        </div>
      )}
    </>
  )
}
