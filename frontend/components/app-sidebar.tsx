"use client";

import Link from "next/link";
import type { TaskStatus } from "@/lib/types";
import { BrandLogo } from "./brand-logo";
import { ThemeToggle } from "./theme-toggle";

interface AppSidebarProps {
  userName: string;
  isAdmin: boolean;
  totalTasks: number;
  activeStatus: TaskStatus | "";
  scopeAll: boolean;
  onFilterStatus: (status: TaskStatus | "") => void;
  onToggleScopeAll: () => void;
  onLogout: () => void;
}

const NAV: { label: string; status: TaskStatus | ""; icon: string }[] = [
  { label: "All tasks", status: "", icon: "◫" },
  { label: "To do", status: "todo", icon: "○" },
  { label: "In progress", status: "in_progress", icon: "◐" },
  { label: "Done", status: "done", icon: "●" },
];

export function AppSidebar({
  userName,
  isAdmin,
  totalTasks,
  activeStatus,
  scopeAll,
  onFilterStatus,
  onToggleScopeAll,
  onLogout,
}: AppSidebarProps) {
  return (
    <aside className="hidden w-60 shrink-0 flex-col border-r border-[var(--border)] bg-[var(--surface)] lg:sticky lg:top-0 lg:h-screen lg:overflow-y-auto lg:flex">
      <div className="border-b border-[var(--border)] px-4 py-4">
        <BrandLogo size="sm" />
        {isAdmin && (
          <span className="mt-2 inline-block rounded-full bg-sky-100 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-[var(--brand-dark)] dark:bg-sky-950 dark:text-[var(--brand)]">
            Admin
          </span>
        )}
      </div>

      <nav className="flex-1 space-y-1 p-3">
        <p className="px-2 pb-2 text-[10px] font-semibold uppercase tracking-wider text-neutral-500">
          Views
        </p>
        {NAV.map((item) => {
          const active = activeStatus === item.status && !scopeAll;
          return (
            <button
              key={item.label}
              onClick={() => onFilterStatus(item.status)}
              className={`flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-left text-sm font-medium transition ${
                active
                  ? "bg-[var(--brand)] text-white"
                  : "text-neutral-600 hover:bg-[var(--surface-muted)] dark:text-neutral-300 dark:hover:bg-neutral-900"
              }`}
            >
              <span className="text-base leading-none">{item.icon}</span>
              {item.label}
            </button>
          );
        })}

        {isAdmin && (
          <button
            onClick={onToggleScopeAll}
            className={`mt-2 flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-left text-sm font-medium transition ${
              scopeAll
                ? "bg-[var(--brand)] text-white"
                : "text-neutral-600 hover:bg-[var(--surface-muted)] dark:text-neutral-300 dark:hover:bg-neutral-900"
            }`}
          >
            <span className="text-base leading-none">◎</span>
            All users
          </button>
        )}
      </nav>

      <div className="border-t border-[var(--border)] p-4">
        <div className="mb-4 rounded-lg bg-[var(--surface-muted)] px-3 py-2 dark:bg-neutral-900">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-neutral-500">
            Total tasks
          </p>
          <p className="text-2xl font-bold tabular-nums">{totalTasks}</p>
        </div>

        <div className="flex items-center justify-between gap-2">
          <div className="min-w-0">
            <p className="truncate text-xs font-medium">{userName}</p>
            <Link href="/tasks" className="text-[10px] text-neutral-500 hover:underline">
              My workspace
            </Link>
          </div>
          <ThemeToggle />
        </div>

        <button
          onClick={onLogout}
          className="mt-3 w-full rounded-lg border border-[var(--border)] px-3 py-2 text-xs font-medium text-neutral-600 transition hover:bg-[var(--surface-muted)] dark:text-neutral-400 dark:hover:bg-neutral-900"
        >
          Log out
        </button>
      </div>
    </aside>
  );
}
