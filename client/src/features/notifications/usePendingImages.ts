import { useCallback, useEffect, useRef, useState } from "react";
import { useWebSocket } from "@/features/websocket/useWebSocket";
import type { Image } from "@/lib/types";

function matchesFilters(image: Image, activeTags: string[]) {
  return (
    activeTags.length === 0 ||
    activeTags.some((tag) => image.tags.includes(tag))
  );
}

function addPendingImage(prev: Image[], nextImage: Image) {
  if (prev.some((pendingImage) => pendingImage.id === nextImage.id)) {
    return prev;
  }

  return [nextImage, ...prev];
}

export function usePendingImages(
  images: Image[],
  filterTags: string[],
  addImage: (img: Image) => void,
) {
  const [pendingImages, setPendingImages] = useState<Image[]>([]);
  const imagesRef = useRef(images);
  const filterTagsRef = useRef(filterTags);
  const addImageRef = useRef(addImage);

  useEffect(() => {
    imagesRef.current = images;
  }, [images]);

  useEffect(() => {
    filterTagsRef.current = filterTags;
  }, [filterTags]);

  useEffect(() => {
    addImageRef.current = addImage;
  }, [addImage]);

  const queuePendingImage = useCallback((img: Image) => {
    if (!matchesFilters(img, filterTagsRef.current)) return;
    if (imagesRef.current.some((existingImage) => existingImage.id === img.id)) {
      return;
    }

    setPendingImages((prev) => addPendingImage(prev, img));
  }, []);

  const flushPendingImages = useCallback(() => {
    setPendingImages((prev) => {
      if (prev.length === 0) return prev;

      for (const img of [...prev].reverse()) {
        addImageRef.current(img);
      }

      return [];
    });
  }, []);

  useWebSocket(queuePendingImage);

  // Auto-flush when user scrolls to top
  useEffect(() => {
    function onScroll() {
      if (window.scrollY === 0) {
        flushPendingImages();
      }
    }

    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, [flushPendingImages]);

  function clear() {
    setPendingImages([]);
  }

  function handleBannerClick() {
    flushPendingImages();
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  function removePending(id: string) {
    setPendingImages((prev) =>
      prev.filter((pendingImage) => pendingImage.id !== id),
    );
  }

  return { pendingImages, handleBannerClick, clear, removePending };
}
