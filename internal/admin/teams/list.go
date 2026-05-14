package teams

import (
	"database/sql"
	"net/http"

	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// List renders the list of all teams.
func List(w http.ResponseWriter, r *http.Request) {
	db, _ := r.Context().Value("db").(*sql.DB)
	rows, err := db.Query("SELECT id, name, created_at, updated_at FROM teams ORDER BY created_at DESC")
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}
	defer rows.Close()

	var teams []admin.Team
	for rows.Next() {
		var t admin.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			server.Errorf(w, http.StatusInternalServerError, "scan error: %s", err.Error())
			return
		}
		teams = append(teams, t)
	}

	rw := render.FromCtx(r.Context())
	rw.Set("teams", teams)
	if err := rw.Render("admin/teams/list.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
