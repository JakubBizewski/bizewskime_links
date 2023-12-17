package memory

import (
	"sync"

	"github.com/JakubBizewski/jakubme_links/domain/model"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driven"
)

type ShortLinkRepository struct {
	shortLinks map[string]model.ShortLink
	mutex      sync.Mutex
}

func CreateMemoryShortLinkRepository() *ShortLinkRepository {
	return &ShortLinkRepository{
		shortLinks: make(map[string]model.ShortLink),
		mutex:      sync.Mutex{},
	}
}

func (r *ShortLinkRepository) Store(shortLink model.ShortLink) error {
	r.mutex.Lock()
	if _, exists := r.shortLinks[shortLink.ShortCode]; exists {
		r.mutex.Unlock()
		return driven.ErrShortCodeAlreadyExists
	}

	r.shortLinks[shortLink.ShortCode] = shortLink
	r.mutex.Unlock()

	return nil
}

func (r *ShortLinkRepository) FindByShortCode(shortCode string) (model.ShortLink, error) {
	shortLink, ok := r.shortLinks[shortCode]
	if !ok {
		return model.ShortLink{}, nil
	}

	return shortLink, nil
}
