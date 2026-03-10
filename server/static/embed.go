package static

import "embed"

//go:embed js/*
//go:embed css/*
//go:embed *
var FS embed.FS
