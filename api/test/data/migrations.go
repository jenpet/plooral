package data

import (
	"embed"
)

//go:embed migrations/*.sql
var TestDataMigrations embed.FS