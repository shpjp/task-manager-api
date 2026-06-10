# Taskflow — Task Management Application

A full-stack task manager with a **Go (Gin + GORM)** REST API, a **Next.js** frontend, and **PostgreSQL** persistence.

## Features

**Core**

- Task CRUD: title, description, status (`todo` / `in_progress` / `done`), priority (`low` / `medium` / `high`), due date
- Filtering by status, title search, sorting (due date, priority, created date) and pagination — all combinable
- JWT authentication (signup / login) with bcrypt-hashed passwords
- Users can only see and modify their own tasks
- Auth persists across page refreshes (httpOnly cookie + bearer token)
- Input validation on every write endpoint with per-field error messages
- Consistent JSON error envelope and proper HTTP status codes
- Responsive UI with loading, empty, and error states

**Bonus**

- **Admin role** — emails listed in `ADMIN_EMAILS` get read-only access to every user's tasks ("All users" toggle in the UI)
- **Real-time updates** — task changes stream live over Server-Sent Events
- **Optimistic UI** — complete/delete update instantly and roll back on failure
- **Attachments** — upload images/documents to tasks (5 MB limit, type allowlist)
- **Activity log** — per-task history of every change
- **Dockerized setup** — one command brings up the whole stack
- **CI pipeline** — GitHub Actions builds, vets, lints and tests on every push
- **Dark mode** — theme toggle persisted across sessions

## Quick start (Docker, one command)

```bash
docker compose up --build
```

Then open **http://localhost:3000** (API at http://localhost:8080).

Sign up with `admin@example.com` to try the admin role (configurable via `ADMIN_EMAILS`).

## Local development setup

Prerequisites: Go 1.26+, Node 22+, Docker (for Postgres).

```bash
# 1. Start PostgreSQL
docker compose up -d postgres

# 2. Configure the backend
cp .env.example .env        # fill in JWT_SECRET (any long random string)

# 3. Run the API
go run cmd/server/main.go   # listens on :8080

# 4. Run the frontend (separate terminal)
cd frontend
cp .env.example .env.local
npm install
npm run dev                 # listens on :3000
```

> **Note:** Postgres is mapped to host port **5433** to avoid clashing with a locally installed Postgres on 5432.

## Environment variables

Backend (`.env`, see `.env.example`):

| Variable | Description | Default |
| --- | --- | --- |
| `DB_HOST` | Postgres host | `localhost` |
| `DB_PORT` | Postgres port | `5433` |
| `DB_USER` | Postgres user | `postgres` |
| `DB_PASSWORD` | Postgres password | — |
| `DB_NAME` | Database name | `taskmanager` |
| `APP_PORT` | API listen port | `8080` |
| `JWT_SECRET` | JWT signing secret (**required**) | — |
| `TOKEN_TTL_HOURS` | JWT lifetime in hours | `24` |
| `FRONTEND_ORIGIN` | Allowed CORS origin | `http://localhost:3000` |
| `COOKIE_SECURE` | Set `true` behind HTTPS | `false` |
| `ADMIN_EMAILS` | Comma-separated emails promoted to admin | — |
| `UPLOAD_DIR` | Attachment storage directory | `./uploads` |
| `MAX_UPLOAD_MB` | Max attachment size (MB) | `5` |

Frontend (`frontend/.env.local`, see `frontend/.env.example`):

| Variable | Description | Default |
| --- | --- | --- |
| `NEXT_PUBLIC_API_URL` | Base URL of the API | `http://localhost:8080` |

## API reference

All `/tasks` routes require authentication (`Authorization: Bearer <token>` or the auth cookie).

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/auth/signup` | Create an account (`name`, `email`, `password`) |
| `POST` | `/auth/login` | Log in, returns user + JWT (also sets httpOnly cookie) |
| `POST` | `/auth/logout` | Clear the auth cookie |
| `GET` | `/auth/me` | Current user |
| `POST` | `/tasks` | Create a task |
| `GET` | `/tasks` | List tasks (`status`, `search`, `sort_by`, `order`, `page`, `limit`, admin: `scope=all`) |
| `GET` | `/tasks/:id` | Fetch one task |
| `PATCH` | `/tasks/:id` | Partially update a task (`"due_date": null` clears it) |
| `DELETE` | `/tasks/:id` | Delete a task |
| `GET` | `/tasks/:id/activity` | Change history for a task |
| `POST` | `/tasks/:id/attachments` | Upload a file (multipart field `file`) |
| `GET` | `/tasks/:id/attachments` | List attachments |
| `GET` | `/tasks/:id/attachments/:attachmentID/download` | Download an attachment |
| `DELETE` | `/tasks/:id/attachments/:attachmentID` | Delete an attachment |
| `GET` | `/events` | Live task events (SSE, cookie auth) |
| `GET` | `/health` | Health check |

Errors always use:

```json
{ "error": { "code": "VALIDATION_ERROR", "message": "...", "fields": { "title": "This field is required" } } }
```

Lists return `{ "data": [...], "meta": { "page", "limit", "total", "total_pages" } }`.

## Tests

```bash
go test ./...
```

11 backend tests cover JWT issue/verify/expiry, password hashing, signup/login flows, route protection, write validation, per-user task isolation, admin access control, the activity log, and combined filter/search/sort/pagination. Handler tests run the full HTTP stack against an in-memory SQLite database.

CI (GitHub Actions) runs build + vet + tests for the backend and lint + build for the frontend on every push and pull request.

## Architecture

```
cmd/server          — entry point, wiring
internal/
  config            — env config + DB connection
  models            — GORM models (User, Task, TaskActivity, Attachment)
  repository        — data access layer
  services          — business logic (auth, tasks, attachments)
  handlers          — HTTP handlers, validation, error mapping
  middleware        — JWT auth middleware
  realtime          — SSE hub
  routes            — route registration
frontend/
  app               — Next.js App Router pages
  components        — UI components
  lib               — API client, auth/theme contexts, types
```

## Assumptions & trade-offs

- **JWT carries the role claim** — avoids a DB lookup per request; a role change takes effect on the next login (tokens live 24h by default).
- **Admin access is read-only** — the spec says admins "view" all tasks, so they cannot modify other users' tasks.
- **Attachments are stored on local disk** (Docker volume in compose). For multi-instance deployments, swap the storage layer for S3 or similar.
- **SSE over WebSockets** — one-directional updates are all the UI needs; SSE is simpler, proxies well, and auto-reconnects.
- **GORM AutoMigrate** instead of versioned migrations — acceptable at this scale; I'd switch to golang-migrate/atlas before a real production launch.
- **Handler tests use in-memory SQLite** for speed and zero test infrastructure; queries stick to portable SQL (the one Postgres-specific bit, `LIKE` escaping, is explicitly `ESCAPE`d for both).
- **httpOnly cookie + bearer token** — the cookie keeps refreshes logged in without exposing the token to XSS; the header path keeps the API curl-friendly.
