"use client";

import { useEffect, useRef, useState } from "react";
import * as api from "@/lib/api";
import type { Attachment } from "@/lib/types";
import { Spinner } from "./ui";

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

const ACCEPT =
  ".png,.jpg,.jpeg,.gif,.webp,.pdf,.txt,.md,.csv,.doc,.docx,.xls,.xlsx";

export function TaskAttachments({ taskId }: { taskId: number }) {
  const [attachments, setAttachments] = useState<Attachment[]>([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [dragOver, setDragOver] = useState(false);
  const [error, setError] = useState("");
  const fileInput = useRef<HTMLInputElement>(null);

  useEffect(() => {
    let cancelled = false;
    api
      .listAttachments(taskId)
      .then((list) => {
        if (!cancelled) setAttachments(list);
      })
      .catch(() => {
        if (!cancelled) setError("Failed to load attachments");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [taskId]);

  async function uploadFile(file: File) {
    setUploading(true);
    setError("");
    try {
      const attachment = await api.uploadAttachment(taskId, file);
      setAttachments((list) => [attachment, ...list]);
    } catch (err) {
      setError(err instanceof api.ApiError ? err.message : "Upload failed");
    } finally {
      setUploading(false);
    }
  }

  async function handleUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    e.target.value = "";
    if (!file) return;
    await uploadFile(file);
  }

  async function handleDelete(attachment: Attachment) {
    setError("");
    const previous = attachments;
    setAttachments((list) => list.filter((a) => a.id !== attachment.id));
    try {
      await api.deleteAttachment(taskId, attachment.id);
    } catch {
      setAttachments(previous);
      setError("Failed to delete attachment");
    }
  }

  function handleDragOver(e: React.DragEvent) {
    e.preventDefault();
    e.stopPropagation();
    setDragOver(true);
  }

  function handleDragLeave(e: React.DragEvent) {
    e.preventDefault();
    e.stopPropagation();
    setDragOver(false);
  }

  async function handleDrop(e: React.DragEvent) {
    e.preventDefault();
    e.stopPropagation();
    setDragOver(false);
    const file = e.dataTransfer.files?.[0];
    if (!file || uploading) return;
    await uploadFile(file);
  }

  return (
    <div>
      <div className="mb-2 flex items-center justify-between">
        <h3 className="text-sm font-medium text-slate-700 dark:text-slate-300">
          Attachments
        </h3>
        <button
          type="button"
          onClick={() => fileInput.current?.click()}
          disabled={uploading}
          className="flex items-center gap-1.5 rounded-lg border border-slate-300 bg-white px-2.5 py-1 text-xs font-medium text-slate-700 transition hover:bg-slate-50 disabled:opacity-60 dark:border-neutral-800 dark:bg-neutral-950 dark:text-neutral-200 dark:hover:bg-neutral-900"
        >
          {uploading ? <Spinner className="size-3" /> : "+"} Browse
        </button>
        <input
          ref={fileInput}
          type="file"
          hidden
          accept={ACCEPT}
          onChange={handleUpload}
        />
      </div>

      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => !uploading && fileInput.current?.click()}
        role="button"
        tabIndex={0}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            fileInput.current?.click();
          }
        }}
        className={`mb-3 cursor-pointer rounded-xl border-2 border-dashed px-4 py-5 text-center transition ${
          dragOver
            ? "border-[var(--brand)] bg-sky-50 dark:border-[var(--brand)] dark:bg-sky-950/30"
            : "border-slate-300 bg-slate-50/80 hover:border-[var(--brand)] hover:bg-sky-50/50 dark:border-neutral-700 dark:bg-neutral-900/50 dark:hover:bg-sky-950/20"
        } ${uploading ? "pointer-events-none opacity-60" : ""}`}
      >
        {uploading ? (
          <div className="flex items-center justify-center gap-2 text-xs text-slate-500 dark:text-neutral-400">
            <Spinner className="size-4" />
            Uploading to Cloudinary…
          </div>
        ) : (
          <>
            <svg
              className="mx-auto mb-2 size-6 text-slate-400 dark:text-neutral-500"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={1.5}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M12 16V4m0 0 4 4m-4-4-4 4M4 16v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-2"
              />
            </svg>
            <p className="text-xs font-medium text-slate-600 dark:text-neutral-300">
              Drag & drop a file here
            </p>
            <p className="mt-0.5 text-[10px] text-slate-400 dark:text-neutral-500">
              or click to browse · stored on Cloudinary
            </p>
          </>
        )}
      </div>

      {error && (
        <p className="mb-2 text-xs font-medium text-red-600 dark:text-red-400">{error}</p>
      )}

      {loading ? (
        <div className="flex justify-center py-3">
          <Spinner className="size-4" />
        </div>
      ) : attachments.length === 0 ? (
        <p className="text-xs text-slate-400 dark:text-neutral-500">No attachments yet.</p>
      ) : (
        <ul className="space-y-1.5">
          {attachments.map((attachment) => (
            <li
              key={attachment.id}
              className="flex items-center gap-2 rounded-lg border border-slate-200 bg-slate-50 px-3 py-1.5 dark:border-neutral-800 dark:bg-neutral-900"
            >
              <a
                href={attachment.url}
                target="_blank"
                rel="noopener noreferrer"
                className="min-w-0 flex-1 truncate text-xs font-medium text-[var(--brand)] hover:underline"
              >
                {attachment.file_name}
              </a>
              <span className="shrink-0 text-xs text-slate-400 dark:text-neutral-500">
                {formatSize(attachment.size)}
              </span>
              <button
                type="button"
                onClick={() => handleDelete(attachment)}
                aria-label={`Delete ${attachment.file_name}`}
                className="shrink-0 rounded p-1 text-slate-400 transition hover:text-red-600 dark:hover:text-red-400"
              >
                <svg className="size-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" d="M6 6l12 12M18 6L6 18" />
                </svg>
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
