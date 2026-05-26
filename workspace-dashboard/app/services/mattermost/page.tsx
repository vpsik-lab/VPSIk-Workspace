"use client"

import IframeServicePage from "@/components/IframeServicePage"

export default function MattermostPage() {
  const url = process.env.NEXT_PUBLIC_MATTERMOST_URL || "http://localhost:8065"
  return <IframeServicePage title="Mattermost" url={url} envVar="NEXT_PUBLIC_MATTERMOST_URL" />
}
