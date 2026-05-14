package admin

import (
	"database/sql"
	"time"
)

// FetchTeams returns all teams for dropdown selection.
func FetchTeams(db *sql.DB) ([]Team, error) {
	rows, err := db.Query("SELECT id, name FROM teams ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

// Hunt represents a scavenger hunt.
type Hunt struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Points      int       `db:"points"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Team represents a scavenger hunt team.
type Team struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Member represents a scavenger hunt participant.
type Member struct {
	ID         int           `db:"id"`
	Name       string        `db:"name"`
	PersonalID string        `db:"personal_id"`
	TeamID     sql.NullInt64 `db:"team_id"`
	IsAdmin    bool          `db:"is_admin"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`

	TeamName string `db:"-"`
}

// Submission represents a hunt entry submitted by a team member.
type Submission struct {
	ID         int       `db:"id"`
	HuntID     int       `db:"hunt_id"`
	MemberID   int       `db:"member_id"`
	TeamID     int       `db:"team_id"`
	Text       string    `db:"text"`
	ImagePath  string    `db:"image_path"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`

	MemberName string `db:"-"`
	TeamName   string `db:"-"`
	HuntTitle  string `db:"-"`
}
