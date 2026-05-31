// Package assets embeds the static frontend assets (CSS, JS) served by the application.
package assets

import (
	"embed"
)

//go:embed *
var Files embed.FS
