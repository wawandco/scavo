package members

import (
	"database/sql"
	"net/http"
	"strconv"

	"go.leapkit.dev/core/server"
)

// Delete removes a member from the database.
func Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid member ID")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	_, err = db.Exec("DELETE FROM members WHERE id = ?", id)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "delete error: %s", err.Error())
		return
	}

	http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
}
