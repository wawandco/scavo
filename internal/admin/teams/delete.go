package teams

import (
	"database/sql"
	"net/http"
	"strconv"

	"go.leapkit.dev/core/server"
)

// Delete removes a team and unassigns its members.
func Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid team ID")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	_, err = db.Exec("UPDATE members SET team_id = NULL WHERE team_id = ?", id)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "unassign error: %s", err.Error())
		return
	}

	_, err = db.Exec("DELETE FROM teams WHERE id = ?", id)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "delete error: %s", err.Error())
		return
	}

	http.Redirect(w, r, "/admin/teams", http.StatusSeeOther)
}
