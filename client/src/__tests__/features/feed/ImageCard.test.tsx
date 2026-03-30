import { render, screen } from "@testing-library/react";
import ImageCard from "@/features/feed/ImageCard";
import type { Image } from "@/lib/types";

const baseImage: Image = {
  id: "1",
  title: "Beautiful sunset",
  tags: ["nature", "sunset"],
  filename: "sunset.jpg",
  url: "/uploads/sunset.jpg",
  width: 1920,
  height: 1080,
  createdAt: "2024-01-01T00:00:00Z",
};

describe("ImageCard", () => {
  it("renders the image with correct src and alt", () => {
    render(<ImageCard image={baseImage} />);
    const img = screen.getByRole("img", { name: "Beautiful sunset" });
    expect(img).toHaveAttribute("src", "/uploads/sunset.jpg");
  });

  it("renders the title as a heading", () => {
    render(<ImageCard image={baseImage} />);
    expect(
      screen.getByRole("heading", { name: "Beautiful sunset" }),
    ).toBeInTheDocument();
  });

  it("renders all tag pills", () => {
    render(<ImageCard image={baseImage} />);
    expect(screen.getByText("nature")).toBeInTheDocument();
    expect(screen.getByText("sunset")).toBeInTheDocument();
  });

  it("does not render any tag pills when tags array is empty", () => {
    render(<ImageCard image={{ ...baseImage, tags: [] }} />);
    expect(screen.queryByText("nature")).not.toBeInTheDocument();
    expect(screen.queryByText("sunset")).not.toBeInTheDocument();
  });

  it("uses lazy loading on the image", () => {
    render(<ImageCard image={baseImage} />);
    expect(screen.getByRole("img")).toHaveAttribute("loading", "lazy");
  });
});
