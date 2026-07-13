import { useState } from "react";
import { invoke } from "@tauri-apps/api/core";
import { X, ChevronLeft } from "lucide-react";
import type { Service } from "../types";
import { CustomSelect } from "./CustomSelect";
import { Toggle } from "./Toggle";

interface ServiceEditorViewProps {
  onClose: () => void;
  onSaved: () => void;
  initialData?: Service | null;
}

export function ServiceEditorView({ onClose, onSaved, initialData }: ServiceEditorViewProps) {
  const isEdit = !!initialData;
  const [name, setName] = useState(initialData?.name || "");
  const [serviceType, setServiceType] = useState(initialData?.config?.serviceType || "Web Server");
  const [executablePath, setExecutablePath] = useState(initialData?.config?.executablePath || "");
  const [argumentsStr, setArgumentsStr] = useState(initialData?.config?.arguments || "");
  const [port, setPort] = useState<number>(initialData?.config?.port || 0);
  const [containerPort, setContainerPort] = useState<number>(initialData?.config?.containerPort || 0);
  const [volumePath, setVolumePath] = useState(initialData?.config?.volumePath || "");
  const [envPairs, setEnvPairs] = useState<{key: string, value: string}[]>(
    initialData?.config?.env 
      ? Object.entries(initialData.config.env).map(([key, value]) => ({ key, value }))
      : [{ key: "", value: "" }]
  );
  const [autoStart, setAutoStart] = useState(initialData?.config?.auto_start || false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  
  const [availableTags, setAvailableTags] = useState<string[]>([]);
  const [tagsLoading, setTagsLoading] = useState(false);
  const [tagSearch, setTagSearch] = useState("");

  const fetchTags = async () => {
    if (!executablePath) return;
    const baseImage = executablePath.split(':')[0];
    if (!baseImage) return;

    setTagsLoading(true);
    setError("");
    try {
      const tags = await invoke<string[]>("fetch_docker_tags", { imageName: baseImage });
      setAvailableTags(tags);
    } catch (e: any) {
      setError(e.toString());
    } finally {
      setTagsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (!name || !executablePath) {
      setError("Name and Docker Image are required.");
      return;
    }
    setLoading(true);
    
    // Force ports to 0/null if it's a Language service
    const finalPort = serviceType === "Language" ? 0 : Number(port);
    const finalContainerPort = serviceType === "Language" ? null : (containerPort ? Number(containerPort) : null);

    // Parse envPairs to Record
    const envRecord: Record<string, string> = {};
    envPairs.forEach(pair => {
      if (pair.key.trim() !== "") {
        envRecord[pair.key.trim()] = pair.value.trim();
      }
    });

    try {
      if (isEdit && initialData) {
        await invoke("edit_service", {
          id: initialData.id,
          name,
          serviceType,
          executablePath,
          arguments: argumentsStr,
          port: finalPort,
          containerPort: finalContainerPort,
          volumePath: volumePath.trim() || null,
          env: Object.keys(envRecord).length > 0 ? envRecord : null,
          autoStart,
        });
      } else {
        await invoke("add_service", {
          name,
          serviceType,
          executablePath,
          arguments: argumentsStr,
          port: finalPort,
          containerPort: finalContainerPort,
          volumePath: volumePath.trim() || null,
          env: Object.keys(envRecord).length > 0 ? envRecord : null,
          autoStart,
        });
      }
      onSaved();
      onClose();
    } catch (e: any) {
      setError(e.toString());
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex-1 flex flex-col bg-[#1e1e2e] overflow-hidden animate-fade-in">
      <header className="h-14 border-b border-[#2a2a35] flex items-center gap-4 px-6 shrink-0 bg-[#1e1e2e]/50 backdrop-blur-sm">
        <button 
          onClick={onClose}
          type="button"
          className="p-1.5 text-gray-400 hover:text-white hover:bg-[#2a2a35] rounded-lg transition-colors"
        >
          <ChevronLeft className="w-5 h-5" />
        </button>
        <h1 className="text-sm font-semibold text-white">
          {isEdit ? "Edit Service" : "Add Custom Service"}
        </h1>
      </header>

        <form onSubmit={handleSubmit} className="flex flex-col flex-1 min-h-0">
          <div className="p-5 overflow-y-auto overflow-x-hidden custom-scrollbar flex-1 space-y-4">
            {error && <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 text-red-400 rounded-lg text-sm shrink-0">{error}</div>}
          <div>
            <label className="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Name</label>
            <input 
              type="text" 
              placeholder="e.g. PostgreSQL 16"
              className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
              value={name}
              onChange={e => setName(e.target.value)}
            />
          </div>
          <div>
            <label className="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Type</label>
            <CustomSelect
              options={[
                { label: "Web Server", value: "Web Server" },
                { label: "Database", value: "Database" },
                { label: "Language Runtime", value: "Language" },
                { label: "Admin Panel", value: "Admin Panel" },
                { label: "Storage", value: "Storage" },
                { label: "Mailer", value: "Mailer" },
                { label: "Other Service", value: "Other" }
              ]}
              value={serviceType}
              onChange={setServiceType}
            />
          </div>
          <div>
            <label className="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">
              Docker Image
            </label>
            <div className="flex gap-2">
              <input 
                type="text" 
                placeholder="e.g. mariadb:10.11"
                className="flex-1 bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                value={executablePath}
                onChange={e => {
                  setExecutablePath(e.target.value);
                  setAvailableTags([]);
                }}
              />
              <button
                type="button"
                onClick={fetchTags}
                disabled={tagsLoading || !executablePath}
                className="px-3 py-2 bg-[#2a2a35] hover:bg-[#32323d] disabled:opacity-50 text-xs font-medium rounded-lg text-white transition-colors"
              >
                {tagsLoading ? "..." : "Browse"}
              </button>
            </div>
            {availableTags.length > 0 && (
              <div className="mt-2 bg-[#1e1e2e] border border-blue-500 rounded-lg overflow-hidden flex flex-col relative shadow-xl shadow-black/50">
                <div className="p-2 border-b border-[#2a2a35] bg-[#181825]">
                  <input
                    type="text"
                    placeholder="Search version..."
                    className="w-full bg-[#181825] text-xs text-white placeholder-gray-500 focus:outline-none"
                    value={tagSearch}
                    onChange={e => setTagSearch(e.target.value)}
                    autoFocus
                  />
                </div>
                <div className="max-h-48 overflow-y-auto custom-scrollbar p-1">
                  {availableTags
                    .filter(tag => tag.toLowerCase().includes(tagSearch.toLowerCase()))
                    .map(tag => (
                    <button
                      key={tag}
                      type="button"
                      onClick={() => {
                        const baseImage = executablePath.split(':')[0];
                        setExecutablePath(`${baseImage}:${tag}`);
                        setAvailableTags([]);
                        setTagSearch("");
                      }}
                      className="w-full text-left px-3 py-1.5 text-xs text-gray-300 hover:text-white hover:bg-blue-600 rounded transition-colors"
                    >
                      {tag}
                    </button>
                  ))}
                  {availableTags.filter(tag => tag.toLowerCase().includes(tagSearch.toLowerCase())).length === 0 && (
                    <div className="px-3 py-2 text-xs text-gray-500 text-center">No versions match "{tagSearch}"</div>
                  )}
                </div>
              </div>
            )}
          </div>
          <div>
            <label className="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">
              Docker Run Args (Optional)
            </label>
            <input 
              type="text" 
              placeholder="e.g. -e MYSQL_ROOT_PASSWORD=secret"
              className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
              value={argumentsStr}
              onChange={e => setArgumentsStr(e.target.value)}
            />
          </div>
          <div className="flex items-center justify-between p-3 bg-[#181825] border border-[#2a2a35] rounded-lg">
            <div>
              <div className="text-sm font-semibold text-white">Start on System Boot</div>
              <div className="text-xs text-gray-500 mt-0.5">Automatically start this service when Lumine launches.</div>
            </div>
            <Toggle checked={autoStart} onChange={setAutoStart} />
          </div>
          {serviceType !== "Language" && (
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Host Port (Optional)</label>
                <input 
                  type="number" 
                  placeholder="e.g. 8080"
                  className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                  value={port || ""}
                  onChange={e => setPort(Number(e.target.value))}
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Container Port (Optional)</label>
                <input 
                  type="number" 
                  placeholder="e.g. 80"
                  className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                  value={containerPort || ""}
                  onChange={e => setContainerPort(Number(e.target.value))}
                />
              </div>
            </div>
          )}
          {serviceType === "Database" && (
            <div>
              <label className="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">
                Internal Volume Path (Optional)
              </label>
              <input 
                type="text" 
                placeholder="e.g. /var/lib/mysql"
                className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                value={volumePath}
                onChange={e => setVolumePath(e.target.value)}
              />
              <p className="mt-1 text-xs text-gray-500">To make data persistent, specify where this database stores its data internally.</p>
            </div>
          )}

          <div>
            <label className="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5 flex justify-between items-center">
              <span>Environment Variables (Optional)</span>
              <button 
                type="button" 
                onClick={() => setEnvPairs([...envPairs, {key: "", value: ""}])}
                className="text-blue-400 hover:text-blue-300 text-xs normal-case font-medium"
              >
                + Add Variable
              </button>
            </label>
            <div className="space-y-2">
              {envPairs.map((pair, index) => (
                <div key={index} className="flex gap-2">
                  <input 
                    type="text" 
                    placeholder="Key (e.g. MYSQL_ROOT_PASSWORD)"
                    className="flex-1 bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500 font-mono"
                    value={pair.key}
                    onChange={e => {
                      const newPairs = [...envPairs];
                      newPairs[index].key = e.target.value;
                      setEnvPairs(newPairs);
                    }}
                  />
                  <input 
                    type="text" 
                    placeholder="Value (e.g. secret)"
                    className="flex-1 bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500 font-mono"
                    value={pair.value}
                    onChange={e => {
                      const newPairs = [...envPairs];
                      newPairs[index].value = e.target.value;
                      setEnvPairs(newPairs);
                    }}
                  />
                  <button 
                    type="button"
                    onClick={() => {
                      const newPairs = envPairs.filter((_, i) => i !== index);
                      if (newPairs.length === 0) newPairs.push({key: "", value: ""});
                      setEnvPairs(newPairs);
                    }}
                    className="px-2 bg-red-500/10 text-red-400 hover:bg-red-500/20 rounded-lg"
                  >
                    <X className="w-4 h-4" />
                  </button>
                </div>
              ))}
            </div>
          </div>

          
          </div>
          <div className="p-5 border-t border-[#2a2a3c] flex justify-end gap-3 bg-[#181825] rounded-b-2xl shrink-0">
            <button 
              type="button"
              onClick={onClose}
              className="px-5 py-2.5 text-gray-400 hover:text-white transition-colors font-medium"
            >
              Cancel
            </button>
            <button 
              type="submit"
              disabled={loading}
              className="px-5 py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-lg transition-colors font-medium disabled:opacity-50"
            >
              {loading ? "Saving..." : (isEdit ? "Save Changes" : "Add Service")}
            </button>
          </div>
        </form>
    </div>
  );
}
