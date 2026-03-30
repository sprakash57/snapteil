import { useEffect, useId, useRef, useState, type SyntheticEvent } from "react";
import type { Image } from "@/lib/types";
import Loader from "@/components/Loader";

interface UploadModalProps {
  onClose: () => void;
  onUploaded: (image: Image) => void;
}

export default function UploadModal({ onClose, onUploaded }: UploadModalProps) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState("");
  const titleInputRef = useRef<HTMLInputElement>(null);
  const titleId = useId();
  const tagsId = useId();
  const fileId = useId();
  const headingId = useId();
  const descriptionId = useId();
  const errorId = useId();

  useEffect(() => {
    const previouslyFocused = document.activeElement as HTMLElement | null;
    const originalOverflow = document.body.style.overflow;

    titleInputRef.current?.focus();
    document.body.style.overflow = "hidden";

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        onClose();
      }
    }

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      document.body.style.overflow = originalOverflow;
      window.removeEventListener("keydown", handleKeyDown);
      previouslyFocused?.focus();
    };
  }, [onClose]);

  async function handleSubmit(e: SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");

    const form = e.currentTarget;
    const formData = new FormData(form);

    const title = (formData.get("title") as string)?.trim();
    const file = formData.get("file") as File | null;

    if (!title) {
      setError("Title is required");
      return;
    }
    if (title.length > 100) {
      setError("Title must be 100 characters or fewer");
      return;
    }
    if (!file || file.size === 0) {
      setError("Please select an image");
      return;
    }

    const rawTags = (formData.get("tags") as string) ?? "";
    const tags = rawTags
      .split(",")
      .map((t) => t.trim().toLowerCase())
      .filter(Boolean);

    if (tags.length > 5) {
      setError("Maximum 5 tags allowed");
      return;
    }
    const invalidTag = tags.find((t) => !/^[a-z0-9-]+$/.test(t));
    if (invalidTag) {
      setError(
        `Tag "${invalidTag}" is invalid — use letters, numbers and hyphens only`,
      );
      return;
    }
    const longTag = tags.find((t) => t.length > 24);
    if (longTag) {
      setError(`Tag "${longTag}" is too long — max 24 characters`);
      return;
    }

    setUploading(true);
    try {
      const res = await fetch("/api/v1/uploads", {
        method: "POST",
        body: formData,
      });
      if (!res.ok) {
        const body = await res.json().catch(() => null);
        throw new Error(body?.message || "Upload failed");
      }
      const image: Image = await res.json();
      onUploaded(image);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setUploading(false);
    }
  }

  return (
    <div
      className="fixed inset-0 z-100 flex items-center justify-center bg-black/70"
      onClick={onClose}
    >
      <div
        className="w-full max-w-md mx-4 rounded-2xl bg-[#e0e5ec] p-6"
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-labelledby={headingId}
        aria-describedby={descriptionId}
      >
        <h2 id={headingId} className="text-lg font-semibold text-gray-700 mb-2">
          Upload Image
        </h2>
        <p id={descriptionId} className="mb-4 text-sm text-gray-500">
          Add a title, optional tags, and an image file to publish a new post.
        </p>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <label
              htmlFor={titleId}
              className="text-sm font-medium text-gray-600"
            >
              Title
            </label>
            <input
              ref={titleInputRef}
              id={titleId}
              name="title"
              placeholder="City lights after rain"
              maxLength={100}
              required
              aria-describedby={error ? errorId : undefined}
              className="px-4 py-2.5 rounded-xl bg-[#e0e5ec] text-gray-700 placeholder-gray-400 shadow-[inset_3px_3px_6px_#b8b9be,inset_-3px_-3px_6px_#ffffff] outline-none focus:shadow-[inset_4px_4px_8px_#b8b9be,inset_-4px_-4px_8px_#ffffff]"
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label
              htmlFor={tagsId}
              className="text-sm font-medium text-gray-600"
            >
              Tags
            </label>
            <input
              id={tagsId}
              name="tags"
              placeholder="nature, street-fight, portrait"
              aria-describedby={error ? errorId : undefined}
              className="w-full px-4 py-2.5 rounded-xl bg-[#e0e5ec] text-gray-700 placeholder-gray-400 shadow-[inset_3px_3px_6px_#b8b9be,inset_-3px_-3px_6px_#ffffff] outline-none focus:shadow-[inset_4px_4px_8px_#b8b9be,inset_-4px_-4px_8px_#ffffff]"
            />
            <p className="mt-1 text-xs text-gray-400">
              Letters, numbers and hyphens only. Max 24 chars per tag.
            </p>
          </div>

          <div className="flex flex-col gap-1.5">
            <label
              htmlFor={fileId}
              className="text-sm font-medium text-gray-600"
            >
              Image File
            </label>
            <input
              id={fileId}
              name="file"
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp,image/avif,image/svg+xml"
              required
              aria-describedby={error ? errorId : undefined}
              className="text-sm p-2 text-gray-600 file:mr-3 file:px-4 file:py-2 file:rounded-xl file:border-0 file:bg-[#e0e5ec] file:text-gray-700 file:shadow-[3px_3px_6px_#b8b9be,-3px_-3px_6px_#ffffff] file:cursor-pointer"
            />
          </div>

          {error && (
            <p id={errorId} role="alert" className="text-red-500 text-sm">
              {error}
            </p>
          )}

          <div className="flex justify-end gap-3 mt-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 rounded-xl text-sm text-[#FF6347] bg-[#e0e5ec] shadow-[3px_3px_6px_#b8b9be,-3px_-3px_6px_#ffffff] hover:shadow-[1px_1px_3px_#b8b9be,-1px_-1px_3px_#ffffff] transition-shadow cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={uploading}
              className="px-5 py-2 rounded-xl text-sm font-medium text-white bg-gray-600 shadow-[3px_3px_6px_#b8b9be,-3px_-3px_6px_#ffffff] hover:bg-gray-700 disabled:opacity-50 transition-all cursor-pointer min-w-20 flex items-center justify-center gap-2"
            >
              {uploading ? <Loader size="sm" shadow={false} /> : "Upload"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
