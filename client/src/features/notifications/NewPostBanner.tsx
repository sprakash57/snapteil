interface NewPostBannerProps {
  count: number;
  onClick: () => void;
}

export default function NewPostBanner({ count, onClick }: NewPostBannerProps) {
  return (
    <div className="fixed top-20 left-1/2 -translate-x-1/2 z-50">
      <button
        onClick={onClick}
        className="flex items-center gap-2 px-5 py-2.5 rounded-2xl text-sm font-medium text-gray-700 bg-[#e0e5ec] shadow-[5px_5px_10px_#b8b9be,-5px_-5px_10px_#ffffff] hover:shadow-[3px_3px_6px_#b8b9be,-3px_-3px_6px_#ffffff] active:shadow-[inset_3px_3px_6px_#b8b9be,inset_-3px_-3px_6px_#ffffff] transition-all duration-200 cursor-pointer"
      >
        <span className="text-base">↑</span>
        {count === 1 ? "New post available" : `${count} new posts available`}
      </button>
    </div>
  );
}
