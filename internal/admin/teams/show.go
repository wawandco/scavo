package teams

import (
	"database/sql"
	"net/http"
	"strconv"

	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// Show renders a single team's details with its members.
func Show(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid team ID")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)

	var t admin.Team
	err = db.QueryRow(
		"SELECT id, name, created_at, updated_at FROM teams WHERE id = ?",
		id,
	).Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		server.Errorf(w, http.StatusNotFound, "team not found")
		return
	}
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}

	rows, err := db.Query(
		"SELECT id, name, personal_id, created_at FROM members WHERE team_id = ? ORDER BY name",
		id,
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}
	defer rows.Close()

	var members []admin.Member
	for rows.Next() {
		var m admin.Member
		if err := rows.Scan(&m.ID, &m.Name, &m.PersonalID, &m.CreatedAt); err != nil {
			server.Errorf(w, http.StatusInternalServerError, "scan error: %s", err.Error())
			return
		}
		members = append(members, m)
	}

	rw := render.FromCtx(r.Context())
	rw.Set("team", t)
	rw.Set("members", members)
	if err := rw.Render("admin/teams/show.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
