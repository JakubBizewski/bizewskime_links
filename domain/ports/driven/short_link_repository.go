package driven

import "github.com/JakubBizewski/jakubme_links/domain/model"

type ShortLinkRepository interface {
	Store(shortLink *model.ShortLink) error
	FindByShortCode(shortCode string) (*model.ShortLink, error)
	ShortCodeExists(shortCode string) (bool, error)
}
