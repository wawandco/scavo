package members

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.leapkit.dev/core/server"
)

// Update updates an existing member.
func Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		server.Errorf(w, http.StatusBadRequest, "form parse error: %s", err.Error())
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid member ID")
		return
	}

	name := r.FormValue("name")
	personalID := r.FormValue("personal_id")
	if name == "" || personalID == "" {
		server.Errorf(w, http.StatusBadRequest, "name and personal ID are required")
		return
	}

	var teamID sql.NullInt64
	if tid := r.FormValue("team_id"); tid != "" {
		tid64, err := strconv.ParseInt(tid, 10, 64)
		if err == nil {
			teamID = sql.NullInt64{Int64: tid64, Valid: true}
		}
	}

	isAdmin := 0
	if r.FormValue("is_admin") == "1" {
		isAdmin = 1
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	_, err = db.Exec(
		"UPDATE members SET name = ?, personal_id = ?, team_id = ?, is_admin = ?, updated_at = ? WHERE id = ?",
		name, personalID, teamID, isAdmin, time.Now(), id,
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "update error: %s", err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/members/%d", id), http.StatusSeeOther)
}
