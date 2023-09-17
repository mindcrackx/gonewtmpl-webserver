package ui

import "embed"

//go:embed static
var EmbeddedContentStatic embed.FS

//go:embed html
var EmbeddedContentHTML embed.FS
