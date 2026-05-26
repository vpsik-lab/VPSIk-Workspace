"use client"

import { useState, useEffect } from "react"
import DashboardLayout from "@/components/DashboardLayout"
import { FadeIn } from "@/components/motion-wrapper"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  BookOpen,
  RefreshCw,
  ArrowUpRight,
  ExternalLink,
  Globe,
  Mail,
  Sparkles,
  Terminal,
} from "lucide-react"

interface UpdateInfo {
  current_version: string
  latest_version: string
  update_available: boolean
  release_url: string
}

const services = [
  ["Dashboard", "Management UI", "dashboard:3000"],
  ["API", "Backend proxy", "api:8081"],
  ["Gitea", "Git hosting", "gitea:3000"],
  ["Ollama", "Local LLM runtime", "ollama:11434"],
  ["OpenWebUI", "AI chat interface", "openwebui:8080"],
  ["Code-Server", "VS Code in browser", "codeserver:8443"],
  ["Mattermost", "Team communication", "mattermost:8065"],
  ["Outline", "Knowledge base", "outline:3000"],
  ["Plane", "Project management", "plane:8080"],
  ["Coolify", "App deployment", "coolify:3000"],
  ["Grafana", "Metrics dashboards", "grafana:3000"],
  ["Prometheus", "Metrics storage", "prometheus:9090"],
  ["PostgreSQL", "Database", "postgres:5432"],
  ["Restic", "Backup engine", "—"],
]

const cliCommands = [
  ["workspace init", "Generate workspace configuration"],
  ["workspace plan", "Scan environment and show installation plan"],
  ["workspace install", "Deploy missing services"],
  ["workspace status", "Show service health"],
  ["workspace doctor", "Check system health"],
  ["workspace backup", "Create backups"],
  ["workspace restore", "Restore from backup"],
  ["workspace upgrade", "Pull latest images and recreate"],
  ["workspace uninstall", "Remove all services"],
]

