"use client"

import { useState, useEffect } from "react"
import { ClipboardCopy } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { toast } from "@/components/ui/use-toast"
import { Progress } from "@/components/ui/progress"

export default function URLShortener() {
  const [longUrl, setLongUrl] = useState("")
  const [shortUrl, setShortUrl] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [cooldown, setCooldown] = useState(0)
  const COOLDOWN_TIME = 10 // 冷却时间（秒）

  useEffect(() => {
    let timer: NodeJS.Timeout
    if (cooldown > 0) {
      timer = setInterval(() => {
        setCooldown((prev) => Math.max(0, prev - 0.1))
      }, 100)
    }
    return () => clearInterval(timer)
  }, [cooldown])

  const isValidUrl = (url: string) => {
    try {
      new URL(url)
      return true
    } catch (e) {
      return false
    }
  }

  const generateShortUrl = async () => {
    if (cooldown > 0) {
      toast({
        title: "请稍后再试",
        description: `还需等待 ${cooldown.toFixed(1)} 秒`,
        variant: "destructive",
      })
      return
    }

    if (!isValidUrl(longUrl)) {
      toast({
        title: "无效的 URL",
        description: "请输入一个有效的 URL，包括 http:// 或 https://",
        variant: "destructive",
      })
      return
    }

    setIsLoading(true)
    try {
      const response = await fetch("http://localhost:8081/api/shorten", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url: longUrl }), // 修改为正确的字段名 url
      });

      if (!response.ok) {
        throw new Error("生成短链接失败");
      }

      const data = await response.json();
      setShortUrl(data.short_url);
    } catch (error: any) {
      toast({
        title: "错误",
        description: error.message || "发生未知错误",
        variant: "destructive",
      });
    }
    setIsLoading(false)
    setCooldown(COOLDOWN_TIME)
  }

  const copyToClipboard = () => {
    navigator.clipboard.writeText(shortUrl)
    toast({
      title: "已复制",
      description: "短链接已复制到剪贴板",
    })
  }

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center px-4">
      {/* 短链接生成器 */}
      <div className="max-w-md w-full space-y-8 bg-white p-6 rounded-xl shadow-md">
        <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">短链接生成器</h2>
        <div className="mt-8 space-y-6">
          <Input
            type="url"
            placeholder="输入长 URL"
            value={longUrl}
            onChange={(e) => setLongUrl(e.target.value)}
            className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
          />
          <div className="space-y-2">
            <Button
              onClick={generateShortUrl}
              disabled={isLoading || cooldown > 0}
              className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              {isLoading ? "生成中..." : cooldown > 0 ? `冷却中 (${cooldown.toFixed(1)}s)` : "生成短链接"}
            </Button>
            {cooldown > 0 && (
              <div className="space-y-1">
                <Progress value={((COOLDOWN_TIME - cooldown) / COOLDOWN_TIME) * 100} />
                <p className="text-xs text-gray-500 text-center">下次生成还需等待 {cooldown.toFixed(1)} 秒</p>
              </div>
            )}
          </div>
          {shortUrl && (
            <div className="mt-4">
              <div className="flex items-center justify-between bg-gray-100 p-3 rounded-md">
                <span className="text-sm text-gray-800">{shortUrl}</span>
                <Button onClick={copyToClipboard} variant="ghost" size="icon" className="ml-2">
                  <ClipboardCopy className="h-4 w-4" />
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* 个人信息 */}
      <div className="mt-8 text-center">
        <p className="text-sm text-gray-500">
          <strong>ShortLinkWizard</strong>，由{" "}
          <a
            href="https://github.com/JimmyPowell/ShortLinkWizard"
            target="_blank"
            rel="noopener noreferrer"
            className="text-indigo-600 hover:underline"
          >
            JimmyPowell
          </a>{" "}
          开发
        </p>
      </div>
    </div>
  )
}
