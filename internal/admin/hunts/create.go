package hunts

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.leapkit.dev/core/server"
)

// Create inserts a new hunt into the database.
func Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		server.Errorf(w, http.StatusBadRequest, "form parse error: %s", err.Error())
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
	res, err := db.Exec(
		"INSERT INTO hunts (title, description, points, updated_at) VALUES (?, ?, ?, ?)",
		title, description, points, time.Now(),
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "insert error: %s", err.Error())
		return
	}

	id, _ := res.LastInsertId()
	http.Redirect(w, r, fmt.Sprintf("/admin/hunts/%d", id), http.StatusSeeOther)
}
