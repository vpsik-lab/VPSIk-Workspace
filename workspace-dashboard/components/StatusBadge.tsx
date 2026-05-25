interface StatusBadgeProps {
  status: 'healthy' | 'unhealthy' | 'unknown'
}

const colors = {
  healthy: 'bg-green-900/50 text-green-400 border-green-800',
  unhealthy: 'bg-red-900/50 text-red-400 border-red-800',
  unknown: 'bg-gray-800 text-gray-400 border-gray-700',
}

const labels = {
  healthy: 'Healthy',
  unhealthy: 'Unhealthy',
  unknown: 'Unknown',
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${colors[status]}`}>
      <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${
        status === 'healthy' ? 'bg-green-400' :
        status === 'unhealthy' ? 'bg-red-400' : 'bg-gray-400'
      }`} />
      {labels[status]}
    </span>
  )
}
