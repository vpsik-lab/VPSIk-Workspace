"use client"

import { useState } from "react"
import DashboardLayout from "@/components/DashboardLayout"
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
  HardDrive,
  Shield,
  Trash2,
  RefreshCw,
  RotateCcw,
  Clock,
  Tag,
} from "lucide-react"
import { toast } from "sonner"
import {
  useResticSnapshots,
  useCreateBackup,
  useRestoreSnapshot,
  useForgetSnapshots,
  useCheckRestic,
} from "@/lib/hooks"

export default function BackupPage() {
  const [backupPaths, setBackupPaths] = useState("/data")
  const [backupTags, setBackupTags] = useState("")
  const [restoreTarget, setRestoreTarget] = useState("/restore")
  const [forgetKeep, setForgetKeep] = useState(7)

  const { data: snapshots, isLoading } = useResticSnapshots()
  const backupMut = useCreateBackup()
  const restoreMut = useRestoreSnapshot()
  const forgetMut = useForgetSnapshots()
  const checkMut = useCheckRestic()

  const running = backupMut.isPending || forgetMut.isPending || checkMut.isPending

  async function handleBackup() {
    try {
      const paths = backupPaths.split(",").map((s) => s.trim()).filter(Boolean)
      const tags = backupTags.split(",").map((s) => s.trim()).filter(Boolean)
      await backupMut.mutateAsync({ paths, tags: tags.length > 0 ? tags : undefined })
      toast.success("Backup completed successfully")
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Backup failed")
    }
  }

  async function handleRestore(snapshotID: string) {
    try {
      await restoreMut.mutateAsync({ snapshotID, target: restoreTarget })
      toast.success(`Restore of ${snapshotID} started`)
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Restore failed")
    }
  }

  async function handleForget() {
    try {
      await forgetMut.mutateAsync({ keepLast: forgetKeep })
      toast.success("Old snapshots pruned")
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Prune failed")
    }
  }

  async function handleCheck() {
    try {
      await checkMut.mutateAsync()
      toast.success("Repository check passed")
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Check failed")
    }
  }

  return (
    <DashboardLayout title="Backup & Recovery" subtitle={`${snapshots?.length ?? 0} snapshots`}>
      <FadeIn>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <HardDrive className="h-4 w-4 text-primary" />
                Create Backup
              </CardTitle>
              <CardDescription>Back up your workspace data</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="space-y-2">
                <Label className="text-xs">Paths (comma separated)</Label>
                <Input value={backupPaths} onChange={(e) => setBackupPaths(e.target.value)} size={1} />
              </div>
              <div className="space-y-2">
                <Label className="text-xs">Tags (optional)</Label>
                <Input value={backupTags} onChange={(e) => setBackupTags(e.target.value)} placeholder="daily,production" />
              </div>
              <Button onClick={handleBackup} disabled={backupMut.isPending} className="w-full">
                <HardDrive className="h-4 w-4 mr-2" />
                {backupMut.isPending ? "Running..." : "Start Backup"}
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <Trash2 className="h-4 w-4 text-amber-500" />
                Prune Snapshots
              </CardTitle>
              <CardDescription>Remove old backups to save space</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="space-y-2">
                <Label className="text-xs">Keep Last</Label>
                <Input type="number" value={forgetKeep} onChange={(e) => setForgetKeep(Number(e.target.value))} min={1} />
              </div>
              <Button onClick={handleForget} disabled={forgetMut.isPending} variant="secondary" className="w-full">
                <Trash2 className="h-4 w-4 mr-2" />
                {forgetMut.isPending ? "Pruning..." : "Prune Old Snapshots"}
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <Shield className="h-4 w-4 text-emerald-500" />
                Maintenance
              </CardTitle>
              <CardDescription>Verify backup integrity</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              <Button onClick={handleCheck} disabled={checkMut.isPending} variant="outline" className="w-full">
                <RefreshCw className={`h-4 w-4 mr-2 ${checkMut.isPending ? "animate-spin" : ""}`} />
                {checkMut.isPending ? "Checking..." : "Check Repository"}
              </Button>
              <div className="space-y-2">
                <Label className="text-xs">Restore Target</Label>
                <Input value={restoreTarget} onChange={(e) => setRestoreTarget(e.target.value)} />
              </div>
            </CardContent>
          </Card>
        </div>
      </FadeIn>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Snapshots ({snapshots?.length ?? 0})</CardTitle>
          <CardDescription>Point-in-time backups available for restore</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <ListSkeleton rows={4} />
          ) : !snapshots || snapshots.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground space-y-3">
              <HardDrive className="h-8 w-8 mx-auto opacity-50" />
              <p className="text-sm">No snapshots found. Create your first backup.</p>
            </div>
          ) : (
            <div className="space-y-3">
              {snapshots.map((snap) => (
                <StaggerItem key={snap.id}>
                  <Card className="bg-muted/30">
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div className="space-y-1">
                          <div className="flex items-center gap-2">
                            <code className="text-sm font-mono text-primary">{snap.short_id}</code>
                            <span className="flex items-center gap-1 text-xs text-muted-foreground">
                              <Clock className="h-3 w-3" />
                              {new Date(snap.time).toLocaleString()}
                            </span>
                            {snap.tags?.map((t) => (
                              <Badge key={t} variant="secondary" className="text-[10px] gap-1">
                                <Tag className="h-2 w-2" /> {t}
                              </Badge>
                            ))}
                          </div>
                          <p className="text-xs text-muted-foreground">
                            {snap.hostname} &mdash; {snap.paths?.join(", ")}
                          </p>
                        </div>
                        <Button size="sm" onClick={() => handleRestore(snap.short_id)} disabled={restoreMut.isPending}>
                          <RotateCcw className="h-3 w-3 mr-1" />
                          {restoreMut.isPending ? "..." : "Restore"}
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                </StaggerItem>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </DashboardLayout>
  )
}
