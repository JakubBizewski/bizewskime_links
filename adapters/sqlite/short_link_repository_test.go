package sqlite_test

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/JakubBizewski/jakubme_links/adapters/sqlite"
	"github.com/JakubBizewski/jakubme_links/domain/model"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driven"
)

func seedExistingShortLink(shortLink model.ShortLink) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO links (short_code, target_url)
		VALUES (?, ?);`,
		shortLink.ShortCode,
		shortLink.TargetURL,
	)
	if err != nil {
		return err
	}

	return nil
}

func getShortLink(shortCode string) (model.ShortLink, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return model.ShortLink{}, err
	}

	defer db.Close()

	row := db.QueryRow(`
		SELECT short_code, target_url
		FROM links
		WHERE short_code = ?;`,
		shortCode,
	)

	var shortLink model.ShortLink
	err = row.Scan(&shortLink.ShortCode, &shortLink.TargetURL)
	if err != nil {
		return model.ShortLink{}, err
	}

	return shortLink, nil
}

func TestSqliteShortLinkRepository(t *testing.T) {
	dbSetupErr := sqlite.Setup(dbPath)
	if dbSetupErr != nil {
		t.Error(dbSetupErr)
	}

	defer os.Remove(dbPath)

	existingShortLink := model.ShortLink{
		ShortCode: "existingShortCode",
		TargetURL: "https://example.com",
	}
	seedRrr := seedExistingShortLink(existingShortLink)
	if seedRrr != nil {
		t.Error(seedRrr)
	}

	repository := sqlite.CreateShortLinkRepository(dbPath)

	t.Run("ShouldStore", func(t *testing.T) {
		shortLink := model.ShortLink{
			ShortCode: "testShortCode",
			TargetURL: "https://example.com",
		}

		err := repository.Store(shortLink)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		storedShortLink, err := getShortLink(shortLink.ShortCode)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		if storedShortLink.ShortCode != shortLink.ShortCode {
			t.Errorf("Expected short code %s, but got %s", shortLink.ShortCode, storedShortLink.ShortCode)
		}
	})

	t.Run("ShouldReturnErrShortCodeAlreadyExists", func(t *testing.T) {
		err := repository.Store(existingShortLink)
		if !errors.Is(err, driven.ErrShortCodeAlreadyExists) {
			t.Errorf("Expected error %s, but got %s", driven.ErrShortCodeAlreadyExists, err)
		}
	})

	t.Run("ShouldFindByShortCode", func(t *testing.T) {
		storedShortLink, err := repository.FindByShortCode(existingShortLink.ShortCode)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		if storedShortLink.ShortCode != existingShortLink.ShortCode {
			t.Errorf("Expected short code %s, but got %s", existingShortLink.ShortCode, storedShortLink.ShortCode)
		}
	})

	t.Run("ShouldReturnEmptyShortLinkIfNotFound", func(t *testing.T) {
		link, err := repository.FindByShortCode("nonExistingShortCode")
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		if link != (model.ShortLink{}) {
			t.Errorf("Expected empty short link, but got %v", link)
		}
	})
}
