interface StatusBadgeProps {
  status: string
}

const colors: Record<string, string> = {
  ok: 'bg-green-900/50 text-green-400 border-green-800',
  healthy: 'bg-green-900/50 text-green-400 border-green-800',
  running: 'bg-green-900/50 text-green-400 border-green-800',
  unhealthy: 'bg-red-900/50 text-red-400 border-red-800',
  error: 'bg-red-900/50 text-red-400 border-red-800',
  down: 'bg-red-900/50 text-red-400 border-red-800',
  degraded: 'bg-yellow-900/50 text-yellow-400 border-yellow-800',
  unknown: 'bg-gray-800 text-gray-400 border-gray-700',
}

const labels: Record<string, string> = {
  ok: 'Healthy',
  healthy: 'Healthy',
  running: 'Healthy',
  unhealthy: 'Unhealthy',
  error: 'Unhealthy',
  down: 'Unhealthy',
  degraded: 'Degraded',
  unknown: 'Unknown',
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const key = status?.toLowerCase() || 'unknown'
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${colors[key] || colors.unknown}`}>
      <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${
        key === 'ok' || key === 'healthy' || key === 'running' ? 'bg-green-400' :
        key === 'error' || key === 'down' || key === 'unhealthy' ? 'bg-red-400' :
        key === 'degraded' ? 'bg-yellow-400' : 'bg-gray-400'
      }`} />
      {labels[key] || labels.unknown}
    </span>
  )
}
