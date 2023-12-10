package domain

type LinksCollection interface {
	AddLink(link *Link) error
	GetLinkByShortCode(shortCode ShortCode) (*Link, error)
	ShortCodeExists(shortCode ShortCode) bool
}
