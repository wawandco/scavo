package auth

import (
	"net/http"

	"go.leapkit.dev/core/server/session"
)

// Logout clears the user session and redirects to the login page.
func Logout(w http.ResponseWriter, r *http.Request) {
	s := session.FromCtx(r.Context())
	delete(s.Values, "user_id")
	delete(s.Values, "is_admin")
	_ = s.Save(r, w)

	http.Redirect(w, r, "/auth/users/login", http.StatusSeeOther)
}
