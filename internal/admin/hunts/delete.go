package hunts

import (
	"database/sql"
	"net/http"
	"strconv"

	"go.leapkit.dev/core/server"
)

// Delete removes a hunt from the database.
func Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid hunt ID")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	_, err = db.Exec("DELETE FROM hunts WHERE id = ?", id)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "delete error: %s", err.Error())
		return
	}

	http.Redirect(w, r, "/admin/hunts", http.StatusSeeOther)
}
