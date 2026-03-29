import { type SyntheticEvent, useRef, useState } from "react";
import type { Image } from "../types";

interface UploadModalProps {
  onClose: () => void;
  onUploaded: (image: Image) => void;
}

export default function UploadModal({ onClose, onUploaded }: UploadModalProps) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState("");
  const fileRef = useRef<HTMLInputElement>(null);

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
    if (!file || file.size === 0) {
      setError("Please select an image");
      return;
    }

    const rawTags = (formData.get("tags") as string) ?? "";
    const tagCount = rawTags
      .split(",")
      .map((t) => t.trim())
      .filter(Boolean).length;
    if (tagCount > 5) {
      setError("Maximum 5 tags allowed");
      return;
    }

    setUploading(true);
    try {
      const res = await fetch("/api/v1/images", {
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
      className="fixed inset-0 z-100 flex items-center justify-center bg-black/40"
      onClick={onClose}
    >
      <div
        className="w-full max-w-md mx-4 rounded-2xl bg-[#e0e5ec] p-6 shadow-[8px_8px_16px_#b8b9be,-8px_-8px_16px_#ffffff]"
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-lg font-semibold text-gray-700 mb-4">
          Upload Image
        </h2>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <input
            name="title"
            placeholder="Title"
            className="px-4 py-2.5 rounded-xl bg-[#e0e5ec] text-gray-700 placeholder-gray-400 shadow-[inset_3px_3px_6px_#b8b9be,inset_-3px_-3px_6px_#ffffff] outline-none focus:shadow-[inset_4px_4px_8px_#b8b9be,inset_-4px_-4px_8px_#ffffff]"
          />

          <input
            name="tags"
            placeholder="Tags (comma separated)"
            className="px-4 py-2.5 rounded-xl bg-[#e0e5ec] text-gray-700 placeholder-gray-400 shadow-[inset_3px_3px_6px_#b8b9be,inset_-3px_-3px_6px_#ffffff] outline-none focus:shadow-[inset_4px_4px_8px_#b8b9be,inset_-4px_-4px_8px_#ffffff]"
          />

          <input
            ref={fileRef}
            name="file"
            type="file"
            accept="image/jpeg,image/png,image/gif,image/webp,image/avif,image/svg+xml"
            className="text-sm text-gray-600 file:mr-3 file:px-4 file:py-2 file:rounded-xl file:border-0 file:bg-[#e0e5ec] file:text-gray-700 file:shadow-[3px_3px_6px_#b8b9be,-3px_-3px_6px_#ffffff] file:cursor-pointer"
          />

          {error && <p className="text-red-500 text-sm">{error}</p>}

          <div className="flex justify-end gap-3 mt-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 rounded-xl text-sm text-gray-600 bg-[#e0e5ec] shadow-[3px_3px_6px_#b8b9be,-3px_-3px_6px_#ffffff] hover:shadow-[1px_1px_3px_#b8b9be,-1px_-1px_3px_#ffffff] transition-shadow cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={uploading}
              className="px-5 py-2 rounded-xl text-sm font-medium text-white bg-gray-600 shadow-[3px_3px_6px_#b8b9be,-3px_-3px_6px_#ffffff] hover:bg-gray-700 disabled:opacity-50 transition-all cursor-pointer"
            >
              {uploading ? "Uploading..." : "Upload"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
