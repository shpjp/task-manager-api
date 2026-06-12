"use client";

import { useEffect, useState } from "react";
import { ApiError } from "@/lib/api";
import type { Task, TaskInput, TaskPriority, TaskStatus } from "@/lib/types";
import { TaskActivityLog } from "./task-activity";
import { TaskAttachments } from "./task-attachments";
import { FieldError, Spinner } from "./ui";

interface FormValues {
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  due_date: string; // yyyy-mm-dd or ""
}

function initialValues(task?: Task | null): FormValues {
  return {
    title: task?.title ?? "",
    description: task?.description ?? "",
    status: task?.status ?? "todo",
    priority: task?.priority ?? "medium",
    due_date: task?.due_date ? task.due_date.slice(0, 10) : "",
  };
}

function validate(values: FormValues): Record<string, string> {
  const errors: Record<string, string> = {};
  if (!values.title.trim()) {
    errors.title = "Title is required";
  } else if (values.title.trim().length > 200) {
    errors.title = "Title must be at most 200 characters";
  }
  if (values.description.length > 2000) {
    errors.description = "Description must be at most 2000 characters";
  }
  if (values.due_date && Number.isNaN(Date.parse(values.due_date))) {
    errors.due_date = "Enter a valid date";
  }
  return errors;
}

export function TaskFormModal({
  task,
  onClose,
  onSubmit,
}: {
  /** When set, the modal edits this task; otherwise it creates a new one. */
  task?: Task | null;
  onClose: () => void;
  onSubmit: (input: TaskInput) => Promise<void>;
}) {
  const [values, setValues] = useState<FormValues>(() => initialValues(task));
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [formError, setFormError] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onClose]);

  const set = <K extends keyof FormValues>(key: K, value: FormValues[K]) =>
    setValues((v) => ({ ...v, [key]: value }));

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const clientErrors = validate(values);
    setErrors(clientErrors);
    setFormError("");
    if (Object.keys(clientErrors).length > 0) return;

    const input: TaskInput = {
      title: values.title.trim(),
      description: values.description.trim(),
      status: values.status,
      priority: values.priority,
      due_date: values.due_date
        ? new Date(`${values.due_date}T00:00:00`).toISOString()
        : null,
    };

    setSaving(true);
    try {
      await onSubmit(input);
      onClose();
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.fields) setErrors(err.fields);
        setFormError(err.fields ? "" : err.message);
      } else {
        setFormError("Something went wrong. Please try again.");
      }
    } finally {
      setSaving(false);
    }
  }

  const inputClass =
    "w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-800 outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100 dark:border-neutral-800 dark:bg-neutral-950 dark:text-neutral-100 dark:focus:ring-indigo-900";

  const labelClass = "mb-1 block text-sm font-medium text-slate-700 dark:text-slate-300";

  return (
    <div
      className="fixed inset-0 z-50 flex items-end justify-center bg-slate-900/40 p-0 sm:items-center sm:p-4 dark:bg-black/60"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-label={task ? "Edit task" : "New task"}
        className="max-h-[92vh] w-full max-w-lg overflow-y-auto rounded-t-2xl bg-white p-6 shadow-xl sm:rounded-2xl dark:bg-neutral-950"
      >
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-slate-800 dark:text-slate-100">
            {task ? "Edit task" : "New task"}
          </h2>
          <button
            onClick={onClose}
            aria-label="Close"
            className="rounded-lg p-1.5 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700 dark:hover:bg-neutral-900 dark:hover:text-neutral-200"
          >
            <svg className="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" d="M6 6l12 12M18 6L6 18" />
            </svg>
          </button>
        </div>

        <form onSubmit={handleSubmit} noValidate className="space-y-4">
          <div>
            <label htmlFor="title" className={labelClass}>
              Title <span className="text-red-500">*</span>
            </label>
            <input
              id="title"
              type="text"
              value={values.title}
              onChange={(e) => set("title", e.target.value)}
              placeholder="What needs to be done?"
              className={inputClass}
              autoFocus
            />
            <FieldError message={errors.title} />
          </div>

          <div>
            <label htmlFor="description" className={labelClass}>
              Description
            </label>
            <textarea
              id="description"
              value={values.description}
              onChange={(e) => set("description", e.target.value)}
              placeholder="Add more details (optional)"
              rows={3}
              className={`${inputClass} resize-none`}
            />
            <FieldError message={errors.description} />
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
            <div>
              <label htmlFor="status" className={labelClass}>
                Status
              </label>
              <select
                id="status"
                value={values.status}
                onChange={(e) => set("status", e.target.value as TaskStatus)}
                className={inputClass}
              >
                <option value="todo">To do</option>
                <option value="in_progress">In progress</option>
                <option value="done">Done</option>
              </select>
            </div>
            <div>
              <label htmlFor="priority" className={labelClass}>
                Priority
              </label>
              <select
                id="priority"
                value={values.priority}
                onChange={(e) => set("priority", e.target.value as TaskPriority)}
                className={inputClass}
              >
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
            </div>
            <div>
              <label htmlFor="due_date" className={labelClass}>
                Due date
              </label>
              <input
                id="due_date"
                type="date"
                value={values.due_date}
                onChange={(e) => set("due_date", e.target.value)}
                className={inputClass}
              />
              <FieldError message={errors.due_date} />
            </div>
          </div>

          {formError && (
            <p className="rounded-lg bg-red-50 px-3 py-2 text-sm font-medium text-red-700 dark:bg-red-950/40 dark:text-red-300">
              {formError}
            </p>
          )}

          <div className="flex justify-end gap-2 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="rounded-lg border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-50 dark:border-neutral-800 dark:bg-neutral-950 dark:text-neutral-200 dark:hover:bg-neutral-900"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={saving}
              className="flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-indigo-700 disabled:opacity-60"
            >
              {saving && <Spinner className="size-4 border-white/40 border-t-white" />}
              {task ? "Save changes" : "Create task"}
            </button>
          </div>
        </form>

        {task && (
          <div className="mt-6 space-y-5 border-t border-slate-200 pt-5 dark:border-neutral-800">
            <TaskAttachments taskId={task.id} />
            <TaskActivityLog taskId={task.id} />
          </div>
        )}
      </div>
    </div>
  );
}
