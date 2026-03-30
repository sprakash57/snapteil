import { renderHook, act } from "@testing-library/react";
import { usePendingImages } from "@/features/notifications/usePendingImages";
import { useWebSocket } from "@/features/websocket/useWebSocket";
import type { Image } from "@/lib/types";

vi.mock("@/features/websocket/useWebSocket");

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

describe("usePendingImages", () => {
  let wsCallback: (img: Image) => void;

  beforeEach(() => {
    vi.mocked(useWebSocket).mockImplementation((cb) => {
      wsCallback = cb;
    });
    vi.stubGlobal("scrollTo", vi.fn());
    Object.defineProperty(window, "scrollY", {
      configurable: true,
      writable: true,
      value: 200,
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
    vi.unstubAllGlobals();
  });

  // ── Queuing ────────────────────────────────────────────────────────────────

  it("starts with no pending images", () => {
    const { result } = renderHook(() => usePendingImages([], [], vi.fn()));
    expect(result.current.pendingImages).toHaveLength(0);
  });

  it("queues a new image received via WebSocket", () => {
    const { result } = renderHook(() => usePendingImages([], [], vi.fn()));
    act(() => {
      wsCallback(makeImage("1"));
    });
    expect(result.current.pendingImages).toHaveLength(1);
  });

  it("does not queue an image that is already in the feed", () => {
    const existingImage = makeImage("1");
    const { result } = renderHook(() =>
      usePendingImages([existingImage], [], vi.fn()),
    );
    act(() => {
      wsCallback(existingImage);
    });
    expect(result.current.pendingImages).toHaveLength(0);
  });

  it("does not queue a duplicate pending image", () => {
    const img = makeImage("1");
    const { result } = renderHook(() => usePendingImages([], [], vi.fn()));
    act(() => {
      wsCallback(img);
      wsCallback(img);
    });
    expect(result.current.pendingImages).toHaveLength(1);
  });

  it("skips images that do not match the active filter", () => {
    const { result } = renderHook(() =>
      usePendingImages([], ["nature"], vi.fn()),
    );
    act(() => {
      wsCallback(makeImage("1", ["street"]));
    });
    expect(result.current.pendingImages).toHaveLength(0);
  });

  it("queues images that match at least one active filter tag", () => {
    const { result } = renderHook(() =>
      usePendingImages([], ["nature", "street"], vi.fn()),
    );
    act(() => {
      wsCallback(makeImage("1", ["nature"]));
    });
    expect(result.current.pendingImages).toHaveLength(1);
  });

  it("queues all images when no filter is active", () => {
    const { result } = renderHook(() => usePendingImages([], [], vi.fn()));
    act(() => {
      wsCallback(makeImage("1", ["anything"]));
    });
    expect(result.current.pendingImages).toHaveLength(1);
  });

  // ── Banner click / flush ───────────────────────────────────────────────────

  it("handleBannerClick calls addImage for each pending image", () => {
    const addImage = vi.fn();
    const { result } = renderHook(() => usePendingImages([], [], addImage));

    act(() => {
      wsCallback(makeImage("1"));
      wsCallback(makeImage("2"));
    });
    expect(result.current.pendingImages).toHaveLength(2);

    act(() => {
      result.current.handleBannerClick();
    });

    expect(addImage).toHaveBeenCalledTimes(2);
    expect(result.current.pendingImages).toHaveLength(0);
  });

  it("handleBannerClick scrolls to the top", () => {
    const { result } = renderHook(() => usePendingImages([], [], vi.fn()));
    act(() => {
      result.current.handleBannerClick();
    });
    expect(window.scrollTo).toHaveBeenCalledWith({
      top: 0,
      behavior: "smooth",
    });
  });

  // ── clear ──────────────────────────────────────────────────────────────────

  it("clear empties the pending queue without calling addImage", () => {
    const addImage = vi.fn();
    const { result } = renderHook(() => usePendingImages([], [], addImage));

    act(() => {
      wsCallback(makeImage("1"));
    });
    act(() => {
      result.current.clear();
    });

    expect(result.current.pendingImages).toHaveLength(0);
    expect(addImage).not.toHaveBeenCalled();
  });

  // ── removePending ──────────────────────────────────────────────────────────

  it("removePending removes only the matching image", () => {
    const { result } = renderHook(() => usePendingImages([], [], vi.fn()));

    act(() => {
      wsCallback(makeImage("1"));
      wsCallback(makeImage("2"));
    });
    act(() => {
      result.current.removePending("1");
    });

    expect(result.current.pendingImages).toHaveLength(1);
    expect(result.current.pendingImages[0].id).toBe("2");
  });

  // ── Auto-flush on scroll to top ────────────────────────────────────────────

  it("auto-flushes pending images when user scrolls to scrollY = 0", () => {
    const addImage = vi.fn();
    const { result } = renderHook(() => usePendingImages([], [], addImage));

    act(() => {
      wsCallback(makeImage("1"));
    });
    expect(result.current.pendingImages).toHaveLength(1);

    act(() => {
      Object.defineProperty(window, "scrollY", {
        value: 0,
        configurable: true,
      });
      window.dispatchEvent(new Event("scroll"));
    });

    expect(addImage).toHaveBeenCalledOnce();
    expect(result.current.pendingImages).toHaveLength(0);
  });

  it("does not flush when scrolling to a non-zero position", () => {
    const addImage = vi.fn();
    const { result } = renderHook(() => usePendingImages([], [], addImage));

    act(() => {
      wsCallback(makeImage("1"));
    });

    act(() => {
      Object.defineProperty(window, "scrollY", {
        value: 100,
        configurable: true,
      });
      window.dispatchEvent(new Event("scroll"));
    });

    expect(addImage).not.toHaveBeenCalled();
    expect(result.current.pendingImages).toHaveLength(1);
  });
});
