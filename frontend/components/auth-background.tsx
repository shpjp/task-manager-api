"use client";

import type { ReactNode } from "react";
import { BrandLogo } from "./brand-logo";
import { WordMarquee } from "./word-marquee";

const HIGHLIGHTS = [
  { title: "Priority & due dates", desc: "Sort by urgency and never miss a deadline." },
  { title: "Team calendar", desc: "See upcoming work across the whole squad." },
  { title: "Collaborate & ship", desc: "Assign tasks, complete todos, move fast." },
  { title: "Weekend to weekday", desc: "Plan sprints, standups, and weekly goals." },
];

export function AuthBackground({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-full flex-1 flex-col bg-[var(--background)]">
      <div className="relative flex min-h-0 flex-1 flex-col overflow-hidden lg:flex-row">
        {/* Background effects — contained, no overlap with marquee */}
        <div aria-hidden className="pointer-events-none absolute inset-0 overflow-hidden">
          <div className="auth-orb auth-orb-1 absolute -left-32 top-0 size-[28rem] rounded-full bg-[var(--brand-glow)] blur-3xl" />
          <div className="auth-orb auth-orb-2 absolute -right-24 top-1/4 size-[24rem] rounded-full bg-sky-400/15 blur-3xl dark:bg-sky-500/10" />
          <div className="auth-grid auth-grid-drift absolute inset-0" />
        </div>

        {/* Left marketing panel */}
        <aside className="relative hidden w-full max-w-xl flex-col justify-center border-r border-[var(--border)] bg-[var(--surface-muted)] p-10 lg:flex dark:bg-[var(--surface)]">
          <BrandLogo size="lg" />
          <p className="mt-6 text-xs font-semibold uppercase tracking-[0.25em] text-[var(--brand)]">
            Built for teams that ship
          </p>
          <h2 className="mt-4 max-w-md text-3xl font-bold leading-tight tracking-tight text-[var(--foreground)]">
            Assign tasks. Hit deadlines. Win the week.
          </h2>
          <p className="mt-3 max-w-sm text-sm leading-relaxed text-neutral-500 dark:text-neutral-400">
            tasktheteam keeps tickets, todos, and due dates in one fast workspace —
            from weekday standups to weekend wrap-ups.
          </p>

          <div className="relative mt-10 max-w-md space-y-3">
            <div className="auth-card-float auth-card-float-1 rounded-xl border border-[var(--border)] bg-[var(--surface)] p-4 shadow-sm">
              <div className="flex items-center gap-2">
                <span className="size-2 rounded-full bg-emerald-500" />
                <span className="text-xs font-medium text-neutral-500">In progress</span>
              </div>
              <p className="mt-2 text-sm font-semibold">Complete sprint tickets</p>
              <p className="mt-1 text-xs text-neutral-500">Due Friday · high priority</p>
            </div>
            <div className="auth-card-float auth-card-float-2 ml-8 rounded-xl border border-[var(--border)] bg-[var(--surface)] p-4 shadow-sm">
              <div className="flex items-center gap-2">
                <span className="size-2 rounded-full bg-[var(--brand)]" />
                <span className="text-xs font-medium text-neutral-500">Todo</span>
              </div>
              <p className="mt-2 text-sm font-semibold">Review weekend backlog</p>
            </div>
          </div>
        </aside>

        {/* Mobile header */}
        <div className="relative border-b border-[var(--border)] px-4 py-5 lg:hidden">
          <div className="flex flex-col items-center gap-2">
            <BrandLogo size="md" />
            <p className="text-center text-xs text-neutral-500 dark:text-neutral-400">
              Tasks · calendar · collaborate
            </p>
          </div>
        </div>

        {/* Right column: hero top + form */}
        <div className="relative flex min-h-0 flex-1 flex-col">
          <section className="relative hidden border-b border-[var(--border)] bg-[var(--surface)] px-10 py-8 lg:block">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-[var(--brand)]">
              Why teams choose tasktheteam
            </p>
            <h3 className="mt-2 text-xl font-bold tracking-tight text-[var(--foreground)]">
              One workspace for every deadline
            </h3>
            <div className="mt-5 grid grid-cols-2 gap-3">
              {HIGHLIGHTS.map((item) => (
                <div
                  key={item.title}
                  className="rounded-xl border border-[var(--border)] bg-[var(--surface-muted)] p-4 dark:bg-neutral-950"
                >
                  <p className="text-sm font-semibold text-[var(--foreground)]">{item.title}</p>
                  <p className="mt-1 text-xs leading-relaxed text-neutral-500 dark:text-neutral-400">
                    {item.desc}
                  </p>
                </div>
              ))}
            </div>
          </section>

          <div className="relative flex flex-1 flex-col justify-center">{children}</div>
        </div>
      </div>

      {/* Single bottom marquee — fixed slot in layout, no overlap */}
      <footer className="shrink-0 border-t border-[var(--border)] bg-[var(--surface)] py-3">
        <WordMarquee speed="normal" />
      </footer>
    </div>
  );
}
