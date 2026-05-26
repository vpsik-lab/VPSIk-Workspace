"use client"

import { Badge } from "@/components/ui/badge"

interface StatusBadgeProps {
  status: string
}

const colorMap: Record<string, "success" | "destructive" | "warning" | "secondary"> = {
  ok: "success",
  healthy: "success",
  running: "success",
  unhealthy: "destructive",
  error: "destructive",
  down: "destructive",
  degraded: "warning",
  unknown: "secondary",
}

const labelMap: Record<string, string> = {
  ok: "Healthy",
  healthy: "Healthy",
  running: "Healthy",
  unhealthy: "Unhealthy",
  error: "Unhealthy",
  down: "Unhealthy",
  degraded: "Degraded",
  unknown: "Unknown",
}

const dotColors: Record<string, string> = {
  success: "bg-emerald-400",
  destructive: "bg-red-400",
  warning: "bg-amber-400",
  secondary: "bg-muted-foreground",
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const key = status?.toLowerCase() || "unknown"
  const variant = colorMap[key] || "secondary"
  const label = labelMap[key] || labelMap.unknown

  return (
    <Badge variant={variant}>
      <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${dotColors[variant]}`} />
      {label}
    </Badge>
  )
}
