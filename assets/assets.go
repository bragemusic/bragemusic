package assets

import "embed"

//go:embed frontend/*
var DistFS embed.FS
