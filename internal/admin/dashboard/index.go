// Package dashboard provides the main admin dashboard handler and overview metrics.
package dashboard

import (
	"database/sql"
	"net/http"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// Index renders the admin landing page.
func Index(w http.ResponseWriter, r *http.Request) {
	db, _ := r.Context().Value("db").(*sql.DB)

	var pendingCount int
	_ = db.QueryRow("SELECT COUNT(*) FROM submissions WHERE status = 'pending'").Scan(&pendingCount)

	rw := render.FromCtx(r.Context())
	rw.Set("pendingCount", pendingCount)
	if err := rw.Render("admin/dashboard/index.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
