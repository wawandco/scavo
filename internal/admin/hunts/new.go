package hunts

import (
	"net/http"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// New renders the form to create a new hunt.
func New(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())
	if err := rw.Render("admin/hunts/new.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
