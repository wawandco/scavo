package hunts

import (
	"database/sql"
	"net/http"
	"strconv"

	"scavo/internal/admin"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// Edit renders the form to edit a hunt.
func Edit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "invalid hunt ID")
		return
	}

	db, _ := r.Context().Value("db").(*sql.DB)
	var h admin.Hunt
	err = db.QueryRow(
		"SELECT id, title, description, points, created_at, updated_at FROM hunts WHERE id = ?",
		id,
	).Scan(&h.ID, &h.Title, &h.Description, &h.Points, &h.CreatedAt, &h.UpdatedAt)
	if err == sql.ErrNoRows {
		server.Errorf(w, http.StatusNotFound, "hunt not found")
		return
	}
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}

	rw := render.FromCtx(r.Context())
	rw.Set("hunt", h)
	if err := rw.Render("admin/hunts/edit.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
