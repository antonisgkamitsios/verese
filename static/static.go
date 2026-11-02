package static

import "embed"

//go:embed "index.js" "style.css"
var StaticFiles embed.FS
