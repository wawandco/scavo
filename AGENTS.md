# AGENTS.md — Scavo (LeapKit / Go)

## Quick Commands

| Task | Command |
|------|---------|
| Dev server (auto-restart + CSS watch) | `go tool dev` |
| Build CSS once | `go tool tailo -i internal/system/assets/tailwind.css -o internal/system/assets/application.css` |
| Run migrations | `go run ./cmd/migrate` |
| Build app binary | `go build -o bin/app ./cmd/app` |
| Build migrate binary | `go build -o bin/migrate ./cmd/migrate` |

`go tool dev` reads the root `Procfile`. It auto-restarts the Go process on `.go` changes and runs Tailwind in watch mode in parallel.

## Environment

- `.env` at repo root is **auto-loaded** by `go.leapkit.dev/core/tools/envload` (blank import in `cmd/app/main.go` and `cmd/migrate/migrate.go`).
- Key vars: `DATABASE_URL` (defaults to `database.db`), `PORT` (3000), `SESSION_SECRET`, `SESSION_NAME`.
- SQLite is configured with WAL mode, busy timeout, and other pragmas in `internal/app.go`.

## Architecture

- **Framework:** LeapKit (`go.leapkit.dev/core`). Server, render middleware, session middleware, DB connection helpers.
- **Frontend:** HTMX 2.0 (CDN), Tailwind CSS v4 (`internal/system/assets/tailwind.css` uses `@import "tailwindcss"` and plugins).
- **Templating:** Plush (`github.com/gobuffalo/plush/v5`), not standard `html/template`. Uses `<%= %>` syntax and helpers like `assetPath(...)` and `yield`.
- **Entrypoints:**
  - `cmd/app/main.go` — HTTP server.
  - `cmd/migrate/migrate.go` — runs embedded SQL migrations.
- **Assets:** `internal/system/assets/` are embedded via `embed.FS` and served at `/internal/system/assets`.
- **Migrations:** `internal/migrations/*.sql` are embedded via `embed.FS` and executed with `db.RunMigrations`.

## Build & Deploy

- Dockerfile is multi-stage Alpine. It downloads the Tailwind binary (`go tool tailo download`), compiles CSS, then builds both `migrate` and `app` with `-tags osusergo,netgo`.
- Container startup runs `migrate && app`.
- No CI workflows or tests currently exist in the repo.

## Conventions

- Routes and handlers live in `internal/` (e.g., `internal/app.go`, `internal/index.go`).
- HTML templates are co-located with handlers in `internal/` and its subdirectories; the render middleware mounts `tmpls embed.FS` under the `internal` prefix.
- The default layout is `internal/system/layout.html`; page templates yield into it with `<%= yield %>`.
- No `_test.go` files exist yet — add them in `internal/` or alongside the code they test.

### Handler Structure

- Admin handlers are organized into per-resource subdirectories under `internal/admin/`: `teams/`, `members/`, `hunts/`, `dashboard/`.
- Each handler lives in its own file named by action: `list.go`, `show.go`, `new.go`, `create.go`, `edit.go`, `update.go`, `delete.go`.
- Handler functions are named by action only (e.g., `teams.List`, `members.Delete`, `hunts.Show`, `dashboard.Index`).
- Templates mirror this structure: each action has a matching `.html` file in the same subdirectory (e.g., `teams/list.html` for `teams/list.go`).

### Database Access

- The DB connection is injected into the request context via `server.InCtxMiddleware("db", conn)` in `internal/app.go`.
- Handlers pull the DB directly from the request context: `db, _ := r.Context().Value("db").(*sql.DB)`.
- No package-level DB variable or `admin.DB()` helper exists. `admin.FetchTeams(db)` accepts `*sql.DB` explicitly.

### Agent Server Port

- When starting the server as an agent, always specify a `PORT` different from 3000 to avoid conflicts with the dev server: `PORT=3001 go run ./cmd/app`.
