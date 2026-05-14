package home

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.leapkit.dev/core/server"
	"go.leapkit.dev/core/server/session"
)

// CreateSubmission handles a team member uploading a hunt submission.
func CreateSubmission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		server.Errorf(w, http.StatusBadRequest, "form parse error: %s", err.Error())
		return
	}

	huntID := r.FormValue("hunt_id")
	if huntID == "" {
		server.Errorf(w, http.StatusBadRequest, "hunt ID is required")
		return
	}

	s := session.FromCtx(r.Context())
	userID, ok := s.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/auth/users/login", http.StatusSeeOther)
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	var teamID sql.NullInt64
	_ = db.QueryRow("SELECT team_id FROM members WHERE id = ?", userID).Scan(&teamID)
	if !teamID.Valid {
		server.Errorf(w, http.StatusBadRequest, "you must be assigned to a team to submit")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "image is required")
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("%d_%d_%s", time.Now().Unix(), userID, header.Filename)
	out, err := os.Create("uploads/" + filename)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "file save error: %s", err.Error())
		return
	}
	defer out.Close()
	_, _ = io.Copy(out, file)

	_, _ = db.Exec(
		"DELETE FROM submissions WHERE hunt_id = ? AND team_id = ?",
		huntID, teamID.Int64,
	)

	_, _ = db.Exec(
		"INSERT INTO submissions (hunt_id, member_id, team_id, text, image_path, status, updated_at) VALUES (?, ?, ?, ?, ?, 'pending', ?)",
		huntID, userID, teamID.Int64, r.FormValue("text"), filename, time.Now(),
	)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
