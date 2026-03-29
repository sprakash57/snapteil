import { useState } from "react";
import Header from "@/components/Header";
import Feed from "@/features/feed/Feed";
import UploadModal from "@/features/upload/UploadModal";
import NewPostBanner from "@/features/notifications/NewPostBanner";
import { useImages } from "@/features/feed/useImages";
import { useInfiniteScroll } from "@/features/feed/useInfiniteScroll";
import { usePendingImages } from "@/features/notifications/usePendingImages";
import Loader from "@/components/Loader";
import type { Image } from "@/lib/types";

function App() {
  const [showUpload, setShowUpload] = useState(false);
  const [filterTags, setFilterTags] = useState<string[]>([]);
  const {
    images,
    allTags,
    loading,
    initialLoading,
    hasMore,
    error,
    loadMore,
    addImage,
    addTags,
    clearError,
  } = useImages(filterTags);

  const scrollRef = useInfiniteScroll(loadMore, hasMore, loading);
  // Pending images from websocket that haven't been added to feed yet (due to filters or own uploads)
  const {
    pendingImages,
    handleBannerClick,
    clear: clearPending,
    removePending,
  } = usePendingImages(images, filterTags, addImage);

  function handleTagsChange(tags: string[]) {
    setFilterTags(tags);
    clearPending();
  }

  function handleUpload(img: Image) {
    addTags(img.tags);
    if (
      filterTags.length === 0 ||
      filterTags.some((t) => img.tags.includes(t))
    ) {
      addImage(img);
    }
    removePending(img.id);
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
          loading={loading}
          error={error}
          onRetry={() => {
            clearError();
            loadMore();
          }}
          allTags={allTags}
          selectedTags={filterTags}
          onTagsChange={handleTagsChange}
        />
      )}

      {showUpload && (
        <UploadModal
          onClose={() => setShowUpload(false)}
          onUploaded={handleUpload}
        />
      )}
    </div>
  );
}

export default App;
