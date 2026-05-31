package home

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"scavo/internal/db"

	"go.leapkit.dev/core/server"
	"go.leapkit.dev/core/server/session"
)

// CreateSubmission handles a team member uploading a hunt submission.
func CreateSubmission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		server.Errorf(w, http.StatusBadRequest, "form parse error: %s", err.Error())
		return
	}

	huntID := r.FormValue("hunt_id")
	if huntID == "" {
		server.Errorf(w, http.StatusBadRequest, "hunt ID is required")
		return
	}

	s := session.FromCtx(r.Context())
	userID, ok := s.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/auth/users/login", http.StatusSeeOther)
		return
	}

	db := db.FromCtx(r.Context())
	var teamID sql.NullInt64
	_ = db.QueryRow("SELECT team_id FROM members WHERE id = ?", userID).Scan(&teamID)
	if !teamID.Valid {
		server.Errorf(w, http.StatusBadRequest, "you must be assigned to a team to submit")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		server.Errorf(w, http.StatusBadRequest, "image is required")
		return
	}
	defer file.Close() //nolint:errcheck // best effort close on request body part

	// Security hardening: turn the untrusted client filename into something safe to store.
	cleanBase := sanitizeUploadFilename(header.Filename)
	filename := fmt.Sprintf("%d_%d_%s", time.Now().Unix(), userID, cleanBase)

	// Ensure uploads dir exists with restricted perms
	if err := os.MkdirAll("uploads", 0o750); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "failed to prepare upload directory")
		return
	}

	out, err := os.Create("uploads/" + filename)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "file save error: %s", err.Error())
		return
	}
	defer out.Close() //nolint:errcheck // best effort close on file write
	_, _ = io.Copy(out, file)

	_, _ = db.Exec(
		"DELETE FROM submissions WHERE hunt_id = ? AND team_id = ?",
		huntID, teamID.Int64,
	)

	_, _ = db.Exec(
		"INSERT INTO submissions (hunt_id, member_id, team_id, text, image_path, status, updated_at) VALUES (?, ?, ?, ?, ?, 'pending', ?)",
		huntID, userID, teamID.Int64, r.FormValue("text"), filename, time.Now(),
	)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// sanitizeUploadFilename takes an untrusted filename provided by a client
// (via multipart form upload) and returns a safe base name suitable for
// writing to disk.
//
// It defends against:
//   - directory traversal (../foo.png, /etc/passwd, etc.)
//   - control characters and most special filesystem characters
//   - missing, extremely long, or dangerous extensions
//   - names that are too long for practical storage
//
// Only the base filename is returned — no directory components.
func sanitizeUploadFilename(original string) string {
	base := filepath.Base(original)
	if base == "" || base == "." || base == ".." {
		return "upload.bin"
	}

	// Keep only characters that are generally safe across filesystems and URLs.
	clean := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_', r == '.':
			return r
		default:
			return '-'
		}
	}, base)

	if clean == "" || clean == "." {
		return "upload.bin"
	}

	// Remove whatever extension the (already sanitized) name had.
	originalExt := filepath.Ext(clean)
	clean = strings.TrimSuffix(clean, originalExt)

	// Decide on the final safe extension.
	ext := strings.ToLower(originalExt)
	if ext == "" || len(ext) > 8 {
		ext = ".bin"
	}

	// Truncate stem if needed, then append the chosen extension.
	const maxStem = 76
	if len(clean) > maxStem {
		clean = clean[:maxStem]
	}
	clean = clean + ext

	return clean
}
