interface NewPostBannerProps {
  count: number;
  onClick: () => void;
}

export default function NewPostBanner({ count, onClick }: NewPostBannerProps) {
  return (
    <div className="fixed top-20 left-1/2 -translate-x-1/2 z-50">
      <button
        onClick={onClick}
        className="flex items-center gap-2 px-5 py-2.5 rounded-2xl text-sm font-medium text-[#e0e5ec] bg-[#5d31cf] cursor-pointer shadow-[0_14px_30px_rgba(39,49,66,0.35)]"
      >
        <span className="text-base">↑</span>
        {count === 1 ? "New post available" : `${count} new posts available`}
      </button>
    </div>
  );
}
