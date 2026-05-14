# Scavo

Scavo is a scavenger hunt web application built for [Wawandco](https://wawandco.com)'s team events ("cumbres"). Teams compete by solving mixed-media clues — text answers, photo submissions, and QR code scans — within a time limit set by the event organizer. A real-time leaderboard tracks progress, and organizers manage hunts through a web admin interface.

## How It Works

- **Login:** Users authenticate with a pre-defined personal ID assigned before the event.
- **Gameplay:** Teams progress through a hunt by solving clues. Clue types include:
  - Text clues with typed answers
  - Photo/image submissions (proof-based tasks)
  - QR code scanning at physical locations
- **Scoring:** Points-based system. The organizer sets a time limit and can close submissions at any time.
- **Leaderboard:** Real-time ranking updates as teams submit answers.
- **Admin:** Organizers create hunts, add clues, define time limits, monitor team progress, and control event flow via a web admin panel.

## Tech Stack

- **Backend:** Go 1.25 with [LeapKit](https://go.leapkit.dev) framework
- **Frontend:** [HTMX 2.0](https://htmx.org/) + [Tailwind CSS v4](https://tailwindcss.com/)
- **Templating:** [Plush](https://github.com/gobuffalo/plush) (`<%= %>` syntax)
- **Database:** SQLite (WAL mode, busy timeout configured)
- **Deployment:** Docker (Alpine Linux multi-stage build)

## Prerequisites

- Go 1.25+
- A `.env` file at repo root (see existing `.env` for defaults)

## Development

### Quick Start

```sh
# Download dependencies
go mod download

# Start dev server (auto-restarts on Go changes + watches CSS)
go tool dev
```

The app will be available at [http://localhost:3000](http://localhost:3000).

### Key Commands

| Task | Command |
|------|---------|
| Dev server (auto-restart + CSS watch) | `go tool dev` |
| Build CSS once | `go tool tailo -i internal/system/assets/tailwind.css -o internal/system/assets/application.css` |
| Run database migrations | `go run ./cmd/migrate` |
| Build app binary | `go build -o bin/app ./cmd/app` |
| Build migrate binary | `go build -o bin/migrate ./cmd/migrate` |

The `go tool dev` command reads the root `Procfile` and runs:
- The Go app (auto-restarted on `.go` file changes)
- Tailwind CSS watcher (rebuilds `application.css` on changes)

### Environment Variables

Key variables (loaded automatically from `.env`):

- `DATABASE_URL` — SQLite database path (default: `database.db`)
- `PORT` — HTTP server port (default: `3000`)
- `SESSION_SECRET` — Secret for session encryption
- `SESSION_NAME` — Session cookie name

## Architecture

- **Entrypoints:**
  - `cmd/app/main.go` — HTTP server
  - `cmd/migrate/migrate.go` — Database migration runner
- **Handlers & Routes:** Located in `internal/` (e.g., `internal/app.go`, `internal/index.go`)
- **Templates:** Co-located with handlers under `internal/`; the render middleware mounts `tmpls embed.FS` under the `internal` prefix
- **Layout:** `internal/system/layout.html` is the default layout; pages yield into it with `<%= yield %>`
- **Assets:** `internal/system/assets/` are embedded via `embed.FS` and served at `/internal/system/assets`
- **Migrations:** `internal/migrations/*.sql` are embedded and executed with `db.RunMigrations`

## Building & Deployment

### Docker

The included `Dockerfile` builds a minimal Alpine image:

1. Downloads the Tailwind binary (`go tool tailo download`)
2. Compiles CSS
3. Builds `migrate` and `app` binaries with `-tags osusergo,netgo`
4. Container startup runs `migrate && app`

### Manual Build

```sh
# 1. Build CSS
go tool tailo -i internal/system/assets/tailwind.css -o internal/system/assets/application.css

# 2. Run migrations
go run ./cmd/migrate

# 3. Build and run the app
go build -o bin/app ./cmd/app
./bin/app
```

## Project Context

This is a **use-once event app** designed for a single Wawandco cumbre. It prioritizes simplicity and reliability over long-term maintainability.

## License

Internal Wawandco project. Not open source.
