"use client"

import IframeServicePage from "@/components/IframeServicePage"

export default function OutlinePage() {
  const url = process.env.NEXT_PUBLIC_OUTLINE_URL || "http://localhost:3000"
  return <IframeServicePage title="Outline" url={url} envVar="NEXT_PUBLIC_OUTLINE_URL" />
}
