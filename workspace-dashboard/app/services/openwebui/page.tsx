"use client"

import IframeServicePage from "@/components/IframeServicePage"

export default function OpenWebUIPage() {
  const url = process.env.NEXT_PUBLIC_OPENWEBUI_URL || "http://localhost:3001"
  return <IframeServicePage title="Open WebUI" url={url} envVar="NEXT_PUBLIC_OPENWEBUI_URL" />
}
