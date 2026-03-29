import type { Image } from "../types";
import ImageCard from "./ImageCard";
import Loader from "./common/Loader";
import ErrorIcon from "../assets/error.svg";

interface FeedProps {
  images: Image[];
  scrollRef: React.RefObject<HTMLDivElement | null>;
  loading: boolean;
  error: string | null;
  onRetry: () => void;
}

export default function Feed({
  images,
  scrollRef,
  loading,
  error,
  onRetry,
}: FeedProps) {
  return (
    <main className="flex-1 overflow-y-auto py-6">
      {error && images.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center py-24 px-4 text-center gap-8">
          <img src={ErrorIcon} alt="Error" width={100} height={100} />
          <h1 className="text-gray-500 text-[40px] font-medium tracking-wider">
            Something went wrong :(
          </h1>
          <button
            onClick={onRetry}
            className="px-5 py-2.5 rounded-xl font-medium text-gray-700 bg-[#e0e5ec] shadow-[4px_4px_8px_#b8b9be,-4px_-4px_8px_#ffffff] hover:shadow-[2px_2px_4px_#b8b9be,-2px_-2px_4px_#ffffff] active:shadow-[inset_3px_3px_6px_#b8b9be,inset_-3px_-3px_6px_#ffffff] transition-shadow cursor-pointer"
          >
            Retry
          </button>
        </div>
      ) : (
        <div className="flex flex-col items-center gap-6">
          {images.map((image) => (
            <ImageCard key={image.id} image={image} />
          ))}

          {/* Inline error for pagination failures */}
          {error && images.length > 0 && (
            <div className="flex flex-col items-center py-4 gap-2">
              <p className="text-gray-400 text-sm">{error}</p>
              <button
                onClick={onRetry}
                className="px-4 py-2 rounded-xl text-sm text-gray-600 bg-[#e0e5ec] shadow-[3px_3px_6px_#b8b9be,-3px_-3px_6px_#ffffff] hover:shadow-[1px_1px_3px_#b8b9be,-1px_-1px_3px_#ffffff] transition-shadow cursor-pointer"
              >
                Try again
              </button>
            </div>
          )}

          {/* Infinite scroll sentinel */}
          <div ref={scrollRef} className="h-1" />

          {loading && <Loader />}
        </div>
      )}
    </main>
  );
}
