"use client"

import IframeServicePage from "@/components/IframeServicePage"

export default function CodeServerPage() {
  const url = process.env.NEXT_PUBLIC_CODESERVER_URL || "http://localhost:8443"
  return <IframeServicePage title="Code Server" url={url} envVar="NEXT_PUBLIC_CODESERVER_URL" />
}
