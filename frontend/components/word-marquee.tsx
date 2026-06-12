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
  reverse?: boolean;
  speed?: "slow" | "normal" | "fast";
  className?: string;
}

export function WordMarquee({
  reverse = false,
  speed = "normal",
  className = "",
}: WordMarqueeProps) {
  const duration =
    speed === "slow" ? "38s" : speed === "fast" ? "18s" : "26s";
  const track = [...WORDS, ...WORDS];

  return (
    <div
      className={`marquee-mask overflow-hidden ${className}`}
      aria-hidden
    >
      <div
        className={`marquee-track flex w-max gap-3 ${reverse ? "marquee-reverse" : ""}`}
        style={{ animationDuration: duration }}
      >
        {track.map((word, i) => (
          <span
            key={`${word}-${i}`}
            className="marquee-pill shrink-0 rounded-full border border-[var(--border)] bg-[var(--surface)]/80 px-4 py-1.5 text-xs font-semibold uppercase tracking-wider text-[var(--brand)] backdrop-blur-sm dark:bg-neutral-950/70"
          >
            {word}
          </span>
        ))}
      </div>
    </div>
  );
}
