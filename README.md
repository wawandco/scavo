# Scavo

**Scavo** is a scavenger hunt web app built for [Wawandco](https://wawandco.com) team events (cumbres).

Teams solve mixed-media clues (text, photos, QR codes) in a timed event. It features a real-time leaderboard and a full admin interface for organizers.

This repo is published as a **showcase** of Wawandco’s engineering practices using Go, HTMX, Tailwind v4, and LeapKit.

> **License:** [O’Saasy](LICENSE) — permissive, with a restriction against competing hosted/SaaS offerings.

## How It Works

- Players log in with a personal ID
- Teams solve clues (text answers, photo uploads, or QR scans)
- Submissions are scored with points
- A live leaderboard updates in real time
- Organizers manage everything through a simple admin interface

The organizer can set time limits and close submissions at any moment.

## Tech Stack

- **Backend:** Go 1.25 with [LeapKit](https://go.leapkit.dev) framework
- **Frontend:** [HTMX 2.0](https://htmx.org/) + [Tailwind CSS v4](https://tailwindcss.com/)
- **Templating:** [Plush](https://github.com/gobuffalo/plush) (`<%= %>` syntax)
- **Database:** SQLite (WAL mode, busy timeout configured)
- **Deployment:** Docker (Alpine) + [Kamal](https://kamal-deploy.org)

## Prerequisites

- Go 1.25+
- Copy `.env.example` → `.env` and adjust values as needed

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

Scavo follows a simple, explicit structure (typical of small LeapKit apps):

- **Entry points** — `cmd/app` (HTTP server) and `cmd/migrate` (runs embedded migrations on startup)
- **Handlers** — Domain-organized under `internal/` (`home/`, `auth/`, `admin/hunts/`, `admin/teams/`, etc.)
- **Templates** — Co-located with handlers using Plush (`<%= %>` syntax)
- **Assets** — Embedded via `embed.FS` (CSS + JS served from `/internal/system/assets`)
- **Database** — SQLite with all migrations in `internal/migrations/` and executed automatically

This design prioritizes clarity and low cognitive overhead over heavy abstraction.

## Deployment

### Docker (Self-hosted)

```bash
# Build the image
docker build -t scavo .

# Run (mount a volume for persistent database + uploads)
docker run -p 3000:3000 \
  -e SESSION_SECRET=your-very-long-random-secret \
  -e DATABASE_URL=/data/scavo.db \
  -v scavo-data:/data \
  scavo
```

The container runs migrations automatically on startup.

### Kamal (Recommended for Production)

This project is deployed using [Kamal 2](https://kamal-deploy.org/).

Configuration lives in:
- `config/deploy.yml`
- `.kamal/secrets` (never commit real secrets)

Typical workflow:

```bash
kamal setup
kamal deploy
```

Make sure `SESSION_SECRET` and any other secrets are provided via Kamal secrets or your CI.

**Required production configuration:**
- Persistent volume for SQLite database and uploaded images
- `SESSION_SECRET` (strong random value)
- Proper `SERVER_IP` / hosts in `deploy.yml`

For local development, see the [Development](#development) section above.

## Project Context

This is a **use-once event app** designed for a single Wawandco cumbre. It prioritizes simplicity and reliability over long-term maintainability.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is released under the [O’Saasy License](LICENSE) — a permissive license with a restriction against offering competing hosted/SaaS versions.

It is published primarily as a showcase of Wawandco’s engineering practices.
