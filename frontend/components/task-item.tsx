"use client";

import type { Task } from "@/lib/types";

const PRIORITY_STYLES: Record<Task["priority"], string> = {
  high: "bg-red-100 text-red-700 dark:bg-red-950 dark:text-red-300",
  medium: "bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300",
  low: "bg-slate-100 text-slate-600 dark:bg-neutral-900 dark:text-neutral-300",
};

const STATUS_LABELS: Record<Task["status"], string> = {
  todo: "To do",
  in_progress: "In progress",
  done: "Done",
};

const STATUS_STYLES: Record<Task["status"], string> = {
  todo: "bg-slate-100 text-slate-600 dark:bg-neutral-900 dark:text-neutral-300",
  in_progress: "bg-blue-100 text-blue-700 dark:bg-blue-950 dark:text-blue-300",
  done: "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300",
};

function formatDueDate(iso: string): { label: string; overdue: boolean } {
  const date = new Date(iso);
  const label = date.toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return { label, overdue: date < today };
}

export function TaskItem({
  task,
  busy,
  readOnly = false,
  onToggleComplete,
  onEdit,
  onDelete,
}: {
  task: Task;
  busy: boolean;
  /** Admins viewing other users' tasks cannot modify them. */
  readOnly?: boolean;
  onToggleComplete: (task: Task) => void;
  onEdit: (task: Task) => void;
  onDelete: (task: Task) => void;
}) {
  const done = task.status === "done";
  const due = task.due_date ? formatDueDate(task.due_date) : null;

  return (
    <li className="group flex items-start gap-3 rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition hover:shadow-md dark:border-neutral-800 dark:bg-neutral-950">
      <input
        type="checkbox"
        checked={done}
        disabled={busy || readOnly}
        onChange={() => onToggleComplete(task)}
        aria-label={done ? "Mark as not done" : "Mark as done"}
        className="mt-1 size-5 shrink-0 cursor-pointer accent-emerald-600 disabled:cursor-not-allowed"
      />

      <div className="min-w-0 flex-1">
        <div className="flex flex-wrap items-center gap-2">
          <h3
            className={`truncate text-sm font-semibold sm:text-base ${
              done
                ? "text-slate-400 line-through dark:text-slate-500"
                : "text-slate-800 dark:text-slate-100"
            }`}
          >
            {task.title}
          </h3>
          <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${STATUS_STYLES[task.status]}`}>
            {STATUS_LABELS[task.status]}
          </span>
          <span
            className={`rounded-full px-2 py-0.5 text-xs font-medium capitalize ${PRIORITY_STYLES[task.priority]}`}
          >
            {task.priority}
          </span>
          {task.owner && (
            <span className="rounded-full bg-indigo-100 px-2 py-0.5 text-xs font-medium text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300">
              {task.owner.email}
            </span>
          )}
        </div>

        {task.description && (
          <p className="mt-1 line-clamp-2 text-sm text-slate-500 dark:text-slate-400">
            {task.description}
          </p>
        )}

        {due && (
          <p
            className={`mt-1.5 text-xs font-medium ${
              due.overdue && !done
                ? "text-red-600 dark:text-red-400"
                : "text-slate-400 dark:text-slate-500"
            }`}
          >
            Due {due.label}
            {due.overdue && !done && " · overdue"}
          </p>
        )}
      </div>

      {!readOnly && (
        <div className="flex shrink-0 gap-1">
          <button
            onClick={() => onEdit(task)}
            disabled={busy}
            className="rounded-lg p-2 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700 disabled:opacity-40 dark:hover:bg-neutral-900 dark:hover:text-neutral-200"
            aria-label={`Edit ${task.title}`}
          >
            <svg className="size-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M16.86 4.49a1.88 1.88 0 1 1 2.65 2.65L7.5 19.14l-3.71 1.06 1.06-3.7L16.86 4.5Z" />
            </svg>
          </button>
          <button
            onClick={() => onDelete(task)}
            disabled={busy}
            className="rounded-lg p-2 text-slate-400 transition hover:bg-red-50 hover:text-red-600 disabled:opacity-40 dark:hover:bg-red-950 dark:hover:text-red-400"
            aria-label={`Delete ${task.title}`}
          >
            <svg className="size-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 7h12M9 7V5a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2m-7 0 .7 12.3A1 1 0 0 0 9.7 20h4.6a1 1 0 0 0 1-0.7L16 7" />
            </svg>
          </button>
        </div>
      )}
    </li>
  );
}
