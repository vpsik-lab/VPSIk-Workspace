import { cn } from "@/lib/utils"

function Pulse({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        "animate-pulse rounded-md bg-muted",
        className
      )}
    />
  )
}

export function CardSkeleton() {
  return (
    <div className="rounded-xl border bg-card p-6">
      <div className="flex items-start justify-between mb-4">
        <Pulse className="h-10 w-10 rounded-lg" />
        <Pulse className="h-6 w-16 rounded-full" />
      </div>
      <Pulse className="h-5 w-24 mb-2" />
      <Pulse className="h-4 w-32" />
    </div>
  )
}

export function ListSkeleton({ rows = 5 }: { rows?: number }) {
  return (
    <div className="space-y-3">
      {Array.from({ length: rows }).map((_, i) => (
        <div key={i} className="rounded-xl border bg-card p-5">
          <div className="flex items-center justify-between">
            <div className="space-y-2 flex-1">
              <Pulse className="h-4 w-1/4" />
              <Pulse className="h-3 w-1/2" />
            </div>
            <Pulse className="h-5 w-20 rounded-full" />
          </div>
        </div>
      ))}
    </div>
  )
}

export function PageSkeleton() {
  return (
    <div className="min-h-screen bg-background">
      <div className="hidden lg:flex flex-col w-60 h-screen fixed left-0 top-0 border-r bg-card">
        <div className="h-16 border-b flex items-center px-4 gap-2">
          <Pulse className="h-8 w-8 rounded-lg" />
          <Pulse className="h-5 w-28" />
        </div>
        <div className="p-3 space-y-2">
          {Array.from({ length: 8 }).map((_, i) => (
            <Pulse key={i} className="h-9 rounded-lg" />
          ))}
        </div>
      </div>
      <div className="lg:pl-60">
        <div className="h-16 border-b bg-background flex items-center justify-between px-8">
          <Pulse className="h-5 w-32" />
          <Pulse className="h-7 w-20 rounded-lg" />
        </div>
        <div className="p-8 space-y-6">
          <Pulse className="h-6 w-48" />
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {Array.from({ length: 6 }).map((_, i) => (
              <CardSkeleton key={i} />
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
