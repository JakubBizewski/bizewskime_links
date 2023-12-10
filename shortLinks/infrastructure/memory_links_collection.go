package infrastructure

import "github.com/JakubBizewski/jakubme_links/shortLinks/domain"

type MemoryLinksCollection struct {
	links map[domain.ShortCode]*domain.Link
}

func CreateMemoryLinksCollection() *MemoryLinksCollection {
	return &MemoryLinksCollection{
		links: make(map[domain.ShortCode]*domain.Link),
	}
}

func (memoryLinksCollection *MemoryLinksCollection) AddLink(link *domain.Link) error {
	memoryLinksCollection.links[link.ShortCode] = link
	return nil
}

func (memoryLinksCollection *MemoryLinksCollection) GetLinkByShortCode(shortCode domain.ShortCode) (*domain.Link, error) {
	return memoryLinksCollection.links[shortCode], nil
}

func (memoryLinksCollection *MemoryLinksCollection) ShortCodeExists(shortCode domain.ShortCode) bool {
	_, ok := memoryLinksCollection.links[shortCode]
	return ok
}
