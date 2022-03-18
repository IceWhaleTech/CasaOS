package web

import "embed"

//go:embed index.html favicon.svg img js browserconfig.xml site.webmanifest robots.txt
var Static embed.FS
