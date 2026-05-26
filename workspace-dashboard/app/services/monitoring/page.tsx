"use client"

import DashboardLayout from "@/components/DashboardLayout"
import StatusBadge from "@/components/StatusBadge"
import { ListSkeleton } from "@/components/LoadingSkeleton"
import { FadeIn, StaggerItem } from "@/components/motion-wrapper"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Activity,
  BarChart3,
  TrendingUp,
  ExternalLink,
  AlertTriangle,
} from "lucide-react"
import { useServices } from "@/lib/hooks"

const GRAFANA_URL = process.env.NEXT_PUBLIC_GRAFANA_URL || "http://localhost:3002"
const PROMETHEUS_URL = process.env.NEXT_PUBLIC_PROMETHEUS_URL || "http://localhost:9090"

export default function MonitoringPage() {
  const { data, isLoading, error } = useServices()
  const services = data?.services ?? []
  const healthyCount = services.filter((s) => s.status === "healthy").length

  return (
    <DashboardLayout
      title="Monitoring"
      subtitle={`${healthyCount}/${services.length} healthy`}
    >
      {/* External Links */}
      <FadeIn>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <a
            href={GRAFANA_URL}
            target="_blank"
            rel="noopener noreferrer"
            className="group"
          >
            <Card className="hover:border-primary/30 transition-all duration-200 h-full">
              <CardContent className="p-6 flex items-start justify-between">
                <div>
                  <div className="p-2 rounded-lg bg-primary/5 w-fit mb-3">
                    <BarChart3 className="h-5 w-5 text-primary" />
                  </div>
                  <CardTitle className="text-base group-hover:text-primary transition-colors">
                    Grafana
                  </CardTitle>
                  <CardDescription className="mt-1">
                    Dashboards, alerts, and metrics visualization
                  </CardDescription>
                </div>
                <ExternalLink className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors mt-1" />
              </CardContent>
            </Card>
          </a>
          <a
            href={PROMETHEUS_URL}
            target="_blank"
            rel="noopener noreferrer"
            className="group"
          >
            <Card className="hover:border-primary/30 transition-all duration-200 h-full">
              <CardContent className="p-6 flex items-start justify-between">
                <div>
                  <div className="p-2 rounded-lg bg-emerald-500/5 w-fit mb-3">
                    <TrendingUp className="h-5 w-5 text-emerald-500" />
                  </div>
                  <CardTitle className="text-base group-hover:text-emerald-500 transition-colors">
                    Prometheus
                  </CardTitle>
                  <CardDescription className="mt-1">
                    Time-series metrics and alerting rules
                  </CardDescription>
                </div>
                <ExternalLink className="h-4 w-4 text-muted-foreground group-hover:text-emerald-500 transition-colors mt-1" />
              </CardContent>
            </Card>
          </a>
        </div>
      </FadeIn>

      {/* Service Health */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Activity className="h-4 w-4 text-primary" />
            Service Health
          </CardTitle>
          <CardDescription>
            Real-time status of all workspace services
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <ListSkeleton rows={6} />
          ) : error ? (
            <div className="flex items-center gap-2 text-destructive text-sm">
              <AlertTriangle className="h-4 w-4" />
              Failed to fetch service status
            </div>
          ) : services.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground space-y-3">
              <Activity className="h-8 w-8 mx-auto opacity-50" />
              <p className="text-sm">No services reported</p>
            </div>
          ) : (
            <div className="divide-y">
              {services.map((svc, i) => (
                <StaggerItem key={svc.name}>
                  <div className="flex items-center justify-between py-3 first:pt-0 last:pb-0">
                    <div className="flex items-center gap-3">
                      <span className="text-sm font-medium capitalize">{svc.name}</span>
                      {svc.error && (
                        <span className="text-xs text-muted-foreground truncate max-w-[200px]">
                          {svc.error}
                        </span>
                      )}
                    </div>
                    <StatusBadge status={svc.status} />
                  </div>
                </StaggerItem>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </DashboardLayout>
  )
}
