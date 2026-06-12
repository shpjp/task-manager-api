"use client";

import type { ListMeta } from "@/lib/types";

export function Pagination({
  meta,
  onPageChange,
}: {
  meta: ListMeta;
  onPageChange: (page: number) => void;
}) {
  if (meta.total_pages <= 1) return null;

  const start = (meta.page - 1) * meta.limit + 1;
  const end = Math.min(meta.page * meta.limit, meta.total);

  const buttonClass =
    "rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40 dark:border-neutral-800 dark:bg-neutral-950 dark:text-neutral-200 dark:hover:bg-neutral-900";

  return (
    <nav
      aria-label="Pagination"
      className="flex flex-col items-center justify-between gap-3 sm:flex-row"
    >
      <p className="text-sm text-slate-500 dark:text-slate-400">
        Showing{" "}
        <span className="font-medium text-slate-700 dark:text-slate-200">
          {start}–{end}
        </span>{" "}
        of{" "}
        <span className="font-medium text-slate-700 dark:text-slate-200">
          {meta.total}
        </span>{" "}
        tasks
      </p>
      <div className="flex items-center gap-2">
        <button
          onClick={() => onPageChange(meta.page - 1)}
          disabled={meta.page <= 1}
          className={buttonClass}
        >
          Previous
        </button>
        <span className="text-sm tabular-nums text-slate-600 dark:text-slate-300">
          {meta.page} / {meta.total_pages}
        </span>
        <button
          onClick={() => onPageChange(meta.page + 1)}
          disabled={meta.page >= meta.total_pages}
          className={buttonClass}
        >
          Next
        </button>
      </div>
    </nav>
  );
}
