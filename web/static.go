package web

import "embed"

//go:embed index.html favicon.ico img js browserconfig.xml site.webmanifest
var Static embed.FS
