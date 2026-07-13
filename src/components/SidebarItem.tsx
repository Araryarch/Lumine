import { Toggle } from "./Toggle"
import type { Service, ServiceStatus } from "../types"
import { Server, Database, Container, Box, Terminal, Mail, HardDrive } from "lucide-react"
import { invoke } from "@tauri-apps/api/core"

interface SidebarItemProps {
  service: Service
  active: boolean
  onClick: () => void
  onToggle: (id: string, start: boolean) => void
  isLoading?: boolean
}

function getIconForService(name: string) {
  const n = name.toLowerCase();
  if (n.includes('mysql') || n.includes('postgres') || n.includes('redis') || n.includes('mongo') || n.includes('mariadb')) return Database;
  if (n.includes('docker') || n.includes('container')) return Container;
  if (n.includes('apache') || n.includes('nginx') || n.includes('caddy') || n.includes('franken')) return Server;
  if (n.includes('mail')) return Mail;
  if (n.includes('ftp') || n.includes('minio') || n.includes('storage')) return HardDrive;
  return Box;
}

function getStatusColor(status: ServiceStatus, active: boolean) {
  switch (status) {
    case "running": return "text-green-400";
    case "error": return "text-red-400";
    case "starting":
    case "stopping": return "text-yellow-400 animate-pulse";
    default: return active ? "text-white" : "text-gray-400 group-hover:text-gray-200";
  }
}

export function SidebarItem({ service, active, onClick, onToggle, isLoading }: SidebarItemProps) {
  const Icon = getIconForService(service.name);
  const isRunning = service.status === "running";
  const statusColor = getStatusColor(service.status, active);
  
  return (
    <div
      className={`group flex items-center justify-between w-full p-2.5 rounded-lg transition-colors cursor-pointer
        ${active ? "bg-[#2a2a35]" : "hover:bg-[#252530]"}
      `}
      onClick={onClick}
    >
      <div className={`flex items-center gap-3 overflow-hidden ${statusColor} transition-colors`}>
        <Icon className="w-4 h-4 shrink-0" />
        <span className="text-sm font-medium truncate">{service.name}</span>
      </div>
      
      <div className="shrink-0 flex items-center ml-2" onClick={e => e.stopPropagation()}>
        {service.config.serviceType === "Language" ? (
          <button 
            onClick={() => invoke("open_service_terminal", { executable: service.config.executablePath })}
            className="p-1.5 rounded-md text-gray-400 hover:text-white hover:bg-[#323240] transition-colors"
            title={`Open Interactive ${service.name} Terminal`}
          >
            <Terminal className="w-4 h-4" />
          </button>
        ) : (
          <Toggle 
            checked={isRunning} 
            onChange={(checked) => onToggle(service.id, checked)}
            disabled={isLoading || service.status === "starting" || service.status === "stopping"} 
          />
        )}
      </div>
    </div>
  )
}
