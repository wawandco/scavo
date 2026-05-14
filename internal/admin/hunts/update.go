package hunts

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.leapkit.dev/core/server"
)

// Update updates an existing hunt.
func Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		server.Errorf(w, http.StatusBadRequest, "form parse error: %s", err.Error())
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid hunt ID")
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	if title == "" || description == "" {
		server.Errorf(w, http.StatusBadRequest, "title and description are required")
		return
	}

	points := 10
	if p := r.FormValue("points"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			points = parsed
		}
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	_, err = db.Exec(
		"UPDATE hunts SET title = ?, description = ?, points = ?, updated_at = ? WHERE id = ?",
		title, description, points, time.Now(), id,
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "update error: %s", err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/hunts/%d", id), http.StatusSeeOther)
}
