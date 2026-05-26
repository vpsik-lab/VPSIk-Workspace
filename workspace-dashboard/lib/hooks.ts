/**
 * TanStack Query hooks for all WorkSpace OS API endpoints.
 *
 * Each hook wraps a raw API function with automatic caching, refetching,
 * loading/error states, and stale-while-revalidate semantics.
 *
 * MUTATION hooks auto-invalidate their related query caches on success
 * so the UI always reflects the latest server state.
 */

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import {
  getStatus,
  getGiteaRepos,
  getOllamaModels,
  pullOllamaModel,
  deleteOllamaModel,
  getCoolifyProjects,
  getCoolifyEnvironments,
  getCoolifyApplications,
  deployCoolifyResource,
  restartCoolifyResource,
  getResticSnapshots,
  createResticBackup,
  restoreResticSnapshot,
  forgetResticSnapshots,
  checkResticRepo,
  openCodeChat,
  type ServiceStatus,
  type GiteaRepo,
  type OllamaModel,
  type CoolifyProject,
  type CoolifyEnvironment,
  type CoolifyApplication,
  type ResticSnapshot,
} from "@/lib/api"

/** Fetch and auto-refresh service health status (every 15s). */
export function useServices() {
  return useQuery<{ services: ServiceStatus[] }>({
    queryKey: ["services"],
    queryFn: getStatus,
    refetchInterval: 15_000,
  })
}

/** Fetch the list of Gitea repositories. */
export function useGiteaRepos() {
  return useQuery<GiteaRepo[]>({
    queryKey: ["gitea", "repos"],
    queryFn: getGiteaRepos,
  })
}

/** Fetch and auto-refresh Ollama models (every 10s). */
export function useOllamaModels() {
  return useQuery<OllamaModel[]>({
    queryKey: ["ollama", "models"],
    queryFn: getOllamaModels,
    refetchInterval: 10_000,
  })
}

/** Pull a new Ollama model (invalidates model list on success). */
export function usePullOllamaModel() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: pullOllamaModel,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["ollama", "models"] })
    },
  })
}

/** Delete an Ollama model (invalidates model list on success). */
export function useDeleteOllamaModel() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: deleteOllamaModel,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["ollama", "models"] })
    },
  })
}

/** Fetch all Coolify projects. */
export function useCoolifyProjects() {
  return useQuery<CoolifyProject[]>({
    queryKey: ["coolify", "projects"],
    queryFn: getCoolifyProjects,
  })
}

/** Fetch environments for a specific Coolify project (enabled only when projectUUID is set). */
export function useCoolifyEnvironments(projectUUID: string | null) {
  return useQuery<CoolifyEnvironment[]>({
    queryKey: ["coolify", "environments", projectUUID],
    queryFn: () => getCoolifyEnvironments(projectUUID!),
    enabled: !!projectUUID,
  })
}

/** Fetch applications for a project+environment pair (enabled only when both are set). */
export function useCoolifyApplications(projectUUID: string | null, envName: string | null) {
  return useQuery<CoolifyApplication[]>({
    queryKey: ["coolify", "apps", projectUUID, envName],
    queryFn: () => getCoolifyApplications(projectUUID!, envName!),
    enabled: !!projectUUID && !!envName,
  })
}

/** Deploy a Coolify resource (invalidates all Coolify caches on success). */
export function useDeployCoolify() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: deployCoolifyResource,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["coolify"] })
    },
  })
}

/** Restart a Coolify resource (invalidates all Coolify caches on success). */
export function useRestartCoolify() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: restartCoolifyResource,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["coolify"] })
    },
  })
}

/** Fetch the list of all restic snapshots. */
export function useResticSnapshots() {
  return useQuery<ResticSnapshot[]>({
    queryKey: ["restic", "snapshots"],
    queryFn: getResticSnapshots,
  })
}

/** Create a new restic backup (invalidates snapshot list on success). */
export function useCreateBackup() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ paths, tags }: { paths: string[]; tags?: string[] }) =>
      createResticBackup(paths, tags),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["restic", "snapshots"] })
    },
  })
}

/** Restore a snapshot to a target path. */
export function useRestoreSnapshot() {
  return useMutation({
    mutationFn: ({ snapshotID, target }: { snapshotID: string; target: string }) =>
      restoreResticSnapshot(snapshotID, target),
  })
}

/** Prune old snapshots (invalidates snapshot list on success). */
export function useForgetSnapshots() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ keepLast, tags }: { keepLast: number; tags?: string[] }) =>
      forgetResticSnapshots(keepLast, tags),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["restic", "snapshots"] })
    },
  })
}

/** Run a restic repository integrity check. */
export function useCheckRestic() {
  return useMutation({
    mutationFn: checkResticRepo,
  })
}

/** Send a message to OpenCode.ai and get an AI response. */
export function useOpenCodeChat() {
  return useMutation({
    mutationFn: ({ message, context, repoPath }: { message: string; context?: string; repoPath?: string }) =>
      openCodeChat(message, context, repoPath),
  })
}
