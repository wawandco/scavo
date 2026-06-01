// Package members provides the admin HTTP handlers for managing
// event participants (team members).
package members

import (
	"database/sql"
	"net/http"

	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// List renders the list of all members.
func List(w http.ResponseWriter, r *http.Request) {
	db, _ := r.Context().Value("db").(*sql.DB)
	rows, err := db.Query(`
		SELECT m.id, m.name, m.personal_id, m.team_id, m.is_admin, m.created_at, m.updated_at, COALESCE(t.name, '')
		FROM members m
		LEFT JOIN teams t ON m.team_id = t.id
		ORDER BY m.created_at DESC`)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}
	defer rows.Close() //nolint:errcheck // best effort cleanup

	var members []admin.Member
	for rows.Next() {
		var m admin.Member
		if err := rows.Scan(&m.ID, &m.Name, &m.PersonalID, &m.TeamID, &m.IsAdmin, &m.CreatedAt, &m.UpdatedAt, &m.TeamName); err != nil {
			server.Errorf(w, http.StatusInternalServerError, "scan error: %s", err.Error())
			return
		}
		members = append(members, m)
	}

	rw := render.FromCtx(r.Context())
	rw.Set("members", members)
	if err := rw.Render("admin/members/list.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
