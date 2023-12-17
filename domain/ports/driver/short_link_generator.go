package driver

import (
	"errors"

	"github.com/JakubBizewski/jakubme_links/domain/model"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driven"
)

var ErrShortCodeGenerationFailed = errors.New("failed to generate unique short code")

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

func (s *ShortLinkService) GenerateShortLink(targetURL string) (string, error) {
	for i := 0; i < shortCodeGenerationMaxAttempts; i++ {
		shortLink := model.CreateRandomShortLink(targetURL, defaultShortCodeLen)

		err := s.shortLinkRepository.Store(shortLink)
		if err == nil {
			return shortLink.ShortCode, nil
		}

		if !errors.Is(err, driven.ErrShortCodeAlreadyExists) {
			return "", err
		}
	}

	return "", ErrShortCodeGenerationFailed
}

func (s *ShortLinkService) GetTargetURL(shortCode string) (string, error) {
	shortLink, err := s.shortLinkRepository.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	return shortLink.TargetURL, nil
}
