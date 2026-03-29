import type { Image } from "../types";

const CARD_WIDTH = 468;
const CARD_HEIGHT = 580;

interface ImageCardProps {
  image: Image;
}

export default function ImageCard({ image }: ImageCardProps) {
  return (
    <article
      className="rounded-2xl bg-[#e0e5ec] shadow-[6px_6px_12px_#b8b9be,-6px_-6px_12px_#ffffff] overflow-hidden"
      style={{ width: CARD_WIDTH }}
    >
      <div
        className="flex items-center justify-center bg-[#d5dae3]"
        style={{ width: CARD_WIDTH, height: CARD_HEIGHT }}
      >
        <img
          src={image.url}
          alt={image.title}
          className="max-w-full max-h-full object-contain"
          loading="lazy"
        />
      </div>
      <div className="px-4 py-3">
        <h2 className="text-base font-semibold text-gray-700 truncate">
          {image.title}
        </h2>
        {image.tags.length > 0 && (
          <div className="mt-1.5 flex flex-wrap gap-1.5">
            {image.tags.map((tag) => (
              <span
                key={tag}
                className="px-2.5 py-0.5 text-xs rounded-full text-gray-600 bg-[#e0e5ec] shadow-[inset_2px_2px_4px_#b8b9be,inset_-2px_-2px_4px_#ffffff]"
              >
                {tag}
              </span>
            ))}
          </div>
        )}
      </div>
    </article>
  );
}
