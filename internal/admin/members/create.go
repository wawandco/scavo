package members

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"go.leapkit.dev/core/server"
)

// Create inserts a new member into the database.
func Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		server.Errorf(w, http.StatusBadRequest, "form parse error: %s", err.Error())
		return
	}

	name := r.FormValue("name")
	personalID := r.FormValue("personal_id")
	if name == "" || personalID == "" {
		server.Errorf(w, http.StatusBadRequest, "name and personal ID are required")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	res, err := db.Exec(
		"INSERT INTO members (name, personal_id, updated_at) VALUES (?, ?, ?)",
		name, personalID, time.Now(),
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "insert error: %s", err.Error())
		return
	}

	id, _ := res.LastInsertId()
	http.Redirect(w, r, fmt.Sprintf("/admin/members/%d", id), http.StatusSeeOther)
}
