import { useState, useRef, useEffect } from 'react';
import { ChevronDown } from 'lucide-react';

export interface SelectOption {
  label: string;
  value: string;
}

interface CustomSelectProps {
  options: SelectOption[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  searchable?: boolean;
  allowCustom?: boolean;
}

export function CustomSelect({ options, value, onChange, placeholder = "Select...", searchable = false, allowCustom = false }: CustomSelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [search, setSearch] = useState("");
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const selectedOption = options.find(o => o.value === value);
  const displayValue = selectedOption ? selectedOption.label : value || placeholder;

  const filteredOptions = searchable || allowCustom
    ? options.filter(o => o.label.toLowerCase().includes(search.toLowerCase()) || o.value.toLowerCase().includes(search.toLowerCase()))
    : options;

  return (
    <div className="relative w-full" ref={dropdownRef}>
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500 flex items-center justify-between transition-colors text-left truncate"
      >
        <span className={value ? "text-white" : "text-gray-500"}>
          {displayValue}
        </span>
        <ChevronDown size={16} className={`text-gray-400 transition-transform flex-shrink-0 ${isOpen ? 'rotate-180' : ''}`} />
      </button>

      {isOpen && (
        <div className="absolute z-50 w-full mt-2 bg-[#1e1e2e] border border-[#2a2a35] rounded-lg overflow-hidden flex flex-col shadow-xl shadow-black/50">
          {(searchable || allowCustom) && (
            <div className="p-2 border-b border-[#2a2a35] bg-[#181825]">
              <input
                type="text"
                placeholder={allowCustom ? "Search or type custom..." : "Search..."}
                className="w-full bg-[#181825] text-xs text-white placeholder-gray-500 focus:outline-none"
                value={search}
                onChange={e => setSearch(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && allowCustom && search.trim() !== '') {
                    onChange(search.trim());
                    setIsOpen(false);
                    setSearch("");
                  }
                }}
                autoFocus
              />
            </div>
          )}
          <div className="max-h-48 overflow-y-auto custom-scrollbar p-1">
            {filteredOptions.length > 0 ? (
              filteredOptions.map(opt => (
                <button
                  key={opt.value}
                  type="button"
                  onClick={() => {
                    onChange(opt.value);
                    setIsOpen(false);
                    setSearch("");
                  }}
                  className={`w-full text-left px-3 py-2 text-xs rounded transition-colors ${
                    value === opt.value 
                      ? 'bg-blue-600 text-white' 
                      : 'text-gray-300 hover:text-white hover:bg-[#32323d]'
                  }`}
                >
                  {opt.label}
                </button>
              ))
            ) : (
              !allowCustom && <div className="px-3 py-2 text-xs text-gray-500 text-center">No options found</div>
            )}
            
            {allowCustom && search.trim() !== '' && !options.some(o => o.value.toLowerCase() === search.trim().toLowerCase()) && (
              <button
                type="button"
                onClick={() => {
                  onChange(search.trim());
                  setIsOpen(false);
                  setSearch("");
                }}
                className="w-full text-left px-3 py-2 text-xs rounded text-blue-400 hover:text-white hover:bg-blue-600/20 transition-colors border-t border-[#2a2a35] mt-1"
              >
                Use custom: "{search.trim()}"
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
