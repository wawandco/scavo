// Package internal integrates the app and loads general settings.
package internal

import (
	"cmp"
	"embed"
	"log/slog"
	"net/http"
	"os"

	"scavo/internal/admin/dashboard"
	"scavo/internal/admin/hunts"
	"scavo/internal/admin/members"
	"scavo/internal/admin/submissions"
	"scavo/internal/admin/teams"
	"scavo/internal/auth"
	"scavo/internal/home"
	"scavo/internal/scoreboard"
	"scavo/internal/system/assets"

	"go.leapkit.dev/core/db"
	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

var (
	//go:embed **/*.html **/*.html
	tmpls embed.FS

	// DBFn is the database connection builder function
	// that will be used by the application based on the driver and
	// connection string.
	DBFn = db.ConnectionFn(
		cmp.Or(os.Getenv("DATABASE_URL"), "database.db"),

		db.WithDriver("sqlite3"),
		db.Params(
			"_busy_timeout", "5000",
			"_journal_mode", "WAL",
			"_sync", "1",
			"_cache_size", "8192",
			"_txlock", "deferred",
		),
	)

	// Server configuration variables loaded from
	host          = cmp.Or(os.Getenv("HOST"), "0.0.0.0")
	port          = cmp.Or(os.Getenv("PORT"), "3000")
	sessionSecret = cmp.Or(os.Getenv("SESSION_SECRET"), "d720c059-9664-4980-8169-1158e167ae57")
	sessionName   = cmp.Or(os.Getenv("SESSION_NAME"), "leapkit_session")
)

// New creates the http handler using the Leapkit server package
// and returns it with the address it is listening on.
func New() (http.Handler, string) {
	// Creating a new server instance with the host and port
	// variables read from the environment or default values.
	r := server.New(
		server.WithHost(host),
		server.WithPort(port),
		server.WithSession(sessionSecret, sessionName),

		// Mounting the assets folder in the /assets URL path
		server.WithAssets(assets.Files, "/internal/system/assets"),
	)

	// Using the render middleware to load HTML templates
	// as well as setting a default layout for the application.
	r.Use(render.Middleware(
		render.TemplateFS(tmpls, "internal"),
		render.WithDefaultLayout("system/layout.html"),
	))

	// Database connection for admin and other handlers.
	conn, err := DBFn()
	if err != nil {
		slog.Error("connecting to database", "error", err)
		os.Exit(1)
	}

	// Inject the database connection into the request context.
	r.Use(server.InCtxMiddleware("db", conn))
	r.Use(auth.AuthContext)

	r.HandleFunc("GET /auth/users/login", auth.ShowLogin)
	r.HandleFunc("POST /auth/users/login", auth.Login)
	r.HandleFunc("POST /auth/users/logout", auth.Logout)

	r.Group("/", func(r server.Router) {
		// User routes (auth required)
		r.Use(auth.RequireAuth)

		r.HandleFunc("GET /{$}", home.Index)
		r.HandleFunc("GET /scoreboard", scoreboard.Index)
		r.HandleFunc("POST /submissions", home.CreateSubmission)

		r.HandleFunc("GET /uploads/{filename}", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "uploads/"+r.PathValue("filename"))
		})
	})

	r.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Admin routes (auth + admin required)
	r.Group("/admin/", func(r server.Router) {
		r.Use(auth.AdminGate)

		r.HandleFunc("GET /", dashboard.Index)

		r.HandleFunc("GET /hunts", hunts.List)
		r.HandleFunc("GET /hunts/new", hunts.New)
		r.HandleFunc("POST /hunts", hunts.Create)
		r.HandleFunc("GET /hunts/{id}", hunts.Show)
		r.HandleFunc("GET /hunts/{id}/edit", hunts.Edit)
		r.HandleFunc("POST /hunts/{id}", hunts.Update)
		r.HandleFunc("POST /hunts/{id}/delete", hunts.Delete)

		r.HandleFunc("GET /submissions", submissions.List)
		r.HandleFunc("POST /hunts/{id}/submissions/{submission_id}/approve", submissions.Approve)
		r.HandleFunc("POST /hunts/{id}/submissions/{submission_id}/reject", submissions.Reject)
		r.HandleFunc("POST /hunts/{id}/submissions/{submission_id}/delete", submissions.Delete)

		r.HandleFunc("GET /teams", teams.List)
		r.HandleFunc("GET /teams/new", teams.New)
		r.HandleFunc("POST /teams", teams.Create)
		r.HandleFunc("GET /teams/{id}", teams.Show)
		r.HandleFunc("GET /teams/{id}/edit", teams.Edit)
		r.HandleFunc("POST /teams/{id}", teams.Update)
		r.HandleFunc("POST /teams/{id}/delete", teams.Delete)

		r.HandleFunc("GET /members", members.List)
		r.HandleFunc("GET /members/new", members.New)
		r.HandleFunc("POST /members", members.Create)
		r.HandleFunc("GET /members/{id}", members.Show)
		r.HandleFunc("GET /members/{id}/edit", members.Edit)
		r.HandleFunc("POST /members/{id}", members.Update)
		r.HandleFunc("POST /members/{id}/delete", members.Delete)
	})

	return r.Handler(), r.Addr()
}
