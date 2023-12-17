package mocks

import "github.com/JakubBizewski/jakubme_links/domain/model"

type MockShortLinkRepository struct {
	StoreFunc           func(shortLink model.ShortLink) error
	FindByShortCodeFunc func(shortCode string) (model.ShortLink, error)
}

func (m *MockShortLinkRepository) Store(shortLink model.ShortLink) error {
	return m.StoreFunc(shortLink)
}

func (m *MockShortLinkRepository) FindByShortCode(shortCode string) (model.ShortLink, error) {
	return m.FindByShortCodeFunc(shortCode)
}
