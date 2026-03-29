import LogoSvg from "../assets/logo.svg";

interface HeaderProps {
  onUploadClick: () => void;
}

export default function Header({ onUploadClick }: HeaderProps) {
  return (
    <header className="sticky top-0 z-50 bg-[#e0e5ec] shadow-[6px_6px_12px_#b8b9be,-6px_-6px_12px_#ffffff]">
      <div className="mx-auto max-w-3xl px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <img src={LogoSvg} alt="Snapteil Logo" width="32" height="32" />
          <h1 className="text-2xl italic font-semibold text-gray-700 tracking-wide">
            Snapteil
          </h1>
        </div>
        <button
          onClick={onUploadClick}
          className="px-5 py-2.5 rounded-xl text-sm font-medium text-gray-700 bg-[#e0e5ec] shadow-[4px_4px_8px_#b8b9be,-4px_-4px_8px_#ffffff] hover:shadow-[2px_2px_4px_#b8b9be,-2px_-2px_4px_#ffffff] active:shadow-[inset_3px_3px_6px_#b8b9be,inset_-3px_-3px_6px_#ffffff] transition-shadow cursor-pointer"
        >
          Upload
        </button>
      </div>
    </header>
  );
}
