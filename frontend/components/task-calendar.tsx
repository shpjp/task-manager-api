"use client";

import { useMemo, useState } from "react";
import type { Task } from "@/lib/types";

interface TaskCalendarProps {
  tasks: Task[];
  layout?: "sidebar" | "inline";
  onSelectDate?: (isoDate: string) => void;
}

function dateKey(d: Date): string {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
}

function parseTaskDate(iso: string): string {
  const d = new Date(iso);
  return dateKey(new Date(d.getUTCFullYear(), d.getUTCMonth(), d.getUTCDate()));
}

function CalendarPanel({
  tasks,
  onSelectDate,
  compact = false,
}: {
  tasks: Task[];
  onSelectDate?: (isoDate: string) => void;
  compact?: boolean;
}) {
  const todayKey = dateKey(new Date());
  const [viewYear, setViewYear] = useState(() => new Date().getFullYear());
  const [viewMonth, setViewMonth] = useState(() => new Date().getMonth());

  const dueByDay = useMemo(() => {
    const map = new Map<string, number>();
    for (const task of tasks) {
      if (!task.due_date) continue;
      const key = parseTaskDate(task.due_date);
      map.set(key, (map.get(key) ?? 0) + 1);
    }
    return map;
  }, [tasks]);

  const monthLabel = new Date(viewYear, viewMonth, 1).toLocaleDateString(undefined, {
    month: "long",
    year: "numeric",
  });

  const firstDay = new Date(viewYear, viewMonth, 1).getDay();
  const daysInMonth = new Date(viewYear, viewMonth + 1, 0).getDate();

  const cells: (number | null)[] = [
    ...Array.from({ length: firstDay }, () => null),
    ...Array.from({ length: daysInMonth }, (_, i) => i + 1),
  ];
  while (cells.length % 7 !== 0) cells.push(null);

  function shiftMonth(delta: number) {
    const d = new Date(viewYear, viewMonth + delta, 1);
    setViewYear(d.getFullYear());
    setViewMonth(d.getMonth());
  }

  const upcoming = useMemo(() => {
    return tasks
      .filter((t) => t.due_date && t.status !== "done")
      .map((t) => ({ task: t, key: parseTaskDate(t.due_date!) }))
      .filter((x) => x.key >= todayKey)
      .sort((a, b) => a.key.localeCompare(b.key))
      .slice(0, compact ? 2 : 4);
  }, [tasks, todayKey, compact]);

  return (
    <>
      <div className="rounded-xl border border-[var(--border)] bg-[var(--surface-muted)] p-4 dark:bg-neutral-900">
        <div className="mb-3 flex items-center justify-between">
          <h2 className="text-sm font-semibold">{monthLabel}</h2>
          <div className="flex gap-1">
            <button
              type="button"
              onClick={() => shiftMonth(-1)}
              className="rounded-md px-2 py-1 text-xs text-neutral-500 hover:bg-[var(--surface)] dark:hover:bg-neutral-800"
              aria-label="Previous month"
            >
              ‹
            </button>
            <button
              type="button"
              onClick={() => shiftMonth(1)}
              className="rounded-md px-2 py-1 text-xs text-neutral-500 hover:bg-[var(--surface)] dark:hover:bg-neutral-800"
              aria-label="Next month"
            >
              ›
            </button>
          </div>
        </div>

        <div className="mb-1 grid grid-cols-7 gap-1 text-center text-[10px] font-medium text-neutral-500">
          {["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"].map((d) => (
            <span key={d}>{d}</span>
          ))}
        </div>

        <div className="grid grid-cols-7 gap-1">
          {cells.map((day, i) => {
            if (day === null) return <span key={`e-${i}`} />;
            const key = dateKey(new Date(viewYear, viewMonth, day));
            const count = dueByDay.get(key) ?? 0;
            const isToday = key === todayKey;
            return (
              <button
                key={key}
                type="button"
                onClick={() => onSelectDate?.(key)}
                className={`relative flex aspect-square items-center justify-center rounded-lg text-xs transition ${
                  isToday
                    ? "bg-indigo-600 font-semibold text-white"
                    : count > 0
                      ? "bg-indigo-100 font-medium text-indigo-800 dark:bg-indigo-950 dark:text-indigo-300"
                      : "text-neutral-600 hover:bg-[var(--surface)] dark:text-neutral-400 dark:hover:bg-neutral-800"
                }`}
              >
                {day}
                {count > 1 && (
                  <span className="absolute bottom-0.5 size-1 rounded-full bg-indigo-500" />
                )}
              </button>
            );
          })}
        </div>
      </div>

      {upcoming.length > 0 && (
        <div className="mt-4 rounded-xl border border-[var(--border)] bg-[var(--surface-muted)] p-4 dark:bg-neutral-900">
          <h3 className="mb-2 text-xs font-semibold uppercase tracking-wider text-neutral-500">
            Upcoming
          </h3>
          <ul className="space-y-2">
            {upcoming.map(({ task, key }) => (
              <li key={task.id} className="text-xs">
                <p className="truncate font-medium">{task.title}</p>
                <p className="text-neutral-500">{key}</p>
              </li>
            ))}
          </ul>
        </div>
      )}
    </>
  );
}

export function TaskCalendar({
  tasks,
  layout = "sidebar",
  onSelectDate,
}: TaskCalendarProps) {
  if (layout === "inline") {
    return (
      <div className="mb-5 xl:hidden">
        <CalendarPanel tasks={tasks} onSelectDate={onSelectDate} compact />
      </div>
    );
  }

  return (
    <aside className="hidden w-72 shrink-0 border-l border-[var(--border)] bg-[var(--surface)] xl:block">
      <div className="sticky top-0 max-h-screen overflow-y-auto p-4">
        <CalendarPanel tasks={tasks} onSelectDate={onSelectDate} />
      </div>
    </aside>
  );
}
