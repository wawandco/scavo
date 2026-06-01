// Package assets embeds the static frontend assets (CSS, JS) served by the application.
package assets

import (
	"embed"
)

//go:embed *

// Files contains the embedded static frontend assets (CSS, JS, images, etc.).
var Files embed.FS
