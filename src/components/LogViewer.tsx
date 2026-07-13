import { useEffect, useRef, useState } from "react"
import { Terminal, Copy, Check, FileText, Trash2 } from "lucide-react"
import { invoke } from "@tauri-apps/api/core"

interface LogViewerProps {
  log: string[]
  title: string
  serviceId?: string
  onClear?: () => void
}

export function LogViewer({ log, title, serviceId, onClear }: LogViewerProps) {
  const scrollContainerRef = useRef<HTMLDivElement>(null)
  const [copied, setCopied] = useState(false)

  useEffect(() => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollTop = scrollContainerRef.current.scrollHeight
    }
  }, [log])

  const handleCopy = () => {
    navigator.clipboard.writeText(log.join('\n'))
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const handleOpenFile = async () => {
    if (serviceId) {
      try {
        await invoke("open_log_file", { id: serviceId })
      } catch (e) {
        console.error("Failed to open log file:", e)
      }
    }
  }

  const handleClearLog = async () => {
    if (onClear) {
      onClear()
    } else if (serviceId) {
      try {
        await invoke("clear_service_log", { id: serviceId })
      } catch (e) {
        console.error("Failed to clear log:", e)
      }
    }
  }

  return (
    <div className="bg-[#242533] rounded-xl border border-[#2a2a35] overflow-hidden flex flex-col h-full min-h-[150px]">
      <div className="px-5 py-3 border-b border-[#2a2a35] flex items-center gap-2 bg-[#242533]">
        <Terminal className="w-4 h-4 text-gray-400 shrink-0" />
        <span className="text-sm font-semibold text-gray-300 truncate">{title}</span>
        <div className="ml-auto flex items-center gap-3 shrink-0">
          <span className="text-xs text-gray-600 hidden sm:inline-block">{log.length} lines</span>
          <div className="h-4 w-px bg-[#2a2a35]" />
          <button 
            onClick={handleCopy}
            title="Copy logs"
            className="text-gray-400 hover:text-white transition-colors bg-[#2a2a35] p-1.5 rounded-md hover:bg-[#323240]"
          >
            {copied ? <Check className="w-3.5 h-3.5 text-green-400" /> : <Copy className="w-3.5 h-3.5" />}
          </button>
          {(serviceId || onClear) && (
            <>
              {serviceId && (
                <button 
                  onClick={handleOpenFile}
                  title="Open full log file"
                  className="text-gray-400 hover:text-white transition-colors bg-[#2a2a35] p-1.5 rounded-md hover:bg-[#323240]"
                >
                  <FileText className="w-3.5 h-3.5" />
                </button>
              )}
              <button 
                onClick={handleClearLog}
                title="Clear logs"
                className="text-red-400/70 hover:text-red-400 transition-colors bg-red-500/10 p-1.5 rounded-md hover:bg-red-500/20"
              >
                <Trash2 className="w-3.5 h-3.5" />
              </button>
            </>
          )}
        </div>
      </div>
      <div ref={scrollContainerRef} className="flex-1 bg-[#1a1b26] p-4 overflow-auto custom-scrollbar select-text">
        {log.length === 0 ? (
          <div className="text-[13px] font-mono text-gray-600 italic">No output yet...</div>
        ) : (
          <div className="space-y-0">
            {log.map((line, i) => {
              // Extract common log patterns
              const renderLine = (text: string) => {
                // Escape HTML first to prevent XSS and rendering issues
                let formatted = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
                
                // Highlight ERROR/FATAL
                formatted = formatted.replace(/\b(ERROR|FATAL|ERR|Exception|Failed)\b/gi, "<span class='text-red-400 font-bold'>$1</span>");
                // Highlight WARN
                formatted = formatted.replace(/\b(WARN|WARNING)\b/gi, "<span class='text-yellow-400 font-bold'>$1</span>");
                // Highlight INFO/SUCCESS
                formatted = formatted.replace(/\b(INFO|SUCCESS|OK|Ready)\b/gi, "<span class='text-green-400 font-bold'>$1</span>");
                
                // Highlight URLs (avoid overlapping with IPs by doing URLs first)
                formatted = formatted.replace(/(https?:\/\/[^\s]+)/g, "<span class='text-blue-300 underline'>$1</span>");
                
                // Highlight IPs (only if not inside an HTML tag, simple workaround: IPs rarely overlap with our span tags)
                formatted = formatted.replace(/(\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b)/g, "<span class='text-cyan-300'>$1</span>");
                
                // Highlight Timestamps
                formatted = formatted.replace(/(\b\d{4}-\d{2}-\d{2}[T\s]\d{2}:\d{2}:\d{2}.*?\b)/g, "<span class='text-gray-500'>$1</span>");
                
                // Highlight Quoted strings (only double quotes, since our HTML uses single quotes)
                formatted = formatted.replace(/(&quot;|")[^"]*(&quot;|")/g, "<span class='text-amber-300'>$&</span>");

                return <span dangerouslySetInnerHTML={{ __html: formatted }} />;
              };

              return (
                <div key={i} className="flex gap-3 text-[13px] font-mono leading-6 hover:bg-[#1e1f2e] -mx-2 px-2 rounded break-all">
                  <span className="text-gray-600 select-none w-8 text-right shrink-0">{i + 1}</span>
                  <span className="text-[#a9b1d6] flex-1">
                    {renderLine(line)}
                  </span>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  )
}
