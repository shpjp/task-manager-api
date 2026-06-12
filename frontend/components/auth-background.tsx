"use client";

import type { ReactNode } from "react";
import { BrandLogo } from "./brand-logo";
import { WordMarquee } from "./word-marquee";

export function AuthBackground({ children }: { children: ReactNode }) {
  return (
    <div className="relative flex min-h-full flex-1 overflow-hidden bg-[var(--background)]">
      <div aria-hidden className="auth-orb auth-orb-1 pointer-events-none absolute -left-32 top-0 size-[28rem] rounded-full bg-[var(--brand-glow)] blur-3xl" />
      <div aria-hidden className="auth-orb auth-orb-2 pointer-events-none absolute -right-24 bottom-0 size-[24rem] rounded-full bg-sky-400/15 blur-3xl dark:bg-sky-500/10" />
      <div aria-hidden className="auth-orb auth-orb-3 pointer-events-none absolute left-[calc(50%-8rem)] top-1/3 size-64 rounded-full bg-[var(--brand-glow)] blur-3xl" />
      <div aria-hidden className="auth-orb auth-orb-4 pointer-events-none absolute right-1/4 top-16 size-40 rounded-full bg-cyan-300/15 blur-2xl dark:bg-cyan-500/10" />

      <div aria-hidden className="auth-grid auth-grid-drift pointer-events-none absolute inset-0" />

      <div aria-hidden className="auth-particle auth-particle-1 pointer-events-none absolute left-[12%] top-[18%] size-2 rounded-full bg-[var(--brand)]/70" />
      <div aria-hidden className="auth-particle auth-particle-2 pointer-events-none absolute left-[78%] top-[28%] size-1.5 rounded-full bg-sky-400/50" />
      <div aria-hidden className="auth-particle auth-particle-3 pointer-events-none absolute left-[65%] top-[72%] size-2.5 rounded-full bg-cyan-400/40" />

      <div className="pointer-events-none absolute inset-x-0 top-[12%] z-10 space-y-3 opacity-90">
        <WordMarquee speed="slow" />
        <WordMarquee reverse speed="fast" />
      </div>

      <div className="pointer-events-none absolute inset-x-0 bottom-[8%] z-10 space-y-3 opacity-80">
        <WordMarquee speed="normal" />
        <WordMarquee reverse speed="slow" />
      </div>

      <div className="relative flex min-h-full flex-1 flex-col lg:flex-row">
        <aside className="relative hidden flex-1 flex-col justify-between overflow-hidden border-r border-[var(--border)] bg-[var(--surface-muted)]/60 p-10 lg:flex dark:bg-[var(--surface)]/90">
          <div className="relative z-10">
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
          </div>

          <div className="relative z-10 mt-10 max-w-md space-y-3">
            <div className="auth-card-float auth-card-float-1 rounded-xl border border-[var(--border)] bg-[var(--surface)] p-4 shadow-sm">
              <div className="flex items-center gap-2">
                <span className="size-2 rounded-full bg-emerald-500" />
                <span className="text-xs font-medium text-neutral-500">In progress</span>
              </div>
              <p className="mt-2 text-sm font-semibold">Complete sprint tickets</p>
              <p className="mt-1 text-xs text-neutral-500">Due Friday · assign @team</p>
            </div>
            <div className="auth-card-float auth-card-float-2 ml-6 rounded-xl border border-[var(--border)] bg-[var(--surface)] p-4 opacity-80 shadow-sm">
              <div className="flex items-center gap-2">
                <span className="size-2 rounded-full bg-[var(--brand)]" />
                <span className="text-xs font-medium text-neutral-500">Todo</span>
              </div>
              <p className="mt-2 text-sm font-semibold">Review weekend backlog</p>
            </div>
            <div className="auth-card-float auth-card-float-3 absolute -right-4 -top-4 rounded-lg border border-[var(--border)] bg-[var(--surface)] px-3 py-2 text-xs shadow-lg">
              ✓ Team velocity +24%
            </div>
          </div>
        </aside>

        <div aria-hidden className="relative overflow-hidden border-b border-[var(--border)] px-4 py-5 lg:hidden">
          <div className="relative z-10 flex flex-col items-center gap-2">
            <BrandLogo size="md" />
            <p className="text-center text-xs text-neutral-500 dark:text-neutral-400">
              Tasks · tickets · deadlines — all in one place
            </p>
          </div>
        </div>

        <div className="relative z-20 flex flex-1 flex-col">{children}</div>
      </div>
    </div>
  );
}
