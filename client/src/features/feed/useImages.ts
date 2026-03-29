import { useCallback, useEffect, useRef, useState } from "react";
import type { Image, PaginatedResponse } from "@/lib/types";

const PER_PAGE = 5;

export function useImages(filterTags: string[] = []) {
  const [images, setImages] = useState<Image[]>([]);
  const [page, setPage] = useState(2);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [allTags, setAllTags] = useState<string[]>([]);
  const loadingRef = useRef(false);

  const tagsParam =
    filterTags.length > 0
      ? `&tags=${encodeURIComponent(filterTags.join(","))}`
      : "";

  // Re-fetch from page 1 when filter tags change
  useEffect(() => {
    let cancelled = false;
    setInitialLoading(true);
    setError(null);

    async function loadInitial() {
      try {
        const res = await fetch(
          `/api/v1/images?page=1&perPage=${PER_PAGE}${tagsParam}`,
        );
        if (!res.ok) throw new Error("Failed to load images");
        const data: PaginatedResponse = await res.json();
        if (!cancelled) {
          setImages(data.images);
          setHasMore(data.hasMore);
          setPage(2);
          // Grow tag list (never shrink — filtered results are a subset)
          setAllTags((prev) => {
            const set = new Set(prev);
            for (const img of data.images) {
              for (const t of img.tags) set.add(t);
            }
            return set.size === prev.length ? prev : Array.from(set).sort();
          });
        }
      } catch (err) {
        console.error(err);
        if (!cancelled) {
          setError("Unable to reach the server. Please check your connection.");
        }
      } finally {
        if (!cancelled) {
          setInitialLoading(false);
        }
      }
    }
    loadInitial();
    return () => {
      cancelled = true;
    };
  }, [tagsParam]);

  // Load more images for infinite scroll
  const loadMore = useCallback(async () => {
    if (loadingRef.current || !hasMore) return;
    loadingRef.current = true;
    setLoading(true);

    try {
      const res = await fetch(
        `/api/v1/images?page=${page}&perPage=${PER_PAGE}${tagsParam}`,
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
      // Grow tag list with newly loaded images
      setAllTags((prev) => {
        const set = new Set(prev);
        for (const img of data.images) {
          for (const t of img.tags) set.add(t);
        }
        return set.size === prev.length ? prev : Array.from(set).sort();
      });
    } catch (err) {
      console.error(err);
      setError("Unable to reach the server. Please check your connection.");
    } finally {
      setLoading(false);
      loadingRef.current = false;
    }
  }, [page, hasMore, tagsParam]);

  // Prepend a newly uploaded image
  const addImage = useCallback((img: Image) => {
    setImages((prev) => [img, ...prev]);
    setAllTags((prev) => {
      const set = new Set(prev);
      for (const t of img.tags) set.add(t);
      return set.size === prev.length ? prev : Array.from(set).sort();
    });
  }, []);

  const clearError = useCallback(() => setError(null), []);

  const addTags = useCallback((tags: string[]) => {
    setAllTags((prev) => {
      const set = new Set(prev);
      for (const t of tags) set.add(t);
      return set.size === prev.length ? prev : Array.from(set).sort();
    });
  }, []);

  return {
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
  };
}
