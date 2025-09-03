package sql

import "embed"

//go:embed script/*.sql
var SQLScript embed.FS
