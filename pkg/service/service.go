package service

import "shortenLink/pkg/repository"

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Link interface {
	CreateShortUrl(originalUrl, date string) (string, error)
	GetShortUrl(url string) (string, error)
	Delete(date string) error
}
type Service struct {
	Link
}

func New(repos *repository.Repository) *Service {
	return &Service{Link: NewLinkService(repos)}
}
