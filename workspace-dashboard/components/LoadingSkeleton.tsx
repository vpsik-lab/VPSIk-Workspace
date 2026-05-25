export function CardSkeleton() {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-5 animate-pulse">
      <div className="flex items-start justify-between">
        <div className="space-y-3 flex-1">
          <div className="h-4 bg-gray-700 rounded w-1/3" />
          <div className="h-3 bg-gray-700 rounded w-2/3" />
        </div>
        <div className="h-6 w-16 bg-gray-700 rounded-full" />
      </div>
    </div>
  )
}

export function ListSkeleton({ rows = 5 }: { rows?: number }) {
  return (
    <div className="space-y-3 animate-pulse">
      {Array.from({ length: rows }).map((_, i) => (
        <div key={i} className="bg-gray-900 border border-gray-800 rounded-xl p-5">
          <div className="flex items-center justify-between">
            <div className="space-y-2 flex-1">
              <div className="h-4 bg-gray-700 rounded w-1/4" />
              <div className="h-3 bg-gray-700 rounded w-1/2" />
            </div>
            <div className="h-5 w-20 bg-gray-700 rounded" />
          </div>
        </div>
      ))}
    </div>
  )
}

export function PageSkeleton() {
  return (
    <div className="flex min-h-screen bg-gray-950">
      <div className="w-60 bg-gray-900 border-r border-gray-800 animate-pulse">
        <div className="h-16 border-b border-gray-800 px-6 flex items-center">
          <div className="h-5 w-20 bg-gray-700 rounded" />
        </div>
        <div className="p-3 space-y-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="h-10 bg-gray-800 rounded-lg" />
          ))}
        </div>
      </div>
      <div className="flex-1">
        <div className="h-16 bg-gray-900 border-b border-gray-800" />
        <div className="p-8 space-y-4">
          <div className="h-6 w-48 bg-gray-700 rounded animate-pulse" />
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="h-32 bg-gray-900 border border-gray-800 rounded-xl animate-pulse" />
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
