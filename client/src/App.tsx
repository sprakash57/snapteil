import { useState } from "react";
import Header from "./components/Header";
import Feed from "./components/Feed";
import UploadModal from "./components/UploadModal";
import { useImages } from "./hooks/useImages";
import { useInfiniteScroll } from "./hooks/useInfiniteScroll";
import Loader from "./components/common/Loader";

function App() {
  const [showUpload, setShowUpload] = useState(false);
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

  return (
    <div className="min-h-screen bg-[#e0e5ec] flex flex-col">
      <Header onUploadClick={() => setShowUpload(true)} />

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
