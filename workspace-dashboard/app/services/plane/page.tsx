"use client"

import IframeServicePage from "@/components/IframeServicePage"

export default function PlanePage() {
  const url = process.env.NEXT_PUBLIC_PLANE_URL || "http://localhost:8080"
  return <IframeServicePage title="Plane" url={url} envVar="NEXT_PUBLIC_PLANE_URL" />
}
