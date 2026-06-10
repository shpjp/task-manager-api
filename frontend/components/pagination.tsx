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

  return (
    <nav
      aria-label="Pagination"
      className="flex flex-col items-center justify-between gap-3 sm:flex-row"
    >
      <p className="text-sm text-slate-500">
        Showing <span className="font-medium text-slate-700">{start}–{end}</span> of{" "}
        <span className="font-medium text-slate-700">{meta.total}</span> tasks
      </p>
      <div className="flex items-center gap-2">
        <button
          onClick={() => onPageChange(meta.page - 1)}
          disabled={meta.page <= 1}
          className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
        >
          Previous
        </button>
        <span className="text-sm tabular-nums text-slate-600">
          {meta.page} / {meta.total_pages}
        </span>
        <button
          onClick={() => onPageChange(meta.page + 1)}
          disabled={meta.page >= meta.total_pages}
          className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
        >
          Next
        </button>
      </div>
    </nav>
  );
}
