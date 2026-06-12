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

export function TaskAttachments({ taskId }: { taskId: number }) {
  const [attachments, setAttachments] = useState<Attachment[]>([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
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

  async function handleUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    e.target.value = "";
    if (!file) return;

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
          {uploading ? <Spinner className="size-3" /> : "+"} Upload file
        </button>
        <input
          ref={fileInput}
          type="file"
          hidden
          accept=".png,.jpg,.jpeg,.gif,.webp,.pdf,.txt,.md,.csv,.doc,.docx,.xls,.xlsx"
          onChange={handleUpload}
        />
      </div>

      {error && <p className="mb-2 text-xs font-medium text-red-600 dark:text-red-400">{error}</p>}

      {loading ? (
        <div className="flex justify-center py-3">
          <Spinner className="size-4" />
        </div>
      ) : attachments.length === 0 ? (
        <p className="text-xs text-slate-400 dark:text-slate-500">No attachments yet.</p>
      ) : (
        <ul className="space-y-1.5">
          {attachments.map((attachment) => (
            <li
              key={attachment.id}
              className="flex items-center gap-2 rounded-lg border border-slate-200 bg-slate-50 px-3 py-1.5 dark:border-neutral-800 dark:bg-neutral-900"
            >
              <a
                href={api.attachmentDownloadUrl(taskId, attachment.id)}
                className="min-w-0 flex-1 truncate text-xs font-medium text-indigo-600 hover:underline dark:text-indigo-400"
              >
                {attachment.file_name}
              </a>
              <span className="shrink-0 text-xs text-slate-400 dark:text-slate-500">
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
