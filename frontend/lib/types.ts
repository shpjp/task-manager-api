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
  /** Present only in admin "all users" listings. */
  owner?: User;
}

export interface User {
  id: number;
  name: string;
  email: string;
  role: "user" | "admin";
  created_at: string;
  updated_at: string;
}

export interface Attachment {
  id: number;
  task_id: number;
  user_id: number;
  file_name: string;
  content_type: string;
  size: number;
  url: string;
  created_at: string;
}

export interface TaskActivity {
  id: number;
  task_id: number;
  user_id: number;
  action: string;
  detail: string;
  created_at: string;
  user?: User;
}

export interface TaskEvent {
  type: "task.created" | "task.updated" | "task.deleted";
  task_id: number;
  task?: Task;
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
  /** Admin only: list every user's tasks. */
  scope?: "all" | "";
}

export interface TaskInput {
  title?: string;
  description?: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  due_date?: string | null;
}
