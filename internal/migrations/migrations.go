// Package migrations contains the embedded SQL migration files and runner for Scavo.
package migrations

import "embed"

//go:embed *.sql

// All contains the embedded SQL migration files.
var All embed.FS
