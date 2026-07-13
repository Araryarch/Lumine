import { create } from "zustand";
import { invoke } from "@tauri-apps/api/core";
import type { RawServiceInfo, RawProjectInfo } from "../types";

interface AppState {
  services: RawServiceInfo[];
  projects: RawProjectInfo[];
  hasDocker: boolean;
  isRefreshing: boolean;
  dockerStats: Record<string, { cpu: string; ram: string; full_name: string }>;
  setHasDocker: (hasDocker: boolean) => void;
  refresh: () => Promise<void>;
  startPolling: () => void;
  stopPolling: () => void;
}

let pollInterval: number | null = null;
let lastServicesJson = "";
let lastProjectsJson = "";

export const useAppStore = create<AppState>((set, get) => ({
  services: [],
  projects: [],
  hasDocker: true,
  isRefreshing: false,
  dockerStats: {},
  setHasDocker: (hasDocker) => set({ hasDocker }),
  
  refresh: async () => {
    if (get().isRefreshing) return;
    set({ isRefreshing: true });
    
    try {
      const svcList = await invoke<RawServiceInfo[]>("get_services");
      const svcJson = JSON.stringify(svcList);
      if (svcJson !== lastServicesJson) {
        lastServicesJson = svcJson;
        set({ services: svcList });
      }

      const projList = await invoke<RawProjectInfo[]>("get_projects");
      const projJson = JSON.stringify(projList);
      if (projJson !== lastProjectsJson) {
        lastProjectsJson = projJson;
        set({ projects: projList });
      }

      // Fetch docker stats silently
      if (get().hasDocker) {
        try {
          const stats = await invoke<Record<string, { cpu: string; ram: string; full_name: string }>>("get_docker_stats");
          set({ dockerStats: stats });
        } catch (e) {
          // ignore stats errors
        }
      }
    } catch (e) {
      console.error("Failed to fetch data from backend:", e);
    } finally {
      set({ isRefreshing: false });
    }
  },

  startPolling: () => {
    if (pollInterval) return;
    
    // Check docker status once on startup
    invoke<boolean>("check_docker").then(status => {
      get().setHasDocker(status);
    });

    get().refresh(); // Initial fetch
    pollInterval = window.setInterval(() => {
      get().refresh();
    }, 3000);
  },

  stopPolling: () => {
    if (pollInterval) {
      window.clearInterval(pollInterval);
      pollInterval = null;
    }
  },
}));
