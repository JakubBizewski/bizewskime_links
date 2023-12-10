package application

import "github.com/JakubBizewski/jakubme_links/shortLinks/domain"

type ShortenedLinksService struct {
	domainRepository domain.DomainRepository
}

func NewShortenedLinksService(domainRepository domain.DomainRepository) *ShortenedLinksService {
	return &ShortenedLinksService{
		domainRepository: domainRepository,
	}
}

func (shortenedLinksService *ShortenedLinksService) CreateShortenedLink(targetUrl domain.Url) (*domain.Link, error) {
	domain, err := shortenedLinksService.domainRepository.GetDomain()
	if err != nil {
		return nil, err
	}

	link, err := domain.CreateShortLink(targetUrl)
	if err != nil {
		return nil, err
	}

	err = shortenedLinksService.domainRepository.SaveDomain(domain)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (shortenedLinksService *ShortenedLinksService) GetLinkByShortCode(shortCode domain.ShortCode) (string, error) {
	domain, err := shortenedLinksService.domainRepository.GetDomain()
	if err != nil {
		return "", err
	}

	link, err := domain.GetLinkByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	if link == nil {
		return "", nil
	}

	return string(link.TargetUrl), nil
}
