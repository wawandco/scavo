package auth

import (
	"net/http"
	"strings"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server/session"
)

// AuthContext sets isAdmin, isLoggedIn, and isAdminUser values in the
// render context so templates can conditionally show navigation.
func AuthContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := render.FromCtx(r.Context())
		rw.Set("isAdmin", strings.HasPrefix(r.URL.Path, "/admin"))

		s := session.FromCtx(r.Context())
		rw.Set("isLoggedIn", s.Values["user_id"] != nil)
		rw.Set("isAdminUser", s.Values["is_admin"] == true)
		next.ServeHTTP(w, r)
	})
}
