package hunts

import (
	"database/sql"
	"net/http"

	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// List renders the list of all hunts.
func List(w http.ResponseWriter, r *http.Request) {
	db, _ := r.Context().Value("db").(*sql.DB)
	rows, err := db.Query("SELECT id, title, description, points, created_at, updated_at FROM hunts ORDER BY created_at DESC")
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}
	defer rows.Close()

	var hunts []admin.Hunt
	for rows.Next() {
		var h admin.Hunt
		if err := rows.Scan(&h.ID, &h.Title, &h.Description, &h.Points, &h.CreatedAt, &h.UpdatedAt); err != nil {
			server.Errorf(w, http.StatusInternalServerError, "scan error: %s", err.Error())
			return
		}
		hunts = append(hunts, h)
	}

	rw := render.FromCtx(r.Context())
	rw.Set("hunts", hunts)
	if err := rw.Render("admin/hunts/list.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
