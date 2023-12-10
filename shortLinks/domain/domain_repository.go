package domain

type DomainRepository interface {
	GetDomain() (*Domain, error)
	SaveDomain(domain *Domain) error
}
