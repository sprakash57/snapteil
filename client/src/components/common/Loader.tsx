import type { CSSProperties } from "react";

export default function Loader() {
  return (
    <div className="flex items-center justify-center py-6 gap-3">
      {[0, 1, 2].map((i) => (
        <span
          key={i}
          className="w-5 h-5 rounded-full bg-[#c8d0dc] shadow-[4px_4px_8px_#a3aab5,-4px_-4px_8px_#ffffff] animate-loader"
          style={{ ["--loader-delay"]: `${i * 0.15}s` } as CSSProperties}
        />
      ))}
    </div>
  );
}
