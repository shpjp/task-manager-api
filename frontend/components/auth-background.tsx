"use client";

import type { ReactNode } from "react";

export function AuthBackground({ children }: { children: ReactNode }) {
  return (
    <div className="relative flex min-h-full flex-1 overflow-hidden bg-[var(--background)]">
      {/* Animated gradient orbs */}
      <div aria-hidden className="auth-orb auth-orb-1 pointer-events-none absolute -left-32 top-0 size-[28rem] rounded-full bg-indigo-400/25 blur-3xl dark:bg-indigo-500/15" />
      <div aria-hidden className="auth-orb auth-orb-2 pointer-events-none absolute -right-24 bottom-0 size-[24rem] rounded-full bg-violet-400/20 blur-3xl dark:bg-violet-500/12" />
      <div aria-hidden className="auth-orb auth-orb-3 pointer-events-none absolute left-[calc(50%-8rem)] top-1/3 size-64 rounded-full bg-sky-300/15 blur-3xl dark:bg-white/8" />
      <div aria-hidden className="auth-orb auth-orb-4 pointer-events-none absolute right-1/4 top-16 size-40 rounded-full bg-fuchsia-300/15 blur-2xl dark:bg-fuchsia-500/10" />

      {/* Drifting grid */}
      <div aria-hidden className="auth-grid auth-grid-drift pointer-events-none absolute inset-0" />

      {/* Floating particles */}
      <div aria-hidden className="auth-particle auth-particle-1 pointer-events-none absolute left-[12%] top-[18%] size-2 rounded-full bg-indigo-400/60 dark:bg-indigo-400/40" />
      <div aria-hidden className="auth-particle auth-particle-2 pointer-events-none absolute left-[78%] top-[28%] size-1.5 rounded-full bg-violet-400/50 dark:bg-violet-400/35" />
      <div aria-hidden className="auth-particle auth-particle-3 pointer-events-none absolute left-[65%] top-[72%] size-2.5 rounded-full bg-sky-400/40 dark:bg-sky-300/30" />
      <div aria-hidden className="auth-particle auth-particle-4 pointer-events-none absolute left-[22%] top-[65%] size-1 rounded-full bg-indigo-500/50 dark:bg-indigo-300/40" />

      <div className="relative flex min-h-full flex-1 flex-col lg:flex-row">
        {/* Landing-style left panel */}
        <aside className="relative hidden flex-1 flex-col justify-between overflow-hidden border-r border-[var(--border)] bg-[var(--surface-muted)]/50 p-10 lg:flex dark:bg-[var(--surface)]/80">
          <div className="relative z-10">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-indigo-600 dark:text-indigo-400">
              Taskflow
            </p>
            <h2 className="mt-4 max-w-md text-3xl font-bold tracking-tight text-[var(--foreground)]">
              Plan work, track progress, ship on time.
            </h2>
            <p className="mt-3 max-w-sm text-sm leading-relaxed text-neutral-500 dark:text-neutral-400">
              A focused task manager with search, filters, due dates, and real-time
              updates — built for individuals and teams getting things done.
            </p>
          </div>

          {/* Animated preview cards */}
          <div className="relative z-10 mt-10 max-w-md space-y-3">
            <div className="auth-card-float auth-card-float-1 rounded-xl border border-[var(--border)] bg-[var(--surface)] p-4 shadow-sm">
              <div className="flex items-center gap-2">
                <span className="size-2 rounded-full bg-emerald-500" />
                <span className="text-xs font-medium text-neutral-500">In progress</span>
              </div>
              <p className="mt-2 text-sm font-semibold">Ship release notes</p>
              <p className="mt-1 text-xs text-neutral-500">Due Friday · High priority</p>
            </div>
            <div className="auth-card-float auth-card-float-2 ml-6 rounded-xl border border-[var(--border)] bg-[var(--surface)] p-4 opacity-80 shadow-sm">
              <div className="flex items-center gap-2">
                <span className="size-2 rounded-full bg-amber-500" />
                <span className="text-xs font-medium text-neutral-500">To do</span>
              </div>
              <p className="mt-2 text-sm font-semibold">Review pull requests</p>
            </div>
            <div className="auth-card-float auth-card-float-3 absolute -right-4 -top-4 rounded-lg border border-[var(--border)] bg-[var(--surface)] px-3 py-2 text-xs shadow-lg">
              ✓ 12 tasks done this week
            </div>
          </div>
        </aside>

        {/* Mobile decorative strip */}
        <div
          aria-hidden
          className="relative overflow-hidden border-b border-[var(--border)] px-4 py-6 lg:hidden"
        >
          <p className="relative z-10 text-center text-xs font-semibold uppercase tracking-[0.2em] text-indigo-600 dark:text-indigo-400">
            Taskflow
          </p>
          <p className="relative z-10 mt-2 text-center text-sm font-medium text-neutral-600 dark:text-neutral-400">
            Your tasks, organized.
          </p>
        </div>

        <div className="relative flex flex-1 flex-col">{children}</div>
      </div>
    </div>
  );
}
