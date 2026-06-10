"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useCallback, useEffect, useRef, useState } from "react";
import * as api from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import type {
  ListMeta,
  SortBy,
  SortOrder,
  Task,
  TaskInput,
  TaskStatus,
} from "@/lib/types";
import { Pagination } from "@/components/pagination";
import { TaskFormModal } from "@/components/task-form-modal";
import { TaskItem } from "@/components/task-item";
import {
  EmptyState,
  ErrorBanner,
  FullPageSpinner,
  Spinner,
} from "@/components/ui";

const PAGE_SIZE = 10;

interface Filters {
  status: TaskStatus | "";
  search: string;
  sortBy: SortBy;
  order: SortOrder;
  page: number;
}

function filtersFromParams(params: URLSearchParams): Filters {
  const status = params.get("status");
  const sortBy = params.get("sort_by");
  const order = params.get("order");
  const page = Number(params.get("page"));
  return {
    status: status === "todo" || status === "in_progress" || status === "done" ? status : "",
    search: params.get("search") ?? "",
    sortBy: sortBy === "due_date" || sortBy === "priority" ? sortBy : "created_at",
    order: order === "asc" ? "asc" : "desc",
    page: Number.isInteger(page) && page > 0 ? page : 1,
  };
}

function paramsFromFilters(filters: Filters): string {
  const params = new URLSearchParams();
  if (filters.status) params.set("status", filters.status);
  if (filters.search) params.set("search", filters.search);
  if (filters.sortBy !== "created_at") params.set("sort_by", filters.sortBy);
  if (filters.order !== "desc") params.set("order", filters.order);
  if (filters.page > 1) params.set("page", String(filters.page));
  return params.toString();
}

