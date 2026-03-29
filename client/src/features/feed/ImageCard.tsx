import type { Image } from "../../lib/types";

interface ImageCardProps {
  image: Image;
}

export default function ImageCard({ image }: ImageCardProps) {
  return (
    <article className="w-full max-w-117 bg-[#e0e5ec] shadow-[6px_6px_12px_#b8b9be,-6px_-6px_12px_#ffffff] overflow-hidden">
      <div className="w-full bg-[#d5dae3]">
        <img
          src={image.url}
          alt={image.title}
          className="w-full h-auto block"
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
