"use client";

import type { TaskStatus } from "@/lib/types";

const ITEMS: { label: string; status: TaskStatus | "" }[] = [
  { label: "All", status: "" },
  { label: "To do", status: "todo" },
  { label: "In progress", status: "in_progress" },
  { label: "Done", status: "done" },
];

interface MobileFilterNavProps {
  activeStatus: TaskStatus | "";
  scopeAll: boolean;
  isAdmin: boolean;
  onFilterStatus: (status: TaskStatus | "") => void;
  onToggleScopeAll: () => void;
}

export function MobileFilterNav({
  activeStatus,
  scopeAll,
  isAdmin,
  onFilterStatus,
  onToggleScopeAll,
}: MobileFilterNavProps) {
  return (
    <nav className="flex gap-2 overflow-x-auto border-b border-[var(--border)] px-4 py-2.5 lg:hidden [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
      {ITEMS.map((item) => {
        const active = activeStatus === item.status && !scopeAll;
        return (
          <button
            key={item.label}
            onClick={() => onFilterStatus(item.status)}
            className={`shrink-0 rounded-full px-3.5 py-1.5 text-xs font-semibold transition ${
              active
                ? "bg-indigo-600 text-white"
                : "border border-[var(--border)] bg-[var(--surface)] text-neutral-600 dark:text-neutral-300"
            }`}
          >
            {item.label}
          </button>
        );
      })}
      {isAdmin && (
        <button
          onClick={onToggleScopeAll}
          className={`shrink-0 rounded-full px-3.5 py-1.5 text-xs font-semibold transition ${
            scopeAll
              ? "bg-indigo-600 text-white"
              : "border border-[var(--border)] bg-[var(--surface)] text-neutral-600 dark:text-neutral-300"
          }`}
        >
          All users
        </button>
      )}
    </nav>
  );
}
