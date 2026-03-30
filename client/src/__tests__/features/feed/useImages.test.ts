import { renderHook, waitFor, act } from "@testing-library/react";
import { useImages } from "@/features/feed/useImages";
import type { Image, PaginatedResponse } from "@/lib/types";

function makeImage(id: string, tags: string[] = []): Image {
  return {
    id,
    title: `Image ${id}`,
    tags,
    filename: `img${id}.jpg`,
    url: `/uploads/img${id}.jpg`,
    width: 800,
    height: 600,
    createdAt: "2024-01-01T00:00:00Z",
  };
}

function mockPage(
  images: Image[],
  overrides: Partial<PaginatedResponse> = {},
): Response {
  const data: PaginatedResponse = {
    images,
    total: images.length,
    page: 1,
    perPage: 5,
    hasMore: false,
    ...overrides,
  };
  return { ok: true, json: () => Promise.resolve(data) } as Response;
}

describe("useImages", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  // ── Initial load ───────────────────────────────────────────────────────────

  it("starts with initialLoading true", () => {
    vi.mocked(fetch).mockReturnValue(new Promise(() => {}));
    const { result } = renderHook(() => useImages());
    expect(result.current.initialLoading).toBe(true);
  });

  it("loads images from page 1 on mount", async () => {
    const images = [makeImage("1"), makeImage("2")];
    vi.mocked(fetch).mockResolvedValue(mockPage(images));

    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));

    expect(result.current.images).toEqual(images);
    expect(vi.mocked(fetch)).toHaveBeenCalledWith(
      expect.stringContaining("page=1"),
    );
  });

  it("sets hasMore based on the API response", async () => {
    vi.mocked(fetch).mockResolvedValue(
      mockPage([makeImage("1")], { hasMore: true }),
    );

    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));

    expect(result.current.hasMore).toBe(true);
  });

  it("collects tags from loaded images into allTags", async () => {
    const images = [
      makeImage("1", ["nature", "landscape"]),
      makeImage("2", ["portrait"]),
    ];
    vi.mocked(fetch).mockResolvedValue(mockPage(images));

    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));

    expect(result.current.allTags).toEqual(["landscape", "nature", "portrait"]);
  });

  it("sets error and clears images on fetch failure", async () => {
    vi.mocked(fetch).mockResolvedValue({ ok: false } as Response);

    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));

    expect(result.current.error).toBeTruthy();
    expect(result.current.images).toHaveLength(0);
  });

  it("sets error when fetch throws a network error", async () => {
    vi.mocked(fetch).mockRejectedValue(new Error("Network error"));

    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));

    expect(result.current.error).toBeTruthy();
  });

  // ── Filter tags ────────────────────────────────────────────────────────────

  it("appends tags query param when filterTags are provided", async () => {
    vi.mocked(fetch).mockResolvedValue(mockPage([]));

    renderHook(() => useImages(["nature", "portrait"]));
    await waitFor(() => expect(vi.mocked(fetch)).toHaveBeenCalled());

    const url = vi.mocked(fetch).mock.calls[0][0] as string;
    expect(decodeURIComponent(url)).toContain("tags=nature,portrait");
  });

  it("refetches from page 1 when filterTags change", async () => {
    vi.mocked(fetch).mockResolvedValue(mockPage([]));

    const { rerender } = renderHook(
      ({ tags }: { tags: string[] }) => useImages(tags),
      { initialProps: { tags: [] } },
    );

    await waitFor(() => expect(vi.mocked(fetch)).toHaveBeenCalledTimes(1));

    rerender({ tags: ["nature"] });
    await waitFor(() => expect(vi.mocked(fetch)).toHaveBeenCalledTimes(2));

    const secondUrl = vi.mocked(fetch).mock.calls[1][0] as string;
    expect(secondUrl).toContain("page=1");
  });

  // ── loadMore ───────────────────────────────────────────────────────────────

  it("loadMore appends the next page of images", async () => {
    vi.mocked(fetch)
      .mockResolvedValueOnce(
        mockPage([makeImage("1")], { hasMore: true, page: 1 }),
      )
      .mockResolvedValueOnce(
        mockPage([makeImage("2")], { hasMore: false, page: 2 }),
      );

    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));
    expect(result.current.images).toHaveLength(1);

    await act(async () => {
      await result.current.loadMore();
    });

    expect(result.current.images).toHaveLength(2);
    expect(result.current.hasMore).toBe(false);
  });

  it("loadMore does not fetch when hasMore is false", async () => {
    vi.mocked(fetch).mockResolvedValue(mockPage([], { hasMore: false }));

    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));

    const callsBefore = vi.mocked(fetch).mock.calls.length;
    await act(async () => {
      await result.current.loadMore();
    });

    expect(vi.mocked(fetch)).toHaveBeenCalledTimes(callsBefore);
  });

  it("loadMore deduplicates images already present", async () => {
    const shared = makeImage("1");
    vi.mocked(fetch)
      .mockResolvedValueOnce(mockPage([shared], { hasMore: true }))
      .mockResolvedValueOnce(
        mockPage([shared, makeImage("2")], { hasMore: false }),
      );

    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));

    await act(async () => {
      await result.current.loadMore();
    });

    expect(result.current.images).toHaveLength(2);
  });

  // ── addImage ───────────────────────────────────────────────────────────────

  it("addImage prepends an image to the feed", async () => {
    vi.mocked(fetch).mockResolvedValue(mockPage([]));
    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));

    const newImage = makeImage("99", ["urban"]);
    act(() => {
      result.current.addImage(newImage);
    });

    expect(result.current.images[0]).toEqual(newImage);
    expect(result.current.allTags).toContain("urban");
  });

  // ── addTags ────────────────────────────────────────────────────────────────

  it("addTags grows allTags without duplicates", async () => {
    vi.mocked(fetch).mockResolvedValue(mockPage([]));
    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.initialLoading).toBe(false));

    act(() => {
      result.current.addTags(["architecture", "urban"]);
    });
    expect(result.current.allTags).toContain("architecture");
    expect(result.current.allTags).toContain("urban");

    act(() => {
      result.current.addTags(["architecture"]);
    });
    expect(
      result.current.allTags.filter((t) => t === "architecture"),
    ).toHaveLength(1);
  });

  // ── clearError ─────────────────────────────────────────────────────────────

  it("clearError removes the error state", async () => {
    vi.mocked(fetch).mockResolvedValue({ ok: false } as Response);
    const { result } = renderHook(() => useImages());
    await waitFor(() => expect(result.current.error).not.toBeNull());

    act(() => {
      result.current.clearError();
    });
    expect(result.current.error).toBeNull();
  });
});
