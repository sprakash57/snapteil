import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import UploadModal from "@/features/upload/UploadModal";
import type { Image } from "@/lib/types";

const baseImage: Image = {
  id: "42",
  title: "Test upload",
  tags: ["nature"],
  filename: "photo.jpg",
  url: "/uploads/photo.jpg",
  width: 1280,
  height: 720,
  createdAt: "2024-01-01T00:00:00Z",
};

function renderModal(
  overrides: Partial<{
    onClose: () => void;
    onUploaded: (i: Image) => void;
  }> = {},
) {
  const props = {
    onClose: vi.fn(),
    onUploaded: vi.fn(),
    ...overrides,
  };
  return { ...render(<UploadModal {...props} />), ...props };
}

/** Bypasses native HTML5 constraint validation to test React handler logic. */
function submitForm(container: HTMLElement) {
  fireEvent.submit(container.querySelector("form")!);
}

const mockFile = new File(["data"], "photo.jpg", { type: "image/jpeg" });

/**
 * jsdom's FormData(form) does not propagate files from <input type="file">.
 * Pass title/tags/file to control what the React handler receives without
 * relying on jsdom's native FormData serialisation.
 */
function spyFormData({
  title = "",
  tags = "",
  file = mockFile,
}: {
  title?: string;
  tags?: string;
  file?: File | null;
} = {}) {
  const original = FormData.prototype.get;
  const spy = vi.spyOn(FormData.prototype, "get").mockImplementation(function (
    this: FormData,
    name: string,
  ) {
    if (name === "title") return title;
    if (name === "tags") return tags;
    if (name === "file") return file ?? new File([], "");
    return original.call(this, name);
  });
  return spy;
}

describe("UploadModal", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  // ── Rendering ─────────────────────────────────────────────────────────────

  it("renders title, tags and file inputs", () => {
    renderModal();
    expect(screen.getByLabelText("Title")).toBeInTheDocument();
    expect(screen.getByLabelText("Tags")).toBeInTheDocument();
    expect(document.querySelector('input[type="file"]')).toBeInTheDocument();
  });

  it("renders Upload and Cancel buttons", () => {
    renderModal();
    expect(screen.getByRole("button", { name: "Upload" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Cancel" })).toBeInTheDocument();
  });

  // ── Keyboard / dismiss ────────────────────────────────────────────────────

  it("closes when Escape key is pressed", async () => {
    const user = userEvent.setup();
    const { onClose } = renderModal();
    await user.keyboard("{Escape}");
    expect(onClose).toHaveBeenCalledOnce();
  });

  it("closes when Cancel button is clicked", async () => {
    const user = userEvent.setup();
    const { onClose } = renderModal();
    await user.click(screen.getByRole("button", { name: "Cancel" }));
    expect(onClose).toHaveBeenCalledOnce();
  });

  // ── Validation ────────────────────────────────────────────────────────────

  it("shows error when submitting without a title", () => {
    const { container } = renderModal();
    const spy = spyFormData({ title: "" });
    submitForm(container);
    spy.mockRestore();
    expect(screen.getByRole("alert")).toHaveTextContent("Title is required");
    expect(vi.mocked(fetch)).not.toHaveBeenCalled();
  });

  it("shows error when no file is selected", () => {
    const { container } = renderModal();
    const spy = spyFormData({ title: "My picture", file: null });
    submitForm(container);
    spy.mockRestore();
    expect(screen.getByRole("alert")).toHaveTextContent(
      "Please select an image",
    );
    expect(vi.mocked(fetch)).not.toHaveBeenCalled();
  });

  it("shows error when more than 5 tags are entered", () => {
    const { container } = renderModal();
    const spy = spyFormData({ title: "My picture", tags: "a, b, c, d, e, f" });
    submitForm(container);
    spy.mockRestore();
    expect(screen.getByRole("alert")).toHaveTextContent(
      "Maximum 5 tags allowed",
    );
  });

  it("shows error for a tag with invalid characters", () => {
    const { container } = renderModal();
    const spy = spyFormData({
      title: "My picture",
      tags: "good-tag, bad tag!",
    });
    submitForm(container);
    spy.mockRestore();
    expect(screen.getByRole("alert")).toHaveTextContent("is invalid");
  });

  it("shows error for a tag exceeding 24 characters", () => {
    const { container } = renderModal();
    const spy = spyFormData({
      title: "My picture",
      tags: "a-tag-that-is-way-too-long-for-the-limit",
    });
    submitForm(container);
    spy.mockRestore();
    expect(screen.getByRole("alert")).toHaveTextContent("too long");
  });

  // ── Successful submit ─────────────────────────────────────────────────────

  it("submits form and calls onUploaded + onClose on success", async () => {
    const { container, onUploaded, onClose } = renderModal();

    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(baseImage),
    } as Response);

    const spy = spyFormData({ title: "Test upload", tags: "nature" });
    submitForm(container);
    spy.mockRestore();

    await waitFor(() => {
      expect(onUploaded).toHaveBeenCalledWith(baseImage);
      expect(onClose).toHaveBeenCalled();
    });
  });

  // ── API error ─────────────────────────────────────────────────────────────

  it("shows server error message on failed upload", async () => {
    const { container } = renderModal();

    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({ message: "File type not supported" }),
    } as Response);

    const spy = spyFormData({ title: "Test upload" });
    submitForm(container);
    spy.mockRestore();

    await waitFor(() => {
      expect(screen.getByRole("alert")).toHaveTextContent(
        "File type not supported",
      );
    });
  });

  // ── Loading state ─────────────────────────────────────────────────────────

  it("shows Loader dots while upload is in progress", async () => {
    const { container } = renderModal();

    // Never resolves — keeps upload state active indefinitely
    vi.mocked(fetch).mockReturnValue(new Promise(() => {}));

    const spy = spyFormData({ title: "Test upload" });
    submitForm(container);
    spy.mockRestore();

    await waitFor(() => {
      expect(container.querySelector(".animate-loader")).toBeInTheDocument();
    });
  });
});
