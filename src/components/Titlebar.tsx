import { getCurrentWindow } from "@tauri-apps/api/window";
import { invoke } from "@tauri-apps/api/core";
import { Minus, Square, X } from "lucide-react";
import React from "react";

interface TitlebarProps {
  children?: React.ReactNode;
}

export function Titlebar({ children }: TitlebarProps) {
  const appWindow = getCurrentWindow();

  const handleClose = async () => {
    try {
      await invoke("stop_all");
    } catch (e) {
      console.error("Failed to stop services", e);
    }
    try {
      await invoke("exit_app");
    } catch (e) {
      console.error("Failed to exit app", e);
    }
  };

  return (
    <div
      data-tauri-drag-region
      className="relative z-[9999] h-9 flex items-center justify-between select-none bg-[#1e1e2e] border-b border-[#2a2a35] shrink-0"
    >
      {/* Left side: App Menu */}
      <div className="flex items-center h-full pl-4 gap-1">
        <div className="text-sm font-semibold text-gray-200 mr-2 tracking-wider">
          Lumine
        </div>
      </div>

      {/* Right side: Window Controls */}
      <div className="flex h-full items-center">
        {children}
        <button
          onClick={() => appWindow.minimize()}
          className="w-11 h-full flex items-center justify-center text-gray-400 hover:bg-white/10 transition-colors"
          title="Minimize"
        >
          <Minus className="w-3.5 h-3.5" />
        </button>
        <button
          onClick={() => appWindow.toggleMaximize()}
          className="w-11 h-full flex items-center justify-center text-gray-400 hover:bg-white/10 transition-colors"
          title="Maximize"
        >
          <Square className="w-3 h-3" />
        </button>
        <button
          onClick={handleClose}
          className="w-11 h-full flex items-center justify-center text-gray-400 hover:bg-red-500 hover:text-white transition-colors"
          title="Close"
        >
          <X className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
}
