import { useState, useEffect, useCallback, useMemo } from "react";
import { invoke } from "@tauri-apps/api/core";
import { listen } from "@tauri-apps/api/event";
import { useConfirm } from "./components/ConfirmProvider";
import { Titlebar } from "./components/Titlebar";
import { SidebarItem } from "./components/SidebarItem";
import { ServiceCard } from "./components/ServiceCard";
import { LogViewer } from "./components/LogViewer"
import { LanguageContainerView } from "./components/LanguageContainerView";
import { HostsView } from "./components/HostsView";
import { SettingsView } from "./components/SettingsView";
import { MkCertView } from "./components/MkCertView";
import { ServiceEditorView } from "./components/ServiceEditorView";
import { ProjectsView } from "./components/ProjectsView";
import { StacksView } from "./components/StacksView";
import { DockerView } from "./components/DockerView";
import { ErrorBoundary } from "./components/ErrorBoundary";

import { useAppStore } from "./store/useAppStore";
import { parseStatus, toService } from "./utils/transformers";
import { Home, Settings, Power, List, Globe, Shield, Plus, X, Folder, Download, Loader2, Check, AlertCircle, Layers, Box } from "lucide-react";
import type { RawServiceInfo, HostEntry } from "./types";

function AppContent() {
  const { confirm } = useConfirm();
  const { services, hasDocker, setHasDocker, startPolling, stopPolling, refresh } = useAppStore();
  const [activePrimary, setActivePrimary] = useState<string>("home");
  const [activeView, setActiveView] = useState<string>("view-hosts");
  const [loadingIds, setLoadingIds] = useState<Set<string>>(new Set());
  const [editingService, setEditingService] = useState<RawServiceInfo | null>(null);
  const [showExitPopup, setShowExitPopup] = useState(false);
  const [showPullDropdown, setShowPullDropdown] = useState(false);
  const [pullStatus, setPullStatus] = useState<Record<string, "idle" | "pulling" | "success" | "error">>({});
  const [pullLogs, setPullLogs] = useState<Record<string, string>>({});
  const [isPullingAll, setIsPullingAll] = useState(false);

  const dockerImages = useMemo(() => {
    const images = new Set<string>();
    services.forEach(s => {
      if (s.config.runner === "docker" && s.config.executable_path) {
        images.add(s.config.executable_path);
      }
    });
    return Array.from(images);
  }, [services]);

  const handlePullImage = async (image: string) => {
    setPullStatus(prev => ({ ...prev, [image]: "pulling" }));
    try {
      await invoke("pull_docker_image", { image });
      setPullStatus(prev => ({ ...prev, [image]: "success" }));
    } catch (e) {
      console.error(e);
      setPullStatus(prev => ({ ...prev, [image]: "error" }));
    }
  };

  const handlePullAll = async () => {
    setIsPullingAll(true);
    for (const img of dockerImages) {
      if (pullStatus[img] !== "success") {
         await handlePullImage(img);
      }
    }
    setIsPullingAll(false);
  };

  useEffect(() => {
    if (dockerImages.length > 0) {
      invoke<Record<string, boolean>>("check_docker_images", { images: dockerImages })
        .then(result => {
          setPullStatus(prev => {
            const next = { ...prev };
            for (const [img, owned] of Object.entries(result)) {
              if (owned && next[img] !== "pulling") next[img] = "success";
            }
            return next;
          });
        })
        .catch(console.error);
    }
  }, [dockerImages, hasDocker]);

  useEffect(() => {
    let unlisten: (() => void) | undefined;
    listen<{ image: string; line: string }>("docker-pull-log", (event) => {
      setPullLogs((prev) => ({
        ...prev,
        [event.payload.image]: event.payload.line,
      }));
    }).then(u => { unlisten = u; });

    return () => {
      if (unlisten) unlisten();
    };
  }, []);

  useEffect(() => {
    startPolling();
    return () => stopPolling();
  }, [startPolling, stopPolling]);

  const activeService = services.find((s) => s.id === activeView);

  const toggleService = useCallback(
    async (id: string, start: boolean) => {
      setLoadingIds(prev => new Set(prev).add(id));
      try {
        if (start) {
          await invoke("start_service", { id });
        } else {
          await invoke("stop_service", { id });
        }

        // Auto-sync hosts
        try {
          const hostsList = await invoke<HostEntry[]>("get_hosts");
          const svcList = await invoke<RawServiceInfo[]>("get_services");
          const svc = svcList.find(s => s.id === id);
          if (svc) {
            // Find match in host comments
            const host = hostsList.find(h => h.comment.includes(svc.name));
            if (host && host.is_enable !== start) {
              await invoke("toggle_host", { name: host.name, isEnable: start });
            }
          }
        } catch (e) {
          console.warn("Failed to auto-sync host", e);
        }

        await refresh();
      } catch (e) {
        console.error("Toggle failed", e);
      } finally {
        setLoadingIds(prev => {
          const next = new Set(prev);
          next.delete(id);
          return next;
        });
      }
    },
    [refresh],
  );

  const updatePort = useCallback(
    async (id: string, port: number) => {
      try {
        await invoke("set_port", { id, port });
        await refresh();
      } catch (e) {
        console.error(e);
      }
    },
    [refresh],
  );

  const deleteService = useCallback(
    async (id: string) => {
      try {
        await invoke("delete_service", { id });
        setActiveView("view-hosts");
        await refresh();
      } catch (e) {
        console.error(e);
      }
    },
    [refresh],
  );

  const groupedServices = useMemo(() => {
    const web = services.filter(s => s.config.service_type === "Web Server");
    const lang = services.filter(s => s.config.service_type === "Language");
    const db = services.filter(s => s.config.service_type === "Database");
    const admin = services.filter(s => s.config.service_type === "Admin Panel" || (s.config.service_type === "Other" && (s.name === "phpMyAdmin" || s.name === "Adminer")));
    const mailers = services.filter(s => s.config.service_type === "Mailer");
    const storage = services.filter(s => s.config.service_type === "Storage");
    const others = services.filter(s => !["Web Server", "Language", "Database", "Admin Panel", "Mailer", "Storage"].includes(s.config.service_type) && !(s.config.service_type === "Other" && (s.name === "phpMyAdmin" || s.name === "Adminer")));
    
    return { "Web Server": web, "Language & Runtime": lang, "Databases": db, "Storage": storage, "Mailers": mailers, "Admin Panel": admin, "Others": others };
  }, [services]);

  const primaryIcons = [
    { id: "home", icon: Home, label: "Home" },
    { id: "list", icon: List, label: "Services" },
    { id: "projects", icon: Folder, label: "Projects" },
    { id: "stacks", icon: Layers, label: "Stacks" },
    { id: "docker", icon: Box, label: "Docker Stats" },
    { id: "power", icon: Power, label: "Quit App" },
  ];

  return (
    <div className="h-screen w-screen bg-[#1e1e2e] text-gray-200 flex flex-col overflow-hidden font-sans select-none tracking-wide relative">
      <Titlebar>
        {hasDocker && dockerImages.length > 0 && (
          <div className="relative mr-4">
            <button
              onClick={() => setShowPullDropdown(!showPullDropdown)}
              disabled={dockerImages.length === 0}
              className="flex items-center gap-2 px-3 py-1.5 bg-[#2a2a35] hover:bg-[#323240] text-gray-300 rounded-md text-xs font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              title="Pull Docker Images"
            >
              {isPullingAll ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Download className="w-3.5 h-3.5" />}
              Pull Images
            </button>
            {showPullDropdown && (
              <>
                <div className="fixed inset-0 z-[9998]" onMouseDown={() => setShowPullDropdown(false)} />
                <div className="absolute top-full right-0 mt-2 w-80 bg-[#1e1e2e] border border-[#2a2a35] rounded-xl shadow-2xl flex flex-col max-h-[60vh] z-[9999] animate-fade-in">
                  <div className="p-3 border-b border-[#2a2a35] flex items-center justify-between shrink-0 bg-[#181825] rounded-t-xl">
                    <h3 className="font-semibold text-white text-sm">Docker Images</h3>
                    <button
                      onClick={handlePullAll}
                      disabled={isPullingAll || dockerImages.every(img => pullStatus[img] === "success")}
                      className="px-3 py-1.5 bg-blue-500 hover:bg-blue-600 text-white rounded text-xs font-medium transition-colors disabled:opacity-50"
                    >
                      Pull All Missing
                    </button>
                  </div>
                  <div className="overflow-y-auto custom-scrollbar p-2 space-y-1">
                    {dockerImages.map((img) => {
                      const status = pullStatus[img];
                      return (
                        <div key={img} className="flex items-center justify-between p-2.5 rounded-lg hover:bg-[#2a2a35] transition-colors group">
                          <span className="font-mono text-xs text-gray-300 truncate pr-2">{img}</span>
                          <button
                            onClick={() => handlePullImage(img)}
                            disabled={status === "pulling" || status === "success"}
                            className={`shrink-0 flex items-center gap-1.5 px-2.5 py-1.5 rounded text-[11px] font-medium transition-all ${
                              status === "pulling" ? "bg-blue-500/10 text-blue-400" :
                              status === "success" ? "bg-green-500/10 text-green-400" :
                              status === "error" ? "bg-red-500/10 text-red-400 hover:bg-red-500/20" :
                              "bg-[#323240] text-gray-300 hover:bg-blue-500 hover:text-white"
                            }`}
                          >
                            {status === "pulling" ? (
                              <><Loader2 className="w-3.5 h-3.5 animate-spin" /> Pulling</>
                            ) : status === "success" ? (
                              <><Check className="w-3.5 h-3.5" /> Owned</>
                            ) : status === "error" ? (
                              <><AlertCircle className="w-3.5 h-3.5" /> Retry</>
                            ) : (
                              <><Download className="w-3.5 h-3.5" /> Pull</>
                            )}
                          </button>
                        </div>
                      );
                    })}
                  </div>
                </div>
              </>
            )}
          </div>
        )}
      </Titlebar>
      <div className="flex-1 flex overflow-hidden">
        <aside className="w-16 bg-[#181825] border-r border-[#2a2a35] flex flex-col items-center py-4 shrink-0 z-20">
        <div className="flex flex-col gap-6 items-center flex-1">
          {primaryIcons.map(item => (
            <div 
              key={item.id}
              title={item.label}
              onClick={() => {
                if (item.id === "power") {
                  setShowExitPopup(true);
                } else {
                  setActivePrimary(item.id);
                }
              }}
              className={`w-10 h-10 rounded-xl flex items-center justify-center cursor-pointer transition-colors ${
                (activePrimary === item.id && item.id !== "power")
                  ? "bg-blue-500/10 text-blue-400" 
                  : "text-gray-400 hover:text-gray-200 hover:bg-[#2a2a35]/50"
              }`}
            >
              <item.icon className="w-5 h-5" />
            </div>
          ))}
        </div>
        <div 
          onClick={() => setActivePrimary("settings")}
          className={`w-10 h-10 rounded-xl flex items-center justify-center cursor-pointer transition-colors ${
            activePrimary === "settings" 
              ? "bg-blue-500/10 text-blue-400" 
              : "text-gray-400 hover:text-gray-200 hover:bg-[#2a2a35]/50"
          }`}
        >
          <Settings className="w-5 h-5" />
        </div>
      </aside>

      {activePrimary === "list" ? (
        <>
          <aside className="w-64 bg-[#1e1e2e] border-r border-[#2a2a35] flex flex-col shrink-0 z-10 overflow-hidden min-h-0">
            <div className="p-4 border-b border-[#2a2a35] shrink-0 flex justify-between items-center">
              <h2 className="text-xs font-semibold text-gray-400 uppercase tracking-widest">Services</h2>
              <button 
                onClick={() => {
                  setEditingService(null);
                  setActiveView("edit-service");
                }}
                className="text-gray-400 hover:text-white bg-[#2a2a35] hover:bg-[#323240] p-1 rounded transition-colors"
                title="Add Custom Service"
              >
                <Plus className="w-4 h-4" />
              </button>
            </div>
            <div className="flex-1 overflow-y-auto p-3 space-y-6 custom-scrollbar">
              <div>
                <h3 className="text-[11px] font-bold text-gray-500 uppercase tracking-widest mb-2 px-2">Site</h3>
                <div className="space-y-0.5">
                  <div
                    className={`group flex items-center gap-3 w-full p-2.5 rounded-lg transition-colors cursor-pointer ${
                      activeView === "view-hosts" ? "bg-[#2a2a35] text-white" : "text-gray-400 hover:bg-[#252530] hover:text-gray-200"
                    }`}
                    onClick={() => setActiveView("view-hosts")}
                  >
                    <Globe className="w-4 h-4 shrink-0" />
                    <span className="text-sm font-medium">Hosts</span>
                  </div>
                  <div
                    className={`group flex items-center gap-3 w-full p-2.5 rounded-lg transition-colors cursor-pointer ${
                      activeView === "view-mkcert" ? "bg-[#2a2a35] text-white" : "text-gray-400 hover:bg-[#252530] hover:text-gray-200"
                    }`}
                    onClick={() => setActiveView("view-mkcert")}
                  >
                    <Shield className="w-4 h-4 shrink-0" />
                    <span className="text-sm font-medium">MkCert</span>
                  </div>
                </div>
              </div>

              {Object.entries(groupedServices).map(([group, svcs]) => svcs.length > 0 && (
                <div key={group}>
                  <h3 className="text-[11px] font-bold text-gray-500 uppercase tracking-widest mb-2 px-2">{group}</h3>
                  <div className="space-y-0.5">
                    {svcs.map(svc => (
                      <SidebarItem
                        key={svc.id}
                        service={toService(svc)}
                        active={activeView === svc.id}
                        onClick={() => setActiveView(svc.id)}
                        onToggle={toggleService}
                        isLoading={loadingIds.has(svc.id)}
                      />
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </aside>

          <main className="flex-1 flex flex-col bg-[#1e1e2e] overflow-hidden">
            {!hasDocker && (
              <div className="bg-orange-500/10 border-b border-orange-500/20 p-3 px-6 flex items-center justify-between shrink-0">
                <div className="flex items-center gap-3 text-orange-400">
                  <Shield className="w-4 h-4 shrink-0" />
                  <span className="text-sm">Docker Engine is not running or not installed. Docker services will fail to start.</span>
                </div>
                <button onClick={() => setHasDocker(true)} className="text-orange-400/50 hover:text-orange-400">
                  <X className="w-4 h-4" />
                </button>
              </div>
            )}
            
            {activeView === "view-hosts" ? (
              <HostsView />
            ) : activeView === "view-mkcert" ? (
              <MkCertView />
            ) : activeView === "edit-service" ? (
              <ServiceEditorView 
                initialData={editingService ? toService(editingService) : null}
                onClose={() => setActiveView(editingService ? editingService.id : "home")} 
                onSaved={() => {
                  refresh();
                  setActiveView("home");
                }} 
              />
            ) : activeService ? (
              <>
                <header className="h-14 border-b border-[#2a2a35] flex items-center justify-between px-6 shrink-0 bg-[#1e1e2e]/50 backdrop-blur-sm">
                  <div className="flex items-center gap-4">
                    <h1 className="text-sm font-semibold text-white">Service Configuration</h1>
                  </div>
                </header>
                <div className="flex-1 flex flex-col p-6 gap-6 overflow-hidden">
                  <div className="shrink-0">
                    <ServiceCard
                      service={toService(activeService)}
                      adminPanelName={
                        services.find(s => s.status === "Running" && (s.name === "phpMyAdmin" || s.name === "Adminer"))?.name || null
                      }
                      onToggle={() => toggleService(activeService.id, parseStatus(activeService) !== 'running')}
                      onPortChange={(port) => updatePort(activeService.id, port)}
                      onEdit={() => {
                        setEditingService(activeService);
                        setActiveView("edit-service");
                      }}
                      onDelete={async () => {
                        const yes = await confirm({
                          message: "Are you sure you want to remove this service?",
                          title: 'Lumine',
                          kind: 'warning',
                        });
                        if (yes) {
                          deleteService(activeService.id);
                        }
                      }}
                      isLoading={loadingIds.has(activeService.id)}
                    />
                  </div>
                  {activeService.config.service_type === 'Language' ? (
                    <LanguageContainerView 
                      executablePath={activeService.config.executable_path}
                      title={activeService.name}
                    />
                  ) : (
                    <LogViewer
                      log={toService(activeService).log}
                      title={`${activeService.name} Output`}
                      serviceId={activeService.id}
                    />
                  )}
                </div>
              </>
            ) : (
              <div className="flex items-center justify-center h-full text-gray-500">
                Select an item from the sidebar
              </div>
            )}
          </main>
        </>
      ) : activePrimary === "docker" ? (
        <DockerView />
      ) : activePrimary === "home" ? (
        <main className="flex-1 flex flex-col items-center justify-center bg-[#1e1e2e] text-gray-200 p-10 overflow-hidden relative">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-blue-500/10 rounded-full blur-[120px] pointer-events-none animate-pulse" style={{ animationDuration: '4s' }} />
          
          {hasDocker === false ? (
            <div className="z-10 flex flex-col items-center max-w-sm text-center">
              <Shield className="w-12 h-12 text-orange-400/80 mb-4" />
              <h2 className="text-xl font-medium text-white mb-2">Docker is not running</h2>
              <p className="text-sm text-gray-400 mb-8">
                Lumine requires Docker Engine to manage services. Please start Docker Desktop to continue.
              </p>
              
              <div className="flex items-center gap-6 text-sm">
                <button 
                  onClick={() => window.open("https://www.docker.com/products/docker-desktop/")}
                  className="text-gray-500 hover:text-white transition-colors"
                >
                  Download Docker
                </button>
                <button 
                  onClick={async () => {
                    const dockerExists = await invoke<boolean>("check_docker");
                    setHasDocker(dockerExists);
                  }}
                  className="bg-blue-500/10 text-blue-400 hover:bg-blue-500/20 px-5 py-2 rounded-full font-medium transition-colors"
                >
                  Retry Connection
                </button>
              </div>
            </div>
          ) : (
            <>
              <div className="text-center z-10 flex flex-col items-center">
                <h1 className="text-5xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-white to-gray-400 tracking-tight mb-4">Lumine</h1>
                <p className="text-gray-400 text-lg max-w-md mx-auto mb-10 leading-relaxed">
                  Your modern, lightning-fast local web development environment.
                </p>
              </div>
              <div className="absolute bottom-4 right-6 text-sm text-[#8c8d9e] font-medium z-10">
                <button 
                  onClick={() => invoke("open_url", { url: "https://tako.id/ararya" })}
                  className="hover:text-white transition-colors cursor-pointer"
                >
                  Donation
                </button>
                <span className="mx-2 opacity-50">|</span>
                <span>Ararya</span>
              </div>
            </>
          )}
        </main>
      ) : activePrimary === "projects" ? (
        <main className="flex-1 flex flex-col bg-[#1e1e2e] overflow-hidden">
            <ProjectsView />
        </main>
      ) : activePrimary === "stacks" ? (
        <main className="flex-1 flex flex-col bg-[#1e1e2e] overflow-hidden">
            <StacksView />
        </main>
      ) : activePrimary === "settings" ? (
        <SettingsView />
      ) : (
        <main className="flex-1 flex items-center justify-center bg-[#1e1e2e] text-gray-500">
          <div className="flex flex-col items-center gap-4">
            <span className="capitalize">{activePrimary} configuration will appear here.</span>
          </div>
        </main>
      )}

      {showExitPopup && (
        <div 
          className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center animate-fade-in"
          onMouseDown={() => setShowExitPopup(false)}
        >
          <div 
            className="bg-[#242533] border border-[#2a2a35] rounded-2xl p-6 shadow-2xl max-w-sm w-full mx-4 animate-modal"
            onMouseDown={e => e.stopPropagation()}
          >
            <h3 className="text-xl font-bold text-white mb-2">Exit Lumine?</h3>
            <p className="text-gray-400 text-sm mb-6">Are you sure you want to close the application? Running services will be stopped.</p>
            <div className="flex justify-end gap-3">
              <button 
                onClick={() => setShowExitPopup(false)}
                className="px-4 py-2 rounded-lg text-sm font-medium text-gray-300 hover:bg-[#2a2a35] transition-colors"
              >
                Cancel
              </button>
              <button 
                onClick={async () => {
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
                }}
                className="px-4 py-2 rounded-lg text-sm font-medium bg-red-500/10 text-red-400 hover:bg-red-500/20 border border-red-500/20 transition-colors"
              >
                Exit App
              </button>
            </div>
          </div>
        </div>
      )}

      

      {/* Pull Progress Toasts */}
      <div className="fixed bottom-4 right-4 z-[9999] flex flex-col gap-2 pointer-events-none">
        {Object.entries(pullStatus).filter(([, status]) => status === "pulling").map(([image]) => (
          <div key={image} className="bg-[#1e1e2e] border border-[#2a2a35] rounded-lg shadow-xl p-3 w-80 animate-slide-up flex flex-col gap-2 pointer-events-auto">
            <div className="flex items-center gap-2">
              <Loader2 className="w-4 h-4 text-blue-400 animate-spin shrink-0" />
              <span className="text-sm font-semibold text-white truncate flex-1" title={image}>{image}</span>
            </div>
            <div className="text-xs text-gray-400 font-mono truncate">
              {pullLogs[image] || "Initializing pull..."}
            </div>
          </div>
        ))}
      </div>
      </div>
    </div>
  );
}

export default function App() {
  return (
    <ErrorBoundary>
      <AppContent />
    </ErrorBoundary>
  );
}
