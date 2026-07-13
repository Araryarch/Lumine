import { useState, useRef, useEffect } from "react";
import { X, Play, Square, ExternalLink, Trash2, Folder, Terminal, Edit2, Check, Loader2, Box, Code2, RefreshCw } from "lucide-react";
import type { Project } from "../types";
import { invoke } from "@tauri-apps/api/core";
import { LogViewer } from "./LogViewer";

interface ProjectDetailsModalProps {
  project: Project;
  onClose: () => void;
  onToggle: () => void;
  onReload?: () => void;
  onDelete: () => void;
  onUpdate?: () => void;
  isLoading?: boolean;
}

export function ProjectDetailsSidebar({ project, onClose, onToggle, onReload, onDelete, onUpdate, isLoading }: ProjectDetailsModalProps) {
  const isRunning = project.status === "running";
  
  const [isEditingCmd, setIsEditingCmd] = useState(false);
  const [editedCmd, setEditedCmd] = useState(project.startCommand);
  const [isSavingCmd, setIsSavingCmd] = useState(false);
  const [width, setWidth] = useState(400);
  const isResizing = useRef(false);

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isResizing.current) return;
      const newWidth = document.body.clientWidth - e.clientX;
      setWidth(Math.max(300, Math.min(800, newWidth)));
    };
    const handleMouseUp = () => {
      if (isResizing.current) {
        isResizing.current = false;
        document.body.style.cursor = 'default';
      }
    };
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, []);

  const handleSaveCmd = async () => {
    setIsSavingCmd(true);
    try {
      await invoke("edit_project", { id: project.id, command: editedCmd });
      setIsEditingCmd(false);
      onUpdate?.();
    } catch (e) {
      console.error("Failed to save command", e);
      alert("Failed to save command: " + e);
    } finally {
      setIsSavingCmd(false);
    }
  };

  const handleOpenBrowser = async () => {
    try {
      await invoke("open_url", { url: `http://localhost:${project.port}` });
    } catch (e) {
      console.error("Failed to open browser", e);
    }
  };

  const handleClearLog = async () => {
    try {
      await invoke("clear_project_log", { id: project.id });
      onUpdate?.();
    } catch (e) {
      console.error("Failed to clear log:", e);
    }
  };

  const [showTerminalMenu, setShowTerminalMenu] = useState(false);
  const [customImages, setCustomImages] = useState<{name: string, image: string}[]>([]);

  useEffect(() => {
    if (showTerminalMenu) {
      invoke<any>("get_settings").then(s => {
        if (s.customTerminalImages) {
          setCustomImages(s.customTerminalImages);
        }
      }).catch(console.error);
    }
  }, [showTerminalMenu]);

  const [showAddImageModal, setShowAddImageModal] = useState(false);
  const [newImageName, setNewImageName] = useState("");
  const [newImageRef, setNewImageRef] = useState("");

  const handleSaveCustomImage = async () => {
    if (!newImageName || !newImageRef) return;
    
    try {
      const settings: any = await invoke("get_settings");
      const updatedImages = [...(settings.customTerminalImages || []), { name: newImageName, image: newImageRef }];
      settings.customTerminalImages = updatedImages;
      await invoke("save_settings", { settings });
      setCustomImages(updatedImages);
      setShowAddImageModal(false);
      setNewImageName("");
      setNewImageRef("");
    } catch (e) {
      console.error("Failed to save custom image", e);
      alert(e);
    }
  };

  const handleOpenTerminalWithOverride = async (imageOverride: string | null) => {
    setShowTerminalMenu(false);
    try {
      await invoke("open_project_terminal", { id: project.id, imageOverride });
    } catch (e) {
      console.error("Failed to open terminal", e);
      alert(e);
    }
  };

  return (
    <div style={{ width }} className="absolute right-0 top-0 bottom-0 z-40 bg-[#1e1e2e] border-l border-[#2a2a35] flex flex-col shrink-0 h-full shadow-2xl animate-in slide-in-from-right-8 duration-200">
      
      <div
        className="absolute left-0 top-0 bottom-0 w-1.5 hover:bg-blue-500/50 cursor-col-resize z-10 transition-colors"
        onMouseDown={(e) => {
          e.preventDefault();
          isResizing.current = true;
          document.body.style.cursor = 'col-resize';
        }}
      />

      <div className="flex items-center justify-between p-5 border-b border-[#2a2a35] shrink-0 bg-[#1e1e2e]/50 backdrop-blur">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-blue-500/10 flex items-center justify-center border border-blue-500/20 text-blue-400 shrink-0">
            <Box className="w-5 h-5" />
          </div>
          <div className="min-w-0">
            <h2 className="text-base font-bold text-white flex items-center gap-2">
              <span className="truncate">{project.name}</span>
              {isLoading ? (
                <span className="flex h-1.5 w-1.5 relative shrink-0">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-yellow-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-1.5 w-1.5 bg-yellow-500"></span>
                </span>
              ) : isRunning ? (
                <span className="flex h-1.5 w-1.5 relative shrink-0">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-1.5 w-1.5 bg-green-500"></span>
                </span>
              ) : null}
            </h2>
            <p className="text-gray-400 text-xs mt-0.5 truncate">{project.framework} Project</p>
          </div>
        </div>
        <button onClick={onClose} className="text-gray-400 hover:text-white transition-colors bg-white/5 hover:bg-white/10 p-1.5 rounded-lg shrink-0 ml-2">
          <X className="w-4 h-4" />
        </button>
      </div>

      <div className="flex-1 overflow-y-auto p-5 flex flex-col gap-5 custom-scrollbar">
        <div className="bg-[#242533] border border-[#2a2a35] rounded-xl p-3 flex flex-col gap-2 shrink-0">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5 text-gray-400">
              <Folder className="w-3.5 h-3.5" />
              <span className="text-[10px] font-semibold uppercase tracking-wider">Workspace Path</span>
            </div>
            <div className="flex items-center gap-1.5">
              <button
                onClick={() => invoke("open_in_explorer", { path: project.path })}
                className="text-[10px] flex items-center gap-1 text-gray-400 hover:text-white bg-white/5 hover:bg-white/10 px-1.5 py-0.5 rounded transition-colors"
                title="Open Folder"
              >
                <Folder className="w-3 h-3" /> Show
              </button>
              <button
                onClick={() => invoke("open_in_editor", { path: project.path })}
                className="text-[10px] flex items-center gap-1 text-blue-400 hover:text-blue-300 bg-blue-500/10 hover:bg-blue-500/20 px-1.5 py-0.5 rounded transition-colors"
                title="Open in Editor"
              >
                <Code2 className="w-3 h-3" /> Editor
              </button>
            </div>
          </div>
          <div className="text-xs font-mono text-gray-200 break-all bg-black/20 p-2 rounded-lg border border-[#2a2a35]">
            {project.path}
          </div>
        </div>

        <div className="bg-[#242533] border border-[#2a2a35] rounded-xl p-3 flex flex-col gap-2 shrink-0">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5 text-gray-400">
              <Terminal className="w-3.5 h-3.5" />
              <span className="text-[10px] font-semibold uppercase tracking-wider">Dev Command</span>
            </div>
            {!isRunning && (
              <button
                onClick={() => {
                  if (isEditingCmd) {
                    handleSaveCmd();
                  } else {
                    setEditedCmd(project.startCommand);
                    setIsEditingCmd(true);
                  }
                }}
                disabled={isSavingCmd}
                className="text-[10px] flex items-center gap-1 text-blue-400 hover:text-blue-300 bg-blue-500/10 hover:bg-blue-500/20 px-1.5 py-0.5 rounded transition-colors"
              >
                {isSavingCmd ? (
                  <Loader2 className="w-3 h-3 animate-spin" />
                ) : isEditingCmd ? (
                  <>
                    <Check className="w-3 h-3" /> Save
                  </>
                ) : (
                  <>
                    <Edit2 className="w-3 h-3" /> Edit
                  </>
                )}
              </button>
            )}
          </div>
          {isEditingCmd ? (
            <textarea
              value={editedCmd}
              onChange={(e) => setEditedCmd(e.target.value)}
              className="w-full min-h-[80px] bg-[#1e1e2e] border border-[#2a2a35] rounded-lg p-2 text-xs font-mono text-gray-200 focus:outline-none focus:border-blue-500/50 resize-y custom-scrollbar"
            />
          ) : (
            <div className="text-xs font-mono text-gray-200 break-all bg-black/20 p-2 rounded-lg border border-[#2a2a35]">
              {project.startCommand}
            </div>
          )}
        </div>

        <div className="flex-1 min-h-[250px] flex flex-col pt-2">
          <div className="flex-1 overflow-hidden relative flex flex-col">
            <LogViewer log={project.log} title="Console Output" onClear={handleClearLog} />
          </div>
        </div>
      </div>

      <div className="p-5 border-t border-[#2a2a35] shrink-0 bg-[#242533]/50 flex flex-col gap-3">
        {isRunning ? (
          <div className="flex items-center gap-2">
            <button
              onClick={handleOpenBrowser}
              className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs font-medium text-blue-400 hover:text-blue-300 bg-blue-500/10 hover:bg-blue-500/20 rounded-lg transition-colors border border-blue-500/20"
            >
              <ExternalLink className="w-3.5 h-3.5" />
              Open (:{project.port})
            </button>
            
            <button
              onClick={() => handleOpenTerminalWithOverride(null)}
              title="Open Terminal"
              className="w-10 h-10 shrink-0 flex items-center justify-center text-gray-400 hover:text-white bg-[#1e1e2e] hover:bg-[#2a2a35] rounded-lg transition-colors border border-[#2a2a35]"
            >
              <Terminal className="w-4 h-4" />
            </button>
          </div>
        ) : project.startCommand.includes('"docker"') ? (
          <>
            <button
              onClick={() => setShowTerminalMenu(true)}
              className="w-full flex items-center justify-center gap-2 px-3 py-2 text-xs font-medium text-gray-400 hover:text-white bg-[#1e1e2e] hover:bg-[#2a2a35] rounded-lg transition-colors border border-[#2a2a35]"
            >
              <Terminal className="w-3.5 h-3.5" />
              Open Terminal...
            </button>
            
            {showTerminalMenu && (
              <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in duration-200">
                <div className="bg-[#1e1e2e] border border-[#2a2a35] rounded-xl shadow-2xl w-full max-w-md overflow-hidden flex flex-col animate-in zoom-in-95 duration-200">
                  <div className="flex items-center justify-between p-4 border-b border-[#2a2a35] bg-[#242533]">
                    <h3 className="text-sm font-bold text-white flex items-center gap-2">
                      <Terminal className="w-4 h-4 text-blue-400" />
                      Select Environment
                    </h3>
                    <button onClick={() => setShowTerminalMenu(false)} className="text-gray-400 hover:text-white p-1 rounded-lg hover:bg-white/10 transition-colors">
                      <X className="w-4 h-4" />
                    </button>
                  </div>
                  
                  <div className="p-2 max-h-[60vh] overflow-y-auto custom-scrollbar flex flex-col gap-1">
                    <button onClick={() => handleOpenTerminalWithOverride(null)} className="w-full text-left px-4 py-3 text-sm text-gray-300 hover:bg-[#2a2a35] rounded-lg transition-colors border border-transparent hover:border-[#3a3a45]">
                      <div className="font-semibold text-white mb-0.5">Default Environment</div>
                      <div className="text-xs text-gray-500 font-mono line-clamp-1">{project.startCommand}</div>
                    </button>
                    
                    <div className="my-2 border-t border-[#2a2a35]"></div>
                    
                    <div className="px-3 pb-1 text-[10px] font-semibold text-gray-500 uppercase tracking-wider">Recommended Templates</div>
                    
                    {(() => {
                      const fw = project.framework.toLowerCase();
                      let templates = [];
                      if (["laravel", "symfony", "php"].includes(fw)) {
                        templates = [
                          { name: "Laravel (PHP + Node)", image: "lorisleiva/laravel-docker:8.4" },
                          { name: "PHP / Composer", image: "composer:latest" },
                          { name: "Node.js / NPM", image: "node:lts-alpine" }
                        ];
                      } else if (["nextjs", "react", "vue", "svelte", "nuxt", "nestjs", "express", "node"].includes(fw)) {
                        templates = [
                          { name: "Node.js / NPM", image: "node:lts-alpine" }
                        ];
                      } else if (["fastapi", "django", "flask", "python"].includes(fw)) {
                        templates = [
                          { name: "Python 3.12", image: "python:3.12-alpine" }
                        ];
                      } else if (["gin", "fiber", "go"].includes(fw)) {
                        templates = [
                          { name: "Golang", image: "golang:1.22-alpine" }
                        ];
                      } else if (["elysia", "bun"].includes(fw)) {
                        templates = [
                          { name: "Bun (Oven)", image: "oven/bun:latest" }
                        ];
                      } else {
                        templates = [
                          { name: "Node.js / NPM", image: "node:lts-alpine" },
                          { name: "Python 3.12", image: "python:3.12-alpine" },
                          { name: "PHP / Composer", image: "composer:latest" },
                          { name: "Golang", image: "golang:1.22-alpine" },
                          { name: "Bun", image: "oven/bun:latest" }
                        ];
                      }
                      
                      return templates.map((t, i) => (
                        <button key={i} onClick={() => handleOpenTerminalWithOverride(t.image)} className="w-full text-left px-4 py-2.5 text-sm text-gray-300 hover:bg-[#2a2a35] rounded-lg transition-colors flex justify-between items-center group">
                          <span>{t.name}</span>
                          <span className="text-[10px] text-gray-500 font-mono group-hover:text-blue-400">{t.image}</span>
                        </button>
                      ));
                    })()}
                    
                    {(customImages.length > 0) && (
                      <>
                        <div className="my-2 border-t border-[#2a2a35]"></div>
                        <div className="px-3 pb-1 text-[10px] font-semibold text-gray-500 uppercase tracking-wider">Custom Images</div>
                      </>
                    )}
                    
                    {customImages.map((img, i) => (
                      <button key={i} onClick={() => handleOpenTerminalWithOverride(img.image)} className="w-full text-left px-4 py-2.5 text-sm text-blue-100 hover:bg-[#2a2a35] rounded-lg transition-colors flex justify-between items-center group">
                        <span>{img.name}</span>
                        <span className="text-[10px] text-blue-400/50 font-mono group-hover:text-blue-400">{img.image}</span>
                      </button>
                    ))}
                  </div>
                  
                  <div className="p-4 border-t border-[#2a2a35] bg-[#242533] flex justify-end">
                    <button onClick={() => { setShowTerminalMenu(false); setShowAddImageModal(true); }} className="flex items-center gap-2 px-4 py-2 bg-[#2a2a35] hover:bg-[#3a3a45] text-gray-300 rounded-lg text-xs font-medium transition-colors">
                      <span className="text-lg leading-none mt-[-2px]">+</span> Add Custom Image
                    </button>
                  </div>
                </div>
              </div>
            )}
            
            {showAddImageModal && (
              <div className="fixed inset-0 z-[110] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in duration-200">
                <div className="bg-[#1e1e2e] border border-[#2a2a35] rounded-xl shadow-2xl w-full max-w-sm overflow-hidden flex flex-col animate-in zoom-in-95 duration-200">
                  <div className="flex items-center justify-between p-4 border-b border-[#2a2a35] bg-[#242533]">
                    <h3 className="text-sm font-bold text-white">Add Custom Image</h3>
                    <button onClick={() => setShowAddImageModal(false)} className="text-gray-400 hover:text-white p-1 rounded-lg hover:bg-white/10 transition-colors">
                      <X className="w-4 h-4" />
                    </button>
                  </div>
                  
                  <div className="p-5 flex flex-col gap-4">
                    <div>
                      <label className="block text-xs font-medium text-gray-400 mb-1.5">Image Label</label>
                      <input
                        type="text"
                        value={newImageName}
                        onChange={(e) => setNewImageName(e.target.value)}
                        placeholder="e.g. Golang Environment"
                        className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                        autoFocus
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-gray-400 mb-1.5">Docker Image Reference</label>
                      <input
                        type="text"
                        value={newImageRef}
                        onChange={(e) => setNewImageRef(e.target.value)}
                        placeholder="e.g. golang:latest"
                        className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm font-mono text-white focus:outline-none focus:border-blue-500"
                      />
                    </div>
                  </div>
                  
                  <div className="p-4 border-t border-[#2a2a35] bg-[#242533] flex justify-end gap-2">
                    <button onClick={() => setShowAddImageModal(false)} className="px-4 py-2 text-xs font-medium text-gray-400 hover:text-white transition-colors">
                      Cancel
                    </button>
                    <button onClick={handleSaveCustomImage} disabled={!newImageName || !newImageRef} className="px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:bg-gray-600 disabled:text-gray-400 text-white rounded-lg text-xs font-medium transition-colors">
                      Save Image
                    </button>
                  </div>
                </div>
              </div>
            )}
          </>
        ) : (
          <button
            onClick={() => handleOpenTerminalWithOverride(null)}
            className="w-full flex items-center justify-center gap-2 px-3 py-2 text-xs font-medium text-gray-400 hover:text-white bg-[#1e1e2e] hover:bg-[#2a2a35] rounded-lg transition-colors border border-[#2a2a35]"
          >
            <Terminal className="w-3.5 h-3.5" />
            Open Terminal
          </button>
        )}
        
        <div className="flex items-center gap-2">
          <button
            onClick={onToggle}
            disabled={isLoading}
            className={`flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs font-medium rounded-lg transition-colors border ${
              isLoading
                ? "bg-gray-500/10 text-gray-400 border-gray-500/20 cursor-not-allowed"
                : isRunning
                ? "bg-red-500/10 hover:bg-red-500/20 text-red-400 border-red-500/20"
                : "bg-green-500/10 hover:bg-green-500/20 text-green-400 border-green-500/20"
            }`}
          >
            {isLoading ? (
              <span className="flex items-center justify-center gap-1.5">
                <div className="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin" />
                Working...
              </span>
            ) : isRunning ? (
              <>
                <Square className="w-3.5 h-3.5" />
                Stop Server
              </>
            ) : (
              <>
                <Play className="w-3.5 h-3.5" />
                Start Server
              </>
            )}
          </button>

          {isRunning && (
            <button
              onClick={onReload}
              className="w-10 h-10 flex items-center justify-center shrink-0 rounded-lg text-yellow-500 bg-yellow-500/10 border border-yellow-500/20 hover:bg-yellow-500/20 transition-colors"
              title="Reload Server"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
          )}
        </div>
        
        <button
          onClick={onDelete}
          className="flex items-center justify-center gap-2 px-3 py-2 text-xs font-medium text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-lg transition-colors mt-2"
        >
          <Trash2 className="w-3.5 h-3.5" />
          Delete Project
        </button>
      </div>
    </div>
  );
}
