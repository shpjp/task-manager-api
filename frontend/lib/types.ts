export type TaskStatus = "todo" | "in_progress" | "done";
export type TaskPriority = "low" | "medium" | "high";

export interface Task {
  id: number;
  user_id: number;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  due_date: string | null;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface ListMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface TaskListResponse {
  data: Task[];
  meta: ListMeta;
}

export type SortBy = "created_at" | "due_date" | "priority";
export type SortOrder = "asc" | "desc";

export interface TaskQuery {
  status?: TaskStatus | "";
  search?: string;
  sort_by?: SortBy;
  order?: SortOrder;
  page?: number;
  limit?: number;
}

export interface TaskInput {
  title?: string;
  description?: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  due_date?: string | null;
}
