package auth

import (
	"net/http"

	"go.leapkit.dev/core/server/session"
)

// RequireAuth redirects unauthenticated users to the login page.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.FromCtx(r.Context())
		if s.Values["user_id"] == nil {
			http.Redirect(w, r, "/auth/users/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
