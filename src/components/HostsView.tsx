import { useState, useEffect, useCallback } from "react";
import { invoke } from "@tauri-apps/api/core";
import { useConfirm } from "./ConfirmProvider";
import { Search, Link as LinkIcon, Plus, Trash2, X, Globe, AlertTriangle, Pencil, Shield } from "lucide-react";
import { Toggle } from "./Toggle";
import type { HostEntry, RawServiceInfo, RawProjectInfo, CertEntry } from "../types";

export function HostsView() {
  const { confirm } = useConfirm();
  const [hosts, setHosts] = useState<HostEntry[]>([]);
  const [services, setServices] = useState<RawServiceInfo[]>([]);
  const [projects, setProjects] = useState<RawProjectInfo[]>([]);
  const [certs, setCerts] = useState<CertEntry[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingHost, setEditingHost] = useState<HostEntry | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  // Add form state
  const [newName, setNewName] = useState("");
  const [newComment, setNewComment] = useState("");
  const [isCustomComment, setIsCustomComment] = useState(false);
  const [isServiceDropdownOpen, setIsServiceDropdownOpen] = useState(false);
  const [newEnable, setNewEnable] = useState(true);
  const [newIpv6, setNewIpv6] = useState(false);
  const [addError, setAddError] = useState<string | null>(null);
  const [addLoading, setAddLoading] = useState(false);

  // Context menu

  const fetchHosts = useCallback(async () => {
    try {
      const [hostsResult, servicesResult, projectsResult, certsResult] = await Promise.all([
        invoke<HostEntry[]>("get_hosts"),
        invoke<RawServiceInfo[]>("get_services"),
        invoke<RawProjectInfo[]>("get_projects"),
        invoke<any[]>("get_certs").catch(() => [])
      ]);
      setHosts(hostsResult);
      setServices(servicesResult);
      setProjects(projectsResult);
      setCerts(certsResult);
      setError(null);
    } catch (e) {
      setError(String(e));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchHosts();
    const interval = setInterval(fetchHosts, 3000);
    return () => clearInterval(interval);
  }, [fetchHosts]);

  const filteredHosts = hosts.filter(host =>
    host.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleSave = async () => {
    if (!newName.trim()) {
      setAddError("Hostname is required");
      return;
    }
    setAddLoading(true);
    setAddError(null);
    try {
      if (editingHost) {
        await invoke("edit_host", {
          oldName: editingHost.name,
          name: newName.trim(),
          php: null,
          comment: newComment.trim(),
          isEnable: newEnable,
          isIpv6: newIpv6,
        });
      } else {
        await invoke("add_host", {
          name: newName.trim(),
          php: null,
          comment: newComment.trim(),
          isEnable: newEnable,
          isIpv6: newIpv6,
        });
      }
      setNewName("");
      setNewComment("");
      setNewEnable(true);
      setNewIpv6(false);
      setIsModalOpen(false);
      setEditingHost(null);
      await fetchHosts();
    } catch (e) {
      setAddError(String(e));
    } finally {
      setAddLoading(false);
    }
  };

  const handleToggle = async (name: string, isEnable: boolean) => {
    try {
      await invoke("toggle_host", { name, isEnable });
      
      const host = hosts.find(h => h.name === name);
      if (host && host.comment && isEnable) {
        const svc = services.find(s => host.comment.includes(s.name));
        if (svc) {
          try {
            await invoke("start_service", { id: svc.id });
          } catch (err: unknown) {
            if (!String(err).includes("already running")) {
              throw err;
            }
          }
        } else {
          const proj = projects.find(p => host.comment.includes(p.name));
          if (proj) {
            try {
              await invoke("start_project", { id: proj.id });
            } catch (err: unknown) {
              if (!String(err).includes("already running")) {
                throw err;
              }
            }
          }
        }
      }

      await fetchHosts();
    } catch (e) {
      setError(String(e));
    }
  };

  const handleDelete = async (name: string) => {
    const yes = await confirm({
      message: `Are you sure you want to delete ${name}?`,
      title: "Delete Host",
      kind: "warning",
    });
    if (!yes) return;
    try {
      await invoke("delete_host", { name });
      await fetchHosts();
    } catch (e) {
      setError(String(e));
    }
  };

  const handleGenerateSSL = async (domain: string) => {
    try {
      await invoke("generate_cert", { domain });
      await fetchHosts();
    } catch (e) {
      setError(String(e));
    }
  };

  const handleOpenHostsFile = async () => {
    try {
      await invoke("open_hosts_file");
    } catch (e) {
      setError(String(e));
    }
  };

  return (
    <div className="flex flex-col h-full bg-[#1e1e2e]">
      {/* Toolbar */}
      <div className="p-5 border-b border-[#2a2a35] flex items-center justify-between shrink-0">
        <div className="flex items-center gap-3">
          <button
            onClick={() => {
              setEditingHost(null);
              setNewName("");
              setNewComment("");
              setIsCustomComment(false);
              setNewEnable(true);
              setNewIpv6(false);
              setIsModalOpen(true);
            }}
            className="bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 px-4 py-1.5 rounded-lg border border-blue-500/20 text-sm font-medium transition-all active:scale-95 flex items-center gap-2"
          >
            <Plus className="w-4 h-4" />
            Add Host
          </button>
          <button
            className="bg-[#242533] hover:bg-[#2a2a35] text-gray-300 px-4 py-1.5 rounded-lg border border-[#2a2a35] text-sm font-medium transition-all active:scale-95"
            onClick={handleOpenHostsFile}
          >
            Open hosts file
          </button>
        </div>

        <div className="flex items-center gap-2 text-xs text-gray-500">
          <Globe className="w-4 h-4" />
          <span>{hosts.filter(h => h.is_enable).length} active / {hosts.length} total</span>
        </div>
      </div>

      {/* Error Banner */}
      {error && (
        <div className="mx-5 mt-4 bg-red-500/10 border border-red-500/20 rounded-xl p-3 flex items-center gap-3">
          <AlertTriangle className="w-4 h-4 text-red-400 shrink-0" />
          <span className="text-sm text-red-300 flex-1">{error}</span>
          <button onClick={() => setError(null)} className="text-red-400 hover:text-red-300">
            <X className="w-4 h-4" />
          </button>
        </div>
      )}

      {/* Content */}
      <div className="flex-1 overflow-auto custom-scrollbar p-5">
        {loading ? (
          <div className="flex items-center justify-center h-full text-gray-500">
            <span className="animate-pulse">Loading hosts...</span>
          </div>
        ) : (
          <div className="bg-[#242533] rounded-xl border border-[#2a2a35]">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-[#2a2a35]">
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm w-1/3">
                    <div className="flex items-center gap-4">
                      <span>Site</span>
                      <div className="relative flex-1 max-w-[200px]">
                        <Search className="w-4 h-4 absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-500" />
                        <input
                          type="text"
                          placeholder="Search"
                          value={searchQuery}
                          onChange={(e) => setSearchQuery(e.target.value)}
                          className="w-full bg-[#1a1b26] border border-[#2a2a35] rounded-md pl-9 pr-3 py-1 text-sm text-white focus:outline-none focus:border-blue-500 transition-colors"
                        />
                      </div>
                    </div>
                  </th>
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm w-24">Status</th>
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm">Comment</th>
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm w-16">IPv6</th>
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm text-right w-20">Action</th>
                </tr>
              </thead>
              <tbody>
                {filteredHosts.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="py-12 text-center text-gray-500 text-sm">
                      {hosts.length === 0 ? "No hosts configured yet. Click \"Add Host\" to get started." : "No hosts match your search."}
                    </td>
                  </tr>
                ) : (
                  filteredHosts.map((host) => {
                    const hasCert = certs.some(c => c.domain === host.name);
                    return (
                    <tr key={host.name} className="border-b border-[#2a2a35] hover:bg-[#2a2a35]/50 transition-colors">
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-2 text-gray-300 text-sm">
                          <LinkIcon className={`w-4 h-4 ${host.is_enable ? 'text-green-500' : 'text-gray-600'}`} />
                          {host.is_enable ? (
                            <button
                              onClick={() => invoke("open_url", { url: `${hasCert ? 'https' : 'http'}://${host.name}` })}
                              className="text-blue-400 hover:text-blue-300 hover:underline cursor-pointer text-left"
                            >
                              {host.name}
                            </button>
                          ) : (
                            <span className="text-gray-600 line-through">
                              {host.name}
                            </span>
                          )}
                        </div>
                      </td>
                      <td className="py-3 px-4">
                        <Toggle
                          checked={host.is_enable}
                          onChange={(checked) => handleToggle(host.name, checked)}
                        />
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-400">
                        {host.comment || <span className="text-gray-600 italic">—</span>}
                      </td>
                      <td className="py-3 px-4">
                        <span className={`text-xs font-mono px-2 py-0.5 rounded ${host.is_ipv6 ? 'bg-blue-500/10 text-blue-400' : 'bg-gray-500/10 text-gray-500'}`}>
                          {host.is_ipv6 ? 'Yes' : 'No'}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-right relative">
                          <div className="flex items-center justify-end gap-1">
                            {hasCert ? (
                              <button
                                className="text-green-500 p-1.5 rounded-lg opacity-100 cursor-default"
                                title="Secured with SSL"
                              >
                                <Shield className="w-4 h-4" />
                              </button>
                            ) : (
                              <button
                                onClick={() => handleGenerateSSL(host.name)}
                                className="text-gray-500 hover:text-green-400 hover:bg-[#2a2a35] p-1.5 rounded-lg transition-colors"
                                title="Generate SSL Certificate"
                              >
                                <Shield className="w-4 h-4" />
                              </button>
                            )}
                            <button
                              onClick={() => {
                                setEditingHost(host);
                                setNewName(host.name);
                                setNewComment(host.comment);
                                setNewEnable(host.is_enable);
                                setNewIpv6(host.is_ipv6);
                                setIsCustomComment(!services.some(s => s.name === host.comment));
                                setIsModalOpen(true);
                              }}
                              className="text-gray-400 hover:text-white hover:bg-[#2a2a35] p-1.5 rounded-lg transition-colors"
                            >
                              <Pencil className="w-4 h-4" />
                            </button>
                            <button
                              onClick={() => handleDelete(host.name)}
                              className="text-gray-400 hover:text-red-400 hover:bg-red-500/10 p-1.5 rounded-lg transition-colors"
                            >
                              <Trash2 className="w-4 h-4" />
                            </button>
                          </div>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Add/Edit Host Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center" onClick={() => setIsModalOpen(false)}>
          <div className="bg-[#242533] border border-[#2a2a35] rounded-2xl p-6 shadow-2xl max-w-md w-full mx-4" onClick={e => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-5">
              <h3 className="text-lg font-bold text-white">{editingHost ? "Edit Host Entry" : "Add Host Entry"}</h3>
              <button onClick={() => setIsModalOpen(false)} className="text-gray-400 hover:text-white transition-colors">
                <X className="w-5 h-5" />
              </button>
            </div>

            {addError && (
              <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3 mb-4 text-sm text-red-300">
                {addError}
              </div>
            )}

            <div className="bg-blue-500/10 border border-blue-500/20 rounded-lg p-3 mb-5 text-sm text-blue-300 flex flex-col gap-1">
              <p className="font-medium">How this works:</p>
              <p className="text-blue-300/80">
                This maps the domain directly to your local machine (<code>127.0.0.1</code>). 
                Make sure your Web Server (e.g., Nginx, Caddy, or Apache) is configured to listen for this domain.
              </p>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1.5">Hostname (Domain)</label>
                <input
                  type="text"
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                  placeholder="mysite.local"
                  className="w-full bg-[#1a1b26] border border-[#2a2a35] rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-blue-500 transition-colors"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1.5">Target Service / Project Name</label>
                {!isCustomComment ? (
                  <div className="relative">
                    <button
                      type="button"
                      onClick={() => setIsServiceDropdownOpen(!isServiceDropdownOpen)}
                      className="w-full bg-[#1a1b26] border border-[#2a2a35] rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-blue-500 transition-colors text-left flex items-center justify-between"
                    >
                      <span className={newComment ? "text-white" : "text-gray-500"}>
                        {newComment || "Select a service..."}
                      </span>
                      <svg className="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7"></path></svg>
                    </button>
                    
                    {isServiceDropdownOpen && (
                      <>
                        <div className="fixed inset-0 z-40" onClick={() => setIsServiceDropdownOpen(false)} />
                        <div className="absolute top-full left-0 right-0 mt-1 z-50 bg-[#1e1e2e] border border-[#2a2a35] rounded-lg shadow-xl max-h-60 overflow-y-auto py-1 custom-scrollbar">
                          {services.map(s => (
                            <button
                              key={s.id}
                              type="button"
                              onClick={() => {
                                setNewComment(`${s.name} (port ${s.config.port})`);
                                setIsServiceDropdownOpen(false);
                              }}
                              className="w-full text-left px-3 py-2 text-sm text-gray-300 hover:bg-[#2a2a35] hover:text-white transition-colors"
                            >
                              <span className="text-gray-500 mr-2">[Service]</span>
                              {s.name} (port {s.config.port})
                            </button>
                          ))}
                          {projects.map(p => (
                            <button
                              key={p.id}
                              type="button"
                              onClick={() => {
                                setNewComment(`${p.name} (Project port ${p.port})`);
                                setIsServiceDropdownOpen(false);
                              }}
                              className="w-full text-left px-3 py-2 text-sm text-gray-300 hover:bg-[#2a2a35] hover:text-white transition-colors"
                            >
                              <span className="text-blue-500/80 mr-2">[Project]</span>
                              {p.name} (port {p.port})
                            </button>
                          ))}
                          <div className="h-px bg-[#2a2a35] my-1" />
                          <button
                            type="button"
                            onClick={() => {
                              setIsCustomComment(true);
                              setNewComment("");
                              setIsServiceDropdownOpen(false);
                            }}
                            className="w-full text-left px-3 py-2 text-sm text-gray-300 hover:bg-[#2a2a35] hover:text-white transition-colors"
                          >
                            Custom / Type manually...
                          </button>
                        </div>
                      </>
                    )}
                  </div>
                ) : (
                  <div className="flex gap-2">
                    <input
                      type="text"
                      value={newComment}
                      onChange={(e) => setNewComment(e.target.value)}
                      placeholder="My local project"
                      className="w-full bg-[#1a1b26] border border-[#2a2a35] rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-blue-500 transition-colors"
                      autoFocus
                    />
                    <button
                      type="button"
                      onClick={() => setIsCustomComment(false)}
                      className="shrink-0 px-3 py-2 rounded-lg bg-[#2a2a35] text-gray-300 hover:text-white transition-colors text-sm font-medium"
                    >
                      Back
                    </button>
                  </div>
                )}
              </div>
              <div className="flex items-center gap-6">
                <div className="flex items-center gap-2">
                  <span className="text-sm text-gray-400">Enable:</span>
                  <Toggle checked={newEnable} onChange={setNewEnable} />
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-gray-400">IPv6:</span>
                  <Toggle checked={newIpv6} onChange={setNewIpv6} />
                </div>
              </div>
            </div>

            <div className="flex justify-end gap-3 mt-6">
              <button
                onClick={() => setIsModalOpen(false)}
                className="px-4 py-2 rounded-lg text-sm font-medium text-gray-300 hover:bg-[#2a2a35] transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleSave}
                disabled={addLoading}
                className="px-4 py-2 rounded-lg text-sm font-medium bg-blue-500 hover:bg-blue-600 text-white transition-colors disabled:opacity-50"
              >
                {addLoading ? "Saving..." : (editingHost ? "Save Changes" : "Add Host")}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
