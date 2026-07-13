import { useState, useEffect, useCallback } from "react";
import { invoke } from "@tauri-apps/api/core";
import { Shield, Plus, Trash2, AlertTriangle, CheckCircle, XCircle, X, FileKey } from "lucide-react";
import type { CertEntry, HostEntry } from "../types";
import { CustomSelect } from "./CustomSelect";

export function MkCertView() {
  const [isInstalled, setIsInstalled] = useState<boolean | null>(null);
  const [certs, setCerts] = useState<CertEntry[]>([]);
  const [availableHosts, setAvailableHosts] = useState<HostEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // Generate form
  const [isGenOpen, setIsGenOpen] = useState(false);
  const [domain, setDomain] = useState("");
  const [genLoading, setGenLoading] = useState(false);
  const [genError, setGenError] = useState<string | null>(null);

  // Install CA loading
  const [caLoading, setCaLoading] = useState(false);
  const [caInstalled, setCaInstalled] = useState(() => localStorage.getItem("lumine_mkcert_ca_installed") === "true");

  const [downloadLoading, setDownloadLoading] = useState(false);

  const refresh = useCallback(async () => {
    try {
      const installed = await invoke<boolean>("check_mkcert");
      setIsInstalled(installed);
      if (installed) {
        const list = await invoke<CertEntry[]>("get_certs");
        setCerts(list);
        const hosts = await invoke<HostEntry[]>("get_hosts");
        // Filter out hosts that already have a cert
        const hostsWithoutCert = hosts.filter(h => !list.some(c => c.domain === h.name));
        setAvailableHosts(hostsWithoutCert);
        if (hostsWithoutCert.length > 0) {
          setDomain(hostsWithoutCert[0].name);
        } else {
          setDomain("");
        }
      }
    } catch (e) {
      setError(String(e));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const handleDownloadMkcert = async () => {
    setDownloadLoading(true);
    setError(null);
    try {
      const msg = await invoke<string>("download_mkcert");
      setSuccess(msg);
      setTimeout(() => setSuccess(null), 5000);
      await refresh();
    } catch (e) {
      setError(String(e));
    } finally {
      setDownloadLoading(false);
    }
  };

  const handleInstallCA = async () => {
    setCaLoading(true);
    setError(null);
    try {
      const msg = await invoke<string>("install_root_ca");
      
      // Auto-generate certificates for all currently available hosts
      try {
        const hosts = await invoke<HostEntry[]>("get_hosts");
        for (const host of hosts) {
          await invoke("generate_cert", { domain: host.name });
        }
      } catch (e) {
        console.warn("Failed to auto-generate certs for hosts:", e);
      }
      
      setSuccess(msg + " (and auto-generated default certificates)");
      setCaInstalled(true);
      localStorage.setItem("lumine_mkcert_ca_installed", "true");
      setTimeout(() => setSuccess(null), 5000);
      refresh();
    } catch (e) {
      setError(String(e));
    } finally {
      setCaLoading(false);
    }
  };

  const handleGenerate = async () => {
    if (!domain.trim()) {
      setGenError("Domain is required");
      return;
    }
    setGenLoading(true);
    setGenError(null);
    try {
      await invoke("generate_cert", { domain: domain.trim() });
      setDomain("");
      setIsGenOpen(false);
      await refresh();
      setSuccess(`Certificate generated for ${domain.trim()}`);
      setTimeout(() => setSuccess(null), 5000);
    } catch (e) {
      setGenError(String(e));
    } finally {
      setGenLoading(false);
    }
  };

  const handleDelete = async (certDomain: string) => {
    try {
      await invoke("delete_cert", { domain: certDomain });
      await refresh();
    } catch (e) {
      setError(String(e));
    }
  };

  const handleGenerateAll = async () => {
    setGenLoading(true);
    setGenError(null);
    try {
      let generatedCount = 0;
      for (const host of availableHosts) {
        await invoke("generate_cert", { domain: host.name });
        generatedCount++;
      }
      setSuccess(`Successfully generated ${generatedCount} certificates!`);
      setTimeout(() => setSuccess(null), 3000);
      refresh();
    } catch (e) {
      setGenError(String(e));
    } finally {
      setGenLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full bg-[#1e1e2e] text-gray-500">
        <span className="animate-pulse">Checking mkcert status...</span>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-[#1e1e2e]">
      {/* Toolbar */}
      <div className="p-5 border-b border-[#2a2a35] flex items-center justify-between shrink-0">
        <div className="flex items-center gap-3">
          <button
            onClick={() => setIsGenOpen(true)}
            disabled={!isInstalled}
            className="bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 px-4 py-1.5 rounded-lg border border-blue-500/20 text-sm font-medium transition-all active:scale-95 flex items-center gap-2 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            <Plus className="w-4 h-4" />
            Generate Certificate
          </button>
          <button
            onClick={handleGenerateAll}
            disabled={!isInstalled || availableHosts.length === 0 || genLoading}
            className="bg-[#242533] hover:bg-[#2a2a35] text-gray-300 px-4 py-1.5 rounded-lg border border-[#2a2a35] text-sm font-medium transition-all active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            {genLoading ? "Generating..." : "Auto-Generate Missing"}
          </button>
          <button
            onClick={handleInstallCA}
            disabled={!isInstalled || caLoading || caInstalled}
            className="bg-[#242533] hover:bg-[#2a2a35] text-gray-300 px-4 py-1.5 rounded-lg border border-[#2a2a35] text-sm font-medium transition-all active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            {caLoading ? "Installing..." : caInstalled ? "Root CA Installed" : "Install Root CA"}
          </button>
        </div>
        <div className="flex items-center gap-2 text-xs">
          {isInstalled ? (
            <span className="flex items-center gap-1.5 text-green-400">
              <CheckCircle className="w-4 h-4" />
              mkcert installed
            </span>
          ) : (
            <span className="flex items-center gap-1.5 text-red-400">
              <XCircle className="w-4 h-4" />
              mkcert not found
            </span>
          )}
        </div>
      </div>

      {/* Status Banners */}
      {error && (
        <div className="mx-5 mt-4 bg-red-500/10 border border-red-500/20 rounded-xl p-3 flex items-center gap-3">
          <AlertTriangle className="w-4 h-4 text-red-400 shrink-0" />
          <span className="text-sm text-red-300 flex-1">{error}</span>
          <button onClick={() => setError(null)} className="text-red-400 hover:text-red-300">
            <X className="w-4 h-4" />
          </button>
        </div>
      )}
      {success && (
        <div className="mx-5 mt-4 bg-green-500/10 border border-green-500/20 rounded-xl p-3 flex items-center gap-3">
          <CheckCircle className="w-4 h-4 text-green-400 shrink-0" />
          <span className="text-sm text-green-300 flex-1">{success}</span>
          <button onClick={() => setSuccess(null)} className="text-green-400 hover:text-green-300">
            <X className="w-4 h-4" />
          </button>
        </div>
      )}

      {/* Content */}
      <div className="flex-1 overflow-auto custom-scrollbar p-5">
        {!isInstalled ? (
          <div className="bg-[#242533] rounded-xl border border-[#2a2a35] p-8 text-center">
            <Shield className="w-16 h-16 text-gray-600 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-white mb-2">mkcert Not Installed</h3>
            <p className="text-gray-400 text-sm mb-6 max-w-md mx-auto">
              mkcert is required to generate locally-trusted SSL certificates.
              Lumine can download and set it up automatically.
            </p>
            <button
              onClick={handleDownloadMkcert}
              disabled={downloadLoading}
              className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-2 rounded-lg font-medium transition-all active:scale-95 disabled:opacity-50 flex items-center gap-2 mx-auto"
            >
              {downloadLoading ? (
                <>
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Downloading...
                </>
              ) : (
                "Download mkcert"
              )}
            </button>
          </div>
        ) : certs.length === 0 ? (
          <div className="bg-[#242533] rounded-xl border border-[#2a2a35] p-8 text-center">
            <FileKey className="w-16 h-16 text-gray-600 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-white mb-2">No Certificates Yet</h3>
            <p className="text-gray-400 text-sm">Click "Generate Certificate" to create your first local SSL cert.</p>
          </div>
        ) : (
          <div className="bg-[#242533] rounded-xl border border-[#2a2a35] overflow-hidden">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-[#2a2a35]">
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm">Domain</th>
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm">Certificate Path</th>
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm">Key Path</th>
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm">Created</th>
                  <th className="py-3 px-4 font-medium text-gray-400 text-sm text-right w-20">Action</th>
                </tr>
              </thead>
              <tbody>
                {certs.map((cert) => (
                  <tr key={cert.domain} className="border-b border-[#2a2a35] hover:bg-[#2a2a35]/50 transition-colors">
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2 text-sm">
                        <Shield className="w-4 h-4 text-green-500 shrink-0" />
                        <span className="text-gray-200 font-medium">{cert.domain}</span>
                      </div>
                    </td>
                    <td className="py-3 px-4 text-xs text-gray-500 font-mono truncate max-w-[200px]" title={cert.cert_path}>
                      {cert.cert_path}
                    </td>
                    <td className="py-3 px-4 text-xs text-gray-500 font-mono truncate max-w-[200px]" title={cert.key_path}>
                      {cert.key_path}
                    </td>
                    <td className="py-3 px-4 text-sm text-gray-400">
                      {cert.created_at}
                    </td>
                    <td className="py-3 px-4 text-right">
                      <button
                        onClick={() => handleDelete(cert.domain)}
                        className="text-gray-400 hover:text-red-400 p-1 rounded hover:bg-red-500/10 transition-colors"
                        title="Delete certificate"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Generate Certificate Modal */}
      {isGenOpen && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center" onClick={() => setIsGenOpen(false)}>
          <div className="bg-[#242533] border border-[#2a2a35] rounded-2xl p-6 shadow-2xl max-w-md w-full mx-4" onClick={e => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-5">
              <h3 className="text-lg font-bold text-white">Generate Certificate</h3>
              <button onClick={() => setIsGenOpen(false)} className="text-gray-400 hover:text-white transition-colors">
                <X className="w-5 h-5" />
              </button>
            </div>

            {genError && (
              <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3 mb-4 text-sm text-red-300">
                {genError}
              </div>
            )}

            <div>
              <label className="block text-sm font-medium text-gray-300 mb-1.5">Select Domain from Hosts</label>
              {availableHosts.length > 0 ? (
                <CustomSelect
                  options={availableHosts.map(h => ({ label: h.name, value: h.name }))}
                  value={domain}
                  onChange={setDomain}
                  searchable={availableHosts.length > 5}
                  placeholder="Select a domain..."
                />
              ) : (
                <div className="text-sm text-gray-500 italic p-3 bg-[#1a1b26] rounded-lg border border-[#2a2a35]">
                  All active hosts already have certificates!
                </div>
              )}
              <p className="text-xs text-gray-500 mt-1.5">Only hosts without certificates are shown here.</p>
            </div>

            <div className="flex justify-end gap-3 mt-6">
              <button
                onClick={() => setIsGenOpen(false)}
                className="px-4 py-2 rounded-lg text-sm font-medium text-gray-300 hover:bg-[#2a2a35] transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleGenerate}
                disabled={genLoading}
                className="px-4 py-2 rounded-lg text-sm font-medium bg-blue-500 hover:bg-blue-600 text-white transition-colors disabled:opacity-50"
              >
                {genLoading ? "Generating..." : "Generate"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
