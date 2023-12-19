package sqlite_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/JakubBizewski/jakubme_links/adapters/sqlite"
)

const dbPath = "../../storage/test.db"

func createDummyLinksTableWithData() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id_dummy INTEGER PRIMARY KEY AUTOINCREMENT
		);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO links (id_dummy)
		VALUES (1);`)
	if err != nil {
		return err
	}

	return nil
}

func linksTableExists() (bool, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return false, err
	}

	defer db.Close()

	row := db.QueryRow(`
		SELECT name
		FROM sqlite_master
		WHERE type='table' AND name='links';`,
	)

	var tableName string
	err = row.Scan(&tableName)
	if err != nil {
		return false, err
	}

	return tableName == "links", nil
}

func dummyTableExistsWithData() (bool, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return false, err
	}

	defer db.Close()

	row := db.QueryRow(`
		SELECT id_dummy
		FROM links`)

	var dummyID int
	err = row.Scan(&dummyID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func TestSqliteDbSetup(t *testing.T) {
	t.Run("ShouldCreateDatabaseFileWithLinksTable", func(t *testing.T) {
		err := sqlite.Setup(dbPath)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		_, err = os.Stat(dbPath)
		if os.IsNotExist(err) {
			t.Errorf("Expected database file to exist, but it doesn't")
		}

		tableExists, err := linksTableExists()
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		if !tableExists {
			t.Errorf("Expected links table to exist, but it doesn't")
		}
	})

	os.Remove(dbPath)

	t.Run("ShouldNotCreateDatabaseFileIfItAlreadyExists", func(t *testing.T) {
		err := createDummyLinksTableWithData()
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		err = sqlite.Setup(dbPath)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		_, err = os.Stat(dbPath)
		if os.IsNotExist(err) {
			t.Errorf("Expected database file to exist, but it doesn't")
		}

		dummyTableExists, err := dummyTableExistsWithData()
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		if !dummyTableExists {
			t.Errorf("Expected dummy table to exist, but it doesn't")
		}
	})

	os.Remove(dbPath)
}
