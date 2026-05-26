"use client"

import { useState } from "react"
import DashboardLayout from "@/components/DashboardLayout"
import { FadeIn } from "@/components/motion-wrapper"
import { useOpenCodeChat } from "@/lib/hooks"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Card,
  CardContent,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  Code2,
  Send,
  Loader2,
  Lightbulb,
} from "lucide-react"

const SUGGESTIONS = [
  "Explain this codebase architecture",
  "How do I add a new API endpoint?",
  "Find potential bugs in the authentication flow",
  "Suggest improvements to the deployment pipeline",
  "Generate a README for this project",
]

export default function OpenCodePage() {
  const [input, setInput] = useState("")
  const [context, setContext] = useState("")
  const [reply, setReply] = useState("")
  const [error, setError] = useState("")
  const chatMut = useOpenCodeChat()

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!input.trim() || chatMut.isPending) return

    setError("")
    setReply("")

    try {
      const result = await chatMut.mutateAsync({ message: input.trim(), context: context.trim() || undefined })
      setReply(result)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : String(err))
    }
  }

  return (
    <DashboardLayout title="OpenCode.ai" showUpgrade={false}>
      <div className="flex-1 flex flex-col -m-4 md:-m-8">
        <main className="flex-1 overflow-y-auto p-4 md:p-8">
          {!reply && !chatMut.isPending && (
            <FadeIn>
              <div className="flex items-center justify-center min-h-[50vh]">
                <div className="text-center space-y-6 max-w-lg">
                  <div className="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center mx-auto">
                    <Code2 className="h-8 w-8 text-primary" />
                  </div>
                  <h2 className="text-xl font-semibold">AI Code Assistant</h2>
                  <p className="text-sm text-muted-foreground">
                    Ask questions about your codebase, generate code, review PRs, and more.
                  </p>
                  <div className="flex flex-wrap gap-2 justify-center">
                    {SUGGESTIONS.map((s, i) => (
                      <Badge key={i} variant="secondary"
                        className="cursor-pointer hover:bg-secondary/80 px-3 py-1.5"
                        onClick={() => setInput(s)}
                      >
                        <Lightbulb className="h-3 w-3 mr-1" /> {s}
                      </Badge>
                    ))}
                  </div>
                </div>
              </div>
            </FadeIn>
          )}

          {error && (
            <div className="mb-4 rounded-md bg-destructive/10 border border-destructive/20 px-4 py-3 text-sm text-destructive">
              {error}
            </div>
          )}

          {chatMut.isPending && (
            <div className="flex items-center justify-center py-12">
              <div className="flex flex-col items-center gap-3">
                <Loader2 className="h-6 w-6 animate-spin text-primary" />
                <p className="text-sm text-muted-foreground">Analyzing your code...</p>
              </div>
            </div>
          )}

          {reply && !chatMut.isPending && (
            <Card>
              <CardContent className="p-5">
                <div className="flex items-center gap-2 mb-3">
                  <Code2 className="h-4 w-4 text-primary" />
                  <span className="text-xs text-muted-foreground font-medium">OpenCode.ai</span>
                </div>
                <p className="text-sm whitespace-pre-wrap leading-relaxed">{reply}</p>
              </CardContent>
            </Card>
          )}
        </main>

        <div className="border-t p-4">
          <form onSubmit={handleSubmit} className="max-w-3xl mx-auto space-y-2">
            <Input
              type="text" value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Ask OpenCode.ai about your code..."
              disabled={chatMut.isPending}
            />
            <div className="flex gap-3">
              <Input type="text" value={context}
                onChange={(e) => setContext(e.target.value)}
                placeholder="Optional context (file path, project name)..."
                disabled={chatMut.isPending} className="flex-1"
              />
              <Button type="submit" disabled={!input.trim() || chatMut.isPending}>
                {chatMut.isPending ? (
                  <><Loader2 className="h-4 w-4 animate-spin mr-2" />Sending...</>
                ) : (
                  <><Send className="h-4 w-4 mr-2" />Send</>
                )}
              </Button>
            </div>
          </form>
        </div>
      </div>
    </DashboardLayout>
  )
}
