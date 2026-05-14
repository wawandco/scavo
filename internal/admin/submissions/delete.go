package submissions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"go.leapkit.dev/core/server"
)

// Delete removes a submission from the database.
func Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid hunt ID")
		return
	}

	submissionIDStr := r.PathValue("submission_id")
	submissionID, err := strconv.Atoi(submissionIDStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid submission ID")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	_, err = db.Exec("DELETE FROM submissions WHERE id = ?", submissionID)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "delete error: %s", err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/hunts/%s", idStr), http.StatusSeeOther)
}
