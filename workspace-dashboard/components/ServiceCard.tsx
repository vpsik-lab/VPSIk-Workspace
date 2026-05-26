"use client"

import Link from "next/link"
import { motion } from "framer-motion"
import {
  Card,
  CardContent,
  CardDescription,
  CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

interface ServiceCardProps {
  name: string
  href: string
  description: string
  status: string
  icon: string
  index?: number
}

export default function ServiceCard({
  name,
  href,
  description,
  status,
  icon,
  index = 0,
}: ServiceCardProps) {
  const isHealthy = status === "healthy"

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, delay: index * 0.06 }}
      whileHover={{ y: -4 }}
    >
      <Link href={href} className="block group">
        <Card className="h-full transition-all duration-300 group-hover:border-primary/30 group-hover:shadow-lg group-hover:shadow-primary/5">
          <CardContent className="p-6">
            <div className="flex items-start justify-between mb-4">
              <span className="text-2xl transition-transform duration-300 group-hover:scale-110">
                {icon}
              </span>
              <Badge variant={isHealthy ? "success" : "destructive"}>
                <span
                  className={`w-1.5 h-1.5 rounded-full mr-1.5 ${
                    isHealthy ? "bg-emerald-400" : "bg-red-400"
                  }`}
                />
                {isHealthy ? "Healthy" : "Unhealthy"}
              </Badge>
            </div>
            <CardTitle className="text-foreground group-hover:text-primary transition-colors duration-200">
              {name}
            </CardTitle>
            <CardDescription className="mt-1">{description}</CardDescription>
          </CardContent>
        </Card>
      </Link>
    </motion.div>
  )
}
