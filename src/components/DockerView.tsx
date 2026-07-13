import { useState, useEffect } from 'react';
import { invoke } from '@tauri-apps/api/core';
import { RefreshCw, Box, Square, Trash2 } from 'lucide-react';

interface DockerStat {
  cpu: string;
  ram: string;
  full_name: string;
}

function parsePercent(val: string): number {
  const n = parseFloat(val.replace('%', ''));
  return isNaN(n) ? 0 : Math.min(n, 100);
}

export function DockerView() {
  const [stats, setStats] = useState<Record<string, DockerStat>>({});
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  const fetchStats = async () => {
    try {
      const data = await invoke<Record<string, DockerStat>>('get_docker_stats');
      setStats(data);
    } catch (e) {
      console.error("Failed to fetch docker stats:", e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
    const interval = setInterval(fetchStats, 5000);
    return () => clearInterval(interval);
  }, []);

  const containers = Object.keys(stats);

  const handleStop = async (name: string) => {
    const fullName = stats[name]?.full_name || name;
    setActionLoading(name);
    try {
      await invoke('stop_container', { name: fullName });
      await fetchStats();
    } catch (e) {
      console.error("Failed to stop container:", e);
    } finally {
      setActionLoading(null);
    }
  };

  const handleRemove = async (name: string) => {
    const fullName = stats[name]?.full_name || name;
    setActionLoading(name);
    try {
      await invoke('remove_container', { name: fullName });
      await fetchStats();
    } catch (e) {
      console.error("Failed to remove container:", e);
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <div className="flex-1 flex flex-col bg-[#1e1e2e] overflow-hidden animate-fade-in">
      <div className="px-8 py-5 border-b border-[#2a2a35] bg-[#181825] shrink-0">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-white flex items-center gap-3">
              <Box className="w-6 h-6 text-blue-400" />
              Docker Containers
            </h1>
            <p className="text-gray-400 mt-2 text-sm">Real-time resource monitoring for active containers</p>
          </div>
          <button 
            onClick={fetchStats} 
            className="p-2.5 bg-[#2a2a35] hover:bg-[#323240] text-gray-300 hover:text-white rounded-lg transition-colors border border-white/5"
            title="Refresh Stats"
          >
            <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin text-blue-400' : ''}`} />
          </button>
        </div>
      </div>
      
      <div className="flex-1 overflow-y-auto custom-scrollbar">
        {containers.length === 0 && !loading ? (
          <div className="flex flex-col items-center justify-center h-full text-gray-500">
            <Box className="w-16 h-16 mb-4 opacity-20" />
            <p className="text-lg">No active Docker containers found.</p>
            <p className="text-sm mt-2">Start some services in the Services or Stacks view.</p>
          </div>
        ) : (
          <div className="flex flex-col">
            {/* Header */}
            <div className="flex items-center gap-4 px-8 py-3 border-b border-[#2a2a35] bg-[#1a1a28] text-xs font-semibold text-gray-500 uppercase tracking-wider sticky top-0 z-10">
              <div className="w-8"></div>
              <div className="flex-1 min-w-0">Container</div>
              <div className="w-64">CPU</div>
              <div className="w-64">Memory</div>
              <div className="w-20 text-right">Actions</div>
            </div>

            {/* Rows */}
            {containers.map(name => {
              const cpuPct = parsePercent(stats[name].cpu);
              const barColor = cpuPct > 80 ? 'bg-red-400' : cpuPct > 50 ? 'bg-amber-400' : 'bg-blue-400';

              return (
                <div 
                  key={name} 
                  className="flex items-center gap-4 px-8 py-4 border-b border-[#2a2a35]/50 hover:bg-[#242533] transition-colors group"
                >
                  <div className="w-8 h-8 rounded-lg bg-[#2a2a35] flex items-center justify-center shrink-0">
                    <Box className="w-4 h-4 text-gray-500 group-hover:text-blue-400 transition-colors" />
                  </div>

                  <div className="flex-1 min-w-0">
                    <span className="text-sm font-medium text-gray-200 truncate block" title={name}>{name}</span>
                  </div>

                  {/* CPU Bar */}
                  <div className="w-64 flex items-center gap-3">
                    <div className="flex-1 bg-[#181825] rounded-full h-2 overflow-hidden border border-[#2a2a35]">
                      <div 
                        className={`${barColor} h-full rounded-full transition-all duration-700 ease-out`}
                        style={{ width: `${cpuPct}%` }}
                      />
                    </div>
                    <span className="text-xs font-mono text-gray-400 w-14 text-right">{stats[name].cpu}</span>
                  </div>

                  {/* RAM Bar */}
                  <div className="w-64 flex items-center gap-3">
                    <div className="flex-1 bg-[#181825] rounded-full h-2 overflow-hidden border border-[#2a2a35]">
                      <div 
                        className="bg-purple-400 h-full rounded-full transition-all duration-700 ease-out"
                        style={{ width: `${Math.min(parsePercent(stats[name].ram), 100)}%` }}
                      />
                    </div>
                    <span className="text-xs font-mono text-gray-400 w-14 text-right">{stats[name].ram}</span>
                  </div>

                  <div className="w-20 flex items-center justify-end gap-1">
                    <button
                      onClick={() => handleStop(name)}
                      disabled={actionLoading === name}
                      className="p-2 rounded-lg text-gray-500 hover:text-amber-400 hover:bg-amber-500/10 transition-colors disabled:opacity-50"
                      title="Stop Container"
                    >
                      {actionLoading === name ? (
                        <span className="w-4 h-4 rounded-full border-2 border-current border-t-transparent animate-spin block" />
                      ) : (
                        <Square className="w-4 h-4" />
                      )}
                    </button>
                    <button
                      onClick={() => handleRemove(name)}
                      disabled={actionLoading === name}
                      className="p-2 rounded-lg text-gray-500 hover:text-red-400 hover:bg-red-500/10 transition-colors disabled:opacity-50"
                      title="Remove Container"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
