import Link from 'next/link'

interface ServiceCardProps {
  name: string
  href: string
  description: string
  status: string
  icon: string
}

export default function ServiceCard({ name, href, description, status, icon }: ServiceCardProps) {
  const isHealthy = status === 'healthy'

  return (
    <Link
      href={href}
      className="bg-gray-900 border border-gray-800 rounded-xl p-6 hover:border-gray-700 transition group"
    >
      <div className="flex items-start justify-between mb-4">
        <span className="text-2xl">{icon}</span>
        <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
          isHealthy
            ? 'bg-green-900/50 text-green-400'
            : 'bg-red-900/50 text-red-400'
        }`}>
          {isHealthy ? 'Healthy' : 'Unhealthy'}
        </span>
      </div>
      <h3 className="text-white font-semibold group-hover:text-vpsik-400 transition">
        {name}
      </h3>
      <p className="text-sm text-gray-500 mt-1">
        {description}
      </p>
    </Link>
  )
}
