import React, { useState, useEffect } from 'react';
import { invoke } from '@tauri-apps/api/core';
import type { Stack, Service } from '../types';
import { Layers, Plus, Play, Square, Trash2, Edit2, Check, X, Database, Server, Monitor, Code, Download, Loader2 } from 'lucide-react';

const getStackIconComponent = (stackId: string) => {
  if (stackId.includes('xampp') || stackId.includes('db') || stackId.includes('lamp') || stackId.includes('lemp')) return Database;
  if (stackId.includes('mern') || stackId.includes('node') || stackId.includes('python')) return Code;
  if (stackId.includes('laravel') || stackId.includes('backend') || stackId.includes('spring') || stackId.includes('rails') || stackId.includes('go')) return Server;
  if (stackId.includes('lemp') || stackId.includes('web')) return Monitor;
  return Layers;
};

import { useConfirm } from './ConfirmProvider';
import { LogViewer } from './LogViewer';

export const StacksView: React.FC = () => {
  const [stacks, setStacks] = useState<Stack[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingStack, setEditingStack] = useState<Stack | null>(null);
  const [selectedStackId, setSelectedStackId] = useState<string | null>(null);
  const [combinedLogs, setCombinedLogs] = useState<string[]>([]);
  const [isPulling, setIsPulling] = useState(false);
  
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [color, setColor] = useState('#3b82f6'); // default blue
  const [selectedServices, setSelectedServices] = useState<string[]>([]);
  
  const { confirm } = useConfirm();

  const loadData = async () => {
    try {
      const fetchedStacks = await invoke<Stack[]>('get_stacks');
      const fetchedServices = await invoke<Service[]>('get_services');
      setStacks(fetchedStacks);
      setServices(fetchedServices);
    } catch (e) {
      console.error('Failed to load data:', e);
    }
  };

  useEffect(() => {
    loadData();
    // Poll services and logs
    const interval = setInterval(async () => {
      try {
        const fetchedServices = await invoke<Service[]>('get_services');
        setServices(fetchedServices);

        if (selectedStackId) {
          const stack = stacks.find(s => s.id === selectedStackId);
          if (stack) {
            let allLogs: string[] = [];
            for (const svcId of stack.services) {
              const svc = fetchedServices.find(s => s.id === svcId);
              if (svc && svc.log && svc.log.length > 0) {
                allLogs.push(`\n--- ${svc.name} ---`);
                allLogs.push(...svc.log.map(l => `[${svc.name}] ${l}`));
              }
            }
            
            // Only update if logs changed to prevent LogViewer from constantly scrolling
            setCombinedLogs(prev => {
              const joinedPrev = prev.join('\n');
              const joinedNew = allLogs.join('\n');
              if (joinedPrev === joinedNew) return prev;
              return allLogs;
            });
          }
        }
      } catch(e) { console.error(e) }
    }, 2000);
    return () => clearInterval(interval);
  }, [stacks, selectedStackId]);

  const handlePullImages = async (stack: Stack) => {
    const dockerImages = stack.services
      .map(svcId => services.find(s => s.id === svcId))
      .filter(s => s?.config?.runner === 'docker')
      .map(s => s?.config?.executablePath)
      .filter((img): img is string => !!img);

    if (dockerImages.length === 0) return;

    setIsPulling(true);
    for (const image of dockerImages) {
      try {
        await invoke('pull_docker_image', { image });
      } catch (e) {
        console.error("Failed to pull", image, e);
      }
    }
    setIsPulling(false);
  };

  const resetForm = () => {
    setName('');
    setDescription('');
    setColor('#3b82f6');
    setSelectedServices([]);
    setEditingStack(null);
  };

  const handleOpenModal = (stack?: Stack) => {
    if (stack) {
      setEditingStack(stack);
      setName(stack.name);
      setDescription(stack.description);
      setColor(stack.color);
      setSelectedServices(stack.services);
    } else {
      resetForm();
    }
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    resetForm();
  };

  const handleSave = async () => {
    if (!name.trim()) return;
    
    const id = editingStack ? editingStack.id : name.toLowerCase().replace(/\s+/g, '_') + '_' + Date.now();
    const payload: Stack = { id, name, description, color, services: selectedServices };

    try {
      if (editingStack) {
        await invoke('edit_stack', { id, info: payload });
      } else {
        await invoke('add_stack', { info: payload });
      }
      handleCloseModal();
      loadData();
    } catch (e) {
      console.error(e);
      alert('Failed to save stack: ' + e);
    }
  };

  const handleDelete = async (id: string, stackName: string) => {
    const yes = await confirm({
      message: `Are you sure you want to delete the stack "${stackName}"?`,
      title: 'Delete Stack',
      confirmText: 'Delete',
      cancelText: 'Cancel',
      kind: 'error'
    });
    if (yes) {
      try {
        await invoke('delete_stack', { id });
        loadData();
      } catch (e) {
        alert('Failed to delete stack: ' + e);
      }
    }
  };

  const handleStartAll = async (id: string) => {
    try {
      await invoke('start_stack', { id });
    } catch (e) {
      console.error(e);
      alert('Failed to start stack: ' + e);
    }
  };

  const handleStopAll = async (id: string) => {
    try {
      await invoke('stop_stack', { id });
    } catch (e) {
      console.error(e);
      alert('Failed to stop stack: ' + e);
    }
  };

  const toggleServiceSelection = (svcId: string) => {
    if (selectedServices.includes(svcId)) {
      setSelectedServices(selectedServices.filter(id => id !== svcId));
    } else {
      setSelectedServices([...selectedServices, svcId]);
    }
  };

  const getServiceStatus = (svcId: string) => {
    const svc = services.find(s => s.id === svcId);
    return svc ? (typeof svc.status === 'object' ? 'error' : svc.status.toLowerCase()) : 'stopped';
  };

  return (
    <div className="p-8 h-full flex flex-col overflow-hidden animate-fade-in">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2 flex items-center gap-3">
            <Layers className="w-8 h-8 text-blue-400" />
            Stacks
          </h1>
          <p className="text-gray-400">Group your services and manage them together.</p>
        </div>
        <button
          onClick={() => handleOpenModal()}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors shadow-lg shadow-blue-500/20 font-medium"
        >
          <Plus className="w-5 h-5" />
          Create Stack
        </button>
      </div>

      <div className="flex flex-1 gap-6 overflow-hidden">
        {/* Left Sidebar (Master) */}
        <div className="w-64 bg-transparent flex flex-col overflow-y-auto custom-scrollbar pr-4">
          <div className="flex justify-between items-center mb-4 px-2">
            <h2 className="text-[11px] font-semibold text-gray-500 uppercase tracking-wider">All Stacks</h2>
            <button onClick={() => handleOpenModal()} className="text-gray-500 hover:text-white transition-colors">
              <Plus className="w-4 h-4" />
            </button>
          </div>
          {stacks.length === 0 ? (
            <div className="text-center text-gray-500 italic mt-4 text-sm">No stacks found.</div>
          ) : (
            <div className="space-y-0.5">
              {stacks.map(stack => {
                const IconComp = getStackIconComponent(stack.id);
                return (
                  <button
                    key={stack.id}
                    onClick={() => setSelectedStackId(stack.id)}
                    className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all ${
                      selectedStackId === stack.id 
                        ? 'bg-[#2a2a3c] text-white' 
                        : 'hover:bg-white/5 text-gray-400'
                    }`}
                  >
                    <div className="shrink-0 flex items-center justify-center w-6 h-6 rounded-md" style={{ backgroundColor: `${stack.color}20` }}>
                      <IconComp className="w-3.5 h-3.5" style={{ color: stack.color }} />
                    </div>
                    <div className="text-left flex-1 min-w-0">
                      <h3 className={`font-medium truncate text-[14px] ${selectedStackId === stack.id ? 'text-white' : 'text-gray-300'}`}>
                        {stack.name}
                      </h3>
                    </div>
                  </button>
                );
              })}
            </div>
          )}
        </div>

        {/* Right Panel (Detail) */}
        <div className="flex-1 bg-transparent border-l border-white/5 pl-8 flex flex-col overflow-y-auto custom-scrollbar">
          {stacks.length === 0 ? (
            <div className="flex-1 flex flex-col items-center justify-center text-gray-500">
              <Layers className="w-12 h-12 mb-4 opacity-20" />
              <p className="text-sm">No stacks found. Create one to get started.</p>
            </div>
          ) : !selectedStackId || !stacks.find(s => s.id === selectedStackId) ? (
            <div className="flex-1 flex items-center justify-center text-gray-500">
              <p className="text-sm">Select a stack from the left panel.</p>
            </div>
          ) : (
            (() => {
              const stack = stacks.find(s => s.id === selectedStackId)!;
              const activeCount = stack.services.filter(id => getServiceStatus(id) === 'running').length;
              const totalCount = stack.services.length;
              const HeaderIconComp = getStackIconComponent(stack.id);

              const dockerImages = stack.services
                .map(svcId => services.find(s => s.id === svcId))
                .filter(s => s?.config?.runner === 'docker')
                .map(s => s?.config?.executablePath)
                .filter(Boolean);

              return (
                <div className="flex flex-col h-full animate-fade-in max-w-4xl">
                  <div className="flex justify-between items-start mb-6 shrink-0">
                    <div className="flex items-center gap-4">
                      <div className="shrink-0 flex items-center justify-center w-10 h-10 rounded-xl shadow-lg" style={{ backgroundColor: `${stack.color}15`, border: `1px solid ${stack.color}30` }}>
                        <HeaderIconComp className="w-5 h-5" style={{ color: stack.color }} />
                      </div>
                      <div>
                        <h3 className="text-2xl font-semibold text-white tracking-tight mb-1">{stack.name}</h3>
                        <p className="text-sm text-gray-400">{stack.description}</p>
                      </div>
                    </div>
                    <div className="flex gap-1">
                      <button onClick={() => handleOpenModal(stack)} className="p-2 text-gray-500 hover:text-white transition-colors rounded-lg hover:bg-white/5">
                        <Edit2 className="w-4 h-4" />
                      </button>
                      <button onClick={() => handleDelete(stack.id, stack.name)} className="p-2 text-gray-500 hover:text-red-400 transition-colors rounded-lg hover:bg-red-500/10">
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>

                  <div className="flex-1 flex flex-col min-h-0 gap-6">
                    <div className="flex-none flex flex-col max-h-[40%] min-h-0">
                      <h4 className="text-[12px] font-semibold text-gray-500 mb-3 uppercase tracking-wider shrink-0">
                        Services ({activeCount}/{totalCount} running)
                      </h4>
                      
                      <div className="border border-white/5 rounded-xl bg-[#161622]/30 divide-y divide-white/5 overflow-y-auto custom-scrollbar shrink-0">
                        {stack.services.length === 0 ? (
                          <div className="px-4 py-4 text-sm text-gray-500 italic">No services assigned</div>
                        ) : (
                          stack.services.map(svcId => {
                            const svcName = services.find(s => s.id === svcId)?.name || svcId;
                            const isRunning = getServiceStatus(svcId) === 'running';
                            return (
                              <div key={svcId} className="group flex items-center justify-between text-[14px] text-gray-300 hover:bg-white/5 px-4 py-2 transition-colors">
                                <div className="flex items-center gap-3">
                                  <div className={`w-2 h-2 rounded-full shrink-0 ${isRunning ? 'bg-green-400 shadow-[0_0_8px_rgba(74,222,128,0.5)]' : 'bg-gray-600'}`} />
                                  {svcName}
                                </div>
                                <div className="opacity-0 group-hover:opacity-100 transition-opacity flex gap-2">
                                  {services.find(s => s.id === svcId) ? (
                                    isRunning ? (
                                      <button 
                                        onClick={() => invoke('stop_service', { id: svcId })}
                                        title="Stop Service"
                                        className="p-1 hover:bg-red-500/10 text-red-400 rounded-md transition-colors"
                                      >
                                        <Square className="w-3.5 h-3.5" />
                                      </button>
                                    ) : (
                                      <button 
                                        onClick={() => invoke('start_service', { id: svcId })}
                                        title="Start Service"
                                        className="p-1 hover:bg-green-500/10 text-green-400 rounded-md transition-colors"
                                      >
                                        <Play className="w-3.5 h-3.5" fill="currentColor" />
                                      </button>
                                    )
                                  ) : (
                                    <span className="text-[11px] text-gray-500 italic">Not Installed</span>
                                  )}
                                </div>
                              </div>
                            );
                          })
                        )}
                      </div>
                    </div>

                    <div className="flex-1 min-h-0 flex flex-col overflow-hidden">
                      <LogViewer 
                        log={combinedLogs.length > 0 ? combinedLogs : ["No logs available for this stack yet."]} 
                        title="Stack Logs" 
                      />
                    </div>
                  </div>

                  <div className="flex items-center gap-3 mt-6 shrink-0 pt-4 border-t border-white/5">
                    {totalCount > 0 && activeCount < totalCount && (
                      <button
                        onClick={() => handleStartAll(stack.id)}
                        className="flex items-center justify-center gap-2 px-6 py-2.5 bg-green-500/10 hover:bg-green-500/20 text-green-400 rounded-lg transition-colors text-sm font-medium border border-transparent hover:border-green-500/20"
                      >
                        <Play className="w-3.5 h-3.5" fill="currentColor" /> Start All
                      </button>
                    )}
                    {totalCount > 0 && activeCount > 0 && (
                      <button
                        onClick={() => handleStopAll(stack.id)}
                        className="flex items-center justify-center gap-2 px-6 py-2.5 bg-red-500/10 hover:bg-red-500/20 text-red-400 rounded-lg transition-colors text-sm font-medium border border-transparent hover:border-red-500/20"
                      >
                        <Square className="w-3.5 h-3.5" /> Stop All
                      </button>
                    )}
                    
                    {dockerImages.length > 0 && (
                      <button
                        onClick={() => handlePullImages(stack)}
                        disabled={isPulling}
                        className="flex items-center justify-center gap-2 px-6 py-2.5 bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 rounded-lg transition-colors text-sm font-medium border border-transparent hover:border-blue-500/20 ml-auto disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        {isPulling ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Download className="w-3.5 h-3.5" />} 
                        {isPulling ? 'Pulling...' : 'Pull Images'}
                      </button>
                    )}
                  </div>
                </div>
              );
            })()
          )}
        </div>
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 animate-fade-in p-4">
          <div className="bg-[#1e1e2e] border border-[#2a2a3c] rounded-2xl w-full max-w-lg shadow-2xl flex flex-col max-h-[90vh]">
            <div className="flex justify-between items-center p-5 border-b border-[#2a2a3c]">
              <h2 className="text-xl font-bold text-white flex items-center gap-2">
                <Layers className="w-5 h-5 text-blue-400" />
                {editingStack ? 'Edit Stack' : 'Create Stack'}
              </h2>
              <button onClick={handleCloseModal} className="text-gray-400 hover:text-white transition-colors bg-[#2a2a3c] rounded-full p-1.5 hover:bg-[#3b3b4f]">
                <X className="w-4 h-4" />
              </button>
            </div>
            
            <div className="p-5 overflow-y-auto flex-1 custom-scrollbar">
              <div className="grid grid-cols-2 gap-4 mb-5">
                <div>
                  <label className="block text-sm font-medium text-gray-400 mb-2">Stack Name</label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="e.g. Production Stack"
                    className="w-full bg-[#181825] border border-[#2a2a3c] rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-blue-500 transition-colors"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-400 mb-2">Theme Color</label>
                  <div className="flex items-center gap-3 bg-[#181825] border border-[#2a2a3c] rounded-lg px-3 py-1">
                    <input
                      type="color"
                      value={color}
                      onChange={(e) => setColor(e.target.value)}
                      className="w-8 h-8 rounded cursor-pointer bg-transparent border-0 p-0"
                    />
                    <span className="text-gray-300 font-mono">{color}</span>
                  </div>
                </div>
              </div>

              <div className="mb-5">
                <label className="block text-sm font-medium text-gray-400 mb-2">Description</label>
                <input
                  type="text"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Optional description"
                  className="w-full bg-[#181825] border border-[#2a2a3c] rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-blue-500 transition-colors"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-400 mb-3 flex justify-between items-center">
                  <span>Assign Services</span>
                  <span className="text-xs bg-[#2a2a3c] px-2 py-1 rounded text-gray-300">{selectedServices.length} selected</span>
                </label>
                <div className="bg-[#181825] border border-[#2a2a3c] rounded-xl p-4 max-h-64 overflow-y-auto custom-scrollbar">
                  {services.length === 0 ? (
                    <div className="text-center text-gray-500 py-4">No services available to assign.</div>
                  ) : (
                    <div className="grid grid-cols-2 gap-3">
                      {services.map(svc => {
                        const isSelected = selectedServices.includes(svc.id);
                        return (
                          <div
                            key={svc.id}
                            onClick={() => toggleServiceSelection(svc.id)}
                            className={`flex items-center gap-3 p-3 rounded-lg cursor-pointer transition-colors border ${isSelected ? 'bg-blue-500/10 border-blue-500/50' : 'bg-[#2a2a3c] border-transparent hover:border-gray-500'}`}
                          >
                            <div className={`w-5 h-5 rounded flex items-center justify-center ${isSelected ? 'bg-blue-500 text-white' : 'bg-[#181825] border border-[#3b3b4f]'}`}>
                              {isSelected && <Check className="w-3 h-3" />}
                            </div>
                            <div className="flex-1 truncate">
                              <div className="text-sm font-medium text-white">{svc.name}</div>
                              <div className="text-xs text-gray-400 truncate">{svc.config.serviceType}</div>
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  )}
                </div>
              </div>
            </div>

            <div className="p-5 border-t border-[#2a2a3c] flex justify-end gap-3 bg-[#181825] rounded-b-2xl">
              <button
                onClick={handleCloseModal}
                className="px-5 py-2.5 text-gray-400 hover:text-white transition-colors font-medium"
              >
                Cancel
              </button>
              <button
                onClick={handleSave}
                disabled={!name.trim()}
                className="flex items-center gap-2 px-6 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors shadow-lg shadow-blue-500/20 font-medium"
              >
                <Check className="w-4 h-4" />
                {editingStack ? 'Save Changes' : 'Create Stack'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
