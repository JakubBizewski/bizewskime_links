package driver

import (
	"errors"

	"github.com/JakubBizewski/jakubme_links/domain/model"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driven"
)

var ErrShortCodeGenerationFailed = errors.New("Failed to generate unique short code")

const defaultShortCodeLen = 3
const shortCodeGenerationMaxAttempts = 10

type ShortLinkService struct {
	shortLinkRepository driven.ShortLinkRepository
}

func NewShortLinkService(shortLinkRepository driven.ShortLinkRepository) *ShortLinkService {
	return &ShortLinkService{
		shortLinkRepository: shortLinkRepository,
	}
}

func (s *ShortLinkService) GenerateShortLink(targetUrl string) (string, error) {
	shortLink, getUniqueShortLinkError := s.getUniqueShortLink(targetUrl)
	if getUniqueShortLinkError != nil {
		return "", getUniqueShortLinkError
	}

	storeLinkError := s.shortLinkRepository.Store(shortLink)
	if storeLinkError != nil {
		return "", storeLinkError
	}

	return shortLink.ShortCode, nil
}

func (s *ShortLinkService) getUniqueShortLink(targetUrl string) (*model.ShortLink, error) {
	for attempt := 0; attempt < shortCodeGenerationMaxAttempts; attempt++ {
		shortLink := model.NewRandomShortLink(targetUrl, defaultShortCodeLen)

		shortCodeExists, err := s.shortLinkRepository.ShortCodeExists(shortLink.ShortCode)
		if err != nil {
			return nil, err
		}

		if !shortCodeExists {
			return &shortLink, nil
		}
	}

	return nil, ErrShortCodeGenerationFailed
}

func (s *ShortLinkService) GetTargetUrl(shortCode string) (string, error) {
	shortLink, err := s.shortLinkRepository.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	if shortLink == nil {
		return "", nil
	}

	return shortLink.TargetUrl, nil
}
