package domain

import (
	"github.com/google/uuid"
)

type Link struct {
	ID        uuid.UUID
	TargetUrl Url
	ShortCode ShortCode
}

func CreateLink(targetUrl Url, shortCode ShortCode) (*Link, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Link{
		ID:        id,
		TargetUrl: targetUrl,
		ShortCode: shortCode,
	}, nil
}
