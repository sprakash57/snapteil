import { useRef, useState } from "react";
import Header from "./components/Header";
import Feed from "./components/Feed";
import UploadModal from "./components/UploadModal";
import NewPostBanner from "./components/NewPostBanner";
import { useImages } from "./hooks/useImages";
import { useInfiniteScroll } from "./hooks/useInfiniteScroll";
import { useWebSocket } from "./hooks/useWebSocket";
import Loader from "./components/common/Loader";
import type { Image } from "./types";

function App() {
  const [showUpload, setShowUpload] = useState(false);
  const [pendingImages, setPendingImages] = useState<Image[]>([]);
  const feedRef = useRef<HTMLElement>(null);
  const {
    images,
    loading,
    initialLoading,
    hasMore,
    error,
    loadMore,
    addImage,
    clearError,
  } = useImages();
  const scrollRef = useInfiniteScroll(loadMore, hasMore, loading);
  useWebSocket((img) => setPendingImages((prev) => [img, ...prev]));

  function handleBannerClick() {
    pendingImages.forEach(addImage);
    setPendingImages([]);
    feedRef.current?.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="min-h-screen bg-[#e0e5ec] flex flex-col">
      <Header onUploadClick={() => setShowUpload(true)} />

      {pendingImages.length > 0 && !initialLoading && (
        <NewPostBanner
          count={pendingImages.length}
          onClick={handleBannerClick}
        />
      )}

      {initialLoading ? (
        <div className="flex-1 flex items-center justify-center text-gray-500">
          <Loader />
        </div>
      ) : (
        <Feed
          images={images}
          scrollRef={scrollRef}
          feedRef={feedRef}
          loading={loading}
          error={error}
          onRetry={() => {
            clearError();
            loadMore();
          }}
        />
      )}

      {showUpload && (
        <UploadModal
          onClose={() => setShowUpload(false)}
          onUploaded={addImage}
        />
      )}
    </div>
  );
}

export default App;
