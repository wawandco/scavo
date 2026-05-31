package auth

import (
	"database/sql"
	"net/http"

	"scavo/internal/db"

	"go.leapkit.dev/core/server"
	"go.leapkit.dev/core/server/session"
)

// Login authenticates a member by personal_id and stores their ID in the session.
func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		server.Errorf(w, http.StatusBadRequest, "form parse error: %s", err.Error())
		return
	}

	personalID := r.FormValue("personal_id")
	if personalID == "" {
		server.Errorf(w, http.StatusBadRequest, "personal ID is required")
		return
	}

	db := db.FromCtx(r.Context())
	var memberID, isAdmin int
	err := db.QueryRow("SELECT id, is_admin FROM members WHERE personal_id = ?", personalID).Scan(&memberID, &isAdmin)
	if err == sql.ErrNoRows {
		server.Errorf(w, http.StatusUnauthorized, "invalid personal ID")
		return
	}
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}

	s := session.FromCtx(r.Context())
	s.Values["user_id"] = memberID
	s.Values["is_admin"] = isAdmin == 1
	_ = s.Save(r, w)

	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = "/"
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