export default function DocsPage() {
  const [update, setUpdate] = useState<UpdateInfo | null>(null)
  const [checking, setChecking] = useState(false)

  useEffect(() => {
    checkUpdate()
  }, [])

  async function checkUpdate() {
    setChecking(true)
    try {
      const res = await fetch("/api/update/check")
      if (res.ok) {
        const data = await res.json()
        setUpdate(data)
      }
    } catch {
      // ignore
    } finally {
      setChecking(false)
    }
  }

  return (
    <DashboardLayout title="Documentation" showUpgrade>
      <div className="max-w-4xl mx-auto space-y-8">
        {/* Version + Update */}
        <FadeIn>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-muted-foreground">
                WorkSpace OS &mdash; Open Source Edition
              </p>
            </div>
            <div className="flex items-center gap-3">
              {update && (
                <Badge
                  variant={update.update_available ? "warning" : "success"}
                  className="gap-1.5"
                >
                  {update.update_available
                    ? `Update: ${update.latest_version} available`
                    : `v${update.current_version} — up to date`}
                </Badge>
              )}
              <Button variant="outline" size="sm" onClick={checkUpdate} disabled={checking}>
                <RefreshCw className={`h-3 w-3 mr-1 ${checking ? "animate-spin" : ""}`} />
                {checking ? "Checking..." : "Check"}
              </Button>
            </div>
          </div>
        </FadeIn>

        {/* What is */}
        <FadeIn delay={0.05}>
          <Card>
            <CardHeader>
              <CardTitle>What is WorkSpace OS?</CardTitle>
              <CardDescription>
                AI-native engineering workspace that deploys on any VPS with a single command.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground leading-relaxed">
                WorkSpace OS bundles Git hosting, local AI, cloud IDE, team communication,
                project management, monitoring, and backup into one isolated Docker network.
              </p>
            </CardContent>
          </Card>
        </FadeIn>

        {/* Architecture */}
        <FadeIn delay={0.1}>
          <Card>
            <CardHeader>
              <CardTitle>Architecture</CardTitle>
              <CardDescription>How the services are organized</CardDescription>
            </CardHeader>
            <CardContent>
              <pre className="text-xs text-muted-foreground font-mono leading-relaxed whitespace-pre bg-muted/30 p-4 rounded-lg">
{`Internet
   ↓
Cloudflare Tunnel (optional)
   ↓
Traefik (reverse proxy, SSL)
   ↓
workspace_net (isolated network)
   ↓
┌─────────────────────────────────────┐
│  Dashboard │ API │ CLI              │
├─────────────────────────────────────┤
│  Gitea │ Ollama │ OpenWebUI         │
│  Code-Server │ Mattermost           │
│  Outline │ Plane │ Coolify          │
├─────────────────────────────────────┤
│  PostgreSQL │ Redis │ Restic        │
│  Grafana │ Prometheus               │
└─────────────────────────────────────┘`}
              </pre>
            </CardContent>
          </Card>
        </FadeIn>

        {/* Services Table */}
        <FadeIn delay={0.15}>
          <Card>
            <CardHeader>
              <CardTitle>Services</CardTitle>
              <CardDescription>All services deployed by WorkSpace OS</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="text-muted-foreground border-b">
                      <th className="text-left py-2 pr-4 font-medium">Service</th>
                      <th className="text-left py-2 pr-4 font-medium">Description</th>
                      <th className="text-left py-2 font-medium">URL</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    {services.map(([name, desc, url]) => (
                      <tr key={name} className="border-b border-muted/50">
                        <td className="py-2 pr-4 font-medium text-foreground">{name}</td>
                        <td className="py-2 pr-4">{desc}</td>
                        <td className="py-2 font-mono text-xs">{url}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </FadeIn>

        {/* CLI Reference */}
        <FadeIn delay={0.2}>
          <Card>
            <CardHeader>
              <CardTitle>CLI Reference</CardTitle>
              <CardDescription>Available commands for the workspace CLI</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {cliCommands.map(([cmd, desc]) => (
                  <div key={cmd} className="flex gap-4 text-sm">
                    <code className="text-primary font-mono whitespace-nowrap">{cmd}</code>
                    <span className="text-muted-foreground">{desc}</span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </FadeIn>

        {/* Editions */}
        <FadeIn delay={0.25}>
          <Card>
            <CardHeader>
              <CardTitle>Editions</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <Card className="border-muted">
                  <CardContent className="p-4 space-y-2">
                    <div className="flex items-center gap-2">
                      <Sparkles className="h-4 w-4 text-primary" />
                      <h3 className="font-semibold">Open Source</h3>
                    </div>
                    <p className="text-xs text-muted-foreground">
                      Single user, self-hosted, free forever.
                    </p>
                  </CardContent>
                </Card>
                <Card className="border-purple-700/50 bg-purple-900/5">
                  <CardContent className="p-4 space-y-2">
                    <h3 className="font-semibold text-purple-400">Pro</h3>
                    <p className="text-xs text-muted-foreground">
                      Multi-user, AI CTO, team roles, notifications.
                    </p>
                    <a href="https://vpsik.com/pro" target="_blank" className="text-purple-400 text-xs hover:underline inline-flex items-center gap-1">
                      Learn more <ArrowUpRight className="h-3 w-3" />
                    </a>
                  </CardContent>
                </Card>
                <Card className="border-blue-700/50 bg-blue-900/5">
                  <CardContent className="p-4 space-y-2">
                    <h3 className="font-semibold text-blue-400">Cloud</h3>
                    <p className="text-xs text-muted-foreground">
                      Unlimited users, hosted, GitHub integration.
                    </p>
                    <a href="https://vpsik.com/cloud" target="_blank" className="text-blue-400 text-xs hover:underline inline-flex items-center gap-1">
                      Learn more <ArrowUpRight className="h-3 w-3" />
                    </a>
                  </CardContent>
                </Card>
              </div>
            </CardContent>
          </Card>
        </FadeIn>

        {/* Resources */}
        <FadeIn delay={0.3}>
          <Card>
            <CardHeader>
              <CardTitle>Resources</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3 text-sm">
                <a href="https://github.com/vpsik-lab/VPSIk-Workspace" target="_blank" className="flex items-center gap-2 text-primary hover:underline">
                  <ExternalLink className="h-4 w-4" /> GitHub Repository
                </a>
                <a href="https://vpsik.com" target="_blank" className="flex items-center gap-2 text-primary hover:underline">
                  <Globe className="h-4 w-4" /> VPSIk Website
                </a>
                <a href="https://github.com/vpsik-lab/VPSIk-Workspace/blob/main/CONTRIBUTING.md" target="_blank" className="flex items-center gap-2 text-primary hover:underline">
                  <ExternalLink className="h-4 w-4" /> Contributing Guide
                </a>
              </div>
            </CardContent>
          </Card>
        </FadeIn>

        {/* Footer */}
        <FadeIn delay={0.35}>
          <Card>
            <CardContent className="p-6 text-center space-y-2">
              <p className="text-sm text-muted-foreground">
                Built by <span className="text-foreground">Youssef Soliman</span>
              </p>
              <div className="flex items-center justify-center gap-4 text-xs text-muted-foreground">
                <a href="https://vpsik.com" className="flex items-center gap-1 hover:text-foreground">
                  <Globe className="h-3 w-3" /> vpsik.com
                </a>
                <a href="mailto:opensource@vpsik.com" className="flex items-center gap-1 hover:text-foreground">
                  <Mail className="h-3 w-3" /> opensource@vpsik.com
                </a>
                <a href="https://github.com/Ymasoli" className="flex items-center gap-1 hover:text-foreground">
                  <ExternalLink className="h-3 w-3" /> Ymasoli
                </a>
              </div>
            </CardContent>
          </Card>
        </FadeIn>
      </div>
    </DashboardLayout>
  )
}
