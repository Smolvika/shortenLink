package service

import (
	"errors"
	"math/rand"
	"time"
)

const set = "_abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (s Service) CreateShortUrl(originalUrl, shortUrl, date string) (string, error) {
	for {
		shortUrl = GenerateShortUrl()
		url, err := s.repos.CreateShortUrl(originalUrl, shortUrl, date)
		if err != nil {
			if err != errors.New("there is already such a key") {
				return "", err
			}
		} else {
			return url, nil
		}
	}
}

func (s Service) GetShortUrl(url string) (string, error) {
	originalURL, err := s.repos.GetShortUrl(url)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (s Service) Delete(date string) error {
	return s.repos.Delete(date)
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
