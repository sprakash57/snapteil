import { useEffect, useState } from "react";

const VISIBILITY_OFFSET = 1200;

export default function ScrollToTopButton() {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    function handleScroll() {
      setVisible(window.scrollY > VISIBILITY_OFFSET);
    }

    handleScroll();
    window.addEventListener("scroll", handleScroll, { passive: true });

    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  function handleClick() {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <button
      type="button"
      onClick={handleClick}
      aria-label="Scroll to top"
      className={`fixed right-6 bottom-6 z-50 flex h-12 w-12 items-center justify-center rounded-full bg-[#5d31cf] text-[#e0e5ec] shadow-[0_14px_30px_rgba(39,49,66,0.35)] cursor-pointer transition-all ${
        visible
          ? "translate-y-0 opacity-100 pointer-events-auto"
          : "translate-y-3 opacity-0 pointer-events-none"
      } hover:-translate-y-0.5`}
    >
      <span className="text-2xl leading-none">↑</span>
    </button>
  );
}
