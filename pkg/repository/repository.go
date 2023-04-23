package repository

import (
	"github.com/jmoiron/sqlx"
	inMemory "shortenLink/pkg/repository/in_memory"
	"shortenLink/pkg/repository/postgres"
)

type Link interface {
	CreateShortUrl(originalUrl, shortUrl, date string) (string, error)
	GetShortUrl(url string) (string, error)
	Delete(date string) error
}

type Repository struct {
	Link
}

func NewDB(db *sqlx.DB) *Repository {
	return &Repository{
		Link: postgres.New(db),
	}
}

func NewIm() *Repository {
	return &Repository{
		Link: inMemory.New(),
	}
}
