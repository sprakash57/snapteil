import { useCallback, useEffect, useRef, useState } from "react";
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
  const imagesRef = useRef(images);
  const filterTagsRef = useRef(filterTags);

  useEffect(() => {
    imagesRef.current = images;
  }, [images]);

  useEffect(() => {
    filterTagsRef.current = filterTags;
  }, [filterTags]);

  // Skip images already in the feed (uploaded by this client)
  // Also skip if filter is active and image doesn't match any active tag
  useWebSocket(
    useCallback((img: Image) => {
      if (imagesRef.current.some((i) => i.id === img.id)) return;
      const activeTags = filterTagsRef.current;
      if (
        activeTags.length > 0 &&
        !activeTags.some((t) => img.tags.includes(t))
      ) {
        return;
      }
      setPendingImages((prev) => {
        if (prev.some((p) => p.id === img.id)) return prev;
        return [img, ...prev];
      });
    }, []),
  );

  // Clear pending when filter changes — handled inside setFilterTags wrapper
  function handleTagsChange(tags: string[]) {
    setFilterTags(tags);
    setPendingImages([]);
  }

  function handleUpload(img: Image) {
    // Always register new tags so they appear in the dropdown
    addTags(img.tags);
    // Only add to visible feed if it matches the active filter (or no filter)
    const active = filterTagsRef.current;
    if (active.length === 0 || active.some((t) => img.tags.includes(t))) {
      addImage(img);
    }
    // Remove from pending in case WS arrived first
    setPendingImages((prev) => prev.filter((p) => p.id !== img.id));
  }

  function flushPending() {
    setPendingImages((prev) => {
      prev.forEach(addImage);
      return [];
    });
  }

  function handleBannerClick() {
    flushPending();
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  // Dismiss banner and load images when user scrolls to top
  useEffect(() => {
    function onScroll() {
      if (window.scrollY === 0) {
        setPendingImages((prev) => {
          if (prev.length === 0) return prev;
          prev.forEach(addImage);
          return [];
        });
      }
    }
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, [addImage]);

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
