"use client"

import { useState } from "react"
import DashboardLayout from "@/components/DashboardLayout"
import { useOllamaModels, usePullOllamaModel, useDeleteOllamaModel } from "@/lib/hooks"
import { ListSkeleton } from "@/components/LoadingSkeleton"
import { FadeIn, StaggerItem } from "@/components/motion-wrapper"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
} from "@/components/ui/card"
import {
  Brain,
  Download,
  Trash2,
  HardDrive,
  Package,
} from "lucide-react"
import { toast } from "sonner"

function formatSize(bytes: number): string {
  const units = ["B", "KB", "MB", "GB"]
  let i = 0
  let size = bytes
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  return `${size.toFixed(i === 0 ? 0 : 1)} ${units[i]}`
}

export default function OllamaPage() {
  const [pullName, setPullName] = useState("")
  const { data: models, isLoading } = useOllamaModels()
  const pullMutation = usePullOllamaModel()
  const deleteMutation = useDeleteOllamaModel()

  async function handlePull(e: React.FormEvent) {
    e.preventDefault()
    if (!pullName.trim() || pullMutation.isPending) return
    try {
      await pullMutation.mutateAsync(pullName.trim())
      setPullName("")
      toast.success(`Model ${pullName.trim()} pulled successfully`)
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err)
      toast.error(msg)
    }
  }

  async function handleDelete(modelName: string) {
    try {
      await deleteMutation.mutateAsync(modelName)
      toast.success(`Model ${modelName} deleted`)
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err)
      toast.error(msg)
    }
  }

  return (
    <DashboardLayout
      title="Ollama"
      subtitle={`${models?.length ?? 0} models`}
    >
      <FadeIn>
        <form onSubmit={handlePull} className="flex gap-3 mb-8">
          <div className="relative flex-1">
            <Package className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              type="text"
              value={pullName}
              onChange={(e) => setPullName(e.target.value)}
              placeholder="Pull a model (e.g., llama3.2, mistral)..."
              disabled={pullMutation.isPending}
              className="pl-9"
            />
          </div>
          <Button type="submit" disabled={!pullName.trim() || pullMutation.isPending}>
            <Download className="h-4 w-4 mr-2" />
            {pullMutation.isPending ? "Pulling..." : "Pull"}
          </Button>
        </form>
      </FadeIn>

      {isLoading ? (
        <ListSkeleton rows={5} />
      ) : !models || models.length === 0 ? (
        <FadeIn>
          <div className="flex items-center justify-center min-h-[40vh]">
            <div className="text-center space-y-4">
              <Brain className="h-12 w-12 text-muted-foreground/50 mx-auto" />
              <h3 className="font-semibold">No models found</h3>
              <p className="text-sm text-muted-foreground">
                Pull a model above to get started with local AI.
              </p>
            </div>
          </div>
        </FadeIn>
      ) : (
        <div className="grid gap-3">
          {models.map((model) => (
            <StaggerItem key={model.name}>
              <Card className="hover:border-primary/30 transition-all duration-200">
                <CardContent className="p-5">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="w-10 h-10 rounded-lg bg-primary/5 flex items-center justify-center">
                        <Brain className="h-5 w-5 text-primary" />
                      </div>
                      <div>
                        <h3 className="font-semibold">{model.name}</h3>
                        <p className="text-xs text-muted-foreground font-mono mt-0.5">
                          {model.digest.slice(0, 19)}...
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-4">
                      <Badge variant="secondary" className="gap-1">
                        <HardDrive className="h-3 w-3" />
                        {formatSize(model.size)}
                      </Badge>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleDelete(model.name)}
                        disabled={deleteMutation.isPending}
                        className="text-destructive hover:text-destructive hover:bg-destructive/10"
                      >
                        <Trash2 className="h-4 w-4" />
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
