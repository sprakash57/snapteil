import { useCallback, useEffect, useRef, useState } from "react";
import type { Image, PaginatedResponse } from "../types";

const PER_PAGE = 5;

export function useImages() {
  const [images, setImages] = useState<Image[]>([]);
  const [page, setPage] = useState(2);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const loadingRef = useRef(false);

  // Load first page (seeds + any in-memory uploads)
  useEffect(() => {
    async function loadInitial() {
      try {
        const res = await fetch(`/api/v1/images?page=1&perPage=${PER_PAGE}`);
        if (!res.ok) throw new Error("Failed to load images");
        const data: PaginatedResponse = await res.json();
        setImages(data.images);
        setHasMore(data.hasMore);
        setError(null);
      } catch (err) {
        console.error(err);
        setError("Unable to reach the server. Please check your connection.");
      } finally {
        setInitialLoading(false);
      }
    }
    loadInitial();
  }, []);

  // Load more images for infinite scroll
  const loadMore = useCallback(async () => {
    if (loadingRef.current || !hasMore) return;
    loadingRef.current = true;
    setLoading(true);

    try {
      const res = await fetch(
        `/api/v1/images?page=${page}&perPage=${PER_PAGE}`,
      );
      if (!res.ok) throw new Error("Failed to load images");
      const data: PaginatedResponse = await res.json();
      setImages((prev) => {
        const existingIds = new Set(prev.map((img) => img.id));
        const newImages = data.images.filter((img) => !existingIds.has(img.id));
        return [...prev, ...newImages];
      });
      setHasMore(data.hasMore);
      setPage((p) => p + 1);
    } catch (err) {
      console.error(err);
      setError("Unable to reach the server. Please check your connection.");
    } finally {
      setLoading(false);
      loadingRef.current = false;
    }
  }, [page, hasMore]);

  // Prepend a newly uploaded image
  const addImage = useCallback((img: Image) => {
    setImages((prev) => [img, ...prev]);
  }, []);

  const clearError = useCallback(() => setError(null), []);

  return {
    images,
    loading,
    initialLoading,
    hasMore,
    error,
    loadMore,
    addImage,
    clearError,
  };
}
