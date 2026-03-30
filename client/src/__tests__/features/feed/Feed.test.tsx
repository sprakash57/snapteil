import { createRef } from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Feed from "@/features/feed/Feed";
import type { Image } from "@/lib/types";

vi.mock("@/assets/error.svg", () => ({ default: "error.svg" }));

const scrollRef = createRef<HTMLDivElement | null>();

const defaultProps = {
  images: [] as Image[],
  scrollRef,
  loading: false,
  error: null,
  onRetry: vi.fn(),
  allTags: [] as string[],
  selectedTags: [] as string[],
  onTagsChange: vi.fn(),
};

const sampleImages: Image[] = [
  {
    id: "1",
    title: "Image One",
    tags: [],
    filename: "img1.jpg",
    url: "/uploads/img1.jpg",
    width: 800,
    height: 600,
    createdAt: "2024-01-01T00:00:00Z",
  },
  {
    id: "2",
    title: "Image Two",
    tags: ["nature"],
    filename: "img2.jpg",
    url: "/uploads/img2.jpg",
    width: 800,
    height: 600,
    createdAt: "2024-01-02T00:00:00Z",
  },
];

describe("Feed", () => {
  it("renders image cards for each image", () => {
    render(<Feed {...defaultProps} images={sampleImages} />);
    expect(
      screen.getByRole("heading", { name: "Image One" }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: "Image Two" }),
    ).toBeInTheDocument();
  });

  it("shows full-page error state when there is an error and no images", () => {
    render(<Feed {...defaultProps} error="Network error" images={[]} />);
    expect(screen.getByText("Something went wrong :(")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Retry" })).toBeInTheDocument();
  });

  it("calls onRetry when the Retry button is clicked", async () => {
    const onRetry = vi.fn();
    const user = userEvent.setup();
    render(
      <Feed
        {...defaultProps}
        error="Network error"
        images={[]}
        onRetry={onRetry}
      />,
    );
    await user.click(screen.getByRole("button", { name: "Retry" }));
    expect(onRetry).toHaveBeenCalledOnce();
  });

  it("shows inline error banner when error occurs after images are loaded", () => {
    render(
      <Feed
        {...defaultProps}
        images={sampleImages}
        error="Failed to load more"
      />,
    );
    expect(screen.getByText("Failed to load more")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: "Try again" }),
    ).toBeInTheDocument();
    // Images stay visible
    expect(
      screen.getByRole("heading", { name: "Image One" }),
    ).toBeInTheDocument();
  });

  it("shows Loader when loading is true", () => {
    const { container } = render(<Feed {...defaultProps} loading={true} />);
    expect(container.querySelector(".animate-loader")).toBeInTheDocument();
  });

  it("does not show loader when not loading", () => {
    const { container } = render(<Feed {...defaultProps} loading={false} />);
    expect(container.querySelector(".animate-loader")).not.toBeInTheDocument();
  });
});
