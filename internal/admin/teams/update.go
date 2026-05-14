package teams

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.leapkit.dev/core/server"
)

// Update updates an existing team.
func Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		server.Errorf(w, http.StatusBadRequest, "form parse error: %s", err.Error())
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid team ID")
		return
	}

	name := r.FormValue("name")
	if name == "" {
		server.Errorf(w, http.StatusBadRequest, "name is required")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	_, err = db.Exec(
		"UPDATE teams SET name = ?, updated_at = ? WHERE id = ?",
		name, time.Now(), id,
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "update error: %s", err.Error())
		return
	}

	_, _ = db.Exec(
		"UPDATE members SET team_id = NULL, updated_at = ? WHERE team_id = ?",
		time.Now(), id,
	)

	for _, mid := range r.Form["member_ids"] {
		memberID, err := strconv.Atoi(mid)
		if err != nil {
			continue
		}
		_, _ = db.Exec(
			"UPDATE members SET team_id = ?, updated_at = ? WHERE id = ?",
			id, time.Now(), memberID,
		)
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/teams/%d", id), http.StatusSeeOther)
}
