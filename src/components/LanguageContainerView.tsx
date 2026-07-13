import { useEffect, useState } from "react"
import { Terminal, Trash2, RotateCw, FileCode2, Info } from "lucide-react"
import { invoke } from "@tauri-apps/api/core"

export interface ContainerInfo {
  id: string;
  name: string;
  status: string;
}

interface LanguageContainerViewProps {
  executablePath: string;
  title: string;
}

export function LanguageContainerView({ executablePath, title }: LanguageContainerViewProps) {
  const [containers, setContainers] = useState<ContainerInfo[]>([])
  const [loading, setLoading] = useState(true)

  const fetchContainers = async () => {
    try {
      const data = await invoke<ContainerInfo[]>("get_language_container_status", { executable: executablePath });
      setContainers(data);
    } catch (e) {
      console.error("Failed to fetch containers", e);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchContainers();
    const interval = setInterval(fetchContainers, 2000);
    return () => clearInterval(interval);
  }, [executablePath]);

  const handleStop = async () => {
    try {
      await invoke("stop_language_container", { executable: executablePath });
      fetchContainers();
    } catch (e) {
      console.error("Failed to stop container", e);
    }
  }

  return (
    <div className="bg-[#242533] rounded-xl border border-[#2a2a35] overflow-hidden flex flex-col h-full min-h-[150px]">
      <div className="px-5 py-3 border-b border-[#2a2a35] flex items-center gap-2 bg-[#242533]">
        <Terminal className="w-4 h-4 text-gray-400 shrink-0" />
        <span className="text-sm font-semibold text-gray-300 truncate">{title} Active Containers</span>
        <div className="ml-auto flex items-center gap-2 shrink-0">
          <button 
            onClick={fetchContainers}
            title="Refresh"
            className="text-gray-400 hover:text-white transition-colors bg-[#2a2a35] p-1.5 rounded-md hover:bg-[#323240]"
          >
            <RotateCw className={`w-3.5 h-3.5 ${loading ? 'animate-spin' : ''}`} />
          </button>
        </div>
      </div>
      <div className="flex-1 bg-[#1a1b26] p-4 overflow-auto custom-scrollbar">
        <div className="mb-4 bg-blue-500/10 border border-blue-500/20 rounded-lg p-3 flex gap-3 text-blue-400 items-start">
          <Info className="w-4 h-4 mt-0.5 shrink-0" />
          <div className="text-xs leading-relaxed">
            <span className="font-semibold block mb-1 text-blue-300">Sandbox Environment</span>
            Lingkungan bahasa ini hanyalah <i>sandbox isolasi</i> yang berjalan secara instan. Jika Anda ingin memetakan direktori lokal ke dalam sebuah <i>full environment</i> (seperti Apache/Nginx), silakan gunakan fitur <b>Projects</b>.
          </div>
        </div>
        {containers.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center text-gray-500">
            <Terminal className="w-8 h-8 mb-3 opacity-20" />
            <p className="text-sm font-medium">No active containers</p>
            <p className="text-xs mt-1 opacity-70">Click "Open Terminal" to start one.</p>
          </div>
        ) : (
          <div className="space-y-3">
            {containers.map((c) => (
              <div key={c.id} className="flex items-center justify-between bg-[#242533] border border-[#2a2a35] p-3 rounded-lg">
                <div>
                  <div className="text-sm font-medium text-white mb-1">{c.name}</div>
                  <div className="flex items-center gap-3 text-xs text-gray-400">
                    <span className="font-mono bg-[#1a1b26] px-1.5 py-0.5 rounded text-gray-500">{c.id.substring(0, 12)}</span>
                    <span className="text-green-400">{c.status}</span>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => invoke("open_service_terminal", { executable: executablePath })}
                    title="Open in Terminal"
                    className="flex items-center gap-1.5 px-3 py-1.5 rounded bg-blue-500/10 text-blue-400 hover:bg-blue-500/20 transition-colors text-xs font-medium"
                  >
                    <Terminal className="w-3.5 h-3.5" />
                    Terminal
                  </button>
                  <button
                    onClick={async () => {
                      try {
                        await invoke("open_language_editor", { executable: executablePath });
                      } catch (e) {
                        alert(e);
                      }
                    }}
                    title="Open in VS Code"
                    className="flex items-center gap-1.5 px-3 py-1.5 rounded bg-purple-500/10 text-purple-400 hover:bg-purple-500/20 transition-colors text-xs font-medium"
                  >
                    <FileCode2 className="w-3.5 h-3.5" />
                    Code Editor
                  </button>
                  <button
                    onClick={handleStop}
                    title="Stop and Delete Container"
                    className="flex items-center justify-center w-8 h-8 rounded bg-red-500/10 text-red-400 hover:bg-red-500/20 hover:text-red-300 transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
