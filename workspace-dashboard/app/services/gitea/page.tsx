"use client"

import { useEffect, useState } from "react"
import DashboardLayout from "@/components/DashboardLayout"
import { getGiteaRepos, API_BASE, type GiteaRepo } from "@/lib/api"
import { ListSkeleton } from "@/components/LoadingSkeleton"
import { FadeIn, StaggerItem } from "@/components/motion-wrapper"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  GitBranch,
  Plus,
  Star,
  GitFork,
  ExternalLink,
} from "lucide-react"

interface WebhookForm {
  repo: string
  url: string
  secret: string
}

export default function GiteaPage() {
  const [repos, setRepos] = useState<GiteaRepo[]>([])
  const [loading, setLoading] = useState(true)
  const [webhookOpen, setWebhookOpen] = useState(false)
  const [webhook, setWebhook] = useState<WebhookForm>({ repo: "", url: "", secret: "" })
  const [webhookMsg, setWebhookMsg] = useState("")

  useEffect(() => {
    getGiteaRepos()
      .then(setRepos)
      .catch((err) => console.error("failed to fetch repos", err))
      .finally(() => setLoading(false))
  }, [])

  async function handleCreateWebhook(e: React.FormEvent) {
    e.preventDefault()
    setWebhookMsg("")
    try {
      const res = await fetch(`${API_BASE}/api/gitea/webhooks`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          repo: webhook.repo,
          url: webhook.url,
          secret: webhook.secret,
          events: ["push", "pull_request"],
        }),
      })
      if (!res.ok) throw new Error("webhook creation failed")
      setWebhookMsg("✅ Webhook created")
      setWebhookOpen(false)
    } catch (err: unknown) {
      setWebhookMsg(`❌ ${err instanceof Error ? err.message : String(err)}`)
    }
  }

  return (
    <DashboardLayout
      title="Gitea"
      subtitle={`${repos.length} repos`}
      actions={
        <Dialog open={webhookOpen} onOpenChange={setWebhookOpen}>
          <DialogTrigger asChild>
            <Button variant="outline" size="sm">
              <Plus className="h-4 w-4 mr-1" />
              Webhook
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create Webhook</DialogTitle>
              <DialogDescription>
                Configure a webhook to trigger deployments on push events.
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleCreateWebhook} className="space-y-4">
              <div className="space-y-2">
                <Label>Repository</Label>
                <Select
                  value={webhook.repo}
                  onValueChange={(v) => setWebhook({ ...webhook, repo: v })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select repository" />
                  </SelectTrigger>
                  <SelectContent>
                    {repos.map((r) => (
                      <SelectItem key={r.full_name} value={r.full_name}>
                        {r.full_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Webhook URL</Label>
                <Input
                  type="url"
                  value={webhook.url}
                  onChange={(e) => setWebhook({ ...webhook, url: e.target.value })}
                  placeholder="http://coolify:8000/api/v1/deploy"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label>Secret (optional)</Label>
                <Input
                  type="text"
                  value={webhook.secret}
                  onChange={(e) => setWebhook({ ...webhook, secret: e.target.value })}
                  placeholder="your-webhook-secret"
                />
              </div>
              {webhookMsg && (
                <p className="text-sm text-muted-foreground">{webhookMsg}</p>
              )}
              <DialogFooter>
                <Button type="submit">Create Webhook</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      }
    >
      {loading ? (
        <ListSkeleton rows={6} />
      ) : repos.length === 0 ? (
        <FadeIn>
          <div className="flex items-center justify-center min-h-[40vh]">
            <div className="text-center space-y-4">
              <GitBranch className="h-12 w-12 text-muted-foreground/50 mx-auto" />
              <h3 className="font-semibold">No repositories found</h3>
              <p className="text-sm text-muted-foreground">
                Create your first repository in Gitea to get started.
              </p>
            </div>
          </div>
        </FadeIn>
      ) : (
        <div className="grid gap-4">
          {repos.map((repo, i) => (
            <StaggerItem key={repo.full_name}>
              <Card className="hover:border-primary/30 transition-all duration-200">
                <CardContent className="p-5">
                  <div className="flex items-start justify-between">
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <h3 className="font-semibold">{repo.full_name}</h3>
                        {repo.private && (
                          <Badge variant="secondary" className="text-[10px] px-1.5 py-0">
                            Private
                          </Badge>
                        )}
                      </div>
                      <p className="text-sm text-muted-foreground">
                        {repo.description || "No description"}
                      </p>
                    </div>
                    <div className="flex items-center gap-3 text-sm text-muted-foreground">
                      {repo.language && (
                        <Badge variant="outline" className="text-xs">
                          {repo.language}
                        </Badge>
                      )}
                      <span className="flex items-center gap-1">
                        <Star className="h-3 w-3" /> {repo.stars_count}
                      </span>
                      <span className="flex items-center gap-1">
                        <GitFork className="h-3 w-3" /> {repo.forks_count}
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </StaggerItem>
          ))}
        </div>
      )}
    </DashboardLayout>
  )
}
