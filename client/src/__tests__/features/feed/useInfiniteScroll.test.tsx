import { render, act } from "@testing-library/react";
import { useInfiniteScroll } from "@/features/feed/useInfiniteScroll";

// Wrap the hook in a component so the ref gets attached to a real DOM element
function Sentinel({
  loadMore,
  hasMore,
  loading,
}: {
  loadMore: () => void;
  hasMore: boolean;
  loading: boolean;
}) {
  const ref = useInfiniteScroll(loadMore, hasMore, loading);
  return <div ref={ref} data-testid="sentinel" />;
}

describe("useInfiniteScroll", () => {
  let observerCallback: IntersectionObserverCallback;
  const observe = vi.fn();
  const disconnect = vi.fn();

  beforeEach(() => {
    observe.mockClear();
    disconnect.mockClear();

    // Arrow functions cannot be used as constructors; use a class instead
    class MockIntersectionObserver {
      constructor(cb: IntersectionObserverCallback) {
        observerCallback = cb;
      }
      observe = observe;
      disconnect = disconnect;
    }

    vi.stubGlobal("IntersectionObserver", MockIntersectionObserver);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("attaches an observer to the sentinel element", () => {
    render(<Sentinel loadMore={vi.fn()} hasMore={true} loading={false} />);
    expect(observe).toHaveBeenCalledOnce();
  });

  it("calls loadMore when sentinel intersects and hasMore is true", () => {
    const loadMore = vi.fn();
    render(<Sentinel loadMore={loadMore} hasMore={true} loading={false} />);

    act(() => {
      observerCallback(
        [{ isIntersecting: true } as IntersectionObserverEntry],
        {} as IntersectionObserver,
      );
    });

    expect(loadMore).toHaveBeenCalledOnce();
  });

  it("does not call loadMore when hasMore is false", () => {
    const loadMore = vi.fn();
    render(<Sentinel loadMore={loadMore} hasMore={false} loading={false} />);

    act(() => {
      observerCallback(
        [{ isIntersecting: true } as IntersectionObserverEntry],
        {} as IntersectionObserver,
      );
    });

    expect(loadMore).not.toHaveBeenCalled();
  });

  it("does not call loadMore while loading is true", () => {
    const loadMore = vi.fn();
    render(<Sentinel loadMore={loadMore} hasMore={true} loading={true} />);

    act(() => {
      observerCallback(
        [{ isIntersecting: true } as IntersectionObserverEntry],
        {} as IntersectionObserver,
      );
    });

    expect(loadMore).not.toHaveBeenCalled();
  });

  it("does not call loadMore when sentinel is not intersecting", () => {
    const loadMore = vi.fn();
    render(<Sentinel loadMore={loadMore} hasMore={true} loading={false} />);

    act(() => {
      observerCallback(
        [{ isIntersecting: false } as IntersectionObserverEntry],
        {} as IntersectionObserver,
      );
    });

    expect(loadMore).not.toHaveBeenCalled();
  });

  it("disconnects the observer on unmount", () => {
    const { unmount } = render(
      <Sentinel loadMore={vi.fn()} hasMore={true} loading={false} />,
    );
    unmount();
    expect(disconnect).toHaveBeenCalledOnce();
  });
});
