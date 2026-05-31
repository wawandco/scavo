// Package home contains the public participant-facing handlers:
// the main game view and submission creation.
package home

import (
	"database/sql"
	"net/http"

	"scavo/internal/db"
	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
	"go.leapkit.dev/core/server/session"
)

// HuntStatus holds a hunt along with the team's submission status.
type HuntStatus struct {
	ID              int
	Title           string
	Description     string
	Points          int
	SubmissionID    int
	SubmissionText  string
	SubmissionImage string
	Status          string
}

// Index renders the public home page with the user's team and all hunts.
func Index(w http.ResponseWriter, r *http.Request) {
	s := session.FromCtx(r.Context())
	userID, ok := s.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/auth/users/login", http.StatusSeeOther)
		return
	}

	db := db.FromCtx(r.Context())

	var teamID sql.NullInt64
	var member admin.Member
	_ = db.QueryRow(
		"SELECT id, name, team_id FROM members WHERE id = ?",
		userID,
	).Scan(&member.ID, &member.Name, &teamID)

	var teamName string
	if teamID.Valid {
		_ = db.QueryRow("SELECT name FROM teams WHERE id = ?", teamID.Int64).Scan(&teamName)
	}

	var totalPoints int
	_ = db.QueryRow(`
		SELECT COALESCE(SUM(h.points), 0)
		FROM submissions s
		JOIN hunts h ON s.hunt_id = h.id
		WHERE s.team_id = ? AND s.status = 'approved'`,
		teamID,
	).Scan(&totalPoints)

	var possiblePoints int
	_ = db.QueryRow(`SELECT COALESCE(SUM(points), 0) FROM hunts`).Scan(&possiblePoints)

	rows, err := db.Query(`
		SELECT h.id, h.title, h.description, h.points,
		       COALESCE(s.id, 0), COALESCE(s.text, ''), COALESCE(s.image_path, ''), COALESCE(s.status, '')
		FROM hunts h
		LEFT JOIN submissions s ON h.id = s.hunt_id AND s.team_id = ?
		ORDER BY h.created_at DESC`,
		teamID,
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}
	defer rows.Close()

	var huntStatuses []HuntStatus
	for rows.Next() {
		var hs HuntStatus
		if err := rows.Scan(&hs.ID, &hs.Title, &hs.Description, &hs.Points, &hs.SubmissionID, &hs.SubmissionText, &hs.SubmissionImage, &hs.Status); err != nil {
			server.Errorf(w, http.StatusInternalServerError, "scan error: %s", err.Error())
			return
		}
		huntStatuses = append(huntStatuses, hs)
	}

	rw := render.FromCtx(r.Context())
	rw.Set("member", member)
	rw.Set("teamName", teamName)
	rw.Set("totalPoints", totalPoints)
	rw.Set("possiblePoints", possiblePoints)
	rw.Set("hunts", huntStatuses)
	if err := rw.Render("home/index.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
