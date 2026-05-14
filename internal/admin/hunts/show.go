package hunts

import (
	"database/sql"
	"net/http"
	"strconv"

	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// Show renders a single hunt's details.
func Show(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid hunt ID")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	var h admin.Hunt
	err = db.QueryRow(
		"SELECT id, title, description, points, created_at, updated_at FROM hunts WHERE id = ?",
		id,
	).Scan(&h.ID, &h.Title, &h.Description, &h.Points, &h.CreatedAt, &h.UpdatedAt)
	if err == sql.ErrNoRows {
		server.Errorf(w, http.StatusNotFound, "hunt not found")
		return
	}
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}

	rows, err := db.Query(`
		SELECT s.id, s.hunt_id, s.member_id, s.team_id, s.text, s.image_path, s.status, s.created_at, s.updated_at, m.name, t.name
		FROM submissions s
		JOIN members m ON s.member_id = m.id
		JOIN teams t ON s.team_id = t.id
		WHERE s.hunt_id = ?
		ORDER BY s.created_at DESC`,
		id,
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}
	defer rows.Close()

	var submissions []admin.Submission
	for rows.Next() {
		var s admin.Submission
		if err := rows.Scan(&s.ID, &s.HuntID, &s.MemberID, &s.TeamID, &s.Text, &s.ImagePath, &s.Status, &s.CreatedAt, &s.UpdatedAt, &s.MemberName, &s.TeamName); err != nil {
			server.Errorf(w, http.StatusInternalServerError, "scan error: %s", err.Error())
			return
		}
		submissions = append(submissions, s)
	}

	rw := render.FromCtx(r.Context())
	rw.Set("hunt", h)
	rw.Set("submissions", submissions)
	if err := rw.Render("admin/hunts/show.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
