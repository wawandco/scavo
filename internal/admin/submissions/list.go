// Package submissions provides the admin HTTP handlers for reviewing,
// approving, and rejecting participant photo and text submissions.
package submissions

import (
	"database/sql"
	"net/http"

	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// List renders all submissions across all hunts.
func List(w http.ResponseWriter, r *http.Request) {
	db, _ := r.Context().Value("db").(*sql.DB)
	rows, err := db.Query(`
		SELECT s.id, s.hunt_id, s.member_id, s.team_id, s.text, s.image_path, s.status, s.created_at, s.updated_at, h.title, m.name, t.name
		FROM submissions s
		JOIN hunts h ON s.hunt_id = h.id
		JOIN members m ON s.member_id = m.id
		JOIN teams t ON s.team_id = t.id
		ORDER BY s.created_at DESC`)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}
	defer rows.Close() //nolint:errcheck // best effort cleanup

	var submissions []admin.Submission
	for rows.Next() {
		var s admin.Submission
		if err := rows.Scan(&s.ID, &s.HuntID, &s.MemberID, &s.TeamID, &s.Text, &s.ImagePath, &s.Status, &s.CreatedAt, &s.UpdatedAt, &s.HuntTitle, &s.MemberName, &s.TeamName); err != nil {
			server.Errorf(w, http.StatusInternalServerError, "scan error: %s", err.Error())
			return
		}
		submissions = append(submissions, s)
	}

	rw := render.FromCtx(r.Context())
	rw.Set("submissions", submissions)
	if err := rw.Render("admin/submissions/list.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
