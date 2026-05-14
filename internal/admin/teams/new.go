package teams

import (
	"database/sql"
	"net/http"

	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// New renders the form to create a new team.
func New(w http.ResponseWriter, r *http.Request) {
	db, _ := r.Context().Value("db").(*sql.DB)
	rows, err := db.Query(
		"SELECT id, name, personal_id, team_id FROM members ORDER BY name",
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}
	defer rows.Close()

	var allMembers []admin.Member
	for rows.Next() {
		var m admin.Member
		if err := rows.Scan(&m.ID, &m.Name, &m.PersonalID, &m.TeamID); err != nil {
			server.Errorf(w, http.StatusInternalServerError, "scan error: %s", err.Error())
			return
		}
		allMembers = append(allMembers, m)
	}

	rw := render.FromCtx(r.Context())
	rw.Set("members", allMembers)
	if err := rw.Render("admin/teams/new.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
