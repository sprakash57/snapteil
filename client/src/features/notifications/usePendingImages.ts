import { useCallback, useEffect, useRef, useState } from "react";
import { useWebSocket } from "@/features/websocket/useWebSocket";
import type { Image } from "@/lib/types";

export function usePendingImages(
  images: Image[],
  filterTags: string[],
  addImage: (img: Image) => void,
) {
  const [pendingImages, setPendingImages] = useState<Image[]>([]);
  const imagesRef = useRef(images);
  const filterTagsRef = useRef(filterTags);

  useEffect(() => {
    imagesRef.current = images;
  }, [images]);

  useEffect(() => {
    filterTagsRef.current = filterTags;
  }, [filterTags]);

  // Queue incoming WS images as pending (skip own uploads + filtered-out tags)
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

  const flush = useCallback(() => {
    setPendingImages((prev) => {
      prev.forEach(addImage);
      return [];
    });
  }, [addImage]);

  const clear = useCallback(() => setPendingImages([]), []);

  const handleBannerClick = useCallback(() => {
    flush();
    window.scrollTo({ top: 0, behavior: "smooth" });
  }, [flush]);

  // Auto-flush when user scrolls to top
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

  const removePending = useCallback((id: string) => {
    setPendingImages((prev) => prev.filter((p) => p.id !== id));
  }, []);

  return { pendingImages, handleBannerClick, clear, removePending };
}
