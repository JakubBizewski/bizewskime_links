package sqlite

import (
	"database/sql"
	"errors"
	"log"

	"github.com/JakubBizewski/jakubme_links/domain/model"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driven"
)

type ShortLinkRepository struct {
	dbPath string
}

func CreateShortLinkRepository(dbPath string) *ShortLinkRepository {
	return &ShortLinkRepository{dbPath: dbPath}
}

func (r *ShortLinkRepository) Store(shortLink model.ShortLink) error {
	db, err := sql.Open(providerName, r.dbPath)
	if err != nil {
		log.Print(err)

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
		if err.Error() == "UNIQUE constraint failed: links.short_code" {
			return driven.ErrShortCodeAlreadyExists
		}

		log.Print(err)
		return err
	}

	return nil
}

func (r *ShortLinkRepository) FindByShortCode(shortCode string) (model.ShortLink, error) {
	db, err := sql.Open(providerName, r.dbPath)
	if err != nil {
		log.Print(err)

		return model.ShortLink{}, err
	}

	defer db.Close()

	row := db.QueryRow(`
		SELECT short_code, target_url
		FROM links
		WHERE short_code = ?;`,
		shortCode,
	)

	var shortCodeResult string
	var targetURL string

	err = row.Scan(&shortCodeResult, &targetURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ShortLink{}, nil
		}

		log.Print(err)
		return model.ShortLink{}, err
	}

	return model.ShortLink{
		ShortCode: shortCodeResult,
		TargetURL: targetURL,
	}, nil
}
