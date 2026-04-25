package assets

import "embed"

//go:embed dist/*
var DistFS embed.FS
