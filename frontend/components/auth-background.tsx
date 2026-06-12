"use client";

import type { ReactNode } from "react";

export function AuthBackground({ children }: { children: ReactNode }) {
  return (
    <div className="relative flex min-h-full flex-1 overflow-hidden bg-[var(--background)]">
      {/* Gradient orbs */}
      <div
        aria-hidden
        className="pointer-events-none absolute -left-32 top-0 size-[28rem] rounded-full bg-indigo-400/20 blur-3xl dark:bg-indigo-500/10"
      />
      <div
        aria-hidden
        className="pointer-events-none absolute -right-24 bottom-0 size-[24rem] rounded-full bg-violet-400/15 blur-3xl dark:bg-violet-500/10"
      />
      <div
        aria-hidden
        className="pointer-events-none absolute left-1/2 top-1/3 size-64 -translate-x-1/2 rounded-full bg-sky-300/10 blur-3xl dark:bg-white/5"
      />

      {/* Grid overlay */}
      <div aria-hidden className="auth-grid pointer-events-none absolute inset-0" />

      <div className="relative flex min-h-full flex-1 flex-col lg:flex-row">
        {/* Landing-style left panel (desktop) */}
        <aside className="hidden flex-1 flex-col justify-between border-r border-[var(--border)] bg-[var(--surface-muted)]/50 p-10 lg:flex dark:bg-[var(--surface)]/80">
          <div>
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

          {/* Decorative preview cards */}
          <div className="relative mt-10 max-w-md space-y-3">
            <div className="rounded-xl border border-[var(--border)] bg-[var(--surface)] p-4 shadow-sm">
              <div className="flex items-center gap-2">
                <span className="size-2 rounded-full bg-emerald-500" />
                <span className="text-xs font-medium text-neutral-500">In progress</span>
              </div>
              <p className="mt-2 text-sm font-semibold">Ship release notes</p>
              <p className="mt-1 text-xs text-neutral-500">Due Friday · High priority</p>
            </div>
            <div className="ml-6 rounded-xl border border-[var(--border)] bg-[var(--surface)] p-4 opacity-80 shadow-sm">
              <div className="flex items-center gap-2">
                <span className="size-2 rounded-full bg-amber-500" />
                <span className="text-xs font-medium text-neutral-500">To do</span>
              </div>
              <p className="mt-2 text-sm font-semibold">Review pull requests</p>
            </div>
            <div className="absolute -right-4 -top-4 rounded-lg border border-[var(--border)] bg-[var(--surface)] px-3 py-2 text-xs shadow-lg">
              ✓ 12 tasks done this week
            </div>
          </div>
        </aside>

        {/* Form column */}
        <div className="relative flex flex-1 flex-col">{children}</div>
      </div>
    </div>
  );
}
