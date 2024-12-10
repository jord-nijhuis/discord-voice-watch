package projectroot

import "embed"

//go:embed migrations/*.sql
var MigrationFiles embed.FS

// Version is the version of the application. This is set at build time to the git tag
// using the -ldflags flag.
var Version = "0.0.0"
