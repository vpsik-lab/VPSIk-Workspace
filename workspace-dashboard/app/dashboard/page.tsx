"use client"

import { motion } from "framer-motion"
import Sidebar from "@/components/Sidebar"
import UserMenu from "@/components/UserMenu"
import ServiceCard from "@/components/ServiceCard"
import StatusBadge from "@/components/StatusBadge"
import ProtectedPage from "@/components/ProtectedPage"
import { useAuth } from "@/lib/auth-context"
import { useServices } from "@/lib/hooks"
import { FadeIn, StaggerItem } from "@/components/motion-wrapper"
import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Activity, ArrowUpRight, Boxes, Cpu, HardDrive, Shield } from "lucide-react"
import { PageSkeleton } from "@/components/LoadingSkeleton"

const SERVICE_CARDS = [
  { name: "Gitea", href: "/services/gitea", description: "Git repositories & issues", icon: "📦", apiName: "gitea" },
  { name: "Ollama", href: "/services/ollama", description: "AI models & management", icon: "🧠", apiName: "ollama" },
  { name: "Coolify", href: "/services/coolify", description: "App deployments & servers", icon: "🚀", apiName: "coolify" },
  { name: "AI Chat", href: "/chat", description: "Chat with local models", icon: "💬", apiName: "" },
  { name: "OpenCode.ai", href: "/services/opencode", description: "AI coding assistant", icon: "🤖", apiName: "opencode" },
  { name: "Open WebUI", href: "/services/openwebui", description: "Full AI chat interface", icon: "🌐", apiName: "" },
]

const statCards = [
  { label: "Services", icon: Boxes, color: "text-primary" },
  { label: "CPU", icon: Cpu, color: "text-emerald-500" },
  { label: "Memory", icon: HardDrive, color: "text-amber-500" },
  { label: "Uptime", icon: Shield, color: "text-sky-500" },
]

const containerVariants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1, transition: { staggerChildren: 0.06 } },
}

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
}

export default function DashboardPage() {
  const { user, authenticated, logout } = useAuth()
  const { data, isLoading } = useServices()
  const services = data?.services ?? []

  if (!authenticated || isLoading) {
    return <PageSkeleton />
  }

  const healthyCount = services.filter((s) => s.status === "healthy").length

  function getStatus(name: string): "healthy" | "unhealthy" | "unknown" {
    const svc = services.find((s) => s.name === name)
    if (!svc) return "unknown"
    return svc.status === "healthy" ? "healthy" : "unhealthy"
  }

  return (
    <ProtectedPage>
      <div className="min-h-screen bg-background flex">
        <Sidebar />
        <div className="flex-1 flex flex-col min-w-0">
          <FadeIn>
            <header className="h-16 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 flex items-center justify-between px-4 md:px-8 sticky top-0 z-30">
              <div className="flex items-center gap-4">
                <div className="hidden lg:flex items-center gap-2">
                  <Activity className="h-5 w-5 text-primary" />
                  <h1 className="text-lg font-semibold">Dashboard</h1>
                </div>
                <Badge variant="secondary" className="gap-1.5">
                  <span className={`w-1.5 h-1.5 rounded-full ${
                    healthyCount === services.length
                      ? "bg-emerald-400"
                      : healthyCount > 0
                        ? "bg-amber-400"
                        : "bg-red-400"
                  }`} />
                  {healthyCount}/{services.length} healthy
                </Badge>
              </div>
              <div className="flex items-center gap-3">
                <a href="https://vpsik.com/pro" target="_blank"
                  className="hidden sm:inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-500 hover:to-blue-500 text-white transition-all duration-200"
                >
                  Upgrade to Pro
                  <ArrowUpRight className="h-3 w-3" />
                </a>
                <UserMenu username={user || "admin"} onLogout={logout} />
              </div>
            </header>
          </FadeIn>

          <main className="flex-1 p-4 md:p-8 space-y-8">
            <motion.div variants={containerVariants} initial="hidden" animate="visible"
              className="grid grid-cols-2 md:grid-cols-4 gap-4"
            >
              {statCards.map((stat, i) => {
                const Icon = stat.icon
                return (
                  <motion.div key={stat.label} variants={itemVariants}>
                    <Card>
                      <CardContent className="p-4 md:p-6 flex items-center gap-4">
                        <div className="p-2 rounded-lg bg-primary/5">
                          <Icon className={`h-5 w-5 ${stat.color}`} />
                        </div>
                        <div>
                          <p className="text-2xl font-bold">
                            {i === 0 ? services.length : "—"}
                          </p>
                          <p className="text-xs text-muted-foreground">{stat.label}</p>
                        </div>
                      </CardContent>
                    </Card>
                  </motion.div>
                )
              })}
            </motion.div>

            <div>
              <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.2 }}>
                <h2 className="text-lg font-semibold mb-4">Services</h2>
              </motion.div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {SERVICE_CARDS.map((card, i) => (
                  <ServiceCard
                    key={card.href}
                    name={card.name}
                    href={card.href}
                    description={card.description}
                    status={card.apiName ? getStatus(card.apiName) : "healthy"}
                    icon={card.icon}
                    index={i}
                  />
                ))}
              </div>
            </div>

            <FadeIn delay={0.4}>
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Activity className="h-5 w-5 text-primary" />
                    Service Status
                  </CardTitle>
                  <CardDescription>Real-time health status of all workspace services</CardDescription>
                </CardHeader>
                <CardContent>
                  {services.length === 0 ? (
                    <div className="text-center py-8 text-muted-foreground">
                      <Activity className="h-8 w-8 mx-auto mb-3 opacity-50" />
                      <p>No services detected</p>
                    </div>
                  ) : (
                    <div className="divide-y">
                      {services.map((svc, i) => (
                        <motion.div key={svc.name}
                          initial={{ opacity: 0, x: -10 }} animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: i * 0.03 }}
                          className="flex items-center justify-between py-3 first:pt-0 last:pb-0"
                        >
                          <div className="flex items-center gap-3">
                            <span className="text-sm font-medium capitalize">{svc.name}</span>
                            {svc.error && (
                              <span className="text-xs text-muted-foreground truncate max-w-[200px]">{svc.error}</span>
                            )}
                          </div>
                          <StatusBadge status={svc.status} />
                        </motion.div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </FadeIn>
          </main>
        </div>
      </div>
    </ProtectedPage>
  )
}
