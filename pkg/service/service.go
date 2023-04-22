package service

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Repository interface {
	CreateShortUrl(originalUrl, shortUrl, date string) (string, error)
	GetShortUrl(url string) (string, error)
	Delete(date string) error
}
type Service struct {
	repos Repository
}

func New(repos Repository) *Service {
	return &Service{repos: repos}
}
