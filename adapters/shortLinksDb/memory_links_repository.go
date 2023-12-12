package shortLinksDb

import "github.com/JakubBizewski/jakubme_links/domain/model"

type MemoryShortLinkRepository struct {
	shortLinks map[string]*model.ShortLink
}

func NewMemoryShortLinkRepository() *MemoryShortLinkRepository {
	return &MemoryShortLinkRepository{
		shortLinks: make(map[string]*model.ShortLink),
	}
}

func (r *MemoryShortLinkRepository) Store(shortLink *model.ShortLink) error {
	r.shortLinks[shortLink.ShortCode] = shortLink

	return nil
}

func (r *MemoryShortLinkRepository) FindByShortCode(shortCode string) (*model.ShortLink, error) {
	shortLink, ok := r.shortLinks[shortCode]
	if !ok {
		return nil, nil
	}

	return shortLink, nil
}

func (r *MemoryShortLinkRepository) ShortCodeExists(shortCode string) (bool, error) {
	_, ok := r.shortLinks[shortCode]

	return ok, nil
}
