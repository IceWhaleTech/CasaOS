package web

import "embed"

//go:embed index.html favicon.svg browserconfig.xml site.webmanifest robots.txt img js fonts css
var Static embed.FS
