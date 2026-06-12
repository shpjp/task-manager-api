"use client";

import Image from "next/image";

interface BrandLogoProps {
  size?: "sm" | "md" | "lg";
  showWordmark?: boolean;
  className?: string;
}

const sizes = {
  sm: { img: 28, text: "text-sm" },
  md: { img: 36, text: "text-base" },
  lg: { img: 48, text: "text-xl" },
};

export function BrandLogo({ size = "md", showWordmark = true, className = "" }: BrandLogoProps) {
  const s = sizes[size];
  return (
    <div className={`flex items-center gap-2.5 ${className}`}>
      <Image
        src="/logo.png"
        alt="tasktheteam logo"
        width={s.img}
        height={s.img}
        className="rounded-md"
        priority
      />
      {showWordmark && (
        <span className={`font-bold tracking-tight text-[var(--foreground)] ${s.text}`}>
          task<span className="text-[var(--brand)]">the</span>team
        </span>
      )}
    </div>
  );
}
