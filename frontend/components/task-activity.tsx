"use client";

import { useEffect, useState } from "react";
import * as api from "@/lib/api";
import type { TaskActivity } from "@/lib/types";
import { Spinner } from "./ui";

function formatTimestamp(iso: string): string {
  return new Date(iso).toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function TaskActivityLog({ taskId }: { taskId: number }) {
  const [activity, setActivity] = useState<TaskActivity[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;
    api
      .listActivity(taskId)
      .then((list) => {
        if (!cancelled) setActivity(list);
      })
      .catch(() => {
        if (!cancelled) setError("Failed to load activity");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [taskId]);

  return (
    <div>
      <h3 className="mb-2 text-sm font-medium text-slate-700 dark:text-slate-300">
        Activity
      </h3>

      {error && <p className="text-xs font-medium text-red-600 dark:text-red-400">{error}</p>}

      {loading ? (
        <div className="flex justify-center py-3">
          <Spinner className="size-4" />
        </div>
      ) : activity.length === 0 ? (
        <p className="text-xs text-slate-400 dark:text-slate-500">No activity yet.</p>
      ) : (
        <ol className="max-h-40 space-y-2 overflow-y-auto pr-1">
          {activity.map((entry) => (
            <li key={entry.id} className="flex gap-2 text-xs">
              <span className="mt-1 size-1.5 shrink-0 rounded-full bg-[var(--brand)]" />
              <div className="min-w-0">
                <p className="text-slate-600 dark:text-slate-300">{entry.detail}</p>
                <p className="text-slate-400 dark:text-slate-500">
                  {entry.user?.name ?? "Someone"} · {formatTimestamp(entry.created_at)}
                </p>
              </div>
            </li>
          ))}
        </ol>
      )}
    </div>
  );
}
