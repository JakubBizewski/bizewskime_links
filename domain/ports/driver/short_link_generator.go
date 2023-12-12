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

func CreateShortLinkService(shortLinkRepository driven.ShortLinkRepository) *ShortLinkService {
	return &ShortLinkService{
		shortLinkRepository: shortLinkRepository,
	}
}

func (s *ShortLinkService) GenerateShortLink(targetUrl string) (string, error) {
	for i := 0; i < shortCodeGenerationMaxAttempts; i++ {
		shortLink := model.CreateRandomShortLink(targetUrl, defaultShortCodeLen)

		err := s.shortLinkRepository.Store(shortLink)
		if err == nil {
			return shortLink.ShortCode, nil
		}

		if err != driven.ErrShortCodeAlreadyExists {
			return "", err
		}
	}

	return "", ErrShortCodeGenerationFailed
}

func (s *ShortLinkService) GetTargetUrl(shortCode string) (string, error) {
	shortLink, err := s.shortLinkRepository.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	return shortLink.TargetUrl, nil
}
