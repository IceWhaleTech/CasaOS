package web

import "embed"

//go:embed index.html favicon.svg img js browserconfig.xml site.webmanifest
var Static embed.FS
