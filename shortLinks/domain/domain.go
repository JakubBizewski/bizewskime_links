package domain

import (
	"errors"

	"github.com/google/uuid"
)

const shortCodeLength = 4
const shortCodeGenerationAttempts = 5

var FailedToGenerateUniqueShortCodeError = errors.New("Failed to generate unique short code")

type Domain struct {
	ID    uuid.UUID
	Name  string
	Links LinksCollection
}

func CreateDomain(name string, linksCollection LinksCollection) (*Domain, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Domain{
		ID:    id,
		Name:  name,
		Links: linksCollection,
	}, nil
}

func (domain *Domain) CreateShortLink(targetUrl Url) (*Link, error) {
	shortCode, err := domain.getNewUniqueShortCode()
	if err != nil {
		return nil, err
	}

	link, err := CreateLink(targetUrl, shortCode)
	if err != nil {
		return nil, err
	}

	err = domain.Links.AddLink(link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (domain *Domain) GetLinkByShortCode(shortCode ShortCode) (*Link, error) {
	link, err := domain.Links.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (domain *Domain) getNewUniqueShortCode() (ShortCode, error) {
	generatedShortCode := CreateRandomShortCode(shortCodeLength)

	for i := 0; i < shortCodeGenerationAttempts; i++ {
		if !domain.Links.ShortCodeExists(generatedShortCode) {
			return generatedShortCode, nil

		}

		generatedShortCode = CreateRandomShortCode(shortCodeLength)
	}

	return "", FailedToGenerateUniqueShortCodeError
}
