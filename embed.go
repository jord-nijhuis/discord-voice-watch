package projectroot

import "embed"

//go:embed migrations/*.sql
var MigrationFiles embed.FS

var Version = "0.0.1"
