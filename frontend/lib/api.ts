import type {
  Task,
  TaskInput,
  TaskListResponse,
  TaskQuery,
  User,
} from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

const TOKEN_KEY = "task-manager-token";

export class ApiError extends Error {
  status: number;
  code: string;
  fields?: Record<string, string>;

  constructor(
    status: number,
    code: string,
    message: string,
    fields?: Record<string, string>
  ) {
    super(message);
    this.status = status;
    this.code = code;
    this.fields = fields;
  }
}

function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string | null) {
  if (typeof window === "undefined") return;
  if (token) {
    window.localStorage.setItem(TOKEN_KEY, token);
  } else {
    window.localStorage.removeItem(TOKEN_KEY);
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };
  const token = getToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  let res: Response;
  try {
    res = await fetch(`${API_URL}${path}`, {
      ...options,
      headers,
      credentials: "include",
    });
  } catch {
    throw new ApiError(0, "NETWORK_ERROR", "Could not reach the server. Is the API running?");
  }

  if (res.status === 204) {
    return undefined as T;
  }

  let body: unknown;
  try {
    body = await res.json();
  } catch {
    body = null;
  }

  if (!res.ok) {
    const err = (body as { error?: { code?: string; message?: string; fields?: Record<string, string> } })?.error;
    throw new ApiError(
      res.status,
      err?.code ?? "UNKNOWN",
      err?.message ?? "Something went wrong",
      err?.fields
    );
  }
  return body as T;
}

// --- Auth ---

interface AuthPayload {
  data: { user: User; token: string };
}

export async function signup(name: string, email: string, password: string): Promise<User> {
  const res = await request<AuthPayload>("/auth/signup", {
    method: "POST",
    body: JSON.stringify({ name, email, password }),
  });
  setToken(res.data.token);
  return res.data.user;
}

export async function login(email: string, password: string): Promise<User> {
  const res = await request<AuthPayload>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
  setToken(res.data.token);
  return res.data.user;
}

export async function logout(): Promise<void> {
  setToken(null);
  await request<void>("/auth/logout", { method: "POST" });
}

export async function me(): Promise<User> {
  const res = await request<{ data: User }>("/auth/me");
  return res.data;
}

// --- Tasks ---

export async function listTasks(query: TaskQuery): Promise<TaskListResponse> {
  const params = new URLSearchParams();
  if (query.status) params.set("status", query.status);
  if (query.search) params.set("search", query.search);
  if (query.sort_by) params.set("sort_by", query.sort_by);
  if (query.order) params.set("order", query.order);
  if (query.page) params.set("page", String(query.page));
  if (query.limit) params.set("limit", String(query.limit));
  const qs = params.toString();
  return request<TaskListResponse>(`/tasks${qs ? `?${qs}` : ""}`);
}

export async function getTask(id: number): Promise<Task> {
  const res = await request<{ data: Task }>(`/tasks/${id}`);
  return res.data;
}

export async function createTask(input: TaskInput): Promise<Task> {
  const res = await request<{ data: Task }>("/tasks", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return res.data;
}

export async function updateTask(id: number, input: TaskInput): Promise<Task> {
  const res = await request<{ data: Task }>(`/tasks/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return res.data;
}

export async function deleteTask(id: number): Promise<void> {
  await request<void>(`/tasks/${id}`, { method: "DELETE" });
}
