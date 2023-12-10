package infrastructure

import "github.com/JakubBizewski/jakubme_links/shortLinks/domain"

type MemoryDomainRepository struct {
	domain *domain.Domain
}

func NewMemoryDomainRepository() *MemoryDomainRepository {
	domain, err := domain.CreateDomain("dummy.net", CreateMemoryLinksCollection())
	if err != nil {
		panic(err)
	}

	return &MemoryDomainRepository{
		domain: domain,
	}
}

func (memoryDomainRepository *MemoryDomainRepository) GetDomain() (*domain.Domain, error) {
	return memoryDomainRepository.domain, nil
}

func (memoryDomainRepository *MemoryDomainRepository) SaveDomain(domain *domain.Domain) error {
	memoryDomainRepository.domain = domain
	return nil
}
