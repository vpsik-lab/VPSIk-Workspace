"use client"

import { useEffect, useState, useRef, useCallback } from "react"
import DashboardLayout from "@/components/DashboardLayout"
import { FadeIn, StaggerItem } from "@/components/motion-wrapper"
import { chatOllamaStream, runOllamaTask } from "@/lib/api"
import { useOllamaModels } from "@/lib/hooks"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
} from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"
import {
  Send,
  Square,
  Sparkles,
  Lightbulb,
  ScrollText,
  Code,
  BookOpen,
  Bot,
  User,
} from "lucide-react"

interface Message {
  role: "user" | "assistant"
  content: string
}

const TASK_TEMPLATES = [
  { id: "explain", label: "Explain", icon: Lightbulb, prompt: "Explain this in simple terms:\n\n" },
  { id: "summarize", label: "Summarize", icon: ScrollText, prompt: "Summarize this concisely:\n\n" },
  { id: "review", label: "Review Code", icon: Code, prompt: "Review this code for bugs and improvements:\n\n" },
  { id: "expand", label: "Expand", icon: BookOpen, prompt: "Expand on this with more detail:\n\n" },
]

export default function ChatPage() {
  const { data: models } = useOllamaModels()
  const [selectedModel, setSelectedModel] = useState("")
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState("")
  const [sending, setSending] = useState(false)
  const [streamingContent, setStreamingContent] = useState("")
  const [showTasks, setShowTasks] = useState(false)
  const abortRef = useRef<AbortController | null>(null)
  const bottomRef = useRef<HTMLDivElement>(null)
  const messagesRef = useRef<Message[]>([])
  messagesRef.current = messages

  useEffect(() => {
    if (models && models.length > 0 && !selectedModel) {
      setSelectedModel(models[0].name)
    }
  }, [models, selectedModel])

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages, streamingContent])

  useEffect(() => {
    return () => {
      if (abortRef.current) {
        abortRef.current.abort()
        abortRef.current = null
      }
    }
  }, [])

  const handleSend = useCallback(
    async (content: string) => {
      if (!content.trim() || !selectedModel || sending) return

      const userMsg: Message = { role: "user", content: content.trim() }
      setMessages((prev) => [...prev, userMsg])
      setInput("")
      setSending(true)
      setStreamingContent("")
      setShowTasks(false)

      const apiMessages = [...messagesRef.current, userMsg].map((m) => ({
        role: m.role,
        content: m.content,
      }))

      let fullReply = ""
      abortRef.current = await chatOllamaStream(
        selectedModel,
        apiMessages,
        (chunk) => {
          fullReply += chunk
          setStreamingContent(fullReply)
        },
        () => {
          setMessages((prev) => [...prev, { role: "assistant", content: fullReply }])
          setStreamingContent("")
          setSending(false)
        },
        (err) => {
          setMessages((prev) => [...prev, { role: "assistant", content: `Error: ${err.message}` }])
          setStreamingContent("")
          setSending(false)
        },
      )
    },
    [selectedModel, sending],
  )

  async function handleTask(taskId: string) {
    const task = TASK_TEMPLATES.find((t) => t.id === taskId)
    if (!task || !input.trim() || !selectedModel || sending) return

    const userMsg: Message = { role: "user", content: `${input.trim()}` }
    setMessages((prev) => [...prev, userMsg])
    setInput("")
    setSending(true)
    setStreamingContent("")
    setShowTasks(false)

    try {
      const reply = await runOllamaTask(selectedModel, taskId, input.trim())
      setMessages((prev) => [...prev, { role: "assistant", content: reply }])
    } catch (err: unknown) {
      setMessages((prev) => [
        ...prev,
        { role: "assistant", content: `Error: ${err instanceof Error ? err.message : String(err)}` },
      ])
    }
    setSending(false)
  }

  function stopStreaming() {
    if (abortRef.current) {
      abortRef.current.abort()
      abortRef.current = null
    }
  }

  return (
    <DashboardLayout
      title="AI Chat"
      showUpgrade={false}
      actions={
        <Select value={selectedModel} onValueChange={setSelectedModel}>
          <SelectTrigger className="w-40 h-8 text-xs">
            <SelectValue placeholder="Select model" />
          </SelectTrigger>
          <SelectContent>
            {(models ?? []).map((m) => (
              <SelectItem key={m.name} value={m.name}>
                {m.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      }
    >
      <div className="flex-1 flex flex-col -m-4 md:-m-8">
        <div className="flex-1 overflow-y-auto p-4 md:p-8 space-y-4">
          {messages.length === 0 && !streamingContent && (
            <FadeIn>
              <div className="flex items-center justify-center min-h-[60vh]">
                <div className="text-center space-y-6 max-w-md">
                  <div className="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center mx-auto">
                    <Sparkles className="h-8 w-8 text-primary" />
                  </div>
                  <h2 className="text-xl font-semibold">Start a conversation</h2>
                  <p className="text-sm text-muted-foreground">
                    Ask questions, analyze code, generate content, or get help with your workspace.
                  </p>
                  <div className="flex flex-wrap gap-2 justify-center">
                    {TASK_TEMPLATES.map((task) => {
                      const Icon = task.icon
                      return (
                        <Badge key={task.id} variant="secondary" className="gap-1.5 px-3 py-1.5 cursor-pointer hover:bg-secondary/80">
                          <Icon className="h-3 w-3" />
                          {task.label}
                        </Badge>
                      )
                    })}
                  </div>
                </div>
              </div>
            </FadeIn>
          )}

          {messages.map((msg, i) => (
            <StaggerItem key={i}>
              <div className={`flex ${msg.role === "user" ? "justify-end" : "justify-start"}`}>
                <div className="flex gap-3 max-w-[80%] md:max-w-[65%]">
                  {msg.role === "assistant" && (
                    <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0 mt-1">
                      <Bot className="h-4 w-4 text-primary" />
                    </div>
                  )}
                  <div>
                    <Card>
                      <CardContent className={`px-4 py-3 ${msg.role === "user" ? "bg-primary/5" : ""}`}>
                        <p className="text-xs text-muted-foreground mb-1">
                          {msg.role === "user" ? "You" : "AI"}
                        </p>
                        <p className="text-sm whitespace-pre-wrap">{msg.content}</p>
                      </CardContent>
                    </Card>
                  </div>
                  {msg.role === "user" && (
                    <div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center flex-shrink-0 mt-1">
                      <User className="h-4 w-4 text-primary-foreground" />
                    </div>
                  )}
                </div>
              </div>
            </StaggerItem>
          ))}

          {streamingContent && (
            <div className="flex justify-start">
              <div className="flex gap-3 max-w-[80%] md:max-w-[65%]">
                <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0 mt-1">
                  <Bot className="h-4 w-4 text-primary" />
                </div>
                <div>
                  <Card>
                    <CardContent className="px-4 py-3">
                      <p className="text-xs text-muted-foreground mb-1">AI</p>
                      <p className="text-sm whitespace-pre-wrap">{streamingContent}</p>
                    </CardContent>
                  </Card>
                </div>
              </div>
            </div>
          )}
          <div ref={bottomRef} />
        </div>

        <div className="border-t p-4">
          <form
            onSubmit={(e) => {
              e.preventDefault()
              handleSend(input)
            }}
            className="max-w-3xl mx-auto space-y-2"
          >
            {showTasks && input.trim() && !sending && (
              <div className="flex gap-2">
                {TASK_TEMPLATES.map((task) => {
                  const Icon = task.icon
                  return (
                    <Tooltip key={task.id}>
                      <TooltipTrigger asChild>
                        <Button
                          type="button"
                          variant="secondary"
                          size="sm"
                          onClick={() => handleTask(task.id)}
                        >
                          <Icon className="h-3 w-3 mr-1" />
                          {task.label}
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>Run AI task</TooltipContent>
                    </Tooltip>
                  )
                })}
              </div>
            )}
            <div className="flex gap-3">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    type="button"
                    variant="outline"
                    size="icon"
                    onClick={() => setShowTasks(!showTasks)}
                    className="flex-shrink-0"
                  >
                    <Sparkles className="h-4 w-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>AI Task templates</TooltipContent>
              </Tooltip>

              <Input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Type your message..."
                disabled={sending}
                className="flex-1"
              />

              {sending ? (
                <Button type="button" variant="destructive" onClick={stopStreaming}>
                  <Square className="h-4 w-4 mr-2" />
                  Stop
                </Button>
              ) : (
                <Button type="submit" disabled={!input.trim()}>
                  <Send className="h-4 w-4 mr-2" />
                  Send
                </Button>
              )}
            </div>
          </form>
        </div>
      </div>
    </DashboardLayout>
  )
}
