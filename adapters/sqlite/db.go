package sqlite

import (
	"database/sql"
)

const providerName = "sqlite3"

func Setup(dbPath string) error {
	db, err := sql.Open(providerName, dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			short_code TEXT NOT NULL UNIQUE,
			target_url TEXT NOT NULL
		);`)

	if err != nil {
		return err
	}

	return nil
}
