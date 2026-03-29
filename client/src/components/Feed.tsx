import type { Image } from "../types";
import ImageCard from "./ImageCard";

interface FeedProps {
  images: Image[];
  sentinelRef: React.RefObject<HTMLDivElement | null>;
  loading: boolean;
}

export default function Feed({ images, sentinelRef, loading }: FeedProps) {
  return (
    <main className="flex-1 overflow-y-auto py-6">
      <div className="flex flex-col items-center gap-6">
        {images.map((image) => (
          <ImageCard key={image.id} image={image} />
        ))}

        {/* Infinite scroll sentinel */}
        <div ref={sentinelRef} className="h-1" />

        {loading && (
          <div className="py-4 text-gray-500 text-sm">Loading more...</div>
        )}
      </div>
    </main>
  );
}
