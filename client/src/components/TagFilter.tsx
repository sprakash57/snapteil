import { useEffect, useRef, useState } from "react";

interface TagFilterProps {
  allTags: string[];
  selectedTags: string[];
  onChange: (tags: string[]) => void;
}

export default function TagFilter({
  allTags,
  selectedTags,
  onChange,
}: TagFilterProps) {
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  // Close dropdown on outside click
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(e.target as Node)
      ) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  function toggle(tag: string) {
    if (selectedTags.includes(tag)) {
      onChange(selectedTags.filter((t) => t !== tag));
    } else {
      onChange([...selectedTags, tag]);
    }
  }

  function remove(tag: string) {
    onChange(selectedTags.filter((t) => t !== tag));
  }

  return (
    <div ref={containerRef} className="flex items-center gap-2 flex-wrap">
      {/* Dropdown trigger */}
      <div className="relative">
        <button
          onClick={() => setOpen((o) => !o)}
          className="px-4 py-2 rounded-xl text-sm text-gray-600 bg-[#e0e5ec] shadow-[3px_3px_6px_#b8b9be,-3px_-3px_6px_#ffffff] hover:shadow-[2px_2px_4px_#b8b9be,-2px_-2px_4px_#ffffff] active:shadow-[inset_2px_2px_4px_#b8b9be,inset_-2px_-2px_4px_#ffffff] transition-shadow cursor-pointer flex items-center gap-1.5"
        >
          Filter by tags
          <svg
            className={`w-3.5 h-3.5 transition-transform ${open ? "rotate-180" : ""}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </button>

        {/* Dropdown menu */}
        {open && allTags.length > 0 && (
          <div className="absolute left-0 top-full mt-2 w-48 max-h-60 overflow-y-auto rounded-xl bg-[#e0e5ec] shadow-[5px_5px_10px_#b8b9be,-5px_-5px_10px_#ffffff] z-50 py-1">
            {allTags.map((tag) => {
              const active = selectedTags.includes(tag);
              return (
                <button
                  key={tag}
                  onClick={() => toggle(tag)}
                  className={`w-full text-left px-4 py-2 text-sm cursor-pointer transition-all ${
                    active
                      ? "text-gray-800 font-medium bg-[#d5dae3] shadow-[inset_2px_2px_4px_#b8b9be,inset_-2px_-2px_4px_#ffffff]"
                      : "text-gray-600 hover:bg-[#d9dee7]"
                  }`}
                >
                  <span className="flex items-center gap-2">
                    <span
                      className={`w-3.5 h-3.5 rounded border shrink-0 flex items-center justify-center ${
                        active
                          ? "bg-gray-600 border-gray-600"
                          : "border-gray-400 bg-transparent"
                      }`}
                    >
                      {active && (
                        <svg
                          className="w-2.5 h-2.5 text-white"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={3}
                            d="M5 13l4 4L19 7"
                          />
                        </svg>
                      )}
                    </span>
                    {tag}
                  </span>
                </button>
              );
            })}
          </div>
        )}
      </div>

      {/* Selected tag pills */}
      {selectedTags.map((tag) => (
        <span
          key={tag}
          className="inline-flex items-center gap-1 px-3 py-1 text-xs rounded-full text-gray-700 bg-[#e0e5ec] shadow-[inset_2px_2px_4px_#b8b9be,inset_-2px_-2px_4px_#ffffff]"
        >
          {tag}
          <button
            onClick={() => remove(tag)}
            className="ml-0.5 text-gray-400 hover:text-gray-700 cursor-pointer leading-none"
          >
            ×
          </button>
        </span>
      ))}
    </div>
  );
}
