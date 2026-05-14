package auth

import (
	"net/http"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// ShowLogin renders the member login page.
func ShowLogin(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())
	rw.Set("redirect", r.URL.Query().Get("redirect"))
	if err := rw.Render("auth/login.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
