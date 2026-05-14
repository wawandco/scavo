package teams

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.leapkit.dev/core/server"
)

// Create inserts a new team into the database.
func Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		server.Errorf(w, http.StatusBadRequest, "form parse error: %s", err.Error())
		return
	}

	name := r.FormValue("name")
	if name == "" {
		server.Errorf(w, http.StatusBadRequest, "name is required")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	res, err := db.Exec(
		"INSERT INTO teams (name, updated_at) VALUES (?, ?)",
		name, time.Now(),
	)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "insert error: %s", err.Error())
		return
	}

	id, _ := res.LastInsertId()

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
