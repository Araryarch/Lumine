import { useState, useEffect } from "react";
import { invoke } from "@tauri-apps/api/core";
import { useConfirm } from "./ConfirmProvider";
import { Monitor, Folder, Settings, Shield, DownloadCloud } from "lucide-react";
import { Toggle } from "./Toggle";
import { CustomSelect } from "./CustomSelect";
import type { AppSettings } from "../types";

export function SettingsView() {
  const { confirm } = useConfirm();
  const [activeTab, setActiveTab] = useState("general");
  const [caLoading, setCaLoading] = useState(false);
  const [caMessage, setCaMessage] = useState<string | null>(null);
  const [availableTools, setAvailableTools] = useState<{terminals: string[], editors: string[], explorers: string[]}>({ terminals: ["cmd"], editors: ["notepad"], explorers: ["explorer"] });
  const [settings, setSettings] = useState({
    startOnBoot: false,
    autoStartServices: true,
    minimizeToTray: true,
    checkUpdates: true,
    defaultPhp: "8.3",
    defaultNode: "20.11",
    documentRoot: "C:\\Lumine\\www",
    terminalEmulator: "cmd",
    codeEditor: "code",
    fileExplorer: "explorer",
  });

  useEffect(() => {
    async function loadSettings() {
      try {
        const data = await invoke<typeof settings>("get_settings");
        setSettings(data);
        
        // Also fetch available tools on mount
        const tools = await invoke<{terminals: string[], editors: string[], explorers: string[]}>("get_available_tools");
        setAvailableTools(tools);
      } catch (e) {
        console.error("Failed to load settings:", e);
      }
    }
    loadSettings();
  }, []);

  const updateSetting = async (key: keyof AppSettings, value: string | number | boolean) => {
    const next = { ...settings, [key]: value };
    setSettings(next as any); // Optimistic UI
    try {
      await invoke("save_settings", { settings: next });
    } catch (e) {
      console.error(e);
      // loadSettings(); // revert
    }
  };

  const tabs = [
    { id: "general", label: "General", icon: Monitor },
    { id: "paths", label: "Directories", icon: Folder },
    { id: "preferences", label: "Preferences", icon: Settings },
    { id: "data", label: "Data & Recovery", icon: DownloadCloud },
  ];

  const terminalLabels: Record<string, string> = {
    cmd: "Command Prompt",
    powershell: "PowerShell",
    wt: "Windows Terminal",
    bash: "Git Bash",
    zsh: "Zsh",
    fish: "Fish",
    "gnome-terminal": "GNOME Terminal",
    konsole: "Konsole",
    "xfce4-terminal": "Xfce Terminal",
    terminator: "Terminator",
    alacritty: "Alacritty",
    kitty: "Kitty",
    wezterm: "WezTerm"
  };

  const editorLabels: Record<string, string> = {
    code: "VS Code",
    cursor: "Cursor",
    zed: "Zed",
    webstorm: "WebStorm",
    phpstorm: "PhpStorm",
    subl: "Sublime Text",
    notepad: "Notepad",
    nvim: "Neovim",
    vim: "Vim",
    nano: "Nano",
    hx: "Helix"
  };

  const explorerLabels: Record<string, string> = {
    explorer: "Windows Explorer",
    dolphin: "Dolphin",
    thunar: "Thunar",
    nautilus: "Nautilus",
    nemo: "Nemo",
    pcmanfm: "PCManFM",
    "xdg-open": "xdg-open (Linux Default)",
    open: "Finder (Mac)"
  };

  return (
    <div className="flex-1 flex flex-col bg-[#1e1e2e] overflow-hidden text-gray-200">
      <header className="h-14 border-b border-[#2a2a35] flex items-center px-6 shrink-0 bg-[#1e1e2e]/50 backdrop-blur-sm">
        <h1 className="text-sm font-semibold text-white">Settings</h1>
      </header>
      
      <div className="flex flex-1 overflow-hidden">
        {/* Settings Sidebar */}
        <aside className="w-56 border-r border-[#2a2a35] p-4 flex flex-col gap-1 overflow-y-auto custom-scrollbar shrink-0">
          {tabs.map((tab) => (
            <div
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-lg cursor-pointer transition-colors ${
                activeTab === tab.id
                  ? "bg-[#2a2a35] text-white"
                  : "text-gray-400 hover:text-gray-200 hover:bg-[#252530]"
              }`}
            >
              <tab.icon className="w-4 h-4 shrink-0" />
              <span className="text-sm font-medium">{tab.label}</span>
            </div>
          ))}
        </aside>

        {/* Settings Content */}
        <div className="flex-1 overflow-y-auto p-8 custom-scrollbar">
          <div className="max-w-2xl">
            {activeTab === "general" && (
              <div className="space-y-8 animate-in fade-in slide-in-from-bottom-2 duration-300">
                <div>
                  <h2 className="text-lg font-semibold text-white mb-4">Application Behavior</h2>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between p-4 bg-[#242533] border border-[#2a2a35] rounded-xl">
                      <div>
                        <h4 className="text-sm font-medium text-white">Start on Boot</h4>
                        <p className="text-xs text-gray-400 mt-1">Automatically launch Lumine when your computer starts.</p>
                      </div>
                      <Toggle checked={settings.startOnBoot} onChange={(v) => updateSetting("startOnBoot", v)} />
                    </div>
                    
                    <div className="flex items-center justify-between p-4 bg-[#242533] border border-[#2a2a35] rounded-xl">
                      <div>
                        <h4 className="text-sm font-medium text-white">Auto-start Services</h4>
                        <p className="text-xs text-gray-400 mt-1">Start all active services automatically when Lumine is opened.</p>
                      </div>
                      <Toggle checked={settings.autoStartServices} onChange={(v) => updateSetting("autoStartServices", v)} />
                    </div>

                    <div className="flex items-center justify-between p-4 bg-[#242533] border border-[#2a2a35] rounded-xl">
                      <div>
                        <h4 className="text-sm font-medium text-white">Minimize to System Tray</h4>
                        <p className="text-xs text-gray-400 mt-1">Keep Lumine running in the background when the window is closed.</p>
                      </div>
                      <Toggle checked={settings.minimizeToTray} onChange={(v) => updateSetting("minimizeToTray", v)} />
                    </div>
                  </div>
                </div>
                
                <div>
                  <h2 className="text-lg font-semibold text-white mb-4">Updates</h2>
                  <div className="flex items-center justify-between p-4 bg-[#242533] border border-[#2a2a35] rounded-xl">
                    <div>
                      <h4 className="text-sm font-medium text-white">Check for Updates</h4>
                      <p className="text-xs text-gray-400 mt-1">Automatically check for application updates on startup.</p>
                    </div>
                    <div className="flex items-center gap-4">
                      <button className="text-xs font-medium text-blue-400 hover:text-blue-300 flex items-center gap-1.5">
                        <DownloadCloud className="w-3.5 h-3.5" />
                        Check Now
                      </button>
                      <Toggle checked={settings.checkUpdates} onChange={(v) => updateSetting("checkUpdates", v)} />
                    </div>
                  </div>
                </div>

                <div>
                  <h2 className="text-lg font-semibold text-white mb-4">Security</h2>
                  <div className="flex flex-col items-center justify-center p-8 bg-[#242533] border border-[#2a2a35] rounded-xl text-center">
                    <Shield className="w-12 h-12 text-gray-500 mb-4" />
                    <h3 className="text-white font-medium mb-1">Local Certificates (MkCert)</h3>
                    <p className="text-gray-400 text-sm max-w-sm mb-6">Lumine uses MkCert to automatically provision valid SSL certificates for your local domains.</p>
                    {caMessage && (
                      <div className="mb-4 text-sm px-3 py-1.5 rounded-lg bg-blue-500/10 text-blue-300 border border-blue-500/20">
                        {caMessage}
                      </div>
                    )}
                    <button 
                      onClick={async () => {
                        setCaLoading(true);
                        setCaMessage(null);
                        try {
                            const msg = await invoke<string>("install_root_ca");
                            setCaMessage(msg);
                        } catch(e) {
                            setCaMessage(String(e));
                        } finally {
                            setCaLoading(false);
                            setTimeout(() => setCaMessage(null), 4000);
                        }
                      }}
                      disabled={caLoading}
                      className="px-5 py-2 bg-blue-500/10 text-blue-400 hover:bg-blue-500/20 border border-blue-500/20 rounded-lg text-sm font-medium transition-colors disabled:opacity-50"
                    >
                      {caLoading ? "Installing..." : "Re-install Root CA"}
                    </button>
                  </div>
                </div>
              </div>
            )}
            
            {activeTab === "paths" && (
              <div className="space-y-8 animate-in fade-in slide-in-from-bottom-2 duration-300">
                <div>
                  <h2 className="text-lg font-semibold text-white mb-4">Directories</h2>
                  <div className="space-y-4">
                    <div className="flex flex-col gap-2 p-4 bg-[#242533] border border-[#2a2a35] rounded-xl">
                      <label className="text-sm font-medium text-white">Default Document Root</label>
                      <p className="text-xs text-gray-400 mb-2">The primary directory where your local websites are stored.</p>
                      <div className="flex items-center gap-3">
                        <input 
                          type="text" 
                          value={settings.documentRoot}
                          onChange={(e) => updateSetting("documentRoot", e.target.value)}
                          className="flex-1 bg-[#1a1b26] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                        />
                        <button 
                          onClick={async () => {
                            try {
                              const { open } = await import("@tauri-apps/plugin-dialog");
                              const selected = await open({
                                directory: true,
                                multiple: false,
                                title: "Select Document Root",
                              });
                              if (selected && typeof selected === "string") {
                                updateSetting("documentRoot", selected);
                              }
                            } catch (e) {
                              console.error("Failed to open directory dialog", e);
                            }
                          }}
                          className="px-4 py-2 bg-[#2a2a35] hover:bg-[#323240] rounded-lg text-sm font-medium transition-colors"
                        >
                          Browse
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {activeTab === "preferences" && (
              <div className="space-y-8 animate-in fade-in slide-in-from-bottom-2 duration-300">
                <div>
                  <h2 className="text-lg font-semibold text-white mb-4">Development Tools</h2>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between p-4 bg-[#242533] border border-[#2a2a35] rounded-xl">
                      <div>
                        <h4 className="text-sm font-medium text-white">Terminal Emulator</h4>
                        <p className="text-xs text-gray-400 mt-1">Default terminal used to open project shells.</p>
                      </div>
                      <div className="w-48">
                        <CustomSelect 
                          allowCustom
                          value={settings.terminalEmulator} 
                          onChange={(v) => updateSetting("terminalEmulator", v)}
                          options={availableTools.terminals.map(t => ({
                            label: terminalLabels[t] || t,
                            value: t
                          }))} 
                        />
                      </div>
                    </div>
                    
                    <div className="flex items-center justify-between p-4 bg-[#242533] border border-[#2a2a35] rounded-xl">
                      <div>
                        <h4 className="text-sm font-medium text-white">Code Editor</h4>
                        <p className="text-xs text-gray-400 mt-1">Default editor to open your projects with.</p>
                      </div>
                      <div className="w-48">
                        <CustomSelect 
                          allowCustom
                          value={settings.codeEditor} 
                          onChange={(v) => updateSetting("codeEditor", v)}
                          options={availableTools.editors.map(e => ({
                            label: editorLabels[e] || e,
                            value: e
                          }))} 
                        />
                      </div>
                    </div>

                    <div className="flex items-center justify-between p-4 bg-[#242533] border border-[#2a2a35] rounded-xl">
                      <div>
                        <h4 className="text-sm font-medium text-white">File Explorer</h4>
                        <p className="text-xs text-gray-400 mt-1">Default file manager to open project paths.</p>
                      </div>
                      <div className="w-48">
                        <CustomSelect 
                          allowCustom
                          value={settings.fileExplorer} 
                          onChange={(v) => updateSetting("fileExplorer", v)}
                          options={availableTools.explorers.map(e => ({
                            label: explorerLabels[e] || e,
                            value: e
                          }))} 
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {activeTab === "data" && (
              <div className="space-y-8 animate-in fade-in slide-in-from-bottom-2 duration-300">
                <div>
                  <h2 className="text-lg font-semibold text-white mb-4">Export & Import</h2>
                  <div className="flex items-center justify-between p-4 bg-[#242533] border border-[#2a2a35] rounded-xl mb-4">
                    <div>
                      <h4 className="text-sm font-medium text-white">Export Configuration</h4>
                      <p className="text-xs text-gray-400 mt-1">Backup your services, projects, and settings to a JSON file.</p>
                    </div>
                    <button 
                      onClick={async () => {
                        try {
                          const jsonData = await invoke<string>("export_config");
                          const { save, message } = await import("@tauri-apps/plugin-dialog");
                          const { writeTextFile } = await import("@tauri-apps/plugin-fs");
                          const path = await save({ defaultPath: "lumine-backup.json" });
                          if (path) {
                            await writeTextFile(path, jsonData);
                            await message("Exported successfully!", { title: "Lumine", kind: "info" });
                          }
                        } catch (e) {
                          console.error(e);
                        }
                      }}
                      className="px-4 py-2 bg-blue-500/10 text-blue-400 hover:bg-blue-500/20 border border-blue-500/20 rounded-lg text-sm font-medium transition-colors"
                    >
                      Export Backup
                    </button>
                  </div>
                  <div className="flex items-center justify-between p-4 bg-[#242533] border border-[#2a2a35] rounded-xl">
                    <div>
                      <h4 className="text-sm font-medium text-white">Import Configuration</h4>
                      <p className="text-xs text-gray-400 mt-1">Restore your setup from a previously saved JSON backup.</p>
                    </div>
                    <button 
                      onClick={async () => {
                        try {
                          const { open, message } = await import("@tauri-apps/plugin-dialog");
                          const { readTextFile } = await import("@tauri-apps/plugin-fs");
                          const path = await open({ filters: [{ name: 'JSON', extensions: ['json'] }] });
                          if (path) {
                            const data = await readTextFile(path);
                            const yes = await confirm({ message: "Are you sure? This will overwrite your current configuration and restart all services.", title: "Lumine", kind: "warning" });
                            if (yes) {
                              await invoke("import_config", { jsonData: data });
                              await message("Imported successfully! Lumine will now reload.", { title: "Lumine", kind: "info" });
                              window.location.reload();
                            }
                          }
                        } catch (e) {
                          console.error(e);
                        }
                      }}
                      className="px-4 py-2 bg-blue-500/10 text-blue-400 hover:bg-blue-500/20 border border-blue-500/20 rounded-lg text-sm font-medium transition-colors"
                    >
                      Import Backup
                    </button>
                  </div>
                </div>

                <div>
                  <h2 className="text-lg font-semibold text-red-500 mb-4">Danger Zone</h2>
                  <div className="flex items-center justify-between p-4 bg-red-500/5 border border-red-500/20 rounded-xl">
                    <div>
                      <h4 className="text-sm font-medium text-red-400">Factory Reset</h4>
                      <p className="text-xs text-red-400/70 mt-1">Permanently delete all Lumine configurations and quit.</p>
                    </div>
                    <button 
                      onClick={async () => {
                        try {
                          const yes = await confirm({ message: "WARNING: This will completely delete all Lumine configurations, stop all services, and exit the app. This cannot be undone. Are you absolutely sure?", title: "Lumine Factory Reset", kind: "warning" });
                          if (yes) {
                            await invoke("factory_reset");
                          }
                        } catch (e) {
                          console.error(e);
                        }
                      }}
                      className="px-4 py-2 bg-red-500 text-white hover:bg-red-600 rounded-lg text-sm font-medium transition-colors"
                    >
                      Factory Reset
                    </button>
                  </div>
                </div>
              </div>
            )}

          </div>
        </div>
      </div>
    </div>
  );
}
