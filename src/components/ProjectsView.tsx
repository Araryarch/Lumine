import { useState } from "react";
import { invoke } from "@tauri-apps/api/core";
import { ask } from "@tauri-apps/plugin-dialog";
import { ProjectCard } from "./ProjectCard";
import { ProjectCreationWizard } from "./ProjectCreationWizard";
import { ProjectDetailsSidebar } from "./ProjectDetailsSidebar";
import { Plus } from "lucide-react";
import { useAppStore } from "../store/useAppStore";
import { toProject } from "../utils/transformers";

export function ProjectsView() {
  const { projects: rawProjects, refresh } = useAppStore();
  const projects = rawProjects.map(toProject);
  
  const [isWizardOpen, setIsWizardOpen] = useState(false);
  const [selectedProjectId, setSelectedProjectId] = useState<string | null>(null);
  const [loadingIds, setLoadingIds] = useState<Set<string>>(new Set());

  const toggleProject = async (id: string, start: boolean) => {
    setLoadingIds(prev => new Set(prev).add(id));
    try {
      if (start) {
        await invoke("start_project", { id });
      } else {
        await invoke("stop_project", { id });
      }
      await refresh();
    } catch (e) {
      console.error("Toggle project failed", e);
    } finally {
      setLoadingIds(prev => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    }
  };

  const deleteProject = async (id: string) => {
    const proj = projects.find(p => p.id === id);
    if (!proj) return;
    
    const yes = await ask("Are you sure you want to remove this project from Lumine?", {
      title: "Remove Project",
      kind: "warning",
    });
    
    if (yes) {
      const deleteFolder = await ask(`Do you also want to permanently delete the project folder from your disk?\n\nPath: ${proj.path}\n\nWARNING: This cannot be undone!`, {
        title: "Delete Folder?",
        kind: "warning",
      });
      
      try {
        await invoke("delete_project", { id, deleteFolder });
        if (selectedProjectId === id) setSelectedProjectId(null);
        await refresh();
      } catch (e) {
        console.error(e);
      }
    }
  };

  const reloadProject = async (id: string) => {
    setLoadingIds(prev => new Set(prev).add(id));
    try {
      await invoke("stop_project", { id });
      // Short delay to let docker cleanup finish
      await new Promise(r => setTimeout(r, 1000));
      await invoke("start_project", { id });
      await refresh();
    } catch (e) {
      console.error("Reload project failed", e);
    } finally {
      setLoadingIds(prev => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    }
  };

  return (
    <div className="flex flex-col h-full">
      <header className="h-14 border-b border-[#2a2a35] flex items-center justify-between px-6 shrink-0 bg-[#1e1e2e]/50 backdrop-blur-sm z-10">
        <div className="flex items-center gap-4">
          <h1 className="text-sm font-semibold text-white">Project Manager</h1>
        </div>
        <div className="flex items-center gap-3">
          {projects.length > 0 && (
            <button
              onClick={() => setIsWizardOpen(true)}
              className="flex items-center gap-2 px-3 py-1.5 bg-blue-500 hover:bg-blue-600 text-white rounded-md text-xs font-medium transition-colors"
            >
              <Plus className="w-4 h-4" />
              Create Project
            </button>
          )}
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden relative">
        <div className="flex-1 overflow-y-auto p-6 custom-scrollbar relative">
          {projects.length === 0 ? (
            <div className="h-full flex flex-col items-center justify-center text-center max-w-sm mx-auto">
              <h2 className="text-xl font-bold text-white mb-2">No Projects Found</h2>
              <p className="text-gray-400 mb-8 leading-relaxed">
                You haven't created any projects yet. Use Lumine to easily scaffold and run modern web frameworks.
              </p>
              <button
                onClick={() => setIsWizardOpen(true)}
                className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-2.5 rounded-lg font-medium transition-colors flex items-center gap-2"
              >
                <Plus className="w-5 h-5" />
                Create Your First Project
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6 transition-all duration-300">
              {projects.map((proj) => (
                <ProjectCard
                  key={proj.id}
                  project={proj}
                  isLoading={loadingIds.has(proj.id)}
                  onClick={() => setSelectedProjectId(proj.id === selectedProjectId ? null : proj.id)}
                  onToggle={(e) => {
                    e.stopPropagation();
                    toggleProject(proj.id, proj.status !== "running");
                  }}
                />
              ))}
            </div>
          )}
        </div>

        {selectedProjectId && projects.some(p => p.id === selectedProjectId) && (
          <ProjectDetailsSidebar
            project={projects.find(p => p.id === selectedProjectId)!}
            onClose={() => setSelectedProjectId(null)}
            onToggle={() => toggleProject(selectedProjectId, projects.find(p => p.id === selectedProjectId)?.status !== "running")}
            onReload={() => reloadProject(selectedProjectId)}
            onDelete={() => deleteProject(selectedProjectId)}
            onUpdate={refresh}
            isLoading={loadingIds.has(selectedProjectId)}
          />
        )}
      </div>

      {isWizardOpen && (
        <ProjectCreationWizard
          onClose={() => setIsWizardOpen(false)}
          onCreated={() => {
            refresh();
          }}
        />
      )}
    </div>
  );
}
