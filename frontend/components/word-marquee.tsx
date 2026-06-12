"use client";

const WORDS = [
  "priority",
  "due date",
  "team work",
  "tasks",
  "calendar",
  "upcoming",
  "weekend",
  "ship",
  "collaborate",
  "assign",
  "complete",
  "backlog",
  "standup",
  "deadlines",
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
    speed === "slow" ? "40s" : speed === "fast" ? "22s" : "30s";
  const track = [...WORDS, ...WORDS];

  return (
    <div className={`overflow-hidden ${className}`} aria-hidden>
      <div
        className={`marquee-track flex w-max gap-3 ${reverse ? "marquee-reverse" : ""}`}
        style={{ animationDuration: duration }}
      >
        {track.map((word, i) => (
          <span
            key={`${word}-${i}`}
            className="shrink-0 rounded-full border border-[var(--border)] bg-[var(--surface)] px-4 py-1.5 text-xs font-semibold uppercase tracking-wider text-[var(--brand)]"
          >
            {word}
          </span>
        ))}
      </div>
    </div>
  );
}
