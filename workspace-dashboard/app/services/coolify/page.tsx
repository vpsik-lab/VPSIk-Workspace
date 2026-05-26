"use client"

import { useState } from "react"
import DashboardLayout from "@/components/DashboardLayout"
import { ListSkeleton } from "@/components/LoadingSkeleton"
import { FadeIn, StaggerItem } from "@/components/motion-wrapper"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
} from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Rocket,
  Play,
  RotateCcw,
  Globe,
  GitBranch,
  Package,
  FolderKanban,
  Layers,
} from "lucide-react"
import {
  useCoolifyProjects,
  useCoolifyEnvironments,
  useCoolifyApplications,
  useDeployCoolify,
  useRestartCoolify,
} from "@/lib/hooks"
import { toast } from "sonner"

export default function CoolifyPage() {
  const [selectedProject, setSelectedProject] = useState("")
  const [selectedEnv, setSelectedEnv] = useState("")
  const { data: projects, isLoading: projectsLoading } = useCoolifyProjects()
  const { data: environments } = useCoolifyEnvironments(selectedProject)
  const { data: apps } = useCoolifyApplications(selectedProject, selectedEnv)
  const deployMut = useDeployCoolify()
  const restartMut = useRestartCoolify()

  function getStatusColor(status: string) {
    switch (status?.toLowerCase()) {
      case "running":
      case "ready":
        return "success" as const
      case "exited":
      case "stopped":
        return "destructive" as const
      default:
        return "warning" as const
    }
  }

  return (
    <DashboardLayout
      title="Coolify"
      subtitle={`${projects?.length ?? 0} projects`}
    >
      <FadeIn>
        <div className="flex gap-4 mb-6">
          <div className="flex items-center gap-2">
            <FolderKanban className="h-4 w-4 text-muted-foreground" />
            <Select value={selectedProject} onValueChange={setSelectedProject}>
              <SelectTrigger className="w-48">
                <SelectValue placeholder="Project" />
              </SelectTrigger>
              <SelectContent>
                {(projects ?? []).map((p) => (
                  <SelectItem key={p.uuid} value={p.uuid}>{p.name}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="flex items-center gap-2">
            <Layers className="h-4 w-4 text-muted-foreground" />
            <Select value={selectedEnv} onValueChange={setSelectedEnv}>
              <SelectTrigger className="w-40">
                <SelectValue placeholder="Environment" />
              </SelectTrigger>
              <SelectContent>
                {(environments ?? []).map((e) => (
                  <SelectItem key={e.name} value={e.name}>{e.name}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
      </FadeIn>

      {projectsLoading ? (
        <ListSkeleton rows={4} />
      ) : !apps || apps.length === 0 ? (
        <FadeIn>
          <div className="flex items-center justify-center min-h-[40vh]">
            <div className="text-center space-y-4">
              <Rocket className="h-12 w-12 text-muted-foreground/50 mx-auto" />
              <h3 className="font-semibold">No applications</h3>
              <p className="text-sm text-muted-foreground">
                No applications found in this project/environment.
              </p>
            </div>
          </div>
        </FadeIn>
      ) : (
        <div className="space-y-4">
          {apps.map((app) => (
            <StaggerItem key={app.uuid}>
              <Card className="hover:border-primary/30 transition-all duration-200">
                <CardContent className="p-5">
                  <div className="flex items-center justify-between">
                    <div className="flex-1 space-y-2">
                      <div className="flex items-center gap-3">
                        <h3 className="font-semibold">{app.name}</h3>
                        <Badge variant={getStatusColor(app.status)}>{app.status || "unknown"}</Badge>
                      </div>
                      <div className="flex gap-4 text-sm text-muted-foreground">
                        {app.fqdn && <span className="flex items-center gap-1"><Globe className="h-3 w-3" /> {app.fqdn}</span>}
                        {app.repository && <span className="flex items-center gap-1"><Package className="h-3 w-3" /> {app.repository}</span>}
                        {app.git_branch && <span className="flex items-center gap-1"><GitBranch className="h-3 w-3" /> {app.git_branch}</span>}
                      </div>
                    </div>
                    <div className="flex gap-2 ml-4">
                      <Button size="sm" onClick={() => deployMut.mutate(app.uuid, {
                        onError: (err) => toast.error(err.message),
                      })} disabled={deployMut.isPending}>
                        <Play className="h-3 w-3 mr-1" /> Deploy
                      </Button>
                      <Button size="sm" variant="outline" onClick={() => restartMut.mutate(app.uuid, {
                        onError: (err) => toast.error(err.message),
                      })} disabled={restartMut.isPending}>
                        <RotateCcw className="h-3 w-3 mr-1" /> Restart
                      </Button>
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
