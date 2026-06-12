"use client";

const WORDS = [
  "deadlines",
  "tickets",
  "tasks",
  "assign",
  "complete",
  "todo",
  "weekend",
  "weekday",
  "ship",
  "sprint",
  "backlog",
  "standup",
];

interface WordMarqueeProps {
  direction?: "left" | "right";
  speed?: "slow" | "normal" | "fast";
  fullBleed?: boolean;
  className?: string;
}

export function WordMarquee({
  direction = "left",
  speed = "normal",
  fullBleed = false,
  className = "",
}: WordMarqueeProps) {
  const duration =
    speed === "slow" ? "20s" : speed === "fast" ? "10s" : "14s";
  const track = [...WORDS, ...WORDS];

  return (
    <div
      className={`overflow-hidden ${fullBleed ? "" : "marquee-mask"} ${className}`}
      aria-hidden
    >
      <div
        className={`marquee-track flex w-max gap-4 ${direction === "right" ? "marquee-to-right" : ""}`}
        style={{ animationDuration: duration }}
      >
        {track.map((word, i) => (
          <span
            key={`${word}-${i}`}
            className="shrink-0 rounded-full border border-[var(--border)] bg-[var(--surface)]/70 px-5 py-2 text-xs font-semibold uppercase tracking-wider text-[var(--brand)] backdrop-blur-md dark:bg-neutral-950/60"
          >
            {word}
          </span>
        ))}
      </div>
    </div>
  );
}
