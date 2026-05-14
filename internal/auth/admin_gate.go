package auth

import (
	"net/http"

	"go.leapkit.dev/core/server"
	"go.leapkit.dev/core/server/session"
)

// AdminGate blocks non-admin users from accessing /admin routes.
// It returns a 404 Not Found to avoid leaking admin route existence.
func AdminGate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.FromCtx(r.Context())
		if s.Values["is_admin"] != true {
			server.Errorf(w, http.StatusNotFound, "not found")
			return
		}

		next.ServeHTTP(w, r)
	})
}