function TasksPageInner() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user, initializing, logout } = useAuth();

  const [filters, setFilters] = useState<Filters>(() =>
    filtersFromParams(new URLSearchParams(searchParams.toString()))
  );
  const [searchInput, setSearchInput] = useState(filters.search);

  const [tasks, setTasks] = useState<Task[]>([]);
  const [meta, setMeta] = useState<ListMeta | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [busyTaskId, setBusyTaskId] = useState<number | null>(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const fetchSeq = useRef(0);

  // Redirect to login when there is no valid session.
  useEffect(() => {
    if (!initializing && !user) router.replace("/login");
  }, [initializing, user, router]);

  // Keep the URL in sync so refresh/back preserve filters.
  useEffect(() => {
    const qs = paramsFromFilters(filters);
    router.replace(qs ? `/tasks?${qs}` : "/tasks", { scroll: false });
  }, [filters, router]);

  // Debounce free-text search.
  useEffect(() => {
    const handle = setTimeout(() => {
      setFilters((f) =>
        f.search === searchInput ? f : { ...f, search: searchInput, page: 1 }
      );
    }, 350);
    return () => clearTimeout(handle);
  }, [searchInput]);

  const loadTasks = useCallback(async () => {
    const seq = ++fetchSeq.current;
    setLoading(true);
    setError("");
    try {
      const res = await api.listTasks({
        status: filters.status,
        search: filters.search,
        sort_by: filters.sortBy,
        order: filters.order,
        page: filters.page,
        limit: PAGE_SIZE,
      });
      if (seq !== fetchSeq.current) return;
      setTasks(res.data);
      setMeta(res.meta);
      // If deleting emptied the current page, step back one page.
      if (res.data.length === 0 && res.meta.total > 0 && filters.page > 1) {
        setFilters((f) => ({ ...f, page: Math.max(1, res.meta.total_pages) }));
      }
    } catch (err) {
      if (seq !== fetchSeq.current) return;
      if (err instanceof api.ApiError && err.status === 401) {
        router.replace("/login");
        return;
      }
      setError(err instanceof api.ApiError ? err.message : "Failed to load tasks");
    } finally {
      if (seq === fetchSeq.current) setLoading(false);
    }
  }, [filters, router]);

  useEffect(() => {
    if (!initializing && user) loadTasks();
  }, [initializing, user, loadTasks]);

  if (initializing || !user) return <FullPageSpinner />;

  const setFilter = (patch: Partial<Filters>) =>
    setFilters((f) => ({ ...f, ...patch, page: patch.page ?? 1 }));

  async function handleToggleComplete(task: Task) {
    setBusyTaskId(task.id);
    try {
      await api.updateTask(task.id, {
        status: task.status === "done" ? "todo" : "done",
      });
      await loadTasks();
    } catch {
      setError("Failed to update task. Please try again.");
    } finally {
      setBusyTaskId(null);
    }
  }

  async function handleDelete(task: Task) {
    if (!window.confirm(`Delete "${task.title}"? This cannot be undone.`)) return;
    setBusyTaskId(task.id);
    try {
      await api.deleteTask(task.id);
      await loadTasks();
    } catch {
      setError("Failed to delete task. Please try again.");
    } finally {
      setBusyTaskId(null);
    }
  }

  async function handleModalSubmit(input: TaskInput) {
    if (editingTask) {
      await api.updateTask(editingTask.id, input);
    } else {
      await api.createTask(input);
    }
    await loadTasks();
  }

  const hasActiveFilters = Boolean(filters.status || filters.search);

  const selectClass =
    "rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-indigo-500";

  return (
    <div className="flex flex-1 flex-col">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-3xl items-center justify-between gap-4 px-4 py-4">
          <h1 className="text-lg font-bold tracking-tight text-slate-800">Taskflow</h1>
          <div className="flex items-center gap-3">
            <span className="hidden text-sm text-slate-500 sm:inline">{user.name}</span>
            <button
              onClick={async () => {
                await logout();
                router.replace("/login");
              }}
              className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 transition hover:bg-slate-50"
            >
              Log out
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto w-full max-w-3xl flex-1 px-4 py-6">
        <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center">
          <div className="relative flex-1">
            <svg
              className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-slate-400"
              fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}
            >
              <path strokeLinecap="round" d="m21 21-4.3-4.3M17 11a6 6 0 1 1-12 0 6 6 0 0 1 12 0Z" />
            </svg>
            <input
              type="search"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              placeholder="Search tasks by title…"
              aria-label="Search tasks by title"
              className="w-full rounded-lg border border-slate-300 bg-white py-2 pl-9 pr-3 text-sm text-slate-800 outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100"
            />
          </div>
          <button
            onClick={() => {
              setEditingTask(null);
              setModalOpen(true);
            }}
            className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-700"
          >
            + New task
          </button>
        </div>

        <div className="mb-5 flex flex-wrap items-center gap-2">
          <select
            value={filters.status}
            onChange={(e) => setFilter({ status: e.target.value as TaskStatus | "" })}
            aria-label="Filter by status"
            className={selectClass}
          >
            <option value="">All statuses</option>
            <option value="todo">To do</option>
            <option value="in_progress">In progress</option>
            <option value="done">Done</option>
          </select>

          <select
            value={filters.sortBy}
            onChange={(e) => setFilter({ sortBy: e.target.value as SortBy })}
            aria-label="Sort by"
            className={selectClass}
          >
            <option value="created_at">Sort: Created</option>
            <option value="due_date">Sort: Due date</option>
            <option value="priority">Sort: Priority</option>
          </select>

          <button
            onClick={() => setFilter({ order: filters.order === "asc" ? "desc" : "asc" })}
            aria-label={`Order ${filters.order === "asc" ? "ascending" : "descending"}`}
            className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-50"
          >
            {filters.order === "asc" ? "↑ Asc" : "↓ Desc"}
          </button>

          {hasActiveFilters && (
            <button
              onClick={() => {
                setSearchInput("");
                setFilter({ status: "", search: "" });
              }}
              className="text-sm font-medium text-indigo-600 hover:underline"
            >
              Clear filters
            </button>
          )}
        </div>

        {error && !loading ? (
          <ErrorBanner message={error} onRetry={loadTasks} />
        ) : loading && tasks.length === 0 ? (
          <div className="flex justify-center py-20">
            <Spinner className="size-7" />
          </div>
        ) : tasks.length === 0 ? (
          <EmptyState
            title={hasActiveFilters ? "No tasks match your filters" : "No tasks yet"}
            hint={
              hasActiveFilters
                ? "Try changing the search or status filter."
                : "Create your first task to get started."
            }
            action={
              !hasActiveFilters && (
                <button
                  onClick={() => {
                    setEditingTask(null);
                    setModalOpen(true);
                  }}
                  className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-700"
                >
                  + New task
                </button>
              )
            }
          />
        ) : (
          <>
            <ul className={`space-y-3 ${loading ? "opacity-60" : ""}`}>
              {tasks.map((task) => (
                <TaskItem
                  key={task.id}
                  task={task}
                  busy={busyTaskId === task.id}
                  onToggleComplete={handleToggleComplete}
                  onEdit={(t) => {
                    setEditingTask(t);
                    setModalOpen(true);
                  }}
                  onDelete={handleDelete}
                />
              ))}
            </ul>
            {meta && (
              <div className="mt-6">
                <Pagination
                  meta={meta}
                  onPageChange={(page) => setFilters((f) => ({ ...f, page }))}
                />
              </div>
            )}
          </>
        )}
      </main>

      {modalOpen && (
        <TaskFormModal
          task={editingTask}
          onClose={() => setModalOpen(false)}
          onSubmit={handleModalSubmit}
        />
      )}
    </div>
  );
}

export default function TasksPage() {
  return (
    <Suspense fallback={<FullPageSpinner />}>
      <TasksPageInner />
    </Suspense>
  );
}
