// Package db provides small context helpers for the application database connection.
package db

import (
	"context"
	"database/sql"
)

// FromCtx extracts the *sql.DB injected by server.InCtxMiddleware("db", conn).
// It panics on programmer error (middleware not mounted). This centralizes the
// previously repeated "db, _ :=" pattern and makes the contract explicit.
func FromCtx(ctx context.Context) *sql.DB {
	database, ok := ctx.Value("db").(*sql.DB)
	if !ok || database == nil {
		panic("database connection not found in request context (missing InCtxMiddleware)")
	}
	return database
}
