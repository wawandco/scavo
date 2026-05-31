// Package migrations contains the embedded SQL migration files and runner for Scavo.
package migrations

import "embed"

//go:embed *.sql
var All embed.FS
