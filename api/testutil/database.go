package testutil

import (
	"database/sql"
	"embed"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"path/filepath"
)

func ApplyFromEmbeddedFS(pgURI string, fsys embed.FS, path string) {
	db, err := sql.Open("postgres", pgURI)
	if err != nil {
		log.Fatalf("Failed to connect to postgres database. Error: %v", err)
	}
	files, err := fs.ReadDir(fsys, path)
	if err != nil {
		log.Fatalf("Failed to read embedded filesystem migrations. Error: %v", err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		b, err := fsys.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			log.Fatalf("Failed to read embedded file '%s'. Error: %v", err, f.Name())
		}
		ApplyTestMigration(db, b)
	}
}

func ApplyTestMigration(db *sql.DB, b []byte) {
	_, err := db.Query(string(b))
	if err != nil {
		log.Fatalf("Failed to apply migration file '%s'. Error: %v", err)
	}
}