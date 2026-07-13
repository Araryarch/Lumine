import type { Service } from "../types"
import { Play, Square, Trash2, Pencil, Terminal, Activity, ExternalLink } from "lucide-react"
import { invoke } from "@tauri-apps/api/core"
import { useAppStore } from "../store/useAppStore"
interface ServiceCardProps {
  service: Service
  onToggle: () => void
  onPortChange: (port: number) => void
  onEdit?: () => void
  onDelete?: () => void
  isLoading?: boolean
  adminPanelName?: string | null
}

export function ServiceCard({ service, onToggle, onPortChange, onEdit, onDelete, isLoading, adminPanelName }: ServiceCardProps) {
  const isRunning = service.status === "running"
  const isPending = isLoading || service.status === "starting" || service.status === "stopping"
  const stats = useAppStore(s => s.dockerStats[service.id])

  return (
    <div className="bg-[#242533] rounded-xl border border-[#2a2a35] overflow-hidden">
      <div className="p-5 border-b border-[#2a2a35] flex items-center justify-between">
        <div>
          <div className="flex items-center gap-3">
            <h3 className="text-lg font-semibold text-white">{service.name}</h3>
            {isRunning && stats && (
              <div className="flex items-center gap-2 px-2.5 py-1 rounded-md bg-[#181825] border border-[#2a2a35] text-xs font-mono text-gray-400">
                <Activity className="w-3 h-3 text-blue-400" />
                <span>CPU: {stats.cpu}</span>
                <span className="w-px h-3 bg-gray-700 mx-1"></span>
                <span>RAM: {stats.ram}</span>
              </div>
            )}
          </div>
          <p className="text-[#8c8d9e] text-sm mt-1">{service.description}</p>
        </div>
        <div className="flex gap-2">
          {onEdit && (
            <button
              onClick={onEdit}
              disabled={isPending}
              className="p-2 rounded-lg bg-gray-500/10 text-gray-400 hover:bg-gray-500/20 disabled:opacity-50 transition-colors"
              title="Edit Service"
            >
              <Pencil className="w-4 h-4" />
            </button>
          )}
          {onDelete && (
            <button
              onClick={onDelete}
              disabled={isPending}
              className="p-2 rounded-lg bg-red-500/10 text-red-400 hover:bg-red-500/20 disabled:opacity-50 transition-colors"
              title="Delete Service"
            >
              <Trash2 className="w-4 h-4" />
            </button>
          )}
          {service.config.serviceType === "Language" ? (
            <button
              onClick={() => invoke("open_service_terminal", { executable: service.config.executablePath })}
              className="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all bg-blue-500/10 text-blue-400 hover:bg-blue-500/20"
            >
              <Terminal className="w-4 h-4" />
              Open Terminal
            </button>
          ) : service.config.serviceType === "Database" && adminPanelName ? (
            <>
              <button
                onClick={() => {
                  const url = `http://localhost/${adminPanelName.toLowerCase()}`;
                  invoke("open_in_browser", { url }).catch(() => window.open(url, '_blank'));
                }}
                className="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20"
              >
                Open {adminPanelName}
              </button>
              <button
                onClick={onToggle}
                disabled={isPending}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                  isRunning
                    ? "bg-red-500/10 text-red-400 hover:bg-red-500/20"
                    : "bg-blue-500/10 text-blue-400 hover:bg-blue-500/20"
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                {isPending ? (
                  <span className="w-4 h-4 rounded-full border-2 border-current border-t-transparent animate-spin" />
                ) : isRunning ? (
                  <Square className="w-4 h-4 fill-current" />
                ) : (
                  <Play className="w-4 h-4 fill-current" />
                )}
                {isPending ? "Pending..." : isRunning ? "Stop Service" : "Start Service"}
              </button>
            </>
          ) : (
            <>
              {isRunning && (
                <button
                  onClick={() => invoke("open_docker_shell", { containerName: `lumine_${service.id}` })}
                  className="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all bg-blue-500/10 text-blue-400 hover:bg-blue-500/20"
                >
                  <Terminal className="w-4 h-4" />
                  Open Shell
                </button>
              )}
              <button
                onClick={onToggle}
                disabled={isPending}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                  isRunning
                    ? "bg-red-500/10 text-red-400 hover:bg-red-500/20"
                    : "bg-blue-500/10 text-blue-400 hover:bg-blue-500/20"
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                {isPending ? (
                  <span className="w-4 h-4 rounded-full border-2 border-current border-t-transparent animate-spin" />
                ) : isRunning ? (
                  <Square className="w-4 h-4 fill-current" />
                ) : (
                  <Play className="w-4 h-4 fill-current" />
                )}
                {isPending ? "Pending..." : isRunning ? "Stop Service" : "Start Service"}
              </button>
            </>
          )}
        </div>
      </div>

      {service.config.serviceType !== "Language" && service.config.port !== 0 && (
        <div className="p-5 flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <label className="text-xs font-semibold text-gray-500 uppercase tracking-wider">Port Configuration</label>
            <div className="flex items-center gap-3">
              <div className="flex items-center">
                <input
                  type="number"
                  value={service.config.port}
                  onChange={(e) => onPortChange(Number(e.target.value))}
                  className="w-24 bg-[#1a1b26] border border-[#2a2a35] rounded-l-lg px-3 py-2 text-white 
                             focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all shadow-inner"
                  min={1}
                  max={65535}
                />
                <div className="bg-[#1a1b26] border border-l-0 border-[#2a2a35] rounded-r-lg px-3 py-2 text-gray-500 text-sm shadow-inner shrink-0">
                  :{service.config.containerPort || service.config.port}
                </div>
                <button
                  onClick={() => invoke("open_url", { url: `http://localhost:${service.config.port}` })}
                  disabled={!isRunning}
                  className="flex items-center gap-2 px-3 py-2 bg-[#2a2a35] hover:bg-[#323240] text-gray-300 rounded-lg text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed ml-2"
                  title="Open in Browser"
                >
                  <ExternalLink className="w-4 h-4" />
                  Open
                </button>
              </div>
              {service.config.executablePath && (
                <div className="text-sm text-gray-400 px-3 py-2 bg-[#1a1b26] border border-[#2a2a35] rounded-lg truncate flex-1 shadow-inner">
                  {service.config.executablePath}
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
