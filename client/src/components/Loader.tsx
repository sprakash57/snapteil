import type { CSSProperties } from "react";

type LoaderSize = "sm" | "md" | "lg";

interface LoaderProps {
  size?: LoaderSize;
  shadow?: boolean;
  className?: string;
}

const sizeMap: Record<
  LoaderSize,
  { dot: string; gap: string; shadow: string }
> = {
  sm: {
    dot: "w-1.5 h-1.5",
    gap: "gap-1",
    shadow: "shadow-[2px_2px_4px_#a3aab5,-2px_-2px_4px_#ffffff]",
  },
  md: {
    dot: "w-2.5 h-2.5",
    gap: "gap-1.5",
    shadow: "shadow-[3px_3px_6px_#a3aab5,-3px_-3px_6px_#ffffff]",
  },
  lg: {
    dot: "w-5 h-5",
    gap: "gap-3",
    shadow: "shadow-[4px_4px_8px_#a3aab5,-4px_-4px_8px_#ffffff]",
  },
};

export default function Loader({
  size = "lg",
  shadow = true,
  className = "",
}: LoaderProps) {
  const { dot, gap, shadow: shadowClass } = sizeMap[size];

  return (
    <span
      className={`inline-flex items-center justify-center ${gap} ${className}`}
    >
      {[0, 1, 2].map((i) => (
        <span
          key={i}
          className={`rounded-full bg-[#c8d0dc] ${dot} ${shadow ? shadowClass : ""} animate-loader`}
          style={{ ["--loader-delay"]: `${i * 0.15}s` } as CSSProperties}
        />
      ))}
    </span>
  );
}
