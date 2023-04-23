package service

import (
	"errors"
	"log"
	"math/rand"
	"shortenLink/pkg/repository"
	"time"
)

const set = "_abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type LinkService struct {
	repo repository.Link
}

func NewLinkService(repo repository.Link) *LinkService {
	return &LinkService{repo: repo}
}
func (s *LinkService) CreateShortUrl(originalUrl, date string) (string, error) {
	for i := 0; i < 10; i++ {
		shortUrl := GenerateShortUrl()
		url, err := s.repo.CreateShortUrl(originalUrl, shortUrl, date)
		if err != nil {
			return "", err
		} else if url != "" {
			return url, nil
		}
	}
	log.Println("exceeded the number of attempts")
	return "", errors.New("there is already such a key")
}

func (s *LinkService) GetShortUrl(url string) (string, error) {
	return s.repo.GetShortUrl(url)
}

func (s *LinkService) Delete(date string) error {
	return s.repo.Delete(date)
}

func GenerateShortUrl() string {
	url := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	for i := range url {
		inx := rand.Intn(len(set))
		url[i] = set[inx]
	}
	return string(url)
}
