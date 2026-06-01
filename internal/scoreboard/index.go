// Package scoreboard renders the public real-time leaderboard
// showing team progress and points during an event.
package scoreboard

import (
	"net/http"

	"scavo/internal/db"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// TeamScore holds a team's rank data.
type TeamScore struct {
	Rank        int
	Name        string
	TotalPoints int
}

// Index renders the public scoreboard.
func Index(w http.ResponseWriter, r *http.Request) {
	db := db.FromCtx(r.Context())

	rows, err := db.Query(`
		SELECT t.name, COALESCE(SUM(h.points), 0)
		FROM teams t
		LEFT JOIN submissions s ON t.id = s.team_id AND s.status = 'approved'
		LEFT JOIN hunts h ON s.hunt_id = h.id
		GROUP BY t.id
		ORDER BY COALESCE(SUM(h.points), 0) DESC`)
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "query error: %s", err.Error())
		return
	}
	defer rows.Close() //nolint:errcheck // best effort cleanup

	var scores []TeamScore
	position := 0
	previousPoints := -1
	for rows.Next() {
		var s TeamScore
		if err := rows.Scan(&s.Name, &s.TotalPoints); err != nil {
			server.Errorf(w, http.StatusInternalServerError, "scan error: %s", err.Error())
			return
		}
		if s.TotalPoints != previousPoints {
			position++
		}
		s.Rank = position
		previousPoints = s.TotalPoints
		scores = append(scores, s)
	}

	var possiblePoints int
	_ = db.QueryRow(`SELECT COALESCE(SUM(points), 0) FROM hunts`).Scan(&possiblePoints)

	rw := render.FromCtx(r.Context())
	rw.Set("scores", scores)
	rw.Set("possiblePoints", possiblePoints)
	if err := rw.Render("scoreboard/index.html"); err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
