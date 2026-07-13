import { createContext, useContext, useState, type ReactNode } from "react";

interface ConfirmOptions {
  title: string;
  message: string;
  kind?: "info" | "warning" | "error";
  confirmText?: string;
  cancelText?: string;
}

interface ConfirmContextType {
  confirm: (options: ConfirmOptions) => Promise<boolean>;
}

const ConfirmContext = createContext<ConfirmContextType | undefined>(undefined);

export function useConfirm() {
  const context = useContext(ConfirmContext);
  if (!context) throw new Error("useConfirm must be used within ConfirmProvider");
  return context;
}

export function ConfirmProvider({ children }: { children: ReactNode }) {
  const [isOpen, setIsOpen] = useState(false);
  const [options, setOptions] = useState<ConfirmOptions | null>(null);
  const [resolvePromise, setResolvePromise] = useState<(value: boolean) => void>();

  const confirm = (opts: ConfirmOptions) => {
    setOptions(opts);
    setIsOpen(true);
    return new Promise<boolean>((resolve) => {
      setResolvePromise(() => resolve);
    });
  };

  const handleConfirm = () => {
    setIsOpen(false);
    if (resolvePromise) resolvePromise(true);
  };

  const handleCancel = () => {
    setIsOpen(false);
    if (resolvePromise) resolvePromise(false);
  };

  return (
    <ConfirmContext.Provider value={{ confirm }}>
      {children}
      {isOpen && options && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4 animate-fade-in" onMouseDown={handleCancel}>
          <div 
            className="bg-[#242533] border border-[#2a2a35] rounded-2xl p-6 shadow-2xl max-w-sm w-full relative animate-modal"
            onMouseDown={e => e.stopPropagation()}
          >
            <h3 className={`text-xl font-bold mb-2 ${options.kind === 'warning' || options.kind === 'error' ? 'text-red-400' : 'text-white'}`}>
              {options.title}
            </h3>
            <p className="text-gray-300 text-sm mb-6 leading-relaxed whitespace-pre-wrap">
              {options.message}
            </p>
            <div className="flex items-center justify-end space-x-3">
              <button 
                onClick={handleCancel}
                className="px-4 py-2 rounded-lg text-sm font-medium text-gray-300 hover:bg-[#2a2a35] transition-colors"
              >
                {options.cancelText || "Cancel"}
              </button>
              <button 
                onClick={handleConfirm}
                className={`px-4 py-2 rounded-lg text-sm font-medium text-white transition-colors ${
                  options.kind === 'warning' || options.kind === 'error' 
                    ? "bg-red-500 hover:bg-red-600" 
                    : "bg-blue-600 hover:bg-blue-500"
                }`}
              >
                {options.confirmText || "Confirm"}
              </button>
            </div>
          </div>
        </div>
      )}
    </ConfirmContext.Provider>
  );
}
