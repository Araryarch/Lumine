import { Play, Square, RefreshCw, ChevronRight, Box, Globe, Copy, Check } from "lucide-react";
import { useState, useEffect } from "react";
import { invoke } from "@tauri-apps/api/core";
import type { Project } from "../types";

interface ProjectCardProps {
  project: Project;
  onClick: () => void;
  onToggle: (e: React.MouseEvent) => void;
  isLoading?: boolean;
}

export function ProjectCard({ project, onClick, onToggle, isLoading }: ProjectCardProps) {
  const isRunning = project.status === "running";
  
  const [isTunneling, setIsTunneling] = useState(false);
  const [tunnelUrl, setTunnelUrl] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    let interval: number;
    if (isRunning && isTunneling && !tunnelUrl) {
      interval = window.setInterval(async () => {
        try {
          const url = await invoke<string | null>("get_tunnel_url", { id: project.id });
          if (url) setTunnelUrl(url);
        } catch (e) {
          console.error("Tunnel check failed:", e);
        }
      }, 2000);
    }
    return () => clearInterval(interval);
  }, [isRunning, isTunneling, tunnelUrl, project.id]);

  useEffect(() => {
    if (!isRunning) {
      setIsTunneling(false);
      setTunnelUrl(null);
    }
  }, [isRunning]);

  const toggleTunnel = async (e: React.MouseEvent) => {
    e.stopPropagation();
    if (isTunneling || tunnelUrl) {
      setIsTunneling(false);
      setTunnelUrl(null);
      await invoke("stop_tunnel", { id: project.id }).catch(console.error);
    } else {
      setIsTunneling(true);
      await invoke("start_tunnel", { id: project.id, port: project.port }).catch(console.error);
    }
  };

  const copyUrl = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (tunnelUrl) {
      navigator.clipboard.writeText(tunnelUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div 
      className="bg-[#242533]/80 hover:bg-[#2a2b3d] border border-[#2a2a35] rounded-2xl p-5 cursor-pointer transition-all flex items-center justify-between group"
      onClick={onClick}
    >
      <div className="flex items-center gap-4">
        <div className={`w-12 h-12 rounded-xl flex items-center justify-center transition-colors ${
          isRunning 
            ? "bg-green-500/10 text-green-400 border border-green-500/20 shadow-[0_0_15px_rgba(34,197,94,0.1)]" 
            : "bg-blue-500/10 text-blue-400 border border-blue-500/20 group-hover:border-blue-500/40"
        }`}>
          <Box className="w-6 h-6" />
        </div>
        <div>
          <h2 className="text-base font-semibold text-white flex items-center gap-2">
            {project.name}
            {isLoading ? (
              <span className="flex h-2 w-2 relative">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-yellow-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-yellow-500 shadow-[0_0_8px_rgba(234,179,8,0.6)]"></span>
              </span>
            ) : isRunning ? (
              <span className="flex h-2 w-2 relative">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.6)]"></span>
              </span>
            ) : null}
          </h2>
          <p className="text-gray-400 text-xs mt-1 font-medium tracking-wide">
            {project.framework} • PORT {project.port}
          </p>
        </div>
      </div>

      <div className="flex items-center gap-3">
        {isRunning && (
          <div className="flex items-center gap-2 mr-2">
            {tunnelUrl ? (
              <button
                onClick={copyUrl}
                className="flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-medium bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20 border border-emerald-500/20 transition-all"
                title="Copy Live Share URL"
              >
                {copied ? <Check className="w-3.5 h-3.5" /> : <Copy className="w-3.5 h-3.5" />}
                {copied ? "Copied!" : "Live"}
              </button>
            ) : null}
            <button
              onClick={toggleTunnel}
              className={`w-9 h-9 rounded-lg flex items-center justify-center transition-all border ${
                isTunneling && !tunnelUrl
                  ? "bg-purple-500/10 text-purple-400 border-purple-500/20 animate-pulse"
                  : tunnelUrl
                  ? "bg-purple-500/10 text-purple-400 border-purple-500/20 hover:bg-purple-500/20"
                  : "bg-[#181825] text-gray-400 border-[#2a2a35] hover:text-purple-400 hover:border-purple-500/30"
              }`}
              title={tunnelUrl ? "Stop Live Share" : "Start Live Share"}
            >
              <Globe className="w-4 h-4" />
            </button>
          </div>
        )}
        <button
          onClick={onToggle}
          disabled={isLoading}
          className={`w-10 h-10 rounded-xl flex items-center justify-center transition-all ${
            isLoading
              ? "bg-gray-500/10 text-gray-500 cursor-not-allowed border border-gray-500/20"
              : isRunning
              ? "bg-red-500/10 text-red-400 hover:bg-red-500/20 border border-red-500/20"
              : "bg-green-500/10 text-green-400 hover:bg-green-500/20 border border-green-500/20"
          }`}
        >
          {isLoading ? (
            <RefreshCw className="w-4 h-4 animate-spin" />
          ) : isRunning ? (
            <Square className="w-4 h-4" />
          ) : (
            <Play className="w-4 h-4 ml-0.5" />
          )}
        </button>
        <div className="w-8 h-8 flex items-center justify-center text-gray-500 opacity-0 group-hover:opacity-100 transition-opacity">
          <ChevronRight className="w-5 h-5" />
        </div>
      </div>
    </div>
  );
}
