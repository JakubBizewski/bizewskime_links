package driven

import (
	"errors"

	"github.com/JakubBizewski/jakubme_links/domain/model"
)

var ErrShortCodeAlreadyExists = errors.New("Short code already exists")

type ShortLinkRepository interface {
	Store(shortLink model.ShortLink) error
	FindByShortCode(shortCode string) (model.ShortLink, error)
}
