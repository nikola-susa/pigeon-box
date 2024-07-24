package assets

import (
	"embed"
	"io/fs"
)

//go:embed static
var _publicFS embed.FS

var PublicFS, _ = fs.Sub(_publicFS, "static")
