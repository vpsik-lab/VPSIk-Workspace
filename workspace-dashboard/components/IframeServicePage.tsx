"use client"

import { useState } from "react"
import DashboardLayout from "@/components/DashboardLayout"
import { Loader2, ExternalLink } from "lucide-react"
import { Button } from "@/components/ui/button"

interface IframeServicePageProps {
  title: string
  url: string
  envVar: string
  sandbox?: string
}

export default function IframeServicePage({
  title,
  url,
  sandbox = "allow-scripts allow-same-origin allow-forms allow-popups",
}: IframeServicePageProps) {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)

  return (
    <DashboardLayout title={title} showUpgrade={false}>
      <div className="flex-1 flex flex-col -m-4 md:-m-8">
        {loading && (
          <div className="flex-1 flex items-center justify-center">
            <div className="flex flex-col items-center gap-3">
              <Loader2 className="h-8 w-8 animate-spin text-primary" />
              <p className="text-sm text-muted-foreground">Loading {title}...</p>
            </div>
          </div>
        )}

        {error ? (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center space-y-4">
              <div className="text-4xl">🔌</div>
              <h2 className="text-lg font-semibold">Could not load {title}</h2>
              <p className="text-sm text-muted-foreground max-w-md">
                Make sure the {title} service is running and accessible via your workspace network.
              </p>
              <div className="flex gap-3 justify-center">
                <Button variant="outline" onClick={() => { setError(false); setLoading(true) }}>
                  Retry
                </Button>
                <Button variant="default" asChild>
                  <a href={url} target="_blank" rel="noopener noreferrer">
                    <ExternalLink className="h-4 w-4 mr-2" />
                    Open in new tab
                  </a>
                </Button>
              </div>
            </div>
          </div>
        ) : (
          <iframe
            src={url}
            className="flex-1 w-full border-0"
            title={title}
            sandbox={sandbox}
            onLoad={() => setLoading(false)}
            onError={() => { setError(true); setLoading(false) }}
            allow="clipboard-read; clipboard-write"
          />
        )}
      </div>
    </DashboardLayout>
  )
}
