package members

import (
	"database/sql"
	"net/http"
	"strconv"

	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// Edit renders the form to edit a member.
func Edit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid member ID")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	var m admin.Member
	err = db.QueryRow(`
		SELECT m.id, m.name, m.personal_id, m.team_id, m.is_admin, m.created_at, m.updated_at, COALESCE(t.name, '')
		FROM members m
		LEFT JOIN teams t ON m.team_id = t.id
		WHERE m.id = ?`,
		id,
	).Scan(&m.ID, &m.Name, &m.PersonalID, &m.TeamID, &m.IsAdmin, &m.CreatedAt, &m.UpdatedAt, &m.TeamName)
	if err == sql.ErrNoRows {
		server.Errorf(w, http.StatusNotFound, "member not found")
		return
	}
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}

	teams, err := admin.FetchTeams(db)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}

	rw := render.FromCtx(r.Context())
	rw.Set("member", m)
	rw.Set("teams", teams)
	if err := rw.Render("admin/members/edit.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
